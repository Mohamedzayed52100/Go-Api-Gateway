package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/goplaceapp/goplace-common/pkg/middleware"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
)

const (
	apiPrefix      string = "/api"
	healthEndpoint string = "/health"
)

func createOpenAPIHandler(_ string, handler openapi.ContextHandler) func(c *gin.Context) {
	return func(c *gin.Context) {
		ctx := common.PropagateKeysFromGinContextToGoContext(c.Request.Context(), c)
		handler(ctx, c)
	}
}

// ClientIP Example IP lookup function that extracts the client IP from the request
func ClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	return ip
}

func (s *Service) setupMiddlewares(apiGroup *gin.RouterGroup, routesWithPermission []struct{ Path, Method, Permissions string }) {
	endpointsWithoutAccessTokenCheck := []string{
		healthEndpoint,
		openapi.LoginEndpoint,
		openapi.RequestDemoEndpoint,
		openapi.GetRealtimeReservationsEndpoint,
		openapi.GetRealtimeWaitingReservationsEndpoint,
		openapi.CreateReservationFeedbackFromWebhookEndpoint,
		openapi.UpdateReservationFromWebhookEndpoint,
		openapi.GetNewMessagesEndpoint,
		openapi.GetWidgetSettingsEndpoint,
		openapi.GetWidgetAvailableTimesEndpoint,
		openapi.GetWidgetAllSeatingAreasEndpoint,
		openapi.GetWidgetAllBranchesEndpoint,
		openapi.GetWidgetAllSpecialOccasionsEndpoint,
		openapi.CreateWidgetReservationEndpoint,
		openapi.GetAllCountriesEndpoint,
		openapi.GetCountryByNameEndpoint,
		openapi.RequestResetPasswordEndpoint,
		openapi.ResendOtpEndpoint,
		openapi.VerifyOtpEndpoint,
		openapi.ResetPasswordEndpoint,
		openapi.UpdatePaymentFromWebhookEndpoint,
		openapi.RequestReservationWebhookEndpoint,
	}

	endpointsUsingClientCredentials := []string{
		openapi.CreateReservationFeedbackFromWebhookEndpoint,
		openapi.UpdateReservationFromWebhookEndpoint,
		openapi.RequestReservationWebhookEndpoint,
	}

	endpointsUsingWidgetClientName := []string{
		openapi.GetWidgetSettingsEndpoint,
		openapi.GetWidgetAvailableTimesEndpoint,
		openapi.GetWidgetAllSeatingAreasEndpoint,
		openapi.GetWidgetAllBranchesEndpoint,
		openapi.CreateWidgetReservationEndpoint,
		openapi.GetWidgetAllSpecialOccasionsEndpoint,
	}

	apiGroup.Use(sentrygin.New(sentrygin.Options{}))
	//ipRateLimiter := middleware.NewIPRateLimiter(rate.Every(5*time.Minute), 500, ClientIP)
	//apiGroup.Use(middleware.RateLimitMiddleware(ipRateLimiter))
	apiGroup.Use(middleware.CheckAccessTokenMiddleware(endpointsWithoutAccessTokenCheck...))
	apiGroup.Use(middleware.CheckClientCredentialsMiddleware(s.TenantClient, endpointsUsingClientCredentials, endpointsWithoutAccessTokenCheck))
	apiGroup.Use(middleware.RBACMiddleware(routesWithPermission, s.UserClient, endpointsWithoutAccessTokenCheck...))
	apiGroup.Use(middleware.HandleWidgetMiddleware(s.TenantClient, endpointsUsingWidgetClientName...))
	apiGroup.Use(middleware.CheckOtpVerified(endpointsWithoutAccessTokenCheck,s.UserClient))
}

func (s *Service) createRoutes(apiGroup *gin.RouterGroup) {
	for _, route := range openapi.CreateRoutes(s) {
		openAPIHandler := createOpenAPIHandler(route.Pattern, route.HandlerFunc)
		switch route.Method {
		case http.MethodGet:
			apiGroup.GET(route.Pattern, openAPIHandler)
		case http.MethodPost:
			apiGroup.POST(route.Pattern, openAPIHandler)
		case http.MethodPut:
			apiGroup.PUT(route.Pattern, openAPIHandler)
		case http.MethodDelete:
			apiGroup.DELETE(route.Pattern, openAPIHandler)
		}
	}
}

func (s *Service) generateRoutesWithPermission() []struct{ Path, Method, Permissions string } {
	var routesWithPermission []struct{ Path, Method, Permissions string }
	for _, route := range openapi.CreateRoutes(s) {
		routesWithPermission = append(routesWithPermission, struct{ Path, Method, Permissions string }{
			Path:        route.Pattern,
			Method:      route.Method,
			Permissions: route.Permissions,
		})
	}

	return routesWithPermission
}

func (s *Service) NewEngine() *gin.Engine {
	ginEngine := gin.New()

	rootGroup := ginEngine.Group("/")
	rootGroup.GET(healthEndpoint, func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	apiGroup := ginEngine.Group(apiPrefix)

	apiGroup.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return false
		},
		MaxAge: 12 * time.Hour,
	}))

	// setup sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://b4cabe43daa1573f398c301cc8e191bd@o4505889732755456.ingest.sentry.io/4505889744093184",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if req, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					event.Request = &sentry.Request{
						URL:         req.URL.String(),
						Method:      req.Method,
						QueryString: req.URL.RawQuery,
					}
				}
			}

			return event
		},
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v", err)
	}

	// generate routes with permission
	routesWithPermission := s.generateRoutesWithPermission()

	// setup middlewares
	s.setupMiddlewares(apiGroup, routesWithPermission)

	// create actual routes
	s.createRoutes(apiGroup)

	return ginEngine
}
