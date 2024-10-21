package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func OpenDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return db, nil
}

func GetOriginalFilePath(dbPath, uuid string) (string, string, error) {
	db, err := OpenDatabase(dbPath)
	if err != nil {
		return "", "", err
	}
	defer db.Close()

	return getOriginalFilePath(db, uuid)
}

func getOriginalFilePath(db *sql.DB, uuid string) (string, string, error) {
	query := `
		SELECT agfile.id_global as uuid, root.absolutePath, agfolder.pathFromRoot, agfile.baseName
		FROM AgLibraryFile agfile
		INNER JOIN AgLibraryFolder agfolder ON agfolder.id_local = agfile.folder
		INNER JOIN AgLibraryRootFolder root ON root.id_local = agfolder.rootFolder
		WHERE agfile.id_global = ?
	`

	var absolutePath, pathFromRoot, baseName string
	err := db.QueryRow(query, uuid).Scan(&uuid, &absolutePath, &pathFromRoot, &baseName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf("no entry found for UUID: %s", uuid)
		}
		return "", "", fmt.Errorf("database query failed: %v", err)
	}

	// Remove first slash from the absolutePath
	fullPath := filepath.Join(absolutePath[1:], pathFromRoot)
	return fullPath, baseName, nil
}
