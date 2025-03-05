package PaymentService

import (
	"fmt"
	"module/Card"
	"module/Database"
	"module/OrderService"
	"module/UserService"
	lg "module/logger"
)

func CheckEligibility(ordres OrderService.OrderResponse, sd *UserService.UserDetails) (bool, error) {
	lg.Log.Info("Entered CheckEligibility Function")
	var bal float64
	k, err := OrderService.CalculateTotal((*OrderService.CompleteOrder)(&ordres))
	if err != nil {
		lg.Log.Error("%s", err)
		fmt.Errorf("%s", err)
		return false, err
	}
	val, err2 := Database.GetUserBalance(sd.Cust_Id)
	if err2 != nil {
		lg.Log.Error("Unable to get the user balance from database ", err2)
		fmt.Println("Unable to get the user balance from database ", err2)
		return false, err2
	}
	lg.Log.Info("Total is :", k)
	fmt.Println("Total is :", k)
	lg.Log.Info("User Balance is :", val)
	fmt.Println("User Balance is :", val)
	lg.Log.Info("Card Number is :", Card.CardD)
	fmt.Println("Card Number is :", Card.CardD)
	lg.Log.Info("Card Flag is :", !(Card.CardD == ""))
	fmt.Println("Card Flag is :", Card.CardD == "")
	if val > k {
		bal = val - k
		sd.Cust_Bal = bal
		lg.Log.Info("After Deduction Balance :", bal)
		fmt.Println("After Deduction Balance :", bal)
		err1 := Database.UpdateUserBalance(Database.DbPool, sd.Cust_Id, bal)
		if err1 != nil {
			lg.Log.Error(err1)
			return false, err1
		}
		lg.Log.Info("User is eligible !!!")
		return true, nil
	} else if !(Card.CardD == "") {
		lg.Log.Info("User provided the card details, He is eligible !!!")
		return true, nil
	} else {
		lg.Log.Info("User is not eligible !!")
		return false, nil
	}
}
