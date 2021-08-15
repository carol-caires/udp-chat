package main

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	"github.com/carol-caires/udp-chat/internal"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func main () {
	log := configs.InitLogs()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	var address = fmt.Sprintf("%s:%d", configs.GetHost(), configs.GetPort())

	server := internal.NewServer(&log)
	err = server.Listen(ctx, address)
	if err != nil && err != context.Canceled {
		log.Fatal("error starting UDP server: ", err.Error())
	}

	log.Info("running UDP server on: ", address)
}
