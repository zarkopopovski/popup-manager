package controllers

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mileusna/useragent"
	"github.com/oschwald/geoip2-golang"
	"github.com/twinj/uuid"
	"github.com/zarkopopovski/popup-manager/db"
	"github.com/zarkopopovski/popup-manager/models"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 0.5 // 500k

type PopupController struct {
	DBManager      *db.DBManager
	AuthController *AuthController
	GeoIPReader    *geoip2.Reader
}

type PopupStats struct {
	NumWebSiteTokens int32 `json:"num_site_tokens"`
	MumWebPopups     int32 `json:"num_web_popups"`
	NumClicks        int64 `json:"num_clicks"`
}

func (popUpController *PopupController) CreatePopopMessage(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	//TODO: CHECK SUBSCRIPTIONS

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if errSize := r.ParseMultipartForm(MAX_UPLOAD_SIZE); errSize != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 50MB in size", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")

	parameter1 := r.FormValue("parameter_1")
	if parameter1 == "" {
		parameter1 = "S3cREtF1L3Up&0@d"
	}

	apiToken := r.FormValue("api_token")
	popupType := r.FormValue("type")
	title := r.FormValue("title")
	description := r.FormValue("description")

	showTime := r.FormValue("show_time")
	closeTime := r.FormValue("close_time")

	popupPos := r.FormValue("popup_pos")

	isTrackable := r.FormValue("is_trackable")

	isFileUploadedError := false

	if err != nil {
		isFileUploadedError = true
	}

	fileName := ""

	if !isFileUploadedError {
		defer file.Close()

		fileName = header.Filename

		randomFloat := strconv.FormatFloat(rand.Float64(), 'E', -1, 64)

		sha1Hash := sha1.New()
		sha1Hash.Write([]byte(time.Now().String() + parameter1 + fileName + randomFloat))
		sha1HashString := sha1Hash.Sum(nil)

		fileNameHASH := fmt.Sprintf("%x", sha1HashString)

		fileName = fileNameHASH + "$" + fileName

		out, err := os.Create("./assets/uploads/" + fileName)

		if err != nil {
			fmt.Fprintf(w, "Unable to create a file for writting. Check your write access privilege")
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)

		if err != nil {
			fmt.Fprintln(w, err)
		}
	}

	queryStr := "INSERT INTO popup_message(user_id, api_token, popup_type, title, description, enabled, date_created, date_modified, show_time, close_time, popup_pos, image_name, is_trackable) VALUES($1, $2, $3, $4, $5, $6, datetime('now'), datetime('now'), $7, $8, $9, $10, $11)"

	_, err = popUpController.DBManager.DB.Exec(queryStr, userID, apiToken, popupType, title, description, true, showTime, closeTime, popupPos, fileName, isTrackable)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "-1"})
}

