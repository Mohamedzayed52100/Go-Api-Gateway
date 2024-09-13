package guest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	openapi "github.com/goplaceapp/goplace-gateway/openapi/go"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/gorilla/websocket"
	"google.golang.org/grpc/metadata"
	"net/http"
	"sync"
)

func (s *SGuest) GetRealtimeReservations() openapi.ContextHandler {
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

		shiftId, date, token := common.ConvertStringToInt(c.Query("shiftId")), c.Query("date"), c.Query("authorization")
		if shiftId == 0 || date == "" {
			c.JSON(http.StatusBadRequest, openapi.GeneralError{Code: http.StatusBadRequest, Message: "Invalid shift id or date"})
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
			accMeta: accMeta,
			ctx:     ctx,
			date:    date,
			shiftId: shiftId,
			cancel:  cancel,
			srv:     s,
			send:    make(chan interface{}, 256),
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			client.listenForUpdates(token, shiftId, date)
		}()

		go client.ReadPump()
		go client.startSendingMessages()

		wg.Wait()
	}
}

func (c *WSClient) listenForUpdates(uToken string, shiftId int32, date string) {
	var stream guestProto.Reservation_GetRealtimeReservationsClient
	var err error

	c.ctx = metadata.NewOutgoingContext(c.ctx, metadata.New(map[string]string{
		meta.TenantDBNameContextKey.String():  c.accMeta.ClientDBName,
		meta.AuthorizationContextKey.String(): uToken,
	}))

	stream, err = c.srv.ReservationClient.GetRealtimeReservations(c.ctx, &guestProto.GetRealtimeReservationsRequest{
		ShiftId: shiftId,
		Date:    date,
	})
	if err != nil {
		logger.Default().Error("Error getting realtime reservations", err)
		return
	}

	// Close the stream when the function returns
	defer stream.CloseSend()

	for {
		select {
		case <-c.ctx.Done():
			logger.Default().Info("Context cancelled, stopping stream")
			return
		default:
			res, err := stream.Recv()
			if err != nil {
				logger.Default().Error("Error receiving realtime reservation", err)
				return
			}

			if res.GetResult().Id == 0 {
				continue
			}

			c.srv.broadcastReservationUpdate(res, c)
		}
	}
}

func (s *SGuest) broadcastReservationUpdate(reservation interface{}, c *WSClient) {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case c.send <- buildReservationAPIResponse(reservation.(*guestProto.GetReservationByIDResponse).GetResult()):
		logger.Default().Info("Sent reservation to client")
	default:
		close(c.send)
	}
}
