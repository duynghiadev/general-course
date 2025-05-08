package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	jwt "github.com/conglt10/web-golang/auth"
	db "github.com/conglt10/web-golang/database"
	"github.com/conglt10/web-golang/models"
	res "github.com/conglt10/web-golang/utils"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type CreatePostRequest struct {
	Title string `json:"title"`
}

type EditPostRequest struct {
	Title string `json:"title"`
}

func GetAllPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := db.ConnectPosts()

	var result []bson.M
	data, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	defer data.Close(context.Background())
	for data.Next(context.Background()) {
		var elem bson.M
		err := data.Decode(&elem)

		if err != nil {
			res.JSON(w, 500, "Internal Server Error")
			return
		}

		result = append(result, elem)
	}

	res.JSON(w, 200, result)
}

func GetMyPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username, err := jwt.ExtractUsernameFromToken(r)

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	collection := db.ConnectPosts()

	var result []bson.M
	data, err := collection.Find(context.Background(), bson.M{"creater": username})

	defer data.Close(context.Background())
	for data.Next(context.Background()) {
		var elem bson.M
		err := data.Decode(&elem)

		if err != nil {
			res.JSON(w, 500, "Internal Server Error")
			return
		}

		result = append(result, elem)
	}

	res.JSON(w, 200, result)
}

func CreatePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	creater, err := jwt.ExtractUsernameFromToken(r)
	if err != nil {
		res.JSON(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	var req CreatePostRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&req); err != nil {
		res.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body", "details": err.Error()})
		return
	}
	defer r.Body.Close()

	if govalidator.IsNull(req.Title) {
		res.JSON(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	title := models.Santize(req.Title)
	uid := uuid.NewV4()
	id := fmt.Sprintf("%x-%x-%x-%x-%x", uid[0:4], uid[4:6], uid[6:8], uid[8:10], uid[10:])

	collection := db.ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newPost := bson.M{"id": id, "creater": creater, "title": title}
	_, err = collection.InsertOne(ctx, newPost)
	if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Failed to create post")
		return
	}

	res.JSON(w, http.StatusCreated, "Post created successfully")
}

func EditPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	if err != nil {
		res.JSON(w, http.StatusUnauthorized, "Invalid or missing token")
		return
	}

	var req EditPostRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&req); err != nil {
		res.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body", "details": err.Error()})
		return
	}
	defer r.Body.Close()

	if govalidator.IsNull(req.Title) {
		res.JSON(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	title := models.Santize(req.Title)
	collection := db.ConnectPosts()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result bson.M
	err = collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		res.JSON(w, http.StatusNotFound, "Post not found")
		return
	}

	creater, ok := result["creater"].(string)
	if !ok || username != creater {
		res.JSON(w, http.StatusForbidden, "Permission denied")
		return
	}

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"title": title}}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		res.JSON(w, http.StatusInternalServerError, "Failed to edit post")
		return
	}

	res.JSON(w, http.StatusOK, "Post updated successfully")
}

func DeletePost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	username, err := jwt.ExtractUsernameFromToken(r)
	collection := db.ConnectPosts()

	if err != nil {
		res.JSON(w, 500, "Internal Server Error")
		return
	}

	var result bson.M
	errFind := collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&result)

	if errFind != nil {
		res.JSON(w, 404, "Post Not Found")
		return
	}

	creater := fmt.Sprintf("%v", result["creater"])

	if username != creater {
		res.JSON(w, 403, "Permission Denied")
		return
	}

	errDelete := collection.FindOneAndDelete(context.TODO(), bson.M{"id": id}).Decode(&result)

	if errDelete != nil {
		res.JSON(w, 500, "Delete has failed")
		return
	}

	res.JSON(w, 200, "Delete Successfully")
}
