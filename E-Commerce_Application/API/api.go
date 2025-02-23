package API

import (
	"fmt"
	"module/Cassandra"
	"module/Database"
	"module/OrderService"
	"module/UserService"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("mysecretkey")
var lc *UserService.LoginCred
var ud *UserService.UserDetails
var brokers = []string{"localhost:9092"}

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
			_, err1 := Database.GetUserByUserDeatils(Database.DbPool, k.UserName)
			if err1 != nil {
				b.JSON(http.StatusBadGateway, err1)
				return
			}
			b.JSON(http.StatusOK, gin.H{"token": token})
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

	o1 := OrderService.Order{}
	if err1 := b.ShouldBindJSON(&o1); err1 != nil {
		b.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := OrderService.AddOrderToList(o1)
	if err != nil {
		b.JSON(http.StatusBadRequest, err)
		return
	}

	k := OrderService.CreateOrder(o1)
	err2 := Cassandra.AddOrder(Cassandra.Session, *k)
	if err2 != nil {
		b.JSON(http.StatusNotFound, err2)
		return
	}

	fmt.Println(*k)

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
