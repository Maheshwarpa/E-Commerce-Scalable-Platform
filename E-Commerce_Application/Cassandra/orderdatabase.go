package Cassandra

import (
	"fmt"
	"log"
	"module/OrderService"
	lg "module/logger"

	gocql "github.com/gocql/gocql"
)

var Session *gocql.Session

func ConnectCDB() {
	lg.Log.Info("Creating Cassandra database connection!!")

	// Connect to the Cassandra cluster (without specifying the keyspace initially)
	cluster := gocql.NewCluster("localhost") // replace with your Cassandra IP
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum

	// Create session (without specifying the keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	Session = session

	// Create the keyspace 'mykeyspace' if it doesn't exist
	createKeyspaceQuery := `
		CREATE KEYSPACE IF NOT EXISTS mykeyspace 
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`
	err = Session.Query(createKeyspaceQuery).Exec()
	if err != nil {
		log.Fatalf("Error creating keyspace: %v", err)
	}
	lg.Log.Info("Keyspace 'mykeyspace' created or already exists ✅")

	// Close the session that doesn't have the keyspace
	Session.Close()

	// Now, create a new session with the 'mykeyspace' specified
	cluster.Keyspace = "mykeyspace"
	session, err = cluster.CreateSession() // Re-create session with the keyspace specified
	if err != nil {
		log.Fatalf("Error creating session with keyspace: %v", err)
	}

	Session = session

	lg.Log.Info("Cassandra DB connected with 'mykeyspace' ✅")
	fmt.Println("Cassandra DB connected with 'mykeyspace' ✅")
}

func ConnectCDBd() {
	lg.Log.Info("Creating Cassandra database connection!!")
	// Connect to the Cassandra cluster (without specifying the keyspace initially)
	cluster := gocql.NewCluster("localhost") // replace with your Cassandra IP
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum

	// Create session (without specifying the keyspace)
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	Session = session

	// Create the keyspace 'mykeyspace' if it doesn't exist
	createKeyspaceQuery := `
		CREATE KEYSPACE IF NOT EXISTS mykeyspace 
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`
	err = Session.Query(createKeyspaceQuery).Exec()
	if err != nil {
		log.Fatalf("Error creating keyspace: %v", err)
	}
	lg.Log.Info("Keyspace 'mykeyspace' created or already exists ✅")

	cluster.Keyspace = "mykeyspace"        // Set the keyspace here
	session, err = cluster.CreateSession() // Re-create session with the keyspace specified
	if err != nil {
		log.Fatalf("Error creating session with keyspace: %v", err)
	}

	Session = session

	lg.Log.Info("Cassandra DB connected with 'mykeyspace' ✅")
	fmt.Println("Cassandra DB connected with 'mykeyspace' ✅")
}

func ConnectCDBs() {
	lg.Log.Info("Creating cassandra database connection !!")
	// Connect to the Cassandra cluster
	cluster := gocql.NewCluster("localhost") // replace with your Cassandra IP
	cluster.Keyspace = "mykeyspace"
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	Session = session
	//defer session.Close()
	/*// Create the keyspace 'mykeyspace' if it doesn't exist
	createKeyspaceQuery := `
		CREATE KEYSPACE IF NOT EXISTS mykeyspace
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`
	err = Session.Query(createKeyspaceQuery).Exec()
	if err != nil {
		log.Fatalf("Error creating keyspace: %v", err)
	}
	lg.Log.Info("Keyspace 'mykeyspace' created or already exists ✅")

	// Set the keyspace for the session

	err = Session.SetKeyspace("mykeyspace")
	if err != nil {
		log.Fatalf("Error setting keyspace: %v", err)
	}*/

	lg.Log.Info("Cassandra DB ✅")
	fmt.Println("Cassandra DB ✅")
}

func CreateOrderTable(session *gocql.Session) error {
	lg.Log.Info("Creating Order table in cassandra database !!")
	createTableQuery := `
		CREATE TABLE IF NOT EXISTS orders (
			order_id TEXT PRIMARY KEY,
			product_id TEXT,
			count INT,
			order_date TEXT
		);`

	if err := session.Query(createTableQuery).Exec(); err != nil {
		lg.Log.Error("Error creating table:", err)
		return fmt.Errorf("Error creating table:", err)
	}
	lg.Log.Info("Order Table ✅")
	fmt.Println("Order Table ✅")
	return nil
}

