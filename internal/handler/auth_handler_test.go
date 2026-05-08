package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dilvi/camp-booking-rest-api-go/internal/app"
	"github.com/dilvi/camp-booking-rest-api-go/internal/config"
)

type registerRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type registerResponse struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

type registerResponseEnvelope struct {
	Data registerResponse `json:"data"`
}

func TestRegisterHandler(t *testing.T) {
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run database integration tests")
	}

	t.Setenv("APP_PORT", "8080")
	t.Setenv("DB_HOST", "127.0.0.1")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "postgres")
	t.Setenv("DB_PASSWORD", "postgres")
	t.Setenv("DB_NAME", "camp_booking")
	t.Setenv("DB_SSLMODE", "disable")
	cfg := config.Load()

	application, err := app.New(cfg)
	if err != nil {
		t.Fatalf("failed to create app: %v", err)
	}
	defer application.DB.Close()

	email := "test_" + time.Now().Format("20060102150405") + "@example.com"
	phone := "+7999" + time.Now().Format("150405")

	reqBody := registerRequest{
		FirstName: "Artem",
		LastName:  "Shmelev",
		Phone:     phone,
		Email:     email,
		Password:  "12345678",
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	application.Router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var envelope registerResponseEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&envelope); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	resp := envelope.Data
	if resp.ID == 0 {
		t.Fatal("expected non-zero user id")
	}

	if resp.Email != email {
		t.Fatalf("expected email %s, got %s", email, resp.Email)
	}

	if resp.Role != "parent" {
		t.Fatalf("expected role parent, got %s", resp.Role)
	}

	_, err = application.DB.Exec("DELETE FROM users WHERE email = $1", email)
	if err != nil {
		t.Fatalf("failed to cleanup test user: %v", err)
	}
}
