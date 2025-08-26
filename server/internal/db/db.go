package db


import (
"context"
"time"


"github.com/jackc/pgx/v5/pgxpool"
)


func Open(dsn string) (*pgxpool.Pool, error) {
cfg, err := pgxpool.ParseConfig(dsn)
if err != nil { return nil, err }
cfg.MaxConns = 10
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
pool, err := pgxpool.NewWithConfig(ctx, cfg)
if err != nil { return nil, err }
// Ping
if err := pool.Ping(ctx); err != nil { return nil, err }
return pool, nil
}