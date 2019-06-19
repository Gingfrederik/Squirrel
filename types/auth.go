package types

type RoleUser struct {
	Role  string
	Users []string
}

type Policy struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}
