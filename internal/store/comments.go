package store

import (
	"context"
	"database/sql"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (c *CommentStore) GetPostById(ctx context.Context, postId int64) ([]Comment, error) {
	query := `
	SELECT
	c.id, c.user_id, c.post_id, c.content, c.created_at, u.id, u.username 
FROM
	comments c
	JOIN users u ON u.id = c.user_id
WHERE
	c.post_id = $1
ORDER BY
	c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.UserID, &c.PostID, &c.Content, &c.CreatedAt, &c.User.ID, &c.User.Username)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (cs *CommentStore) Create(ctx context.Context, comment *Comment) error {
	query := `INSERT INTO comments (post_id, user_id, content)
	VALUES ($1, $2, $3) RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(ctx, QUERY_WRITE_TIME_OUR_DURATION)
	defer cancel()

	err := cs.db.QueryRowContext(ctx, query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
