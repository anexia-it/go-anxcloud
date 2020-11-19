package provisioning_test

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
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/powercontrol"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/ips"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/progress"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/vm"
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
	if location, set = os.LookupEnv(client.VsphereLocationEnvName); !set || location == "" {
		panic(fmt.Sprintf("could not find environment variable %s, which is required for testing", client.VsphereLocationEnvName))
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

func TestVMProvisioningDeprovisioningIntegration(t *testing.T) { //nolint:funlen // Flake prevention needs space.
	t.Parallel()

	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	ips, err := ips.NewAPI(c).GetFree(ctx, location, vlan)
	defer cancel()
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}
	if len(ips) < 1 {
		t.Fatalf("no IPs left for testing in vlan")
	}

	networkInterfaces := []vm.Network{{NICType: "vmxnet3", IPs: []string{ips[0].Identifier}, VLAN: vlan}}
	definition := vm.NewAPI(c).NewDefinition(location, templateType, templateID, randomHostname(), cpus, memory, disk, networkInterfaces)
	definition.SSH = randomPublicSSHKey()

	provisionResponse, err := vm.NewAPI(c).Provision(ctx, definition)
	if err != nil {
		t.Fatalf("provisioning vm failed: %v", err)
	}

	vmID, err := progress.NewAPI(c).AwaitCompletion(ctx, provisionResponse.Identifier)
	if err != nil {
		t.Fatalf("waiting for VM provisioning failed: %v", err)
	}

	time.Sleep(10 * time.Second)

	task, err := powercontrol.NewAPI(c).Set(ctx, vmID, powercontrol.RebootRequest)
	if err != nil {
		t.Fatalf(fmt.Sprintf("could not set power: %v", err))
	}

	if err = powercontrol.NewAPI(c).AwaitCompletion(ctx, vmID, task.TaskIdentifier); err != nil {
		t.Fatalf("could not await power state change of created VM: %v", err)
	}

	if _, err = powercontrol.NewAPI(c).Get(ctx, vmID); err != nil {
		t.Fatalf("could not get power state of created VM: %v", err)
	}

	if err = vm.NewAPI(c).Deprovision(ctx, vmID, false); err != nil {
		t.Fatalf(fmt.Sprintf("could not deprovision VM: %v", err))
	}
}
