package database

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestSqlxSelect(t *testing.T) {
	type Snippet struct {
		Name string `db:"name"`
	}
	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockSnippets := []interface{}{
			Snippet{
				Name: "test1",
			},
			Snippet{
				Name: "test2",
			},
		}
		mockQuery := "^select name from snippet$"
		mockRows := NewMockRows(mock, mockSnippets...)
		mock.ExpectQuery(mockQuery).WillReturnRows(mockRows)

		snippets := []Snippet{}
		err := SqlxSelect(db, &snippets, "select name from snippet")

		assert.NoError(t, err)
		for index, mockSnippet := range mockSnippets {
			mockSnippet := mockSnippet.(Snippet)
			assert.Equal(t, mockSnippet, snippets[index])
		}
	})

	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockQuery := "^select name from snippet$"
		mockRows := NewMockRowsWithoutData(mock, &Snippet{})
		mock.ExpectQuery(mockQuery).WillReturnRows(mockRows)

		snippets := []Snippet{}
		err := SqlxSelect(db, &snippets, "select name from snippet")

		assert.NoError(t, err)
		assert.Equal(t, []Snippet{}, snippets)
	})

	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockQuery := "^select name from snippet$"
		mock.ExpectQuery(mockQuery).WillReturnError(fmt.Errorf("some error"))

		snippets := []Snippet{}
		err := SqlxSelect(db, &snippets, "select name from snippet")

		assert.Error(t, err)
		assert.NotEqual(t, sql.ErrNoRows, err)
	})
}

func TestSqlxGet(t *testing.T) {
	type Snippet struct {
		Name string `db:"name"`
	}
	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockSnippets := []interface{}{
			Snippet{
				Name: "test",
			},
		}
		mockQuery := "^select name from snippet where name=(.*)$"
		mockRows := NewMockRows(mock, mockSnippets...)
		mock.ExpectQuery(mockQuery).WithArgs("test").WillReturnRows(mockRows)

		snippet := Snippet{}
		err := SqlxGet(db, &snippet, "select name from snippet where name=?", "test")

		assert.NoError(t, err)
		mockSnippet := mockSnippets[0].(Snippet)
		assert.Equal(t, mockSnippet, snippet)
	})

	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockQuery := "^select name from snippet where name=(.*)$"
		mockRows := NewMockRowsWithoutData(mock, &Snippet{})
		mock.ExpectQuery(mockQuery).WithArgs("test").WillReturnRows(mockRows)

		snippets := Snippet{}
		err := SqlxGet(db, &snippets, "select name from snippet where name=?", "test")

		assert.Error(t, err)
		assert.Equal(t, sql.ErrNoRows, err)
	})

	RunWithMock(t, func(db *sqlx.DB, mock sqlmock.Sqlmock, t *testing.T) {
		mockQuery := "^select name from snippet where name=(.*)$"
		mock.ExpectQuery(mockQuery).WithArgs("test").WillReturnError(fmt.Errorf("some error"))

		snippets := Snippet{}
		err := SqlxGet(db, &snippets, "select name from snippet where name=?", "test")

		assert.Error(t, err)
		assert.NotEqual(t, sql.ErrNoRows, err)
	})
}
