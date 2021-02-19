package sqs_broker

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"os"
	"strings"
	"strconv"
)

var(
	sqsInstance *sqs.SQS
)

type onSuccess func(string, []byte, string, string)

type CallbackConsumerInterface interface {
	OnSuccess(queueUrl string, data []byte, messageId string, receivedHandle string)
}

func (c *callbackConsumer) OnSuccess(queueUrl string, data []byte, messageId string, receivedHandle string){
	c.success(queueUrl, data, messageId, receivedHandle);
}

type callbackConsumer struct {
	success onSuccess
}

func CreateCallBack(_onSucess onSuccess) CallbackConsumerInterface{
	return &callbackConsumer{
		success: _onSucess,
	}
}

type SqsSenderObject interface {
	ToAwsSender() *MessageSender
}

type MessageSender struct {
	ID string
	DuplicationID string
	MessageBody []byte
}

func initSqs(){
	if sqsInstance == nil {
		sess, err := awsSession()
		sqsInstance = sqs.New(sess)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func DeleteCallbackConsumer() CallbackConsumerInterface {
	return &callbackConsumer{
		success: func(queueUrl string, content []byte, messageId string, receiptHandle string) {
			param := buildDeleteMessageInput(aws.String(queueUrl), aws.String(receiptHandle))
			if param != nil {
				_, err := sqsInstance.DeleteMessage(param)
				if err != nil {
					fmt.Println("Error to Delete Message on Queue")
					return
				}
			}
		},
	}
}

func awsSession() (*session.Session, error) {
	awsSqsKey := os.Getenv("AWS_SQS_KEY")
	awsSqsSecret := os.Getenv("AWS_SQS_SECRET")
	region := os.Getenv("AWS_REGION")
	return session.NewSession(&aws.Config{
		Region: aws.String(strings.ToLower(region)),
		Credentials: credentials.NewStaticCredentials(awsSqsKey, awsSqsSecret, ""),},
	)
}

func buildQueueUrl(queueName *string) string {
	region := strings.ToLower(os.Getenv("AWS_REGION"))
	account := os.Getenv("AWS_ACCOUNT_NUMBER")
	envi := strings.ToLower(os.Getenv("ENVI"))
	stmt:=fmt.Sprintf("https://sqs.%s.amazonaws.com/%s/%s-%s", region, account, envi, strings.ReplaceAll(strconv.Quote(*queueName), "\"", ""))
	return stmt
}

func buildDeleteMessageInput(queueURL *string, receiptHandle *string) *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl:		aws.String(buildQueueUrl(queueURL)),
		ReceiptHandle: 	receiptHandle,
	}
}

func buildReceiveMessageInput(queueURL *string, size int64) *sqs.ReceiveMessageInput {
	return &sqs.ReceiveMessageInput{
		QueueUrl:    			aws.String(buildQueueUrl(queueURL)),
		MaxNumberOfMessages: 	aws.Int64(size),
	}
}

func buildSendMessageInput(queueURL *string, duplicationID *string, messageBody *string, messageGroupId *string, attributes map[string]*sqs.MessageAttributeValue) *sqs.SendMessageInput {
	return &sqs.SendMessageInput{
		MessageDeduplicationId: duplicationID,
		MessageAttributes: 		attributes,
		MessageBody: 			messageBody,
		QueueUrl:    			aws.String(buildQueueUrl(queueURL)),
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

func sqsConsume(queueUrl string, size int64, callbacks ...CallbackConsumerInterface) (bool, error) {
	result, err := sqsInstance.ReceiveMessage(buildReceiveMessageInput(aws.String(queueUrl), size))
	if err == nil{
		for _, message := range result.Messages {
			if callbacks != nil {
				for _, callback := range callbacks {
					if message.Body != nil {
						bytes:=[]byte(*message.Body)
						if callback.OnSuccess != nil {
							callback.OnSuccess(queueUrl, bytes, *message.MessageId, *message.ReceiptHandle)
						}
					}
				}
			}
		}
		return true, nil
	}
	return false, err
}


func SqsPublish(queueUrl string, sender *MessageSender, _messageGroupId string, attributes map[string]*sqs.MessageAttributeValue) (*string, error) {
	initSqs()
	duplicationID := aws.String(sender.DuplicationID)
	messageBody := aws.String(string(sender.MessageBody))
	messageGroupId := aws.String("GROUP_DEFAULT")
	if _messageGroupId != ""{
		messageGroupId = aws.String(_messageGroupId)
	}
	result, err := sqsInstance.SendMessage(buildSendMessageInput(aws.String(queueUrl), duplicationID, messageBody, messageGroupId, attributes))
	if err != nil {
		return nil, err
	}
	return result.MessageId, err
}

func SqsPublishBatch(queueUrl string, senders []*MessageSender, _messageGroupId *string, attributes map[string]*sqs.MessageAttributeValue) ([]string, []string, error) {
	initSqs()
	messageGroupId := aws.String("GROUP_DEFAULT")
	if _messageGroupId != nil{
		messageGroupId = aws.String(*_messageGroupId)
	}
	result, err := sqsInstance.SendMessageBatch(buildSendMessageBatchInput(queueUrl, senders, messageGroupId, attributes))
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

func SqsConsume(queueUrl string, callback ...CallbackConsumerInterface) {
	initSqs()
	SqsConsumeMultiple(queueUrl, 1, callback...)
}

func SqsConsumeWithoutDelete(queueUrl string, callback ...CallbackConsumerInterface) {
	initSqs()
	SqsConsumeWithoutDeleteMultiple(queueUrl, 1, callback...)
}

func SqsConsumeMultiple(queueUrl string, size int64, callback ...CallbackConsumerInterface) {
	initSqs()
	if size > 10 {
		size = 10
	}
	//var _callBacks []CallbackConsumerInterface
	//if callback != nil {
	//	_callBacks = make([]CallbackConsumerInterface, len(callback)+1)
	//	for i, consumer := range callback {
	//		_callBacks[i] = consumer
	//	}
	//} else {
	//	_callBacks = make([]CallbackConsumerInterface, 1)
	//}

	sqsConsume(queueUrl, size, append(callback, DeleteCallbackConsumer())...)
}

func SqsConsumeWithoutDeleteMultiple(queueUrl string, size int64, callback ...CallbackConsumerInterface) {
	initSqs()
	if size > 10 {
		size = 10
	}
	sqsConsume(queueUrl, size, callback...)
}
