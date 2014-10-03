package config

type RepositoryLayout struct {
	Distributions []string `json:"distributions"`
	Categories    []string `json:"categories"`
	Architectures []string `json:"architectures"`
}

func (this RepositoryLayout) validate() error {
	return nil
}
