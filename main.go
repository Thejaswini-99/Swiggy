package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type order struct {
	Id       int    `json:"id"`
	Item     string `json:"item"`
	Customer string `json:"customer"`
	Status   string `json:"status"`
}

var ordersList = []order{
	{1, "Pizza", "Thejaswini", "Pending"},
	{2, "Burger", "Khushi", "Completed"},
	{3, "Chicken", "Veera", "Completed"},
}

func main() {

	http.HandleFunc("/getOrders", getOrder)
	http.HandleFunc("/getOrdersbyId", getOrderbyId)
	http.HandleFunc("/createOrder", createOrder)
	http.HandleFunc("/updateOrder", updateOrder)
	http.HandleFunc("/deleteOrder", deleteOrder)
	http.ListenAndServe(":8080", nil)

}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		fmt.Fprintf(w, "Only DELETE method is allowed")
		return
	}
	idstr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Id is required to update the order")
		return
	}

	for i, val := range ordersList {
		if val.Id == id {
			ordersList = append(ordersList[:i], ordersList[i+1:]...)
			fmt.Fprint(w, "Order Deleted!")
			return
		}
	}
	fmt.Fprint(w, "Order not Deleted as the given id was not found in the list!")

}

func updateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		fmt.Fprintf(w, "Only PUT method is allowed")
		return
	}
	var updateOrder order
	idstr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Id is required to update the order")
		return
	}
	json.NewDecoder(r.Body).Decode(&updateOrder)

	for i, val := range ordersList {
		if val.Id == id {
			ordersList[i] = updateOrder
			fmt.Fprint(w, "Order Updated!")
			return
		}
	}
	fmt.Fprint(w, "Order not Updated!")

}

func createOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only POST method is allowed")
		return
	}
	var newOrder order

	err := json.NewDecoder(r.Body).Decode(&newOrder)
	if err != nil {
		fmt.Fprintf(w, "Request Body not Found")
		return
	}
	ordersList = append(ordersList, newOrder)
	fmt.Fprintf(w, "Hey! Updated the Order")

}

func getOrderbyId(w http.ResponseWriter, r *http.Request) {
	idstr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idstr)
	if err != nil {
		fmt.Fprintf(w, "Id should be needed to fetch the order")
		return
	}

	for _, val := range ordersList {
		if id == val.Id {
			json.NewEncoder(w).Encode(val)
			return
		}

	}

	fmt.Fprintf(w, "Order Not FOund")

}

func getOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Only Get Method is allowed")
		return
	}

	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(ordersList)
}
