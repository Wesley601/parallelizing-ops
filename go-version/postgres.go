package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostGres struct {
	conn *pgxpool.Pool
}

func NewPostGres() *PostGres {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn, err := pgxpool.New(ctx, "postgres://erickwendel:mypassword@localhost:5432/school")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	conn.Exec(ctx, "DELETE FROM students")
	return &PostGres{
		conn: conn,
	}
}

func (p *PostGres) Close() {
	p.conn.Close()
}

func (p *PostGres) Insert(items []Student) error {
	var batch pgx.Batch
	for _, i := range items {
		batch.Queue("INSERT INTO students(name, email, age, registered_at) VALUES($1, $2, $3, $4)", i.Name, i.Email, i.Age, i.RegisteredAt)
	}

	return p.conn.SendBatch(context.Background(), &batch).Close()
}

func (p *PostGres) Count() (int, error) {
	var total int
	row := p.conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM students")
	if err := row.Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}
