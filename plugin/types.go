package plugin

type PluginMetadata struct {
	Name      string `json:"name" validate:"required"`
	Version   string `json:"version" validate:"required,semver"`
	Downloads []struct {
		OS       string `json:"os" validate:"required"`
		Arch     string `json:"arch" validate:"required"`
		URL      string `json:"url" validate:"required,http_url"`
		Checksum string `json:"checksum"`
	} `json:"downloads" validate:"required,dive"`
}
