package UserService

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/vault/api"
)

type UserDetails struct {
	Cust_Id    int     `json:"cust_is"`
	Cust_Name  string  `json:"cust_name"`
	Cust_Email string  `json:"cust_email"`
	Cust_PNum  string  `json:"cust_pnum"`
	Cust_Bal   float64 `json:"cust_bal"`
}

var UserDetailList []UserDetails

var x int = 1
var ind *int = &x

func CreateUser() {
	//reader := bufio.NewReader(os.Stdin)
	var nm, em, pnum = "", "", ""
	var bal float64
	fmt.Println("Please enter the follwoing details to set up your account")
	fmt.Println("Please enter the Customer Name: ")
	fmt.Scan(&nm)
	fmt.Println("Please enter the Customer Email:")
	fmt.Scan(&em)
	fmt.Println("Please enter the Customer PhoneNumber")
	fmt.Scan(&pnum)
	fmt.Println("Please enter the Customer Balance")
	fmt.Scan(&bal)
	Cust := UserDetails{*ind, nm, em, pnum, bal}
	*ind = *ind + 1
	CreateCredentials(Cust)
	UserDetailList = append(UserDetailList, Cust)
	fmt.Println("Customer_Account & Credentials Created Successfully ✅")
}

func CreateCredentials(cust UserDetails) {
	var pass string
	var unm string = "cust"
	unm += cust.Cust_Name[:3] + strconv.Itoa(cust.Cust_Id)
	fmt.Println("Please make a note that UserName for your Account is:  ", unm)
	fmt.Println("Please set up your password as per NIST standards ")
	fmt.Println("Please enter the password below :")
	fmt.Scan(&pass)
	lc := LoginCred{unm, pass}
	LoginCredList = append(LoginCredList, lc)
}

func StoreVault(str string, lc LoginCred) {
	// Create Vault client

	client, err := api.NewClient(&api.Config{Address: "http://127.0.0.1:8200"})
	if err != nil {
		log.Fatal(err)
	}
	client.SetNamespace("admin")

	// Set authentication token (replace with your actual token)
	client.SetToken(str)

	// Store credentials
	data := map[string]interface{}{
		"username": lc.UserName,
		"password": lc.Password,
	}

	_, err = client.Logical().Write("secret/data/myapp", data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Credentials stored successfully")

	// Retrieve credentials
	secret, err := client.Logical().Read("secret/data/myapp")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Retrieved secret:", secret.Data)
}
