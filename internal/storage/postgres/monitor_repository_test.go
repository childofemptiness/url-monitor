package postgres

import (
	"context"
	"errors"
	"testing"
	"time"
	"url-monitor/internal/monitor"

	"github.com/jackc/pgx/v5/pgxpool"
)

func TestRepositoryCreate(t *testing.T) {
	pool := setupTestDatabase(t)
	repo := NewMonitorRepository(pool)

	nextCheckAt := time.Now().UTC().Add(45 * time.Second).Truncate(time.Microsecond)

	created, err := repo.Create(context.Background(), monitor.Monitor{
		URL:             "https://example.com",
		IntervalSeconds: 45,
		NextCheckAt:     &nextCheckAt,
	})
	if err != nil {
		t.Fatalf("create monitor: %v", err)
	}

	if created.ID == 0 {
		t.Fatal("expected created monitor id to be set")
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("expected created_at to be set")
	}
	if created.NextCheckAt == nil {
		t.Fatal("expected next_check_at to be set")
	}
	if !created.NextCheckAt.Equal(nextCheckAt) {
		t.Fatalf("expected next_check_at %s, got %s", nextCheckAt, created.NextCheckAt)
	}
}

func TestRepositoryList(t *testing.T) {
	pool := setupTestDatabase(t)
	repo := NewMonitorRepository(pool)

	firstNextCheckAt := time.Now().UTC().Add(10 * time.Second).Truncate(time.Microsecond)
	secondNextCheckAt := time.Now().UTC().Add(20 * time.Second).Truncate(time.Microsecond)

	insertMonitorRow(t, pool, "https://first.example.com", 10, firstNextCheckAt)
	insertMonitorRow(t, pool, "https://second.example.com", 20, secondNextCheckAt)

	monitors, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list monitors: %v", err)
	}

	if len(monitors) != 2 {
		t.Fatalf("expected 2 monitors, got %d", len(monitors))
	}
	if monitors[0].URL != "https://first.example.com" {
		t.Fatalf("expected first monitor url to be https://first.example.com, got %s", monitors[0].URL)
	}
	if monitors[0].CreatedAt.IsZero() {
		t.Fatal("expected created_at to be populated")
	}
	if monitors[1].NextCheckAt == nil || !monitors[1].NextCheckAt.Equal(secondNextCheckAt) {
		t.Fatalf("expected second next_check_at %s, got %v", secondNextCheckAt, monitors[1].NextCheckAt)
	}
}

func TestRepositoryListDue(t *testing.T) {
	pool := setupTestDatabase(t)
	repo := NewMonitorRepository(pool)

	now := time.Now().UTC().Truncate(time.Microsecond)

	firstDueID := insertMonitorRow(t, pool, "https://due-one.example.com", 10, now.Add(-2*time.Minute))
	secondDueID := insertMonitorRow(t, pool, "https://due-two.example.com", 20, now)
	insertMonitorRow(t, pool, "https://future.example.com", 30, now.Add(2*time.Minute))
	insertMonitorWithoutNextCheckAt(t, pool, "https://missing-next-check.example.com", 40)

	monitors, err := repo.ListDue(context.Background(), now, 10)
	if err != nil {
		t.Fatalf("list due monitors: %v", err)
	}

	if len(monitors) != 2 {
		t.Fatalf("expected 2 due monitors, got %d", len(monitors))
	}
	if monitors[0].ID != firstDueID {
		t.Fatalf("expected first due monitor id %d, got %d", firstDueID, monitors[0].ID)
	}
	if monitors[1].ID != secondDueID {
		t.Fatalf("expected second due monitor id %d, got %d", secondDueID, monitors[1].ID)
	}
}

func TestRepositoryListDueRespectsLimit(t *testing.T) {
	pool := setupTestDatabase(t)
	repo := NewMonitorRepository(pool)

	now := time.Now().UTC().Truncate(time.Microsecond)
	firstDueID := insertMonitorRow(t, pool, "https://due-limit-one.example.com", 10, now.Add(-time.Minute))
	insertMonitorRow(t, pool, "https://due-limit-two.example.com", 20, now.Add(-30*time.Second))

	monitors, err := repo.ListDue(context.Background(), now, 1)
	if err != nil {
		t.Fatalf("list due monitors with limit: %v", err)
	}

	if len(monitors) != 1 {
		t.Fatalf("expected 1 due monitor, got %d", len(monitors))
	}
	if monitors[0].ID != firstDueID {
		t.Fatalf("expected monitor id %d, got %d", firstDueID, monitors[0].ID)
	}
}

func TestRepositoryCreateReturnsDuplicateError(t *testing.T) {
	pool := setupTestDatabase(t)
	repo := NewMonitorRepository(pool)

	createUniqueConstraint(t, pool)

	nextCheckAt := time.Now().UTC().Add(time.Minute).Truncate(time.Microsecond)
	input := monitor.Monitor{
		URL:             "https://duplicate.example.com",
		IntervalSeconds: 60,
		NextCheckAt:     &nextCheckAt,
	}

	if _, err := repo.Create(context.Background(), input); err != nil {
		t.Fatalf("create first monitor: %v", err)
	}

	_, err := repo.Create(context.Background(), input)
	if !errors.Is(err, monitor.ErrMonitorAlreadyExists) {
		t.Fatalf("expected ErrMonitorAlreadyExists, got %v", err)
	}
}

func insertMonitorRow(t *testing.T, pool *pgxpool.Pool, rawURL string, intervalSeconds int, nextCheckAt time.Time) int64 {
	t.Helper()

	var id int64
	err := pool.QueryRow(context.Background(), `
		INSERT INTO monitors (url, interval_seconds, next_check_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`, rawURL, intervalSeconds, nextCheckAt).Scan(&id)
	if err != nil {
		t.Fatalf("insert monitor row: %v", err)
	}

	return id
}

func insertMonitorWithoutNextCheckAt(t *testing.T, pool *pgxpool.Pool, rawURL string, intervalSeconds int) int64 {
	t.Helper()

	var id int64
	err := pool.QueryRow(context.Background(), `
		INSERT INTO monitors (url, interval_seconds)
		VALUES ($1, $2)
		RETURNING id
	`, rawURL, intervalSeconds).Scan(&id)
	if err != nil {
		t.Fatalf("insert monitor without next_check_at: %v", err)
	}

	return id
}

func createUniqueConstraint(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	if _, err := pool.Exec(context.Background(), `
		ALTER TABLE monitors
		ADD CONSTRAINT monitors_url_key UNIQUE (url)
	`); err != nil {
		t.Fatalf("create unique constraint: %v", err)
	}
}
