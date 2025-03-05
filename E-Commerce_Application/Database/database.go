package Database

import (
	"context"
	"fmt"
	"log"
	"net"

	"module/Database/files"
	ps "module/ProductService"
	"module/UserService"
	lg "module/logger"
	"os"

	pgxpool "github.com/jackc/pgx/v5/pgxpool"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const dbURL = "postgres://user:password@localhost:5432/orders?sslmode=disable"

var (
	DbPool *pgxpool.Pool
	logger = logrus.New()
)

var SampleData *UserService.UserDetails

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{}) // Logs in JSON format for ELK
	logger.SetOutput(os.Stdout)                  // Output to console (or a file)
	logger.SetLevel(logrus.InfoLevel)            // Set log level
}

type server struct {
	files.UnimplementedDatabaseServiceServer
}

func ConnectDB() (*pgxpool.Pool, error) {
	lg.Log.Info("Entered Connect DB Function")
	Dbp, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"ConnectDB error": err,
		}).Fatal("Unable to connect to database")
		lg.Log.Error("Unable to connect to DB")
		//log.Fatalf("Unable to connect to database: %s\n", err)
		return nil, fmt.Errorf("Unable to connect to DB")
	}
	//defer dbpool.Close()

	//fmt.Println("Connected to PostgreSQL successfully!")
	lg.Log.Info("Postgresql DB ✅")
	logger.Info("Connected to PostgreSQL successfully  ✅")
	fmt.Println("Postgresql DB ✅")

	//getTableSchema(dbpool, "order")
	//createTable(dbpool)
	DbPool = Dbp
	return Dbp, nil
}

func StartServer() {
	// Create a listener on port 50051
	lg.Log.Info("Create a listener on port 50051")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		lg.Log.Error("Failed to listen: %v", err)
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	files.RegisterDatabaseServiceServer(grpcServer, &server{})
	lg.Log.Info("gRPC server started on port 50051")
	fmt.Println("gRPC server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		lg.Log.Error("Failed to serve: %v", err)
		log.Fatalf("Failed to serve: %v", err)
	}

}

func (s *server) CheckLoginCredentials(ctx context.Context, req *files.LoginRequest) (*files.Empty, error) {
	log.Printf("Checking login credentials for username: %s", req.Username)
	lg.Log.Info("Checking login credentials for username: %s", req.Username)
	err := CheckLoginCredentials(req.Username)

	return &files.Empty{}, err
}

func (s *server) LoadLoginCred(ctx context.Context, req *files.LoginCredRequest) (*files.Empty, error) {
	lg.Log.Info("Loading login credentials for username: %s", req.Username)
	log.Printf("Loading login credentials for username: %s", req.Username)
	var lc UserService.LoginCred
	lc.UserName = req.Username
	lc.Password = req.Password
	err := LoadLoginCred(lc)

	return &files.Empty{}, err
}

func (s *server) GetAllLoginCred(ctx context.Context, req *files.Empty) (*files.LoginCredList, error) {
	lg.Log.Info("Fetching all login credentials")
	log.Println("Fetching all login credentials")
	k, err := GetAllLoginCred()

	var loginList []*files.LoginCred
	lc := &files.LoginCred{}
	for _, i := range k {

		lc.Username = i.UserName
		lc.Password = i.Password
		loginList = append(loginList, lc) // Convert to pointer
	}

	return &files.LoginCredList{Credentials: loginList}, err
}

func (s *server) GetAllUserData(ctx context.Context, req *files.Empty) (*files.UserDetailsList, error) {
	lg.Log.Info("Fetching all user details")
	log.Println("Fetching all user details")
	var ulist []*files.UserDetails
	u := files.UserDetails{}
	k, err := GetALLUserData()
	for _, k := range k {
		u.Cust_Id = int32(k.Cust_Id)
		u.Cust_Bal = float32(k.Cust_Bal)
		u.Cust_Email = k.Cust_Email
		u.Cust_Name = k.Cust_Name
		u.Cust_PNum = k.Cust_PNum
		u.UserName = k.UserName
		ulist = append(ulist, &u)
	}
	return &files.UserDetailsList{Users: ulist}, err
}

