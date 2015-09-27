package main

import (
	"fmt"
	"os"

	"github.com/arussellsaw/workbench/runner"
	"github.com/arussellsaw/workbench/server"
	"github.com/howeyc/fsnotify"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("workbench")

func main() {
	do := make(chan bool)
	results := make(chan runner.ResultSet)
	go watcher(do)
	go runner.Start(do, results)

	var api = server.APIServer{}
	go api.Run()

	for {
		select {
		case rs := <-results:
			log.Debug(fmt.Sprintf("got results lol %v", rs))
			api.Lock()
			api.Results = append(api.Results, rs)
			api.Unlock()
		}
	}
}

func watcher(do chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	do <- true

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsModify() || ev.IsCreate() || ev.IsDelete() {
					if ev.IsAttrib() == false {
						if len(do) == 0 {
							do <- true
						}
					}
				}
			case err := <-watcher.Error:
				log.Debug("error:", err)
			}
		}
	}()

	dir, err := os.Getwd()
	if err != nil {
		log.Error("failed to get watcher dir: ", err)
	}
	err = watcher.Watch(dir)
	if err != nil {
		log.Fatal(err)
	}

	<-done

	watcher.Close()
}
