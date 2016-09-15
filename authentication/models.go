package authentication

// loginAuth is an smtp.Auth that implements the LOGIN authentication mechanism.
type loginAuth struct {
	username string
	password string
	host     string
}
