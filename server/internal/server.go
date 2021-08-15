package internal

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	"github.com/carol-caires/udp-chat/internal/utils/cache"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type Server struct {
	cache cache.Client
}

func NewServer (cache cache.Client) Server {
	return Server{cache}
}

func (s *Server) Listen(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen to packets")
		return
	}
	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, configs.GetMaxBufferSize())

	go func() {
		for {
			bytesRead, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				log.Error().Err(err).Msg("failed to read from buffer")
				doneChan <- err
				return
			}

			// todo: write message in cache

			log.Debug().Msgf("packet-received: bytes=%d from=%s", bytesRead, addr.String())

			deadline := time.Now().Add(time.Second * configs.GetBlockingDeadline())
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				log.Error().Err(err).Msg("failed set write blocking deadline")
				doneChan <- err
				return
			}

			bytesWritten, err := conn.WriteTo(buffer[:bytesRead], addr)
			if err != nil {
				log.Error().Err(err).Msg("failed set send packet to clients")
				doneChan <- err
				return
			}

			log.Debug().Msgf("packet-written: bytes=%d to=%s", bytesWritten, addr.String())
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
