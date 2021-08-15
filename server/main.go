package main

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	"github.com/carol-caires/udp-chat/internal/service"
	"github.com/carol-caires/udp-chat/internal/infrastructure/cache"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

func main () {
	log.Info().Msg("starting...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	redisConn, err := cache.NewRedisConn()
	if err != nil {
		log.Fatal().Err(err).Msg("error starting redis connection")
	}

	var address = fmt.Sprintf("%s:%d", configs.GetHost(), configs.GetPort())
	log.Info().Msgf("running UDP service on address %s", address)

	server := service.NewServer(redisConn)
	err = server.Listen(ctx, address)
	if err != nil && err != context.Canceled {
		log.Fatal().Err(err).Msg("error starting UDP service")
	}
}
