package client

type AuthorizationCodeClientCreate struct {
	ClientId          string   `json:"clientId"`
	ClientName        string   `json:"clientName"`
	AllowedScopes     []string `json:"allowedScopes"`
	RedirectUris      []string `json:"redirectUris"`
	AllowedGrantTypes []string `json:"allowedGrantTypes"`
}

type AuthorizationCodeClientUpdate struct {
	ClientName        string   `json:"clientName,omitempty"`
	AllowedScopes     []string `json:"allowedScopes,omitempty"`
	RedirectUris      []string `json:"redirectUris,omitempty"`
	Enabled           bool     `json:"enabled"`
	AllowedGrantTypes []string `json:"allowedGrantTypes"`
}

type AuthorizationCodeClient struct {
	ClientId            string         `json:"clientId"`
	ClientName          string         `json:"clientName"`
	AllowedScopes       []string       `json:"allowedScopes"`
	RedirectUris        []string       `json:"redirectUris"`
	Enabled             bool           `json:"enabled"`
	RequireClientSecret bool           `json:"requireClientSecret"`
	RequirePkce         bool           `json:"requirePkce"`
	ClientSecrets       []ClientSecret `json:"clientSecrets"`
	AllowOfflineAccess  bool           `json:"allowOfflineAccess"`
	AllowedGrantTypes   []string       `json:"allowedGrantTypes"`
}

type ClientSecret struct {
	Value string
	Type  string
}

func (source AuthorizationCodeClient) ToUpdateModel() AuthorizationCodeClientUpdate {
	updateModel := AuthorizationCodeClientUpdate{
		ClientName:        source.ClientName,
		AllowedScopes:     source.AllowedScopes,
		RedirectUris:      source.RedirectUris,
		AllowedGrantTypes: source.AllowedGrantTypes,
		Enabled:           source.Enabled,
	}

	return updateModel
}

func (source ApiResource) ToUpdateModel() ApiResourceUpdate {
	updateModel := ApiResourceUpdate{
		DisplayName: source.DisplayName,
		Scopes:      source.Scopes,
	}

	return updateModel
}
