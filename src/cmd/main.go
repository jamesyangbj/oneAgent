package main

import (
	"flag"
	"fmt"
	"oneProxy/src/server"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var host *string = flag.String("l", "0.0.0.0", "proxy listen address")
var port *int = flag.Int("port", 9527, "proxy listen port")

func main() {
	flag.Parse()
	fmt.Println(fmt.Sprintf("start oneProxy server... listen %s:%d\n", *host, *port))
	runtime.GOMAXPROCS(runtime.NumCPU())
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGTERM|syscall.SIGQUIT)

	go func() {
		s := <-sc
		switch s {
		case syscall.SIGTERM | syscall.SIGQUIT:
			fmt.Println("receive signal SIGTERM, need to close the server. ")
			os.Exit(1)
			break
		default:
			fmt.Println("default signal handler")
		}
	}()

	proxy := new(server.ProxyServer)
	proxy.Start()

}
