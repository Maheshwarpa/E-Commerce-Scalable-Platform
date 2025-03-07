package OrderService

import (
	"fmt"
	"module/Database"
	"module/ProductService"
	lg "module/logger"
	"strconv"
	"time"
)

type Order struct {
	Product_Id string `json:"product_id"`
	Count      int    `json:"count"`
}

type PaymentDetails struct {
	CardDetails string `json:"carddetails"`
}

type OrderWithPay struct {
	Ord Order          `json:"order"`
	Crd PaymentDetails `json:"card"`
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
	OrderDts    Order
	OrderStatus string
	OrderDte    string
}

var FullOrderList []CompleteOrder
var FinalOrderList []FinalOrder

func CreateOrder(Ord Order) *CompleteOrder {
	lg.Log.Info("Started creating an order !")
	var oid string
	currentTime := time.Now()
	od := currentTime.Format(time.DateOnly)
	//od := time.DateOnly
	k := strconv.Itoa(*ind)
	oid = "O" + k + Ord.Product_Id
	*ind = *ind + 1
	nco := CompleteOrder{oid, Ord, od}
	FullOrderList = append(FullOrderList, nco)
	lg.Log.Info("Complete order has been created Successfully !")
	return &nco
}

func GetStockStatus(Ord Order) bool {
	lg.Log.Info("Entered getstock status function")
	flag := false
	for ind, k := range ProductService.Inventory {
		if k.ProductID == Ord.Product_Id {
			ProductService.Inventory[ind].Quantity -= Ord.Count
			if ProductService.Inventory[ind].Quantity > 0 {
				lg.Log.Info("Stock is available !!")
				_, err := Database.UpdateProdTbPostCompletion(k.ProductID, ProductService.Inventory[ind].Quantity)
				if err != nil {
					lg.Log.Info("Unable to update the value of the product quantity in product table")
				}
				return !flag
			} else {
				ProductService.Inventory[ind].Quantity += Ord.Count
				lg.Log.Info("Stock is NOT available !!")
				return flag
			}
		}
	}
	lg.Log.Error("Stock is NOT available !!")
	return flag
}

func ValidateOrder(Ord Order) (bool, bool) {
	lg.Log.Info("Entered the validate order function")
	flag := false
	for _, k := range ProductService.Inventory {
		if k.ProductID == Ord.Product_Id {
			lg.Log.Info("Valid Order !!")
			return !flag, GetStockStatus(Ord)
		}
	}
	lg.Log.Error("Order is not valid !!")
	return flag, GetStockStatus(Ord)
}

func AddOrderToList(Ord Order) error {
	lg.Log.Info("Entered Addorderlist function")
	flag1, flag2 := ValidateOrder(Ord)

	switch flag1 {
	case true:
		switch flag2 {
		case true:
			// Do something when both are true
			OrderList = append(OrderList, Ord)
			lg.Log.Info("Order has been added to the list")
			return nil
		default:
			lg.Log.Error("Product is Out of stock")
			return fmt.Errorf("Product is Out Of Stock")

		}
	default:
		lg.Log.Error("Invalid product_id found !!!")
		return fmt.Errorf("Invalid Product_Id Found!!!!")

	}
}

func CalculateTotal(ordl *CompleteOrder) (float64, error) {
	Total := 0.0
	lg.Log.Info("Entered the calculate total function !")
	val, err := Database.GetProdPrice(Database.DbPool, ordl.PlacedOrder.Product_Id)
	if err != nil {
		lg.Log.Error("Unable to fetch the price", err)
		return Total, fmt.Errorf("Unable to fetch the price", err)
	}
	Total += (float64(ordl.PlacedOrder.Count) * val)
	lg.Log.Info("Total is :", Total)
	return Total, nil
}

func UpdateStockInDb(pid string, cnt int, str string) (bool, error) {

	var numb int
	var flag bool = false
	for _, k := range ProductService.Inventory {
		if k.ProductID == pid {
			numb = k.Quantity
			//fmt.Println(numb)
		}
	}

	if str == "SUCCESS" {
		flag, err := Database.UpdateProdTbPostCompletion(pid, numb)
		return flag, err
	} else {
		for ind, k := range ProductService.Inventory {
			if k.ProductID == pid {
				ProductService.Inventory[ind].Quantity += cnt
				//fmt.Println(numb)
				flag = true
			}
		}
		return flag, nil
	}

}
