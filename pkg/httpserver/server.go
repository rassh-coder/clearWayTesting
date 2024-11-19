package httpserver

import (
	"clearWayTest/config"
	"clearWayTest/pkg/handler"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func Server(cfg *config.Config, h *handler.Handler) {
	router := h.InitRoutes()

	server := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.ListenAndServeTLS(fmt.Sprintf("./certs/%s.crt", cfg.HTTP.CertName), fmt.Sprintf("./certs/%s.key", cfg.HTTP.CertKey)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Listen and serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutdown server...")
	time.Sleep(time.Second * 5)
	if err := server.Shutdown(context.TODO()); err != nil {
		log.Printf("Server shutdown returned an err: %v\n", err)
	}
}
