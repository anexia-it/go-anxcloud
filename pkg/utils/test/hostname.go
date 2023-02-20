package test

const (
	hostnameCharset      = "abcdefghijklmnopqrstuvwxyz"
	randomHostnameLength = 8
)

// RandomHostname generates a random hostname, useful for e2e tests and primarily used there.
func RandomHostname() string {
	hostname := make([]byte, randomHostnameLength)
	for i := range hostname {
		hostname[i] = hostnameCharset[getRandom().Intn(len(hostnameCharset))]
	}

	return string(hostname)
}
