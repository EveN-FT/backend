package db

import (
	"context"
	"fmt"

	"github.com/EveN-FT/backend/config"
	"github.com/EveN-FT/backend/models"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

func init() {
	dbConf, err := pgxpool.ParseConfig(config.Conf.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("Malformed DB connection string: %v", config.Conf.DatabaseURL))
	}
	dbConf.MaxConns = 50

	db, err = pgxpool.ConnectConfig(context.Background(), dbConf)
	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database at %v\n", config.Conf.DatabaseURL))
	}
}

// CreateEvent creates a new event in the database.
func CreateEvent(ctx context.Context, event *models.Event) (uint64, error) {
	var id uint64
	err := db.QueryRow(
		ctx,
		`
		INSERT INTO events (address, owner_address)
		VALUES($1, $2)
		RETURNING id
		`,
		event.Address,
		event.OwnerAddress,
	).Scan(&id)
	return id, err
}

// ListEvents list events in the database.
func ListEvents(ctx context.Context) ([]*models.Event, error) {
	events := make([]*models.Event, 0)
	err := pgxscan.Select(
		ctx, db, &events,
		`
		SELECT *
		FROM events
		`,
	)
	return events, err
}

// ListEventsByOwnerAddress list events in the database given an owner address
func ListEventsByOwnerAddress(ctx context.Context, ownerAddress string) ([]*models.Event, error) {
	events := make([]*models.Event, 0)
	err := pgxscan.Select(
		ctx, db, &events,
		`
		SELECT *
		FROM events
		WHERE owner_address = $1
		`,
		ownerAddress,
	)
	return events, err
}
