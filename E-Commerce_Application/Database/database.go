package Database

import (
	"context"
	"fmt"
	"log"
	ps "module/ProductService"
	"os"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	logrus "github.com/sirupsen/logrus"
)

const dbURL = "postgres://user:password@localhost:5432/orders?sslmode=disable"

var (
	DbPool *pgxpool.Pool
	logger = logrus.New()
)

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{}) // Logs in JSON format for ELK
	logger.SetOutput(os.Stdout)                  // Output to console (or a file)
	logger.SetLevel(logrus.InfoLevel)            // Set log level
}

func ConnectDB() (*pgxpool.Pool, error) {

	Dbp, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ConnectDB error": err,
		}).Fatal("Unable to connect to database")
		//log.Fatalf("Unable to connect to database: %s\n", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}
	//defer dbpool.Close()

	//fmt.Println("Connected to PostgreSQL successfully!")
	logger.Info("Connected to PostgreSQL successfully  ✅")
	fmt.Println("Postgresql DB ✅")

	//getTableSchema(dbpool, "order")
	//createTable(dbpool)
	DbPool = Dbp
	return Dbp, nil
}

func CreateUserTable(Dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS userdb (
    cust_id SERIAL PRIMARY KEY,
    cust_name VARCHAR(255) NOT NULL,
    cust_email VARCHAR(255) UNIQUE NOT NULL,
    cust_pnum VARCHAR(15) NOT NULL,
	cust_bal FLOAT DEFAULT 0.0
);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create userdb table: %v\n", err)
	}

	fmt.Println("User Table ✅")
}

func CreateProductTable(Dbpool *pgxpool.Pool) {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS product (
    product_id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100),
    subcategory VARCHAR(100),
    brand VARCHAR(100),
    price FLOAT NOT NULL,
    quantity INT NOT NULL
	);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create product table: %v\n", err)
	}

	fmt.Println("Product Table ✅")
}

// Load all the products in the Product Table

func LoadProdData(Dbpool *pgxpool.Pool) {
	//data := ps.Inventory

	insertQuery := `
		INSERT INTO product (product_id, name, description, category, subcategory, brand, price, quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (product_id) DO NOTHING;
	`

	// Loop over the inventory slice and insert each product
	for _, product := range ps.Inventory {
		_, err := Dbpool.Exec(context.Background(), insertQuery, product.ProductID, product.Name, product.Description, product.Category, product.Subcategory, product.Brand, product.Price, product.Quantity)
		if err != nil {
			log.Printf("Failed to insert product %s: %v\n", product.ProductID, err)
		} else {
			fmt.Printf("Successfully inserted product %s\n", product.ProductID)
		}
	}
	fmt.Println("Successfully Loaded")

}

func GetAllProdData(Dbpool *pgxpool.Pool) ([]ps.Product, error) {
	selectQuery := `SELECT * from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil

}

func GetProdByCateg(Dbpool *pgxpool.Pool, str string) ([]ps.Product, error) {
	selectQuery := `SELECT * from product where category = $1`

	rows, err := Dbpool.Query(context.Background(), selectQuery, str)
	if err != nil {
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetProdByBrand(Dbpool *pgxpool.Pool, str string) ([]ps.Product, error) {
	selectQuery := `SELECT * from product where brand = $1`

	rows, err := Dbpool.Query(context.Background(), selectQuery, str)
	if err != nil {
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetUniqueBrand(Dbpool *pgxpool.Pool) ([]string, error) {
	selectQuery := `SELECT DISTINCT(brand) from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		log.Printf("Failed to fetch product from GetUniqueBrand function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []string
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetUniqueCategory(Dbpool *pgxpool.Pool) ([]string, error) {
	selectQuery := `SELECT DISTINCT(category) from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		log.Printf("Failed to fetch product from GetUniqueBrand function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []string
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetProdPrice(Dbpool *pgxpool.Pool, id string) (float64, error) {

	inputQuery := `SELECT price from product where product_id=$1`
	var price float64
	err := Dbpool.QueryRow(context.Background(), inputQuery, id).Scan(&price)
	if err != nil {
		return price, fmt.Errorf("Error in fetching the price", err)

	}
	return price, nil
}
