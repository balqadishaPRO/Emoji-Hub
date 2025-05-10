package repo

import (
	"context"

	"github.com/balqadishaPRO/Emoji-Hub/internal/model"
	"github.com/lib/pq"
)

func (r *Repo) AddFav(ctx context.Context, sessID, emojiID string) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO favorites(session_id,emoji_id)
		 VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		sessID, emojiID)
	return err
}

func (r *Repo) DelFav(ctx context.Context, sessID, emojiID string) error {
	_, err := r.DB.ExecContext(ctx,
		`DELETE FROM favorites WHERE session_id=$1 AND emoji_id=$2`,
		sessID, emojiID)
	return err
}

func (r *Repo) ListFav(ctx context.Context, sessID string) ([]model.Emoji, error) {
	const q = `
	SELECT e.id,e.name,e.category,e."group",e.html_code,e.unicode
	FROM   favorites f
	JOIN   emoji     e ON e.id = f.emoji_id
	WHERE  f.session_id=$1`
	rows, err := r.DB.QueryContext(ctx, q, sessID)
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
