package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

var db *gorm.DB
var err error

// Product is representation of a product / Http Res & Req
type Product struct {
	ID    int             `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Code  string          `json:"code"`
	Name  string          `json:"name"`
	Price decimal.Decimal `json:"price" sql:"type:decimal(16,2)"`
}

// Http Response / Result
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// Request Http handler
func handleRequest() {
	log.Println("Start development server at http://127.0.0.1:8080")

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/api/add", createProduct).Methods("POST")
	myRouter.HandleFunc("/api/read", getProduct).Methods("GET")
	myRouter.HandleFunc("/api/read/{id}", getOneProduct).Methods("GET")
	myRouter.HandleFunc("/api/update/{id}", updateProduct).Methods("PUT")
	myRouter.HandleFunc("/api/delete/{id}", deleteProduct).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

// Server Connection Test
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Connected to server")
}

// Func ADD New Product
func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)

	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success CREATE Product"}
	result, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(result)

}

// Func READ Product
func getProduct(w http.ResponseWriter, r *http.Request) {
	products := []Product{}

	db.Find(&products)
	res := Result{Code: 200, Data: products, Message: "Success READ Product"}
	result, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(result)
}

// Func READ ONE Product
func getOneProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product

	db.First(&product, productID)
	res := Result{Code: 200, Data: product, Message: "Success READ SPECIFY Product"}
	result, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if product.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(result)

}

// Func UPDATE Product
func updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	payloads, _ := ioutil.ReadAll(r.Body)
	var productUpdate Product
	json.Unmarshal(payloads, &productUpdate)

	var product Product
	db.First(&product, productID)
	db.Model(&product).Updates(&productUpdate)

	res := Result{Code: 200, Data: product, Message: "Success UPDATE SPECIFY Product"}
	result, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Func DELETE Product
func deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	productID := vars["id"]

	var product Product
	db.First(&product, productID)
	db.Delete(&product)

	res := Result{Code: 200, Data: product, Message: "Success DELETE Product"}
	result, err := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

// Main Function
func main() {
	// Load .env file
	errENV := godotenv.Load()

	// check env is connect or not
	if errENV != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	// declare .env
	DB_ROOT := os.Getenv("DB_ROOT")
	DB_NAME := os.Getenv("DB_NAME")
	DB_PASS := os.Getenv("DB_PASS")

	// connect to db
	db, err = gorm.Open("mysql", DB_ROOT+":"+DB_PASS+"@/"+DB_NAME+"?charset=utf8&parseTime=true")

	// check db connection
	if err != nil {
		log.Println("Connection failed ", err)
	} else {
		log.Println("Connection Success")
	}

	// run automigrate with gorm
	db.AutoMigrate(&Product{})
	handleRequest()
}
