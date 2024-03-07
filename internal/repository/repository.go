package repository

import (
	"database/sql"
	"forum/internal/models"
)

type Authorization interface {
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(userID int) (*models.User, error)
	CreateSession(session *models.Session) error
	GetSessionByToken(token string) (*models.Session, error)
	DeleteAllSessionsByUserID(userID int) error
	ClearExpiredSessions() error
}

type Post interface {
	CreatePost(post *models.Post) error
	GetAllPosts() (*[]models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	GetPostByIDAndUserID(postID, userID int) (*models.Post, error)
	GetAllPostsByUserID(id int) (*[]models.Post, error)
	DeletePost(userID, postID int) error
	AddCategoryToPost(postID int, category string) error
	GetAllPostsByCategoryName(name string) (*[]models.Post, error)
}

type Category interface {
	CreateCategory(post *models.Category) error
	GetAllCategories() (*[]models.Category, error)
}

type Comment interface {
	CreateComment(comment *models.Comment) error
	GetCommentsByPostID(id int) (*[]models.Comment, error)
}

type Like interface {
	CreateHeart(heart *models.Heart) error
	DeleteHeart(postID, userID int) error
	GetLikesByPostID(postID int) (*models.Like, error)
	GetLikeByUserIDAndPostID(userID, postID int) (*models.Heart, error)
	GetLikesByUserID(userID int) ([]models.Heart, error)
}

type CommentLike interface {
	CreateCommentHeart(heart *models.CommentHeart) error
	DeleteCommentHeart(commentID, userID int) error
	GetCommentLikesByCommentID(commentID int) (*models.Like, error)
	GetCommentLikeByUserIDAndCommentID(userID, commentID int) (*models.CommentHeart, error)
	GetCommentLikesByUserID(userID int) ([]models.CommentHeart, error)
}

type Repository struct {
	Authorization
	Post
	Category
	Comment
	Like
	CommentLike
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSQLite(db),
		Post:          NewPostSQLite(db),
		Category:      NewCategorySQLite(db),
		Comment:       NewCommentSQLite(db),
		Like:          NewLikeSQLite(db),
		CommentLike:   NewCommentLikeSQLite(db),
	}
}
