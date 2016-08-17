package structs

type FacebookUserStruct struct {
	Email        string `json:"email"`
	FacebookID   string `json:"facebook_id"`
	FacebookName string `json:"facebook_name"`
}

type UserStruct struct {
	UUID     string `json:"uuid"`
	Email    string `json:"email"`
	Username string `json:"username"`
}
