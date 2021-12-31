// Package domain contains all the specific objects related to the application
package domain

import (
	"net/url"
)

type Repo struct {
	Remote *url.URL
	Local  string
	Server string
	Owner  string
	Name   string
}
