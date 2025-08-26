package handlers


import (
"context"
"net/http"
"strconv"
"time"


"github.com/jackc/pgx/v5/pgxpool"
)


type MessageDTO struct {
ID int64 `json:"id"`
UserID int64 `json:"userId"`
DisplayName string `json:"displayName"`
Content string `json:"content"`
CreatedAt time.Time `json:"createdAt"`
}


func ListMessagesHandler(pool *pgxpool.Pool) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
limit := 50
if q := r.URL.Query().Get("limit"); q != "" {
if n, err := strconv.Atoi(q); err == nil && n > 0 && n <= 200 { limit = n }
}
ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
defer cancel()
rows, err := pool.Query(ctx, `
SELECT m.id, m.user_id, u.display_name, m.content, m.created_at
FROM messages m
JOIN users u ON u.id = m.user_id
ORDER BY m.created_at DESC
LIMIT $1`, limit)
if err != nil { writeJSON(w, http.StatusInternalServerError, jsonResp{"error":"db"}); return }
defer rows.Close()
var list []MessageDTO
for rows.Next() {
var d MessageDTO
if err := rows.Scan(&d.ID, &d.UserID, &d.DisplayName, &d.Content, &d.CreatedAt); err == nil {
list = append(list, d)
}
}
// Return ascending order (oldest first) for UI
for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 { list[i], list[j] = list[j], list[i] }
writeJSON(w, http.StatusOK, list)
})
}