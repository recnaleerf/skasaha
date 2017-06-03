package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KuroiKitsu/skasaha"
)

const (
	ConfigFile = "skasaha.json"
	IndexTTL   = 24 * time.Hour
)

func main() {
	var (
		err error

		config *Config
	)

	config, err = LoadConfigFile(ConfigFile)
	if err != nil {
		panic(err)
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	if config.Token == "" {
		logger.Fatal("missing token")
	}
	if config.Prefix == "" {
		logger.Fatal("missing prefix")
	}

	s := &skasaha.Skasaha{
		Token:    config.Token,
		Prefix:   config.Prefix,
		EmojiDir: config.EmojiDir,
		Logger:   logger,
	}

	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	err = s.Sync()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-time.After(IndexTTL):
				err = s.Sync()
				if err != nil {
					log.Print(err)
				}
			case <-interrupt:
				break
			}
		}
	}()

	<-interrupt

	err = s.Close()
	if err != nil {
		log.Fatal(err)
	}
}
