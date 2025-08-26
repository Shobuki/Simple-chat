package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourname/go-react-chat/internal/auth"
)

type jsonResp map[string]any

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

// CORS middleware
func CORS(origin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// auth context keys
type ctxKeyUserID struct{}
type ctxKeyDisplayName struct{}

func UserIDFromCtx(r *http.Request) int64 {
	v := r.Context().Value(ctxKeyUserID{})
	if v == nil {
		return 0
	}
	return v.(int64)
}

func DisplayNameFromCtx(r *http.Request) string {
	v := r.Context().Value(ctxKeyDisplayName{})
	if v == nil {
		return ""
	}
	return v.(string)
}

// bearer auth middleware
func AuthRequired(jwt *auth.JWT, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			writeJSON(w, http.StatusUnauthorized, jsonResp{"error": "missing bearer"})
			return
		}
		token := strings.TrimPrefix(h, "Bearer ")
		claims, err := jwt.Parse(token)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, jsonResp{"error": "invalid token"})
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyUserID{}, claims.UserID)
		ctx = context.WithValue(ctx, ctxKeyDisplayName{}, claims.DisplayName)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// --- handlers ---

func RegisterHandler(pool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, jsonResp{"error": "method"})
			return
		}
		var in struct {
			Email       string `json:"email"`
			Password    string `json:"password"`
			DisplayName string `json:"displayName"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, jsonResp{"error": "bad json"})
			return
		}
		if in.Email == "" || in.Password == "" || in.DisplayName == "" {
			writeJSON(w, http.StatusBadRequest, jsonResp{"error": "missing fields"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, jsonResp{"error": "hash"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var id int64
		err = pool.QueryRow(ctx,
			`INSERT INTO users(email, display_name, password_hash) VALUES($1,$2,$3) RETURNING id`,
			in.Email, in.DisplayName, string(hash),
		).Scan(&id)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				writeJSON(w, http.StatusConflict, jsonResp{"error": "email exists"})
				return
			}
			writeJSON(w, http.StatusInternalServerError, jsonResp{"error": "db"})
			return
		}
		writeJSON(w, http.StatusCreated, jsonResp{"id": id})
	})
}

func LoginHandler(pool *pgxpool.Pool, jwt *auth.JWT) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, jsonResp{"error": "method"})
			return
		}
		var in struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeJSON(w, http.StatusBadRequest, jsonResp{"error": "bad json"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var id int64
		var displayName, hash string
		err := pool.QueryRow(ctx,
			`SELECT id, display_name, password_hash FROM users WHERE email=$1`,
			in.Email,
		).Scan(&id, &displayName, &hash)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, jsonResp{"error": "invalid credentials"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(hash), []byte(in.Password)) != nil {
			writeJSON(w, http.StatusUnauthorized, jsonResp{"error": "invalid credentials"})
			return
		}

		token, err := jwt.Sign(id, displayName)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, jsonResp{"error": "jwt"})
			return
		}
		writeJSON(w, http.StatusOK, jsonResp{"token": token, "displayName": displayName})
	})
}

func MeHandler(pool *pgxpool.Pool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := UserIDFromCtx(r)
		if id == 0 {
			writeJSON(w, http.StatusUnauthorized, jsonResp{"error": "no ctx"})
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		var email, display string
		err := pool.QueryRow(ctx,
			`SELECT email, display_name FROM users WHERE id=$1`, id,
		).Scan(&email, &display)
		if err != nil {
			writeJSON(w, http.StatusNotFound, jsonResp{"error": "not found"})
			return
		}
		writeJSON(w, http.StatusOK, jsonResp{
			"id":          id,
			"email":       email,
			"displayName": display,
		})
	})
}
