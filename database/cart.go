package database

import (
	"context"
	"errors"
	"log"

	"github.com/akhil/ecommerce-yt/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddProductToCart(ctx context.Context, prodCollection, userCollection *mongo.Collection, productId primitive.ObjectID, userID string) error {
	searchProduct, err := prodCollection.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		return errors.New("Product not found")
	}
	var productCart []models.ProductUser
	err = searchProduct.All(ctx, &productCart)
	if err != nil {
		return errors.New("UserId is not valid")
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("UserId is not valid")
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "usercart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.New("Add Product to Cart failed")
	}
	return nil
}
func RemoveCartItem(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, userID string) error {
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return errors.New("UserId is not valid")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.M{"$pull": bson.M{"usercart": bson.M{"_id": productID}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return errors.New("Remove Item failed")
	}
	return nil

}
