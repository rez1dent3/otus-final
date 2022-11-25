package usecases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rez1dent3/otus-final/internal/pkg/fetcher"
	"github.com/rez1dent3/otus-final/internal/pkg/fs"
	"github.com/rez1dent3/otus-final/internal/pkg/hsum"
	"github.com/rez1dent3/otus-final/internal/pkg/lru"
	"github.com/rez1dent3/otus-final/internal/pkg/transformer"
)

type PreviewUseCaseInterface interface {
	FillCenter(
		ctx context.Context,
		originalURL string,
		width int,
		height int,
		header http.Header,
	) ([]byte, error)
}

type PreviewItem struct {
	Key  string
	size uint64
}

func (p *PreviewItem) Size() uint64 {
	return p.size
}

func New(
	fm fs.FileInterface,
	hash hsum.HashInterface,
	cache lru.CacheInterface,
	transform transformer.TransformInterface,
	fetch fetcher.FetchInterface,
) PreviewUseCaseInterface {
	return &impl{fm: fm, hash: hash, cache: cache, transform: transform, fetch: fetch}
}

type impl struct {
	fm        fs.FileInterface
	hash      hsum.HashInterface
	cache     lru.CacheInterface
	fetch     fetcher.FetchInterface
	transform transformer.TransformInterface
}

func (i *impl) cacheKey(originalURL string, width int, height int) string {
	return i.hash.HashByString(fmt.Sprintf("fill:%s:%d:%d", originalURL, width, height))
}

func (i *impl) FillCenter(
	ctx context.Context,
	originalURL string,
	width int,
	height int,
	header http.Header,
) ([]byte, error) {
	cacheKey := i.cacheKey(originalURL, width, height)
	if _, ok := i.cache.Get(cacheKey); ok {
		if body, err := i.fm.Content(cacheKey); err == nil {
			return body, nil
		}
	}

	source, err := i.fetch.Get(ctx, originalURL, header)
	if err != nil {
		return nil, err
	}

	resp, err := i.transform.FillCenter(source, width, height)
	if err != nil {
		return nil, err
	}

	if i.cache.Put(cacheKey, &PreviewItem{Key: cacheKey, size: uint64(len(resp))}) {
		err = i.fm.Create(cacheKey, resp)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
