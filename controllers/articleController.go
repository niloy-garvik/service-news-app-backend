package controller

import (
	"context"
	"encoding/json"
	"net/http"
	PostgresInstance "service-news-app-backend/Postgres_Instance"
	schemas "service-news-app-backend/Schemas"
	"service-news-app-backend/utils"
)

var ctx = context.Background()

func ExtractMetaDataHandler(w http.ResponseWriter, r *http.Request) {

	var body schemas.ExtractMetaDataHandlerBody

	// decode body
	json.NewDecoder(r.Body).Decode(&body)

	// body validation
	validationError := schemas.ValidateInput(body)
	if validationError != nil {
		utils.SendErrorResponse(w, http.StatusBadRequest, "bodyValidationFailed", validationError.Error(), nil)
		return
	}

	// Call OpenAI API to get categories
	responseFromOpenAI, err := utils.GetResponseFromChatGPT(ctx, body.Content)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "openAIError", err.Error(), nil)
		return
	}

	// generating summary
	summary, err := utils.GenerateSummary(ctx, body.Content)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "openAIError", err.Error(), nil)
		return
	}
	publicationDate, _ := utils.ConvertStringToTimestamp(body.PublicationDate)

	// building article info struct
	articleInfoObj := schemas.ArticleSchema{
		ArticleId:       body.ArticleId,
		Title:           body.Title,
		Publisher:       body.Publisher,
		PublicationDate: publicationDate,
		Url:             body.Url,
		Content:         body.Content,
		Summary:         summary,
		Tags:            body.Tags,
		Entities:        responseFromOpenAI.Entities,
		SentimentScore:  responseFromOpenAI.SentimentScore,
		Categories:      responseFromOpenAI.Categories,
		ContentS3Path:   body.ContentS3Path,
		Status:          "published",
	}

	// get article info
	articleInfo, err := schemas.GetArticleByID(ctx, PostgresInstance.GetPostgresInstance(), body.ArticleId)
	if err != nil {
		utils.SendErrorResponse(w, http.StatusInternalServerError, "internalServerError", err.Error(), nil)
		return
	}

	if articleInfo == nil {
		// store article info
		err = schemas.InsertArticleData(ctx, PostgresInstance.GetPostgresInstance(), articleInfoObj)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "internalServerError", err.Error(), nil)
			return
		}

		utils.SendSuccessResponse(w, http.StatusOK, "Aritcle Created successfully", nil)

	} else {
		// update article info
		err = schemas.UpdateArticleByID(ctx, PostgresInstance.GetPostgresInstance(), articleInfoObj)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, "internalServerError", err.Error(), nil)
			return
		}

		utils.SendSuccessResponse(w, http.StatusOK, "Aritcle Updated successfully", nil)
		return
	}

	// Create the prompt for OpenAI
	// prompt := "Categorize the following news article into relevant topics separated by commas:\n\n" + body.Content

	// // Call OpenAI API to get categories
	// _, err := utils.GenerateEmbeddings(ctx, prompt)
	// if err != nil {
	// 	log.Printf("Error getting categories from OpenAI: %v", err)
	// }

}
