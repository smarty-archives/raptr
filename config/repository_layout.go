package config

type RepositoryLayout struct {
	LayoutKey     string   `json:-`
	Distributions []string `json:"distributions"`
	Categories    []string `json:"categories"`
	Architectures []string `json:"architectures"`
}

func (this RepositoryLayout) Validate() error {
	return nil
}
