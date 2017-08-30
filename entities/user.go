package entities

type User struct {
	firstname         string
	lastname          string
	stravaId          int64
	email             string
	stravaAccessToken string
}

func NewUser(firstname, lastname, email string, stravaId int64, stravaAccessToken string) User {
	return User{
		firstname:         firstname,
		lastname:          lastname,
		email:             email,
		stravaId:          stravaId,
		stravaAccessToken: stravaAccessToken,
	}
}

func (u *User) GetFirstname() string {
	return u.firstname
}

func (u *User) GetStravaId() int64 {
	return u.stravaId
}

func (u *User) GetLastname() string {
	return u.lastname
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetStravaAccessToken() string {
	return u.stravaAccessToken
}
