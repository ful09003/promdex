package main

import (
	"errors"
	"atamedomain.name/promdex/internal/metastorage"
	"github.com/urfave/cli/v2"
)

func initStore(c *cli.Context) error {
	if !c.Bool("accept-data-loss") {
		return errors.New("expilcit confirmation of data loss required. Re-run with --accept-data-loss")
	}

	switch c.String("store-type") {
	case metastorage.NilStore.String():
		//Do something special for NilStore
	case metastorage.SQLiteStore.String():
		if _, err := metastorage.NewSQLiteMetastore(c.String("store-path"), true); err != nil {
			return err
		}
	}
	return nil
}