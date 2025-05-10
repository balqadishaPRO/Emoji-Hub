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
	const q = `
SELECT id, name, category, "group", html_code, unicode
FROM   emoji
ORDER  BY id::text               -- сортируем по id
OFFSET $1
LIMIT  $2`

	rows, err := r.DB.QueryContext(ctx, q,
		p.Offset, p.Limit)
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
		`SELECT id,name,category,"group",html_code,unicode
		  FROM emoji WHERE id=$1`,
		id).Scan(&e.ID, &e.Name, &e.Category,
		&e.Group, pq.Array(&e.HtmlCode), pq.Array(&e.Unicode))
	return e, err
}
