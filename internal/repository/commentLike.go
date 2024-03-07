package repository

import (
	"database/sql"
	"forum/internal/models"
)

type CommentLikeSQLite struct {
	db *sql.DB
}

func NewCommentLikeSQLite(db *sql.DB) *CommentLikeSQLite {
	return &CommentLikeSQLite{db: db}
}

func (l *CommentLikeSQLite) CreateCommentHeart(heart *models.CommentHeart) error {
	query := "INSERT INTO comment_likes (user_id, comment_id, type) VALUES (?, ?, ?)"
	result, err := l.db.Exec(query, heart.UserID, heart.CommentID, heart.LikeType)
	if err != nil {
		return err
	}

	likeID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	heart.ID = int(likeID)

	return nil

}

func (l *CommentLikeSQLite) DeleteCommentHeart(commentID, userID int) error {
	query := "DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?"
	_, err := l.db.Exec(query, commentID, userID)
	return err
}

func (l *CommentLikeSQLite) GetCommentLikesByCommentID(commentID int) (*models.Like, error) {
	query := `
	SELECT 
    	COALESCE(SUM(CASE WHEN type = ? THEN 1 ELSE 0 END), 0) AS like_count,
    	COALESCE(SUM(CASE WHEN type = ? THEN 1 ELSE 0 END), 0) AS dislike_count
	FROM comment_likes
	WHERE comment_id = ?
	`

	row := l.db.QueryRow(query, models.LikeTypes.Like, models.LikeTypes.Dislike, commentID)

	var like models.Like
	err := row.Scan(&like.Like, &like.Dislike)
	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (l *CommentLikeSQLite) GetCommentLikeByUserIDAndCommentID(userID, commentID int) (*models.CommentHeart, error) {
	query := "SELECT id, user_id, comment_id, type FROM comment_likes WHERE user_id = ? AND comment_id = ?"
	row := l.db.QueryRow(query, userID, commentID)

	var like models.CommentHeart
	err := row.Scan(&like.ID, &like.UserID, &like.CommentID, &like.LikeType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &like, nil
}

func (l *CommentLikeSQLite) GetCommentLikesByUserID(userID int) ([]models.CommentHeart, error) {
	query := "SELECT id, user_id, comment_id, type FROM comment_likes WHERE user_id = ?"
	rows, err := l.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.CommentHeart
	for rows.Next() {
		var like models.CommentHeart
		err := rows.Scan(&like.ID, &like.UserID, &like.CommentID, &like.LikeType)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}
