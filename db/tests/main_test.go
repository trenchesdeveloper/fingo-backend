package db_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	db "github.com/trenchesdeveloper/fingo/db/sqlc"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	conn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5433/fingo?sslmode=disable")

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = db.New(conn)
	

	os.Exit(m.Run())
}
