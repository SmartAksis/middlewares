package sqs_broker

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
)

type convert func([]byte)

type SqsSenderObject interface {
	ToAwsSender() *MessageSender
}

type MessageSender struct {
	ID string
	DuplicationID string
	MessageBody []byte
}

func awsSession() (*session.Session, error) {
	awsSqsKey := os.Getenv("AWS_SQS_KEY")
	awsSqsSecret := os.Getenv("AWS_SQS_SECRET")
	return session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(awsSqsKey, awsSqsSecret, ""),},
	)
}

func buildQueueUrl(queueName string) string {
	region := os.Getenv("AWS_REGION")
	account := os.Getenv("AWS_ACCOUNT_NUMBER")
	envi := os.Getenv("ENVI")
	return fmt.Sprintf("https://sqs.%d.amazonaws.com/%d/%d-%d", region, account, envi, queueName)
}

func buildDeleteMessageInput(queueURL *string, receiptHandle *string) *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl:		queueURL,
		ReceiptHandle: 	receiptHandle,
	}
}

func buildReceiveMessageInput(queueURL *string, size int64) *sqs.ReceiveMessageInput {
	return &sqs.ReceiveMessageInput{
		QueueUrl:    			queueURL,
		MaxNumberOfMessages: 	aws.Int64(size),
	}
}

func buildSendMessageInput(queueURL *string, duplicationID *string, messageBody *string, messageGroupId *string, attributes map[string]*sqs.MessageAttributeValue) *sqs.SendMessageInput {
	return &sqs.SendMessageInput{
		MessageDeduplicationId: duplicationID,
		MessageAttributes: 		attributes,
		MessageBody: 			messageBody,
		QueueUrl:    			queueURL,
		MessageGroupId: 		messageGroupId,
	}
}

func buildSendMessageBatchInput(queueURL string, senders []*MessageSender, messageGroupId *string, attributes map[string]*sqs.MessageAttributeValue) *sqs.SendMessageBatchInput {
	batch:=make([]*sqs.SendMessageBatchRequestEntry, len(senders))

	for i, input := range senders {
		batch[i] = &sqs.SendMessageBatchRequestEntry{
			Id:							aws.String(input.ID),
			MessageAttributes:	       	attributes,
			MessageBody:        	    aws.String(string(input.MessageBody)),
			MessageDeduplicationId: 	aws.String(input.DuplicationID),
			MessageGroupId:          	messageGroupId,
		}
	}
	return &sqs.SendMessageBatchInput{
		Entries:  batch,
		QueueUrl: aws.String(queueURL),
	}
}

func sqsConsume(session *session.Session, queueUrl string, constructor convert, size int64, messageManipulation func(message sqs.Message, _sqs *sqs.SQS)) (*sqs.SQS, error) {
	if session != nil {
		sqs := sqs.New(session)
		result, err := sqs.ReceiveMessage(buildReceiveMessageInput(aws.String(queueUrl), size))
		if err != nil{
			return sqs, err
		}
		if sqs != nil && result != nil {
			for _, message := range result.Messages {
				if message.Body != nil {
					bytes:=[]byte(*message.Body)
					constructor(bytes)
				}
				if messageManipulation != nil {
					messageManipulation(*message, sqs)
				}
			}
			return sqs, nil
		}
	}
	return nil, nil
}

func sqsDelete(_sqs *sqs.SQS, queueUrl string, receiptHandle *string) {
	param := buildDeleteMessageInput(aws.String(queueUrl), receiptHandle)
	if param != nil {
		_, err := _sqs.DeleteMessage(param)
		if err != nil {
			fmt.Println("Error to Delete Message on Queue")
			return
		}
	}
}

func SqsPublish(queueUrl string, sender *MessageSender, _messageGroupId string, attributes map[string]*sqs.MessageAttributeValue) (*string, error) {
	sess, err := awsSession()
	if err != nil {
		return nil, err
	}
	if sess != nil {
		sqs := sqs.New(sess)
		queueUrl := aws.String(queueUrl)
		duplicationID := aws.String(sender.DuplicationID)
		messageBody := aws.String(string(sender.MessageBody))
		messageGroupId := aws.String("GROUP_DEFAULT")
		attributes := attributes
		if _messageGroupId != ""{
			messageGroupId = aws.String(_messageGroupId)
		}

		result, err := sqs.SendMessage(buildSendMessageInput(queueUrl, duplicationID, messageBody, messageGroupId, attributes))
		if err != nil {
			return nil, err
		}
		return result.MessageId, err
	}
	return nil, nil
}

func SqsPublishBatch(queueUrl string, senders []*MessageSender, _messageGroupId *string, attributes map[string]*sqs.MessageAttributeValue) ([]string, []string, error) {
	sess, err := awsSession()
	if len(senders) > 10 {

	}

	if sess != nil {
		sqs := sqs.New(sess)
		messageGroupId := aws.String("GROUP_DEFAULT")
		if _messageGroupId != nil{
			messageGroupId = aws.String(*_messageGroupId)
		}
		result, err := sqs.SendMessageBatch(buildSendMessageBatchInput(queueUrl, senders, messageGroupId, attributes))
		if err != nil {
			fmt.Println(err)
			return nil, nil, err
		}

		var success []string
		var failed []string

		if len(result.Successful) > 0 {
			success = make([]string, len(result.Successful))
			for index, entry := range result.Successful {
				success[index] = *entry.MessageId
			}
		}

		if len(result.Failed) > 0 {
			failed = make([]string, len(result.Failed))
			for index, entry := range result.Successful {
				failed[index] = *entry.MessageId
			}
		}

		return success, failed, err
	}
	fmt.Println(err)
	return nil, nil, err
}

func SqsConsume(queueUrl string, constructor convert) {
	SqsConsumeMultiple(queueUrl, 1, constructor)
}

func SqsConsumeWithoutDelete(queueUrl string, constructor convert) {
	SqsConsumeWithoutDeleteMultiple(queueUrl, 1, constructor)
}

func SqsConsumeMultiple(queueUrl string, size int64, constructor convert) {
	sess, err := awsSession()
	if size > 10 {
		size = 10
	}
	if err == nil {
		sqsConsume(sess, queueUrl, constructor, size, func(message sqs.Message, _sqs *sqs.SQS) {
			sqsDelete(_sqs, queueUrl, message.ReceiptHandle)
		})
	}
}

func SqsConsumeWithoutDeleteMultiple(queueUrl string, size int64, constructor convert) {
	sess, err := awsSession()
	if size > 10 {
		size = 10
	}
	if err == nil {
		sqsConsume(sess, queueUrl, constructor, size, nil)
	}
}
