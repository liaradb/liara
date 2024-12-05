package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	db, err := ConnectPostgresDB("postgresql://postgres:example@localhost:5432/liara?sslmode=disable")
	if err != nil {
		return err
	}

	ctx := context.Background()

	// err = DropTable(ctx, db)

	err = CreateTable(ctx, db)
	if err != nil {
		return err
	}

	// err = InsertRow(ctx, db)

	err = SelectRows(ctx, db)
	if err != nil {
		return err
	}

	return nil
}

func ConnectPostgresDB(uri string) (*sql.DB, error) {
	return sql.Open("postgres", uri)
}

func InsertRow(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
INSERT INTO outbox (
    transaction_id,
    publication_id,
    message_id,
    message_type,
    data,
    scheduled
) VALUES (
    pg_current_xact_id(),
    'p1',
    'm1',
    'type1',
    '""',
    now()
);
	`)
	return err
}

func SelectRows(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `
SELECT
	position,
	transaction_id,
	publication_id,
	message_id,
	message_type,
	data,
	scheduled
FROM outbox
WHERE
	(
		(
			transaction_id = $1
			AND position > 0
		)
		OR
		(
			transaction_id > $1
		)
	)
	AND transaction_id < pg_snapshot_xmin(pg_current_snapshot())
ORDER BY
	transaction_id ASC,
	position ASC
LIMIT 100;`,
		0)
	if err != nil {
		return err
	}

	defer func() { rows.Close() }()

	for rows.Next() {
		row := outboxRow{}
		err := rows.Scan(
			&row.Position,
			&row.TransactionID,
			&row.PublicationID,
			&row.MessageID,
			&row.MessageType,
			&row.Data,
			&row.Scheduled,
		)
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(row, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))
	}
	return err
}

type outboxRow struct {
	Position      int64     `json:"position"`
	TransactionID uint64    `json:"transactionId"`
	PublicationID string    `json:"publicationId"`
	MessageID     string    `json:"messageId"`
	MessageType   string    `json:"messageType"`
	Data          string    `json:"data"`
	Scheduled     time.Time `json:"scheduled"`
}

func CreateTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS outbox (
   position        SERIAL                    PRIMARY KEY,
   transaction_id  xid8                      NOT NULL,
   publication_id  VARCHAR(250)              NOT NULL,
   message_id      VARCHAR(250)              NOT NULL,
   message_type    VARCHAR(250)              NOT NULL,
   data            JSONB                     NOT NULL,
   scheduled       TIMESTAMP WITH TIME ZONE  NOT NULL    default (now())
);
`)
	return err
}

func DropTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
DROP TABLE outbox;
	`)
	return err
}
