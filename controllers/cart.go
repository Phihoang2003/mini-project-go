package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/akhil/ecommerce-yt/database"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}
func (app *Application) AddToCart(c *gin.Context) {
	productQueryID := c.Query("id")
	if productQueryID == "" {
		log.Println("product id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}
	userQueryID := c.Query("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}
	productID, err := primitive.ObjectIDFromHex(productQueryID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(200, "Successfully Added to the cart")
}
