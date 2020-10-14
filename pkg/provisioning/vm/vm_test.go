package vm_test

import (
	"context"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/provisioning/ips"
	"github.com/anexia-it/go-anxcloud/pkg/provisioning/progress"
	"github.com/anexia-it/go-anxcloud/pkg/provisioning/vm"
	"golang.org/x/crypto/ssh"
)

const (
	hostnameCharset = "abcdefghijklmnopqrstuvwxyz"
	templateType    = "templates"
	templateID      = "44b38284-6adb-430e-b4a4-1553e29f352f"
	cpus            = 2
	memory          = 2048
	disk            = 10
)

var (
	location        = ""
	vlan            = ""
	skipIntegration = true
)

func init() {
	var set bool
	if _, set = os.LookupEnv(client.IntegrationTestEnvName); !set {
		return
	}
	skipIntegration = false
	if location, set = os.LookupEnv(client.LocationEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.LocationEnvName))
	}
	if vlan, set = os.LookupEnv(client.VLANEnvName); !set {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VLANEnvName))
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

func TestVMProvisioningDeprovisioningIntegration(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.NewAnyClientFromEnvs(false, nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	ips, err := ips.GetFree(ctx, location, vlan, c)
	defer cancel()
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}
	if len(ips) < 1 {
		t.Fatalf("no IPs left for testing in vlan")
	}

	networkInterfaces := []vm.Network{{
		NICType: "vmxnet3",
		IPs:     []string{ips[0].Identifier},
		VLAN:    vlan,
	}}

	definition := vm.NewDefinition(location, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
	definition.SSH = randomPublicSSHKey()

	provisionResponse, err := vm.Provision(ctx, definition, c)
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}

	vmID, err := progress.AwaitCompletion(ctx, provisionResponse.Identifier, c)
	if err != nil {
		t.Fatalf("waiting for VM provisioning failed: %v", err)
	}

	if err = vm.Deprovision(ctx, vmID, false, c); err != nil {
		t.Fatalf(fmt.Sprintf("could not deprovision VM: %v", err))
	}
}
