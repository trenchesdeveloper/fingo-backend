package db_test

import (
	"database/sql"
	"github.com/trenchesdeveloper/fingo-backend/utils"
	"log"
	"os"
	"testing"

	db "github.com/trenchesdeveloper/fingo-backend/db/sqlc"

	_ "github.com/lib/pq"
)

var testQueries *db.Queries

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	conn, err := sql.Open(config.DBdriver, config.DB_source)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = db.New(conn)

	os.Exit(m.Run())
}
