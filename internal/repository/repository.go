package repository

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// ErrorDuplicate указывает на конфликт данных в хранилище
var ErrorDuplicate = errors.New("data conflict duplicate")

type Repository struct {
	db *sql.DB
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *Repository) Bootstrap(ctx context.Context) error {

	sqlStatement := `
	CREATE TABLE IF NOT EXISTS users (
	    id UUID NOT NULL primary key, 
	    login VARCHAR(200) NOT NULL UNIQUE, 
	    password VARCHAR(200) NOT NULL
	    )`

	_, err := r.db.ExecContext(
		ctx,
		sqlStatement,
	)

	return err
}

func NewRepository(dsn string) (*Repository, error) {

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
	}, nil
}
