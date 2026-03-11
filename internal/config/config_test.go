package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"github.com/topokrat/topoclean/internal/config"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Wir simulieren ein nicht vorhandenes Config-File
	tempDir, _ := os.MkdirTemp("", "topoclean_config_test")
	defer os.RemoveAll(tempDir)
	
	configPath := filepath.Join(tempDir, "nonexistent.json")
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Laden der Standard-Config fehlgeschlagen: %v", err)
	}

	if cfg.HeptagonRoot == "" {
		t.Error("HeptagonRoot sollte standardmäßig gesetzt sein")
	}
	
	// Standardmäßig sollte zumindest das Home-Verzeichnis als Zone existieren (oder leer sein für manuelle Scans)
	if len(cfg.Zones) == 0 {
		t.Log("Standard-Config hat keine vordefinierten Zonen, das ist okay.")
	}
}

func TestLoadJsonConfig(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "topoclean_config_json")
	defer os.RemoveAll(tempDir)
	
	jsonContent := `{
		"heptagon_root": "~/Order",
		"zones": [
			{"path": "~/Downloads", "name": "Downloads"},
			{"path": "~/Desktop", "name": "Desktop"}
		],
		"mapping": {
			"preserve_origin": true
		}
	}`
	
	configPath := filepath.Join(tempDir, "config.json")
	os.WriteFile(configPath, []byte(jsonContent), 0644)
	
	cfg, err := config.Load(configPath)
	if err != nil {
		t.Fatalf("Laden der JSON-Config fehlgeschlagen: %v", err)
	}
	
	if len(cfg.Zones) != 2 {
		t.Errorf("erwartete 2 Zonen, erhalten %d", len(cfg.Zones))
	}
	
	if cfg.Mapping.PreserveOrigin != true {
		t.Error("Mapping.PreserveOrigin sollte true sein")
	}
}
