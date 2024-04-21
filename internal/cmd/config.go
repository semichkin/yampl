package cmd

type Config struct {
	Template string         `json:"template"`
	Params   map[string]any `json:"params"`
}
