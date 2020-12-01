package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// StartServer starts the http server
func StartServer() {
	addr := fmt.Sprintf("%s:%d", "localhost", 1232)

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	httpServer := http.Server{
		Addr:    addr,
		Handler: GetRouter(),
	}

	go func() {
		log.Println("Http server has been started on", addr)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("Failed to start http server on :", err)
			os.Exit(-1)
		}
	}()
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("Http server couldn't shutdown gracefully")
	}
	log.Println("Http server has been shutdown gracefully")
}
