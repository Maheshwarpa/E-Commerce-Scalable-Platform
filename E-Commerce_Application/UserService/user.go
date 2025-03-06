package UserService

import (
	"fmt"
	lg "module/logger"
	"strconv"
	//"module/Database"
)

var Sd UserDetails

type UserDetails struct {
	Cust_Id    int     `json:"cust_id"`
	Cust_Name  string  `json:"cust_name"`
	Cust_Email string  `json:"cust_email"`
	Cust_PNum  string  `json:"cust_pnum"`
	Cust_Bal   float64 `json:"cust_bal"`
	UserName   string  `json:"username"`
}

/*
{
    "username": "custMah1",
    "password": "1111"
}
*/

var UserDetailList []UserDetails

var x int = 1
var ind *int = &x

func CreateUser() {

	//reader := bufio.NewReader(os.Stdin)
	var nm, em, pnum = "", "", ""
	var id int
	var bal float64
	fmt.Println("Please enter the follwoing details to set up your account")
	fmt.Println("Please enter the Customer Id: ")
	fmt.Scan(&id)
	fmt.Println("Please enter the Customer Name: ")
	fmt.Scan(&nm)
	fmt.Println("Please enter the Customer Email:")
	fmt.Scan(&em)
	fmt.Println("Please enter the Customer PhoneNumber")
	fmt.Scan(&pnum)
	fmt.Println("Please enter the Customer Balance")
	fmt.Scan(&bal)
	Cust := UserDetails{id, nm, em, pnum, bal, ""}
	//*ind = *ind + 2
	lk := CreateCredentials(Cust)
	Cust.UserName = lk.UserName
	UserDetailList = append(UserDetailList, Cust)
	/*err := Database.LoadUserData(Database.DbPool, Cust)
	if err != nil {
		fmt.Errorf("Unable to Load User data to Database ", err)
		return
	}*/
	Sd = Cust
	//Database.LoadUserData(Database.DbPool,Sd)
	lg.Log.Info("User created successfully ✅")
	fmt.Println("Customer_Account & Credentials Created Successfully ✅")
}

func CreateCredentials(cust UserDetails) LoginCred {
	lg.Log.Info("Entered the create credentials function")
	var pass string
	var unm string = "cust"
	unm += cust.Cust_Name[:3] + strconv.Itoa(cust.Cust_Id)
	fmt.Println("Please make a note that UserName for your Account is:  ", unm)
	fmt.Println("Please set up your password as per NIST standards ")
	fmt.Println("Please enter the password below :")
	fmt.Scan(&pass)
	lc := LoginCred{unm, pass}
	LoginCredList = append(LoginCredList, lc)
	lg.Log.Info("Login credentials created successfully ✅")
	return lc
}

/*
func StoreVault(str string, lc LoginCred) {
	// Create Vault client
	lg.Log.Info("Entered Store Vault Section !!")
	client, err := api.NewClient(&api.Config{Address: "http://127.0.0.1:8200"})
	if err != nil {
		log.Fatal(err)
	}
	client.SetNamespace("admin")

	// Set authentication token (replace with your actual token)
	client.SetToken("hvs.femyXdNg4m2pl3rgsy0sqYVJ")

	// Store credentials
	fmt.Println(lc)
	lg.Log.Info("Login details is :", lc)
	data := map[string]interface{}{
		"data": map[string]interface{}{ // Vault expects the data to be nested under "data"
			"username": lc.UserName,
			"password": lc.Password,
		},
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
*/
