package main

import (
	"context"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/aabbuukkaarr8/internal/apiserver"
	"github.com/aabbuukkaarr8/internal/config"
	"github.com/aabbuukkaarr8/internal/handler"
	"github.com/aabbuukkaarr8/internal/kafka"
	"github.com/aabbuukkaarr8/internal/machine"
	"github.com/aabbuukkaarr8/internal/repository"
	"github.com/aabbuukkaarr8/internal/service"
	"github.com/aabbuukkaarr8/internal/storage"
	"github.com/aabbuukkaarr8/internal/storage/minio"
	image "github.com/aabbuukkaarr8/internal/upload"
	"github.com/wb-go/wbf/retry"
	"github.com/wb-go/wbf/zlog"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	configPath string
)

func main() {

	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
	flag.Parse()
	zlog.Init()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	config, err := config.LoadConfig(configPath)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load config file")
	}
	_, err = toml.DecodeFile(configPath, config)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("config load error")
	}
	db := storage.New()
	err = db.Open(config.DB.DSN())
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("db open error")
		return
	}
	strategy := retry.Strategy{
		Attempts: 2,
		Delay:    300 * time.Millisecond,
		Backoff:  2.0,
	}
	store, err := minio.NewMinio("minio:9000", "minioadmin", "minioadmin", "images")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("minio init error")
	}
	//repo
	repo := repository.NewRepository(db)
	//kafka producer
	prod := kafka.NewProducer(&config.Kafka, strategy)
	//image machine
	imgM := machine.New(store) //service
	service := service.NewService(store, prod, imgM, repo)

	//handle
	handle := image.NewdHandler(service)

	//handler route
	handler := handler.NewHandler(service)

	//kafka consumer
	c := kafka.NewConsumer(&config.Kafka, strategy, handle)

	//starting kafka
	var wg sync.WaitGroup
	wg.Add(1)
	go c.Consume(ctx, &wg)
	//HTTP
	s := apiserver.New(config)
	s.ConfigureRouter(handler)
	if err = s.Run(); err != nil {
		panic(err)
	}

	//context cancel
	<-ctx.Done()
	zlog.Logger.Info().Msg("context done")

	wg.Wait()

}
