package repository

import (
	"database/sql"
	"forum/internal/models"
	"time"
)

type CommentSQLite struct {
	db *sql.DB
}

func NewCommentSQLite(db *sql.DB) *CommentSQLite {
	return &CommentSQLite{db: db}
}

func (c *CommentSQLite) CreateComment(comment *models.Comment) error {
	query := "INSERT INTO comments (content, user_id, post_id, created_at) VALUES (?, ?, ?, ?)"
	result, err := c.db.Exec(query, comment.Content, comment.AuthorID, comment.PostID, time.Now())
	if err != nil {
		return err
	}

	commentID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = int(commentID)

	return nil
}

func (c *CommentSQLite) GetCommentsByPostID(postID int) (*[]models.Comment, error) {
	query := `SELECT comments.id, comments.content, users.username, comments.created_at 
	FROM comments 
	JOIN users ON comments.user_id = users.id 
	WHERE post_id = ?
	`
	rows, err := c.db.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment
		err := rows.Scan(&comment.ID, &comment.Content, &comment.AuthorName, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}
