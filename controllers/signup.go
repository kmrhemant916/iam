package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/authorization"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Organization string `json:"organization"`
}

type MailPayload struct {
	To string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

type App struct {
	DB *gorm.DB
	Conn *amqp.Connection
	JWTKey []byte
}

const (
	DefaultRootGroup = "Administrator"
	DefaultUserGroup = "Reader"
)

func (app *App)Signup(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload
	err := json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user := models.User{ID: id, Email: requestPayload.Email, Password: string(hashedPassword), IsRoot: "true"}
	organization := models.Organization{Name: requestPayload.Organization}
	userResult := app.DB.Create(&user)
	if userResult.Error != nil {
		response := map[string]interface{}{
			"message": "Internal server error",
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	} else {
		organizationResult := app.DB.Create(&organization)
		if organizationResult.Error != nil {
			response := map[string]interface{}{
				"message": "Internal server error",
			}
			helpers.SendResponse(w,response, http.StatusInternalServerError)
		}
	}
	mailPayload := MailPayload{
		To: "hemank",
		From: "ddd",
		Body: "sss",
	}
	body, err := json.Marshal(mailPayload)
	if err != nil {
		log.Fatalf("Failed to marshal message: %s", err)
	}
	app.SendEmail(body)
	rbac := &authorization.Rbac{
		DB: app.DB,
	}
	var groups []string
	groups = append(groups, DefaultRootGroup)
	rbac.AssignGroups(id, groups)
	response := map[string]interface{}{
		"message": "User stored successfully",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (app *App) SendEmail(body []byte) {
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
		ContentType: "application/json",
		Body:        []byte(body),
	  })
	if err != nil {
		panic(err)
	}
	log.Printf(" [x] Sent %s\n", body)
}