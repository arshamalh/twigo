package twigo

type OAuthType int8

const (
	OAuth_Default OAuthType = 0
	OAuth_1a      OAuthType = 1
	OAuth_2       OAuthType = 2
)

// TODO: Should be completed, for example,
// Search tweets and Search spaces or muting have OAuth 1 as default,
// but others, have version 2 as default
var DefaultOAuthes = map[string]OAuthType{
	"": OAuth_Default,
}
