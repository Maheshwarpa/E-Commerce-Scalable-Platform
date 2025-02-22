package InventoryServices

import (
	"module/OrderService"
	"module/ProductService"
)

//var List OrderService.OrderList

func ApprovePayment(ol *[]OrderService.Order) {

	for _, k := range *ol {
		for key, v := range ProductService.Inventory {
			if v.ProductID == k.Product_Id {
				ProductService.Inventory[key].Quantity -= k.Count
			}
		}
	}

}

/*
func FailPayment(ol *OrderService.OrderList) {
	for _, k := range *ol {
		for key, v := range ProductService.Inventory {
			if v.ProductID == k.Product_Id {
				ProductService.Inventory[key].Quantity -= k.Count
			}
		}
	}
}
*/
