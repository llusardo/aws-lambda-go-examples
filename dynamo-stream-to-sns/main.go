package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/llusardo/go-lambda-examples/dynamo-stream-to-sns/sns"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, e events.DynamoDBEvent) error {
	var err error

	// If the processing of one event of the batch fails, the entire batch will be reprocessing
	for _, record := range e.Records {
		err = processDynamoDBRecord(record)
		if err != nil {
			return err
		}
	}

	return err
}

func processDynamoDBRecord(record events.DynamoDBEventRecord) error {
	fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

	var item map[string]events.DynamoDBAttributeValue

	switch record.EventName {
	case "INSERT":
		fallthrough
	case "MODIFY":
		item = record.Change.NewImage
	case "REMOVE":
		item = record.Change.OldImage
	}

	itemJson, unMarshalErr := UnmarshalStreamImage(item)
	if unMarshalErr != nil {
		log.Fatal("Error unmarshalling dynamo item", unMarshalErr)
		return unMarshalErr
	}

	itemJsonString, marshalErr := json.Marshal(itemJson)
	if marshalErr != nil {
		log.Fatal("Error marshalling retrieved item", marshalErr)
		return marshalErr
	}

	err := sns.SendToSNS(record.EventName, itemJsonString, marshalErr)
	if err != nil {
		log.Fatal("Error sending to SNS", err)
		return err
	}

	return nil
}

// UnmarshalStreamImage converts events.DynamoDBAttributeValue to struct
func UnmarshalStreamImage(attribute map[string]events.DynamoDBAttributeValue) (map[string]interface{}, error) {
	dbAttrMap := make(map[string]*dynamodb.AttributeValue)

	for k, v := range attribute {

		var dbAttr dynamodb.AttributeValue

		bytes, marshalErr := v.MarshalJSON()
		if marshalErr != nil {
			log.Fatal("Error marshaling stream image attribute", marshalErr)
			return nil, marshalErr
		}

		unmarshalErr := json.Unmarshal(bytes, &dbAttr)
		if unmarshalErr != nil {
			log.Fatal("Error unMarshaling stream image attribute", unmarshalErr)
			return nil, unmarshalErr
		}

		dbAttrMap[k] = &dbAttr
	}

	var anyJson map[string]interface{}

	err := dynamodbattribute.UnmarshalMap(dbAttrMap, &anyJson)

	if err != nil {
		log.Fatal("Error unMarshaling stream image attribute Map", err)
		return nil, err
	}

	return anyJson, nil
}
