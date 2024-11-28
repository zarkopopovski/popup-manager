package controllers

import (
	"errors"
	"fmt"
	"io"

	"net/http"

	"strconv"

	"encoding/json"
	"strings"
	"time"

	"crypto/sha1"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"

	"github.com/zarkopopovski/popup-manager/db"
	"github.com/zarkopopovski/popup-manager/models"
)

type AuthController struct {
	DBManager     *db.DBManager
	AccessSecret  string
	RefreshSecret string
	AdminUser     string
	AdminPassword string
}

type Exception struct {
	Message string `json:"message"`
}

func (aController *AuthController) CheckUserCredentials(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var postMap map[string]interface{}

	json.Unmarshal([]byte(b), &postMap)

	email := postMap["email"].(string)
	password := postMap["password"].(string)

	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(password))
	sha1HashString := sha1Hash.Sum(nil)

	passwordEnc := fmt.Sprintf("%x", sha1HashString)

	query := "SELECT id, email, confirmed, last_login FROM user WHERE email=$1 AND password=$2"

	newUser := new(models.User)

	err = aController.DBManager.DB.QueryRowx(query, email, passwordEnc).StructScan(newUser)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)

		json.NewEncoder(w).Encode(map[string]string{"error": "The user is not found"})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if !newUser.Confirmed {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{"error": "The user is not confirmed, please confirm your registartion first."})
		return
	}

	ts, err := aController.CreateToken(strconv.Itoa(int(newUser.Id)))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return
	}

	err = aController.CreateAuth(newUser.Id, ts)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return
	}

	newUser.Tokens = ts

	_, err = aController.DBManager.DB.Exec("UPDATE user SET last_login=datetime('now') WHERE id=$1", newUser.Id)

	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{"data": newUser}); err != nil {
		panic(err)
	}

}

func (aController *AuthController) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := aController.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(aController.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (aController *AuthController) ExtractTokenMetadata(r *http.Request) (*models.AccessDetails, error) {
	token, err := aController.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &models.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func (aController *AuthController) FetchAuth(authD *models.AccessDetails) (int64, error) {
	var userID int64 = -1

	query := "SELECT user_id FROM tokens WHERE uuid=$1"

	err := aController.DBManager.DB.QueryRow(query, authD.AccessUuid).Scan(&userID)

	if err != nil {
		return 0, err
	}

	//userID, _ = strconv.ParseUint(userid, 10, 64)
	if authD.UserId != userID {
		return 0, errors.New("unauthorized")
	}
	return int64(userID), nil
}

func (aController *AuthController) TokenValid(r *http.Request) error {
	token, err := aController.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return err
	}
	return nil
}

func (aController *AuthController) Refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.URL.Query().Get("refreshToken")

	w.Header().Set("Content-Type", "application/json")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(aController.RefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		fmt.Println("the error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: "Refresh token expired"})
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}
		userId, err := strconv.ParseInt(fmt.Sprintf("%v", claims["user_id"]), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(Exception{Message: "Error occurred"})
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := aController.DeleteAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(Exception{Message: "unauthorized"})
			return
		}
		//Create new pairs of refresh and access tokens
		userID := claims["user_id"]

		ts, err := aController.CreateToken(userID.(string))
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: err.Error()})
			return
		}
		//save the tokens metadata to redis
		saveErr := aController.CreateAuth(userId, ts)
		if saveErr != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(Exception{Message: saveErr.Error()})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(tokens)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(Exception{Message: "refresh expired"})
	}
}

func (aController *AuthController) CreateToken(userID string) (*models.TokenDetails, error) {
	td := &models.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.AccessUuid + "++" + userID

	var err error
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(aController.AccessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(aController.RefreshSecret))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func (aController *AuthController) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (aController *AuthController) CreateAuth(userid int64, td *models.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	query := "INSERT INTO tokens(type, uuid, user_id, date_created) VALUES($1, $2, $3, $4);"

	_, err := aController.DBManager.DB.Exec(query, "ACCESS_UUID", td.AccessUuid, strconv.Itoa(int(userid)), at.Sub(now))
	if err != nil {
		return err
	}
	_, err = aController.DBManager.DB.Exec(query, "REFRESH_UUID", td.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now))
	if err != nil {
		return err
	}
	return nil
}

func (aController *AuthController) DeleteTokens(authD *models.AccessDetails) error {
	//get the refresh uuid
	refreshUuid := fmt.Sprintf("%s++%d", authD.AccessUuid, authD.UserId)
	//delete access token
	query := "DELETE FROM tokens WHERE uuid=$1 AND type=$2;"

	_, err := aController.DBManager.DB.Exec(query, authD.AccessUuid, "ACCESS_UUID")
	if err != nil {
		return err
	}

	_, err = aController.DBManager.DB.Exec(query, refreshUuid, "REFRESH_UUID")
	if err != nil {
		return err
	}

	return nil
}

func (aController *AuthController) DeleteAuth(givenUuid string) (int64, error) {
	query := "DELETE FROM tokens WHERE uuid=$1 AND type=$2;"

	_, err := aController.DBManager.DB.Exec(query, givenUuid, "REFRESH_UUID")
	if err != nil {
		return 0, err
	}

	return -1, nil
}

func (aController *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	metadata, err := aController.ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return
	}
	err = aController.DeleteTokens(metadata)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(Exception{Message: err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"status": "success", "error_code": "-1"}); err != nil {
		panic(err)
	}
}
