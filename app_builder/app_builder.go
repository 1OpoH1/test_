package app_builder

import (
	"anti-scam-trainer/backend/models"
	"anti-scam-trainer/backend/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-pg/pg"
	"github.com/lpernett/godotenv"
)

type App struct {
	DB     *pg.DB
	Router *http.ServeMux
	Port   string
}

func NewApp() (*App, error) {
	_ = godotenv.Load()

	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_NAME"),
	})

	if db == nil {
		return nil, fmt.Errorf("failed to connect to database")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app := &App{
		DB:     db,
		Router: http.NewServeMux(),
		Port:   port,
	}

	app.registerRoutes()
	return app, nil
}

func (a *App) Run() error {
	return http.ListenAndServe(":"+a.Port, a.Router)
}

func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}

func (a *App) registerRoutes() {
	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "AntiScamTrainer backend is running")
	})

	a.Router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	a.Router.HandleFunc("/api/lessons", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			idParam := r.URL.Query().Get("id")
			if idParam == "" {
				idParam = r.URL.Query().Get("lesson_id")
			}
			if idParam != "" {
				id, err := strconv.Atoi(idParam)
				if err != nil {
					http.Error(w, "invalid lesson id", http.StatusBadRequest)
					return
				}
				lesson, err := services.GetLessonByID(a.DB, id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				writeJSON(w, lesson)
				return
			}

			lessons, err := services.ListLessons(a.DB)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, lessons)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/lessons/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/lessons/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		parts := strings.Split(path, "/")
		lessonID, err := strconv.Atoi(parts[0])
		if err != nil {
			http.Error(w, "invalid lesson id", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			switch r.Method {
			case http.MethodGet:
				lesson, err := services.GetLessonByID(a.DB, lessonID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				writeJSON(w, lesson)
			default:
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		if len(parts) == 3 && parts[1] == "quiz" && parts[2] == "submit" {
			if r.Method != http.MethodPost {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			var payload struct {
				UserID int  `json:"user_id"`
				Score  int  `json:"score"`
				Passed bool `json:"passed"`
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if payload.UserID == 0 {
				http.Error(w, "user_id is required", http.StatusBadRequest)
				return
			}

			if err := services.SubmitLessonQuiz(a.DB, payload.UserID, lessonID, payload.Score, payload.Passed); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if len(parts) == 2 && parts[1] == "chats" {
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}

			userID := 0
			userIDParam := r.URL.Query().Get("user_id")
			if userIDParam != "" {
				userID, err = strconv.Atoi(userIDParam)
				if err != nil {
					http.Error(w, "invalid user_id", http.StatusBadRequest)
					return
				}
			}

			chats, err := services.ListLessonChats(a.DB, lessonID, userID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, chats)
			return
		}

		http.NotFound(w, r)
	})

	a.Router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			users, err := services.ListUsers(a.DB)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, users)
		case http.MethodPost:
			var user models.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			created, err := services.CreateUser(a.DB, &user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, created)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/users/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			user, err := services.GetUserByID(a.DB, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, user)
		case http.MethodPut:
			var user models.User
			if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			user.ID = id
			if err := services.UpdateUser(a.DB, &user); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, user)
		case http.MethodDelete:
			if err := services.DeleteUser(a.DB, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/chats", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			chats, err := services.ListChats(a.DB)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, chats)
		case http.MethodPost:
			var chat models.Chat
			if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			created, err := services.CreateChat(a.DB, &chat)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, created)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/chats/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/chats/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid chat id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			chat, err := services.GetChatByID(a.DB, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, chat)
		case http.MethodPut:
			var chat models.Chat
			if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			chat.ID = id
			if err := services.UpdateChat(a.DB, &chat); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, chat)
		case http.MethodDelete:
			if err := services.DeleteChat(a.DB, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/chat-sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			sessions, err := services.ListChatSessions(a.DB)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, sessions)
		case http.MethodPost:
			var session models.ChatSession
			if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			created, err := services.CreateChatSession(a.DB, &session)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, created)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	a.Router.HandleFunc("/api/chat-sessions/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/api/chat-sessions/")
		if path == "" {
			http.NotFound(w, r)
			return
		}

		id, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid chat session id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			session, err := services.GetChatSessionByID(a.DB, id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, session)
		case http.MethodPut:
			var session models.ChatSession
			if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			session.ID = id
			if err := services.UpdateChatSession(a.DB, &session); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			writeJSON(w, session)
		case http.MethodDelete:
			if err := services.DeleteChatSession(a.DB, id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func writeJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		fmt.Printf("failed to encode response: %v\n", err)
	}
}
