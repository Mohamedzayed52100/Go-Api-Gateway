package main

import (
	"fmt"
	"log"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/httphelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-gateway/pkg/api/guest"
	"github.com/goplaceapp/goplace-gateway/pkg/service"
	"github.com/joho/godotenv"
	"github.com/pusher/pusher-http-go/v5"
)

func init() {
	// Load .env file
	godotenv.Load()

	fmt.Println("zizoooooooooo")
	log.Println("zayedLog")

	// Initialize global variables or configurations here
	meta.TokenSymmetricKey = os.Getenv("TOKEN_SYMMETRIC_KEY")
	if meta.TokenSymmetricKey == "" {
		log.Fatal("TOKEN_SYMMETRIC_KEY must be set in the environment variables")
	}

	guest.PusherClient = &pusher.Client{
		AppID:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_KEY"),
		Secret:  os.Getenv("PUSHER_SECRET"),
		Cluster: os.Getenv("PUSHER_CLUSTER"),
		Secure:  true,
	}
}

func main() {
	s := service.New()
	httphelper.Start(s)
}
