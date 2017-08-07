package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	//TestLine = "[$time_local] $status $upstream_addr $upstream_status $upstream_response_time $scheme $request_method $request_uri $uri $request_completion $host"
	TestLine = "[%s]	200	192.168.100.39:9000	200	0.036	https	POST	/vkadmin/queues_data/	/vkadmin/queues_data/	OK	topface.com"
)

var (
	file = flag.String("file", "", "logs file name")
	done = make(chan struct{})
)

func main() {
	flag.Parse()

	go writeToFile()

	c := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Notify(c, signals...)
	for {
		<-c
		close(done)
		return
	}
}

func writeToFile() error {
	f, err := os.OpenFile(*file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		return err
	}
	defer f.Close()

	ticker := time.NewTicker(time.Second)
	ticker2 := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			//RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
			fmt.Fprintf(f, TestLine+"\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"))
		case <-ticker2.C:
			// truncate file every minute
			f.Truncate(0)
			f.Seek(0, 0)
		case <-done:
			ticker.Stop()
			break
		}
	}

	return nil
}
