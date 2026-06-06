package oauth

import "errors"

var (
	InvalidResponseCodeError      = errors.New("invalid_response_type")
	InvalidClientError            = errors.New("invalid_client")
	InvalidClientIDError          = errors.New("invalid_cliend_id")
	InvalidRedirectURIError       = errors.New("invalid_redirect_uri")
	InvalidScopeError             = errors.New("invalid_scope")
	InvalidStateError             = errors.New("invalid_state")
	InvalidCodeChallenge          = errors.New("invalid_code_challenge")
	InvalidGrantTypeError         = errors.New("invalid_grant_type")
	MissingClientCredentialsError = errors.New("missing_client_credential")
	ClientIDMisMatchError         = errors.New("client_id_mismatch")
	NoClientSecretConfiguredError = errors.New("no_secret_configured_error")
	InvalidPKCEError              = errors.New("invalid_pkce")
)
