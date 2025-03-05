package API

import (
	"context"
	"fmt"
	"log"
	"module/Card"
	"module/Cassandra"
	"module/Consumer"
	"module/Database"
	"module/Database/files"
	"module/OrderService"
	"module/Publisher"
	"module/UserService"
	lg "module/logger"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
{
    "product_id": "P12387",
    "count": 1
}

{
"username":"custMah1",
"password": "1111"
}

{
    "order": {
        "product_id": "P12387",
        "count": 1
    },
    "card": {
         "carddetails": "a123456789"
    }
}

*/

var secretKey = []byte("mysecretkey")
var lc *UserService.LoginCred
var ud *UserService.UserDetails
var brokers = []string{"localhost:9092"}
var topic = "Orders"
var wg sync.WaitGroup
var SampleData *UserService.UserDetails
var x int = 0
var Ind *int = &x
var Cardd OrderService.PaymentDetails
var grpcClient files.DatabaseServiceClient

func initGRPCClient() {
	lg.Log.Info("Connect to the gRPC server")
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		lg.Log.Error("Failed to connect to gRPC server: %v", err)
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	grpcClient = files.NewDatabaseServiceClient(conn)
}

func StartServer() {
	lg.Log.Info("Starting API Server at port 8080")
	initGRPCClient()
	router := gin.Default()
	router.POST("/login", userlogin)

	// Protected Routes
	protected := router.Group("/SKHT")
	protected.Use(authMiddleware()) // Middleware applied here
	//otected.GET("/dashboard", protectedRoute)

	protected.POST("/Orders", saveItem)
	protected.GET("/Products", getProducts)
	protected.GET("/Products/category/:category", getCategory)
	protected.GET("/Products/uniquecategory", getUniqueCategory)
	protected.GET("/Products/brand/:brand", getBrand)
	protected.GET("/Products/uniquebrand", getUniqueBrand)
	protected.GET("/Orders/status/:status", getOrderByStatus)
	protected.GET("/Orders/date/:date", getOrderByDate)
	protected.GET("/Orders/productid/:product", getOrderByProductId)
	protected.GET("/Orders/List", getAllFinalOrders)
	router.Run("localhost:8080")
}

func userlogin(b *gin.Context) {
	lg.Log.Info("Enetered User login request !!")
	var req UserService.LoginCred
	flag := false
	if err := b.ShouldBindJSON(&req); err != nil {
		lg.Log.Error("error: Invalid request")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//grpcClient.CheckLoginCredentials(ctx, &files.LoginRequest{Username: req.UserName})
	//grpcClient.LoadLoginCred(ctx, &files.LoginCredRequest{Username: req.UserName, Password: req.Password})
	_, checkerr := grpcClient.CheckLoginCredentials(ctx, &files.LoginRequest{Username: req.UserName}) //Database.CheckLoginCredentials(req.UserName)
	if checkerr == nil {
		lg.Log.Error("error is:", checkerr)
		fmt.Println("error is:", checkerr)
		_, lerr := grpcClient.LoadLoginCred(ctx, &files.LoginCredRequest{Username: req.UserName, Password: req.Password}) //Database.LoadLoginCred(req)
		if lerr != nil {
			lg.Log.Error("error : Unable to load Data")
			b.JSON(http.StatusNotFound, gin.H{"error": "Unable to load Data"})
			return
		}
	}

	lcl, err0 := grpcClient.GetAllLoginCred(ctx, &files.Empty{}) //Database.GetAllLoginCred()
	if err0 != nil {
		lg.Log.Error("error : Unable to fetch login Data")
		b.JSON(http.StatusNotFound, gin.H{"error": "Unable to fetch login Data"})
		return
	}
	fmt.Println("Hi", lcl)
	for _, lk := range lcl.Credentials {
		var lc UserService.LoginCred
		lc.UserName = lk.Username
		lc.Password = lk.Password
		UserService.LoginCredList = append(UserService.LoginCredList, lc)
	}

	fmt.Println(UserService.UserDetailList)

	for _, k := range UserService.UserDetailList {
		fmt.Println(k.UserName)
		_, err4 := grpcClient.GetUserByUserDetails(ctx, &files.UserRequest{Username: k.UserName})
		fmt.Println("EEError is", err4)
		if err4 != nil {
			lg.Log.Error("ENtered loading state")
			fmt.Println("ENtered loading state")
			err3 := Database.LoadUserData(k)
			if err3 != nil {
				lg.Log.Error("Unable to Load User data to Database ", err3)
				fmt.Errorf("Unable to Load User data to Database ", err3)
				return
			}
			if req.UserName == k.UserName {
				lg.Log.Info("Present")
				//fmt.Println("Present")
			}
		}

	}

	kl, err4 := grpcClient.GetAllUserData(ctx, &files.Empty{}) //Database.GetALLUserData()
	if err4 != nil {
		lg.Log.Error("Unable to fetch error", err4)
		fmt.Errorf("Unable to fetch error", err4)
		return
	}

	fmt.Println("GeALL func is:")
	for _, k := range kl.Users {

		fmt.Println(k)
	}
	lg.Log.Info("Cred is :", req)
	//fmt.Println("Cred is :", req)
	lg.Log.Info("List is ", UserService.LoginCredList)
	//fmt.Println("List is ", UserService.LoginCredList)
	// Dummy authentication check
	for _, k := range UserService.LoginCredList {
		if k.UserName == req.UserName && k.Password == req.Password {
			flag = true

			token, err := generateToken(req.UserName)
			if err != nil {
				b.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
				return
			}
			//UserService.StoreVault(token, k)

			sd, err1 := Database.GetUserByUserDeatils(k.UserName)
			fmt.Println(sd)
			if err1 != nil {
				b.JSON(http.StatusBadGateway, err1)
				return
			}
			SampleData = &sd
			Database.SampleData = SampleData

			/*	k, err0 := Database.GetALLUserData()
				if err0 != nil {
					return
				}
				fmt.Println(k)
				fmt.Println(req.UserName)
				ljk, err9 := Database.GetUserByUserDeatils(req.UserName)
				if err9 != nil {
					return
				}
				fmt.Println(ljk)*/

			b.JSON(http.StatusOK, gin.H{"token": token, "SampleData": SampleData})
			lg.Log.Info("**************************************************************\n")
		}
	}
	if !flag {
		lg.Log.Error("error : Invalid credentials")
		b.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return

	}

}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})
	return token.SignedString(secretKey)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		// Ensure the token starts with "Bearer "
		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Extract the token part after "Bearer "
		tokenString := authHeader[len(prefix):]

		// Validate Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next() // Proceed to the next handler
	}
}

