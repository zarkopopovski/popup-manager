package controllers

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"crypto/sha1"
	"encoding/json"
	"time"

	"math/rand"

	"crypto/tls"

	"gopkg.in/gomail.v2"

	//"github.com/teris-io/shortid"

	// "github.com/google/uuid"
	"github.com/zarkopopovski/popup-manager/db"
	"github.com/zarkopopovski/popup-manager/models"
)

type UserController struct {
	DBManager      *db.DBManager
	AuthController *AuthController
}

func (uController *UserController) Index(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("./templates/index.html"))

	tmpl.Execute(w, nil)
}

func (uController *UserController) RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	emailAddress := postMap["email"].(string)
	password := postMap["password"].(string)

	roles := "USER"
	//r.ParseForm()
	if r.PostForm.Has("roles") {
		roles = "USER:ADMIN"
	}

	mailContact := os.Getenv("MAIL_CONTACT")

	mailPort := os.Getenv("SMTP_PORT")
	mailServer := os.Getenv("MAIL_SERVER")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	hostname := os.Getenv("HOSTNAME")

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	sha1Hash = sha1.New()
	sha1Hash.Write([]byte(time.Now().String() + password + emailAddress))
	cfgToken := sha1Hash.Sum(nil)

	confirmationToken := fmt.Sprintf("%x", cfgToken)

	queryTest := "SELECT id, email, confirmed, last_login FROM user WHERE email=$1 LIMIT 1;"

	existingUser := new(models.User)

	err = uController.DBManager.DB.Get(existingUser, queryTest, emailAddress)

	if err == nil {
		w.WriteHeader(http.StatusForbidden)

		json.NewEncoder(w).Encode(map[string]string{"error": "The user already exists."})
		return
	}

	query := "INSERT INTO user(email, password, confirmed, last_login, date_created, date_modified, confirmation_token,roles) VALUES($1, $2, false, datetime('now'), datetime('now'), datetime('now'), $3, $4)"

	_, err = uController.DBManager.DB.Exec(query, emailAddress, passwordEnc, confirmationToken, roles)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if err == nil {
		go func() {
			m := gomail.NewMessage()

			m.SetHeader("From", mailContact)
			m.SetHeader("To", emailAddress)
			m.SetHeader("Subject", "Popup Message Ads: System notification")

			m.SetBody("text/html", `
				<p>Thank you for your registration.</p>
				<p>Please confirm your registration on the following link:</p>
				<p><a href="`+hostname+`"/confirm-registartion/`+confirmationToken+`>`+hostname+`/confirm-registartion/`+confirmationToken+`<a/></p>
			`)

			port, _ := strconv.Atoi(mailPort)

			d := gomail.NewDialer(mailServer, port, mailUsername, mailPassword)
			d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
			
			if err := d.DialAndSend(m); err != nil {
				log.Println(err.Error())
			}
		}()

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "-1"}); err != nil {
			log.Println(err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(map[string]interface{}{"status": "ok"}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Printf("%s", err)
	}
}

func (uController *UserController) RegisterAdminUser(emailAddress string, password string) error {
	
	roles := "USER:ADMIN"

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	sha1Hash = sha1.New()
	sha1Hash.Write([]byte(time.Now().String() + password + emailAddress))
	cfgToken := sha1Hash.Sum(nil)

	confirmationToken := fmt.Sprintf("%x", cfgToken)

	queryTest := "SELECT id, email, confirmed, last_login FROM user WHERE email=$1 LIMIT 1;"

	existingUser := new(models.User)

	err := uController.DBManager.DB.Get(existingUser, queryTest, emailAddress)

	if err == nil {
		return nil
	}

	query := "INSERT INTO user(email, password, confirmed, last_login, date_created, date_modified, confirmation_token,roles) VALUES($1, $2, true, datetime('now'), datetime('now'), datetime('now'), $3, $4)"

	_, err = uController.DBManager.DB.Exec(query, emailAddress, passwordEnc, confirmationToken, roles)

	if err == nil {
		return err
	}

	return nil
}

func (uController *UserController) ConfirmRegistration(w http.ResponseWriter, r *http.Request) {
	confirmationKey := r.PathValue("confirmationKey")

	mailContact := os.Getenv("MAIL_CONTACT")

	mailPort := os.Getenv("SMTP_PORT")
	mailServer := os.Getenv("MAIL_SERVER")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	user := models.User{}

	err := uController.DBManager.DB.Get(&user, "SELECT * FROM user WHERE confirmation_token=$1;", confirmationKey)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Not existing confirmation token"}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	if user.Confirmed {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusNotAcceptable)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Confirmation token already used; the user is likely already confirmed."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	_, err = uController.DBManager.DB.Exec("UPDATE user SET confirmed=true WHERE confirmation_token=$1", confirmationKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something got wrong..."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully confirmed"}); err != nil {
		log.Printf("%s", err)
	}

	go func() {
		m := gomail.NewMessage()

		m.SetHeader("From", mailContact)
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Popup Message Ads: System notification")

		m.SetBody("text/html", `
			<p>Thank you for your registration confirmation.</p>
			<p>You can now use your account.</p>
		`)

		port, _ := strconv.Atoi(mailPort)

		d := gomail.NewDialer(mailServer, port, mailUsername, mailPassword)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		
		if err := d.DialAndSend(m); err != nil {
			log.Println(err.Error())
		}
	}()
}

