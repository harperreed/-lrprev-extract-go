package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetOriginalFilePathSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT agfile.id_global as uuid, root.absolutePath, agfolder.pathFromRoot, agfile.baseName").
		WithArgs("test-uuid").
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "absolutePath", "pathFromRoot", "baseName"}).
			AddRow("test-uuid", "/absolute/path", "path/from/root", "test.jpg"))

	fullPath, baseName, err := getOriginalFilePath(db, "test-uuid")

	if err != nil {
		t.Errorf("error was not expected while getting original file path: %s", err)
	}

	expectedFullPath := "absolute/path/path/from/root"
	if fullPath != expectedFullPath {
		t.Errorf("expected full path %s, got %s", expectedFullPath, fullPath)
	}

	expectedBaseName := "test.jpg"
	if baseName != expectedBaseName {
		t.Errorf("expected base name %s, got %s", expectedBaseName, baseName)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetOriginalFilePathNoEntry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT agfile.id_global as uuid, root.absolutePath, agfolder.pathFromRoot, agfile.baseName").
		WithArgs("non-existent-uuid").
		WillReturnError(sql.ErrNoRows)

	_, _, err = getOriginalFilePath(db, "non-existent-uuid")

	if err == nil {
		t.Error("expected an error, but got none")
	}

	expectedError := "no entry found for UUID: non-existent-uuid"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetOriginalFilePathDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT agfile.id_global as uuid, root.absolutePath, agfolder.pathFromRoot, agfile.baseName").
		WithArgs("test-uuid").
		WillReturnError(errors.New("database connection error"))

	_, _, err = getOriginalFilePath(db, "test-uuid")

	if err == nil {
		t.Error("expected an error, but got none")
	}

	expectedError := "database query failed: database connection error"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetOriginalFilePathQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT agfile.id_global as uuid, root.absolutePath, agfolder.pathFromRoot, agfile.baseName").
		WithArgs("test-uuid").
		WillReturnError(errors.New("query execution error"))

	_, _, err = getOriginalFilePath(db, "test-uuid")

	if err == nil {
		t.Error("expected an error, but got none")
	}

	expectedError := "database query failed: query execution error"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

package database

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)