package main

func main() {
	broker := NewBroker(":8080")
	broker.StartBroker()
}
