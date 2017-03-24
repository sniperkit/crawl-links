package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/golang/glog"
	"github.com/k0kubun/pp"
	"gopkg.in/urfave/cli.v2"

	"github.com/crackcomm/crawl"
	"github.com/crackcomm/crawl/nsq/consumer"

	links "github.com/crackcomm/crawl-links/spider"
)

func init() {
	consumer.Flags = append(consumer.Flags, &cli.StringFlag{
		Name:   "output-topic",
		EnvVars: []string{"OUTPUT_TOPIC"},
		Usage:  "search results output nsq topic (required)",
		Value:  "links_results",
	})
}

func main() {
	defer glog.Flush()

	// CRAWL_DEBUG environment variable turns on debug mode
	// crawler then can spit out logs using glog.V(3)
	var verbosity string
	if yes, _ := strconv.ParseBool(os.Getenv("CRAWL_DEBUG")); yes {
		verbosity = "-v=3"
	}

	// We are setting glog to log to stderr
	flag.CommandLine.Parse([]string{"-logtostderr", verbosity})

	// Start consumer
	app := consumer.New(
		consumer.WithSpiderConstructor(func(app *consumer.App) consumer.Spider {
			// Get NSQ topic for results
			outputTopic := app.Ctx.String("output-topic")
			// Spider constructor
			return func(crawler crawl.Crawler) {
				// links spider
				spider := &links.Spider{
					Crawler: app.Crawler(),
					Output: func(result *links.Result) error {
						// Pretty print result to stdout
						pp.Print(result)
						// Publish result to NSQ on a given topic
						return app.Producer.PublishJSON(outputTopic, result)
					},
				}
				spider.Register()
			}
		}),
	)

	// Command line usage
	app.Name = "crawl-links"
	app.Usage = "links crawler"
	app.Version = "0.0.1"

	if err := app.Run(os.Args); err != nil {
		glog.Fatal(err)
	}
}
