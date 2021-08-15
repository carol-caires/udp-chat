package internal

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Server struct {
	log *log.Logger
}

func NewServer (log *log.Logger) Server {
	return Server{log}
}

func (s *Server) Listen(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Error("failed to listen to packets: ", err.Error())
		return
	}
	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, configs.GetMaxBufferSize())

	go func() {
		for {
			bytesRead, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				log.Error("failed to read from buffer: ", err.Error())
				doneChan <- err
				return
			}

			log.Debug(fmt.Sprintf("packet-received: bytes=%d from=%s\n",
				bytesRead, addr.String()))

			deadline := time.Now().Add(time.Second * configs.GetBlockingDeadline())
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				log.Error("failed set write blocking deadline: ", err.Error())
				doneChan <- err
				return
			}

			bytesRead, err = conn.WriteTo(buffer[:bytesRead], addr)
			if err != nil {
				log.Error("failed set send packet to clients: ", err.Error())
				doneChan <- err
				return
			}

			log.Debug(fmt.Sprintf("packet-written: bytes=%d to=%s\n", bytesRead, addr.String()))
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
