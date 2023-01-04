package storage

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/0xc00000f/shortener-tpl/internal/models"
)

type DatabaseStorage struct {
	db *pgxpool.Pool

	l *zap.Logger
}

func NewDatabaseStorage(
	ctx context.Context,
	pgxConnPool *pgxpool.Pool,
	l *zap.Logger,
) (DatabaseStorage, error) {
	query := `CREATE TABLE IF NOT EXISTS url_mapping
		(
		id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
		user_id uuid NOT NULL,
		short_url text NOT NULL,
		long_url text NOT NULL,
		is_active boolean DEFAULT true,
		PRIMARY KEY (id));
		CREATE UNIQUE INDEX IF NOT EXISTS long_unique_idx on url_mapping (long_url);`

	_, err := pgxConnPool.Exec(ctx, query)
	if err != nil {
		return DatabaseStorage{}, err
	}

	return DatabaseStorage{db: pgxConnPool, l: l}, nil
}

func (ds DatabaseStorage) Get(ctx context.Context, short string) (string, error) {
	var m models.URL

	err := ds.db.QueryRow(
		ctx,
		"SELECT user_id, short_url, long_url, is_active FROM url_mapping WHERE short_url = $1::text",
		short,
	).Scan(&m.UserID, &m.Short, &m.Long, &m.IsActive)
	if err != nil {
		return "", err
	}

	if !m.IsActive {
		return m.Long, URLDeletedError{}
	}

	return m.Long, nil
}

func (ds DatabaseStorage) GetAll(
	ctx context.Context,
	userID uuid.UUID,
) (result map[string]string, err error) {
	m := make(map[string]string)

	rows, err := ds.db.Query(
		ctx,
		"SELECT user_id, short_url, long_url FROM url_mapping WHERE user_id = $1::uuid",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u models.URL

		err = rows.Scan(&u.UserID, &u.Short, &u.Long)
		if err != nil {
			return nil, err
		}

		m[u.Short] = u.Long
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ds DatabaseStorage) Store(ctx context.Context, userID uuid.UUID, short string, long string) error {
	query := `INSERT INTO url_mapping(user_id, short_url, long_url) VALUES ($1::uuid, $2::text, $3::text)`

	_, err := ds.db.Exec(ctx, query, userID, short, long)
	if err != nil {
		var pgError *pgconn.PgError
		if !errors.As(err, &pgError) {
			return err
		}

		pgErr, ok := err.(*pgconn.PgError) //nolint:errorlint
		if !ok {
			return err
		}

		if pgErr.Code != pgerrcode.UniqueViolation {
			return err
		}

		var m models.URL
		if err := ds.db.QueryRow(
			ctx,
			"SELECT user_id, short_url, long_url FROM url_mapping WHERE long_url = $1::text",
			long,
		).Scan(&m.UserID, &m.Short, &m.Long); err != nil {
			return err
		}

		return &encoder.UniqueViolationError{
			Err:    err,
			UserID: m.UserID,
			Short:  m.Short,
			Long:   m.Long,
		}
	}

	return nil
}

func (ds DatabaseStorage) IsKeyExist(ctx context.Context, short string) (bool, error) {
	var i bool

	row := ds.db.QueryRow(
		ctx,
		`SELECT COUNT(1)>0 AS N FROM url_mapping WHERE short_url = $1`,
		short,
	)

	err := row.Scan(&i)
	if err != nil {
		return false, err
	}

	return i, nil
}

func (ds DatabaseStorage) Delete(ctx context.Context, data []models.URL) error {
	log.Printf("db: delete in, data len: %d, data: %v", len(data), data)
	defer log.Printf("db: delete out")

	conn, err := ds.db.Acquire(ctx)
	if err != nil {
		log.Printf("Unable to acquire a database connection: %v\n", err)
		return err
	}
	defer conn.Release()

	for _, url := range data {
		sqlStatement := `
				UPDATE url_mapping
				SET is_active = false
				WHERE user_id = $1::uuid AND short_url = $2::text;
				`

		_, err := conn.Exec(ctx,
			sqlStatement, url.UserID, url.Short)
		if err != nil {
			log.Printf("Unable to DELETE: %v\n", err)
			return err
		}
	}

	return nil
}
