package utils

// SwaggerBasicResponse digunakan HANYA untuk dokumentasi contoh response di Swagger
type SwaggerBasicResponse struct {
	Data    string `json:"data" example:"string"`
	Error   string `json:"error" example:"detail error"`
	Message string `json:"message" example:"detail pesan"`
}

// Swagger401Response digunakan HANYA untuk dokumentasi contoh response 401 di Swagger
type Swagger401Response struct {
	Message string `json:"message" example:"token tidak valid atau sudah expired"`
}

// Swagger403Response digunakan HANYA untuk dokumentasi contoh response 403 di Swagger
type Swagger403Response struct {
	Message string `json:"message" example:"user tidak memiliki akses untuk fitur ini"`
}
