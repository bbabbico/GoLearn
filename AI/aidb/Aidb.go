package aidb

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
	ID    int
	Name  string
	Email string
	Age   int
}

func Init() error {
	dsn := "root:ehgus2003@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	db, err = NewDB(dsn)
	if err != nil {
		return fmt.Errorf("DB 커넥션 풀 생성 실패: %w", err)
	}

	fmt.Println("✓ MySQL 연결 성공")
	fmt.Println("✓ DB 커넥션 풀 생성 성공")

	printDBStats(db)
	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func Load() error {
	user, err := getUserByID(db, 1)
	if err != nil {
		return fmt.Errorf("DB 조회 실패: %w", err)
	}

	fmt.Printf("회원정보\nID : %d\nName : %s\nEmail : %s\nAge : %d\n",
		user.ID, user.Name, user.Email, user.Age)
	return nil
}

// NewDB 커넥션 풀 관리
func NewDB(dsn string) (*sql.DB, error) {
	newdb, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	newdb.SetMaxOpenConns(25)
	newdb.SetMaxIdleConns(10)
	newdb.SetConnMaxLifetime(time.Minute * 3)
	newdb.SetConnMaxIdleTime(time.Minute)

	if err := newdb.Ping(); err != nil {
		newdb.Close()
		return nil, err
	}

	return newdb, nil
}

func printDBStats(db *sql.DB) {
	stats := db.Stats()
	fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
	fmt.Printf("InUse: %d\n", stats.InUse)
	fmt.Printf("Idle: %d\n", stats.Idle)
	fmt.Printf("WaitCount: %d\n", stats.WaitCount)
}

func getUserByID(db *sql.DB, id int) (*User, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	var u User
	err := db.QueryRow(
		"SELECT id, name, email, age FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Age)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return nil, fmt.Errorf("user id=%d not found", id)
	case err != nil:
		return nil, fmt.Errorf("getUserByID: %w", err)
	}

	return &u, nil
}
