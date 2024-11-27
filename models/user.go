package models

type User struct {
	Id                int64         `json:"-" db:"id"`
	Email             string        `json:"email" db:"email"`
	Password          string        `json:"-" db:"password"`
	Confirmed         bool          `json:"-" db:"confirmed"`
	LastLogin         string        `json:"last_login" db:"last_login"`
	DateCreated       string        `json:"-" db:"date_created"`
	DateModified      string        `json:"-" db:"date_modified"`
	ConfirmationToken string        `json:"-" db:"confirmation_token"`
	Tokens            *TokenDetails `json:"tokens"`
	Roles             string        `json:"roles" db:"roles"`
}
