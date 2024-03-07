package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"time"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment *models.Comment) error {
	comment.CreatedAt = time.Now()
	return s.repo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostID(postID int) (*[]models.Comment, error) {
	return s.repo.GetCommentsByPostID(postID)
}
