package handlers

import (
	"encoding/json"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := h.services.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to extract userID from token", http.StatusUnauthorized)
		return
	}

	var post models.Post
	err = json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post.AuthorID = session.UserID

	err = h.services.CreatePost(&post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.services.AddCategoriesToPost(post.ID, post.Categories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) readPostById(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.services.GetPostByID(postIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(post)
	if err != nil {
		http.Error(w, "Failed to marshal post to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}

func (h *Handler) deletePost(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := h.services.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to extract userID from token", http.StatusUnauthorized)
		return
	}

	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	// Convert postID to integer
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	_, err = h.services.GetPostByID(postIDInt)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	_, err = h.services.GetPostByIDAndUserID(postIDInt, session.UserID)
	if err != nil {
		http.Error(w, "Post does not belong to the user", http.StatusForbidden)
		return
	}

	err = h.services.DeletePost(session.UserID, postIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) readAllPostsByUserID(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := h.services.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to extract userID from token", http.StatusUnauthorized)
		return
	}

	userID := session.UserID

	postsPtr, likes, err := h.fetchMyPostsAndLikes(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PostWithLikes struct {
		Post  models.Post  `json:"post"`
		Likes *models.Like `json:"likes"`
	}

	posts := *postsPtr
	var postsWithLikes []PostWithLikes

	for _, post := range posts {
		postWithLikes := PostWithLikes{
			Post:  post,
			Likes: likes[post.ID],
		}
		postsWithLikes = append(postsWithLikes, postWithLikes)
	}

	jsonResponse, err := json.Marshal(postsWithLikes)
	if err != nil {
		http.Error(w, "Failed to marshal posts with likes to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}

func (h *Handler) readAllPostsByLikes(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	session, err := h.services.GetSessionByToken(cookie.Value)
	if err != nil {
		http.Error(w, "Failed to extract userID from token", http.StatusUnauthorized)
		return
	}

	postsPtr, likes, err := h.fetchLikedPostsAndLikes(session.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PostWithLikes struct {
		Post  models.Post  `json:"post"`
		Likes *models.Like `json:"likes"`
	}

	posts := *postsPtr
	var postsWithLikes []PostWithLikes

	for _, post := range posts {
		postWithLikes := PostWithLikes{
			Post:  post,
			Likes: likes[post.ID],
		}
		postsWithLikes = append(postsWithLikes, postWithLikes)
	}

	jsonResponse, err := json.Marshal(postsWithLikes)
	if err != nil {
		http.Error(w, "Failed to marshal posts with likes to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}

func (h *Handler) fetchLikedPostsAndLikes(userID int) (*[]models.Post, map[int]*models.Like, error) {
	likes, err := h.services.GetLikesByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	var likedPostIDs []int
	for _, like := range likes {
		if like.LikeType == models.LikeTypes.Like {
			likedPostIDs = append(likedPostIDs, like.PostID)
		}
	}

	posts := make([]models.Post, 0)
	for _, postID := range likedPostIDs {
		post, err := h.services.GetPostByID(postID)
		if err != nil {
			continue
		}
		posts = append(posts, *post)
	}

	likesMap := make(map[int]*models.Like)

	for _, post := range posts {
		like, err := h.services.GetLikesByPostID(post.ID)
		if err != nil {
			return nil, nil, err
		}
		likesMap[post.ID] = like
	}

	return &posts, likesMap, nil
}

func (h *Handler) fetchMyPostsAndLikes(userID int) (*[]models.Post, map[int]*models.Like, error) {
	postsPtr, err := h.services.GetAllPostsByUserID(userID)
	if err != nil {
		return nil, nil, err
	}

	likesMap := make(map[int]*models.Like)

	posts := *postsPtr

	for _, post := range posts {
		like, err := h.services.GetLikesByPostID(post.ID)
		if err != nil {
			return nil, nil, err
		}
		likesMap[post.ID] = like
	}

	return &posts, likesMap, nil
}

func (h *Handler) fetchPostsAndLikes() (*[]models.Post, map[int]*models.Like, error) {
	postsPtr, err := h.services.GetAllPosts()
	if err != nil {
		return nil, nil, err
	}

	posts := *postsPtr
	likes := make(map[int]*models.Like)

	for _, post := range posts {
		like, err := h.services.GetLikesByPostID(post.ID)
		if err != nil {
			return nil, nil, err
		}
		likes[post.ID] = like // Use address of like
	}

	return &posts, likes, nil
}

func (h *Handler) fetchPostsAndLikesByCategory(category string) (*[]models.Post, map[int]*models.Like, error) {
	postsPtr, err := h.services.GetAllPostsByCategory(category)
	if err != nil {
		return nil, nil, err
	}

	posts := *postsPtr
	likes := make(map[int]*models.Like)

	for _, post := range posts {
		like, err := h.services.GetLikesByPostID(post.ID)
		if err != nil {
			return nil, nil, err
		}
		likes[post.ID] = like // Use address of like
	}

	return &posts, likes, nil
}

func (h *Handler) displayPostsAndLikes(w http.ResponseWriter, r *http.Request) {
	postsPtr, likes, err := h.fetchPostsAndLikes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PostWithLikes struct {
		Post  models.Post  `json:"post"`
		Likes *models.Like `json:"likes"`
	}

	posts := *postsPtr
	var postsWithLikes []PostWithLikes

	for _, post := range posts {
		postWithLikes := PostWithLikes{
			Post:  post,
			Likes: likes[post.ID],
		}
		postsWithLikes = append(postsWithLikes, postWithLikes)
	}

	jsonResponse, err := json.Marshal(postsWithLikes)
	if err != nil {
		http.Error(w, "Failed to marshal posts with likes to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}

func (h *Handler) displayPostsAndLikesByCategory(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}

	postsPtr, likes, err := h.fetchPostsAndLikesByCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type PostWithLikes struct {
		Post  models.Post  `json:"post"`
		Likes *models.Like `json:"likes"`
	}

	posts := *postsPtr
	var postsWithLikes []PostWithLikes

	for _, post := range posts {
		postWithLikes := PostWithLikes{
			Post:  post,
			Likes: likes[post.ID],
		}
		postsWithLikes = append(postsWithLikes, postWithLikes)
	}

	jsonResponse, err := json.Marshal(postsWithLikes)
	if err != nil {
		http.Error(w, "Failed to marshal posts with likes to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}
