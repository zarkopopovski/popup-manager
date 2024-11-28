package models

type BasicStat struct {
	Id          int64  `json:"-" db:"id"`
	UserId      int64  `json:"-" db:"user_id"`
	ApiToken    string `json:"-" db:"api_token"`
	PopupId     int64  `json:"-" db:"popup_id"`
	Os          string `json:"os" db:"os"`
	Browser     string `json:"browser" db:"browser"`
	Country     string `json:"country" db:"country"`
	Area        string `json:"area" db:"area"`
	City        string `json:"city" db:"city"`
	DateCreated string `json:"date_created" db:"date_created"`
}
