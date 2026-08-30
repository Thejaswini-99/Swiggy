package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type order struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Item   string `json:"item"`
	Status string `json:"status"`
}

func main() {
	initDB()
	http.HandleFunc("/getOrders", getOrders)
	http.HandleFunc("/getOrderbyID", getOrdersbyID)
	http.HandleFunc("/createOrder", createOrder)
	http.HandleFunc("/updateOrder", updateOrder)
	http.HandleFunc("/deleteOrder", deleteOrder)
	http.ListenAndServe(":8080", nil)

}

func getOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please use the GET Method")
		return
	}
	row, err := db.Query("Select Id,Name,Item,Status from orders")
	if err != nil {
		fmt.Fprintf(w, "Error fetching orders")
		return
	}
	defer row.Close()

	var orderList []order
	for row.Next() {
		var o order
		row.Scan(&o.Id, &o.Name, &o.Item, &o.Status)
		orderList = append(orderList, o)
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(orderList)

}
func getOrdersbyID(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("Id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Please enter the id to fetch")
		return
	}
	var o order
	err = db.QueryRow("Select Id,Name,Item,Status from orders where id = ? ", id).Scan(&o.Id, &o.Name, &o.Item, &o.Status)
	if err != nil {
		fmt.Fprintf(w, "Error fetching orders")
		return
	}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(o)

	fmt.Fprintf(w, "Order Id Fetched")

}

func createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please use the POST Method")
		return
	}
	var newOrder order

	w.Header().Set("content-type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&newOrder)

	if err != nil {
		fmt.Fprintf(w, "Not able to decode the value")
	}
	_, err = db.Exec("INSERT INTO orders (name, item, status) VALUES (?, ?, ?)",
		newOrder.Name, newOrder.Item, newOrder.Status)

	//orders = append(orders, newOrder)
	fmt.Fprintf(w, "New Order Created")

}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please use the POST Method")
		return
	}
	idstr := r.URL.Query().Get("Id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Please enter the proper id to update")
	}

	var newOrder order

	err = json.NewDecoder(r.Body).Decode(&newOrder)

	_, err = db.Exec("update orders set name= ?,item=?,status=? where id=?", newOrder.Name, newOrder.Item, newOrder.Status, id)

	fmt.Fprintf(w, "Order got updated")

}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Please use the DELETE Method")
		return
	}
	idstr := r.URL.Query().Get("Id")
	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Please enter the proper id to update")
	}
	_, err = db.Exec("delete from orders where id=?", id)

	fmt.Fprintf(w, "Order got Deleted")

}