func (popUpController *PopupController) UpdatePopopMessage(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	apiToken := r.PathValue("apiToken")

	notificationID := r.PathValue("notificationID")

	popUpMessage := models.PopUpMessage{}

	notifQuery := "SELECT * FROM popup_message WHERE id=$1 AND api_token=$2 AND user_id=$3"

	err = popUpController.DBManager.DB.Get(&popUpMessage, notifQuery, notificationID, apiToken, userID)
	if err != nil {
		http.Error(w, "This popup doesnt exist", http.StatusNotFound)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if errSize := r.ParseMultipartForm(MAX_UPLOAD_SIZE); errSize != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than 500k in size", http.StatusBadRequest)
		return
	}

	file, header, errFile := r.FormFile("file")

	parameter1 := r.FormValue("parameter_1")
	if parameter1 == "" {
		parameter1 = "S3cREtF1L3Up&0@d"
	}

	popupType := r.FormValue("type")
	title := r.FormValue("title")
	description := r.FormValue("description")

	showTime := r.FormValue("show_time")
	closeTime := r.FormValue("close_time")

	popupPos := r.FormValue("popup_pos")

	isTrackable := r.FormValue("is_trackable")

	isFileUploadedError := false

	if errFile != nil {
		isFileUploadedError = true
	}

	fileName := ""

	if !isFileUploadedError {
		_, err := os.Stat("./assets/uploads/" + popUpMessage.ImageName.String)
		if popUpMessage.ImageName.Valid && err == nil {
			err = os.Remove("./assets/uploads/" + popUpMessage.ImageName.String)
			if err != nil {
				fmt.Printf("File %s can't be deleted", popUpMessage.ImageName.String)
			}
		}

		defer file.Close()

		fileName = header.Filename

		randomFloat := strconv.FormatFloat(rand.Float64(), 'E', -1, 64)

		sha1Hash := sha1.New()
		sha1Hash.Write([]byte(time.Now().String() + parameter1 + fileName + randomFloat))
		sha1HashString := sha1Hash.Sum(nil)

		fileNameHASH := fmt.Sprintf("%x", sha1HashString)

		fileName = fileNameHASH + "$" + fileName

		out, err := os.Create("./assets/uploads/" + fileName)

		if err != nil {
			fmt.Fprintf(w, "Unable to create a file for writting. Check your write access privilege")
			return
		}

		defer out.Close()

		_, err = io.Copy(out, file)

		if err != nil {
			fmt.Fprintln(w, err)
		}
	} else {
		_, err := os.Stat("./assets/uploads/" + popUpMessage.ImageName.String)
		if popUpMessage.ImageName.Valid && err == nil {
			err = os.Remove("./assets/uploads/" + popUpMessage.ImageName.String)
			if err != nil {
				fmt.Printf("File %s can't be deleted", popUpMessage.ImageName.String)
			}
		}
	}
	//TODO: CHECK IF THE REPETITION TIMES ARE CHANGED
	queryStr := "UPDATE popup_message SET popup_type=$1, title=$2, description=$3, enabled=$4, date_modified=datetime('now'), show_time=$5, close_time=$6, popup_pos=$7, image_name=$8, is_trackable=$9 WHERE id=$10 AND api_token=$11 AND user_id=$12"

	_, err = popUpController.DBManager.DB.Exec(queryStr, popupType, title, description, true, showTime, closeTime, popupPos, fileName, isTrackable, notificationID, apiToken, userID)

	if err != nil {
		log.Println(err.Error())

		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "-1"})
}

func (popUpController *PopupController) DeletePopopMessage(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	apiToken := r.PathValue("apiToken")

	notificationID := r.PathValue("notificationID")

	queryStr := "DELETE FROM popup_message WHERE id=$1 AND api_token=$2 AND user_id=$3"

	_, err = popUpController.DBManager.DB.Exec(queryStr, notificationID, apiToken, userID)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "-1"})
}

func (popUpController *PopupController) ListPopopMessages(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	apiToken := r.PathValue("apiToken")

	popupMessages := make([]models.PopUpMessage, 0)

	err = popUpController.DBManager.DB.Select(&popupMessages, "SELECT * FROM popup_message WHERE user_id=$1 AND api_token=$2;", userID, apiToken)
	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "3", "message": "Not Found"})
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "error_code": "-1", "data": popupMessages})
}

func (popUpController *PopupController) ListPopopMessagesPerApiToken(w http.ResponseWriter, r *http.Request) {
	popupMessages := make([]models.PopUpMessage, 0)

	apiToken := r.PathValue("apiToken")

	err := popUpController.DBManager.DB.Select(&popupMessages, "SELECT * FROM popup_message WHERE api_token=$1;", apiToken)
	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "3", "message": "Not Found"})
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "error_code": "-1", "data": popupMessages})
}

func (popUpController *PopupController) PushInstantPopUpMessagePerApiToken(w http.ResponseWriter, r *http.Request) {
	popupMessage := models.PopUpMessage{}

	apiToken := r.PathValue("apiToken")
	notificationID := r.PathValue("notificationID")

	err := popUpController.DBManager.DB.Get(&popupMessage, "SELECT * FROM popup_message WHERE api_token=$1 AND id=$2;", apiToken, notificationID)
	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "3", "message": "Not Found"})
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "error_code": "-1", "data": popupMessage})
}

