package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
}

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

func (ps *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (title, content, tags, user_id)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := ps.db.QueryRowContext(ctx, query,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
		post.UserID,
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PostStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `
	SELECT id, user_id, title, content, created_at, updated_at, tags, version  from 
	posts 
	where id = $1
	`
	var post Post

	err := ps.db.QueryRowContext(ctx, query, postId).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

func (ps *PostStore) Delete(ctx context.Context, postId int64) error {
	query := `DELETE from posts where id = $1`

	result, err := ps.db.ExecContext(ctx, query, postId)
	if err != nil {
		return err
	}

	num, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if num == 0 {
		return ErrNotFound
	}

	return nil
}

func (ps *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
	SET title = $1, content = $2, version = version + 1
	WHERE id = $3 and version = $4 RETURNING version`

	row := ps.db.QueryRowContext(ctx, query, post.Title, post.Content, post.ID, post.Version)

	err := row.Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
