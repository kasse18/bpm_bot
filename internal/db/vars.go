package db

const (
	queryUser        = "SELECT * FROM users WHERE id=$1"
	insertUser       = "INSERT INTO users (id, username, state, balance) values ($1, $2, $3, $4)"
	updateUser       = "UPDATE users SET state=$1, balance=$2 WHERE id=$3"
	queryLeaderboard = "SELECT id, username, balance FROM users ORDER BY balance DESC"
	queryUsers       = "SELECT username FROM users"
)

const (
	queryInitUsers = `CREATE TABLE IF NOT EXISTS users (
		id bigint NOT NULL,
		username text NOT NULL,
		state text NOT NULL, 
		balance bigint NOT NULL DEFAULT 100,
		PRIMARY KEY (id)
	  )`
)
