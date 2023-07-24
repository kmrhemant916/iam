package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/authorization"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/repositories"
	"github.com/kmrhemant916/iam/service"
	"github.com/kmrhemant916/iam/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/crypto/bcrypt"
)

type SignupPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=32"`
	Organization string `json:"organization" validate:"required"`
}

type MailPayload struct {
	To string `json:"to"`
	From string `json:"from"`
	Body string `json:"body"`
}

const (
	DefaultRootGroup = "Administrator"
	DefaultUserGroup = "Reader"
)

func (app *App)Signup(w http.ResponseWriter, r *http.Request) {
	var signupPayload SignupPayload
	err := json.NewDecoder(r.Body).Decode(&signupPayload)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	errorsList, err := utils.ValidateJSON(signupPayload)
	if err != nil {
		for _, e := range errorsList {
			switch {
				case e.FailedField == "SignupPayload.Password" && (e.Tag == "min" || e.Tag == "max"):
					response := map[string]interface{}{
						"message": "password should be in between 8 and 32 characters",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				case e.FailedField == "SignupPayload.Email" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "email field is required",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				case e.FailedField == "SignupPayload.Organization" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "organization field is required",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				default:
					helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
					return
			}
		}
		helpers.SendResponse(w, global.InvalidRequestPayloadMessage, http.StatusBadRequest)
		return
	}
	userId := uuid.New()
	organizationId := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupPayload.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	organization := models.Organization{OrganizationID: organizationId, Name: signupPayload.Organization}
	user := models.User{UserID: userId, Email: signupPayload.Email, Password: string(hashedPassword), IsRoot: true, OrganizationID: organizationId}
	signupRepository := repositories.NewSignupRepository(app.DB)
	signupService := service.NewSignupService(signupRepository)
	err = signupService.CreateRootAccount(utils.UserToEntity(&user), utils.OrganizationToEntity(&organization))
    if err != nil {
        switch {
			case errors.Is(err, global.ErrOrgExists):
				response := map[string]interface{}{
					"message": "org already exist",
				}
				helpers.SendResponse(w,response, http.StatusForbidden)
				return
			case errors.Is(err, global.ErrUserExists):
				response := map[string]interface{}{
					"message": "user already exist",
				}
				helpers.SendResponse(w,response, http.StatusForbidden)
				return
			default:
				helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			}
        return
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
	rbac.AssignGroups(userId, groups)
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