package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	portFlag := flag.Int("p", 8080, "port number used into the http server")
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())

	mux := http.NewServeMux()
	mux.HandleFunc("/", handle)

	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", *portFlag),
		Handler:     mux,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	httpServer.RegisterOnShutdown(cancel)

	go func() {
		log.Printf("starting server on :%d\n", *portFlag)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("error on listen and serve: %v\n", err)
		}
	}()

	// capture os interrupt
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		os.Interrupt,
		syscall.SIGTERM,
	)

	// waiting for the os interrupt
	<-signalCh

	log.Println("received interrupt signal - shutting down")

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	}

	log.Printf("service stopped\n")
	defer os.Exit(0)
}

func handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello! You found me ;)"))
}
