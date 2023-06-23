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
	    );
	
	CREATE TABLE IF NOT EXISTS orders (
	    id UUID NOT NULL primary key, 
	    number BIGINT NOT NULL UNIQUE,
	    status VARCHAR(200) NOT NULL DEFAULT 'NEW', 
	    user_id UUID NOT NULL,
		uploaded_at TIMESTAMP,
	    CONSTRAINT fk_user
	        FOREIGN KEY (user_id)
	    		REFERENCES users(id)
	    );

	CREATE TABLE IF NOT EXISTS accounts (
	    id UUID NOT NULL primary key, 
	    current BIGINT NOT NULL UNIQUE DEFAULT 0,
	    user_id UUID NOT NULL UNIQUE,
	    created_at TIMESTAMP,
		updated_at TIMESTAMP,
	    CONSTRAINT fk_user
	        FOREIGN KEY (user_id)
	    		REFERENCES users(id)
	    );

	CREATE TABLE IF NOT EXISTS transacts (
	    id UUID NOT NULL primary key, 
	    amount BIGINT NOT NULL,
	    debit BOOL DEFAULT false,
	    order_number BIGINT NOT NULL UNIQUE,
	    account_id UUID NOT NULL,
		processed_at TIMESTAMP,
	    CONSTRAINT fk_account
	        FOREIGN KEY (account_id)
	    		REFERENCES accounts(id)
	    );
	    `

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
