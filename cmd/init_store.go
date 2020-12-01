package main

import (
	"errors"

	"atamedomain.name/promdex/internal"
	"github.com/urfave/cli/v2"
)

func initStore(c *cli.Context) error {
	if !c.Bool("accept-data-loss") {
		return errors.New("expilcit confirmation of data loss required. Re-run with --accept-data-loss")
	}

	switch c.String("store-type") {
	case internal.NilStore.String():
		//Do something special for NilStore
	case internal.SQLiteStore.String():
		if _, err := internal.NewSQLiteMetastore(c.String("store-path"), true); err != nil {
			return err
		}
	}
	return nil
}
