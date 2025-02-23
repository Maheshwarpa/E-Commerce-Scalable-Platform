package OrderService

import (
	"fmt"
	"module/Database"
	"module/ProductService"
	"strconv"
	"time"
)

type Order struct {
	Product_Id string `json:"product_id"`
	Count      int    `json:"count"`
}

var OrderList []Order

// Create a variable
var x int = 1

// Create a pointer to x
var ind *int = &x

type CompleteOrder struct {
	OrderId     string `json:"order_id"`
	PlacedOrder Order  `json:"placedorder"`
	OrderDate   string `json:"orderdate"`
}

type OrderResponse struct {
	OrderId     string `json:"order_id"`
	PlacedOrder Order  `json:"placedorder"`
	OrderDate   string `json:"orderdate"`
}

type FinalOrder struct {
	OrderId     string
	OrderStatus string
}

var FullOrderList []CompleteOrder
var FinalOrderList []FinalOrder

func CreateOrder(Ord Order) *CompleteOrder {
	var oid string
	currentTime := time.Now()
	od := currentTime.Format(time.DateOnly)
	//od := time.DateOnly
	k := strconv.Itoa(*ind)
	oid = "O" + k + Ord.Product_Id
	*ind = *ind + 1
	nco := CompleteOrder{oid, Ord, od}
	FullOrderList = append(FullOrderList, nco)
	return &nco
}

func GetStockStatus(Ord Order) bool {
	flag := false
	for ind, k := range ProductService.Inventory {
		if k.ProductID == Ord.Product_Id {
			ProductService.Inventory[ind].Quantity -= Ord.Count
			if ProductService.Inventory[ind].Quantity > 0 {
				return !flag
			}
		}
	}
	return flag
}

func ValidateOrder(Ord Order) (bool, bool) {
	flag := false
	for _, k := range ProductService.Inventory {
		if k.ProductID == Ord.Product_Id {
			return !flag, GetStockStatus(Ord)
		}
	}
	return flag, GetStockStatus(Ord)
}

func AddOrderToList(Ord Order) error {
	flag1, flag2 := ValidateOrder(Ord)

	switch flag1 {
	case true:
		switch flag2 {
		case true:
			// Do something when both are true
			OrderList = append(OrderList, Ord)
			return nil
		default:
			return fmt.Errorf("Product is Out Of Stock")

		}
	default:
		return fmt.Errorf("Invalid Product_Id Found!!!!")

	}
}

func CalculateTotal(ordl *CompleteOrder) (float64, error) {
	Total := 0.0

	val, err := Database.GetProdPrice(Database.DbPool, ordl.PlacedOrder.Product_Id)
	if err != nil {
		return Total, fmt.Errorf("Unable to fetch the price", err)
	}
	Total += (float64(ordl.PlacedOrder.Count) * val)

	return Total, nil
}
