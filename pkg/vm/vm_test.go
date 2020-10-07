package vm_test

import (
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

const (
	hostnameCharset = "abcdefghijklmnopqrstuvwxyz"
	templateType    = "templates"
	templateID      = "44b38284-6adb-430e-b4a4-1553e29f352f"
	cpus            = 2
	memory          = 2048
	disk            = 10
	locationEnvName = "ANXCLOUD_LOCATION"
	vlanEnvName     = "ANXCLOUD_VLAN"
)

var (
	location = ""
	vlan     = ""
)

func init() {
	var set bool
	if location, set = os.LookupEnv(locationEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", locationEnvName))
	}
	if vlan, set = os.LookupEnv(vlanEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", vlanEnvName))
	}
}

func randomPublicSSHKey() string {
	private, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("could not create ssh private key: %v", err))
	}
	public, err := ssh.NewPublicKey(&private.PublicKey)
	if err != nil {
		panic(fmt.Sprintf("could not create ssh public key: %v", err))
	}

	return string(ssh.MarshalAuthorizedKey(public))
}

func randomHostname() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // No crypto needed here.
	hostnameSuffix := make([]string, 8)
	for i := range hostnameSuffix {
		hostnameSuffix[i] = string(hostnameCharset[r.Intn(len(hostnameCharset))])
	}

	return fmt.Sprintf("go-test-%s", strings.Join(hostnameSuffix, ""))
}
