package main

//Исходники задания для первого занятия у других групп https://github.com/t0pep0/GB_best_go

import (
	"Best-GO/internal/scann"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	depth int64 = 2
)

func main() {
	fmt.Printf("My PID: %d\n", os.Getpid())
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	const wantExt = ".go"
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	sc := scann.New(ctx, wd, wantExt, depth)

	go sc.FindFiles()
	listenChannels(ctx, cancel, &sc)
}

func listenChannels(ctx context.Context, cancel context.CancelFunc, s scann.Scanner) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	usr1Chan := make(chan os.Signal, 1)
	signal.Notify(usr1Chan, syscall.SIGUSR1)

	usr2Chan := make(chan os.Signal, 1)
	signal.Notify(usr2Chan, syscall.SIGUSR2)

	for {
		select {
		case <-ctx.Done():
			log.Println("Done")
			return
		case err := <-s.ErrChan():
			log.Printf("Erroe %s\n", err)
		case result := <-s.ResChan():
			fmt.Printf("\tName: %s\t\t Path: %s\n", result.Name(), result.Dir())
		case <-stopCh:
			log.Printf("Принят сигнал к завершению программы\n")
			cancel()
		case <-usr1Chan:
			log.Printf("Текущая дириктория %s, Текущая глубина поиска %v\n", s.CurDir(), s.Depth())
		case <-usr2Chan:
			s.IncDepth()
			log.Printf("Текущая глубина поиска измениниа = %d\n", s.Depth())
		}
	}
}
