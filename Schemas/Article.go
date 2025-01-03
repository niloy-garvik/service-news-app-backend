package schemas

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"service-news-app-backend/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type ArticleSchema struct {
	ArticleId       string         `json:"articleId"` // UUID as string
	Title           string         `json:"title"`
	Publisher       string         `json:"publisher"`
	PublicationDate time.Time      `json:"publicationDate"` // TIMESTAMP
	Url             string         `json:"url"`
	Content         string         `json:"content"`
	Summary         string         `json:"summary"`
	Tags            pq.StringArray `json:"tags"`           // -- e.g., ["Indian Army", "Kashmir"]
	Entities        any            `json:"entities"`       // -- e.g., {"organizations": ["Indian Army"], "locations": ["Kashmir"]}
	SentimentScore  string         `json:"sentimentScore"` // -- e.g., "Positive", "Negative", "Neutral"
	Categories      pq.StringArray `json:"categories"`     // -- e.g., ["National Security", "Conflict"]
	ContentS3Path   string         `json:"contentS3Path"`  // -- e.g., S3 path to the full article
	Status          string         `json:"status"`         // -- e.g., "published" or "unpublished"
	CreatedAt       time.Time      `json:"createdAt"`      // TIMESTAMP
	UpdatedAt       time.Time      `json:"updatedAt"`      // TIMESTAMP
}

// CreateArticlesTable creates the articles table in the database
func CreateArticlesTable(ctx context.Context, pool *pgxpool.Pool) error {
	articleTableName := config.GetEnvironmentVariable("ARTICLE_TABLE_NAME")

	createTableSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		article_id UUID PRIMARY KEY,
		title TEXT NOT NULL,
		publisher TEXT NOT NULL,
		publication_date TIMESTAMP NOT NULL,
		url TEXT NOT NULL,
		content TEXT NOT NULL,
		summary TEXT,
		tags TEXT[],              -- Changed to TEXT[] for array of tags
		entities JSONB,           -- Keeping entities as JSONB for flexibility
		sentiment_score VARCHAR(10), 
		categories TEXT[],         -- Changed to TEXT[] for array of categories
		content_s3_path TEXT,     
		status TEXT NOT NULL,      
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`, articleTableName)

	// Execute the SQL command to create the table
	_, err := pool.Exec(ctx, createTableSQL)
	if err != nil {
		fmt.Println("Error creating articles table: ", err)
		return err
	}

	log.Printf("%s table created successfully or already exists\n", articleTableName)
	return nil
}
func InsertArticleData(ctx context.Context, pool *pgxpool.Pool, article ArticleSchema) error {

	articleTableName := config.GetEnvironmentVariable("ARTICLE_TABLE_NAME")

	entitiesDataJSON, err := json.Marshal(article.Entities)
	if err != nil {
		return fmt.Errorf("error marshaling entities: %v", err)
	}

	insertSQL := fmt.Sprintf(`
    INSERT INTO %s (article_id, title, publisher, publication_date, url, content, summary, tags, entities, sentiment_score, categories, content_s3_path, status)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13);`, articleTableName)

	// categoriesArray := pq.Array(article.Categories) // Convert categories slice to PostgreSQL array
	// tagsArray := pq.Array(article.Tags)

	_, err = pool.Exec(ctx, insertSQL,
		article.ArticleId,
		article.Title,
		article.Publisher,
		article.PublicationDate,
		article.Url,
		article.Content,
		article.Summary,
		article.Tags,           // Insert categories as an array
		entitiesDataJSON,       // Insert entities as JSONB
		article.SentimentScore, // Insert sentiment score as VARCHAR
		article.Categories,     // Insert tags as JSONB
		article.ContentS3Path,
		article.Status) // Insert status

	return err
}

// GetArticleByID retrieves an article by its ID
func GetArticleByID(ctx context.Context, pool *pgxpool.Pool, articleId string) (*ArticleSchema, error) {
	articleTableName := config.GetEnvironmentVariable("ARTICLE_TABLE_NAME")

	var article ArticleSchema

	query := fmt.Sprintf(`
	SELECT *
	FROM %s
	WHERE article_id = $1;`, articleTableName)

	err := pool.QueryRow(ctx, query, articleId).Scan(
		&article.ArticleId,
		&article.Title,
		&article.Publisher,
		&article.PublicationDate,
		&article.Url,
		&article.Content,
		&article.Summary,
		&article.Tags, // Scan tags as an array of strings
		&article.Entities,
		&article.SentimentScore,
		&article.Categories, // Scan categories as an array of strings
		&article.ContentS3Path,
		&article.Status,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No record found, return nil for both article and error
		}
		return nil, fmt.Errorf("error fetching article: %v", err)
	}

	return &article, nil
}
func UpdateArticleByID(ctx context.Context, pool *pgxpool.Pool, article ArticleSchema) error {

	articleTableName := config.GetEnvironmentVariable("ARTICLE_TABLE_NAME")

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

	updateSQL := fmt.Sprintf(`
    UPDATE %s
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
    WHERE article_id = $13;`, articleTableName)

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
