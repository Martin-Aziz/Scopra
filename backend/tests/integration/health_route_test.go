package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/martin-aziz/scopra/backend/src/api"
	"github.com/martin-aziz/scopra/backend/src/services"
	"github.com/martin-aziz/scopra/backend/src/utils"
)

func TestLivenessEndpointReturnsOK(t *testing.T) {
	app := fiber.New()
	handler := api.NewHandler(nil, nil, nil, nil, nil, utils.NewLogger("test"))
	tokenService := services.NewTokenService("this-is-a-test-secret-with-at-least-32", "nexus", "clients", 15*time.Minute, 7*24*time.Hour)
	api.RegisterRoutes(app, handler, tokenService)

	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response, err := app.Test(request)
	if err != nil {
		t.Fatalf("expected liveness request to succeed: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.StatusCode)
	}
}
