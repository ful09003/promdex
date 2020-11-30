package main

import (
	"time"
	"os"

	"github.com/urfave/cli/v2"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	app := &cli.App{
		Name: "promdex-server",
		Version: "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name: "Michael Fuller",
				Email: "mfuller@digitalocean.com",
			},
		},
		Usage: "Promdex Server binary",
		Commands: []*cli.Command{
			&cli.Command{
				Name: "storage",
				Description: "Promdex storage-related functions. Beware -- these functions may result in Promdex data loss. Proceed with caution and review subcommand help.",
				Subcommands: []*cli.Command{
					{
						Name: "init",
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

/*	promdexCli, err := internal.NewPromdexTarget("http://localhost:9090")
	if err != nil {
		panic(err)
	}
	m, e := promdexCli.QueryTargetMetadata()
	if e != nil {
		panic(e)
	}

	for k, m1 := range m {
		fmt.Printf("%s, %+v\n", k, m1[0])
	}*/
}
