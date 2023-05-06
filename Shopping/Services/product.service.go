package services

import (
	"Shopping/Models"
	"Shopping/db"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func CreateNewProduct(product Models.Product) error {
	log.Println("IN Service: CreateNewProduct")
	defer log.Println("OUT Service: CreateNewProduct")

	product.CreatedOn = time.Now()
	product.ModifiedOn = time.Now()
	fmt.Println("product: ", product)

	session, dberr := db.GetMongoConnection()
	ctx := context.Background()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}
	cd := session.Database("shopping").Collection("Product")
	_, err := cd.InsertOne(ctx, product)
	if err != nil {
		log.Println("err: ", err)
		return err
	}

	return nil
}

func GetAllProducts(pageNumber int) ([]Models.Product, error) {
	log.Println("IN Service: GetAllProducts")
	defer log.Println("OUT Service: GetAllProducts")

	var products []Models.Product

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return products, dberr
	}
	query := bson.M{}
	collection := session.Database("shopping").Collection("Product")

	// Pagination of data
	options := options.Find()
	options.SetSkip(int64(pageNumber) * int64(10))
	options.SetLimit(int64(10))

	cursor, ferr := collection.Find(context.Background(), query, options)
	if nil != ferr {
		if ferr.Error() == "mongo: no documents in result" {
			log.Println("ERR_PRODUCT_NOT_FOUND: ", ferr)
			return products, nil
		}
		log.Println("ERR_FEATCHING_PRODUCTS: ", ferr)
		return products, ferr
	}

	defer cursor.Close(context.Background())

	err := cursor.All(context.Background(), &products)
	if err != nil {
		log.Println("ERR_FEATCHING_PRODUCTS", err)
		return products, err
	}

	return products, nil
}

func GetProductById(productId string) (Models.Product, error) {
	log.Println("IN Service: GetProductById() ")
	defer log.Println("OUT Service: GetProductById ")

	var product Models.Product

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return product, dberr
	}

	query := bson.M{"productId": productId}
	collection := session.Database("shopping").Collection("Product")
	ferr := collection.FindOne(context.Background(), query).Decode(&product)
	if nil != ferr {
		if ferr.Error() == "mongo: no documents in result" {
			log.Println("ERR_PRODUCT_NOT_FOUND: ", ferr)
			return product, nil
		}
		log.Println("ERR_FEATCHING_PRODUCTS: ", ferr)
		return product, ferr
	}

	return product, nil
}

func DeleteProductById(productId string) error {
	log.Println("IN Service: GetProductById() ")
	defer log.Println("OUT Service: GetProductById ")

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}
	query := bson.M{"productId": productId}
	collection := session.Database("shopping").Collection("Product")
	_, dbInsrtErr := collection.DeleteOne(context.Background(), query)
	if dbInsrtErr != nil {
		log.Println("Error Removing Product Document", dbInsrtErr)
		return dbInsrtErr
	}

	return nil
}
