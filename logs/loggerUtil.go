package logs

import (
	"context"
	"log"
	"service-news-app-backend/config"
	envUtil "service-news-app-backend/config"

	"cloud.google.com/go/logging"
	"google.golang.org/api/option"
)

type ErrorPayloadData struct {
	ApiEndPoint     string      `json:"apiEndPoint,omitempty" bson:"apiEndPoint,omitempty"`
	ReqHeaders      interface{} `json:"reqHeaders,omitempty" bson:"reqHeaders,omitempty"`
	ReqBody         interface{} `json:"reqBody,omitempty" bson:"reqBody,omitempty"`
	ReqQueryParams  interface{} `json:"reqQueryParams,omitempty" bson:"reqQueryParams,omitempty"`
	ResponseStatus  int64       `json:"responseStatus,omitempty" bson:"responseStatus,omitempty"`
	Response        interface{} `json:"response,omitempty" bson:"response,omitempty"`
	ErrorDetails    string      `json:"errorDetails,omitempty" bson:"errorDetails,omitempty"`
	ErrorStackTrace string      `json:"errorStackTrace,omitempty" bson:"errorStackTrace,omitempty"`
}

func LogRequest(apiEndPoint string, reqHeaders interface{}, reqBody interface{}, reqQueryParams interface{}, responseStatus int64, response interface{}, errorDetails string, functionName string, errorNumber string) {
	var stackTrace string
	if errorNumber != "" {
		stackTrace = "service-news-app-backend:" + functionName + ":" + errorNumber
	}

	CreateLogEntry(
		ErrorPayloadData{
			ApiEndPoint:     apiEndPoint,
			ReqHeaders:      reqHeaders,
			ReqBody:         reqBody,
			ReqQueryParams:  reqQueryParams,
			ResponseStatus:  responseStatus,
			Response:        response,
			ErrorDetails:    errorDetails,
			ErrorStackTrace: stackTrace,
		},
	)
}

func CreateLogEntry(payloadData interface{}) {
	ctx := context.Background()

	// fmt.Println("logging started")

	projectID := config.GetEnvironmentVariable("projectId")

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsFile("loggingConfig.json"))
	if err != nil {
		log.Println("Failed to create client:", err)
		return
	}

	logName := envUtil.GetEnvironmentVariable("ahaanCrudAPILog")

	// Selects the log to write to.
	logger := client.Logger(logName)

	lables := map[string]string{
		"bucket_name": "ahaan-backend-error-logs",
		"location":    "asia-south1",
	}

	// Adds an entry to the log buffer.
	logger.Log(logging.Entry{Payload: payloadData, Severity: logging.Error, Labels: lables})

	// Closes the client and flushes the buffer to the Cloud Logging
	// service.
	if err := client.Close(); err != nil {
		log.Println("Failed to close client:", err)
		return
	}

	// fmt.Println("logging end")

	return
}

type ErrorPayload struct {
	StackTrace   string      `json:"stackTrace,omitempty" bson:"stackTrace,omitempty"` //serviceName:functionName:LineNumberOfCode
	ErrorDetails string      `json:"errorDetails,omitempty" bson:"errorDetails,omitempty"`
	ReqDetails   interface{} `json:"reqDetails,omitempty" bson:"reqDetails,omitempty"`
}

func CreateLogEntryForCasdoor(payloadData interface{}) {
	ctx := context.Background()

	// fmt.Println("logging started")

	projectID := config.GetEnvironmentVariable("projectId")

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID, option.WithCredentialsFile("loggingConfig.json"))
	if err != nil {
		log.Println("Failed to create client:", err)
		return
	}

	logName := envUtil.GetEnvironmentVariable("casdoorAPILogName")

	// Selects the log to write to.
	logger := client.Logger(logName)

	lables := map[string]string{
		"bucket_name": "ahaan-backend-error-logs",
		"location":    "asia-south1",
	}

	// Adds an entry to the log buffer.
	logger.Log(logging.Entry{Payload: payloadData, Severity: logging.Error, Labels: lables})

	// Closes the client and flushes the buffer to the Cloud Logging
	// service.
	if err := client.Close(); err != nil {
		log.Println("Failed to close client:", err)
		return
	}

	// fmt.Println("logging end")

	return
}
