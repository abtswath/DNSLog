package main

import (
	"dnslog/api"
	"dnslog/api/middleware"
	"dnslog/api/router"
	"dnslog/api/router/login"
	"dnslog/api/router/logs"
	"dnslog/configs"
	"dnslog/internal/dns"
	"dnslog/internal/store"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func runDaemon(opts configs.Option, db string) error {
	logrus.SetLevel(opts.LogLevel)
	store, err := store.NewDatabaseStore(db)
	if err != nil {
		return err
	}
	apiServer := api.New(opts.HTTPAddr)
	apiServer.Group("/api", []router.Router{
		login.New(opts.Username, opts.Password),
	})
	apiServer.Group("/api", []router.Router{
		logs.New(store),
	}, middleware.Auth(opts.Tokens))
	err = apiServer.Serve()
	if err != nil {
		return err
	}
	dnsServer := dns.NewServer(dns.Option{
		Addr:   opts.DNSAddr,
		Domain: opts.Domain,
		A:      opts.A,
	}, store)
	return dnsServer.Serve()
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	logrus.SetLevel(logrus.DebugLevel)

	logrus.SetOutput(os.Stdout)

	dir, err := initApp()
	if err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}

	cfg := configs.NewDaemonOptions("$HOME/")
	if err := runDaemon(cfg, path.Join(dir, "dnslog.db")); err != nil {
		logrus.Errorln(err)
		os.Exit(1)
	}
}
