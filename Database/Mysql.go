package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	// 드라이버를 blank import로 등록 (init() 함수 실행)
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB // 데이터베이스 핸들, 커넥션 풀 관리

type User struct { // DB 전용 유저 구조체
	ID    int
	Name  string
	Email string
	Age   int
}

func main() {
	var err error
	var user *User
	var id int64

	// DSN 형식: user:password@tcp(host:port)/dbname?params
	dsn := "root:ehgus2003@tcp(127.0.0.1:3306)/godb?charset=utf8mb4&parseTime=True&loc=Local"

	// sql.Open은 실제 연결하지 않고 DB 핸들만 생성
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Ping으로 실제 연결 확인
	if err = db.Ping(); err != nil {
		log.Fatal("DB 연결 실패:", err)
	}
	fmt.Println("✓ MySQL 연결 성공")

	// DB 커넥션 풀 내부적으로 관리함
	if _, err = NewDB(dsn); err != nil {
		log.Fatal("DB 커넥션 풀 생성 실패", err)
	}
	fmt.Println("✓ DB 커넥션 풀 생성 성공")

	// DB 상태 모니터링
	printDBStats(db)

	// 삽입
	id, err = createUser(db, "qwe", "qwe@@@", 10)
	if err != nil {
		log.Fatal("INSERT 실패:", err)
	}
	fmt.Printf("삽입한 회원의 ID : %d\n", id)

	// 단일 SELECT
	if user, err = getUserByID(db, int(id)); err != nil {
		log.Fatal("DB 조회 결과가 없습니다:", err)
	}
	fmt.Printf("회원정보\nID : %d\nName : %s\nEmail : %s\nAge : %d\n", user.ID, user.Name, user.Email, user.Age)

	// 트랜젝션
	err = transfer(db, 1, 2, 1000)
	if err != nil {
		log.Fatal(err)
	}
}

// NewDB 커넥션 풀 관리
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// 최대 열린 커넥션 수 (MySQL max_connections 미만으로 설정)
	db.SetMaxOpenConns(25)

	// 최대 유휴(idle) 커넥션 수
	db.SetMaxIdleConns(10)

	// 커넥션 최대 재사용 시간 (MySQL wait_timeout 이하로 설정)
	db.SetConnMaxLifetime(time.Minute * 3)

	// 유휴 커넥션 최대 대기 시간 (Go 1.15+)
	db.SetConnMaxIdleTime(time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// 풀 상태 모니터링
func printDBStats(db *sql.DB) {
	stats := db.Stats()
	fmt.Printf("OpenConnections: %d\n", stats.OpenConnections)
	fmt.Printf("InUse: %d\n", stats.InUse)
	fmt.Printf("Idle: %d\n", stats.Idle)
	fmt.Printf("WaitCount: %d\n", stats.WaitCount)
}

// 다중 행 SELECT
func getUsers(db *sql.DB) ([]User, error) {
	// ? 플레이스홀더로 SQL Injection 방지
	rows, err := db.Query(
		"SELECT id, name, email, age FROM users WHERE age > ? ORDER BY id",
		18,
	)
	if err != nil {
		return nil, fmt.Errorf("getUsers: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows) // 반드시 닫아야 커넥션이 풀로 반환됨

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Age); err != nil {
			return nil, fmt.Errorf("getUsers scan: %w", err)
		}
		users = append(users, u)
	}

	// 반복 중 발생한 에러 확인 (필수!)
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("getUsers rows: %w", err)
	}
	return users, nil
}

// 단일 행 SELECT
func getUserByID(db *sql.DB, id int) (*User, error) {
	var u User

	// QueryRow는 항상 *sql.Row를 반환 (에러는 Scan에서 확인)
	err := db.QueryRow(
		"SELECT id, name, email, age FROM users WHERE id = ?",
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Age)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		// 결과 없음 — nil 반환 또는 커스텀 에러
		return nil, fmt.Errorf("user id=%d not found", id)
	case err != nil:
		return nil, fmt.Errorf("getUserByID: %w", err)
	}
	return &u, nil
}

