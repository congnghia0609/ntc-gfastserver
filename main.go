/**
 *
 * @author nghiatc
 * @since Oct 06, 2020
 */

package main

import (
	"log"
	"ntc-gfastserver/server"
	"os"
	"os/signal"
)

func main() {
	////// -------------------- Start WebServer -------------------- //////
	// StartWebServer
	go server.StartWebServer("webserver")

	// Hang thread Main.
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("################# End Main #################")
}
