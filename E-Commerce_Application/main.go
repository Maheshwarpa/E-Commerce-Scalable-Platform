package main

import (
	"fmt"
	"log"
	"module/API"
	"module/Cassandra"
	"module/Database"
	"module/UserService"

	//ns "module/NotificationService"
	"module/ProductService"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func init() {
	var err error
	dbpool, err = Database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Call table creation function once the DB connection is established
	Database.CreateProductTable(dbpool)
	time.Sleep(3 * time.Second)
	ProductService.LoadProducts()
	Database.LoadProdData(dbpool)
	Database.CreateUserTable(dbpool)

	Cassandra.ConnectCDB()
	Cassandra.CreateOrderTable(Cassandra.Session)
}

func main() {
	fmt.Println("Hello World")
	var opt int
	fmt.Println("Welcome to SKHT E-Commerce Site\nPlease login,\to get into our wonderful shopping experience")
	for {
		fmt.Println("PLease select the following menu\n1. Login\n2. New User")
		fmt.Scan(&opt)
		if opt == 1 {
			API.StartServer()
		} else if opt == 2 {
			UserService.CreateUser()
		} else {
			return
		}
	}

	//ns.SendSuccessEmail()
}
