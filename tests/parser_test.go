package tests

import (
	core "DNS-server/core"
	models "DNS-server/models"
	"testing"
)

func assertEqual(t *testing.T, got, want models.URL) {
	t.Helper()

	if got.Scheme != want.Scheme {
		t.Errorf("Scheme: got %q, want %q", got.Scheme, want.Scheme)
	}

	if got.Authority.Username != want.Authority.Username {
		t.Errorf("Username: got %q, want %q", got.Authority.Username, want.Authority.Username)
	}

	if got.Authority.Password != want.Authority.Password {
		t.Errorf("Password: got %q, want %q", got.Authority.Password, want.Authority.Password)
	}

	if got.Authority.Host != want.Authority.Host {
		t.Errorf("Host: got %q, want %q", got.Authority.Host, want.Authority.Host)
	}

	if got.Authority.Port != want.Authority.Port {
		t.Errorf("Port: got %q, want %q", got.Authority.Port, want.Authority.Port)
	}

	if got.Path != want.Path {
		t.Errorf("Path: got %q, want %q", got.Path, want.Path)
	}

	if got.Query != want.Query {
		t.Errorf("Query: got %q, want %q", got.Query, want.Query)
	}

	if got.Fragment != want.Fragment {
		t.Errorf("Fragment: got %q, want %q", got.Fragment, want.Fragment)
	}
}

func TestFullURL(t *testing.T) {
	parsedURL := core.ParseURL("https://username:password@api.example.com:8443/v2/users/search?name=john&active=true#results")
	
	expected := models.URL{
		Scheme: "https",
		Authority: models.Authority{
			Username: "username",
			Password: "password",
			Host:     "api.example.com",
			Port:     "8443",
		},
		Path:     "/v2/users/search",
		Query:    "name=john&active=true",
		Fragment: "results",
	}

	assertEqual(t, parsedURL, expected)
}

func TestNoAuthURL(t *testing.T) {
	parsedURL := core.ParseURL("https://api.example.com:8443/v2/users/search?name=john&active=true#results")
	
	expected := models.URL{
		Scheme: "https",
		Authority: models.Authority{
			Username: "",
			Password: "",
			Host:     "api.example.com",
			Port:     "8443",
		},
		Path:     "/v2/users/search",
		Query:    "name=john&active=true",
		Fragment: "results",
	}

	assertEqual(t, parsedURL, expected)
}

func TestDefaultURL(t *testing.T) {
	parsedURL := core.ParseURL("https://api.example.com/v2/users")
	
	expected := models.URL{
		Scheme: "https",
		Authority: models.Authority{
			Username: "",
			Password: "",
			Host:     "api.example.com",
			Port:     "443",
		},
		Path:     "/v2/users",
		Query:    "",
		Fragment: "",
	}

	assertEqual(t, parsedURL, expected)
}

func TestLocalURL(t *testing.T) {
	parsedURL := core.ParseURL("ftp://ftp.example.com/readme.md")
	
	expected := models.URL{
		Scheme: "ftp",
		Authority: models.Authority{
			Username: "",
			Password: "",
			Host:     "ftp.example.com",
			Port:     "21",
		},
		Path:     "/readme.md",
		Query:    "",
		Fragment: "",
	}

	assertEqual(t, parsedURL, expected)
}