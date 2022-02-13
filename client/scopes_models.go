package client

type ApiScopeCreate struct {
	DisplayName             string            `json:"displayName"`
	Description             string            `json:"description"`
	ShowInDiscoveryDocument bool              `json:"showInDiscoveryDocument"`
	UserClaims              []string          `json:"userClaims"`
	Properties              map[string]string `json:"properties"`
	Enabled                 bool              `json:"enabled"`
	Required                bool              `json:"required"`
	Emphasize               bool              `json:"emphasize"`
}

type ApiScope struct {
	Name                    string            `json:"name"`
	DisplayName             string            `json:"displayName"`
	Description             string            `json:"description"`
	ShowInDiscoveryDocument bool              `json:"showInDiscoveryDocument"`
	UserClaims              []string          `json:"userClaims"`
	Properties              map[string]string `json:"properties"`
	Enabled                 bool              `json:"enabled"`
	Required                bool              `json:"required"`
	Emphasize               bool              `json:"emphasize"`
}
