package cuckooc

import "testing"

func TestConfig_loadConfig(t *testing.T) {
	file := "./testdata/config_example.json"
	config, err := loadConfig(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if config.BackupFolder != "./testdata/backups" {
		t.Fatalf("expected %s but got %s", "./testdata/backups", config.BackupFolder)
	}
}
