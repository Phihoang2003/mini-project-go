package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/akhil/ecommerce-yt/database"
	"github.com/akhil/ecommerce-yt/models"
	generate "github.com/akhil/ecommerce-yt/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userpassword string, givenpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpassword), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Passowrd is Incorerct"
		valid = false
	}
	return valid, msg
}
func SignUp(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "Email has been already"})

		return
	}

	count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"messsage": "Phone has been already"})
		return
	}
	password := HashPassword(*user.Password)
	user.Password = &password
	user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.User_ID = user.ID.Hex()
	user.ID = primitive.NewObjectID()
	user.User_ID = user.ID.Hex()
	token, refreshToken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
	user.Token = &token
	user.Refresh_Token = &refreshToken
	user.UserCart = make([]models.ProductUser, 0)
	user.Address_Details = make([]models.Address, 0)
	user.Order_Status = make([]models.Order, 0)
	if _, err := UserCollection.InsertOne(ctx, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Sign up failed"})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"message": "Sign up successful"})

}
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var foundUser models.User
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email or password not correct"})
		return
	}
	passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
	defer cancel()
	if !passwordIsValid {
		c.JSON(http.StatusBadRequest, gin.H{"message": msg})
		return
	}
	token, refreshToken, _ := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
	defer cancel()
	generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
	c.JSON(http.StatusFound, foundUser)
}
func AddProductViewAdmin(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var products models.Product
	defer cancel()
	if err := c.BindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	products.Product_ID = primitive.NewObjectID()
	_, anyerr := ProductCollection.InsertOne(ctx, products)
	if anyerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, "Successfully added our Product Admin!!")
}