func (popUpController *PopupController) TriggerNotification(w http.ResponseWriter, r *http.Request) {
	apiToken := r.PathValue("apiToken")
	notificationID := r.PathValue("notificationID")

	userAgent := r.UserAgent()

	clientIP := r.Header.Get("X-FORWARDED-FOR")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}

	if clientIP == "" {
		clientIP = "8.8.8.8"
	}

	ip := net.ParseIP(clientIP)
	record, err := popUpController.GeoIPReader.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	countryCode := strings.ToLower(record.Country.IsoCode)
	countryName := record.Country.Names[countryCode]
	cityName := record.City.Names[countryCode]
	countryArea := ""

	if len(record.Subdivisions) > 0 {
		countryArea = record.Subdivisions[0].Names[countryCode]
	}

	parsedUA := useragent.Parse(userAgent)
	opSystem := parsedUA.OS
	browserName := parsedUA.Name

	queryWeb := "SELECT * FROM web_tokens WHERE api_token=$1;"

	webToken := models.WebTokenShort{}

	err = popUpController.DBManager.DB.Get(&webToken, queryWeb, apiToken)
	if err == nil {
		queryStr := "INSERT INTO basic_stats(user_id, api_token, popup_id, os, browser, country, area, city, date_created) VALUES($1, $2, $3, $4, $5, $6, $7, $8, datetime('now'));"

		_, _ = popUpController.DBManager.DB.Exec(queryStr, webToken.UserId, apiToken, notificationID, opSystem, browserName, countryName, countryArea, cityName)
	}

	w.WriteHeader(http.StatusOK)
}

func (popUpController *PopupController) GetLastXStatsSortedByDate(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	numRecords := r.PathValue("numRecords")

	basicStats := make([]models.BasicStat, 1)

	queryStr := "SELECT * FROM basic_stats WHERE user_id=$1 ORDER BY date_created DESC LIMIT " + numRecords

	err = popUpController.DBManager.DB.Select(&basicStats, queryStr, userID)

	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "error_code": "-1", "data": basicStats})
}

func (popUpController *PopupController) CreateApiToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := popUpController.AuthController.FetchAuth(metaData)
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

	title := postMap["title"].(string)
	description := postMap["description"].(string)
	webURL := postMap["web_url"].(string)

	// sid, _ := shortid.New(1, shortid.DefaultABC, 2342)
	// shortid.SetDefault(sid)

	// webApiToken, _ := sid.Generate()
	webApiToken := uuid.NewV4().String()

	queryStr := "INSERT INTO web_tokens(user_id, web_url, api_token, is_valid, date_created, date_modified, title, description) VALUES($1, $2, $3, $4, datetime('now'), datetime('now'), $5, $6)"

	_, err = popUpController.DBManager.DB.Exec(queryStr, userID, webURL, webApiToken, true, title, description)

	if err != nil {
		log.Printf("%s", err.Error())

		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something got wrong..."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully created"}); err != nil {
		log.Printf("%s", err)
	}
}

func (popUpController *PopupController) UpdateApiToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	apiToken := r.PathValue("apiToken")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	title := postMap["title"].(string)
	description := postMap["description"].(string)
	webURL := postMap["web_url"].(string)
	// apiToken := postMap["api_token"].(string)
	// isValid := postMap["is_valid"].(bool)

	queryStr := "UPDATE web_tokens SET web_url=$1, date_modified=datetime('now'), title=$2, description=$3 WHERE api_token=$4 AND user_id=$5"

	_, err = popUpController.DBManager.DB.Exec(queryStr, webURL, title, description, apiToken, userID)

	if err != nil {
		log.Println(err.Error())

		w.Header().Set("Content-Type", "application/json; charset=UTF8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "Something got wrong..."}); err != nil {
			log.Printf("%s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully updated"}); err != nil {
		log.Printf("%s", err)
	}
}

func (popUpController *PopupController) DeleteApiToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	apiToken := r.PathValue("apiToken")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	popupMessages := make([]models.PopUpMessage, 0)

	err = popUpController.DBManager.DB.Select(&popupMessages, "SELECT * FROM popup_message WHERE user_id=$1 AND api_token=$2;", userID, apiToken)
	if err != nil {
		log.Println(err.Error())
	}

	if len(popupMessages) > 0 {
		for _, popupMessage := range popupMessages {
			queryStrDelete := "DELETE FROM pop_timing WHERE user_id=$1, popup_id=$2"

			_, err = popUpController.DBManager.DB.Exec(queryStrDelete, userID, popupMessage.Id)

			if err != nil {
				log.Println(err.Error())
			}
		}

		queryStr := "DELETE FROM popup_message WHERE api_token=$1 AND user_id=$2"

		_, err = popUpController.DBManager.DB.Exec(queryStr, apiToken, userID)

		if err != nil {
			log.Println(err.Error())
		}
	}

	queryStr := "DELETE FROM web_tokens WHERE api_token=$1 AND user_id=$2"

	_, err = popUpController.DBManager.DB.Exec(queryStr, apiToken, userID)

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
	if err := json.NewEncoder(w).Encode(map[string]string{"error": "Successfully deleted"}); err != nil {
		log.Printf("%s", err)
	}
}

