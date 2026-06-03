package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RunMigration(db *sql.DB, dirName string) error {
	directory, err := os.ReadDir(dirName)
	if err != nil {
		return err
	}

	migrationFiles := []string{}
	for _, entry := range directory {
		fileName := entry.Name()
		if strings.HasSuffix(fileName, ".sql") {
			migrationFiles = append(migrationFiles, fileName)
		}
	}

	for _, migrationFile := range migrationFiles {
		filePath := filepath.Join(dirName, migrationFile)
		file, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		_, err = db.Exec(string(file))
		if err != nil {
			return err
		}
		log.Printf("%s migrated", migrationFile)
	}

	return nil
}