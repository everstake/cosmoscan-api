package clickhouse

import (
	"database/sql"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/golang-migrate/migrate/v4"
	goclickhouse "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mailru/go-clickhouse"
	"strings"
)

const migrationsPath = "./dao/clickhouse/migrations"

type DB struct {
	conn *sqlx.DB
}

func NewDB(cfg config.Clickhouse) (*DB, error) {
	conn, err := sql.Open("clickhouse", makeSource(cfg))
	if err != nil {
		return nil, fmt.Errorf("can`t make connection: %s", err.Error())
	}
	//err = makeMigration(conn, migrationsPath, cfg.Database)
	//if err != nil {
	//	return nil, fmt.Errorf("can`t make makeMigration: %s", err.Error())
	//}
	return &DB{
		conn: sqlx.NewDb(conn, "clickhouse"),
	}, nil
}

func (db *DB) Find(dest interface{}, b squirrel.SelectBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	err = db.conn.Select(dest, q, params...)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) FindFirst(dest interface{}, b squirrel.SelectBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	err = db.conn.Get(dest, q, params...)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Insert(b squirrel.InsertBuilder) error {
	q, params, err := b.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.Exec(q, params...)
	if err != nil {
		return err
	}
	return nil
}

func makeSource(cfg config.Clickhouse) string {
	return fmt.Sprintf("%s://%s:%d/%s?password=%s&user=%s",
		strings.Trim(cfg.Protocol, "://"),
		strings.Trim(cfg.Host, "/"),
		cfg.Port,
		cfg.Database,
		cfg.Password,
		cfg.User,
	)
}

func makeMigration(conn *sql.DB, migrationDir string, dbName string) error {
	driver, err := goclickhouse.WithInstance(conn, &goclickhouse.Config{})
	if err != nil {
		return fmt.Errorf("clickhouse.WithInstance: %s", err.Error())
	}
	mg, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		dbName, driver)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %s", err.Error())
	}
	if err := mg.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		}
	}
	return nil
}
