package main

import (
	"database/sql"
	"os"
	"path"

	_ "github.com/mattn/go-sqlite3"
)

var recordTableStmt = `
CREATE TABLE IF NOT EXISTS "records" (
	"id" integer NOT NULL,
	"domain" TEXT NOT NULL DEFAULT '',
	"ip" TEXT NOT NULL DEFAULT '',
	"created_at" text NOT NULL DEFAULT '',
	PRIMARY KEY ("id")
);
CREATE INDEX IF NOT EXISTS "domain_ip" ON "records" (
  	"domain" ASC,
  	"ip" ASC
);
CREATE INDEX IF NOT EXISTS "ip" ON "records" (
  	"ip" ASC
);
`

func initApp() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	appDir := path.Join(homeDir, ".dnslog")
	_, err = os.ReadDir(appDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err := os.Mkdir(appDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	if err = createDatabase(appDir); err != nil {
		return "", err
	}

	return appDir, nil
}

func createDatabase(appDir string) error {
	dbPath := path.Join(appDir, "dnslog.db")
	db, err := os.Open(dbPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		db, err = os.Create(dbPath)
		if err != nil {
			return err
		}
	}
	defer db.Close()

	return createTable(dbPath)
}

func createTable(dbPath string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	_, err = db.Exec(recordTableStmt)
	return err
}
