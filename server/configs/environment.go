package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func GetHost() string {
	host := os.Getenv("HOST")
	if host != "" {
		return host
	}
	panic("the HOST environment variable must be declared")
}

func GetPort() int {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println(err.Error())
		panic("the PORT environment variable must be declared")
	}
	return int(port)
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

