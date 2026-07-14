package config

import "testing"

func TestLoadMailConfiguration(t *testing.T) {
	t.Setenv("MAIL_HOST", "mail.tino.vn")
	t.Setenv("MAIL_PORT", "587")
	t.Setenv("MAIL_USER", "admin@codelife138.io.vn")
	t.Setenv("MAIL_PASS", "secret")
	t.Setenv("MAIL_FROM", "VNDoctor <admin@codelife138.io.vn>")
	t.Setenv("MAIL_SECURE", "false")

	config := Load()
	if config.MailHost != "mail.tino.vn" || config.MailPort != "587" {
		t.Fatal("unexpected mail server configuration")
	}
	if config.MailUser != "admin@codelife138.io.vn" || config.MailPass != "secret" {
		t.Fatal("mail credentials were not loaded")
	}
	if config.MailFrom != "VNDoctor <admin@codelife138.io.vn>" || config.MailSecure {
		t.Fatal("unexpected sender configuration")
	}
}

func TestLoadSupportsLegacySMTPNames(t *testing.T) {
	t.Setenv("MAIL_HOST", "")
	t.Setenv("MAIL_PORT", "")
	t.Setenv("MAIL_USER", "")
	t.Setenv("MAIL_PASS", "")
	t.Setenv("MAIL_FROM", "")
	t.Setenv("SMTP_HOST", "legacy.example.com")
	t.Setenv("SMTP_PORT", "2525")
	t.Setenv("SMTP_USERNAME", "legacy-user")
	t.Setenv("SMTP_PASSWORD", "legacy-pass")
	t.Setenv("SMTP_FROM", "legacy@example.com")

	config := Load()
	if config.MailHost != "legacy.example.com" || config.MailPort != "2525" || config.MailUser != "legacy-user" {
		t.Fatal("legacy SMTP configuration was not loaded")
	}
}
