package Models

import "time"

type Product struct {
	ProductId    string       `json:"productId" bson:"productId"`
	ProductName  string       `json:"productName" bson:"productName"`
	Price        int          `json:"price" bson:"price"`
	Category     Category     `json:"category" bson:"category"`
	Availability Availability `json:"availability" bson:"availability"`
	CreatedBy    string       `json:"createdBy" bson:"createdBy"`
	CreatedOn    time.Time    `json:"createdOn,omitempty" bson:"createdOn"`
	ModifiedBy   string       `json:"modifiedBy" bson:"modifiedBy"`
	ModifiedOn   time.Time    `json:"modifiedOn,omitempty" bson:"modifiedOn"`
}

type Category struct {
	CategoryId   int    `json:"categoryId" bson:"categoryId"`
	CategoryName string `json:"categoryName" bson:"categoryName"`
}

type Availability struct {
	IsAvailable bool `json:"isAvailable" bson:"isAvailable"`
	Count       int  `json:"count" bson:"count"`
}

type Order struct {
	OrderId            string         `json:"orderId" bson:"orderId"`
	CustomerId         string         `json:"customerId" bson:"customerId"`
	Products           []ProductModel `json:"products" bson:"products"`
	DispatchDate       time.Time      `json:"dispatchDate" bson:"dispatchDate"`
	OrderStatus        string         `json:"orderStatus" bson:"orderStatus"`
	TotalCost          int            `json:"totalCost," bson:"totalCost,"`
	GrandTotal         int            `json:"grandTotal," bson:"grandTotal,"`
	DiscountPercentage int            `json:"discountPercentage," bson:"discountPercentage,"`
	Address            string         `json:"address" bson:"address"`
	CreatedBy          string         `json:"createdBy" bson:"createdBy"`
	CreatedOn          time.Time      `json:"createdOn" bson:"createdOn"`
	ModifiedBy         string         `json:"modifiedBy" bson:"modifiedBy"`
	ModifiedOn         time.Time      `json:"modifiedOn" bson:"modifiedOn"`
}

type ProductModel struct {
	ProductId   string   `json:"productId" bson:"productId"`
	ProductName string   `json:"productName" bson:"productName"`
	Price       int      `json:"price" bson:"price"`
	Category    Category `json:"category" bson:"category"`
	Count       int      `json:"count" bson:"count"`
}

type OrderDetails struct {
	CustomerId   string         `json:"customerId" bson:"customerId"`
	CustomerName string         `json:"customerName" bson:"customerName"`
	Products     []ProductModel `json:"products" bson:"products"`
	Address      string         `json:"address" bson:"address"`
}
