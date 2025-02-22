package Cassandra

import (
	"fmt"
	"log"
	"module/OrderService"

	gocql "github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectCDB() {
	// Connect to the Cassandra cluster
	cluster := gocql.NewCluster("localhost") // replace with your Cassandra IP
	cluster.Keyspace = "mykeyspace"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	Session = session
	//defer session.Close()
	fmt.Println("Cassandra DB ✅")
}

func CreateOrderTable(session *gocql.Session) error {
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			order_id TEXT PRIMARY KEY,
			product_id TEXT,
			count INT,
			order_date TEXT
		);`

	if err := session.Query(createTableQuery).Exec(); err != nil {
		return fmt.Errorf("Error creating table:", err)
	}
	fmt.Println("Order Table ✅")
	return nil
}

func AddOrder(session *gocql.Session, Ord OrderService.CompleteOrder) error {
	query := `INSERT INTO orders (order_id, product_id, count,order_date) VALUES (?, ?, ?, ?)`

	if err := session.Query(query, Ord.OrderId, Ord.PlacedOrder.Product_Id, Ord.PlacedOrder.Count, Ord.OrderDate).Exec(); err != nil {
		fmt.Println("Unable to add a record in Order Table in Cassandra ", err)
		return fmt.Errorf("Unable to add a record in Order Table in Cassandra ", err)
	}
	return nil
}
