package auth

type AuthError struct {
	reason string
}

func (e AuthError) Error() string {
	return e.reason
}

type ConfigurationError struct {
	reason string
}

func (e ConfigurationError) Error() string {
	return e.reason
}
