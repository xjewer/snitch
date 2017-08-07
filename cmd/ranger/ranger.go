package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/xjewer/ranger"
	"github.com/xjewer/ranger/lib/stats"
)

var (
	done = make(chan struct{})
	wg   = sync.WaitGroup{}

	file           = flag.String("file", "", "logs file name")
	noflow         = flag.Bool("noflow", false, "no flow file")
	offsetFile     = flag.String("offest", "", "file for preserving offset")
	statsdEndpoint = flag.String("statsd", "", "statsd endpoint")
	statsdPrefix   = flag.String("prefix", "balancer.%HOST%.", "statsd metrics prefix")
)

func main() {
	flag.Parse()

	s := stats.NewStatsd(*statsdEndpoint, *statsdPrefix, 0)
	err := s.CreateSocket()
	defer s.Close()

	if err != nil {
		panic(err)
	}

	p := ranger.NewParser(s)

	r, err := ranger.New(*file, *noflow, *offsetFile)
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go handle(r, p)

	c := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
	//signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT}
	signal.Notify(c, signals...)
	for {
		sig := <-c
		fmt.Printf("Got %q signal\n", sig)
		close(done)
		r.Close()
		break
	}

	wg.Wait()

}

func handle(r *ranger.Reader, p ranger.Parser) {
	errOutput := os.Stdout
	lines := r.GetLines()
	for {
		select {
		case l := <-lines:
			if l == nil {
				fmt.Fprintln(errOutput, "Empty line")
				continue
			}
			err := p.HandleLine(l.Text)
			if err != nil {
				fmt.Fprintln(errOutput, err)
			}
		case <-done:
			fmt.Println("Closing...")
			errOutput.Close()
			wg.Done()
			return
		}
	}
}
