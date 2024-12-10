package main

import (
	"flag"
	"log"
	"os"
	"sync"

	"github.com/KirantiLoh/ssh-automator/internal/commands"
	"github.com/KirantiLoh/ssh-automator/internal/model"
	"github.com/KirantiLoh/ssh-automator/internal/parser"
)

func main() {
	args := os.Args[1:]
	showHelp := flag.Bool("help", false, "Show help")

	fileName := flag.String("config", "", "Path to config file (.json)")

	flag.Parse()

	if *showHelp || len(args) <= 0 || *fileName == "" {
		commands.ShowHelp()
		return
	}

	config, err := parser.ParseConfigFile(*fileName)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	var wg sync.WaitGroup

	limiter := make(chan struct{}, 4)

	updates := make(chan model.Update)

	for _, server := range config.Servers {
		wg.Add(1)
		limiter <- struct{}{}
		go func() {
			commands.RunCommandsSSH(server, config.DefaultConfig, updates, &wg)
			<-limiter
		}()
	}

	go func() {
		for update := range updates {
			if update.IsComplete {
				log.Printf("[DONE] %s: %s\n", update.Host, update.Message)
			} else if update.IsError {
				log.Printf("[ERROR] %s: %s\n", update.Host, update.Message)
			} else {
				log.Printf("[INFO] %s: %s\n", update.Host, update.Message)
			}
		}
	}()

	wg.Wait()
	close(updates)
	close(limiter)

}
