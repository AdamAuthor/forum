package handlers

import (
	"encoding/json"
	"fmt"
	"forum/internal/models"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) likeComment(w http.ResponseWriter, r *http.Request) {
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

	var like models.CommentHeart
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

	err = h.services.CreateCommentHeart(&like)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Like created successfully")
}

func (h *Handler) unlikeComment(w http.ResponseWriter, r *http.Request) {
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

	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}

	// Convert postID to integer
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	h.services.DeleteCommentHeart(commentIDInt, session.UserID)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Like removed successfully")
}

func (h *Handler) getCommentLikes(w http.ResponseWriter, r *http.Request) {
	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}
	like, err := h.services.GetCommentLikesByCommentID(commentIDInt)
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

func (h *Handler) isCommentLiked(w http.ResponseWriter, r *http.Request) {
	commentID := r.URL.Query().Get("comment_id")
	if commentID == "" {
		http.Error(w, "Comment ID is required", http.StatusBadRequest)
		return
	}
	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
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

	like, err := h.services.GetCommentLikeByUserIDAndCommentID(session.UserID, commentIDInt)
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
