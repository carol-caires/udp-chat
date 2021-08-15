package configs

import (
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"time"
)

func GetHost() string {
	host := os.Getenv("HOST")
	if host != "" {
		return host
	}
	log.Fatal().Msg("the HOST environment variable must be declared")
	return ""
}

func GetPort() int {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatal().Err(err).Msg("the PORT environment variable must be declared")
		return 0
	}
	return port
}

func GetBlockingDeadline() time.Duration {
	deadlineStr := os.Getenv("BLOCKING_DEADLINE_SECONDS")
	deadline, err := strconv.Atoi(deadlineStr)
	if err != nil {
		return 15
	}
	return time.Duration(deadline)
}

func GetMaxBufferSize() int {
	bufferStr := os.Getenv("MAX_BUFFER_SIZE_BYTES")
	buffer, err := strconv.Atoi(bufferStr)
	if err != nil {
		return 1024
	}
	return buffer
}

func GetRedisAddr() string {
	addr := os.Getenv("REDIS_ADDR")
	if addr != "" {
		return addr
	}
	log.Fatal().Msg("the HOST environment variable must be declared")
	return ""
}