func (uController *UserController) SendTempPassPerMail(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	emailAddress := postMap["email"].(string)

	mailContact := os.Getenv("MAIL_CONTACT")

	mailPort := os.Getenv("SMTP_PORT")
	mailServer := os.Getenv("MAIL_SERVER")
	mailUsername := os.Getenv("MAIL_USERNAME")
	mailPassword := os.Getenv("MAIL_PASSWORD")

	user := models.User{}

	err = uController.DBManager.DB.Get(&user, "SELECT * FROM user WHERE email=$1;", emailAddress)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Not existing confirmation token"}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	rand.Seed(time.Now().UnixNano())

	password := uController.generatePassword(8, true, true)

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	query := "UPDATE user SET password=$1 WHERE id=$2;"

	_, err = uController.DBManager.DB.Exec(query, passwordEnc, user.Id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something got wrong..."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully confirmed"}); err != nil {
		log.Printf("%s", err)
	}

	go func() {
		m := gomail.NewMessage()

		m.SetHeader("From", mailContact)
		m.SetHeader("To", user.Email)
		m.SetHeader("Subject", "Popup Message Ads: System notification")

		m.SetBody("text/html", `
			<p>Your temporary password is:`+password+`</p>
			<p>Please change it!!!</p>
		`)

		port, _ := strconv.Atoi(mailPort)

		d := gomail.NewDialer(mailServer, port, mailUsername, mailPassword)
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		
		if err := d.DialAndSend(m); err != nil {
			log.Println(err.Error())
		}
	}()
}

func (uController *UserController) ChangePassword(w http.ResponseWriter, r *http.Request) {
	metaData, err := uController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := uController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	password := postMap["password"].(string)

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	query := "UPDATE user SET password=? WHERE id=?;"

	_, err = uController.DBManager.DB.Exec(query, passwordEnc, userID)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something got wrong..."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully changed"}); err != nil {
		log.Printf("%s", err)
	}
}

// TODO
func (uController *UserController) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	//USER
}

func (uController *UserController) generatePassword(length int, includeNumber bool, includeSpecial bool) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var password []byte
	var charSource string

	if includeNumber {
		charSource += "0123456789"
	}
	if includeSpecial {
		charSource += "!@#$%^&*()_+=-"
	}
	charSource += charset

	for i := 0; i < length; i++ {
		randNum := rand.Intn(len(charSource))
		password = append(password, charSource[randNum])
	}
	return string(password)
}
