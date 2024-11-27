package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
	"github.com/rs/cors"

	"github.com/zarkopopovski/popup-manager/controllers"
	"github.com/zarkopopovski/popup-manager/db"
)

type Handlers struct {
	Authentication  *controllers.AuthController
	UserController  *controllers.UserController
	PopupController *controllers.PopupController
}

type FileSystem struct {
	fs http.FileSystem
}

func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("The .env config file doesnt exist")
	}

	portNumber := os.Getenv("PORT")
	//database := os.Getenv("DATABASE")

	database := os.Getenv("DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")

	isUseGEOIP2 := os.Getenv("USE_GEOIP2")
	geoIP2Base := os.Getenv("GEOIP2_BASE")

	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	var isInitialStart bool = false

	if _, err := os.Stat("assets"); os.IsNotExist(err) {
		err := os.Mkdir("assets", 0777)
		if err != nil {
			log.Fatalln(err)
		}

		err = os.Mkdir("assets/uploads", 0777)
		if err != nil {
			log.Fatalln(err)
		}

		isInitialStart = true
	} else {
		if _, err := os.Stat("assets/uploads"); os.IsNotExist(err) {
			err = os.Mkdir("assets/uploads", 0777)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	databaseDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, database)

	httpRouter := http.NewServeMux()

	dbHandler := db.NewDBConnection(databaseDSN)

	var geoIP2DBBase *geoip2.Reader

	if isUseGEOIP2 != "" && !isInitialStart {
		useGEOIP2, err := strconv.ParseBool(isUseGEOIP2)

		if err == nil && useGEOIP2 && len(geoIP2Base) > 0 {

			if _, err = os.Stat("assets/" + geoIP2Base); os.IsNotExist(err) {
				log.Fatal(err)
			}

			geoIP2DBBase, err = geoip2.Open("assets/" + geoIP2Base)
			if err != nil {
				log.Fatal(err)
			}
			defer geoIP2DBBase.Close()
		}
	}

	authController := &controllers.AuthController{
		DBManager: dbHandler,
	}

	handlers := &Handlers{
		Authentication: authController,
		UserController: &controllers.UserController{
			DBManager:      dbHandler,
			AuthController: authController,
		},
		PopupController: &controllers.PopupController{
			DBManager:      dbHandler,
			AuthController: authController,
			GeoIPReader:    geoIP2DBBase,
		},
	}

	_ = handlers.UserController.RegisterAdminUser(adminUser, adminPassword)

	//PUBLIC
	httpRouter.HandleFunc("GET /api/v1/js/{apiToken}", handlers.PopupController.JSHandler)
	httpRouter.HandleFunc("POST /api/v1/login", handlers.Authentication.CheckUserCredentials)
	httpRouter.HandleFunc("POST /api/v1/logout", handlers.Authentication.Logout)
	httpRouter.HandleFunc("POST /api/v1/register-user", handlers.UserController.RegisterNewUser)
	httpRouter.HandleFunc("POST /api/v1/reset-password", handlers.UserController.SendTempPassPerMail)
	httpRouter.HandleFunc("GET /api/v1/confirm-registartion/{confirmationKey}", handlers.UserController.ConfirmRegistration)

	//PUBLIC NOTIFICATIONS
	httpRouter.HandleFunc("GET /api/v1/notification/{apiToken}", handlers.PopupController.ListPopopMessagesPerApiToken)
	httpRouter.HandleFunc("GET /api/v1/notification/{apiToken}/{notificationID}", handlers.PopupController.PushInstantPopUpMessagePerApiToken)
	httpRouter.HandleFunc("GET /api/v1/notification/{apiToken}/{notificationID}/trigger", handlers.PopupController.TriggerNotification)

	//USER
	httpRouter.HandleFunc("GET /api/v1/user/refresh-token/{refreshToken}", handlers.Authentication.Refresh)
	httpRouter.HandleFunc("POST /api/v1/user/change-password", handlers.UserController.ChangePassword)
	httpRouter.HandleFunc("POST /api/v1/user/user-details", handlers.UserController.UpdateUserDetails)

	//USER NOTIFICATIONS
	httpRouter.HandleFunc("POST /api/v1/user/web-site", handlers.PopupController.CreateApiToken)
	httpRouter.HandleFunc("PUT /api/v1/user/web-site/{apiToken}", handlers.PopupController.UpdateApiToken)
	httpRouter.HandleFunc("DELETE /api/v1/user/web-site/{apiToken}", handlers.PopupController.DeleteApiToken)
	httpRouter.HandleFunc("GET /api/v1/user/web-site", handlers.PopupController.ListAllApiToken)
	httpRouter.HandleFunc("POST /api/v1/user/notification", handlers.PopupController.CreatePopopMessage) // API_TOKEN as JSON parameter
	httpRouter.HandleFunc("PUT /api/v1/user/notification/{apiToken}/{notificationID}", handlers.PopupController.UpdatePopopMessage)
	httpRouter.HandleFunc("DELETE /api/v1/user/notification/{apiToken}/{notificationID}", handlers.PopupController.DeletePopopMessage)
	httpRouter.HandleFunc("GET /api/v1/user/notification/{apiToken}", handlers.PopupController.ListPopopMessages)
	httpRouter.HandleFunc("GET /api/v1/user/simple-stats", handlers.PopupController.GetBasicStatistics)
	httpRouter.HandleFunc("GET /api/v1/user/simple-stats-visits/{numRecords}", handlers.PopupController.GetLastXStatsSortedByDate)

	fileServer := http.FileServer(FileSystem{http.Dir("assets/uploads/")})
	httpRouter.Handle("/static/", http.StripPrefix(strings.TrimRight("/static/", "/"), fileServer))

	handler := cors.AllowAll().Handler(httpRouter)

	logger := log.New(os.Stdout, "popup-backend", log.LstdFlags)
	logger.Println("Start Listening on port:" + portNumber)

	thisServer := &http.Server{
		Addr:         ":" + portNumber,
		Handler:      handler,
		IdleTimeout:  120 + time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		err := thisServer.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	thisSignalChan := <-sigChan

	logger.Println("Graceful Shutdown", thisSignalChan)

	timeOutContext, canFunct := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer canFunct()

	thisServer.Shutdown(timeOutContext)
}
