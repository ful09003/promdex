package main

import (
	"atamedomain.name/promdex/internal"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func loadSingleTarget(c *cli.Context) error {
	var store internal.Metastorer

	promdexCli, err := internal.NewPromdexTarget(c.String("target"))
	if err != nil {
		log.WithField("target", c.String("target")).Errorf("failed discovery, original: %s", err)
		return err
	}

	switch c.String("store-type") {
	case internal.NilStore.String():
		//Do something special for NilStore
		store = internal.NullMetastore{}
	case internal.SQLiteStore.String():
		s, err := internal.NewSQLiteMetastore(c.String("store-path"), false)
		if err != nil {
			log.WithField("storePath", c.String("store-path")).Errorf("failed init of sqlite store, original: %s", err)
			return err
		}

		store = s
	}

	var gs = internal.PromdexDiscoverStorer{
		TargetDiscoverer: promdexCli,
		Metastorer:       store,
	}

	actioned, err := gs.LoadIn()

	if err != nil {
		log.WithField("ingested", actioned).Errorf("failed loading metrics, original: %s", err)
	}

	log.WithField("ingested", actioned).Info("loaded metadata in")
	return nil
}
