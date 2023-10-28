package main

type CreateUserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

type CreateFlashcardRequest struct {
	Front string `json:"front"`
	Back  string `json:"back"`
	SetID uint64 `json:"flashcard_set"`
}

type CreateFlashcardSetRequest struct {
	Status bool   `json:"status"`
	Title  string `json:"name"`
	Public bool   `json:"is_public"`
}
