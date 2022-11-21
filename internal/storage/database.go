package storage

import (
	"context"

	"github.com/0xc00000f/shortener-tpl/internal/encoder"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type DatabaseStorage struct {
	db *pgxpool.Pool

	l *zap.Logger
}

type urlMapping struct {
	userID uuid.UUID `db:"user_id"`
	short  string    `db:"short_url"`
	long   string    `db:"long_url"`
}

func NewDatabaseStorage(connStr string, l *zap.Logger) (DatabaseStorage, error) {
	pgxConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		l.Error("unable to parsing config", zap.Error(err))
		return DatabaseStorage{}, err
	}

	pgxConnPool, err := pgxpool.ConnectConfig(context.TODO(), pgxConfig)
	if err != nil {
		l.Error("Unable to connect to database", zap.Error(err))
		return DatabaseStorage{}, err
	}

	query := `CREATE TABLE IF NOT EXISTS url_mapping
		(
		id integer NOT NULL GENERATED ALWAYS AS IDENTITY,
		user_id uuid NOT NULL,
		short_url text NOT NULL,
		long_url text NOT NULL,
		PRIMARY KEY (id));
		CREATE UNIQUE INDEX IF NOT EXISTS long_unique_idx on url_mapping (long_url);`

	_, err = pgxConnPool.Exec(context.TODO(), query)
	if err != nil {
		return DatabaseStorage{}, err
	}

	return DatabaseStorage{db: pgxConnPool, l: l}, nil
}

func (ds DatabaseStorage) Get(short string) (string, error) {
	var m urlMapping
	err := ds.db.QueryRow(context.TODO(), "SELECT user_id, short_url, long_url FROM url_mapping WHERE short_url = $1::text", short).Scan(&m.userID, &m.short, &m.long)
	if err != nil {
		return "", err
	}

	return m.long, nil
}

func (ds DatabaseStorage) GetAll(userID uuid.UUID) (result map[string]string, err error) {
	m := make(map[string]string)

	rows, err := ds.db.Query(context.TODO(), "SELECT user_id, short_url, long_url FROM url_mapping WHERE user_id = $1::uuid", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u urlMapping
		err = rows.Scan(&u.userID, &u.short, &u.long)
		if err != nil {
			return nil, err
		}

		m[u.short] = u.long
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (ds DatabaseStorage) Store(userID uuid.UUID, short string, long string) error {
	query := `INSERT INTO url_mapping(user_id, short_url, long_url) VALUES ($1::uuid, $2::text, $3::text)`
	_, err := ds.db.Exec(context.TODO(), query, userID, short, long)

	if pqErr, ok := err.(*pgconn.PgError); ok {
		if pqErr.Code == pgerrcode.UniqueViolation {
			var m urlMapping
			if err := ds.db.QueryRow(context.TODO(), "SELECT user_id, short_url, long_url FROM url_mapping WHERE long_url = $1::text", long).Scan(&m.userID, &m.short, &m.long); err != nil {
				return err
			}
			return &encoder.UniqueViolationError{
				Err:    err,
				UserID: m.userID,
				Short:  m.short,
				Long:   m.long,
			}
		}
		return err
	}
	return nil
}

func (ds DatabaseStorage) IsKeyExist(short string) (bool, error) {
	var i bool
	row := ds.db.QueryRow(context.TODO(), `SELECT COUNT(1)>0 AS N FROM url_mapping WHERE short_url = $1`, short)

	err := row.Scan(&i)
	if err != nil {
		return false, err
	}
	return i, nil
}
