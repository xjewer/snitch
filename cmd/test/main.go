package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	testLineOk    = "[%s]	200	192.168.1.1:80	200	0.036	https	POST	/test	/test	OK	hostname"
	testLineError = "[%s]	200	192.168.1.1:80 : 127.0.0.1:8000	504 : 200	0.7 : 0.002	https	POST	/test	/test	Error	hostname"
)

var (
	file   = flag.String("file", "", "logs file name")
	events = flag.Int("events", 1, "event per second")
	done   = make(chan struct{})
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

	ticker := time.NewTicker(time.Second / time.Duration(*events))
	ticker2 := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			getLogLine(f)
		case <-ticker2.C:
			// truncate file every minute
			f.Truncate(0)
			f.Seek(0, 0)
		case <-done:
			ticker.Stop()
			return nil
		}
	}
}

func getLogLine(f *os.File) {
	if rand.Intn(2) == 0 {
		fmt.Fprintf(f, testLineOk+"\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"))
	} else {
		fmt.Fprintf(f, testLineError+"\n", time.Now().Format("02/Jan/2006:15:04:05 -0700"))
	}
}
