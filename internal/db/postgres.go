package db

import (
	"bmp-tgbot/internal/sdk/models"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresClient struct {
	conn   *sql.Conn
	logger *zap.Logger
}

func NewPostgresClient(ctx context.Context, connStr string) *PostgresClient {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	conn, err := db.Conn(ctx)
	if err != nil {
		panic(err)
	}

	_, err = conn.ExecContext(ctx, queryInitUsers)
	if err != nil {
		panic(err)
	}

	return &PostgresClient{
		conn: conn,
	}
}

type Client interface {
	GetUser(ctx context.Context, user *models.User) error
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetLeaderboard(ctx context.Context) (map[string]int64, error)
	GetUsers(ctx context.Context) string
}

func (r *PostgresClient) GetUser(ctx context.Context, user *models.User) error {
	row := r.conn.QueryRowContext(ctx, queryUser, user.ID)

	if err := row.Scan(&user.ID, &user.Username, &user.State, &user.Balance); err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) CreateUser(ctx context.Context, user *models.User) error {
	user.Balance = 100
	_, err := r.conn.ExecContext(ctx, insertUser, user.ID, user.Username, user.State, user.Balance)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.conn.ExecContext(ctx, updateUser, user.State, user.Balance, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresClient) GetLeaderboard(ctx context.Context) (map[string]int64, error) {
	out := make(map[string]int64)
	rows, err := r.conn.QueryContext(ctx, queryLeaderboard)
	fmt.Println(rows)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		temp := models.User{}

		if err := rows.Scan(&temp.ID, &temp.Username, &temp.Balance); err != nil {
			r.logger.Error("failed to scan row", zap.Error(err))
			continue
		}

		out["@"+temp.Username] = temp.Balance
	}
	fmt.Println("LEADERBOARD:", out)
	return out, nil
}

//
//func (r *PostgresClient) GetUsers(ctx context.Context) string {
//	out := ""
//	rows, err := r.conn.QueryRowContext(ctx, queryUsers)
//}
