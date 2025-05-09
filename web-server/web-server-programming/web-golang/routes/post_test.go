package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	jwt "github.com/conglt10/web-golang/auth"
	db "github.com/conglt10/web-golang/database"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGetAllPosts(t *testing.T) {
	// Thiết lập dữ liệu test
	collection := db.ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Tạo dữ liệu test
	testPosts := []interface{}{
		bson.M{
			"id":      "test-id-1",
			"creater": "testuser1",
			"title":   "Test post 1",
		},
		bson.M{
			"id":      "test-id-2",
			"creater": "testuser2",
			"title":   "Test post 2",
		},
	}

	_, err := collection.InsertMany(ctx, testPosts)
	if err != nil {
		t.Fatal(err)
	}

	// Thực hiện test
	req := httptest.NewRequest("GET", "/posts", nil)
	rr := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/posts", GetAllPosts)
	router.ServeHTTP(rr, req)

	// Kiểm tra status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler trả về status code không đúng: nhận được %v muốn %v",
			status, http.StatusOK)
	}

	// Kiểm tra response body
	var response []bson.M
	if err = json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Errorf("không thể decode response body: %v", err)
	}

	if len(response) != 8 {
		t.Errorf("số lượng bài đăng không đúng: nhận được %v muốn %v",
			len(response), 8)
	}

	// Dọn dẹp dữ liệu test
	_, err = collection.DeleteMany(ctx, bson.M{"id": bson.M{"$in": []string{"test-id-1", "test-id-2"}}})
	if err != nil {
		t.Errorf("không thể xóa dữ liệu test: %v", err)
	}
}

func TestCreatePost(t *testing.T) {
	testUser := "testuser"
	validToken, err := jwt.Create(testUser)
	if err != nil {
		t.Fatalf("không thể tạo token: %v", err)
	}

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		token          string
		expectedStatus int
	}{
		{
			name: "Tạo bài đăng thành công",
			requestBody: map[string]interface{}{
				"title": "Test post",
			},
			token:          validToken,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Không có token",
			requestBody: map[string]interface{}{
				"title": "Test post",
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Title rỗng",
			requestBody: map[string]interface{}{
				"title": "",
			},
			token:          validToken,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonBody))

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()

			router := httprouter.New()
			router.POST("/posts", CreatePost)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler trả về status code không đúng cho test case '%s': nhận được %v muốn %v",
					tt.name, status, tt.expectedStatus)
			}
		})
	}
}

func TestEditPost(t *testing.T) {
	// Thiết lập dữ liệu test
	testUser := "testuser"
	testPostID := "test-post-id"
	validToken, err := jwt.Create(testUser)
	if err != nil {
		t.Fatalf("không thể tạo token: %v", err)
	}

	// Tạo bài đăng test
	collection := db.ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testPost := bson.M{
		"id":      testPostID,
		"creater": testUser,
		"title":   "Original title",
	}

	_, err = collection.InsertOne(ctx, testPost)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		postID         string
		requestBody    EditPostRequest
		token          string
		expectedStatus int
	}{
		{
			name:   "Sửa bài đăng thành công",
			postID: testPostID,
			requestBody: EditPostRequest{
				Title: "Updated title",
			},
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Không có token",
			postID: testPostID,
			requestBody: EditPostRequest{
				Title: "Updated title",
			},
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "Bài đăng không tồn tại",
			postID: "non-existent-id",
			requestBody: EditPostRequest{
				Title: "Updated title",
			},
			token:          validToken,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/posts/"+tt.postID, bytes.NewBuffer(jsonBody))

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()

			router := httprouter.New()
			router.PUT("/posts/:id", EditPost)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler trả về status code không đúng: nhận được %v muốn %v",
					status, tt.expectedStatus)
			}
		})
	}

	// Dọn dẹp dữ liệu test
	_, err = collection.DeleteOne(ctx, bson.M{"id": testPostID})
	if err != nil {
		t.Errorf("không thể xóa dữ liệu test: %v", err)
	}
}

func TestDeletePost(t *testing.T) {
	// Thiết lập dữ liệu test
	testUser := "testuser"
	testPostID := "test-post-id"
	validToken, err := jwt.Create(testUser)
	if err != nil {
		t.Fatalf("không thể tạo token: %v", err)
	}

	// Tạo bài đăng test
	collection := db.ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testPost := bson.M{
		"id":      testPostID,
		"creater": testUser,
		"title":   "Test post",
	}

	_, err = collection.InsertOne(ctx, testPost)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name           string
		postID         string
		token          string
		expectedStatus int
	}{
		{
			name:           "Xóa bài đăng thành công",
			postID:         testPostID,
			token:          validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Không có token",
			postID:         testPostID,
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Bài đăng không tồn tại",
			postID:         "non-existent-id",
			token:          validToken,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", "/posts/"+tt.postID, nil)

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()

			router := httprouter.New()
			router.DELETE("/posts/:id", DeletePost)
			router.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler trả về status code không đúng: nhận được %v muốn %v",
					status, tt.expectedStatus)
			}
		})
	}
}
