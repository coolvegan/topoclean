package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Version      string  `json:"version"`
	HeptagonRoot string  `json:"heptagon_root"`
	Zones        []Zone  `json:"zones"`
	Mapping      Mapping `json:"mapping"`
}

type Zone struct {
	Path     string `json:"path"`
	Name     string `json:"name"`
	Strategy string `json:"strategy"`
	KeepDays int    `json:"keep_days"`
}

type Mapping struct {
	PreserveOrigin bool   `json:"preserve_origin"`
	DateFormat     string `json:"date_format"`
}

func Load(path string) (*Config, error) {
	// 1. Standardwerte setzen
	home, _ := os.UserHomeDir()
	cfg := &Config{
		Version:      "1.0",
		HeptagonRoot: home,
		Mapping: Mapping{
			PreserveOrigin: true,
			DateFormat:     "2006-01",
		},
	}

	// 2. JSON laden, falls Datei existiert
	if _, err := os.Stat(path); err == nil {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, err
		}
	}

	// 3. Pfade auflösen (Tilde-Expansion)
	cfg.HeptagonRoot = expandPath(cfg.HeptagonRoot)
	for i := range cfg.Zones {
		cfg.Zones[i].Path = expandPath(cfg.Zones[i].Path)
	}

	return cfg, nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}
	return path
}
