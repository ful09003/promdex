package main

import (
	"net/http"
	"encoding/json"
	"atamedomain.name/promdex/internal"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/gorilla/mux"
)

func runPromdexServer(c *cli.Context) error {
	log.WithField("bind-addr", c.String("bind-addr")).Info("Spinning up a Promdex server...")
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

	r := mux.NewRouter()
	r.Use(logReq)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Yooo\n"))
	})
	r.HandleFunc("/m/{job}/{metric}", func(w http.ResponseWriter, r *http.Request){
		v := mux.Vars(r)
		d := store.RetrieveMetric(v["metric"], v["job"])
		b, e := json.Marshal(d)
		if e != nil {
			log.WithField("original", e).Warn("failed sending promdex server response")
		}
		w.Write(b)
	})

	log.Fatal(http.ListenAndServe(c.String("bind-addr"), r))
	return nil
}

func logReq(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		v := mux.Vars(r)

		log.WithFields(log.Fields{
			"request_path": r.RequestURI,
			"remote_addr": r.RemoteAddr,
			"request_method": r.Method,
			"job": v["job"],
			"metric": v["metric"],
		}).Debug("handling promdex request")

		next.ServeHTTP(w, r)
	})
}