package mapmarkpg

import (
	"context"
	"database/sql"
	"echoFramework/internal/domain"
	"errors"
	"fmt"
	trmsql "github.com/avito-tech/go-transaction-manager/drivers/sql/v2"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db     *sql.DB
	getter *trmsql.CtxGetter
}

func New(db *sql.DB, c *trmsql.CtxGetter) *Repository {
	return &Repository{db: db, getter: c}
}

func (r *Repository) Get(ctx context.Context) ([]domain.Mark, error) {
	const q = `SELECT id, name, lat, lng FROM mapmarks`

	marks := make([]domain.Mark, 0)

	row, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, q)
	if err != nil {
		return marks, err
	}
	defer row.Close()

	for row.Next() {
		mark := domain.Mark{}
		err := row.Scan(&mark.Id, &mark.Name, &mark.Lat, &mark.Lng)
		if err != nil {
			return marks, err
		}
		marks = append(marks, mark)
	}

	return marks, nil
}

func (r *Repository) Save(ctx context.Context, mark domain.Mark) (int64, error) {
	const q = `INSERT INTO mapmarks (name, lat, lng) VALUES (?, ?, ?)`

	res, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, q, mark.Name, mark.Lat, mark.Lng)

	var sqliteErr sqlite3.Error
	if err != nil {
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintNotNull) {
				return 0, fmt.Errorf("%w: %s", domain.ErrInvalidParams, sqliteErr.Error())
			}
			return 0, sqliteErr
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	const q = `DELETE FROM mapmarks WHERE id = ?`

	_, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetById(ctx context.Context, id int) (domain.Mark, error) {
	const q = `SELECT id, name, description, lat, lng, image FROM mapmarks WHERE id = ?`

	var mark domain.Mark
	row, err := r.getter.DefaultTrOrDB(ctx, r.db).QueryContext(ctx, q, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return mark, domain.ErrNotFound
		}
		return mark, err
	}
	defer row.Close()

	if row.Next() {
		err := row.Scan(&mark.Id, &mark.Name, &mark.Description, &mark.Lat, &mark.Lng, &mark.Image)
		if err != nil {
			return mark, err
		}
		return mark, nil
	}

	return mark, domain.ErrNotFound
}

func (r *Repository) Update(ctx context.Context, mark domain.Mark) error {
	const q = `UPDATE mapmarks SET name = ?, description = ?, lat = ?, lng = ?, image = ? WHERE id = ?`

	res, err := r.getter.DefaultTrOrDB(ctx, r.db).ExecContext(ctx, q,
		mark.Name, mark.Description,
		mark.Lat, mark.Lng,
		mark.Image,
		mark.Id,
	)
	if err != nil {
		return err
	}

	if n, _ := res.RowsAffected(); n == 0 {
		return domain.ErrNotFound
	}

	return nil
}
