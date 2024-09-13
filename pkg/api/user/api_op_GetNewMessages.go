package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/auth"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

func (s *SUser) GetNewMessages() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		token := c.Query("authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Invalid authorization token"})
			return
		}

		accMeta := auth.GetAccountMetadataFromToken(token)
		if accMeta.ClientDBName == "" {
			c.JSON(http.StatusUnauthorized, openapi.GeneralError{Code: http.StatusUnauthorized, Message: "Invalid authorization token"})
			return
		}

		channelName := strings.ReplaceAll(accMeta.Email, ".", "-")
		channelName = "messages:" + strings.ReplaceAll(channelName, "@", "-")

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Could not open WebSocket connection"})
			return
		}
		defer ws.Close()

		//wg := sync.WaitGroup{}
		//wg.Add(1)
		//go readNewMessages(ws, accMeta.Email, channelName)
		//go readNewMessages(ws, "broadcast", "")
		//wg.Wait()
	}
}

func readNewMessages(ws *websocket.Conn, queueName, channelName string) {
	ctx := context.Background()
	rdb := common.RdbInstance

	sess := common.AwsSession

	msgQueue, err := getQueueURL(sess, queueName)
	if err != nil {
		if _, err := createQueue(sess, queueName); err != nil {
			ws.WriteJSON(gin.H{"error": err.Error()})
			return
		}
	}

	sqsClient := common.SqsClient

	for {
		result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            &msgQueue,
			MaxNumberOfMessages: aws.Int64(10),
			WaitTimeSeconds:     aws.Int64(10),
		})
		if err != nil {
			ws.WriteJSON(openapi.GeneralError{Code: http.StatusInternalServerError, Message: "Could not read messages from SQS queue: " + queueName})
			return
		}

		if len(result.Messages) == 0 {
			time.Sleep(1 * time.Second)
		}

		for _, message := range result.Messages {
			var jsonMessage map[string]interface{}

			err = json.Unmarshal([]byte(*message.Body), &jsonMessage)
			if err != nil {
				ws.WriteJSON(openapi.GeneralError{Code: http.StatusInternalServerError, Message: "Could not parse message"})
			} else {
				ws.WriteJSON(jsonMessage)
			}

			if err := rdb.ZAdd(ctx, channelName, redis.Z{
				Member: *message.Body,
			}).Err(); err != nil {
				ws.WriteJSON(openapi.GeneralError{Code: http.StatusInternalServerError, Message: err.Error()})
			}

			if _, err := sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      &msgQueue,
				ReceiptHandle: message.ReceiptHandle,
			}); err != nil {
				ws.WriteJSON(openapi.GeneralError{Code: http.StatusInternalServerError, Message: err.Error()})
			}
		}
	}
}
