package handlers

import (
	"forum/internal/middleware"
	"forum/internal/service"
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	router := http.NewServeMux()

	indexTemp, err := template.ParseFiles("./web/public/index.html")
	if err != nil {
		log.Fatalln(err)
	}

	fileServer := http.FileServer(http.Dir("./web/src"))

	router.Handle("/src/", http.StripPrefix("/src/", fileServer))
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := indexTemp.Execute(w, "localhost:8080")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	})

	router.HandleFunc("/register", middleware.CorsMiddleware(h.registerUser))
	router.HandleFunc("/login", middleware.CorsMiddleware(h.loginUser))
	router.HandleFunc("/logout", middleware.CorsMiddleware(h.logoutUser))
	router.HandleFunc("/is-login", middleware.CorsMiddleware(h.isLogin))

	router.HandleFunc("/create-post", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.createPost, h.services)))
	router.HandleFunc("/posts", middleware.CorsMiddleware(h.displayPostsAndLikes))
	router.HandleFunc("/my-posts", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.readAllPostsByUserID, h.services)))
	router.HandleFunc("/post", middleware.CorsMiddleware(h.readPostById))
	router.HandleFunc("/posts-by-category", middleware.CorsMiddleware(h.displayPostsAndLikesByCategory))
	router.HandleFunc("/posts-by-like", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.readAllPostsByLikes, h.services)))
	router.HandleFunc("/delete-post", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.deletePost, h.services)))
	router.HandleFunc("/create-category", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.createCategory, h.services)))
	router.HandleFunc("/categories", middleware.CorsMiddleware(h.getAllCategories))

	router.HandleFunc("/create-comment", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.createComment, h.services)))
	router.HandleFunc("/comments", middleware.CorsMiddleware(h.readAllCommentsByPostID))
	router.HandleFunc("/likes", middleware.CorsMiddleware(h.getPostLikes))
	router.HandleFunc("/like-post", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.likePost, h.services)))
	router.HandleFunc("/unlike-post", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.unlikePost, h.services)))
	router.HandleFunc("/is-liked", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.isLiked, h.services)))

	router.HandleFunc("/comment-likes", middleware.CorsMiddleware(h.getCommentLikes))
	router.HandleFunc("/like-comment", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.likeComment, h.services)))
	router.HandleFunc("/unlike-comment", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.unlikeComment, h.services)))
	router.HandleFunc("/is-comment-liked", middleware.CorsMiddleware(middleware.AuthMiddleWare(h.isCommentLiked, h.services)))
	return router
}
