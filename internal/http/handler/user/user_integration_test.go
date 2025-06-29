package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SoraDaibu/go-clean-starter/builder"
	"github.com/SoraDaibu/go-clean-starter/config"
	"github.com/SoraDaibu/go-clean-starter/internal/http/handler/user"
	"github.com/SoraDaibu/go-clean-starter/migration"
)

func TestMain(m *testing.M) {
	// Setup test database
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}

	// Run migrations
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name, cfg.DB.SSLMode)

	if err := migration.Up(dbURL); err != nil {
		panic(fmt.Sprintf("failed to run migrations: %v", err))
	}

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupTestDependencies(t *testing.T) (*builder.Dependency, func()) {
	cfg, err := config.Load()
	require.NoError(t, err)

	dependency, err := builder.InitializeDependency(cfg)
	require.NoError(t, err)

	cleanup := func() {
		if dependency.DB != nil {
			dependency.DB.Close()
		}
	}

	return dependency, cleanup
}

func TestUserHandler_CreateUser(t *testing.T) {
	dependency, cleanup := setupTestDependencies(t)
	defer cleanup()

	handler := builder.InitializeUserHandler(dependency)
	e := echo.New()

	// Add duplicate email test separately to handle shared email properly
	duplicateEmail := fmt.Sprintf("duplicate-%s@example.com", uuid.New().String())

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		preTestFunc    func(t *testing.T)
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name: "successful user creation",
			requestBody: map[string]string{
				"name":     "John Doe",
				"email":    fmt.Sprintf("john.doe-%s@example.com", uuid.New().String()),
				"password": "password123",
			},
			expectedStatus: http.StatusCreated,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var user map[string]interface{}
				err := json.Unmarshal(response.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, "John Doe", user["name"])

				// Check if user ID exists and is valid
				userID, exists := user["id"]
				require.True(t, exists, "User ID should exist in response")
				require.NotNil(t, userID, "User ID should not be nil")

				userIDStr, ok := userID.(string)
				require.True(t, ok, "User ID should be a string")
				require.NotEmpty(t, userIDStr, "User ID should not be empty")

				// Validate that it's a valid UUID format
				_, err = uuid.Parse(userIDStr)
				require.NoError(t, err, "User ID should be a valid UUID")

				assert.Equal(t, response.Code, http.StatusCreated)
			},
		},
		{
			name: "missing required fields",
			requestBody: map[string]string{
				"name": "John Doe",
				// missing email and password
			},
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var errorResp struct {
					Status  int    `json:"status"`
					Title   string `json:"title"`
					Details []struct {
						Field string `json:"field,omitempty"`
						Text  string `json:"text"`
					} `json:"details"`
				}
				err := json.Unmarshal(response.Body.Bytes(), &errorResp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, errorResp.Status)
				assert.Equal(t, "Bad Request", errorResp.Title)
				assert.NotEmpty(t, errorResp.Details)
				for _, detail := range errorResp.Details {
					assert.True(t,
						detail.Field == "email" || detail.Field == "password",
						"Expected field to be 'email' or 'password', got: %s", detail.Field)
					assert.True(t,
						strings.Contains(detail.Text, "is required"),
						"Expected validation error message, got: %s", detail.Text)
				}
			},
		},
		{
			name:           "invalid JSON",
			requestBody:    "{invalid json",
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var errorResp struct {
					Status  int    `json:"status"`
					Title   string `json:"title"`
					Details []struct {
						Text string `json:"text"`
					} `json:"details"`
				}
				err := json.Unmarshal(response.Body.Bytes(), &errorResp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, errorResp.Status)
				assert.Equal(t, "Bad Request", errorResp.Title)
				assert.NotEmpty(t, errorResp.Details)
				assert.Equal(t, "invalid parameter", errorResp.Details[0].Text)
			},
		},
		{
			name: "duplicate email",
			requestBody: map[string]string{
				"name":     "Jane Doe",
				"email":    duplicateEmail,
				"password": "password123",
			},
			expectedStatus: http.StatusConflict,
			preTestFunc: func(t *testing.T) {
				// Create a user first with the same email
				createTestUser(t, handler, e, map[string]string{
					"name":     "Jane Doe",
					"email":    duplicateEmail,
					"password": "password123",
				})
			},
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var errorResp struct {
					Status  int    `json:"status"`
					Title   string `json:"title"`
					Details []struct {
						Text string `json:"text"`
					} `json:"details"`
				}
				err := json.Unmarshal(response.Body.Bytes(), &errorResp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusConflict, errorResp.Status)
				assert.Equal(t, "Conflict", errorResp.Title)
				assert.NotEmpty(t, errorResp.Details)
				assert.Equal(t, "Resource already exists", errorResp.Details[0].Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			if tt.preTestFunc != nil {
				tt.preTestFunc(t)
			}

			rec := httptest.NewRecorder()
			ctx := e.NewContext(httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)), rec)
			ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			handler.CreateUser(ctx)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.validateFunc(t, rec)
		})
	}
}

