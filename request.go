package main

type Request struct {
	Method string                 // get, post, put, delete
	Model  string                 // person, invoice, order
	Data   map[string]interface{} // data for PUT, POST
}

type Response struct {
	Status  string
	Message string
	Data    map[string]interface{} // data for PUT, POST
}
