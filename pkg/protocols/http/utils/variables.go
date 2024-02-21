package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"Ni/pkg/protocols/common/contextargs"
	"Ni/pkg/protocols/common/generators"
	"Ni/pkg/protocols/dns"
)

// GenerateVariables will create default variables with context args
func GenerateVariablesWithContextArgs(input *contextargs.Context, trailingSlash bool) map[string]interface{} {
	parsed, err := url.Parse(input.MetaInput.Input)
	if err != nil {
		return nil
	}
	return GenerateVariablesWithURL(parsed, trailingSlash, contextargs.GenerateVariables(input))
}

// GenerateVariables will create default variables after parsing a url with additional variables
func GenerateVariablesWithURL(parsed *url.URL, trailingSlash bool, additionalVars map[string]interface{}) map[string]interface{} {
	domain := parsed.Host
	if strings.Contains(parsed.Host, ":") {
		domain = strings.Split(parsed.Host, ":")[0]
	}

	port := parsed.Port()
	if port == "" {
		if parsed.Scheme == "https" {
			port = "443"
		} else if parsed.Scheme == "http" {
			port = "80"
		}
	}

	if trailingSlash {
		parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	}

	escapedPath := parsed.EscapedPath()
	directory := path.Dir(escapedPath)
	if directory == "." {
		directory = ""
	}
	base := path.Base(escapedPath)
	if base == "." {
		base = ""
	}
	httpVariables := map[string]interface{}{
		"BaseURL":  parsed.String(),
		"RootURL":  fmt.Sprintf("%s://%s", parsed.Scheme, parsed.Host),
		"Hostname": parsed.Host,
		"Host":     domain,
		"Port":     port,
		"Path":     directory,
		"File":     base,
		"Scheme":   parsed.Scheme,
	}
	return generators.MergeMaps(httpVariables, dns.GenerateVariables(domain), additionalVars)
}
