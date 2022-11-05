package handlers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"

	"github.com/rez1dent3/otus-final/internal/imgprev"
	"github.com/rez1dent3/otus-final/internal/usecases"
	"github.com/rez1dent3/otus-final/pkg/bytesize"
	"github.com/rez1dent3/otus-final/pkg/fs"
	"github.com/rez1dent3/otus-final/pkg/hsum"
	"github.com/rez1dent3/otus-final/pkg/lru"
)

var ErrParseURL = errors.New("can't parse URL")

type PreviewHandler struct {
	app     imgprev.AppInterface
	useCase usecases.PreviewUseCaseInterface
	cache   lru.CacheInterface
}

func NewPreviewer(app imgprev.AppInterface) *PreviewHandler {
	hash := hsum.New()
	config := app.Config()
	commandBus := app.CommandBus()

	fm := fs.New(config.Preview.CacheDir, config.Preview.CachePrefix)
	previewerCache := lru.New(bytesize.Parse(config.Preview.CacheSize), commandBus)

	commandBus.Subscribe(lru.EventEvict, func(input any) {
		if val, ok := input.(usecases.PreviewItem); ok {
			if err := fm.Delete(val.Key); err != nil {
				app.Logger().Error(err.Error())
			}
		}
	})

	useCase := usecases.New(
		fm,
		hash,
		previewerCache,
		app.Transform(),
		app.Fetcher(),
	)

	return &PreviewHandler{app: app, useCase: useCase, cache: previewerCache}
}

func (p *PreviewHandler) PreviewerFillHandle(
	originalURL string,
	width int,
	height int,
	w http.ResponseWriter,
	r *http.Request,
) {
	resp, err := p.useCase.FillCenter(r.Context(), originalURL, width, height, r.Header)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	w.Header().Add("Content-Type", http.DetectContentType(resp))
	_, _ = w.Write(resp)
}

var reFillRoute = regexp.MustCompile(`^\/fill\/(\d+)\/(\d+)\/(.+)$`)

func (p *PreviewHandler) ParseURL(r *http.Request) (string, int, int, error) {
	results := reFillRoute.FindStringSubmatch(r.URL.Path)
	if len(results) != 4 {
		return "", 0, 0, ErrParseURL
	}

	width, err := strconv.Atoi(results[1])
	if err != nil {
		return "", 0, 0, err
	}

	height, err := strconv.Atoi(results[2])
	if err != nil {
		return "", 0, 0, err
	}

	return results[3], width, height, nil
}

func (p *PreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		p.app.Logger().Info("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte("Method not allowed"))
		return
	}

	originalURL, width, height, err := p.ParseURL(r)
	if err != nil {
		p.app.Logger().Info(err.Error())
		http.NotFound(w, r)
		return
	}

	p.PreviewerFillHandle(originalURL, width, height, w, r)
}

func (p *PreviewHandler) Purge() {
	p.cache.Purge()
}
