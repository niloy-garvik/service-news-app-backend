package controller

// import (
// 	"service-news-app-backend/utils"
// 	"fmt"
// 	"io"
// 	"net/http"

// 	"github.com/mmcdole/gofeed"
// )

// func FetchNewsController(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	fp := gofeed.NewParser()
// 	feed, err := fp.ParseURL("https://www.cnbc.com/id/100003114/device/rss/rss.html")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Feed Title:", feed.Title)

// 	// news, _ := FetchNewsDetails(feed.Items[0].Link)
// 	// fmt.Println("\n\n News: ", string(utils.ConvertToJson(item.Link)))
// 	// for _, item := range feed.Items {
// 	// 	news, _ := FetchNewsDetails(item.Link)
// 	// 	fmt.Println("\n\n News: ", string(utils.ConvertToJson(item.Link)))
// 	// 	fmt.Println("news: ", news)
// 	// }

// 	successResponse := utils.GenerateSuccessResponse("Data added successfully", feed.Items)
// 	apiResponse := utils.ConvertToJson(successResponse)

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(apiResponse))
// 	return

// }

// func FetchNews() {
// 	fp := gofeed.NewParser()
// 	feed, err := fp.ParseURL("https://www.cnbc.com/id/100003114/device/rss/rss.html")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("Feed Title:", feed.Title)

// 	news, _ := FetchNewsDetails(feed.Items[0].Link)
// 	// fmt.Println("\n\n News: ", string(utils.ConvertToJson(item.Link)))
// 	fmt.Println("news: ", news)
// 	// for _, item := range feed.Items {
// 	// 	news, _ := FetchNewsDetails(item.Link)
// 	// 	fmt.Println("\n\n News: ", string(utils.ConvertToJson(item.Link)))
// 	// 	fmt.Println("news: ", news)
// 	// }
// }

// // FetchNewsDetails fetches and returns the content of a news article from its URL.
// func FetchNewsDetails(url string) (string, error) {
// 	// Make an HTTP GET request to the news article's URL
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return "", err // Return error if the request fails
// 	}
// 	defer resp.Body.Close() // Ensure the response body is closed after reading

// 	// Check if the response status is OK (200)
// 	if resp.StatusCode != http.StatusOK {
// 		return "", fmt.Errorf("failed to fetch news details: %s", resp.Status)
// 	}

// 	// Read the response body
// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err // Return error if reading fails
// 	}

// 	return string(body), nil // Return the content as a string
// }