func CreateFinalOrderTable(session *gocql.Session) error {
	lg.Log.Info("Creating FinalOrder table in Cassandra database!!")

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS finalorder (
			order_id TEXT PRIMARY KEY,
			product_id TEXT,
			count INT,
			order_date TEXT,
			order_status TEXT
		);`

	if err := session.Query(createTableQuery).Exec(); err != nil {
		lg.Log.Error("Error creating FinalOrder table:", err)
		return fmt.Errorf("error creating FinalOrder table: %v", err)
	}

	lg.Log.Info("FinalOrder Table ✅")
	fmt.Println("FinalOrder Table ✅")
	return nil
}

func AddOrder(session *gocql.Session, Ord OrderService.CompleteOrder) error {
	lg.Log.Info("Adding order to cassandra database !!")
	query := `INSERT INTO orders (order_id, product_id, count,order_date) VALUES (?, ?, ?, ?)`

	if err := session.Query(query, Ord.OrderId, Ord.PlacedOrder.Product_Id, Ord.PlacedOrder.Count, Ord.OrderDate).Exec(); err != nil {
		lg.Log.Error("Unable to add a record in Order Table in Cassandra ", err)
		fmt.Println("Unable to add a record in Order Table in Cassandra ", err)
		return fmt.Errorf("Unable to add a record in Order Table in Cassandra ", err)
	}
	return nil
}

func AddFinalOrder(session *gocql.Session, order OrderService.FinalOrder) error {
	lg.Log.Info("Adding order to the finalorder table in Cassandra database!!")

	query := `INSERT INTO finalorder (order_id, product_id, count, order_date, order_status) 
	          VALUES (?, ?, ?, ?, ?)`

	if err := session.Query(query,
		order.OrderId,
		order.OrderDts.Product_Id,
		order.OrderDts.Count,
		order.OrderDte,
		order.OrderStatus,
	).Exec(); err != nil {
		lg.Log.Error("Unable to add a record in FinalOrder Table in Cassandra", err)
		fmt.Println("Unable to add a record in FinalOrder Table in Cassandra", err)
		return fmt.Errorf("unable to add a record in FinalOrder Table in Cassandra: %v", err)
	}

	lg.Log.Info("Order added successfully to FinalOrder Table ✅")
	fmt.Println("Order added successfully to FinalOrder Table ✅")
	return nil
}

func GetOrdersByProductId(productId string) ([]OrderService.FinalOrder, error) {
	lg.Log.Info("Fetching orders by product_id from the finalorder table!!")

	query := `SELECT order_id, product_id, count, order_date, order_status FROM finalorder WHERE product_id = ?ALLOW FILTERING`
	iter := Session.Query(query, productId).Iter()

	var orders []OrderService.FinalOrder
	var order OrderService.FinalOrder

	for iter.Scan(&order.OrderId, &order.OrderDts.Product_Id, &order.OrderDts.Count, &order.OrderDte, &order.OrderStatus) {
		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		lg.Log.Error("Error fetching orders by product_id:", err)
		return nil, fmt.Errorf("error fetching orders by product_id: %v", err)
	}

	lg.Log.Info("Successfully fetched orders by product_id ✅")
	return orders, nil
}

func GetOrdersByDate(orderDate string) ([]OrderService.FinalOrder, error) {
	lg.Log.Info("Fetching orders by order_date from the finalorder table!!")

	query := `SELECT order_id, product_id, count, order_date, order_status FROM finalorder WHERE order_date = ?ALLOW FILTERING`
	iter := Session.Query(query, orderDate).Iter()

	var orders []OrderService.FinalOrder
	var order OrderService.FinalOrder

	for iter.Scan(&order.OrderId, &order.OrderDts.Product_Id, &order.OrderDts.Count, &order.OrderDte, &order.OrderStatus) {
		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		lg.Log.Error("Error fetching orders by order_date:", err)
		return nil, fmt.Errorf("error fetching orders by order_date: %v", err)
	}

	lg.Log.Info("Successfully fetched orders by order_date ✅")
	return orders, nil
}

func GetOrdersByStatus(status string) ([]OrderService.FinalOrder, error) {
	lg.Log.Info("Fetching orders by status from the finalorder table!!")

	query := `SELECT order_id, product_id, count, order_date, order_status FROM finalorder WHERE order_status = ?ALLOW FILTERING`
	iter := Session.Query(query, status).Iter()

	var orders []OrderService.FinalOrder
	var order OrderService.FinalOrder

	for iter.Scan(&order.OrderId, &order.OrderDts.Product_Id, &order.OrderDts.Count, &order.OrderDte, &order.OrderStatus) {
		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		lg.Log.Error("Error fetching orders by status:", err)
		return nil, fmt.Errorf("error fetching orders by status: %v", err)
	}

	lg.Log.Info("Successfully fetched orders by status ✅")
	return orders, nil
}

func GetALLOrdersList() ([]OrderService.FinalOrder, error) {
	lg.Log.Info("Fetching orders by status from the finalorder table!!")

	query := `SELECT order_id, product_id, count, order_date, order_status FROM finalorder`
	iter := Session.Query(query).Iter()

	var orders []OrderService.FinalOrder
	var order OrderService.FinalOrder

	for iter.Scan(&order.OrderId, &order.OrderDts.Product_Id, &order.OrderDts.Count, &order.OrderDte, &order.OrderStatus) {
		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		lg.Log.Error("Error fetching orders placed:", err)
		return nil, fmt.Errorf("error fetching orders placed: %v", err)
	}

	lg.Log.Info("Successfully fetched orders list ✅")
	return orders, nil
}
