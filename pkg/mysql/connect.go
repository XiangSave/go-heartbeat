package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBInfo struct {
	DBType             string
	Host               string
	UserName           string
	Password           string
	Charset            string
	DatabaseName       string
	MaxIdleConnections int
	MaxOpenConnections int
}

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

func NewDBModel(info *DBInfo) *DBModel {
	return &DBModel{DBInfo: info}
}

func (m *DBModel) Connect() error {
	var err error
	s := "%s:%s@tcp(%s)/%s?" +
		"charset=%s&parseTime=True&loc=Local&timeout=1000ms"
	dsn := fmt.Sprintf(
		s,
		m.DBInfo.UserName,
		m.DBInfo.Password,
		m.DBInfo.Host,
		m.DBInfo.DatabaseName,
		m.DBInfo.Charset,
	)

	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	m.DBEngine.SetMaxIdleConns(m.DBInfo.MaxIdleConnections)
	m.DBEngine.SetMaxOpenConns(m.DBInfo.MaxOpenConnections)
	if err != nil {
		return err
	}
	return nil
}
