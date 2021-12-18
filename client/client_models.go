package client

type AuthorizationCodeClientCreate struct {
	ClientId      string   `json:"clientId"`
	ClientName    string   `json:"clientName"`
	AllowedScopes []string `json:"allowedScopes"`
	RedirectUris  []string `json:"redirectUris"`
}

type AuthorizationCodeClientUpdate struct {
	ClientName    string   `json:"clientName,omitempty"`
	AllowedScopes []string `json:"allowedScopes,omitempty"`
	RedirectUris  []string `json:"redirectUris,omitempty"`
	Enabled       bool     `json:"enabled,omitempty"`
}

type AuthorizationCodeClient struct {
	ClientId            string
	ClientName          string
	AllowedScopes       []string
	RedirectUris        []string
	Enabled             bool
	RequireClientSecret bool
	RequirePkce         bool
	ClientSecrets       []ClientSecret
	AllowOfflineAccess  bool
	AllowedGrantTypes   []string
}

type ClientSecret struct {
	Value string
	Type  string
}