func (s *server) GetUserByUserDetails(ctx context.Context, req *files.UserRequest) (*files.UserDetails, error) {
	log.Printf("Fetching user details for username: %s", req.Username)
	lg.Log.Info("Fetching user details for username: %s", req.Username)
	k, err := GetUserByUserDeatils(req.Username)
	fmt.Println("k is:", k)
	if err == nil {
		ud := &files.UserDetails{
			Cust_Bal:   float32(k.Cust_Bal),
			Cust_Email: k.Cust_Email,
			Cust_Id:    int32(k.Cust_Id),
			Cust_Name:  k.Cust_Name,
			Cust_PNum:  k.Cust_PNum,
			UserName:   k.UserName,
		}
		lg.Log.Info("Userdetails is : %v", ud)
		log.Printf("Userdetails is : %v", ud)
		return ud, err
	}

	return nil, err // Return nil if user is not found
}

func CreateUserTable(Dbpool *pgxpool.Pool) {
	lg.Log.Info("Creating user table !")
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS userdb (
    cust_id INT PRIMARY KEY,
    cust_name VARCHAR(255) NOT NULL,
    cust_email VARCHAR(255) UNIQUE NOT NULL,
    cust_pnum VARCHAR(15) NOT NULL,
	cust_bal FLOAT DEFAULT 0.0,
	cust_uname VARCHAR(255) NOT NULL
);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createTableQuery)
	if err != nil {
		lg.Log.Error("Failed to create userdb table: %v\n", err)
		log.Fatalf("Failed to create userdb table: %v\n", err)
	}
	lg.Log.Info("User Table Created Successfully !!")
	fmt.Println("User Table ✅")
}

func CreateLoginTable(Dbpool *pgxpool.Pool) {
	lg.Log.Info("Creating login table !")
	createLoginTableQuery := `
	CREATE TABLE IF NOT EXISTS logincredentials (
	   	
	    username VARCHAR(255) PRIMARY KEY,
	    password TEXT NOT NULL
	);
	`
	var err error
	_, err = Dbpool.Exec(context.Background(), createLoginTableQuery)
	if err != nil {
		lg.Log.Error("Failed to create logincredentials table: %v\n", err)
		log.Fatalf("Failed to create logincredentials table: %v\n", err)
	}
	lg.Log.Info("LoginCRedential ✅")
	fmt.Println("LoginCRedential ✅")
}

func CheckLoginCredentials(uname string) error {
	lg.Log.Info("Checking login credentials !!!")
	checkQuery := `SELECT COUNT(*) FROM logincredentials WHERE username = $1;`
	var count int
	if DbPool == nil {
		fmt.Println("DbPool is nil")
		lg.Log.Error("DbPool is nil")
		return fmt.Errorf("failed to check username existence")
	}
	err := DbPool.QueryRow(context.Background(), checkQuery, uname).Scan(&count)
	if err != nil {
		lg.Log.Error("failed to check username existence: %v", err)
		return fmt.Errorf("failed to check username existence: %v", err)
	}

	if count > 0 {
		lg.Log.Error("exists")
		return fmt.Errorf("exists")
	} else {
		return nil
	}

}

func LoadLoginCred(lc UserService.LoginCred) error {
	lg.Log.Info("Entered login credentials table")
	insertQuery := `INSERT INTO logincredentials (username, password) VALUES ($1, $2);`
	var err error
	if DbPool == nil {
		lg.Log.Error("DbPool is nil")
		fmt.Println("DbPool is nil")
		return fmt.Errorf("failed to check username existence")
	}
	_, err = DbPool.Exec(context.Background(), insertQuery, lc.UserName, lc.Password)
	if err != nil {
		lg.Log.Error("failed to insert user: %v", err)
		return fmt.Errorf("failed to insert user: %v", err)
	}
	lg.Log.Info("User Login inserted successfully ✅")
	fmt.Println("User Login inserted successfully ✅")
	return nil
}

