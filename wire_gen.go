// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package part

import (
	"fmt"
	"github.com/nuvemex/commons/config"
	"github.com/nuvemex/commons/log"
	"github.com/nuvemex/commons/queue"
	"github.com/nuvemex/goseanto"
	"os"
	"sync"
)

// Injectors from wire.go:

func provideHinterService(cfg *config.Config) *Hinter {
	elasticSearch := goseanto.MustElasticSearch(cfg)
	queueQueue := queue.MustQueue(cfg)
	v := goseanto.ProviderSuppliers(cfg)
	logger := log.MustLogger(cfg)
	hinter := &Hinter{
		ElasticService: elasticSearch,
		Queue:          queueQueue,
		Suppliers:      v,
		Logger:         logger,
	}
	return hinter
}

func MustSearchLambda(appConfig2 *config.Config) *SearchLambda {
	searchService := goseanto.MustSearchService(appConfig2)
	logger := log.MustLogger(appConfig2)
	searchLambda := &SearchLambda{
		Service: searchService,
		Logger:  logger,
	}
	return searchLambda
}

func MustHinterLambda(appConfig2 *config.Config) *HinterLambda {
	hinter := MustHinterService(appConfig2)
	logger := log.MustLogger(appConfig2)
	hinterLambda := &HinterLambda{
		Service: hinter,
		Logger:  logger,
	}
	return hinterLambda
}

func MustDetailsLambda(appConfig2 *config.Config) *DetailsLambda {
	searchService := goseanto.MustSearchService(appConfig2)
	logger := log.MustLogger(appConfig2)
	detailsLambda := &DetailsLambda{
		Service: searchService,
		Logger:  logger,
	}
	return detailsLambda
}

// wire.go:

var onceAppConfig sync.Once

var appConfig *config.Config

func MustConfig() *config.Config {
	onceAppConfig.Do(func() {
		appConfig = config.LoadFromDirectory("./resources/config",
			"goseanto.yml", fmt.Sprintf("goseanto-%s.yml", os.Getenv("app_env")), "config.yml", fmt.Sprintf("%s.yml", os.Getenv("app_env")),
		)
	})

	return appConfig
}

var onceHinterService sync.Once

var hinterService *Hinter

func MustHinterService(cfg *config.Config) *Hinter {
	onceHinterService.Do(func() {
		hinterService = provideHinterService(cfg)
	})

	return hinterService
}
