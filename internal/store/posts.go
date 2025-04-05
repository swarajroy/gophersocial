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
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{db: db}
}

func (ps *PostStore) GetUserFeed(ctx context.Context, userID int64) ([]PostWithMetadata, error) {
	query := `
		SELECT p.id, p.user_id, p.title, p.content, p.tags, p.created_at, COUNT(c.id) AS comments_count, u.username
			FROM 
		posts p
			LEFT JOIN comments c ON c.post_id = p.id
			LEFT JOIN users u ON p.user_id = u.id
			JOIN followers f ON f.follower_id = p.user_id or p.user_id = $1
			WHERE f.user_id = $1 or p.user_id = $1
			GROUP BY p.id, u.username
			ORDER BY p.created_at DESC;
	`
	ctx, cancel := context.WithTimeout(ctx, QUERY_READ_TIME_OUR_DURATION)
	defer cancel()

	rows, err := ps.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetadata

	for rows.Next() {
		var singleFeed PostWithMetadata
		err := rows.Scan(
			&singleFeed.ID,
			&singleFeed.UserID,
			&singleFeed.Title,
			&singleFeed.Content,
			pq.Array(&singleFeed.Tags),
			&singleFeed.CreatedAt,
			&singleFeed.CommentsCount,
			&singleFeed.User.Username,
		)
		if err != nil {
			return nil, err
		}

		feed = append(feed, singleFeed)

	}

	return feed, nil
}

func (ps *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (title, content, tags, user_id)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_READ_TIME_OUR_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

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

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

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
