package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gplume/gokit-error-handling/api"

	kitlog "github.com/go-kit/kit/log"
)

func main() {
	logger := kitlog.NewLogfmtLogger(os.Stderr)

	// SERVICES
	svc := api.StringService{}

	// svc = loggingMiddleware{logger, svc}
	// svc = api.Instrumenting{requestCount, requestLatency, countResult, svc}

	// ENDPOINTS
	e := api.Endpoints{
		Uppercase: api.MakeUppercaseEndpoint(svc),
		Count:     api.MakeCountEndpoint(svc),
	}

	// TRANSPORT
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
		// Handler:      middle.RecoverFromPanic(logger, api),
		Handler: r,
	}

	// Run our server in a goroutine so that it doesn't block.

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Log(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	logger.Log("---ready---", "¡GO!")
	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	logger.Log("shutting down")
	os.Exit(0)
}
