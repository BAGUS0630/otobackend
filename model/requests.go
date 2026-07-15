package model

// ProfileUpdateRequest represents payload to update user profile
type ProfileUpdateRequest struct {
    FullName    string `json:"full_name" example:"Budi Santoso"`
    Email       string `json:"email" example:"budi@example.com"`
    PhoneNumber string `json:"phone_number" example:"08123456789"`
}

// AccountDeleteRequest represents payload to confirm account deletion
type AccountDeleteRequest struct {
    Password string `json:"password" example:"your_password"`
}
