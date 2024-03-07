package repository

import (
	"database/sql"
	"forum/internal/models"
)

type CategorySQLite struct {
	db *sql.DB
}

func NewCategorySQLite(db *sql.DB) *CategorySQLite {
	return &CategorySQLite{db: db}
}

func (c *CategorySQLite) CreateCategory(category *models.Category) error {
	query := "INSERT INTO categories (name) VALUES (?)"
	result, err := c.db.Exec(query, category.Name)
	if err != nil {
		return err
	}

	categoryID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = int(categoryID)

	return nil
}

func (c *CategorySQLite) GetAllCategories() (*[]models.Category, error) {
	query := "SELECT id, name FROM categories"
	rows, err := c.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &categories, nil
}
