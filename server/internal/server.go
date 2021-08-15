package internal

import (
	"context"
	"fmt"
	"github.com/carol-caires/udp-chat/configs"
	"net"
	"time"
)

func Server(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, configs.GetMaxBufferSize())

	go func() {
		for {
			bytesRead, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n",
				bytesRead, addr.String())

			deadline := time.Now().Add(time.Second * configs.GetBlockingDeadline())
			err = conn.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}

			bytesRead, err = conn.WriteTo(buffer[:bytesRead], addr)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", bytesRead, addr.String())
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
