package guest

import (
	"context"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"google.golang.org/grpc/metadata"
)

func (s *SGuest) GetRealtimeWaitingReservations() openapi.ContextHandler {
	return func(ctx context.Context, c *gin.Context) {
		var upgrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Could not open WebSocket connection"})
			return
		}
		defer ws.Close()

		shiftId, token := common.ConvertStringToInt(c.Query("shiftId")), c.Query("authorization")
		if shiftId == 0 {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Invalid shift id"})
			return
		} else if token == "" {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Invalid authorization token"})
			return
		}

		accMeta := auth.GetAccountMetadataFromToken(token)
		if accMeta.ClientDBName == "" {
			c.JSON(http.StatusUnauthorized, openapi.GeneralError{Code: http.StatusUnauthorized, Message: "Invalid authorization token"})
			return
		}

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		client := &WSClient{
			id:      uuid.New().String(),
			conn:    ws,
			send:    make(chan interface{}, 256),
			accMeta: accMeta,
			ctx:     ctx,
			cancel:  cancel,
			srv:     s,
		}

		date := c.Query("date")

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.listenForWaitingReservationsUpdates(token, &guestProto.GetWaitingReservationRequest{
				ShiftId: shiftId,
				Date:    date,
			})
		}()

		go client.ReadPump()
		go client.startSendingMessages()

		wg.Wait()
	}
}

func (c *WSClient) listenForWaitingReservationsUpdates(uToken string, req *guestProto.GetWaitingReservationRequest) {
	c.ctx = metadata.NewOutgoingContext(c.ctx, metadata.New(map[string]string{
		meta.TenantDBNameContextKey.String():  c.accMeta.ClientDBName,
		meta.AuthorizationContextKey.String(): uToken,
	}))

	stream, err := c.srv.ReservationWaitlistClient.GetRealtimeWaitingReservations(c.ctx, &guestProto.GetWaitingReservationRequest{
		ShiftId: req.ShiftId,
		Date:    req.Date,
	})
	if err != nil {
		logger.Default().Error("Error while listening for waiting reservations updates", err)
		return
	}

	// Close the stream when the function returns
	defer stream.CloseSend()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			res, err := stream.Recv()
			if err != nil {
				logger.Default().Error("Error while receiving waiting reservations updates", err)
				return
			}

			if res.GetResult().Id == 0 {
				continue
			}

			c.srv.broadcastWaitingReservationUpdate(res, c)
		}
	}
}

func (s *SGuest) broadcastWaitingReservationUpdate(reservation interface{}, c *WSClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case c.send <- buildReservationWaitlistResponse(reservation.(*guestProto.GetWaitingReservationResponse).GetResult()):
		logger.Default().Info("Sent reservation to client")
	default:
		close(c.send)
	}
}
