package repo

import (
	"context"

	"github.com/balqadishaPRO/Emoji-Hub/internal/model"
	"github.com/lib/pq"
)

type (
	ListParams struct {
		Search, Category, Group, Sort string
		Limit, Offset                 int
	}
)

func (r *Repo) ListEmoji(ctx context.Context, p ListParams) ([]model.Emoji, error) {
	query := `
		SELECT id, name, category, "group", html_code, unicode
		FROM emoji
		WHERE ($1 = '' OR name ILIKE $1)
		AND ($2 = '' OR category = $2)
		AND ($3 = '' OR "group" = $3)
		ORDER BY 
			CASE WHEN $4 = 'name' THEN name END,
			CASE WHEN $4 = 'category' THEN category END`

	rows, err := r.DB.QueryContext(ctx, query,
		"%"+p.Search+"%", p.Category, p.Group, p.Sort)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Emoji
	for rows.Next() {
		var e model.Emoji
		if err := rows.Scan(&e.ID, &e.Name, &e.Category,
			&e.Group, pq.Array(&e.HtmlCode), pq.Array(&e.Unicode)); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *Repo) GetEmoji(ctx context.Context, id string) (model.Emoji, error) {
	var e model.Emoji
	err := r.DB.QueryRowContext(ctx,
		`SELECT id, name, category, "group", html_code, unicode
		FROM emoji WHERE id = $1`,
		id).Scan(&e.ID, &e.Name, &e.Category,
		&e.Group, pq.Array(&e.HtmlCode), pq.Array(&e.Unicode))
	return e, err
}

func (r *Repo) ImportEmojis(ctx context.Context, emojis []model.Emoji) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO emoji(id, name, category, "group", html_code, unicode)
		VALUES($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			category = EXCLUDED.category,
			"group" = EXCLUDED."group",
			html_code = EXCLUDED.html_code,
			unicode = EXCLUDED.unicode`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, e := range emojis {
		_, err = stmt.ExecContext(ctx,
			e.ID, e.Name, e.Category, e.Group,
			pq.Array(e.HtmlCode), pq.Array(e.Unicode))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
