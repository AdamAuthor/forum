package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) likePost(w http.ResponseWriter, r *http.Request) {
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

	var like models.Heart
	err = json.NewDecoder(r.Body).Decode(&like)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if like.LikeType != models.LikeTypes.Like && like.LikeType != models.LikeTypes.Dislike {
		http.Error(w, "Invalid like type", http.StatusBadRequest)
		return
	}

	like.UserID = session.UserID

	err = h.services.CreateHeart(&like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Like created successfully")
}

func (h *Handler) unlikePost(w http.ResponseWriter, r *http.Request) {
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

	h.services.DeleteHeart(postIDInt, session.UserID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Like removed successfully")
}

func (h *Handler) getPostLikes(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	like, err := h.services.GetLikesByPostID(postIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(like)
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

func (h *Handler) isLiked(w http.ResponseWriter, r *http.Request) {
	postID := r.URL.Query().Get("post_id")
	if postID == "" {
		http.Error(w, "Post ID is required", http.StatusBadRequest)
		return
	}
	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
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

	like, err := h.services.GetLikeByUserIDAndPostID(session.UserID, postIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse, err := json.Marshal(like)
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
