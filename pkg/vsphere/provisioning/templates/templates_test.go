package templates_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/anexia-it/go-anxcloud/pkg/client"
	"github.com/anexia-it/go-anxcloud/pkg/vsphere/provisioning/templates"
)

var (
	location        = ""
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
}

func TestListLocations(t *testing.T) {
	if skipIntegration {
		t.Skip("integration tests disabled")
	}
	c, err := client.New(client.AuthFromEnv(false))
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), client.DefaultRequestTimeout)
	defer cancel()

	_, err = templates.NewAPI(c).List(ctx, location, templates.TemplateTypeTemplates, 1, 50)
	if err != nil {
		t.Errorf("could not get templates: %v", err)
	}
}
