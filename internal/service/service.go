package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"time"
)

type Authorization interface {
	CreateUser(user *models.User) error
	GetUserByID(userID int) (*models.User, error)
	CreateSession(username, password string) (*models.Session, int, error)
	GetSessionByToken(token string) (*models.Session, error)
	ClearExpiredSessionsPeriodically(interval time.Duration)
}

type Post interface {
	CreatePost(post *models.Post) error
	GetAllPosts() (*[]models.Post, error)
	GetPostByID(id int) (*models.Post, error)
	GetPostByIDAndUserID(postID, userID int) (*models.Post, error)
	DeletePost(userID, postID int) error
	GetAllPostsByUserID(id int) (*[]models.Post, error)
	AddCategoriesToPost(postID int, categories []string) error
	GetAllPostsByCategory(name string) (*[]models.Post, error)
}

type Category interface {
	CreateCategory(category *models.Category) error
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

type Service struct {
	Authorization
	Post
	Category
	Comment
	Like
	CommentLike
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Post:          NewPostService(repos.Post),
		Category:      NewCategoryService(repos.Category),
		Comment:       NewCommentService(repos.Comment),
		Like:          NewLikeService(repos.Like),
		CommentLike:   NewCommentLikeService(repos.CommentLike),
	}
}
