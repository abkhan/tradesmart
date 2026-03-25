package models

import "time"

// Trade represents a unified model for orders and sales.
type Trade struct {
	OrderDate       time.Time `bson:"order_date" json:"order_date"`
	DeliveryDate    time.Time `bson:"delivery_date" json:"delivery_date"`
	OrderType       string    `bson:"order_type" json:"order_type"`
	SellerOrderID   string    `bson:"seller_order_id" json:"seller_order_id"`
	OrderID         string    `bson:"order_id" json:"order_id"`
	ProductName     string    `bson:"product_name" json:"product_name"`
	Inventory       string    `bson:"inventory" json:"inventory"`
	Status          string    `bson:"status" json:"status"`
	Tracking        string    `bson:"tracking" json:"tracking"`
	PackageCount    int       `bson:"package_count" json:"package_count"`
	Qty             int       `bson:"qty" json:"qty"`
	SalePrice1      float64   `bson:"sale_price_1" json:"sale_price_1"`
	SalePrice2      float64   `bson:"sale_price_2" json:"sale_price_2"`
	ShippingCharged float64   `bson:"shipping_charged" json:"shipping_charged"`
	AMZFee          float64   `bson:"amz_fee" json:"amz_fee"`
	TotalRate       float64   `bson:"total_rate" json:"total_rate"`
	TaxFees         float64   `bson:"tax_fees" json:"tax_fees"`
	Shipping        float64   `bson:"shipping" json:"shipping"`
	UnitCost        float64   `bson:"unit_cost" json:"unit_cost"`
	TotalCost       float64   `bson:"total_cost" json:"total_cost"`
	Refund          float64   `bson:"refund" json:"refund"`
	Profit          float64   `bson:"profit" json:"profit"`
	Loss            float64   `bson:"loss" json:"loss"`
	ROI             float64   `bson:"roi" json:"roi"`
	ReturnRefund    float64   `bson:"return_refund" json:"return_refund"`
	VeeqoCredits    float64   `bson:"veeqo_credits" json:"veeqo_credits"`
	ThreePLCost     float64   `bson:"3pl_cost" json:"3pl_cost"`
	Net             float64   `bson:"net" json:"net"`
	Created         time.Time `bson:"created" json:"created"`
}
