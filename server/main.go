package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const maxBufferSize = 1024

var (
	port     = flag.Uint("port", 1337, "port to send to or receive from")
	host     = flag.String("host", "127.0.0.1", "address to send to or receive from")
	timeout  = flag.Duration("timeout", 15*time.Second, "read and write blocking deadlines")
)

func server(ctx context.Context, address string) (err error) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}

	defer conn.Close()

	doneChan := make(chan error, 1)
	buffer := make([]byte, maxBufferSize)

	go func() {
		for {
			bytesRead, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n",
				bytesRead, addr.String())

			deadline := time.Now().Add(*timeout)
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

func main () {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		err     error
		address = fmt.Sprintf("%s:%d", *host, *port)
	)

	// Gracefully handle signals so that we can finalize any of our
	// blocking operations by cancelling their contexts.
	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	fmt.Println("running as a server on " + address)
	err = server(ctx, address)
	if err != nil && err != context.Canceled {
		panic(err)
	}
}
