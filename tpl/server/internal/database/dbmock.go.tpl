package database

import (
	"database/sql/driver"
	"github.com/xpfo-go/logs"
	"reflect"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

// NewMockSqlxDB ...
func NewMockSqlxDB() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		logs.Error("An error was not expected when opening a stub database connection", err.Error())
	}
	sqlxDB := sqlx.NewDb(db, "mysql")
	sqlxDB.SetMaxOpenConns(10)
	return sqlxDB, mock
}

// RunWithMock ...
func RunWithMock(t *testing.T, test func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T)) {
	runner := func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		test(db, mock, t)
	}

	db, mock := NewMockSqlxDB()
	runner(db, mock, t)

	if err := mock.ExpectationsWereMet(); err != nil {
		logs.Error("there were unfulfilled expectations", err.Error())
	}
}

// NewMockRowsWithoutData ...
func NewMockRowsWithoutData(mock sqlmock.Sqlmock, arg interface{}) *sqlmock.Rows {
	var mockRows *sqlmock.Rows

	// 根据 Struct 的 db 标签，获取 columns
	objType := reflect.TypeOf(arg)
	if objType.Kind() == reflect.Ptr {
		objType = objType.Elem()
	}
	var columns []string
	for i := 0; i < objType.NumField(); i++ {
		dbTagName := objType.Field(i).Tag.Get("db")
		if dbTagName != "" {
			columns = append(columns, dbTagName)
		}
	}
	// log.Infof("columns len: %d, %+v", len(columns), columns)

	mockRows = sqlmock.NewRows(columns)
	return mockRows
}

// NewMockRows ...
func NewMockRows(mock sqlmock.Sqlmock, args ...interface{}) *sqlmock.Rows {
	mockRows := NewMockRowsWithoutData(mock, args[0])

	objType := reflect.TypeOf(args[0])

	// 获取数据并写入
	for _, obj := range args {
		objValue := reflect.ValueOf(obj)
		if objType.Kind() == reflect.Ptr {
			objValue = objValue.Elem()
		}
		values := []driver.Value{}
		for i := 0; i < objType.NumField(); i++ {
			dbTagName := objType.Field(i).Tag.Get("db")
			if dbTagName != "" {
				values = append(values, objValue.Field(i).Interface())
			}
		}
		// log.Infof("values len: %d, %+v", len(values), values)

		mockRows = mockRows.AddRow(values...)
	}
	return mockRows
}