func (popUpController *PopupController) ListAllApiToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}
	queryStr := "SELECT * FROM web_tokens WHERE user_id=$1 ORDER BY date_created DESC"

	webTokens := make([]models.WebTokens, 0)

	err = popUpController.DBManager.DB.Select(&webTokens, queryStr, userID)

	if err != nil {
		log.Println(err.Error())

		w.WriteHeader(http.StatusNotFound)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "3", "message": "Not Found"})
		return
	}

	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "error_code": "-1", "data": webTokens})
}

func (popUpController *PopupController) CreateDirIfNotExist(dir string) {
	if _, err := os.Stat("store/" + dir); os.IsNotExist(err) {
		err = os.MkdirAll("store/"+dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func (popUpController *PopupController) GetBasicStatistics(w http.ResponseWriter, r *http.Request) {
	metaData, err := popUpController.AuthController.ExtractTokenMetadata(r)
	if err != nil {
		fmt.Println(err.Error())

		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	userID, err := popUpController.AuthController.FetchAuth(metaData)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)

		_ = json.NewEncoder(w).Encode(map[string]string{"status": "error", "error_code": "1", "message": "Forbidden access"})
		return
	}

	var numWebSiteTokens int32
	var numWebPopups int32
	var numClick int64

	queryNumWebSites := "SELECT COUNT(*) AS num_web_site_tokens FROM web_tokens WHERE user_id=$1;"
	queryNumWebPopups := "SELECT COUNT(*) AS num_web_popups FROM popup_message WHERE user_id=$1;"
	queryNumClicks := "SELECT COUNT(*) AS num_clicks FROM basic_stats WHERE user_id=$1;"
	_ = popUpController.DBManager.DB.Get(&numWebSiteTokens, queryNumWebSites, userID)
	_ = popUpController.DBManager.DB.Get(&numWebPopups, queryNumWebPopups, userID)
	_ = popUpController.DBManager.DB.Get(&numClick, queryNumClicks, userID)

	popupStats := &PopupStats{
		NumWebSiteTokens: numWebSiteTokens,
		MumWebPopups:     numWebPopups,
		NumClicks:        numClick,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"status": "error", "error_code": "-1", "data": popupStats})
}

func (popUpController *PopupController) JSHandler(w http.ResponseWriter, r *http.Request) {
	apiToken := r.PathValue("apiToken")

	apiHostname := os.Getenv("HOSTNAME_API")

	workload := `
		(()=>{
			function getAjax(url, success) {
				var xhr = window.XMLHttpRequest ? new XMLHttpRequest() : new ActiveXObject('Microsoft.XMLHTTP');
				xhr.open('GET', url);
				xhr.onreadystatechange = function() {
					if (xhr.readyState>3 && xhr.status==200) success(xhr.responseText);
				};
				xhr.setRequestHeader('X-Requested-With', 'XMLHttpRequest');
				xhr.send();
				return xhr;
			}

			var sheet = document.createElement('style');
			sheet.innerHTML = ".toast-container{position:fixed;width:300px;max-height:calc(100vh - 40px);overflow-y:auto;z-index:999999999}.toast{background-color:#333;color:#fff;padding:10px;border-radius:5px;margin-bottom:10px;display:flex;align-items:center;justify-content:flex-start}.toast-icon{margin-right:10px;width:48px;height:48px;border-radius:5px}.toast-message{text-align:left}.toast-message div{margin-top:5px}";

			document.body.appendChild(sheet); // append in body
			document.head.appendChild(sheet); // append in head

			let toastId = 0;

			showToast = (popupID = -1, position = 'top-right', message1 = '', message2 = '', image = '', showTime = 0, closeTime = 0, popupType = -1, isTrackable = false) => {
				//const toastContainer = document.getElementById('toastContainer');
				// Create toast container element
				let toastContainer = document.querySelector('.toast-container.' + position);
				if (!toastContainer) {
					toastContainer = document.createElement('div');
					toastContainer.classList.add('toast-container', position);
					document.body.appendChild(toastContainer);
				}

				// Create toast element
				const toast = document.createElement('div');
				toast.classList.add('toast');
				toast.id = 'toast-${toastId++}';
				
				let fileName = Object.values(image)[0];

				const icon = document.createElement('div');
				if (fileName !== '') {
					// Create icon element
					icon.classList.add('toast-icon');
					icon.style.backgroundImage = 'url("$hostName$/static/'+fileName+'")'; // Set your image URL here
					icon.style.backgroundSize = 'contain';
				}
				// Create message elements
				const messageContainer = document.createElement('div');
				messageContainer.classList.add('toast-message');
				const messageLine1 = document.createElement('div');
				messageLine1.textContent = message1;
				const messageLine2 = document.createElement('div');
				messageLine2.textContent = message2;

				// Append icon and messages to toast
				messageContainer.appendChild(messageLine1);
				messageContainer.appendChild(messageLine2);

				if (fileName !== '') {
					toast.appendChild(icon); 
				}
				
				toast.appendChild(messageContainer);

				// Append toast to container
				toastContainer.appendChild(toast);

				// Apply animation
				setTimeout(() => {
					toast.classList.add('show');
					
					if (closeTime > 0) {
						setTimeout(() => {
							toast.classList.remove('show');
							setTimeout(() => {
								toast.remove(); // Remove toast after animation
							}, 300);
						}, closeTime); // Close after time > 0
					}
				}, showTime); // Delay before showing toast
				
				// Close toast on click
				toast.addEventListener('click', () => {
					if (isTrackable) {
						getAjax('$hostName$/api/v1/notification/$apiToken$/' + popupID + '/trigger', (res) => { });
					}

					toast.remove();
					if (toastContainer.childNodes.length === 0) {
						toastContainer.remove(); // Remove toast container if empty
					}
				});

				// Set toast position
				switch (position) {
					case 'top-right':
						toastContainer.style.top = '20px';
						toastContainer.style.right = '20px';
						break;
					case 'top-left':
						toastContainer.style.top = '20px';
						toastContainer.style.left = '20px';
						break;
					case 'bottom-left':
						toastContainer.style.bottom = '20px';
						toastContainer.style.left = '20px';
						break;
					case 'bottom-right':
						toastContainer.style.bottom = '20px';
						toastContainer.style.right = '20px';
						break;
					case 'center':
						toastContainer.style.top = '50%';
						toastContainer.style.left = '50%';
						toastContainer.style.transform = 'translate(-50%, -50%)';
						break;
					default:
						console.error('Invalid position specified for toast notification.');
				}

				switch (popupType) {
					case 1: 
						toast.style.backgroundColor = '#333';
						break;
					case 2: 
						toast.style.backgroundColor = '#8BC34A';
						break;
					case 3: 
						toast.style.backgroundColor = '#FFEB3B';
						break;
					case 4: 
						toast.style.backgroundColor = '#FF9800';
						break;
					case 5: 
						toast.style.backgroundColor = '#F44336';
						break;
					default:	
						toast.style.backgroundColor = '#333';
						break;						
				}
			}

			
			checkPopupPosition = (id) => {
				switch(id) {
					case 1: return 'top-right';
					case 2: return 'top-left';
					case 3: return 'bottom-right';
					case 4: return 'bottom-left';
					case 5: return 'center';
					default: return '';
					}
			}
			
			getAjax('$hostName$/api/v1/notification/$apiToken$', (res) => {
				const jsonObj = JSON.parse(res);
				if (jsonObj.error_code === '-1') {
					const popupsArray = jsonObj.data;
					if (popupsArray.length > 0) {
						for (let i = 0; i < popupsArray.length; i++) {
							const obj1 = popupsArray[i];
							setTimeout(() => {
								showToast(
									obj1.id,
									checkPopupPosition(obj1.popup_pos),
									obj1.title,
									obj1.description,
									obj1.image_name,
									obj1.show_time,
									obj1.close_time,
									obj1.pop_type,
									obj1.is_trackable
								);
							}, obj1.show_time);
						}
					}
				}
			});
		})();
	`

	workload = strings.ReplaceAll(workload, "$apiToken$", apiToken)
	workload = strings.ReplaceAll(workload, "$hostName$", apiHostname)

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(workload))
}
