package main

import (
	"database/sql"
	"dwk/common"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var port = ":3001"
var (
	db_host     = common.GetEnv("DB_HOST", "localhost")
	db_user     = common.GetEnv("DB_USERNAME", "app")
	db_password = common.GetEnv("DB_PASSWORD", "example")
)
var dbUrl = fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", db_host, db_user, db_password)

func main() {
	db, err := sql.Open("postgres", dbUrl)
	common.CheckErr(err, "Failed to connect to the database")
	defer db.Close()

	err = initDb(db)
	common.CheckErr(err, "Failed initialize database")

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		count, err := updatePingCount(db)
		if err != nil {
			common.HandleErr(w, "Failed to update pings", http.StatusInternalServerError, err)
			return
		}
		fmt.Fprintf(w, "%d\n", count)
	})

	log.Printf("Server listening on port %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v\n", err)
	}
}

func initDb(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS pings (pinged_at TIMESTAMPTZ NOT NULL)`)
	return err
}

func updatePingCount(db *sql.DB) (int, error) {
	var count int
	err := db.QueryRow(`INSERT INTO pings (pinged_at) VALUES (now()) RETURNING (SELECT COUNT(*) + 1 FROM pings)`).Scan(&count)
	return count, err
}
