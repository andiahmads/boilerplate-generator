package main

import(
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"MyApp/infra"
	"MyApp/commons/helper"
	"MyApp/commons/logger"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"context"
	"crypto/tls"
	"fmt"
)


func main(){
  fmt.Println("Welcome to MyApp!")
	fmt.Println("Version: 1.0.0")
	failOnError := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s : %s", msg, err)
		}
	}

	dir := helper.DynamicDir()
	envPath := fmt.Sprintf("%s.env", dir)
	err := godotenv.Load(envPath)
	failOnError(err, "error load .env")

	db, err := infra.NewMySQLConn()
	failOnError(err, "error create mysql connection")
	defer db.Close()

	client, err := infra.NewRedisConn(0)
	failOnError(err, "failed to connect to redis")
	defer client.Close()

	ctx := context.Background()

	port := os.Getenv("SERVICE_PORT")
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})


	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second, //
		TLSConfig:    &tls.Config{InsecureSkipVerify: true},
	}

	serverErr := make(chan error, 1)
	server.SetKeepAlivesEnabled(true)
	go func() {
		logger.Slogger().Info("starting server", slog.String("port", port))
		serverErr <- server.ListenAndServe()
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, syscall.SIGINT)
	select {
	case sig := <-shutdownChannel:
		log.Printf("get signal -> %s", sig)
		logger.Slogger().Info("[GRACEFUL_SHUTDOWN_SPAWNED]", slog.Any("signal type", sig))
		timewait := 10 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timewait)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			server.Close()
		}
	case err := <-serverErr:
		failOnError(err, "server error")
	}


}


