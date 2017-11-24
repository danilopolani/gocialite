package structs

// User struct
type User struct {
	ID        string                 `json:"id"`
	Username  string                 `json:"username"`
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	FullName  string                 `json:"full_name"`
	Email     string                 `json:"email"`
	Avatar    string                 `json:"avatar"`
	Raw       map[string]interface{} `json:"raw"` // Raw data
}
