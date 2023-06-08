package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Organization string `json:"organization"`
}

type App struct {
	DB *gorm.DB
	Conn *amqp.Connection
}

func (app *App)Signup(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user := models.User{ID: id, Email: input.Email, Password: string(hashedPassword)}
	organization := models.Organization{Name: input.Organization}
	userResult := app.DB.Create(&user)
	organizationResult := app.DB.Create(&organization)
	if userResult.Error != nil ||  organizationResult.Error != nil {
		response := map[string]interface{}{
			"message": "Internal server error",
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}
	app.SendEmail("Welcome "+input.Email)
	response := map[string]interface{}{
		"message": "User stored successfully",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (app *App) SendEmail(body string) {
	ch, err := app.Conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"mail", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
	  "",     // exchange
	  q.Name, // routing key
	  false,  // mandatory
	  false,  // immediate
	  amqp.Publishing {
		ContentType: "text/plain",
		Body:        []byte(body),
	  })
	if err != nil {
		panic(err)
	}
	log.Printf(" [x] Sent %s\n", body)
}