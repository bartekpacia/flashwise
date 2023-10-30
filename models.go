package main

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type CreateFlashcardSetRequest struct {
	Title      string `json:"name"`
	Public     bool   `json:"is_public"`
	CategoryID uint64 `json:"category"`
}
