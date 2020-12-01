package main

import (
	"atamedomain.name/promdex/internal"
	"errors"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func addFlavor(c *cli.Context) error {
	log.Info("wtf")
	var store internal.Metastorer

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

	m := c.String("metric")
	log.WithField("metric", m).Debug("attempting to split input metric name")
	parts, err := splitMetric(m)

	if err != nil {
		return err
	}

	flavString := c.String("flavor")
	flavStruct := internal.PromdexFlavor{
		CtxString: flavString,
	}

	err = store.AddMetricFlavor(parts[0], parts[1], flavStruct)

	if err != nil {
		log.WithFields(log.Fields{
			"job":    parts[0],
			"metric": parts[1],
		}).Errorf("error during metric flavor addition, original: %s", err)
	}
	return nil
}

func splitMetric(in string) ([]string, error) {
	var ret []string
	ret = strings.Split(in, "/")
	log.WithField("split", ret).Info("split input metric name")

	if len(ret) == 1 {
		return ret, errors.New("input metric name was invalid")
	}

	return ret, nil
}
