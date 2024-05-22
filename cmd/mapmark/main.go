package main

import (
	"echoFramework/internal/app/mapmark"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := mapmark.New()
	app.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
}
