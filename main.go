package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type User struct {
	UserID    string
	UserName  string
	CreatedAt time.Time
}

func (u *User) String() string {
	return fmt.Sprintf("UserID: %s, UserName: %s, CreatedAt: %s", u.UserID, u.UserName, u.CreatedAt)
}

func main() {
	db, err := sql.Open("pgx", "host=localhost port=5432 user=testuser dbname=testdb password=pass sslmode=disable")
	if err != nil {
		log.Fatalf("main sql.Open error err:%v", err)
	}
	defer db.Close()

	ctx := context.Background()

	rows, err := db.QueryContext(ctx, `SELECT user_id, user_name, created_at FROM users ORDER BY user_id;`)
	if err != nil {
		log.Fatalf("query all users: %v", err)
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var (
			userID, userName string
			createdAt        time.Time
		)

		if err := rows.Scan(&userID, &userName, &createdAt); err != nil {
			log.Fatalf("scan the user: %v", err)
		}

		users = append(users, &User{
			UserID:    userID,
			UserName:  userName,
			CreatedAt: createdAt,
		})
	}

	if err := rows.Close(); err != nil {
		log.Fatalf("rows close: %v", err)
	}

	if err := rows.Err(); err != nil {
		log.Fatalf("scan users: %v", err)
	}

	for _, user := range users {
		fmt.Printf("%v\n", user)
	}
}
