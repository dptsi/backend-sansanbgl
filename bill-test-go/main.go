package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/denisenkom/go-mssqldb"
)

type Item struct {
	BillDataId string         `json:"billdataid"`
	BillKey    string         `json:"billkey"`
	Semester   string         `json:"semester"`
	BillRef1   string         `json:"nama"`
	BillRef2   string         `json:"prodi"`
	BillAmount string         `json:"amount"`
	BillFlag   string         `json:"flag"`
	BillPaidAt sql.NullString `json:"paidat"`
}

var db *sql.DB

func init() {
	// Connect to the SQL Server database
	connString := ""
	var err error
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
}

func getItemsByBillKeyAndSemester(w http.ResponseWriter, r *http.Request) {
	// Get the billkey and semester parameters from the URL
	billKey := r.URL.Query().Get("billkey")
	semester := r.URL.Query().Get("semester")
	if billKey == "" || semester == "" {
		http.Error(w, "billkey and semester parameters are required", http.StatusBadRequest)
		return
	}

	// Query the database for the items with the specified billkey and semester
	rows, err := db.Query("SELECT BILL_DATA_ID as BillDataId, BILL_KEY1 as BillKey, SEMESTER_DIKTI as Semester, BILL_REF1 as BillRef1,  BILL_REF2 as BillRef2, BILL_AMOUNT as BillAmount, PAID_DATE as BillPaidAt, BILL_FLAG as BillFlag FROM bill_data WHERE BILL_KEY1 = @billKey AND SEMESTER_DIKTI = @semester", sql.Named("billkey", billKey), sql.Named("semester", semester))
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a slice of Item structs from the results
	var items []Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.BillDataId, &item.BillKey, &item.Semester, &item.BillRef1, &item.BillRef2, &item.BillAmount, &item.BillPaidAt, &item.BillFlag)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}
	err = rows.Err()
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode the items as JSON and return them
	json.NewEncoder(w).Encode(items)
}

func main() {
	// Register the endpoint that retrieves items by billkey and semester
	http.HandleFunc("/items", getItemsByBillKeyAndSemester)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
