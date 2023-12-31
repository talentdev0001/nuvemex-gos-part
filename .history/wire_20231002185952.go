//+build wireinject

package part

import (
	"fmt"
	"github.com/Montrealist-cPunto/commons/config"
	"github.com/Montrealist-cPunto/commons/log"
	"github.com/Montrealist-cPunto/commons/queue"
	"github.com/Montrealist-cPunto/goseanto"
	"github.com/google/wire"
	"os"
	"sync"
)

var onceAppConfig sync.Once
var appConfig *config.Config

func MustConfig() *config.Config {
	onceAppConfig.Do(func() {
		appConfig = config.LoadFromDirectory("./resources/config",
			"goseanto.yml",
			fmt.Sprintf("goseanto-%s.yml", os.Getenv("app_env")),
			"config.yml",
			fmt.Sprintf("%s.yml", os.Getenv("app_env")),
		)
	})

	return appConfig
}

var onceHinterService sync.Once
var hinterService *Hinter

func provideHinterService(cfg *config.Config) *Hinter {
	panic(wire.Build(
		goseanto.MustElasticSearch,
		log.MustLogger,
		queue.MustQueue,
		goseanto.ProviderSuppliers,
		wire.Struct(new(Hinter), "*")))
}

func MustHinterService(cfg *config.Config) *Hinter {
	onceHinterService.Do(func() {
		hinterService = provideHinterService(cfg)
	})

	return hinterService
}

func MustSearchLambda(appConfig *config.Config) *SearchLambda {
	panic(wire.Build(
		goseanto.MustSearchService,
		log.MustLogger,
		wire.Struct(new(SearchLambda), "*")))
}

func MustHinterLambda(appConfig *config.Config) *HinterLambda {
	panic(wire.Build(
		MustHinterService,
		log.MustLogger,
		wire.Struct(new(HinterLambda), "*")))
}

func MustDetailsLambda(appConfig *config.Config) *DetailsLambda {
	panic(wire.Build(
		goseanto.MustSearchService,
		log.MustLogger,
		wire.Struct(new(DetailsLambda), "*")))
}
