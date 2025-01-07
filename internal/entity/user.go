package entity

type User struct {
	ID             string `json:"id"`
	UserType       string `json:"user_type"`
	UserRole       string `json:"user_role"`
	FullName       string `json:"full_name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	Bio            string `json:"bio"`
	Gender         string `json:"gender"`
	ProfilePicture string `json:"profile_picture"`
	AccessToken    string `json:"access_token"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type UserSingleRequest struct {
	ID       string `json:"id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
