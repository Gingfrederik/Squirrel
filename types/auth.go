package types

type RoleUser struct {
	Role  string   `json:"role"`
	Users []string `json:"users"`
}

type RoleAllUser struct {
	Role  string  `json:"role"`
	Users []*User `json:"users"`
}

type Policy struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}
