package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"time"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(post *models.Post) error {
	post.CreatedAt = time.Now()
	return p.repo.CreatePost(post)
}

func (p *PostService) GetAllPosts() (*[]models.Post, error) {
	return p.repo.GetAllPosts()
}

func (p *PostService) GetPostByID(id int) (*models.Post, error) {
	return p.repo.GetPostByID(id)
}

func (p *PostService) GetPostByIDAndUserID(postID, userID int) (*models.Post, error) {
	return p.repo.GetPostByIDAndUserID(postID, userID)
}

func (p *PostService) DeletePost(userID, postID int) error {
	return p.repo.DeletePost(userID, postID)
}

func (p *PostService) GetAllPostsByUserID(id int) (*[]models.Post, error) {
	return p.repo.GetAllPostsByUserID(id)
}

func (p *PostService) AddCategoriesToPost(postID int, categories []string) error {
	for _, category := range categories {
		err := p.repo.AddCategoryToPost(postID, category)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PostService) GetAllPostsByCategory(name string) (*[]models.Post, error) {
	return p.repo.GetAllPostsByCategoryName(name)
}
