package cmd

import "github.com/vvvkkkggg/kubeconomist-core/internal/config"

func Run() error {
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		return err
	}

	_ = cfg

	panic("implement me")

	return nil
}
