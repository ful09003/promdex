package main

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	app := &cli.App{
		Name:     "promdex-server",
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Michael Fuller",
				Email: "mfuller@digitalocean.com",
			},
		},
		Usage: "Promdex Server binary",
		Commands: []*cli.Command{
			{
				Name:        "storage",
				Description: "Promdex storage-related functions. Beware -- these functions may result in Promdex data loss. Proceed with caution and review subcommand help.",
				Subcommands: []*cli.Command{
					{
						Name:        "init",
						Description: "(re-)generates a Promdex data store. DANGEROUS -- this will lead to Promdex data loss. Use this only on a new installation or when comfortable with losing data.",
						Flags: []cli.Flag{
							&cli.BoolFlag{Name: "accept-data-loss", Required: true, Usage: "User confirmation that destructive actions are acceptable"},
							&cli.StringFlag{Name: "store-type", Aliases: []string{"t"}, Required: true, Usage: "Used to signify what type of storage backend will be used", Value: "null"},
							&cli.StringFlag{Name: "store-path", Aliases: []string{"sp"}, Usage: "For path-based storage backends (e.g. sqlite), location to use for the backing store"},
						},
						Action: initStore,
					},
				},
			},
			{
				Name:        "target",
				Description: "Promdex target-related functions. Individually probably not very useful",
				Subcommands: []*cli.Command{
					{
						Name:        "load",
						Description: "Takes a singular target and a singular backing store, then loads target metadata into the store.",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "target", Required: true, Usage: "Prometheus endpoint to scrape from"},
							&cli.StringFlag{Name: "store-type", Aliases: []string{"t"}, Required: true, Usage: "Used to signify what type of storage backend will be used", Value: "null"},
							&cli.StringFlag{Name: "store-path", Aliases: []string{"sp"}, Usage: "For path-based storage backends (e.g. sqlite), location to use for the backing store"},
						},
						Action: loadSingleTarget,
					},
					{
						Name:        "add-flavor",
						Description: "Adds flavor (user context) to a single metric.",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "store-type", Aliases: []string{"t"}, Required: true, Usage: "Used to signify what type of storage backend will be used", Value: "null"},
							&cli.StringFlag{Name: "store-path", Aliases: []string{"sp"}, Usage: "For path-based storage backends (e.g. sqlite), location to use for the backing store"},
							&cli.StringFlag{Name: "metric", Usage: "Forward-slash-separated {job/metric} pair e.g. 'node-exporter/node_load5'", Required: true},
							&cli.StringFlag{Name: "flavor", Usage: "Flavor text to add to the metric", Required: true},
						},
						Action: addFlavor,
					},
				},
			},
			{
				Name: "server",
				Description: "Promdex server-related commands",
				Subcommands: []*cli.Command{
					{
						Name: "start",
						Description: "Starts a Promdex server",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "bind-addr", Required: true, Usage: "Address to bind Promdex server to", Value: ":19918"},
							&cli.StringFlag{Name: "store-type", Aliases: []string{"t"}, Required: true, Usage: "Used to signify what type of storage backend will be used", Value: "null"},
							&cli.StringFlag{Name: "store-path", Aliases: []string{"sp"}, Usage: "For path-based storage backends (e.g. sqlite), location to use for the backing store"},
						},
						Action: runPromdexServer,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
