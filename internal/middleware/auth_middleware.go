package middleware

import (

)

// AuthRequired memvalidasi header "Authorization: Bearer <token>".
// Jika valid, user_id & email disimpan di gin.Context agar bisa dipakai handler.