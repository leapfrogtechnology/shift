package http

import "github.com/go-resty/resty/v2"

// Client is the HTTP client used throughout the CLI
var Client = resty.New()
