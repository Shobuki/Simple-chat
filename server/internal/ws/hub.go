package ws

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Hub struct {
	pool       *pgxpool.Pool
	clients    map[*Client]bool
	broadcast  chan Outgoing
	register   chan *Client
	unregister chan *Client
}

type Outgoing struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userId"`
	DisplayName string    `json:"displayName"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"createdAt"`
}

func NewHub(pool *pgxpool.Pool) *Hub {
	return &Hub{
		pool:       pool,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan Outgoing, 128),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}

		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
				}
			}
		}
	}
}

func (h *Hub) SaveAndBroadcast(userID int64, displayName, content string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int64
	var created time.Time
	err := h.pool.QueryRow(ctx,
		`INSERT INTO messages(user_id, content) VALUES($1,$2) RETURNING id, created_at`,
		userID, content,
	).Scan(&id, &created)
	if err != nil {
		return err
	}

	h.broadcast <- Outgoing{
		ID:          id,
		UserID:      userID,
		DisplayName: displayName,
		Content:     content,
		CreatedAt:   created,
	}
	return nil
}
