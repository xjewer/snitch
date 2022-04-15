package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Topface/snitch/internal"
	"github.com/Topface/snitch/internal/lib/config"
	"github.com/Topface/snitch/internal/lib/simplelog"
	"github.com/Topface/snitch/internal/lib/stats"
	"github.com/quipo/statsd"
)

var (
	wg         = sync.WaitGroup{}
	processors = make([]*internal.Processor, 0)

	cfg             = flag.String("config", config.DefaultConfigPath, "config file name")
	statsdEndpoint  = flag.String("statsd", "", "statsd endpoint")
	statsdKeyPrefix = flag.String("prefix", "", "statsd global key prefix")
	buffer          = flag.Int("buffer", 0, "statsd buffer interval")
)

func main() {
	flag.Parse()

	l := log.New(os.Stderr, "", log.LstdFlags)

	c, err := config.Parse(*cfg)
	if err != nil {
		l.Fatal(err)
	}

	s := stats.NewStatsd(*statsdEndpoint, *statsdKeyPrefix, *buffer)
	err = s.CreateSocket()
	defer s.Close()

	if err != nil {
		l.Fatal(err)
	}

	runProcessors(c, s, l)

	cs := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Notify(cs, signals...)
	for {
		sig := <-cs
		fmt.Printf("Got %q signal\n", sig)
		closeProcessors(l)
		break
	}

	wg.Wait()
}

// runProcessors run all processors
func runProcessors(c *config.Data, s statsd.Statsd, l simplelog.Logger) {
	for _, source := range c.Sources {
		reader, err := internal.NewFileReader(source, l)
		if err != nil {
			l.Println(err)
			continue
		}

		parser, err := internal.NewParser(reader, s, source, l)
		if err != nil {
			l.Println(err)
			continue
		}

		p := internal.NewProcessor(parser, reader, l)
		processors = append(processors, p)
		go p.Run()
	}
	wg.Add(len(processors))
}

// closeProcessors closes all established processors
func closeProcessors(l simplelog.Logger) {
	for _, p := range processors {
		err := p.Close()
		if err != nil {
			l.Println(err)
		}

		wg.Done()
	}
}
