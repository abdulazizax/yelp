package user_entity

type (
	User struct {
		ID             string `json:"id" binding:"required,uuid"` // UUID for the primary key
		UserType       string `json:"user_type" binding:"required,oneof=user admin business_owner"`
		UserRole       string `json:"user_role" binding:"required,oneof=user admin business_owner super_admin"`
		Name           string `json:"name" binding:"required,max=100"`
		Email          string `json:"email" binding:"required,email,max=255"`
		PasswordHash   string `json:"password_hash" binding:"required"`
		Bio            string `json:"bio" binding:"omitempty"`
		Gender         string `json:"gender" binding:"required,oneof=male female" default:"male"`
		ProfilePicture string `json:"profile_picture" binding:"omitempty,url"`
		Status         string `json:"status" binding:"required,oneof=active blocked inverify" default:"inverify"`
		CreatedAt      string `json:"created_at" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
		UpdatedAt      string `json:"updated_at" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	}

	CreateUser struct {
		Name     string `json:"name" binding:"required,max=100"`
		Email    string `json:"email" binding:"required,email,max=255"`
		Password string `json:"password" binding:"required"`
		Gender   string `json:"gender" binding:"required,oneof=male female" default:"male"`
	}

	SignInUser struct {
		Email    string `json:"email" binding:"required,email,max=255"`
		Password string `json:"password" binding:"required"`
	}

	CreateSession struct {
		UserID    string `json:"user_id" binding:"required,uuid"`
		UserAgent string `json:"user_agent" binding:"required"`
		Platform  string `json:"platform" binding:"required,oneof=web mobile admin_panel"`
		IPAddress string `json:"ip_address" binding:"required"`
	}

	Session struct {
		ID        string `json:"id" binding:"required,uuid"`
		UserID    string `json:"user_id" binding:"required,uuid"`
		UserAgent string `json:"user_agent" binding:"required"`
		Platform  string `json:"platform" binding:"required,oneof=web mobile admin_panel"`
		IPAddress string `json:"ip_address" binding:"required"`
		CreatedAt string `json:"created_at" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
		UpdatedAt string `json:"updated_at" binding:"required" time_format:"2006-01-02T15:04:05Z07:00"`
	}
)