func saveItem(b *gin.Context) {
	lg.Log.Info("Entered SaveItem function !!")
	//	wg := sync.WaitGroup{}
	//brokers := []string{"localhost:9092"}
	err9 := Publisher.InitKafkaProducer(brokers)
	if err9 != nil {
		lg.Log.Error("Kafka producer initialization failed: %v", err9)
		log.Fatalf("Kafka producer initialization failed: %v", err9)
	}
	//fmt.Println("Waiting for waitgroup")
	//b.JSON(http.StatusAccepted, gin.H{"message": "Waiting for waitgroup at beginning"})
	//wg.Wait()

	o1 := OrderService.Order{}
	owp := OrderService.OrderWithPay{}
	fmt.Println("SampleData is ", Database.SampleData)
	if Database.SampleData.Cust_Bal > 0 {
		fmt.Println("111111")
		if err1 := b.ShouldBindJSON(&o1); err1 != nil {
			lg.Log.Error("error : Invalid JSON data")
			b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
			return
		}
		if (o1 == OrderService.Order{}) {
			lg.Log.Error("error : Invalid JSON data")
			b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Data, Please provide the input in proper format"})
			return
		}
	} else {
		fmt.Println("222222")
		if err98 := b.ShouldBindJSON(&owp); err98 != nil {
			lg.Log.Error("error : Invalid JSON data")
			b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Data, Please provide the input in proper format"})
			return
		}
		o1 = owp.Ord
		Cardd = owp.Crd
		Card.CardD = owp.Crd.CardDetails
		fmt.Println("Card Number is :", owp.Crd.CardDetails)
		fmt.Println("owp is", owp)
		if (owp == OrderService.OrderWithPay{}) {
			lg.Log.Error("error : Invalid JSON data")
			b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Data, Please provide the input in proper format"})
			return
		}
	}

	err8 := Database.CheckValidProduct(Database.DbPool, o1.Product_Id)
	if err8 != nil {
		lg.Log.Error("error : Invalid Product Id found")
		b.JSON(http.StatusNotFound, gin.H{"error": "Invalid Product Id found"})
		return
	}

	err := OrderService.AddOrderToList(o1)
	if err != nil {
		lg.Log.Error("error", err)
		b.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	k := OrderService.CreateOrder(o1)
	err2 := Cassandra.AddOrder(Cassandra.Session, *k)
	if err2 != nil {
		lg.Log.Error("error : Unable to add Order to Database")
		b.JSON(http.StatusNotFound, gin.H{"error": "Unable to add Order to Database"})
		return
	}

	fmt.Println(*k)
	lg.Log.Info("Before Add: WaitGroup count: %d\n", wg)
	//fmt.Printf("Before Add: WaitGroup count: %d\n", wg)
	wg.Add(1)
	lg.Log.Info("After Add: WaitGroup count: %d\n", wg)
	//fmt.Printf("After Add: WaitGroup count: %d\n", wg)
	go func() {
		//defer wg.Done()
		Consumer.ConsumeMessages(brokers, "orders", "Group_Id", &wg)
	}()

	//defer wg.Done()

	/*
		go Consumer.ConsumeMessages(brokers, topic, "Group_Id", &wg)
		wg.Add(1)
		*Ind += 1
		go Publisher.PublishOrderCreatedEvent(brokers, topic, k, &wg)
		Publisher.Ipnd = Ind
	*/

	err = Publisher.PublishMessage("orders", k)
	if err != nil {
		lg.Log.Error("Failed to send order message: %v", err)
		log.Printf("Failed to send order message: %v", err)
	}
	lg.Log.Info("Producer sent an order message.")
	//fmt.Println("Producer sent an order message.")
	lg.Log.Info("*******************Waiting for Go ROutine main to join")
	//fmt.Println("*******************Waiting for Go ROutine main to join")
	wg.Wait()
	lg.Log.Info("Go ROutine main has joined******************")
	//fmt.Println("Go ROutine main has joined******************")
	fmt.Println(OrderService.FinalOrderList)
	lg.Log.Info("********************************************************************\n\n\n\n")
	b.IndentedJSON(http.StatusOK, *k)

}

func getProducts(b *gin.Context) {
	lg.Log.Info("entered the get products function post request call")
	k, err := Database.GetAllProdData(Database.DbPool)
	if err != nil {
		lg.Log.Error("error : Unable to fetch all product data")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch all product data"})
	}

	b.IndentedJSON(http.StatusOK, k)
}

func getCategory(b *gin.Context) {
	lg.Log.Info("entered the get category function post request call")
	category := b.Param("category")

	list, err := Database.GetProdByCateg(Database.DbPool, category)
	if err != nil {
		lg.Log.Error("error : Unable to fetch the data by category")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the data by category"})
	}

	b.JSON(http.StatusOK, list)
}

func getBrand(b *gin.Context) {
	lg.Log.Info("entered the get brand function post request call")
	brand := b.Param("brand")

	list, err := Database.GetProdByBrand(Database.DbPool, brand)
	if err != nil {
		lg.Log.Error("error : Unable to fetch the data by brand")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the data by brand"})
	}

	b.JSON(http.StatusOK, list)
}

func getUniqueBrand(b *gin.Context) {
	//brand := b.Param("brand")
	lg.Log.Info("entered the get unique brand function post request call")
	list, err := Database.GetUniqueBrand(Database.DbPool)
	if err != nil {
		lg.Log.Error("error : Unable to fetch the unique data by brand")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the unique data by brand"})
	}

	b.JSON(http.StatusOK, list)
}

func getUniqueCategory(b *gin.Context) {
	//brand := b.Param("brand")
	lg.Log.Info("entered the get unique category function post request call")
	list, err := Database.GetUniqueCategory(Database.DbPool)
	if err != nil {
		lg.Log.Error("error : Unable to fetch the unique data by category")
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the unique data by category"})
	}

	b.JSON(http.StatusOK, list)
}

func getOrderByStatus(b *gin.Context) {

}

func getOrderByDate(b *gin.Context) {

}

func getOrderByProductId(b *gin.Context) {

}

func getAllFinalOrders(b *gin.Context) {

}
