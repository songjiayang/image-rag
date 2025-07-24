package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"image-rag-backend/cmd/mcp/tool"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

var (
	mode         string
	serverlisten string
)

var (
	serverName    = "Image Rag Tools"
	serverVersion = "v0.1.0"
)

func init() {
	flag.StringVar(&mode, "transport", "sse", "The transport to use, should be \"stdio\" or \"sse\"")
	flag.StringVar(&serverlisten, "server_listen", "127.0.0.1:8081", "The sse server listen address")
	flag.Parse()
}

func getTransport() (t transport.ServerTransport) {
	if mode == "stdio" {
		log.Println("start mcp server with stdio transport")
		t = transport.NewStdioServerTransport()
	} else {
		log.Printf("start mcp server with sse transport, listen %s", serverlisten)
		t, _ = transport.NewSSEServerTransport(serverlisten)
	}

	return t
}

func main() {
	svr, _ := server.NewServer(
		getTransport(),
		server.WithServerInfo(protocol.Implementation{
			Name:    serverName,
			Version: serverVersion,
		}),
	)

	svr.RegisterTool(tool.NewImageRagTool(), tool.ImageRagHandler)

	// register poster tool
	runWithSignalWaiter(svr)
}

func runWithSignalWaiter(srv *server.Server) {
	errCh := make(chan error)
	go func() {
		errCh <- srv.Run()
	}()

	if err := signalWaiter(errCh); err != nil {
		log.Fatalf("signal waiter: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
}

func signalWaiter(errCh chan error) error {
	signalToNotify := []os.Signal{syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM}
	if signal.Ignored(syscall.SIGHUP) {
		signalToNotify = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, signalToNotify...)

	select {
	case sig := <-signals:
		switch sig {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
			log.Printf("Received signal: %s\n", sig)
			// graceful shutdown
			return nil
		}
	case err := <-errCh:
		log.Printf("Run server with error: %v", err)
		return err
	}

	return nil
}
