package schemas

import "github.com/go-playground/validator"

func ValidateInput(input interface{}) error {
	validate := validator.New()
	err := validate.Struct(input)
	return err
}

type ExtractMetaDataHandlerBody *struct {
	ArticleId       string   `validate:"required" json:"articleId,omitempty" bson:"articleId,omitempty"`
	Title           string   `validate:"required" json:"title,omitempty" bson:"title,omitempty"`
	Publisher       string   `validate:"required" json:"publisher,omitempty" bson:"publisher,omitempty"`
	PublicationDate string   `validate:"required" json:"publicationDate,omitempty" bson:"publicationDate,omitempty"`
	Url             string   `validate:"required" json:"url,omitempty" bson:"url,omitempty"`
	Content         string   `validate:"required" json:"content,omitempty" bson:"content,omitempty"`
	Summary         string   `validate:"required" json:"summary,omitempty" bson:"summary,omitempty"`
	Tags            []string `validate:"required" json:"tags,omitempty" bson:"tags,omitempty"`
	ContentS3Path   string   `validate:"required" json:"contentS3Path,omitempty" bson:"contentS3Path,omitempty"`
}
