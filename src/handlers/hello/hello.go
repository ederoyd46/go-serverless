package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	var buf bytes.Buffer

	addDynamoDBEntry()

	body, err := json.Marshal(map[string]interface{}{
		"message": "Go Serverless v1.0! Your function executed successfully!",
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "hello-handler",
		},
	}

	return resp, nil
}

// Item struct to hold info about new item
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func addDynamoDBEntry() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	item := Item{ID: 1, Name: "Hello"}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		panic("Marshal went wrong")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("matt-go-dev-HelloTable-Z9T2M7Y0IRIT")}

	_, err = svc.PutItem(input)

	if err != nil {
		fmt.Println("Error calling PutItem")
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Added item")
}

func main() {
	lambda.Start(Handler)
}
