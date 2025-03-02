package PaymentService

import (
	"fmt"
	"module/Database"
	"module/OrderService"
	"module/UserService"
)

func CheckEligibility(ordres OrderService.OrderResponse, sd *UserService.UserDetails) (bool, error) {
	var bal float64
	k, err := OrderService.CalculateTotal((*OrderService.CompleteOrder)(&ordres))
	if err != nil {
		fmt.Errorf("%s", err)
		return false, err
	}
	fmt.Println("Total is :", k)
	fmt.Println("User Balance is :", sd.Cust_Bal)
	if sd.Cust_Bal > k {
		bal = sd.Cust_Bal - k
		fmt.Println("After Deduction Balance :", bal)
		err1 := Database.UpdateUserBalance(Database.DbPool, sd.Cust_Id, bal)
		if err1 != nil {
			return false, err1
		}
		return true, nil
	}

	return false, nil
}
