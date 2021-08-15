package service

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	"github.com/carol-caires/udp-chat/internal/infrastructure/cache"
	"github.com/carol-caires/udp-chat/internal/model"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type Server struct {
	cache cache.Client
	conn net.PacketConn
}

func NewServer (cache cache.Client) Server {
	return Server{cache: cache}
}

func (s *Server) Listen(ctx context.Context, address string) (err error) {
	s.conn, err = net.ListenPacket("udp", address)
	if err != nil {
		log.Error().Err(err).Msg("failed to listen to packets")
		return
	}
	defer s.conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, configs.GetMaxBufferSize())

	go func() {
		for {
			bytesRead, addr, err := s.conn.ReadFrom(buffer)
			if err != nil {
				log.Error().Err(err).Msg("failed to read from buffer")
				doneChan <- err
				return
			}

			log.Debug().Msgf("packet-received: bytes=%d from=%s", bytesRead, addr.String())
			log.Info().Msg("trying to sync messages cache")

			message, jsonMessage, err := model.NewMessage(string(buffer[:bytesRead]))
			if err != nil {
				log.Error().Err(err).Msgf("message have incorrect format: %s")
				doneChan <- err
				return
			}

			switch message.MessageType {
			case model.MessageTypeMessage:
				err = s.getConnectedClientsAndBroadcastMessage(ctx, message, jsonMessage)
				break
			case model.MessageTypeNewClient:
				// todo: sync with clients array in redis
				break
			case model.MessageTypeCloseClient:
				// todo: remove from clients array in redis
				break
			}

			if err != nil {
				doneChan <- err
				return
			}
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}

func (s *Server) getConnectedClientsAndBroadcastMessage(ctx context.Context, message *model.Message, jsonMessage string) (err error) {
	clients, err := s.getConnectedClients(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed get connected clients")
		return
	}

	err = s.broadcastMessage(clients, message.Body)
	if err != nil {
		log.Error().Err(err).Msg("failed to broadcast to clients")
		return
	}

	err = s.cache.Set(ctx, fmt.Sprintf("message:%s", message.Id), jsonMessage)
	if err != nil {
		log.Error().Err(err).Msg("failed to broadcast to sync message in cache")
	}
	return
}

func (s *Server) getConnectedClients(ctx context.Context) (clients []model.Client, err error) {
	connClientsStr, err := s.cache.Get(ctx, "clients")
	if err != nil {
		return
	}

	clients, err = model.ParseClientsArray(connClientsStr)
	return
}

func (s *Server) broadcastMessage(clients []model.Client, message string) (err error) {
	deadline := time.Now().Add(time.Second * configs.GetBlockingDeadline())
	err = s.conn.SetWriteDeadline(deadline)
	if err != nil {
		log.Error().Err(err).Msg("failed set write blocking deadline")
		return
	}

	for _, client := range clients {
		tcpAddr, _ := net.ResolveTCPAddr("ip", client.IpAddress)
		bytesWritten, err := s.conn.WriteTo([]byte(message), tcpAddr)
		if err != nil {
			log.Error().Err(err).Msg("failed set send packet to clients")
			return
		}
		log.Debug().Msgf("packet-written: bytes=%d to=%s", bytesWritten, client.IpAddress)
	}

	log.Info().Msg("finished to broadcast a new message to clients")
	return
}
