package database

import (
	"testing"
	"testing/fstest"
)

func TestLoadMigrationsSortsByTimestamp(t *testing.T) {
	files := fstest.MapFS{
		"20260714120000_second.sql": {Data: []byte("SELECT 2;")},
		"20260714110000_first.sql":  {Data: []byte("SELECT 1;")},
		"ignored.txt":               {Data: []byte("ignored")},
	}

	items, err := loadMigrations(files)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].name != "20260714110000_first.sql" || items[1].name != "20260714120000_second.sql" {
		t.Fatalf("unexpected migration order: %#v", items)
	}
	if items[0].checksum == "" || items[0].checksum == items[1].checksum {
		t.Fatalf("unexpected migration checksums: %#v", items)
	}
}

func TestLoadMigrationsRejectsInvalidNames(t *testing.T) {
	files := fstest.MapFS{
		"004_progress.sql": {Data: []byte("SELECT 1;")},
	}
	if _, err := loadMigrations(files); err == nil {
		t.Fatal("expected invalid migration name error")
	}
}

func TestLoadMigrationsRejectsDuplicateTimestamps(t *testing.T) {
	files := fstest.MapFS{
		"20260714110000_first.sql": {Data: []byte("SELECT 1;")},
		"20260714110000_other.sql": {Data: []byte("SELECT 2;")},
	}
	if _, err := loadMigrations(files); err == nil {
		t.Fatal("expected duplicate timestamp error")
	}
}
