package main

import (
	"context"
	"encoding/json"
	pizza "workshop/exercises/non-retryable-error-types/practice"
	"flag"
	"fmt"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	var cardNumber = flag.String("creditcard", "1234123412341234", "Provide a credit card number")
	flag.Parse()
	order := *createPizzaOrder(cardNumber)

	workflowID := fmt.Sprintf("pizza-workflow-order-%s", order.OrderNumber)

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: pizza.TaskQueueName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, pizza.PizzaWorkflow, order)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	var result pizza.OrderConfirmation
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable to deliver pizza", err)
	}

	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalln("Unable to format order confirmation as JSON", err)
	}
	log.Printf("Workflow result: %s\n", string(data))
}

func createPizzaOrder(cardNumber *string) *pizza.PizzaOrder {
	customer := pizza.Customer{
		CustomerID: 12983,
		Name:       "María García",
		Email:      "maria1985@example.com",
		Phone:      "415-555-7418",
	}

	address := pizza.Address{
		Line1:      "1 Main St",
		Line2:      "Apartment 9C",
		City:       "San Francisco",
		State:      "CA",
		PostalCode: "94103",
		CardNumber: *cardNumber,
	}

	p1 := pizza.Pizza{
		Description: "Large, with mushrooms and onions",
		Price:       1500,
	}

	p2 := pizza.Pizza{
		Description: "Small, with pepperoni",
		Price:       1200,
	}

	p3 := pizza.Pizza{
		Description: "Medium, with extra cheese",
		Price:       1300,
	}

	items := []pizza.Pizza{p1, p2, p3}

	order := pizza.PizzaOrder{
		OrderNumber: "Z1238",
		Customer:    customer,
		Items:       items,
		Address:     address,
		IsDelivery:  true,
	}

	return &order
}
