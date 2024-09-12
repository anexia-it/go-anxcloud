package v1

import (
	"context"
	"net/url"
	"path"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
)

func (c *VirtualMachine) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch op {
	case apiTypes.OperationList:
		query := url.Values{}
		u, _ := url.Parse("/api/vsphere/v1/vmlist/list.json")
		u.RawQuery = query.Encode()
		return u, nil
	case apiTypes.OperationCreate:
		return url.Parse("/api/provisioning/v1/vm.json/" + c.Location.Identifier + "/templates/" + c.TemplateID)
	}

	return url.Parse("/api/vsphere/v1/info.json")
}

func (c *VirtualMachine) FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == apiTypes.OperationGet {
		url.Path = path.Join(url.Path, "/info")
	}
	return url, nil
}

func (c *VirtualMachine) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Creating a Virtual Machine is done with only a few fields.
	if op == apiTypes.OperationCreate {
		// TODO: support multiple disks via `additional_disks`
		data := struct {
			Hostname   string  `json:"hostname"`
			CustomName string  `json:"custom_name"`
			MemoryMB   int     `json:"memory_mb"`
			CPUs       int     `json:"cpus"`
			DiskGB     float64 `json:"disk_gb"`
		}{
			Hostname:   c.Name,
			CustomName: c.CustomName,
			MemoryMB:   c.RAM,
			CPUs:       c.CPU,
			DiskGB:     c.DiskInfo[0].DiskGB,
		}

		return data, nil
	}
	if op == apiTypes.OperationGet {

		return nil, nil
	}

	return c, nil
}
