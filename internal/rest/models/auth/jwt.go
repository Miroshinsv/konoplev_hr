package auth

type JWT struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}
