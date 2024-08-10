package config

type Plugin struct {
	Name string `json:"name,omitempty"`
	Cmd  string `json:"cmd,omitempty"`
}
type Config struct {
	Version   string   `json:"version,omitempty"`
	ConfigDir *string  `json:"configDir,omitempty"`
	Plugins   []Plugin `json:"plugins,omitempty"`
}
