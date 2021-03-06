package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/wikidistance/wikidist/pkg/crawler"
	"github.com/wikidistance/wikidist/pkg/db"
)

func main() {

	args := os.Args[1:]

	log.Println(args)

	if len(args) != 3 {
		log.Println("Usage: crawler <prefix> <startTitle> <nWorkers>")
		return
	}

	nWorkers, err := strconv.Atoi(args[2])

	if err != nil {
		log.Println("nWorkers should be an integer")
		return
	}

	client, _ := db.NewDGraph()
	c := crawler.NewCrawler(nWorkers, args[0], args[1], client)

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	c.Start()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("exiting")
}
