package ws


import (
"time"
)


type Client struct {
userID int64
displayName string
hub *Hub
send chan Outgoing
}


const (
writeWait = 10 * time.Second
pongWait = 60 * time.Second
)