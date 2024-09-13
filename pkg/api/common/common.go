package common

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	reportsProto "github.com/goplaceapp/shared-protobufs/reports/go_out"

	"github.com/redis/go-redis/v9"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/h2non/filetype"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"

	"net/http"

	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/httphelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
)

var RdbInstance = redis.NewClient(&redis.Options{
	Addr: os.Getenv("REDIS_ADDRESS") + ":" + os.Getenv("REDIS_PORT"),
})

var AwsSession, _ = session.NewSession(&aws.Config{
	Region: aws.String(os.Getenv("S3_REGION")),
})

var SqsClient = sqs.New(AwsSession)

type Config struct{}

type Resources struct {
	Log         *zap.SugaredLogger
	BaseCfg     *httphelper.BaseConfig
	SvcCfg      *Config
	Engine      *gin.Engine
	GrpcClients map[string]*grpchelper.GrpcClient

	// User service clients
	UserClient       userProto.UserClient
	TenantClient     userProto.TenantClient
	RoleClient       userProto.RoleClient
	DepartmentClient userProto.DepartmentClient

	// Guest service clients
	GuestClient                      guestProto.GuestClient
	GuestLogClient                   guestProto.GuestLogClient
	ReservationClient                guestProto.ReservationClient
	ReservationLogClient             guestProto.ReservationLogClient
	ReservationSpecialOccasionClient guestProto.ReservationSpecialOccasionClient
	ReservationFeedbackClient        guestProto.ReservationFeedbackClient
	ReservationFeedbackCommentClient guestProto.ReservationFeedbackCommentClient
	ReservationWaitlistClient        guestProto.ReservationWaitlistClient
	ReservationFeedbackWebhookClient guestProto.ReservationFeedbackWebhookClient
	DayOperationsClient              guestProto.DayOperationsClient
	ReservationWidgetClient          guestProto.ReservationWidgetClient
	PaymentClient                    guestProto.PaymentClient

	// Settings service clients
	SettingsClient          settingsProto.SettingsClient
	ShiftClient             settingsProto.ShiftClient
	TableClient             settingsProto.TableClient
	SeatingAreaClient       settingsProto.SeatingAreaClient
	AutomatedReportClient   settingsProto.AutomatedReportClient
	FloorPlanClient         settingsProto.FloorPlanClient
	GuestTagClient          settingsProto.GuestTagsClient
	ReservationTagClient    settingsProto.ReservationTagsClient
	ReservationStatusClient settingsProto.ReservationStatusClient
	WidgetSettingsClient    settingsProto.WidgetSettingsClient
	RestaurantItemClient    settingsProto.RestaurantItemClient

	// Customer service clients
	CSSInquiryClient        cssProto.InquiryServiceClient
	CSSInquiryCommentClient cssProto.InquiryCommentServiceClient
	CSSEmployeeClient       cssProto.EmployeeServiceClient
	CSSReportClient         cssProto.ReportServiceClient

	// Reports service clients
	OverviewReportClient  reportsProto.OverviewReportsServiceClient
	FinancialReportClient reportsProto.FinancialReportsServiceClient
	GuestReportClient     reportsProto.GuestReportsServiceClient
	CustomizeReportClient reportsProto.CustomizeReportServiceClient
}

