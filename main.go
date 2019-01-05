package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gplume/gokit-error-handling/api"
	"github.com/gplume/gokit-error-handling/errs"
	"github.com/gplume/gokit-error-handling/middle"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitlog "github.com/go-kit/kit/log"
)

func main() {
	logger := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stderr)) // for use with sumologic
	// logger := kitlog.NewLogfmtLogger(os.Stderr) // preferable for local dev

	if err := errs.SetDefaults(errs.End, false); err != nil {
		logger.Log("init_error", err)
		os.Exit(1)
	}

	// SERVICE
	svc, err := api.NewStringService()
	if err != nil {
		logger.Log("init_error", err)
		os.Exit(1)
	}

	// ENDPOINTS
	e := api.Endpoints{
		Uppercase: api.MakeUppercaseEndpoint(svc),
		Count:     api.MakeCountEndpoint(svc),
	}

	r := api.MakeHTTPHandler(
		e,
		logger,
	)

	srv := &http.Server{
		Addr: ":8080",
		// It is good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      middle.RecoverFromPanic(logger, r),
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Debug
	go func() {
		debugAddr := ":8081"
		// mux := http.NewServeMux()
		http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
		http.Handle("/metrics", promhttp.Handler())
		panic(http.ListenAndServe(debugAddr, nil))
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	//  SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGQUIT)
	logger.Log("---ready---", "¡GO!")
	// Block until we receive our signal.
	<-c
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	go func() {
		srv.Shutdown(ctx)
	}()
	<-ctx.Done() // if your application should wait for other services
	logger.Log("shutting down", ctx)
	os.Exit(0)
}
