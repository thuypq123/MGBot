package repository

import (
	"discord_bot/internal/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetQueueURL(sess *session.Session, queue *string) (*sqs.GetQueueUrlOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	result, err := svc.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: queue,
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func SendMsg(sess *session.Session, queueURL *string, m string) error {
	svc := sqs.New(sess)

	idMessage := utils.Uuid()

	// create a new message with uuid
	_, err := svc.SendMessage(&sqs.SendMessageInput{
		MessageBody:            aws.String(m),
		QueueUrl:               queueURL,
		MessageGroupId:         aws.String("messageGroup1"),
		MessageDeduplicationId: aws.String(idMessage),
	})
	if err != nil {
		return err
	}

	return nil
}

func GetMessages(sess *session.Session, queueURL *string) (*sqs.ReceiveMessageOutput, error) {
	// Create an SQS service client
	svc := sqs.New(sess)

	// snippet-start:[sqs.go.receive_messages.call]
	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            queueURL,
		MaxNumberOfMessages: aws.Int64(1),
	})
	// snippet-end:[sqs.go.receive_messages.call]
	if err != nil {
		return nil, err
	}

	return msgResult, nil
}

func DeleteMessage(sess *session.Session, queueURL *string, messageHandle *string) error {
	// snippet-start:[sqs.go.delete_message.call]
	svc := sqs.New(sess)

	_, err := svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      queueURL,
		ReceiptHandle: messageHandle,
	})
	// snippet-end:[sqs.go.delete_message.call]
	if err != nil {
		return err
	}

	return nil
}
