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
	"github.com/xjewer/snitch/lib/stats"
)

var (
	wg         = sync.WaitGroup{}
	processors = make([]*snitch.Processor, 0)

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

	runProcessors(c, s)

	cs := make(chan os.Signal, 1)
	signals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGPIPE}
	signal.Notify(cs, signals...)
	for {
		sig := <-cs
		fmt.Printf("Got %q signal\n", sig)
		closeProcessors()
		break
	}

	wg.Wait()
}

// runProcessors run all processors
func runProcessors(c *config.Data, s statsd.Statsd) {
	for _, source := range c.Sources {
		reader, err := snitch.NewFileReader(source)
		if err != nil {
			log.Println(err)
			continue
		}

		parser, err := snitch.NewParser(reader, s, source)
		if err != nil {
			log.Println(err)
			continue
		}

		p := snitch.NewProcessor(parser, reader)
		processors = append(processors, p)
		go p.Run()
	}
	wg.Add(len(processors))
}

// closeProcessors closes all established processors
func closeProcessors() {
	for _, p := range processors {
		err := p.Close()
		if err != nil {
			log.Println(err)
		}

		wg.Done()
	}
}
