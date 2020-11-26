package sns

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/llusardo/go-lambda-examples/dynamo-stream-to-sns/sns/models"
)

var awsSession *session.Session
var snsClient *sns.SNS
var SNSTopic = os.Getenv("SNS_TOPIC")

func init() {
	log.Println("SNS client Initialization")
	awsSession = session.New()
	snsClient = sns.New(awsSession)
}

func SendToSNS(eventType string, itemJsonString []byte, marshalErr error) error {
	event := models.Event{
		EventType: eventType,
		Data:      string(itemJsonString),
	}
	eventStr, marshalErr := json.Marshal(event)
	if marshalErr != nil {
		log.Fatal("Error marshalling event to send SNS", marshalErr)
		return marshalErr
	}

	messageStr := string(eventStr)

	req, _ := snsClient.PublishRequest(&sns.PublishInput{
		TopicArn: aws.String(SNSTopic),
		Message:  aws.String(messageStr),
	})

	fmt.Printf("Sending %s to SNS.\n", messageStr)

	sendErr := req.Send()
	if sendErr != nil {
		log.Fatal("Error sending message to SNS", sendErr)
		return sendErr
	}

	return nil
}
