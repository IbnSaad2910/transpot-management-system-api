package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type docket struct {
	OrderNo       string  `json:"orderNo"`
	Customer      string  `json:"customer"`
	PickUpPoint   string  `json:"pickUpPoint"`
	DeliveryPoint string  `json:"deliveryPoint"`
	Quantity      int     `json:"quantity"`
	Volume        float32 `json:"volume"`
	Status        string  `json:"status"`
	TruckNo       string  `json:"truckNo"`
	LogsheetNo    string  `json:"logsheetNo"`
}

type logsheet struct {
	LogsheetNo   string   `json:"logsheetNo"`
	DocketsList  []string `json:"docketsList"`
	DocketsSlice []docket `json:"docketsSlice"`
	TruckNo      string   `json:"truckNo"`
}

// Docket slice
var dockets = []docket{}
var orderNo string = "TDN0000"

// Logsheet slice
var logsheets = []logsheet{}
var logsheetNo string = "DT0000"

func main() {
	router := gin.Default()
	router.POST("/docket", createDocket)
	router.GET("/docket/:orderNo", getDocketByOrderNo)
	router.GET("/docket", getDockets)
	router.POST("/logsheet", createLogsheet)
	router.GET("/logsheet/:logsheetNo", getLogsheetByLogsheetNo)

	router.Run("localhost:8080")
}

// Fetch all dockets
func getDockets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, dockets)
}

// Generate unique OrderNo
func generateOrderNo() string {
	var temp int
	if _, err := fmt.Sscanf(orderNo, "TDN%4d", &temp); err != nil {
		return orderNo
	}
	temp++
	orderNo = fmt.Sprintf("TDN%04d", temp)

	return orderNo
}

// Create new docket
func createDocket(c *gin.Context) {
	var newDocket docket

	newOrderNo := generateOrderNo()

	if err := c.BindJSON(&newDocket); err != nil {
		return
	}

	updateDocket(&newDocket, "", "", newOrderNo)
	dockets = append(dockets, newDocket)
	c.IndentedJSON(http.StatusCreated, newDocket)
}

// Fetch a docket based on order number
func getDocketByOrderNo(c *gin.Context) {
	orderNo := c.Param("orderNo")

	for _, docket := range dockets {
		if docket.OrderNo == orderNo {
			c.IndentedJSON(http.StatusOK, docket)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "docket not found"})
}

// Update docket
func updateDocket(d *docket, truckNo, logsheetNo, orderNo string) {
	d.TruckNo = truckNo
	d.LogsheetNo = logsheetNo
	d.OrderNo = orderNo
	d.Status = "Created"
}

// Generate unique LogsheetNo
func generateLogsheetNo() string {
	var temp int
	if _, err := fmt.Sscanf(logsheetNo, "DT%4d", &temp); err != nil {
		return logsheetNo
	}
	temp++
	logsheetNo = fmt.Sprintf("DT%04d", temp)

	return logsheetNo
}

// Create a new logsheet
func createLogsheet(c *gin.Context) {
	var newLogsheet logsheet

	newLogsheetNo := generateLogsheetNo()

	if err := c.BindJSON(&newLogsheet); err != nil {
		return
	}

	updateLogsheet(&newLogsheet, newLogsheetNo)

	for i := 0; i < len(dockets); i++ {
		for _, orderNo := range newLogsheet.DocketsList {
			if dockets[i].OrderNo == orderNo {
				// updateDocket(&docket, newLogsheet.TruckNo, newLogsheet.LogsheetNo, docket.OrderNo)
				dockets[i].LogsheetNo = newLogsheet.LogsheetNo
				dockets[i].TruckNo = newLogsheet.TruckNo
				newLogsheet.DocketsSlice = append(newLogsheet.DocketsSlice, dockets[i])
			}
		}
	}

	logsheets = append(logsheets, newLogsheet)
	c.IndentedJSON(http.StatusCreated, newLogsheet.DocketsSlice)
}

// Fetch a logsheet based on logsheet number
func getLogsheetByLogsheetNo(c *gin.Context) {
	logsheetNo := c.Param("logsheetNo")

	for _, logsheet := range logsheets {
		if logsheet.LogsheetNo == logsheetNo {
			c.IndentedJSON(http.StatusOK, logsheet.DocketsSlice)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "logsheet not found"})
}

// Update logsheet
func updateLogsheet(l *logsheet, logsheetNo string) {
	l.LogsheetNo = logsheetNo
}
