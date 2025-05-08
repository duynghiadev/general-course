package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/conglt10/web-golang/auth"
	db "github.com/conglt10/web-golang/database"
	"github.com/conglt10/web-golang/errors"
	"github.com/conglt10/web-golang/models"
	res "github.com/conglt10/web-golang/utils"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var req LoginRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		res.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body", "details": err.Error()})
		return
	}
	defer r.Body.Close()

	// Validate input
	if strings.TrimSpace(req.Username) == "" {
		res.JSON(w, http.StatusBadRequest, errors.ErrEmptyUsername)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		res.JSON(w, http.StatusBadRequest, errors.ErrEmptyPassword)
		return
	}

	// Sanitize input
	username := models.Santize(req.Username)
	password := models.Santize(req.Password)

	// Connect to database
	collection := db.ConnectUsers()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		res.JSON(w, http.StatusUnauthorized, "Username or Password incorrect")
		return
	} else if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Verify password
	hashedPassword, ok := result["password"].(string)
	if !ok {
		res.JSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err = models.CheckPasswordHash(hashedPassword, password); err != nil {
		res.JSON(w, http.StatusUnauthorized, "Username or Password incorrect")
		return
	}

	// Generate JWT token
	token, err := jwt.Create(username)
	if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	res.JSON(w, http.StatusOK, map[string]string{"token": token})
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		res.JSON(w, http.StatusBadRequest, map[string]string{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	defer r.Body.Close()

	// Validate input
	if strings.TrimSpace(req.Username) == "" {
		res.JSON(w, http.StatusBadRequest, errors.ErrEmptyUsername)
		return
	}
	if strings.TrimSpace(req.Password) == "" {
		res.JSON(w, http.StatusBadRequest, errors.ErrEmptyPassword)
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		res.JSON(w, http.StatusBadRequest, errors.ErrEmptyEmail)
		return
	}
	if !govalidator.IsEmail(req.Email) {
		res.JSON(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Sanitize input
	username := models.Santize(req.Username)
	email := models.Santize(req.Email)
	password := models.Santize(req.Password)

	// Connect to database
	collection := db.ConnectUsers()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if user or email exists
	var result bson.M
	err := collection.FindOne(ctx, bson.M{"$or": []bson.M{
		{"username": username},
		{"email": email},
	}}).Decode(&result)
	if err == nil {
		res.JSON(w, http.StatusConflict, "Username or email already exists")
		return
	} else if err != mongo.ErrNoDocuments {
		res.JSON(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	// Hash password
	hashedPassword, err := models.Hash(password)
	if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create new user
	newUser := bson.M{
		"username": username,
		"email":    email,
		"password": hashedPassword,
	}

	_, err = collection.InsertOne(ctx, newUser)
	if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	res.JSON(w, http.StatusCreated, "Registration successful")
}