// HandleGrpcError handles grpc errors returned by the gateway and writes it to the gin context
func HandleGrpcError(c *gin.Context, err error) {
	logger.Default().With(
		"code", c.Writer.Status(),
		"error", err,
	).Error("grpc error")

	s, ok := status.FromError(err)
	if !ok {
		WriteErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.JSON(int(s.Code()), openapi.GeneralError{
		Code:    int32(s.Code()),
		Message: s.Message(),
	})
}

// PropagateKeysFromGinContextToGoContext propagates the keys from the gin context to the go context
func PropagateKeysFromGinContextToGoContext(goCtx context.Context, ginCtx *gin.Context) context.Context {
	for _, key := range grpchelper.KeysToPropagate {
		value, ok := ginCtx.Value(key.String()).(string)
		if ok && len(value) > 0 {
			goCtx = context.WithValue(goCtx, key.String(), value)
		}
	}
	return goCtx
}

func writeErrorJSON(c *gin.Context, code int, message string) {
	errorResponse := openapi.GeneralError{Code: int32(code), Message: message}
	c.AbortWithStatusJSON(code, errorResponse)
}

// GetRequiredContextParam enforces having the given URL param and rerun it
func GetRequiredContextParam(c *gin.Context, name string) (string, error) {
	param := c.Param(name)
	if param == "" {
		return "", fmt.Errorf("missing URL param %s", name)
	}

	return param, nil
}

// WriteErrorResponse writes the error with the given code and message to the gin context
func WriteErrorResponse(c *gin.Context, code int, message string) {
	logger.Default().With(
		"code", code,
		meta.EndpointContextKey.String(), c.Request.URL.Path,
	).Error(message)

	writeErrorJSON(c, code, message)
}

// ResolveRequestBody bind the body to the given struct, fails if the request body is empty
func ResolveRequestBody(c *gin.Context, input interface{}) (ok bool) {
	if c.Request.Body == http.NoBody {
		WriteErrorResponse(c, http.StatusBadRequest, "empty request")
		return false
	}

	if err := c.ShouldBindJSON(input); err != nil {
		WriteErrorResponse(c, http.StatusBadRequest, err.Error())
		return false
	}

	return true
}

// ValidateRequestBody validate given request struct and sends validation errors response
func ValidateRequestBody(input interface{}) error {
	v := validator.New()

	err := v.Struct(input)
	if err != nil {
		return err
	}

	return nil
}

// HandleValidationErrors handles validation errors
func HandleValidationErrors(c *gin.Context, validationErrors error) {
	if validationErrors != nil {
		c.JSON(http.StatusBadRequest,
			openapi.ValidationErrors{
				GeneralError: openapi.GeneralError{
					Code:    http.StatusBadRequest,
					Message: "validation error",
				},
			})
		return
	}

	WriteErrorResponse(c, http.StatusInternalServerError, "internal server error")
}

func ConvertStringToInt(str string) int32 {
	value, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0
	}
	return int32(value)
}

func ConvertStringToBool(str string) bool {
	value, _ := strconv.ParseBool(str)
	return value
}

func ConvertIntToString(value int32) string {
	return strconv.FormatInt(int64(value), 10)
}

func ParseQueryCriteriaToIntArray(criteria string) []int32 {
	if criteria == "" {
		return []int32{}
	}

	criteriaArray := strings.Split(criteria, ",")
	result := make([]int32, len(criteriaArray))
	for i := range criteriaArray {
		result[i] = ConvertStringToInt(criteriaArray[i])
	}

	return result
}

func ParseQueryCriteriaToStringArray(criteria string) []string {
	if criteria == "" {
		return []string{}
	}

	criteriaArray := strings.Split(criteria, ",")
	result := make([]string, len(criteriaArray))
	for i := range criteriaArray {
		result[i] = criteriaArray[i]
	}

	return result
}

func generateUniqueTimestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func UploadFileToS3(file io.Reader, folder string, allowedExtensions []string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)

	buf, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	kind, _ := filetype.Match(buf)
	if !contains(allowedExtensions, kind.Extension) {
		return "", fmt.Errorf("unsupported file type. please upload one of: %v", allowedExtensions)
	}

	extension := kind.Extension
	fileName := generateUniqueTimestamp() + "." + extension

	upParams := &s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(fmt.Sprintf("%s/%s", folder, fileName)),
		Body:        bytes.NewReader(buf),
		ContentType: aws.String(kind.MIME.Value),
	}

	result, err := uploader.Upload(upParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	return result.Location, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Capitalize first letter
func Capitalize(s string) string {
	s = strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
	for i := range s {
		if s[i] == ' ' {
			s = strings.ToLower(s[0:i]) + strings.ToUpper(s[i+1:i+2]) + strings.ToLower(s[i+2:])
		}
	}

	return s
}

// RemoveHyphen from string
func RemoveHyphen(s string) string {
	return strings.ReplaceAll(s, "-", "")
}

// CapitalizeWords in string
func CapitalizeWords(s string) string {
	words := strings.Split(s, "-")
	for i := range words {
		words[i] = Capitalize(words[i])
	}

	return strings.Join(words, "-")
}

// ConvertSortingCriteria converts sorting criteria from string to []string
func ConvertSortingCriteria(criteria string) []string {
	if criteria == "" {
		return []string{}
	}

	criteria = RemoveHyphen(CapitalizeWords(criteria))

	result := strings.Split(criteria, ",")
	for i := range result {
		if result[i] == "" {
			continue
		}

		result[i] = strings.ToLower(result[i][:1]) + result[i][1:]
	}

	return result
}
