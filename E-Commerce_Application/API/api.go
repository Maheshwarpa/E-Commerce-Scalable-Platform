package API

import (
	"fmt"
	"log"
	"module/Cassandra"
	"module/Consumer"
	"module/Database"
	"module/OrderService"
	"module/Publisher"
	"module/UserService"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func StartServer() {
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
	router.Run("localhost:8080")
}

func userlogin(b *gin.Context) {
	var req UserService.LoginCred
	flag := false
	if err := b.ShouldBindJSON(&req); err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	checkerr := Database.CheckLoginCredentials(Database.DbPool, req.UserName)
	if checkerr != nil {
		lerr := Database.LoadLoginCred(Database.DbPool, req)
		if lerr != nil {
			b.JSON(http.StatusNotFound, gin.H{"error": "Unable to load Data"})
			return
		}
	}

	lcl, err0 := Database.GetAllLoginCred(Database.DbPool)
	if err0 != nil {
		b.JSON(http.StatusNotFound, gin.H{"error": "Unable to fetch login Data"})
		return
	}
	for _, k := range lcl {
		UserService.LoginCredList = append(UserService.LoginCredList, k)
	}
	fmt.Println(UserService.UserDetailList)
	for _, k := range UserService.UserDetailList {
		err3 := Database.LoadUserData(Database.DbPool, k)
		if err3 != nil {
			fmt.Errorf("Unable to Load User data to Database ", err3)
			return
		}
		if req.UserName == k.UserName {
			fmt.Println("Present")
		}
	}

	kl, err4 := Database.GetALLUserData(Database.DbPool)
	if err4 != nil {
		fmt.Errorf("Unable to fetch error", err4)
	}
	fmt.Println("GeALL func is:")
	for _, k := range kl {

		fmt.Println(k)
	}

	fmt.Println("Cred is :", req)
	fmt.Println("List is ", UserService.LoginCredList)
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

			sd, err1 := Database.GetUserByUserDeatils(Database.DbPool, k.UserName)
			fmt.Println(sd)
			if err1 != nil {
				b.JSON(http.StatusBadGateway, err1)
				return
			}
			SampleData = &sd
			Database.SampleData = SampleData

			b.JSON(http.StatusOK, gin.H{"token": token, "SampleData": SampleData})
		}
	}
	if !flag {
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

	//	wg := sync.WaitGroup{}
	//brokers := []string{"localhost:9092"}
	err9 := Publisher.InitKafkaProducer(brokers)
	if err9 != nil {
		log.Fatalf("Kafka producer initialization failed: %v", err9)
	}
	//fmt.Println("Waiting for waitgroup")
	//b.JSON(http.StatusAccepted, gin.H{"message": "Waiting for waitgroup at beginning"})
	//wg.Wait()

	o1 := OrderService.Order{}
	if err1 := b.ShouldBindJSON(&o1); err1 != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err8 := Database.CheckValidProduct(Database.DbPool, o1.Product_Id)
	if err8 != nil {
		b.JSON(http.StatusNotFound, gin.H{"error": "Invalid Product Id found"})
		return
	}

	err := OrderService.AddOrderToList(o1)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	k := OrderService.CreateOrder(o1)
	err2 := Cassandra.AddOrder(Cassandra.Session, *k)
	if err2 != nil {
		b.JSON(http.StatusNotFound, gin.H{"error": "Unable to add Order to Database"})
		return
	}

	fmt.Println(*k)
	fmt.Printf("Before Add: WaitGroup count: %d\n", wg)
	wg.Add(1)
	fmt.Printf("After Add: WaitGroup count: %d\n", wg)
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
		log.Printf("Failed to send order message: %v", err)
	}

	fmt.Println("Producer sent an order message.")

	fmt.Println("*******************Waiting for Go ROutine main to join")
	wg.Wait()
	fmt.Println("Go ROutine main has joined******************")
	b.IndentedJSON(http.StatusOK, *k)

}

func getProducts(b *gin.Context) {

	k, err := Database.GetAllProdData(Database.DbPool)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch all product data"})
	}

	b.IndentedJSON(http.StatusOK, k)
}

func getCategory(b *gin.Context) {
	category := b.Param("category")

	list, err := Database.GetProdByCateg(Database.DbPool, category)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the data by category"})
	}

	b.JSON(http.StatusOK, list)
}

func getBrand(b *gin.Context) {
	brand := b.Param("brand")

	list, err := Database.GetProdByBrand(Database.DbPool, brand)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the data by brand"})
	}

	b.JSON(http.StatusOK, list)
}

func getUniqueBrand(b *gin.Context) {
	//brand := b.Param("brand")

	list, err := Database.GetUniqueBrand(Database.DbPool)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the unique data by brand"})
	}

	b.JSON(http.StatusOK, list)
}

func getUniqueCategory(b *gin.Context) {
	//brand := b.Param("brand")

	list, err := Database.GetUniqueCategory(Database.DbPool)
	if err != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Unable to fetch the unique data by category"})
	}

	b.JSON(http.StatusOK, list)
}
