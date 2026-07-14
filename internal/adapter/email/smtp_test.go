package email

import (
	"strings"
	"testing"
)

func TestNewSMTPMailerSeparatesHeaderAndEnvelopeSender(t *testing.T) {
	mailer, err := NewSMTPMailer(
		"mail.tino.vn",
		"587",
		"admin@codelife138.io.vn",
		"secret",
		"VNDoctor <admin@codelife138.io.vn>",
		false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if mailer.fromAddress != "admin@codelife138.io.vn" {
		t.Fatalf("unexpected envelope sender: %s", mailer.fromAddress)
	}
	if !strings.Contains(mailer.fromHeader, "VNDoctor") {
		t.Fatalf("display name missing from header: %s", mailer.fromHeader)
	}
	message := string(mailer.otpMessage("student@example.com", "123456"))
	if !strings.Contains(message, "VNDoctor") || !strings.Contains(message, "<admin@codelife138.io.vn>") {
		t.Fatal("unexpected From header")
	}
}

func TestNewSMTPMailerRejectsInvalidConfiguration(t *testing.T) {
	if _, err := NewSMTPMailer("", "587", "", "", "admin@example.com", false); err == nil {
		t.Fatal("expected missing host error")
	}
	if _, err := NewSMTPMailer("mail.example.com", "587", "", "", "invalid", false); err == nil {
		t.Fatal("expected invalid sender error")
	}
}
