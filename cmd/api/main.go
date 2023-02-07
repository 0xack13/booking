package main

import (
	"context"
	"fmt"
	"github.com/xsolrac87/booking/api"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	// HTTP Server
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	timeout, _ := strconv.Atoi(os.Getenv("SERVER_TIMEOUT"))
	httpServer, err := api.NewHTTPServer(
		"",
		api.WithPort(port),
		api.WithTimeout(time.Duration(timeout)*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	// APi
	serverAPI, err := api.New(httpServer)
	if err != nil {
		log.Fatal(err)
	}

	// Context
	ctx, cancel := context.WithCancel(context.TODO())
	go applicationEnd(cancel)

	// Run server
	if err := serverAPI.RunAPI(ctx); err != nil {
		log.Fatal(err)
	}
}

func applicationEnd(cf context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	sig := <-sigChan
	fmt.Printf("context cancel form signal %v \n", sig)
	cf()
}
