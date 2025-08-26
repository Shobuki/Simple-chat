package main


import (
"log"
"net/http"
"os"
"time"


"github.com/joho/godotenv"


"github.com/yourname/go-react-chat/internal/auth"
"github.com/yourname/go-react-chat/internal/db"
"github.com/yourname/go-react-chat/internal/handlers"
"github.com/yourname/go-react-chat/internal/ws"
)


func main() {
_ = godotenv.Load()


dsn := os.Getenv("DATABASE_URL")
jwtSecret := os.Getenv("JWT_SECRET")
port := os.Getenv("PORT")
if port == "" { port = "8080" }
corsOrigin := os.Getenv("CORS_ORIGIN")
if corsOrigin == "" { corsOrigin = "http://localhost:5173" }


pool, err := db.Open(dsn)
if err != nil { log.Fatal(err) }
defer pool.Close()


jwt := auth.NewJWT(jwtSecret)


hub := ws.NewHub(pool)
go hub.Run()


mux := http.NewServeMux()


// REST endpoints
mux.Handle("/api/register", handlers.CORS(corsOrigin)(handlers.RegisterHandler(pool)))
mux.Handle("/api/login", handlers.CORS(corsOrigin)(handlers.LoginHandler(pool, jwt)))
mux.Handle("/api/me", handlers.CORS(corsOrigin)(handlers.AuthRequired(jwt, handlers.MeHandler(pool))))
mux.Handle("/api/messages", handlers.CORS(corsOrigin)(handlers.AuthRequired(jwt, handlers.ListMessagesHandler(pool))))


// WebSocket (token di query: /ws?token=...)
mux.Handle("/ws", handlers.CORS(corsOrigin)(ws.ServeWS(jwt, hub)))


srv := &http.Server{
Addr: ":" + port,
Handler: mux,
ReadHeaderTimeout: 10 * time.Second,
IdleTimeout: 60 * time.Second,
}


log.Printf("server listening on :%s", port)
if err := srv.ListenAndServe(); err != nil { log.Fatal(err) }
}