func TestUserHandler_GetUser(t *testing.T) {
	dependency, cleanup := setupTestDependencies(t)
	defer cleanup()

	handler := builder.InitializeUserHandler(dependency)
	e := echo.New()

	// Create a user first
	createdUser := createTestUser(t, handler, e, map[string]string{
		"name":     "Test User",
		"email":    fmt.Sprintf("test-%s@example.com", uuid.New().String()),
		"password": "password123",
	})

	tests := []struct {
		name           string
		userID         string
		expectedStatus int
		validateFunc   func(t *testing.T, response *httptest.ResponseRecorder)
	}{
		{
			name:           "successful user retrieval",
			userID:         createdUser["id"].(string),
			expectedStatus: http.StatusOK,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var user map[string]interface{}
				err := json.Unmarshal(response.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, "Test User", user["name"])
				assert.Equal(t, createdUser["id"], user["id"])
			},
		},
		{
			name:           "user not found",
			userID:         uuid.New().String(),
			expectedStatus: http.StatusNotFound,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var errorResp struct {
					Status  int    `json:"status"`
					Title   string `json:"title"`
					Details []struct {
						Text string `json:"text"`
					} `json:"details"`
				}
				err := json.Unmarshal(response.Body.Bytes(), &errorResp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusNotFound, errorResp.Status)
				assert.Equal(t, "Not Found", errorResp.Title)
				assert.NotEmpty(t, errorResp.Details)
				assert.Equal(t, "User not found", errorResp.Details[0].Text)
			},
		},
		{
			name:           "invalid UUID",
			userID:         "invalid-uuid",
			expectedStatus: http.StatusBadRequest,
			validateFunc: func(t *testing.T, response *httptest.ResponseRecorder) {
				var errorResp struct {
					Status  int    `json:"status"`
					Title   string `json:"title"`
					Details []struct {
						Text string `json:"text"`
					} `json:"details"`
				}
				err := json.Unmarshal(response.Body.Bytes(), &errorResp)
				assert.NoError(t, err)

				assert.Equal(t, http.StatusBadRequest, errorResp.Status)
				assert.Equal(t, "Bad Request", errorResp.Title)
				assert.NotEmpty(t, errorResp.Details)
				assert.Equal(t, "invalid UUID format", errorResp.Details[0].Text)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			ctx := e.NewContext(httptest.NewRequest(http.MethodGet, "/users/"+tt.userID, nil), rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.userID)

			handler.GetUser(ctx)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			tt.validateFunc(t, rec)
		})
	}
}

func createTestUser(t *testing.T, handler *user.UserHandler, e *echo.Echo, requestBody map[string]string) map[string]interface{} {
	body, _ := json.Marshal(requestBody)

	rec := httptest.NewRecorder()
	ctx := e.NewContext(httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body)), rec)
	ctx.Request().Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	handler.CreateUser(ctx)

	var user map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &user)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)

	return user
}

// TestUserHandler_Integration tests the full flow of creating and retrieving users
func TestUserHandler_Integration(t *testing.T) {
	dependency, cleanup := setupTestDependencies(t)
	defer cleanup()

	handler := builder.InitializeUserHandler(dependency)
	e := echo.New()

	// Test full integration flow
	t.Run("create and retrieve user flow", func(t *testing.T) {
		// Step 1: Create a user
		createReq := map[string]string{
			"name":     "Integration Test User",
			"email":    fmt.Sprintf("integration-%s@example.com", uuid.New().String()),
			"password": "integrationpassword123",
		}

		// Create a user
		createdUser := createTestUser(t, handler, e, createReq)

		// Retrieve the created user
		userID := createdUser["id"].(string)
		getHttpReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", userID), nil)
		getRec := httptest.NewRecorder()
		getCtx := e.NewContext(getHttpReq, getRec)
		getCtx.SetParamNames("id")
		getCtx.SetParamValues(userID)

		err := handler.GetUser(getCtx)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, getRec.Code)

		var retrievedUser map[string]interface{}
		err = json.Unmarshal(getRec.Body.Bytes(), &retrievedUser)
		require.NoError(t, err)

		// Verify the retrieved user matches the created user
		assert.Equal(t, createdUser["id"], retrievedUser["id"])
		assert.Equal(t, "Integration Test User", retrievedUser["name"])
	})
}