func GetAllLoginCred() ([]UserService.LoginCred, error) {
	lg.Log.Info("Entered GetAllLogin cred function")
	selectQuery := `SELECT * from logincredentials`

	rows, err := DbPool.Query(context.Background(), selectQuery)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetALLProdData function %v\n", err)
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []UserService.LoginCred
	for rows.Next() {
		var p UserService.LoginCred
		err := rows.Scan(&p.UserName, &p.Password)
		if err != nil {
			lg.Log.Error("Error scanning login cred: %v", err)
			log.Printf("Error scanning login cred: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func CheckValidProduct(Dbpool *pgxpool.Pool, pid string) error {
	lg.Log.Info("Entered check valid product function !!")
	query := `select product_id from product where product_id= $1`

	var err error
	_, err = Dbpool.Exec(context.Background(), query, pid)
	if err != nil {
		lg.Log.Error("Failed to extract product from the table: %v\n", err)
		log.Fatalf("Failed to extract product from the table: %v\n", err)
		return err
	}
	return nil
}

func CreateProductTable(Dbpool *pgxpool.Pool) {
	lg.Log.Info("Entered create product table function")
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
		lg.Log.Error("Failed to create product table: %v\n", err)
		log.Fatalf("Failed to create product table: %v\n", err)
	}
	lg.Log.Info("Product Table ✅")
	fmt.Println("Product Table ✅")
}

// Load User data in the user table

func LoadUserData(user UserService.UserDetails) error {
	lg.Log.Info("Loading User data !!!")
	insertQuery := `
		INSERT INTO userdb (cust_id, cust_name, cust_email, cust_pnum, cust_bal, cust_uname)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := DbPool.Exec(context.Background(), insertQuery, user.Cust_Id, user.Cust_Name, user.Cust_Email, user.Cust_PNum, user.Cust_Bal, user.UserName)
	if err != nil {
		lg.Log.Error("Failed to insert user %s: %v\n", user.Cust_Id, err)
		log.Printf("Failed to insert user %s: %v\n", user.Cust_Id, err)
		return err
	} else {
		lg.Log.Info("Successfully inserted product %s \n", user.Cust_Name)
		fmt.Printf("Successfully inserted product %s \n", user.Cust_Name)
		return nil
	}

}

func GetUserBalance(id int) (float64, error) {
	lg.Log.Info("Entered get user balance function !!")
	selectQuery := `SELECT cust_bal FROM userdb WHERE cust_id = $1`

	var custBal float64
	err := DbPool.QueryRow(context.Background(), selectQuery, id).Scan(&custBal)
	if err != nil {
		lg.Log.Error("Failed to fetch balance for customer %d: %v\n", id, err)
		log.Printf("Failed to fetch balance for customer %d: %v\n", id, err)
		return 0, err
	}

	return custBal, nil
}

func GetALLUserData() ([]UserService.UserDetails, error) {
	lg.Log.Info("Entered getall user data function !!")
	selectQuery := `SELECT * from userdb`

	rows, err := DbPool.Query(context.Background(), selectQuery)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetALLProdData function %v\n", err)
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []UserService.UserDetails
	for rows.Next() {
		var p UserService.UserDetails
		err := rows.Scan(&p.Cust_Id, &p.Cust_Name, &p.Cust_Email, &p.Cust_PNum, &p.Cust_Bal, &p.UserName)
		if err != nil {
			lg.Log.Error("Error scanning user: %v", err)
			log.Printf("Error scanning user: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil

}

// Load all the products in the Product Table

func LoadProdData(Dbpool *pgxpool.Pool) {
	//data := ps.Inventory
	lg.Log.Info("Loading product data !!!")
	insertQuery := `
		INSERT INTO product (product_id, name, description, category, subcategory, brand, price, quantity)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (product_id) DO NOTHING;
	`

	// Loop over the inventory slice and insert each product
	for _, product := range ps.Inventory {
		_, err := Dbpool.Exec(context.Background(), insertQuery, product.ProductID, product.Name, product.Description, product.Category, product.Subcategory, product.Brand, product.Price, product.Quantity)
		if err != nil {
			lg.Log.Error("Failed to insert product %s: %v\n", product.ProductID, err)
			log.Printf("Failed to insert product %s: %v\n", product.ProductID, err)
		} else {
			lg.Log.Info("Successfully inserted product %s\n", product.ProductID)
			fmt.Printf("Successfully inserted product %s\n", product.ProductID)
		}
	}
	lg.Log.Info("Successfully Loaded !!!")
	fmt.Println("Successfully Loaded")

}

func UpdateUserBalance(Dbpool *pgxpool.Pool, str int, newbal float64) error {
	lg.Log.Info("Updating the user balance")
	updateQuery := `UPDATE userdb SET cust_bal=$1 WHERE cust_id=$2`
	_, err := Dbpool.Exec(context.Background(), updateQuery, newbal, str)
	if err != nil {
		lg.Log.Error("Failed to update balance for user %s: %v\n", str, err)
		log.Printf("Failed to update balance for user %s: %v\n", str, err)
		return err
	}
	return nil
}

func GetUserByUserDeatils(str string) (UserService.UserDetails, error) {
	lg.Log.Info("Get user details by username !!")
	selectQuery := `SELECT cust_id, cust_name, cust_email, cust_pnum, cust_bal, cust_uname FROM userdb WHERE cust_uname=$1`
	var p UserService.UserDetails

	// Execute the query
	rows, err := DbPool.Query(context.Background(), selectQuery, str)
	if err != nil {
		lg.Log.Error("Failed to fetch user: %v\n", err)
		log.Printf("Failed to fetch user: %v\n", err)
		return p, err
	}
	defer rows.Close()

	// Check if any row exists
	if !rows.Next() {
		lg.Log.Error("No user found for username: %s", str)
		log.Printf("No user found for username: %s", str)
		return p, fmt.Errorf("no user found for username: %s", str)
	}

	// Scan the row into the struct
	err = rows.Scan(&p.Cust_Id, &p.Cust_Name, &p.Cust_Email, &p.Cust_PNum, &p.Cust_Bal, &p.UserName)
	if err != nil {
		lg.Log.Error("Error scanning user details: %v", err)
		log.Printf("Error scanning user details: %v", err)
		return p, err
	}

	// Assign to global variable (if necessary)
	SampleData = &p
	fmt.Println(*SampleData)

	return p, nil
}

func GetAllProdData(Dbpool *pgxpool.Pool) ([]ps.Product, error) {
	lg.Log.Info("Get all products data function is called !!")
	selectQuery := `SELECT * from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetALLProdData function %v\n", err)
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			lg.Log.Error("Error scanning product: %v", err)
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil

}

func GetProdByCateg(Dbpool *pgxpool.Pool, str string) ([]ps.Product, error) {
	lg.Log.Info("Entered get product by category functions")
	selectQuery := `SELECT * from product where category = $1`

	rows, err := Dbpool.Query(context.Background(), selectQuery, str)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetALLProdData function %v\n", err)
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			lg.Log.Error("Error scanning product: %v", err)
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetProdByBrand(Dbpool *pgxpool.Pool, str string) ([]ps.Product, error) {
	lg.Log.Info("Entered the get product by brand function !!")
	selectQuery := `SELECT * from product where brand = $1`

	rows, err := Dbpool.Query(context.Background(), selectQuery, str)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetALLProdData function %v\n", err)
		log.Printf("Failed to fetch product from GetALLProdData function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []ps.Product
	for rows.Next() {
		var p ps.Product
		err := rows.Scan(&p.ProductID, &p.Name, &p.Description, &p.Category, &p.Subcategory, &p.Brand, &p.Price, &p.Quantity)
		if err != nil {
			lg.Log.Error("Error scanning product: %v", err)
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetUniqueBrand(Dbpool *pgxpool.Pool) ([]string, error) {
	lg.Log.Info("Entered unique brand function")
	selectQuery := `SELECT DISTINCT(brand) from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetUniqueBrand function %v\n", err)
		log.Printf("Failed to fetch product from GetUniqueBrand function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []string
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			lg.Log.Error("Error scanning product: %v", err)
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetUniqueCategory(Dbpool *pgxpool.Pool) ([]string, error) {
	lg.Log.Info("enetered unique category function")
	selectQuery := `SELECT DISTINCT(category) from product`

	rows, err := Dbpool.Query(context.Background(), selectQuery)
	if err != nil {
		lg.Log.Error("Failed to fetch product from GetUniqueBrand function %v\n", err)
		log.Printf("Failed to fetch product from GetUniqueBrand function %v\n", err)
		return nil, err
	}

	defer rows.Close()

	var pad []string
	for rows.Next() {
		var p string
		err := rows.Scan(&p)
		if err != nil {
			lg.Log.Error("Error scanning product: %v", err)
			log.Printf("Error scanning product: %v", err)
			//c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning product"})
			return nil, err
		}
		pad = append(pad, p)
	}

	return pad, nil
}

func GetProdPrice(Dbpool *pgxpool.Pool, id string) (float64, error) {
	lg.Log.Info("Entered the get product price function !!!")
	inputQuery := `SELECT price from product where product_id=$1`
	var price float64
	err := Dbpool.QueryRow(context.Background(), inputQuery, id).Scan(&price)
	if err != nil {
		lg.Log.Error("Error in fetching the price", err)
		return price, fmt.Errorf("Error in fetching the price", err)

	}
	return price, nil
}

func UpdateProdTbPostCompletion(pid string, cnt int) (bool, error) {

	lg.Log.Info("Updating product table after completion...")

	updateQuery := `UPDATE product SET quantity = $1 WHERE product_id = $2`

	_, err := DbPool.Exec(context.Background(), updateQuery, cnt, pid)
	if err != nil {
		lg.Log.Error("Error updating product quantity: ", err)
		return false, fmt.Errorf("error updating product quantity: %v", err)
	}

	lg.Log.Info("Product quantity updated successfully for product_id:", pid)
	return true, nil

}
