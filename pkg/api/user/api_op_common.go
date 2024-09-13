package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/goplaceapp/goplace-gateway/pkg/api/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type SUser struct {
	*common.Resources
}

type CreateUserForm struct {
	EmployeeId  string `form:"employee_id" binding:"required"`
	FirstName   string `form:"first_name" binding:"required"`
	LastName    string `form:"last_name" binding:"required"`
	Email       string `form:"email" binding:"required"`
	PhoneNumber string `form:"phone_number" binding:"required"`
	Role        string `form:"role" binding:"required"`
	JoinedAt    string `form:"joined_at" binding:"required"`
	Birthdate   string `form:"birthdate" binding:"required"`
	BranchIds   string `form:"branch_ids" binding:"required"`
}

type UpdateUserForm struct {
	EmployeeId  string `form:"employee_id"`
	FirstName   string `form:"first_name"`
	LastName    string `form:"last_name"`
	Email       string `form:"email"`
	PhoneNumber string `form:"phone_number"`
	Role        string `form:"role"`
	JoinedAt    string `form:"joined_at"`
	Birthdate   string `form:"birthdate"`
	BranchIds   string `form:"branch_ids"`
}

func getQueueURL(sess *session.Session, queue string) (string, error) {
	sqsClient := sqs.New(sess)
	queue = strings.ReplaceAll(queue, ".", "-")
	queue = strings.ReplaceAll(queue, "@", "-")

	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queue,
	})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return *result.QueueUrl, nil
}

func createQueue(sess *session.Session, queueName string) (*sqs.CreateQueueOutput, error) {
	sqsClient := sqs.New(sess)
	queueName = strings.ReplaceAll(queueName, ".", "-")
	queueName = strings.ReplaceAll(queueName, "@", "-")
	result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &queueName,
		Attributes: map[string]*string{
			"DelaySeconds":      aws.String("0"),
			"VisibilityTimeout": aws.String("60"),
		},
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
