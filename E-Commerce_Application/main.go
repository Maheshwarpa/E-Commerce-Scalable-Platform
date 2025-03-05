package main

import (
	"fmt"
	"log"
	"module/API"
	"module/Cassandra"
	"module/Database"
	"module/ProductService"
	"module/UserService"
	lg "module/logger"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func init() {
	lg.InitLogger()
	var err error
	lg.Log.Info("Connecting to the Postgresql database")
	dbpool, err = Database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
		return
	}
	lg.Log.Info("Successfully Connected to the Postgresql Database")
	Database.DbPool = dbpool
	lg.Log.Info("Assigned dbpool value to the global variable")
	// Call table creation function once the DB connection is established
	lg.Log.Info("Connecting to the Product table in  database")
	Database.CreateProductTable(dbpool)
	lg.Log.Info("Successfully Connected to the Product table in  database")
	time.Sleep(3 * time.Second)
	lg.Log.Info("Loading products to the Products table in database")
	ProductService.LoadProducts()
	Database.LoadProdData(dbpool)
	lg.Log.Info("Creating the User table in database")
	Database.CreateUserTable(dbpool)
	lg.Log.Info("Creating login table in database")
	Database.CreateLoginTable(dbpool)
	lg.Log.Info("Creating a connection to cassandra database")
	Cassandra.ConnectCDB()
	lg.Log.Info("Creating order table in cassandra database")
	Cassandra.CreateOrderTable(Cassandra.Session)
	lg.Log.Info("Creating final order table in cassandra database")
	Cassandra.CreateFinalOrderTable(Cassandra.Session)

	lg.Log.Info("Starting GRPC Server as a Go Routinue")
	go Database.StartServer()
	time.Sleep(3 * time.Second)
}

func main() {
	//fmt.Println("Hello World")
	//ns.SendSuccessEmail()
	var opt int
	lg.Log.Info("Application Started Here !")
	fmt.Println("Welcome to SKHT E-Commerce Site\nPlease login to get into our wonderful shopping experience")
	for {
		fmt.Println("PLease select the following menu\n1. Login\n2. New User")
		fmt.Scan(&opt)
		if opt == 1 {
			lg.Log.Info("Connecting to 8080 port through RESTAPI")
			API.StartServer()
		} else if opt == 2 {
			lg.Log.Info("Creating a new User !")
			UserService.CreateUser()
		} else {
			return
		}
	}

	//ns.SendSuccessEmail()
}
