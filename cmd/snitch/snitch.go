package main // import "github.com/xjewer/snitch/cmd/snitch"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/xjewer/snitch"
	"github.com/xjewer/snitch/lib/stats"
)

var (
	done = make(chan struct{})
	wg   = sync.WaitGroup{}

	file           = flag.String("file", "", "logs file name")
	noflow         = flag.Bool("noflow", false, "no flow file")
	mustExists     = flag.Bool("exists", false, "true, if file must exists")
	offsetFile     = flag.String("offest", "", "file for preserving offset")
	statsdEndpoint = flag.String("statsd", "", "statsd endpoint")
	statsdPrefix   = flag.String("prefix", "balancer.%HOST%.", "statsd metrics prefix")
	buffer         = flag.Int("buffer", 0, "statsd buffer interval")
)

func main() {
	flag.Parse()

	if *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	errOutput := os.Stderr
	log.SetOutput(errOutput)
	defer errOutput.Close()

	s := stats.NewStatsd(*statsdEndpoint, *statsdPrefix, *buffer)
	err := s.CreateSocket()
	defer s.Close()

	if err != nil {
		log.Fatal(err)
	}

	p := snitch.NewParser(s)

	r, err := snitch.New(*file, *noflow, *mustExists, *offsetFile)
	if err != nil {
		log.Fatal(err)
	}
	wg.Add(1)
	go handle(r, p)

	c := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
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

func handle(r *snitch.Reader, p snitch.Parser) {
	lines := r.GetLines()
	for {
		select {
		case l := <-lines:
			if l == nil {
				log.Println("Empty line")
				continue
			}
			err := p.HandleLine(l.Text)
			if err != nil {
				log.Println(err)
			}
		case <-done:
			fmt.Println("Closing...")
			wg.Done()
			return
		}
	}
}
