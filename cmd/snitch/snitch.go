package main // import "github.com/xjewer/snitch/cmd/snitch"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/quipo/statsd"
	"github.com/xjewer/snitch"
	"github.com/xjewer/snitch/lib/config"
	"github.com/xjewer/snitch/lib/simplelog"
	"github.com/xjewer/snitch/lib/stats"
)

var (
	wg         = sync.WaitGroup{}
	processors = make([]*snitch.Processor, 0)

	cfg             = flag.String("config", config.DefaultConfigPath, "config file name")
	statsdEndpoint  = flag.String("statsd", "", "statsd endpoint")
	statsdKeyPrefix = flag.String("prefix", "test", "statsd global key prefix")
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
		reader, err := snitch.NewFileReader(source, l)
		if err != nil {
			l.Println(err)
			continue
		}

		parser, err := snitch.NewParser(reader, s, source)
		if err != nil {
			l.Println(err)
			continue
		}

		p := snitch.NewProcessor(parser, reader, l)
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
