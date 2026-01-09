package core

import (
	models "DNS-server/models"
	"strings"
)

func ParseURL(url string) models.URL {
	var parsedURL models.URL
	remaining := url

	// Step 1: Extract Scheme
	if idx := strings.Index(remaining, "://"); idx != -1 {
		parsedURL.Scheme = remaining[:idx]
		remaining = remaining[idx+3:] // Skip "://"
	} else {
		// assume https
		parsedURL.Scheme = "https"
	}

	// Step 2: Extract Fragment
	if idx := strings.Index(remaining, "#"); idx != -1 {
		parsedURL.Fragment = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Step 3: Extract Query String
	if idx := strings.Index(remaining, "?"); idx != -1 {
		parsedURL.Query = remaining[idx+1:]
		remaining = remaining[:idx]
	}

	// Step 4: Extract Path
	if idx := strings.Index(remaining, "/"); idx != -1 {
		parsedURL.Path = remaining[idx:] // Include the leading /
		remaining = remaining[:idx]
	} else {
		// assume no path
		parsedURL.Path = "/"
	}

	// Step 5: Extract Host and Port
	hostPort := remaining
	
	// Check for IPv6 (enclosed in brackets)
	if strings.HasPrefix(hostPort, "[") {
		if idx := strings.Index(hostPort, "]"); idx != -1 {
			parsedURL.Authority.Host = hostPort[1:idx] // Remove brackets
			if len(hostPort) > idx+1 && hostPort[idx+1] == ':' {
				parsedURL.Authority.Port = hostPort[idx+2:]
			}
		}
	} else {
		// IPv4 or domain name
		if idx := strings.LastIndex(hostPort, ":"); idx != -1 {
			parsedURL.Authority.Host = hostPort[:idx]
			parsedURL.Authority.Port = hostPort[idx+1:]
		} else {
			parsedURL.Authority.Host = hostPort
			// Set default port based on scheme
			if parsedURL.Scheme == "https" {
				parsedURL.Authority.Port = "443"
			} else if parsedURL.Scheme == "http" {
				parsedURL.Authority.Port = "80"
			} else if parsedURL.Scheme == "ftp" {
				parsedURL.Authority.Port = "21"
			}
		}
	}

	return parsedURL
}