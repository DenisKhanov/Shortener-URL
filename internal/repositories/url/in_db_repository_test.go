package url

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"testing"
)

func TestNewURLInDBRepo(t *testing.T) {
	// Создаем поддельный пул соединений
	db, err := pgxpool.New(context.Background(), "user=test password=1111111 dbname=shortenerURL sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer db.Close()

	// Вызываем функцию, которую хотим протестировать
	repo, _ := NewURLInDBRepo(db)

	// Проверяем, что UserID был инициализирован как uuid.Nil
	if repo.UserID != uuid.Nil {
		t.Errorf("Expected UserID to be uuid.Nil, got %v", repo.UserID)
	}

	// Проверяем, что DB был инициализирован
	if repo.DB == nil {
		t.Error("Expected DB to be initialized, got nil")
	}
}
