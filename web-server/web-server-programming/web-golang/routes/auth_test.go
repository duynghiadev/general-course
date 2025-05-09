package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	db "github.com/conglt10/web-golang/database"
	"github.com/conglt10/web-golang/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func TestLogin(t *testing.T) {
	// Thiết lập dữ liệu test
	testUser := "testuser"
	testPass := "testpass123"
	hashedPassword, _ := models.Hash(testPass)

	// Kết nối database và tạo user test
	collection := db.ConnectUsers()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Xóa user test nếu tồn tại
	collection.DeleteOne(ctx, bson.M{"username": testUser})

	// Tạo user test mới
	_, err := collection.InsertOne(ctx, bson.M{
		"username": testUser,
		"password": hashedPassword,
		"email":    "test@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		requestBody    LoginRequest
		expectedStatus int
	}{
		{
			name: "Successful login",
			requestBody: LoginRequest{
				Username: testUser,
				Password: testPass,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid username",
			requestBody: LoginRequest{
				Username: "nonexistent",
				Password: testPass,
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid password",
			requestBody: LoginRequest{
				Username: testUser,
				Password: "wrongpass",
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Empty username",
			requestBody: LoginRequest{
				Username: "",
				Password: testPass,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty password",
			requestBody: LoginRequest{
				Username: testUser,
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(jsonBody))
			rr := httptest.NewRecorder()

			router := httprouter.New()
			router.POST("/login", Login)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}
		})
	}

	// Cleanup
	collection.DeleteOne(ctx, bson.M{"username": testUser})
}

func TestRegister(t *testing.T) {
	collection := db.ConnectUsers()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tests := []struct {
		name           string
		requestBody    RegisterRequest
		expectedStatus int
	}{
		{
			name: "Successful registration",
			requestBody: RegisterRequest{
				Username: "newuser",
				Password: "newpass123",
				Email:    "new@example.com",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Empty username",
			requestBody: RegisterRequest{
				Username: "",
				Password: "pass123",
				Email:    "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty password",
			requestBody: RegisterRequest{
				Username: "testuser",
				Password: "",
				Email:    "test@example.com",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Empty email",
			requestBody: RegisterRequest{
				Username: "testuser",
				Password: "pass123",
				Email:    "",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			requestBody: RegisterRequest{
				Username: "testuser",
				Password: "pass123",
				Email:    "invalid-email",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Xóa user test nếu tồn tại
			collection.DeleteOne(ctx, bson.M{"username": tt.requestBody.Username})

			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(jsonBody))
			rr := httptest.NewRecorder()

			router := httprouter.New()
			router.POST("/register", Register)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			// Cleanup sau mỗi test
			collection.DeleteOne(ctx, bson.M{"username": tt.requestBody.Username})
		})
	}
}
