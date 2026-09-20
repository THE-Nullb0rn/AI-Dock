package model

type Variant struct {
	Label      string `json:"label,omitempty"`
	Type       string `json:"type"`
	Path       string `json:"path"`
	LaunchCmd  string `json:"launch_cmd"`
	LaunchMode string `json:"launch_mode"`
}

type Tool struct {
	Name     string    `json:"name"`
	Variants []Variant `json:"variants"`
}
