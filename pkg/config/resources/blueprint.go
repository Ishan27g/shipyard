package resources

import (
	"fmt"
	"net/url"
)

// Blueprint defines a stack blueprint for defining yard configs
type Blueprint struct {
	Title              string            `hcl:"title,optional" json:"title,omitempty"`
	Author             string            `hcl:"author,optional" json:"author,omitempty"`
	Slug               string            `hcl:"slug,optional" json:"slug,omitempty"`
	Intro              string            `hcl:"intro,optional" json:"intro,omitempty"`
	BrowserWindows     []string          `hcl:"browser_windows,optional" json:"browser_windows,omitempty"`
	HealthCheckTimeout string            `hcl:"health_check_timeout,optional" json:"health_check_timeout,omitempty"`
	Environment        map[string]string `hcl:"env,optional" json:"environment,omitempty"`
	ShipyardVersion    string            `hcl:"shipyard_version,optional" json:"shipyard_version,omitempty"`
}

// Validate the Blueprint and return errors
func (b *Blueprint) Validate() []error {
	errors := make([]error, 0)
	// ensure BrowserWindows are valid URIs
	for _, i := range b.BrowserWindows {
		uri, err := url.Parse(i)
		if err != nil {
			errors = append(
				errors,
				fmt.Errorf("invalid BrowserWindow URI: %s, %s", i, err),
			)
		}

		if uri.String() == "" {
			errors = append(
				errors,
				fmt.Errorf("invalid BrowserWindow URI, uri is empty: %s", i),
			)
		}
	}

	return errors
}
