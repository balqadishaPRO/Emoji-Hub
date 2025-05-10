package service

import (
	"context"
	"database/sql"

	"github.com/balqadishaPRO/Emoji-Hub/internal/llm"
	"github.com/balqadishaPRO/Emoji-Hub/internal/model"
	"github.com/balqadishaPRO/Emoji-Hub/internal/repo"
)

type EmojiService struct{ Repo *repo.Repo }

func (s *EmojiService) List(ctx context.Context, p repo.ListParams) ([]model.Emoji, error) {
	if p.Limit == 0 || p.Limit > 100 {
		p.Limit = 50
	}
	return s.Repo.ListEmoji(ctx, p)
}

func (s *EmojiService) Detail(ctx context.Context, id string) (model.EmojiDetail, error) {
	e, err := s.Repo.GetEmoji(ctx, id)
	if err != nil {
		return model.EmojiDetail{}, err
	}

	var mood string
	err = s.Repo.DB.QueryRowContext(ctx,
		`SELECT mood FROM llm_cache WHERE emoji_id=$1`, id).Scan(&mood)

	if err == sql.ErrNoRows {
		mood, err = llm.GenerateMood(e.Name)
		if err != nil {
			return model.EmojiDetail{Emoji: e}, nil
		}
		_, _ = s.Repo.DB.ExecContext(ctx,
			`INSERT INTO llm_cache(emoji_id,mood)
			 VALUES($1,$2) ON CONFLICT (emoji_id)
			 DO UPDATE SET mood=$2,updated=NOW()`,
			id, mood)
	}

	return model.EmojiDetail{Emoji: e, Mood: mood}, nil
}

func (s *EmojiService) ListFav(ctx context.Context, sess string) ([]model.Emoji, error) {
	return s.Repo.ListFav(ctx, sess)
}

func (s *EmojiService) AddFav(ctx context.Context, sess, id string) error {
	return s.Repo.AddFav(ctx, sess, id)
}

func (s *EmojiService) DelFav(ctx context.Context, sess, id string) error {
	return s.Repo.DelFav(ctx, sess, id)
}