// INSERT — LastInsertId 활용
func createUser(db *sql.DB, name string, email string, age int) (int64, error) {
	result, err := db.Exec( // 단일 실행
		"INSERT INTO users (name, email, age) VALUES (?, ?, ?)",
		name, email, age,
	)
	if err != nil {
		return 0, fmt.Errorf("createUser: %w", err)
	}

	id, err := result.LastInsertId() // AUTO_INCREMENT ID
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 다중 삽입
func bulkInsertUsers(db *sql.DB, users []User) error {
	if len(users) == 0 {
		return nil
	}

	// VALUES 플레이스홀더 동적 생성
	valueStrings := make([]string, 0, len(users))
	valueArgs := make([]interface{}, 0, len(users)*3)

	for _, u := range users {
		valueStrings = append(valueStrings, "(?, ?, ?)")
		valueArgs = append(valueArgs, u.Name, u.Email, u.Age)
	}

	query := fmt.Sprintf(
		"INSERT INTO users (name, email, age) VALUES %s",
		strings.Join(valueStrings, ","),
	)

	_, err := db.Exec(query, valueArgs...)
	return err
}

// UPDATE — RowsAffected 활용
func updateUserEmail(db *sql.DB, id int, newEmail string) (int64, error) {
	result, err := db.Exec(
		"UPDATE users SET email = ? WHERE id = ?",
		newEmail, id,
	)
	if err != nil {
		return 0, fmt.Errorf("updateUserEmail: %w", err)
	}

	affected, err := result.RowsAffected() // 변경된 행 수
	if err != nil {
		return 0, err
	}
	return affected, nil
}

// DELETE
func deleteUser(db *sql.DB, id int) error {
	result, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return fmt.Errorf("user id=%d not found", id)
	}
	return nil
}

type UserUpdate struct {
	ID   int
	Name string
	Age  int
}

// Prepared Statement 쿼리 미리 컴파일
func batchUpdate(db *sql.DB, updates []UserUpdate) error {
	// Prepare: 쿼리를 미리 컴파일
	stmt, err := db.Prepare("UPDATE users SET name = ?, age = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("prepare: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt) // 반드시 닫기

	for _, upd := range updates {
		// 컴파일된 쿼리에 파라미터만 바인딩해 실행
		if _, err := stmt.Exec(upd.Name, upd.Age, upd.ID); err != nil {
			return fmt.Errorf("exec id=%d: %w", upd.ID, err)
		}
	}
	return nil
}

// Stmt를 이용한 SELECT
func preparedSelect(db *sql.DB) error {
	stmt, err := db.Prepare("SELECT id, name FROM users WHERE status = ? LIMIT ?")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {

		}
	}(stmt)

	rows, err := stmt.Query("active", 100)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Printf("%d: %s\n", id, name)
	}
	return rows.Err()
}

// 계좌 이체 예제 — 원자성 보장 트렌젝션
func transfer(db *sql.DB, fromID, toID, amount int) error {
	// 트랜잭션 시작
	tx, err := db.Begin() // Begin - sql.Tx 객체 생성해줌.
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	// defer로 롤백 보장 — Commit 이후엔 no-op이 됨
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// 출금 계좌에서 차감
	res, err := tx.Exec(
		"UPDATE accounts SET balance = balance - ? WHERE id = ? AND balance >= ?",
		amount, fromID, amount,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("debit: %w", err)
	}
	if n, _ := res.RowsAffected(); n == 0 {
		_ = tx.Rollback()
		return fmt.Errorf("잔액 부족 또는 계좌 없음")
	}

	// 입금 계좌에 추가
	if _, err = tx.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE id = ?",
		amount, toID,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("credit: %w", err)
	}

	// 트랜잭션 로그 기록
	if _, err = tx.Exec(
		"INSERT INTO tx_log (from_id, to_id, amount) VALUES (?, ?, ?)",
		fromID, toID, amount,
	); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("log: %w", err)
	}

	// 모든 작업 성공 시 커밋
	return tx.Commit()
}
