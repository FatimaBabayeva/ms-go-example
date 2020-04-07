package repo

import (
	"database/sql"
	"github.com/go-pg/pg"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
	"ms-go-example/properties"
	"strings"
)

var Db *pg.DB

func InitDb() {
	Db = pg.Connect(&pg.Options{
		Addr:     strings.Split(properties.Props.DbUrl, "/")[0],
		Database: strings.Split(properties.Props.DbUrl, "/")[1],
		User:     properties.Props.DbUser,
		Password: properties.Props.DbPass,
	})
}

func MigrateDb() error {
	log.Info("MigrateDb.start")

	connStr := properties.Props.DbConnStr() + "?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Info("Applied ", n, " migrations")
	log.Info("MigrateDb.end")
	return nil
}
