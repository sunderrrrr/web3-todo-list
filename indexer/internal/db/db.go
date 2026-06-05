package db

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	_ "github.com/lib/pq"
)

const ZeroAddress = "0x0000000000000000000000000000000000000000"

type Postgres struct {
	Conn *sql.DB
}

func New(ctx context.Context, dsn string) (*Postgres, error) {
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	db := &Postgres{Conn: conn}
	if err := db.migrate(ctx); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return db, nil
}

func (db *Postgres) migrate(ctx context.Context) error {
	_, err := db.Conn.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS transfers (
			id SERIAL PRIMARY KEY,
			from_addr TEXT NOT NULL,
			to_addr TEXT NOT NULL,
			value NUMERIC NOT NULL,
			tx_hash TEXT NOT NULL UNIQUE,
			block_number BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS todos (
			id NUMERIC PRIMARY KEY,
			text TEXT NOT NULL,
			owner TEXT NOT NULL,
			tx_hash TEXT NOT NULL UNIQUE,
			block_number BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS rewards (
			id SERIAL PRIMARY KEY,
			user_addr TEXT NOT NULL,
			amount NUMERIC NOT NULL,
			tx_hash TEXT NOT NULL UNIQUE,
			block_number BIGINT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS balances (
			addr TEXT PRIMARY KEY,
			balance NUMERIC NOT NULL DEFAULT 0
		);
	`)
	return err
}

func (db *Postgres) SaveTransfer(from, to string, value *big.Int, txHash string, blockNum uint64) error {
	_, err := db.Conn.Exec(
		`INSERT INTO transfers (from_addr, to_addr, value, tx_hash, block_number)
		 VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING`,
		from, to, value.String(), txHash, blockNum,
	)
	if err != nil {
		return fmt.Errorf("save transfer: %w", err)
	}

	db.Conn.Exec(`INSERT INTO balances (addr, balance) VALUES ($1, 0) ON CONFLICT DO NOTHING`, from)
	db.Conn.Exec(`INSERT INTO balances (addr, balance) VALUES ($1, 0) ON CONFLICT DO NOTHING`, to)

	if from != ZeroAddress {
		db.Conn.Exec(`UPDATE balances SET balance = balance - $1 WHERE addr = $2`, value.String(), from)
	}
	if to != ZeroAddress {
		db.Conn.Exec(`UPDATE balances SET balance = balance + $1 WHERE addr = $2`, value.String(), to)
	}
	return nil
}

func (db *Postgres) SaveTodo(id uint64, text, owner, txHash string, blockNum uint64) error {
	_, err := db.Conn.Exec(
		`INSERT INTO todos (id, text, owner, tx_hash, block_number)
		 VALUES ($1,$2,$3,$4,$5) ON CONFLICT DO NOTHING`,
		id, text, owner, txHash, blockNum,
	)
	return err
}

func (db *Postgres) SaveReward(userAddr string, amount *big.Int, txHash string, blockNum uint64) error {
	_, err := db.Conn.Exec(
		`INSERT INTO rewards (user_addr, amount, tx_hash, block_number)
		 VALUES ($1,$2,$3,$4) ON CONFLICT DO NOTHING`,
		userAddr, amount.String(), txHash, blockNum,
	)
	return err
}

func (db *Postgres) Close() error {
	return db.Conn.Close()
}
