package postgres

import (
	"LostAndFound/internal/domain/entity"
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

type CardRepository struct {
	db *sql.DB
}

func (l *CardRepository) Create(ctx context.Context, card *entity.Card) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	insertCardQuery := `
		INSERT INTO cards (id, title, description, owner_id, preview_url, location, city, street, status)
		VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($6, $7), 4326), $8, $9, $10)
	`

	_, err = tx.ExecContext(ctx, insertCardQuery,
		card.ID,
		card.Title,
		card.Description,
		card.Owner.ID,
		card.PreviewURL,
		card.Longitude,
		card.Latitude,
		card.City,
		card.Street,
		card.Status,
	)
	if err != nil {
		return fmt.Errorf("failed to insert card: %w", err)
	}

	insertImageQuery := `
		INSERT INTO card_images (id, card_id, url)
		VALUES ($1, $2, $3)
	`

	for _, url := range card.Images {
		_, err = tx.ExecContext(ctx, insertImageQuery, uuid.New().String(), card.ID, url)
		if err != nil {
			return fmt.Errorf("failed to insert image: %w", err)
		}
	}

	return tx.Commit()
}

func (l *CardRepository) GetByID(ctx context.Context, id string) (*entity.Card, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
	SELECT 
		l.id, l.title, l.description, l.city, l.street, l.status, l.created_at,
		ST_Y(l.location::geometry),
		ST_X(l.location::geometry),
		u.id, u.name, u.surname, u.phone, u.telegram
	FROM cards l
	JOIN users u ON l.owner_id = u.id
	WHERE l.id = $1;
	`

	var card entity.Card
	var owner entity.Owner

	if err = tx.QueryRowContext(ctx, query, id).Scan(
		&card.ID,
		&card.Title,
		&card.Description,
		&card.City,
		&card.Street,
		&card.Status,
		&card.CreatedAt,
		&card.Latitude,
		&card.Longitude,
		&owner.ID,
		&owner.Name,
		&owner.Surname,
		&owner.Phone,
		&owner.Telegram,
	); err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	card.Owner = owner
	card.OwnerID = owner.ID

	imgQuery := `SELECT url FROM card_images WHERE card_id = $1`
	rows, err := tx.QueryContext(ctx, imgQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err = rows.Scan(&url); err != nil {
			return nil, err
		}
		card.Images = append(card.Images, url)
	}

	return &card, tx.Commit()
}

func (l *CardRepository) FindAll(ctx context.Context, status string) ([]*entity.Card, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT 
			l.id, l.title, l.description, l.city, l.street, l.status, 
			l.preview_url, l.created_at,
			ST_Y(l.location::geometry),
			ST_X(l.location::geometry),
			u.id, u.name, u.surname, u.phone, u.telegram
		FROM cards l
		JOIN users u ON l.owner_id = u.id
	`
	var args []interface{}
	if status != "" {
		query += " WHERE l.status = $1 ORDER BY l.created_at DESC"
		args = append(args, status)
	} else {
		query += " ORDER BY l.created_at DESC"
	}

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying cards: %w", err)
	}
	defer rows.Close()

	var cards []*entity.Card
	for rows.Next() {
		var card entity.Card
		var owner entity.Owner

		err = rows.Scan(
			&card.ID,
			&card.Title,
			&card.Description,
			&card.City,
			&card.Street,
			&card.Status,
			&card.PreviewURL,
			&card.CreatedAt,
			&card.Latitude,
			&card.Longitude,
			&owner.ID,
			&owner.Name,
			&owner.Surname,
			&owner.Phone,
			&owner.Telegram,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning card row: %w", err)
		}

		card.Owner = owner
		cards = append(cards, &card)
	}

	return cards, tx.Commit()
}

func (l *CardRepository) Update(ctx context.Context, card *entity.Card) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE cards SET 
			title = $1, 
			description = $2, 
			city = $3, 
			street = $4,
			status = $5,
			preview_url = $6,
			location = ST_SetSRID(ST_MakePoint($7, $8), 4326)
		WHERE id = $9;
	`
	if _, err = tx.ExecContext(ctx, query,
		card.Title,
		card.Description,
		card.City,
		card.Street,
		card.Status,
		card.PreviewURL,
		card.Longitude,
		card.Latitude,
		card.ID,
	); err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	if len(card.Images) > 0 {
		delQuery := `DELETE FROM card_images WHERE card_id = $1;`
		if _, err = tx.ExecContext(ctx, delQuery, card.ID); err != nil {
			return fmt.Errorf("failed to delete card images: %w", err)
		}

		insertQuery := `INSERT INTO card_images (id, card_id, url) VALUES ($1, $2, $3);`
		for _, img := range card.Images {
			if _, err = tx.ExecContext(ctx, insertQuery, uuid.New().String(), card.ID, img); err != nil {
				return fmt.Errorf("failed to insert card image: %w", err)
			}
		}
	}

	return tx.Commit()
}

func (l *CardRepository) Delete(ctx context.Context, id string) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	deleteImagesQuery := `DELETE FROM card_images WHERE card_id = $1`
	if _, err = tx.ExecContext(ctx, deleteImagesQuery, id); err != nil {
		return fmt.Errorf("failed to delete card images: %w", err)
	}

	deleteCardQuery := `DELETE FROM cards WHERE id = $1`
	if _, err = tx.ExecContext(ctx, deleteCardQuery, id); err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return tx.Commit()
}

func (l *CardRepository) FindNearLocation(ctx context.Context, lat, lon, radius float64, status string) ([]*entity.Card, error) {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		SELECT 
			l.id, l.title, l.description, l.preview_url, l.status, l.created_at, l.city, l.street,
			ST_Distance(l.location, ST_MakePoint($1, $2)::geography) as distance_m,
			u.id, u.name, u.surname, u.telegram
		FROM cards l
		JOIN users u ON l.owner_id = u.id
		WHERE ST_DWithin(l.location, ST_MakePoint($1, $2)::geography, $3)
	`

	var args []interface{}
	args = append(args, lon, lat, radius)

	if status != "" {
		query += " AND l.status = $4"
		args = append(args, status)
	}

	query += " ORDER BY distance_m ASC"

	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying nearby cards: %w", err)
	}
	defer rows.Close()

	var cards []*entity.Card
	for rows.Next() {
		var card entity.Card
		var owner entity.Owner

		err = rows.Scan(
			&card.ID,
			&card.Title,
			&card.Description,
			&card.PreviewURL,
			&card.Status,
			&card.CreatedAt,
			&card.City,
			&card.Street,
			&card.DistanceM,
			&owner.ID,
			&owner.Name,
			&owner.Surname,
			&owner.Telegram,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning nearby card row: %w", err)
		}

		card.Owner = owner
		cards = append(cards, &card)
	}

	return cards, tx.Commit()
}

func NewCardRepo(db *sql.DB) *CardRepository {
	return &CardRepository{db: db}
}
