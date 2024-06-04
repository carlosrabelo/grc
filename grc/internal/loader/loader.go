package loader

import (
	"fmt"
	"os"

	"github.com/carlosrabelo/grc/grc/internal/validate"
	"github.com/carlosrabelo/grc/grc/pkg/types"
)

func LoadConfig(path string) (*types.FiltersConfig, error) {
	if err := validate.Extension(path); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	cfg, err := validate.YAML(data)
	if err != nil {
		return nil, err
	}
	if err := validate.Config(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
