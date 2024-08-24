package auth

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
	"github.com/npushpakumara/go-backend-template/internal/config"
)

// NewOAuthProviders initializes and registers the OAuth providers using the Goth library.
// It accepts a configuration object that contains the necessary credentials and settings for each OAuth provider.
func NewOAuthProviders(cfg *config.Config) {
	// goth.UseProviders registers the OAuth providers that Goth will use for authentication.
	// In this case, we are setting up Google and Microsoft as the providers.

	goth.UseProviders(
		google.New(
			cfg.OAuth.Google.ClientID,
			cfg.OAuth.Google.ClientSecret,
			cfg.OAuth.Google.RedirectURL,
			cfg.OAuth.Google.GetScopes()...,
		),
		microsoftonline.New(
			cfg.OAuth.Microsoft.ClientID,
			cfg.OAuth.Microsoft.ClientSecret,
			cfg.OAuth.Microsoft.RedirectURL,
			cfg.OAuth.Microsoft.GetScopes()...,
		),
	)
}
