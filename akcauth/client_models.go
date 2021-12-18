package akcauth

type AuthorizationCodeClientCreate struct {
	ClientId      string   `json:"clientId"`
	ClientName    string   `json:"clientNAme"`
	AllowedScopes []string `json:"allowedScopes"`
	RedirectUris  []string `json:"redirectUris"`
}
