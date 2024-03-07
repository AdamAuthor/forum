package repository

import (
	"database/sql"
	"forum/internal/models"
)

type LikeSQLite struct {
	db *sql.DB
}

func NewLikeSQLite(db *sql.DB) *LikeSQLite {
	return &LikeSQLite{db: db}
}

func (l *LikeSQLite) CreateHeart(heart *models.Heart) error {
	query := "INSERT INTO likes (user_id, post_id, type) VALUES (?, ?, ?)"
	result, err := l.db.Exec(query, heart.UserID, heart.PostID, heart.LikeType)
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

func (l *LikeSQLite) DeleteHeart(postID, userID int) error {
	query := "DELETE FROM likes WHERE post_id = ? AND user_id = ?"
	_, err := l.db.Exec(query, postID, userID)
	return err
}

func (l *LikeSQLite) GetLikesByPostID(postID int) (*models.Like, error) {
	query := `
	SELECT 
    	COALESCE(SUM(CASE WHEN type = ? THEN 1 ELSE 0 END), 0) AS like_count,
    	COALESCE(SUM(CASE WHEN type = ? THEN 1 ELSE 0 END), 0) AS dislike_count
	FROM likes
	WHERE post_id = ?
	`

	row := l.db.QueryRow(query, models.LikeTypes.Like, models.LikeTypes.Dislike, postID)

	var like models.Like
	err := row.Scan(&like.Like, &like.Dislike)
	if err != nil {
		return nil, err
	}

	return &like, nil
}

func (l *LikeSQLite) GetLikeByUserIDAndPostID(userID, postID int) (*models.Heart, error) {
	query := "SELECT id, user_id, post_id, type FROM likes WHERE user_id = ? AND post_id = ?"
	row := l.db.QueryRow(query, userID, postID)

	var like models.Heart
	err := row.Scan(&like.ID, &like.UserID, &like.PostID, &like.LikeType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &like, nil
}

func (l *LikeSQLite) GetLikesByUserID(userID int) ([]models.Heart, error) {
	query := "SELECT id, user_id, post_id, type FROM likes WHERE user_id = ?"
	rows, err := l.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var likes []models.Heart
	for rows.Next() {
		var like models.Heart
		err := rows.Scan(&like.ID, &like.UserID, &like.PostID, &like.LikeType)
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
