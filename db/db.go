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

// Check if ticket has been redeemed
func CheckRedeem(ctx context.Context, ticketID uint64) (bool, error) {
	var redeemed bool
	err := db.QueryRow(
		ctx,
		`
		SELECT redeemed
		FROM redeems
		WHERE id = $1
		`,
		ticketID,
	).Scan(&redeemed)
	return redeemed, err
}

// Check if address owns ticket
func CheckTicketAddress(ctx context.Context, userAddress string, ticketID uint64) (bool, error) {
	var address string
	err := db.QueryRow(
		ctx,
		`
		SELECT owner_address
		FROM redeems
		WHERE id = $1
		`,
		ticketID,
	).Scan(&address)
	var owns bool
	if address == userAddress {
		owns = true
	}
	return owns, err
}

// redeem ticket
func Redeem(ctx context.Context, ticketID uint64, userAddress string) error {
	var id uint64
	err := db.QueryRow(
		ctx,
		`
		UPDATE redeems
		SET redeemed = TRUE
		WHERE id = $1 AND owner_address = $2 AND NOT redeemed
		RETURNING id
		`,
		ticketID,
		userAddress,
	).Scan(&id)
	return err
}

// transfer ticket
func Transfer(ctx context.Context, ticketID uint64, oldAddress string, newAddress string) error {
	var id uint64
	err := db.QueryRow(
		ctx,
		`
		UPDATE redeems
		SET owner_address = $1
		WHERE id = $2 AND owner_address = $3 AND NOT redeemed
		RETURNING id
		`,
		newAddress,
		ticketID,
		oldAddress,
	).Scan(&id)
	return err
}

// Create Redeem for a Ticket
func CreateRedeemForTicket(ctx context.Context, ticketIDs []uint64, address string) error {
	var id uint64
	err := db.QueryRow(
		ctx,
		`
		INSERT INTO redeems (id, owner_address)
		VALUES (UNNEST($1::INTEGER[]), $2)
		RETURNING id
		`,
		ticketIDs,
		address,
	).Scan(&id)
	return err
}
