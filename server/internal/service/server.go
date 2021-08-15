package service

import (
	"context"
	"encoding/json"
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
				log.Error().Err(err).Msgf("message have incorrect format: %s", jsonMessage)
				doneChan <- err
				return
			}

			switch message.MessageType {
			case model.MessageTypeNewMessage:
				err = s.getConnectedClientsAndBroadcastMessage(ctx, message, jsonMessage)
				break
			case model.MessageTypeDeleteMessage:
				// todo: delete message
				break
			case model.MessageTypeNewClient:
				err = s.saveClient(ctx, message.From.Name, addr.String())
				break
			case model.MessageTypeDeleteClient:
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

func (s *Server) getConnectedClientsAndBroadcastMessage(ctx context.Context, message model.Message, jsonMessage string) (err error) {
	clients, err := s.getConnectedClients(ctx)
	if err != nil {
		return
	}

	if len(clients) > 0 {
		err = s.broadcastMessage(clients, message.Body)
		if err != nil {
			log.Error().Err(err).Msg("failed to broadcast to clients")
			return
		}

		err = s.cache.Set(ctx, fmt.Sprintf("message:%s", message.Id), jsonMessage)
		if err != nil {
			log.Error().Err(err).Msg("failed to broadcast to sync message in cache")
		}
	}
	return
}

func (s *Server) getConnectedClients(ctx context.Context) (clients []model.Client, err error) {
	connClientsStr, _ := s.cache.Get(ctx, "clients")

	if connClientsStr == "" {
		log.Info().Msg("there is no connected clients, flushing messages db")
		// todo: remove all messages from cache
		return
	}

	clients, err = model.ParseClientsArray(connClientsStr)
	if err != nil {
		log.Error().Err(err).Msg("failed get connected clients")
	}
	return
}

func (s *Server) saveClient(ctx context.Context, name, address string) (err error) {
	clients, err := s.getConnectedClients(ctx)
	if err != nil {
		return
	}

	clients = append(clients, model.Client{
		Name:      name,
		IpAddress: address,
	})

	clientsStr, _ := json.Marshal(clients)
	err = s.cache.Set(ctx, "clients", string(clientsStr))
	if err != nil {
		log.Error().Err(err).Msg("failed save client")
		return
	}
	log.Info().Msgf("saved client %s with ip %s", name, address)
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
		log.Info().Msgf("sending message to client %s (address: %s): %s", client.Name, client.IpAddress, message)
		udpAddr, _ := net.ResolveUDPAddr("udp4", client.IpAddress)
		bytesWritten, errWrite := s.conn.WriteTo([]byte(message), udpAddr)
		if errWrite != nil {
			log.Error().Err(errWrite).Msg("failed sending packet to clients")
			return
		}
		log.Debug().Msgf("packet-written: bytes=%d to=%s", bytesWritten, client.IpAddress)
	}

	log.Info().Msg("finished to broadcast a new message to clients")
	return
}
