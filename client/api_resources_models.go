package client

type ApiResourceCreate struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Scopes      []string `json:"scopes"`
}

type ApiResourceUpdate struct {
	DisplayName string   `json:"displayName,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
}

type ApiResource struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	Scopes      []string `json:"scopes"`
}
