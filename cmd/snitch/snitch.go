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
	"github.com/xjewer/snitch/lib/config"
	"github.com/quipo/statsd"
)

var (
	wg      = sync.WaitGroup{}
	parsers = make([]snitch.Parser, 0)

	cfg            = flag.String("config", config.DefaultConfigPath, "config file name")
	statsdEndpoint = flag.String("statsd", "", "statsd endpoint")
	statsdPrefix   = flag.String("prefix", "balancer.%HOST%.", "statsd metrics prefix")
	buffer         = flag.Int("buffer", 0, "statsd buffer interval")
)

func main() {
	flag.Parse()

	c, err := config.Parse(*cfg)
	if err != nil {
		log.Fatal(err)
	}

	errOutput := os.Stderr
	log.SetOutput(errOutput)
	defer errOutput.Close()

	s := stats.NewStatsd(*statsdEndpoint, *statsdPrefix, *buffer)
	err = s.CreateSocket()
	defer s.Close()

	if err != nil {
		log.Fatal(err)
	}

	runReaders(c, s)

	cs := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Notify(cs, signals...)
	for {
		sig := <-cs
		fmt.Printf("Got %q signal\n", sig)
		closeParsers()
		break
	}

	wg.Wait()
}

func runReaders(c *config.Data, s statsd.Statsd) {
	for _, source := range c.Sources {
		r, err := snitch.NewFileReader(source)
		if err != nil {
			log.Println(err)
			continue
		}

		p, err := snitch.NewParser(r, s, source)
		if err != nil {
			log.Println(err)
			continue
		}
		parsers = append(parsers, p)
		go p.Run()
	}
	wg.Add(len(parsers))
}

func closeParsers() {
	for _, p := range parsers {
		err := p.Close()
		if err != nil {
			log.Println(err)
		}

		wg.Done()
	}
}
