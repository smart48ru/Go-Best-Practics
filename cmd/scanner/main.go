package main

// Исходники задания для первого занятия у других групп https://github.com/t0pep0/GB_best_go

import (
	"Best-GO/internal/scann"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Best-GO/internal/config"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := configuration.New()
	if cfg.Helper() {
		cfg.PrintHelp()
		os.Exit(1)
	}
	log.Info().Msg("Starting filescanner")
	log.Debug().Msgf("My id: %d", os.Getpid())
	wd, err := os.Getwd()
	if err != nil {
		msg := fmt.Sprintf("%s", err)
		log.Error().Msg(msg)
	}

	log.Trace().Msg("Make context WitchTimeOut")
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	log.Debug().Msgf("Make NewScanner extension = *%s, depth = %d", cfg.FileExt(), cfg.MaxDepth())
	sc := scann.New(ctx, wd, cfg.FileExt(), cfg.MaxDepth())

	log.Trace().Msg("Start FindFile in goroutine")
	go sc.FindFiles()

	log.Trace().Msg("Start ListChanel in goroutine")
	listenChannels(ctx, cancel, &sc, &cfg)
}

func listenChannels(ctx context.Context, cancel context.CancelFunc, s scann.Scanner, cfg configuration.Configuration) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	usr1Chan := make(chan os.Signal, 1)
	signal.Notify(usr1Chan, syscall.SIGUSR1)

	usr2Chan := make(chan os.Signal, 1)
	signal.Notify(usr2Chan, syscall.SIGUSR2)

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Done")
			return
		case err := <-s.ErrChan():
			log.Error().Msgf("%s", err)
		case result := <-s.ResChan():
			if cfg.JsonLog() {
				log.Info().Msgf("Name: %s Path: %s", result.Name(), result.Dir())
			} else {
				log.Info().Msgf("\tName: %s\t\t Path: %s", result.Name(), result.Dir())
			}
		case <-stopCh:
			log.Info().Msg("Принят сигнал к завершению программы")
			cancel()
		case <-usr1Chan:
			log.Info().Msgf("Текущая директория %s, Текущая глубина поиска = %v", s.CurDir(), s.Depth())
		case <-usr2Chan:
			s.IncDepth()
			log.Info().Msgf("Текущая глубина поиска изменена = %d", s.Depth())
		}
	}
}
