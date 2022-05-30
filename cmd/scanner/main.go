package main

// Исходники задания для первого занятия у других групп https://github.com/t0pep0/GB_best_go

import (
	"Best-GO/internal/config"
	"Best-GO/internal/scann"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
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
		log.Error().Msgf("%s", err)
	}

	log.Debug().Msgf("Make NewScanner extension = *%s, depth = %d", cfg.FileExt(), cfg.MaxDepth())
	sc := scann.New(10, wd, cfg.FileExt(), cfg.MaxDepth())

	log.Trace().Msg("Start FindFile in goroutine")
	go sc.FindFiles()

	log.Trace().Msg("Start ListChanel in goroutine")
	listenChannels(&sc, &cfg)
}

func listenChannels(s scann.Scanner, cfg configuration.Configuration) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	usr1Chan := make(chan os.Signal, 1)
	signal.Notify(usr1Chan, syscall.SIGUSR1)

	usr2Chan := make(chan os.Signal, 1)
	signal.Notify(usr2Chan, syscall.SIGUSR2)

	for {
		select {
		case <-s.Ctx().Done():
			log.Info().Msg("Done")
			s.CtxCancel()
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
			s.CtxCancel()
		case <-usr1Chan:
			log.Info().Msgf("Текущая директория %s, Текущая глубина поиска = %v", s.CurDir(), s.Depth())
		case <-usr2Chan:
			s.IncDepth()
			log.Info().Msgf("Текущая глубина поиска изменена = %d", s.Depth())
		}
	}
}
