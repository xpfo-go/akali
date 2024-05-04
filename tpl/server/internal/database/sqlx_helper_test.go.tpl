package database

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func Test_bindArray(t *testing.T) {
	type Place struct {
		Country string `db:"country"`
		TelCode string `db:"telcode"`
	}
	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		query := "INSERT INTO place (country, telcode) VALUES (:country, :telcode)"
		places := []Place{{Country: "china", TelCode: "86"}, {Country: "us", TelCode: "001"}}

		q, args, err := bindArray(sqlx.BindType(db.DriverName()), query, places, db.Mapper)
		assert.NoError(t, err)
		assert.Equal(t, q, "INSERT INTO place (country, telcode) VALUES (?, ?),(?, ?)")
		assert.Equal(t, args, []interface{}{"china", "86", "us", "001"})
	})
}
