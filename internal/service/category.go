package service

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (c *CategoryService) CreateCategory(category *models.Category) error {
	return c.repo.CreateCategory(category)
}

func (c *CategoryService) GetAllCategories() (*[]models.Category, error) {
	return c.repo.GetAllCategories()
}
