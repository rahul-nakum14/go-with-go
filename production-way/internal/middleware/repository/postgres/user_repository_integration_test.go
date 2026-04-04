//go:build integration

// Run with: make test/integration  (requires TEST_DATABASE_URL env var)
package postgres_test

import (
	"context"
	"os"
	"testing"

	"github.com/yourorg/go-user-crud/internal/domain"
	"github.com/yourorg/go-user-crud/internal/repository/postgres"
)

func setupDB(t *testing.T) *postgres.UserRepository {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}
	pool, err := postgres.NewPool(dsn)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	t.Cleanup(pool.Close)
	return postgres.NewUserRepository(pool)
}

func TestUserRepository_CreateAndGet(t *testing.T) {
	repo := setupDB(t)
	ctx := context.Background()

	input := domain.CreateUserInput{
		Name:  "Integration Test User",
		Email: "integration@example.com",
		Role:  domain.RoleUser,
	}

	created, err := repo.Create(ctx, input)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == 0 {
		t.Fatal("expected non-zero ID")
	}

	fetched, err := repo.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if fetched.Email != input.Email {
		t.Errorf("email = %q, want %q", fetched.Email, input.Email)
	}

	// Cleanup
	if err := repo.Delete(ctx, created.ID); err != nil {
		t.Errorf("Delete: %v", err)
	}
}
