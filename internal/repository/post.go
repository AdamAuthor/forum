package repository

import (
	"database/sql"
	"forum/internal/models"
	"strings"
)

type PostSQLite struct {
	db *sql.DB
}

func NewPostSQLite(db *sql.DB) *PostSQLite {
	return &PostSQLite{db: db}
}

func (p *PostSQLite) CreatePost(post *models.Post) error {
	query := "INSERT INTO posts (title, content, user_id, created_at) VALUES (?, ?, ?, ?)"
	result, err := p.db.Exec(query, post.Title, post.Content, post.AuthorID, post.CreatedAt)
	if err != nil {
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = int(postID)

	return nil
}

func (p *PostSQLite) GetAllPosts() (*[]models.Post, error) {
	query := `
	SELECT posts.id, posts.title, posts.content, posts.user_id, users.username, posts.created_at, GROUP_CONCAT(categories.name)
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
	LEFT JOIN posts_categories ON posts.id = posts_categories.post_id
	LEFT JOIN categories ON posts_categories.category_id = categories.id
	GROUP BY posts.id
`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var username string
		var categories string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &username, &post.CreatedAt, &categories)
		if err != nil {
			return nil, err
		}
		post.AuthorName = username
		post.Categories = strings.Split(categories, ",")
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (p *PostSQLite) GetPostByID(id int) (*models.Post, error) {
	query := `
	SELECT posts.id, posts.title, posts.content, posts.user_id, users.username, posts.created_at, GROUP_CONCAT(categories.name)
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
	LEFT JOIN posts_categories ON posts.id = posts_categories.post_id
	LEFT JOIN categories ON posts_categories.category_id = categories.id
	WHERE posts.id = ?
	GROUP BY posts.id
`
	row := p.db.QueryRow(query, id)

	var post models.Post
	var username string
	var categories string
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &username, &post.CreatedAt, &categories)
	if err != nil {
		return nil, err
	}
	post.AuthorName = username
	post.Categories = strings.Split(categories, ",")

	return &post, nil
}

func (p *PostSQLite) GetPostByIDAndUserID(postID, userID int) (*models.Post, error) {
	query := `
        SELECT
            posts.id,
            posts.title,
            posts.content,
            posts.user_id,
            users.username,
            posts.created_at,
            GROUP_CONCAT(categories.name) AS category_names
        FROM
            posts
        INNER JOIN
            users ON posts.user_id = users.id
        LEFT JOIN
            posts_categories ON posts.id = posts_categories.post_id
        LEFT JOIN
            categories ON posts_categories.category_id = categories.id
        WHERE
            posts.id = ? AND
            posts.user_id = ?
        GROUP BY
            posts.id
    `
	row := p.db.QueryRow(query, postID, userID)

	var post models.Post
	var username string
	var categories string
	err := row.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &username, &post.CreatedAt, &categories)
	if err != nil {
		return nil, err
	}
	post.AuthorName = username
	post.Categories = strings.Split(categories, ",")

	return &post, nil
}

func (p *PostSQLite) DeletePost(userID, postID int) error {
	query := "DELETE FROM posts WHERE user_id = ? AND id = ?"
	_, err := p.db.Exec(query, userID, postID)
	return err
}

func (p *PostSQLite) GetAllPostsByUserID(id int) (*[]models.Post, error) {
	query := `
	SELECT posts.id, posts.title, posts.content, posts.user_id, users.username, posts.created_at, GROUP_CONCAT(categories.name)
	FROM posts
	INNER JOIN users ON posts.user_id = users.id
	LEFT JOIN posts_categories ON posts.id = posts_categories.post_id
	LEFT JOIN categories ON posts_categories.category_id = categories.id
	WHERE users.id = ?
	GROUP BY posts.id
`
	rows, err := p.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var username string
		var categories string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &username, &post.CreatedAt, &categories)
		if err != nil {
			return nil, err
		}
		post.AuthorName = username
		post.Categories = strings.Split(categories, ",")
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}

func (p *PostSQLite) AddCategoryToPost(postID int, category string) error {
	query := "INSERT INTO posts_categories (post_id, category_id) VALUES (?, (SELECT id FROM categories WHERE name = ?))"
	_, err := p.db.Exec(query, postID, category)
	return err
}

func (p *PostSQLite) GetAllPostsByCategoryName(name string) (*[]models.Post, error) {
	query := `
	SELECT 
    posts.id, 
    posts.title, 
    posts.content, 
    posts.user_id, 
    users.username, 
    posts.created_at, 
    GROUP_CONCAT(categories.name) AS category_names
FROM 
    posts
INNER JOIN 
    users ON posts.user_id = users.id
LEFT JOIN 
    posts_categories ON posts.id = posts_categories.post_id
LEFT JOIN 
    categories ON posts_categories.category_id = categories.id
WHERE 
    posts.id IN (SELECT post_id FROM posts_categories WHERE category_id = (SELECT id FROM categories WHERE name = ?))
GROUP BY 
    posts.id

`
	rows, err := p.db.Query(query, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		var username string
		var categories string
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.AuthorID, &username, &post.CreatedAt, &categories)
		if err != nil {
			return nil, err
		}
		post.AuthorName = username
		post.Categories = strings.Split(categories, ",")
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &posts, nil
}
