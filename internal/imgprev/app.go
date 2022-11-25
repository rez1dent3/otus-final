package imgprev

import (
	"os"
	"time"

	"github.com/rez1dent3/otus-final/internal/pkg/bus"
	"github.com/rez1dent3/otus-final/internal/pkg/bytesize"
	"github.com/rez1dent3/otus-final/internal/pkg/fetcher"
	"github.com/rez1dent3/otus-final/internal/pkg/fs"
	"github.com/rez1dent3/otus-final/internal/pkg/hsum"
	"github.com/rez1dent3/otus-final/internal/pkg/logger"
	"github.com/rez1dent3/otus-final/internal/pkg/lru"
	"github.com/rez1dent3/otus-final/internal/pkg/transformer"
	"github.com/rez1dent3/otus-final/internal/transport"
)

var supportedContentTypes = []string{
	"image/jpeg",
	"image/png",
}

type AppInterface interface {
	CommandBus() bus.CommandBusInterface
	Transform() transformer.TransformInterface
	Fetcher() fetcher.FetchInterface
	Logger() logger.LogInterface
	Config() *Config
	Purge()
}

type Config struct {
	Server struct {
		Addr string
	}

	Logger struct {
		Level string
	}

	Original struct {
		CacheDir    string `yaml:"cacheDir"`
		CachePrefix string `yaml:"cachePrefix"`
		CacheSize   string `yaml:"cacheSize"`
	}

	Preview struct {
		CacheDir    string `yaml:"cacheDir"`
		CachePrefix string `yaml:"cachePrefix"`
		CacheSize   string `yaml:"cacheSize"`
	}
}

type impl struct {
	fetch      fetcher.FetchInterface
	commandBus bus.CommandBusInterface
	log        logger.LogInterface
	config     *Config

	transform transformer.TransformInterface

	fetcherCache lru.CacheInterface
}

func New(config *Config) AppInterface {
	hash := hsum.New()
	commandBus := bus.NewSyncBus()
	log := logger.New(config.Logger.Level, os.Stdout)

	// fetcher
	fm := fs.New(config.Original.CacheDir, config.Original.CachePrefix)
	fetcherCache := lru.New(bytesize.Parse(config.Original.CacheSize), commandBus)
	fetcherTransport := transport.New(hash, fetcherCache, fm, log)

	// cleanup original images
	commandBus.Subscribe(lru.EventEvict, func(input any) {
		if val, ok := input.(transport.ResponseItem); ok {
			if err := fm.Delete(hash.HashByString(val.URL)); err != nil {
				log.Error(err.Error())
			}
		}
	})

	return &impl{
		fetch:        fetcher.NewHTTPFetcher(fetcherTransport, time.Second, supportedContentTypes),
		commandBus:   commandBus,
		log:          log,
		config:       config,
		fetcherCache: fetcherCache,
		transform:    transformer.NewStack(),
	}
}

func (i *impl) Fetcher() fetcher.FetchInterface {
	return i.fetch
}

func (i *impl) Transform() transformer.TransformInterface {
	return i.transform
}

func (i *impl) CommandBus() bus.CommandBusInterface {
	return i.commandBus
}

func (i *impl) Logger() logger.LogInterface {
	return i.log
}

func (i *impl) Config() *Config {
	return i.config
}

func (i *impl) Purge() {
	i.fetcherCache.Purge()
}
