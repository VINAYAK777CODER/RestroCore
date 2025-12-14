package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/VINAYAK777CODER/RestroCore/database"
	"github.com/VINAYAK777CODER/RestroCore/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var invoiceCollection *mongo.Collection=database.OpenCollection(database.Client,"invoice")

type InvoiceView struct {
	Invoice_id       string    `json:"invoice_id"`
	Order_id         string    `json:"order_id"`
	Payment_method   *string   `json:"payment_method,omitempty"`
	Payment_status   *string   `json:"payment_status,omitempty"`
	Payment_due_date time.Time `json:"payment_due_date,omitempty"`
	Created_at       time.Time `json:"created_at"`
}


func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {

		// create context with timeout for DB operations
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// fetch all invoices from database
		cursor, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while fetching invoices"})
			return
		}
		// close cursor after function execution
		defer cursor.Close(ctx)

		// slice to store all invoice documents from DB
		var invoices []models.Invoice

		// decode DB records into slice
		if err := cursor.All(ctx, &invoices); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error decoding invoice data"})
			return
		}

		// slice to store API response data
		var invoiceViews []InvoiceView

		// convert DB model to response view model
		for _, invoice := range invoices {
			invoiceViews = append(invoiceViews, InvoiceView{
				Invoice_id:       invoice.Invoice_id,
				Order_id:         invoice.Order_id,
				Payment_method:   invoice.Payment_method,
				Payment_status:   invoice.Payment_status,
				Payment_due_date: invoice.Payment_due_date,
				Created_at:       invoice.Created_at,
			})
		}

		// send final response
		c.JSON(http.StatusOK, invoiceViews)
	}
}


func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		invoiceId := c.Param("invoice_id")
		if invoiceId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invoice_id is required"})
			return
		}

		var invoice models.Invoice

		err := invoiceCollection.FindOne(
			ctx,
			bson.M{"invoice_id": invoiceId},
		).Decode(&invoice)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
			return
		}

		// Map DB model → View model
		invoiceView := InvoiceView{
			Invoice_id:       invoice.Invoice_id,
			Order_id:         invoice.Order_id,
			Payment_method:   invoice.Payment_method,
			Payment_status:   invoice.Payment_status,
			Payment_due_date: invoice.Payment_due_date,
			Created_at:       invoice.Created_at,
		}

		c.JSON(http.StatusOK, invoiceView)
	}
}


func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

		// create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var invoice models.Invoice

		// bind request body (ONLY order_id & payment_method allowed)
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}

		// validate required fields (order_id)
		if err := validate.Struct(invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// generate MongoDB ID & invoice_id
		invoice.ID = primitive.NewObjectID()
		invoice.Invoice_id = invoice.ID.Hex()

		// RESTAURANT LOGIC (backend controlled)
		status := "PENDING"
		invoice.Payment_status = &status

		// payment due in 30 minutes
		invoice.Payment_due_date = time.Now().Add(30 * time.Minute)

		// timestamps
		invoice.Created_at = time.Now()
		invoice.Updated_at = time.Now()

		// insert invoice into DB
		_, err := invoiceCollection.InsertOne(ctx, invoice)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create invoice"})
			return
		}

		// prepare response (view model)
		invoiceView := InvoiceView{
			Invoice_id:       invoice.Invoice_id,
			Order_id:         invoice.Order_id,
			Payment_method:   invoice.Payment_method,
			Payment_status:   invoice.Payment_status,
			Payment_due_date: invoice.Payment_due_date,
			Created_at:       invoice.Created_at,
		}

		// send response
		c.JSON(http.StatusCreated, invoiceView)
	}
}


func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		// get invoice_id from URL
		invoiceId := c.Param("invoice_id")
		if invoiceId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invoice_id is required"})
			return
		}

		// struct for partial updates
		type UpdateInvoiceInput struct {
			Payment_method   *string    `json:"payment_method,omitempty"`
			Payment_status   *string    `json:"payment_status,omitempty"`
			Payment_due_date *time.Time `json:"payment_due_date,omitempty"`
		}

		var input UpdateInvoiceInput

		// bind JSON body
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
			return
		}

		// build update fields dynamically
		updateFields := bson.D{}

		if input.Payment_method != nil {
			updateFields = append(updateFields, bson.E{
				Key:   "payment_method",
				Value: input.Payment_method,
			})
		}

		if input.Payment_status != nil {
			updateFields = append(updateFields, bson.E{
				Key:   "payment_status",
				Value: input.Payment_status,
			})
		}

		if input.Payment_due_date != nil {
			updateFields = append(updateFields, bson.E{
				Key:   "payment_due_date",
				Value: input.Payment_due_date,
			})
		}

		// if no fields provided
		if len(updateFields) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no fields provided to update"})
			return
		}

		// always update updated_at
		updateFields = append(updateFields, bson.E{
			Key:   "updated_at",
			Value: time.Now(),
		})

		// update invoice
		result, err := invoiceCollection.UpdateOne(
			ctx,
			bson.M{"invoice_id": invoiceId},
			bson.M{"$set": updateFields},
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update invoice"})
			return
		}

		// invoice not found
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "invoice not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "invoice updated successfully"})
	}
}
