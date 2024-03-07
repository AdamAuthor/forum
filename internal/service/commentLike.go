package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CommentLikeService struct {
	repo repository.CommentLike
}

func NewCommentLikeService(repo repository.CommentLike) *CommentLikeService {
	return &CommentLikeService{repo: repo}
}

func (c *CommentLikeService) CreateCommentHeart(heart *models.CommentHeart) error {
	return c.repo.CreateCommentHeart(heart)
}

func (c *CommentLikeService) DeleteCommentHeart(commentID, userID int) error {
	return c.repo.DeleteCommentHeart(commentID, userID)
}

func (c *CommentLikeService) GetCommentLikesByCommentID(commentID int) (*models.Like, error) {
	return c.repo.GetCommentLikesByCommentID(commentID)
}

func (c *CommentLikeService) GetCommentLikeByUserIDAndCommentID(userID, commentID int) (*models.CommentHeart, error) {
	return c.repo.GetCommentLikeByUserIDAndCommentID(userID, commentID)
}

func (c *CommentLikeService) GetCommentLikesByUserID(userID int) ([]models.CommentHeart, error) {
	return c.repo.GetCommentLikesByUserID(userID)
}
