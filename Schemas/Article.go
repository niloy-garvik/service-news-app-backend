package schemas

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type ArticleSchema struct {
	ArticleId       string   `json:"articleId"` // UUID as string
	Title           string   `json:"title"`
	Publisher       string   `json:"publisher"`
	PublicationDate string   `json:"publicationDate"` // TIMESTAMP
	Url             string   `json:"url"`
	Content         string   `json:"content"`
	Summary         string   `json:"summary"`
	Tags            []string `json:"tags"`           // -- e.g., ["Indian Army", "Kashmir"]
	Entities        any      `json:"entities"`       // -- e.g., {"organizations": ["Indian Army"], "locations": ["Kashmir"]}
	SentimentScore  string   `json:"sentimentScore"` // -- e.g., "Positive", "Negative", "Neutral"
	Categories      []string `json:"categories"`     // -- e.g., ["National Security", "Conflict"]
	ContentS3Path   string   `json:"contentS3Path"`  // -- e.g., S3 path to the full article
	Status          string   `json:"status"`         // -- e.g., "published" or "unpublished"
	CreatedAt       string   `json:"createdAt"`      // TIMESTAMP
	UpdatedAt       string   `json:"updatedAt"`      // TIMESTAMP
}

func CreateArticlesTable(ctx context.Context, pool *pgxpool.Pool) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS articles (
		article_id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		publisher TEXT NOT NULL,
		publication_date TIMESTAMP NOT NULL,
		url TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT,
		tags JSONB,              
		entities JSONB,           
		sentiment_score VARCHAR(10), 
		categories JSONB,         
		content_s3_path TEXT,     
		status TEXT NOT NULL,      
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);`

	_, err := pool.Exec(ctx, createTableSQL)
	return err
}

func InsertArticleData(ctx context.Context, pool *pgxpool.Pool, article ArticleSchema) error {
	tagsJSON, err := json.Marshal(article.Tags)
	if err != nil {
		return fmt.Errorf("error marshaling tags: %v", err)
	}

	entitiesDataJSON, err := json.Marshal(article.Entities)
	if err != nil {
		return fmt.Errorf("error marshaling entities: %v", err)
	}

	insertSQL := `
    INSERT INTO articles (article_id, title, publisher, publication_date, url, content, summary, tags, entities, sentiment_score, categories, content_s3_path, status)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	categoriesArray := pq.Array(article.Categories) // Convert categories slice to PostgreSQL array

	_, err = pool.Exec(ctx, insertSQL,
		article.ArticleId,
		article.Title,
		article.Publisher,
		article.PublicationDate,
		article.Url,
		article.Content,
		article.Summary,
		tagsJSON,               // Insert tags as JSONB
		entitiesDataJSON,       // Insert entities as JSONB
		article.SentimentScore, // Insert sentiment score as VARCHAR
		categoriesArray,        // Insert categories as an array
		article.ContentS3Path,
		article.Status) // Insert status

	return err
}

func GetArticleByID(ctx context.Context, pool *pgxpool.Pool, articleId string) (*ArticleSchema, error) {
	var article ArticleSchema

	query := `
	SELECT *
	FROM articles
	WHERE article_id = $1;`

	err := pool.QueryRow(ctx, query, articleId).Scan(
		&article.ArticleId,
		&article.Title,
		&article.Publisher,
		&article.PublicationDate,
		&article.Url,
		&article.Content,
		&article.Summary,
		&article.Tags,
		&article.Entities,
		&article.SentimentScore,
		pq.Array(&article.Categories), // Assuming pq is imported for array handling
		&article.ContentS3Path,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error fetching article: %v", err)
	}

	return &article, nil
}

func UpdateArticleByID(ctx context.Context, pool *pgxpool.Pool, article ArticleSchema) error {
	// Marshal tags to JSONB
	tagsJSON, err := json.Marshal(article.Tags)
	if err != nil {
		return fmt.Errorf("error marshaling tags: %v", err)
	}

	// Marshal entities to JSONB
	entitiesJSON, err := json.Marshal(article.Entities)
	if err != nil {
		return fmt.Errorf("error marshaling entities: %v", err)
	}

	updateSQL := `
    UPDATE articles
    SET title = $1,
        publisher = $2,
        publication_date = $3,
        url = $4,
        content = $5,
        summary = $6,
        tags = $7,
        entities = $8,
        sentiment_score = $9,
        categories = $10,
        content_s3_path = $11,
        status = $12,
        updated_at = CURRENT_TIMESTAMP  -- Automatically set updated_at to current time
    WHERE article_id = $13;`

	_, err = pool.Exec(ctx, updateSQL,
		article.Title,
		article.Publisher,
		article.PublicationDate,
		article.Url,
		article.Content,
		article.Summary,
		tagsJSON,                     // Insert tags as JSONB
		entitiesJSON,                 // Insert entities as JSONB
		article.SentimentScore,       // Insert sentiment score as VARCHAR
		pq.Array(article.Categories), // Insert categories as an array
		article.ContentS3Path,        // Insert content S3 path
		article.Status,               // Insert status
		article.ArticleId)            // Article ID to identify the row to update

	return err
}
