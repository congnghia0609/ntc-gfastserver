/**
 *
 * @author nghiatc
 * @since Oct 06, 2020
 */

package main

import (
	"fmt"
	"github.com/congnghia0609/ntc-gfastserver/server"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/natefinch/lumberjack"
)

// initNConf init file config
func initNConf() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

// https://github.com/natefinch/lumberjack
func initLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/data/log/ntc-gfastserver/ntc-gfastserver.log",
		MaxSize:    10,   // 10 megabytes. Defaults to 100 MB.
		MaxBackups: 3,    // maximum number of old log files to retain.
		MaxAge:     28,   // maximum number of days to retain old log files
		Compress:   true, // disabled by default
	})
}

// increaseLimit increase resources limitations: ulimit -aH
func increaseLimit() {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	rlimit.Cur = rlimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	log.Printf("rlimit.Max = %d\n", rlimit.Max)
	log.Printf("rlimit.Cur = %d\n", rlimit.Cur)
}

func main() {
	////// -------------------- Init System -------------------- //////
	// Increase resources limitations
	increaseLimit()

	// Init NConf
	initNConf()

	//// Init Logger
	if "development" != nconf.GetEnv() {
		log.Printf("============== LogFile: /data/log/ntc-gfastserver/ntc-gfastserver.log")
		initLogger()
	}

	// Enable pprof hooks
	go func() {
		// http://localhost:6060/debug/pprof/
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

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

// func main() {
// 	////// -------------------- Init System -------------------- //////
// 	// Init NConf
// 	initNConf()

// 	//// Init Logger
// 	if "development" != nconf.GetEnv() {
// 		log.Printf("============== LogFile: /data/log/ntc-gfastserver/ntc-gfastserver.log")
// 		initLogger()
// 	}

// 	// Code test here.
// 	// 1. Insert
// 	// p := post.Post{
// 	// 	ID:        1,
// 	// 	Title:     "title1",
// 	// 	Body:      "body1",
// 	// 	CreatedAt: time.Now(),
// 	// 	UpdatedAt: time.Now(),
// 	// }
// 	// err := post.InsertPost(p)
// 	// fmt.Println("err:", err)
// 	// bp, err := json.Marshal(p)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(string(bp))

// 	// 2. Get post
// 	// p := post.GetPost(3)
// 	// bp, err := json.Marshal(p)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(string(bp))

// 	// 3. Get all post
// 	// p := post.GetAllPost()
// 	// bp, err := json.Marshal(p)
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(string(bp))

// 	// 4. gen id
// 	id, _ := mdb.Next("aaa")
// 	log.Println("id gen:", id)

// 	// 5. Benchmark nid gen
// 	// wg := new(sync.WaitGroup)
// 	// wg.Add(100)
// 	// for i := 0; i < 100; i++ {
// 	// 	go func() {
// 	// 		for j := 0; j < 100; j++ {
// 	// 			n, err := mdb.Next("benchmark")
// 	// 			if err != nil {
// 	// 				panic(err)
// 	// 			}
// 	// 			fmt.Println(n)
// 	// 		}
// 	// 		wg.Done()
// 	// 		fmt.Println("Done!!!")
// 	// 	}()
// 	// }
// 	// wg.Wait()

// 	// 6. Reset id gen
// 	// count, _ := mdb.ResetID("aaa", 0)
// 	// log.Println("resetid count:", count)

// 	log.Println("################# End Main #################")
// }
