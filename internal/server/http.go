package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rez1dent3/otus-final/internal/imgprev"
	"github.com/rez1dent3/otus-final/internal/server/handlers"
)

type HTTPServerInterface interface {
	ListenAndServe(context.Context) error
	HTTPHandler() http.Handler
}

type impl struct {
	app    imgprev.AppInterface
	server *http.Server

	previewer *handlers.PreviewHandler
}

func New(appImpl imgprev.AppInterface) HTTPServerInterface {
	return &impl{app: appImpl, previewer: handlers.NewPreviewer(appImpl)}
}

func (i *impl) ListenAndServe(ctx context.Context) error {
	httpHandler := i.HTTPHandler()

	i.server = &http.Server{
		Addr:              i.app.Config().Server.Addr,
		Handler:           i.middleware(httpHandler),
		ReadHeaderTimeout: time.Second,
	}

	go func() {
		<-ctx.Done()

		err := i.Stop(ctx)
		if err != nil {
			i.app.Logger().Error(err.Error())
			return
		}
	}()

	return i.server.ListenAndServe()
}

func (i *impl) Stop(ctx context.Context) error {
	if i.previewer != nil {
		i.previewer.Purge()
	}

	if i.server == nil {
		return nil
	}

	return i.server.Shutdown(ctx)
}

func (i *impl) HTTPHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", (&handlers.Health{}).Handle)
	mux.Handle("/fill/", i.previewer)

	return mux
}

func (i *impl) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		next.ServeHTTP(w, r)
		latency := time.Since(now)

		i.app.Logger().Info(fmt.Sprintf(
			"%s [%s] %s %s %s %d %s",
			r.Header.Get("X-Forwarded-For"),
			now.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			latency.Microseconds(),
			r.Header.Get("User-Agent"),
		))
	})
}
