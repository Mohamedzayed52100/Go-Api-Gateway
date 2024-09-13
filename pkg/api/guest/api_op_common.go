package guest

import (
	"context"
	"fmt"

	"sort"
	"strings"
	"sync"

	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-gateway/pkg/api/common"
	guestProto "github.com/goplaceapp/goplace-guest/api/v1"
	"github.com/gorilla/websocket"
	"github.com/pusher/pusher-http-go/v5"
)

var PusherClient *pusher.Client

type WSClient struct {
	id      string
	conn    *websocket.Conn
	send    chan interface{}
	date    string
	shiftId int32
	accMeta auth.AccountMetaData
	ctx     context.Context
	cancel  context.CancelFunc
	srv     *SGuest
}

type SGuest struct {
	*common.Resources
	mu sync.Mutex
}

// SortReservationsByGuestsNumber sorts reservations by guests number
func SortReservationsByGuestsNumber(reservations []*guestProto.Reservation, criteria string) {
	sort.Slice(reservations, func(i, j int) bool {
		if criteria == "asc" {
			return reservations[i].GuestsNumber < reservations[j].GuestsNumber
		} else {
			return reservations[i].GuestsNumber > reservations[j].GuestsNumber
		}
	})
}

// SortReservationsByName sorts reservations by name
func SortReservationsByName(reservations []*guestProto.Reservation, criteria string) {
	sort.Slice(reservations, func(i, j int) bool {
		var primaryGuest1 *guestProto.ReservationGuest
		for _, g := range reservations[i].Guests {
			if g.IsPrimary {
				primaryGuest1 = g
				break
			}
		}

		var primaryGuest2 *guestProto.ReservationGuest
		for _, g := range reservations[j].Guests {
			if g.IsPrimary {
				primaryGuest2 = g
				break
			}
		}

		firstFullName := primaryGuest1.FirstName + " " + primaryGuest1.LastName
		secondFullName := primaryGuest2.FirstName + " " + primaryGuest2.LastName
		if criteria == "asc" {
			return firstFullName < secondFullName
		} else {
			return firstFullName > secondFullName
		}
	})
}

// SortReservationsByTime sorts reservations by time
func SortReservationsByTime(reservations []*guestProto.Reservation, criteria string) {
	sort.Slice(reservations, func(i, j int) bool {
		firstTime := reservations[i].Time.AsTime().Hour()*60 + reservations[i].Time.AsTime().Minute()
		secondTime := reservations[j].Time.AsTime().Hour()*60 + reservations[j].Time.AsTime().Minute()

		if reservations[i].Time.AsTime().Hour() == 0 {
			firstTime += 24 * 60
		}

		if reservations[j].Time.AsTime().Hour() == 0 {
			secondTime += 24 * 60
		}

		if criteria == "asc" {
			return firstTime < secondTime
		} else {
			return firstTime > secondTime
		}
	})
}

// SortGuestsByName sorts guests by name
func SortGuestsByName(guests []*guestProto.Guest, criteria string) {
	sort.Slice(guests, func(i, j int) bool {
		firstFullName := guests[i].FirstName + " " + guests[i].LastName
		secondFullName := guests[j].FirstName + " " + guests[j].LastName
		if criteria == "asc" {
			return strings.Compare(strings.ToUpper(firstFullName), strings.ToUpper(secondFullName)) == -1
		} else {
			return strings.Compare(strings.ToUpper(firstFullName), strings.ToUpper(secondFullName)) != -1
		}
	})
}

// SortGuestsByLastVisit sorts guests by last visit
func SortGuestsByLastVisit(guests []*guestProto.Guest, criteria string) {
	sort.Slice(guests, func(i, j int) bool {
		if criteria == "asc" {
			return guests[i].LastVisit.AsTime().Before(guests[j].LastVisit.AsTime())
		} else {
			return guests[i].LastVisit.AsTime().After(guests[j].LastVisit.AsTime())
		}
	})
}

// ReadPump reads messages from the client
func (c *WSClient) ReadPump() {
	defer func() {
		c.srv.mu.Lock()
		c.srv.mu.Unlock()
		c.conn.Close()
	}()

	messageChannel := make(chan []byte)
	errorChannel := make(chan error)

	go func() {
		for {
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				errorChannel <- err
				return
			}
			messageChannel <- message
		}
	}()

	for {
		select {
		case message := <-messageChannel:
			if string(message) == "close" {
				c.cancel()
				return
			} else if string(message) == "ping" {
				c.conn.WriteMessage(websocket.TextMessage, []byte("pong"))
			}
		case err := <-errorChannel:
			fmt.Println("An error occurred while reading message: ", err)
			c.cancel()
			return
		case <-c.ctx.Done():
			fmt.Println("Context is done")
			return
		}
	}
}

// WritePump writes messages to the client
func (c *WSClient) WritePump() {
	go func() {
		<-c.ctx.Done()
		c.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		c.conn.Close()
	}()

	for {
		select {
		case message := <-c.send:
			err := c.conn.WriteJSON(message)
			if err != nil {
				fmt.Println("An error occurred while writing message: ", err)
				return
			}
		case <-c.ctx.Done():
			fmt.Println("Context is done")
			return
		}
	}
}

func parseQueryCriteriaToIntArray(criteria string) []int32 {
	if criteria == "" {
		return []int32{}
	}

	criteriaArray := strings.Split(criteria, ",")
	result := make([]int32, len(criteriaArray))
	for i := range criteriaArray {
		result[i] = common.ConvertStringToInt(criteriaArray[i])
	}

	return result
}

func parseQueryCriteriaToStringArray(criteria string) []string {
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

func (c *WSClient) startSendingMessages() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				logger.Default().Error("Error listening to messages")
				return
			}
			if err := c.conn.WriteJSON(message); err != nil {
				logger.Default().Error("Error sending message", err)
				return
			}
		case <-c.ctx.Done():
			logger.Default().Info("Context cancelled, stopping sending messages")
			return
		}
	}
}
