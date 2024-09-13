package service

import (
	"fmt"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-gateway/pkg/api/reports"
	reportsProto "github.com/goplaceapp/shared-protobufs/reports/go_out"

	"github.com/caarlos0/env/v6"

	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/httphelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	cssProto "github.com/goplaceapp/goplace-customer-service/api/v1"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	"github.com/goplaceapp/goplace-gateway/pkg/api/customerservice"
	"github.com/goplaceapp/goplace-gateway/pkg/api/general"
	"github.com/goplaceapp/goplace-gateway/pkg/api/guest"
	"github.com/goplaceapp/goplace-gateway/pkg/api/settings"
	"github.com/goplaceapp/goplace-gateway/pkg/api/user"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	settingsProto "github.com/goplaceapp/goplace-settings/api/v1"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	"go.uber.org/zap"
)

type Service struct {
	*common.Resources
	*general.SGeneral
	*user.SUser
	*guest.SGuest
	*settings.SSettings
	*customerservice.SCustomerService
	*reports.SReports
}

func New() *Service {
	// setup logger
	log, err := logger.New(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Errorf("failed to initialize the logger, %w", err))
	}

	// parse service config
	cfg := &common.Config{}
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Errorf("failed to read service config, %w", err))
	}

	resources := &common.Resources{
		Log:    log,
		SvcCfg: cfg,
		GrpcClients: map[string]*grpchelper.GrpcClient{
			meta.UserService:     {},
			meta.GuestService:    {},
			meta.SettingsService: {},
			meta.CSSService:      {},
			meta.ReportsService:  {},
		},
	}

	// register API resources
	s := &Service{
		Resources: resources,
		SUser: &user.SUser{
			Resources: resources,
		},
		SGuest: &guest.SGuest{
			Resources: resources,
		},
		SSettings: &settings.SSettings{
			Resources: resources,
		},
		SCustomerService: &customerservice.SCustomerService{
			Resources: resources,
		},
		SReports: &reports.SReports{
			Resources: resources,
		},
	}

	// prepare gRPC dependencies
	grpchelper.PrepareGrpcClientConnections(s.Log, s.GrpcClients)
	s.UserClient = userProto.NewUserClient(s.GrpcClients[meta.UserService].Conn)
	s.RoleClient = userProto.NewRoleClient(s.GrpcClients[meta.UserService].Conn)
	s.DepartmentClient = userProto.NewDepartmentClient(s.GrpcClients[meta.UserService].Conn)
	s.TenantClient = userProto.NewTenantClient(s.GrpcClients[meta.UserService].Conn)

	s.GuestClient = guestProto.NewGuestClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationClient = guestProto.NewReservationClient(s.GrpcClients[meta.GuestService].Conn)
	s.DayOperationsClient = guestProto.NewDayOperationsClient(s.GrpcClients[meta.GuestService].Conn)
	s.GuestLogClient = guestProto.NewGuestLogClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationLogClient = guestProto.NewReservationLogClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationFeedbackClient = guestProto.NewReservationFeedbackClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationFeedbackCommentClient = guestProto.NewReservationFeedbackCommentClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationFeedbackWebhookClient = guestProto.NewReservationFeedbackWebhookClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationSpecialOccasionClient = guestProto.NewReservationSpecialOccasionClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationWaitlistClient = guestProto.NewReservationWaitlistClient(s.GrpcClients[meta.GuestService].Conn)
	s.ReservationWidgetClient = guestProto.NewReservationWidgetClient(s.GrpcClients[meta.GuestService].Conn)
	s.PaymentClient = guestProto.NewPaymentClient(s.GrpcClients[meta.GuestService].Conn)

	s.SettingsClient = settingsProto.NewSettingsClient(s.GrpcClients[meta.SettingsService].Conn)
	s.ShiftClient = settingsProto.NewShiftClient(s.GrpcClients[meta.SettingsService].Conn)
	s.AutomatedReportClient = settingsProto.NewAutomatedReportClient(s.GrpcClients[meta.SettingsService].Conn)
	s.TableClient = settingsProto.NewTableClient(s.GrpcClients[meta.SettingsService].Conn)
	s.FloorPlanClient = settingsProto.NewFloorPlanClient(s.GrpcClients[meta.SettingsService].Conn)
	s.SeatingAreaClient = settingsProto.NewSeatingAreaClient(s.GrpcClients[meta.SettingsService].Conn)
	s.ReservationStatusClient = settingsProto.NewReservationStatusClient(s.GrpcClients[meta.SettingsService].Conn)
	s.ReservationTagClient = settingsProto.NewReservationTagsClient(s.GrpcClients[meta.SettingsService].Conn)
	s.GuestTagClient = settingsProto.NewGuestTagsClient(s.GrpcClients[meta.SettingsService].Conn)
	s.WidgetSettingsClient = settingsProto.NewWidgetSettingsClient(s.GrpcClients[meta.SettingsService].Conn)
	s.RestaurantItemClient = settingsProto.NewRestaurantItemClient(s.GrpcClients[meta.SettingsService].Conn)

	s.CSSInquiryClient = cssProto.NewInquiryServiceClient(s.GrpcClients[meta.CSSService].Conn)
	s.CSSInquiryCommentClient = cssProto.NewInquiryCommentServiceClient(s.GrpcClients[meta.CSSService].Conn)
	s.CSSEmployeeClient = cssProto.NewEmployeeServiceClient(s.GrpcClients[meta.CSSService].Conn)
	s.CSSReportClient = cssProto.NewReportServiceClient(s.GrpcClients[meta.CSSService].Conn)

	s.OverviewReportClient = reportsProto.NewOverviewReportsServiceClient(s.GrpcClients[meta.ReportsService].Conn)
	s.FinancialReportClient = reportsProto.NewFinancialReportsServiceClient(s.GrpcClients[meta.ReportsService].Conn)
	s.GuestReportClient = reportsProto.NewGuestReportsServiceClient(s.GrpcClients[meta.ReportsService].Conn)
	s.CustomizeReportClient = reportsProto.NewCustomizeReportServiceClient(s.GrpcClients[meta.ReportsService].Conn)

	// setup API engine
	s.Engine = s.NewEngine()

	return s
}

func (s *Service) GetGrpcClients() map[string]*grpchelper.GrpcClient {
	return s.GrpcClients
}

func (s *Service) GetEngine() *gin.Engine {
	return s.Engine
}

func (s *Service) SetBaseConfig(cfg *httphelper.BaseConfig) {
	s.BaseCfg = cfg
}

func (s *Service) GetLog() *zap.SugaredLogger {
	return s.Log
}

func (s *Service) GetHealthzHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) GetReadyzHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
