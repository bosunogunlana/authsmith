package handlers

import (
	"database/sql"
	"html/template"
)

type Handlers struct {
	DB        *sql.DB
	Templates *template.Template
}
