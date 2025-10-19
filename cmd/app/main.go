package main

import (
	"flag"
	"github.com/aabbuukkaarr8/internal/apiserver"
	"github.com/aabbuukkaarr8/internal/storage"
	"github.com/wb-go/wbf/kafka"

	"github.com/BurntSushi/toml"
	"github.com/wb-go/wbf/zlog"
)

var (
	configPath string
)

func main() {

	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
	flag.Parse()
	zlog.Init()
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("config load error")
	}
	db := storage.New()
	err = db.Open(config.Store.DatabaseURL)
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("db open error")
		return
	}
	store, err := storage.NewMinio("minio:9000", "minioadmin", "minioadmin", "images")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("minio init error")
	}
	// Kafka
	prod := kafka.NewProducer("kafka:9092")
	cons := kafka.NewConsumer("kafka:9092", "image.uploaded", svc)
	go cons.Run(ctx)

}
