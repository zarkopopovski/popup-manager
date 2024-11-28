package models

import "database/sql"

type PopUpMessage struct {
	Id           int64          `json:"id" db:"id"`
	UserId       int64          `json:"-" db:"user_id"`
	ApiToken     string         `json:"-" db:"api_token"`
	PopupType    int16          `json:"pop_type" db:"popup_type"`
	Title        string         `json:"title" db:"title"`
	Description  string         `json:"description" db:"description"`
	Enabled      bool           `json:"enabled" db:"enabled"`
	DateCreated  string         `json:"date_created" db:"date_created"`
	DateModified string         `json:"-" db:"date_modified"`
	ShowTime     int64          `json:"show_time" db:"show_time"`
	CloseTime    int64          `json:"close_time" db:"close_time"`
	PopupPos     int16          `json:"popup_pos" db:"popup_pos"`
	ImageName    sql.NullString `json:"image_name" db:"image_name"`
	IsTrackable  bool           `json:"is_trackable" db:"is_trackable"`
}
