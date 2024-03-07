package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type LikeService struct {
	repo repository.Like
}

func NewLikeService(repo repository.Like) *LikeService {
	return &LikeService{repo: repo}
}

func (l *LikeService) CreateHeart(heart *models.Heart) error {
	return l.repo.CreateHeart(heart)
}

func (l *LikeService) DeleteHeart(postID, userID int) error {
	return l.repo.DeleteHeart(postID, userID)
}

func (l *LikeService) GetLikesByPostID(postID int) (*models.Like, error) {
	return l.repo.GetLikesByPostID(postID)
}

func (l *LikeService) GetLikeByUserIDAndPostID(userID, postID int) (*models.Heart, error) {
	return l.repo.GetLikeByUserIDAndPostID(userID, postID)
}

func (l *LikeService) GetLikesByUserID(userID int) ([]models.Heart, error) {
	return l.repo.GetLikesByUserID(userID)
}
