package models

type WebTokens struct {
	Id           int64  `json:"-" db:"id"`
	UserId       int64  `json:"-" db:"user_id"`
	WebURL       string `json:"web_url" db:"web_url"`
	ApiToken     string `json:"api_token" db:"api_token"`
	IsValid      int32  `json:"is_valid" db:"is_valid"`
	DateCreated  string `json:"date_created" db:"date_created"`
	DateModified string `json:"-" db:"date_modified"`
	Title        string `json:"title" db:"title"`
	Description  string `json:"description" db:"description"`
}

type WebTokenShort struct {
	Id           int64  `json:"-" db:"id"`
	UserId       int64  `json:"user_id" db:"user_id"`
	WebURL       string `json:"-" db:"web_url"`
	ApiToken     string `json:"api_token" db:"api_token"`
	IsValid      int32  `json:"is_valid" db:"is_valid"`
	DateCreated  string `json:"date_created" db:"date_created"`
	DateModified string `json:"-" db:"date_modified"`
	Title        string `json:"-" db:"title"`
	Description  string `json:"-" db:"description"`
}
