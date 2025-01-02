package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type Follower struct {
	UserID    int64  `json:"user_id"`
	FollwerID int64  `json:"follwer_id"`
	CreatedAt string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {
	query := `
	INSERT INTO followers (user_id,follower_id) VALUES ($1,$2);
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrConflict
			}
		}
	}
	return nil
}
func (s *FollowerStore) UnFollow(ctx context.Context, followedID, userID int64) error {
	query := `
	DELETE FROM followers WHERE user_id=$1 AND follower_id=$2;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)

	return err
}
