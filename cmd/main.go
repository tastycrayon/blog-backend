package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tastycrayon/blog-backend/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("❗could not load .env file")
	}
	// The HTTP Server
	idleConnClosed := make(chan struct{})
	var router http.Handler
	router, err = Service(config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	srv := &http.Server{Addr: "0.0.0.0:" + config.HTTPServerPort, Handler: router}

	go GracefullyShutdown(srv, idleConnClosed)

	//Run the server
	log.Printf("🚀 Server ready at: http://localhost:%s", config.HTTPServerPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to start server %v", err)
	}

	<-idleConnClosed
	log.Println("Service Stop ✌")
}

func GracefullyShutdown(srv *http.Server, idleConnClosed chan struct{}) {
	sig := make(chan os.Signal, 1)
	defer close(sig)
	// Listen for syscall signals for process to interrupt/quit
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sig

	log.Println("interrupt received: server shutting down 🤷")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("http server shoutdown error %v", err)
	}
	log.Println("Shutdown Gracefully 🥁")
	close(idleConnClosed)
}
