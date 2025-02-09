package plugin

// NotFoundMsg contains the name of the plugin that is not found
type NotFoundMsg struct {
	PluginName string
	Source     string
}

// VersionMsg contains the version that was found
type VersionMsg string

type Metadata struct {
	Name      string `json:"name" validate:"required"`
	Version   string `json:"version" validate:"required,semver"`
	Downloads []struct {
		OS       string `json:"os" validate:"required"`
		Arch     string `json:"arch" validate:"required"`
		URL      string `json:"url" validate:"required,http_url"`
		Checksum string `json:"checksum"`
	} `json:"downloads" validate:"required,dive"`
}

type Entry struct {
	Name        string `json:"name" validate:"required"`
	Version     string `json:"version" validate:"required,semver"`
	MetadataUrl string `json:"metadataUrl" validate:"required,http_url"`
}
