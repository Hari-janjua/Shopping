package services

import (
	"Shopping/Models"
	"Shopping/db"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func GetProductsDetail(productIds []string) ([]Models.Product, error) {
	log.Println("IN Service: GetProductsDetail")
	defer log.Println("IN Service: GetProductsDetail")

	var products []Models.Product
	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return products, dberr
	}
	query := bson.M{"productId": bson.M{"$in": productIds}}
	collection := session.Database("shopping").Collection("Product")
	cursor, ferr := collection.Find(context.Background(), query)
	if nil != ferr {
		if ferr.Error() == "mongo: no documents in result" {
			log.Println("ERR_IFR_APPLICATIONS_NOT_FOUND: ", ferr)
			return products, nil
		}
		log.Println("ERR_FEATCHING_IFR_APPLICATIONS: ", ferr)
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

func CheckProductCountLessThanTen(products []Models.ProductModel) error {
	log.Println("IN Service: CheckProductCountLessThanTen")
	defer log.Println("OUT Service: CheckProductCountLessThanTen")

	for _, product := range products {
		if product.Count > 10 {
			log.Println("Product count is more than 10.")
			return errors.New("PRODUCT_COUNT_IS_MORE_THAN_10")
		}
	}

	return nil
}

func GetOrderedProductDetail(productInStock Models.Product, orderedDetails Models.OrderDetails) (Models.ProductModel, error) {
	log.Println("IN Service: GetOrderedProductDetail")
	defer log.Println("OUT Service: GetOrderedProductDetail")

	for i := 0; i < len(orderedDetails.Products); i++ {

		// log.Println("product: ", orderProduct)
		if orderedDetails.Products[i].ProductId == productInStock.ProductId {
			return orderedDetails.Products[i], nil
		}
	}

	return Models.ProductModel{}, errors.New("OUT_OF_STOCK")

}

func UpdateProductCatalog(products []Models.ProductModel, productStock []Models.Product) error {
	log.Println("IN Service: UpdateProductCatalog")
	defer log.Println("OUT Service: UpdateProductCatalog")

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}

	coll := session.Database("shopping").Collection("Product")

	// Bulk update the product catalog
	models := []mongo.WriteModel{}
	for _, product := range products {
		fmt.Println("product: ", product)
		productCountInStock := getCountOfProductInStock(product.ProductId, productStock)
		newAvailableCount := productCountInStock - product.Count
		models = append(models, mongo.NewUpdateOneModel().SetFilter(bson.M{"productId": product.ProductId}).SetUpdate(bson.M{"$set": bson.M{"availability.count": newAvailableCount}}))
	}

	opts := options.BulkWrite().SetOrdered(true)
	_, err := coll.BulkWrite(context.TODO(), models, opts)
	if err != nil {
		log.Println("ERROR WHILE BULK WRITE: ", err)
		return err
	}
	// fmt.Println("results: ", results)

	return nil
}

func getCountOfProductInStock(productId string, productStock []Models.Product) int {
	log.Println("IN service: getCountOfProductInStock")
	log.Println("OUT service: getCountOfProductInStock")

	for _, product := range productStock {
		if productId == product.ProductId {
			return product.Availability.Count
		}
	}

	return 0
}

func InsertOrderDetails(orderedData Models.Order) error {
	log.Println("IN service: InsertOrderDetails")
	log.Println("OUT service: InsertOrderDetails")

	session, dberr := db.GetMongoConnection()
	ctx := context.Background()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}
	cd := session.Database("shopping").Collection("Order")
	_, err := cd.InsertOne(ctx, orderedData)
	if err != nil {
		log.Println("err: ", err)
		return err
	}

	return nil
}

func GetOrderDetailsById(orderId string) (Models.Order, error) {
	log.Println("IN service: GetOrderDetailsById")
	log.Println("OUT service: GetOrderDetailsById")

	var order Models.Order

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return order, dberr
	}

	query := bson.M{"orderId": orderId}
	collection := session.Database("shopping").Collection("Order")
	ferr := collection.FindOne(context.Background(), query).Decode(&order)
	if nil != ferr {
		if ferr.Error() == "mongo: no documents in result" {
			log.Println("ERR_IFR_APPLICATIONS_NOT_FOUND: ", ferr)
			return order, nil
		}
		log.Println("ERR_FEATCHING_IFR_APPLICATIONS: ", ferr)
		return order, ferr
	}

	return order, nil
}

func GetOrderDetailsByUserId(userId string) ([]Models.Order, error) {
	log.Println("IN service: GetOrderDetailsByUserId")
	log.Println("OUT service: GetOrderDetailsByUserId")

	var orders []Models.Order

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return orders, dberr
	}
	query := bson.M{}
	collection := session.Database("shopping").Collection("Order")

	cursor, ferr := collection.Find(context.Background(), query)
	if nil != ferr {
		if ferr.Error() == "mongo: no documents in result" {
			log.Println("ERR_ORDER_NOT_FOUND: ", ferr)
			return orders, nil
		}
		log.Println("ERR_FEATCHING_ORDERS: ", ferr)
		return orders, ferr
	}

	defer cursor.Close(context.Background())

	err := cursor.All(context.Background(), &orders)
	if err != nil {
		log.Println("ERR_FEATCHING_ORDER", err)
		return orders, err
	}

	return orders, nil
}

func DispatchOrderByOrderId(orderId string) error {
	log.Println("IN service: DispatchOrderByOrderId")
	log.Println("OUT service: DispatchOrderByOrderId")

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}
	collection := session.Database("shopping").Collection("Order")
	selector := bson.M{"orderId": orderId}
	updator := bson.M{"$set": bson.M{"orderStatus": "Dispatched", "dispatchDate": time.Now(), "modifiedOn": time.Now()}}

	_, updateErr := collection.UpdateOne(context.Background(), selector, updator)
	if updateErr != nil {
		log.Println("Error while decoding order details : DispatchOrderByOrderId", updateErr)
		return updateErr
	}

	return nil
}

func CompleteOrderByOrderId(orderId string) error {
	log.Println("IN service: DispatchOrderByOrderId")
	log.Println("OUT service: DispatchOrderByOrderId")

	session, dberr := db.GetMongoConnection()
	if dberr != nil {
		log.Println("ERROR IN GETTING DB CONNECTION: ", dberr.Error())
		return dberr
	}
	collection := session.Database("shopping").Collection("Order")
	selector := bson.M{"orderId": orderId}
	updator := bson.M{"$set": bson.M{"orderStatus": "Completed", "modifiedOn": time.Now()}}

	_, updateErr := collection.UpdateOne(context.Background(), selector, updator)
	if updateErr != nil {
		log.Println("Error while decoding order details : CompleteOrderByOrderId", updateErr)
		return updateErr
	}

	return nil
}
