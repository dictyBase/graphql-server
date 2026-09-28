package resolver

// errParseIDFormat is the shared message for identifier parsing failures.
const errParseIDFormat = "error in parsing string %s to int %s"

// GraphQL type discriminators carried on the JSON:API style payloads.
const (
	typePermission = "permission"
	typeRole       = "role"
	typeUser       = "user"
)
