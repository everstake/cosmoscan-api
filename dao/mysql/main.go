package mysql

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/everstake/cosmoscan-api/dao/derrors"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
	"os"
	"path/filepath"
	"time"
)

const migrationsDir = "./dao/mysql/migrations"

type DB struct {
	config config.Mysql
	db     *sqlx.DB
}

func NewDB(cfg config.Mysql) (*DB, error) {
	m := &DB{
		config: cfg,
	}
	m.tryOpenConnection()
	err := m.migrate()
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *DB) tryOpenConnection() {
	for {
		err := m.openConnection()
		if err != nil {
			log.Error("cant open connection to mysql: %s", err.Error())
		} else {
			log.Info("mysql connection success")
			return
		}
		time.Sleep(time.Second)
	}
}

func (m *DB) openConnection() error {
	source := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true&parseTime=true",
		m.config.User,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.DB,
	)
	var err error
	m.db, err = sqlx.Connect("mysql", source)
	if err != nil {
		return err
	}
	err = m.db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func (m DB) find(dest interface{}, sb squirrel.SelectBuilder) error {
	sql, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	err = m.db.Select(dest, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m DB) first(dest interface{}, sb squirrel.SelectBuilder) error {
	sql, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	err = m.db.Get(dest, sql, args...)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(derrors.ErrNotFound)
		}
		return err
	}
	return nil
}

func (m DB) insert(sb squirrel.InsertBuilder) (id uint64, err error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return id, err
	}
	result, err := m.db.Exec(sql, args...)
	if err != nil {
		mErr, ok := err.(*mysql.MySQLError)
		if ok && mErr.Number == 1062 {
			return 0, errors.New(derrors.ErrDuplicate)
		}
		return id, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return id, err
	}
	return uint64(lastID), nil
}

func (m DB) update(sb squirrel.UpdateBuilder) (err error) {
	sql, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}

func (m DB) migrate() error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	dir := filepath.Join(filepath.Dir(ex), migrationsDir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		dir = migrationsDir
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			return errors.New("Migrations dir does not exist: " + dir)
		}
	}
	migrations := &migrate.FileMigrationSource{
		Dir: dir,
	}
	_, err = migrate.Exec(m.db.DB, "mysql", migrations, migrate.Up)
	return err
}

func field(table string, column string, alias ...string) string {
	s := fmt.Sprintf("%s.%s", table, column)
	if len(alias) == 1 {
		return fmt.Sprintf("%s as %s", s, alias)
	}
	return s
}

func joiner(rightTable string, leftTable string, field string) string {
	return fmt.Sprintf("%s ON %s.%s = %s.%s", rightTable, leftTable, field, rightTable, field)
}
