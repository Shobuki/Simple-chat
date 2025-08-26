package ws


import (
"log"
"net/http"
"time"


"github.com/gorilla/websocket"


"github.com/yourname/go-react-chat/internal/auth"
)


var upgrader = websocket.Upgrader{ CheckOrigin: func(r *http.Request) bool { return true } }


func ServeWS(jwt *auth.JWT, hub *Hub) http.Handler {
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
token := r.URL.Query().Get("token")
if token == "" { http.Error(w, "missing token", http.StatusUnauthorized); return }
claims, err := jwt.Parse(token)
if err != nil { http.Error(w, "invalid token", http.StatusUnauthorized); return }


conn, err := upgrader.Upgrade(w, r, nil)
if err != nil { log.Println("upgrade:", err); return }


client := &Client{ userID: claims.UserID, displayName: claims.DisplayName, hub: hub, send: make(chan Outgoing, 32) }
hub.register <- client


// Writer
go func() {
defer conn.Close()
conn.SetWriteDeadline(time.Now().Add(writeWait))
for msg := range client.send {
conn.SetWriteDeadline(time.Now().Add(writeWait))
if err := conn.WriteJSON(msg); err != nil { return }
}
}()


// Reader
go func() {
defer func(){ hub.unregister <- client; conn.Close() }()
conn.SetReadLimit(4096)
conn.SetReadDeadline(time.Now().Add(pongWait))
conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
for {
var incoming struct{ Content string `json:"content"` }
if err := conn.ReadJSON(&incoming); err != nil { return }
if incoming.Content == "" { continue }
if err := hub.SaveAndBroadcast(client.userID, client.displayName, incoming.Content); err != nil { log.Println("save:", err); return }
}
}()
})
}