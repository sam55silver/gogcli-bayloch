package bayloch

import "fmt"

var serviceToProvider = map[string]string{
	"gmail":     "gmail",
	"calendar":  "google_calendar",
	"drive":     "google_drive",
	"docs":      "google_drive",
	"sheets":    "google_drive",
	"contacts":  "contacts",
	"tasks":     "tasks",
	"people":    "people",
	"chat":      "chat",
	"classroom": "classroom",
	"groups":    "groups",
	"keep":      "keep",
}

// ProviderForService maps a gogcli service name to the Bayloch provider name.
func ProviderForService(service string) (string, error) {
	provider, ok := serviceToProvider[service]
	if !ok {
		return "", fmt.Errorf("unsupported service %q for Bayloch backend", service)
	}
	return provider, nil
}
