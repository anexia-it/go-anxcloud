package v1

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	v1 "go.anx.io/go-anxcloud/pkg/apis/core/v1"

	apiTypes "go.anx.io/go-anxcloud/pkg/api/types"
)

func (v *VirtualMachine) EndpointURL(ctx context.Context) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch op {
	case apiTypes.OperationCreate:
		return url.Parse("/api/vsphere/v1/provisioning/vm.json/" + v.Location.Identifier + "/templates/" + v.TemplateID)
	case apiTypes.OperationDestroy:
		return url.Parse("/api/vsphere/v1/provisioning/vm.json")
	case apiTypes.OperationList:
		query := url.Values{}
		u, _ := url.Parse("/api/vsphere/v1/vmlist/list.json")
		u.RawQuery = query.Encode()
		return u, nil
	}

	return url.Parse("/api/vsphere/v1/info.json")
}

func (v *VirtualMachine) FilterRequestURL(ctx context.Context, url *url.URL) (*url.URL, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == apiTypes.OperationGet {
		url.Path = path.Join(url.Path, "/info")
	}
	return url, nil
}

func (v *VirtualMachine) FilterAPIRequestBody(ctx context.Context) (interface{}, error) {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if op == apiTypes.OperationCreate {
		return createVMBody(v)
	}
	if op == apiTypes.OperationGet {
		return nil, nil
	}

	return v, nil
}

func createVMBody(vm *VirtualMachine) (interface{}, error) {
	// some checks upfront
	if len(vm.DiskInfo) == 0 {
		return nil, errors.New("at least one disk with size >= 10 must be specified")
	}
	if vm.Cores > 16 {
		return nil, errors.New("at most 16 cores must be specified")
	}

	// Creating a Virtual Machine is done with only a few and different fields.
	data := struct {
		Hostname           string                   `json:"hostname"`
		AdditionalDisks    []map[string]interface{} `json:"additional_disks,omitempty"`
		CPUs               int                      `json:"cpus"`
		CPUPerformanceType *string                  `json:"cpu_performance_type,omitempty"`
		CustomName         *string                  `json:"custom_name,omitempty"`
		DiskGB             int                      `json:"disk_gb"`
		DiskType           *string                  `json:"disk_type,omitempty"`
		MemoryMB           int                      `json:"memory_mb"`
		Networks           []map[string]interface{} `json:"network,omitempty"`
		Script             *string                  `json:"script,omitempty"`
		Sockets            *int                     `json:"sockets,omitempty"`
		SSH                *string                  `json:"ssh,omitempty"`
		Password           *string                  `json:"password,omitempty"`
	}{
		Hostname: vm.Name,
		MemoryMB: vm.RAM,
		CPUs:     vm.Cores, // upon creation, we can only provide cores.
		DiskGB:   vm.DiskInfo[0].DiskGB,
	}

	if vm.CPU != 0 {
		data.Sockets = &vm.CPU
	}

	if vm.CPUPerformanceType != "" {
		data.CPUPerformanceType = &vm.CPUPerformanceType
	}

	if vm.CustomName != "" {
		data.CustomName = &vm.CustomName
	}

	if vm.DiskInfo[0].DiskType != "" {
		data.DiskType = &vm.DiskInfo[0].DiskType
	}

	for _, disk := range vm.DiskInfo[1:] {
		data.AdditionalDisks = append(data.AdditionalDisks, map[string]interface{}{
			"gb":   disk.DiskGB,
			"type": disk.DiskType,
		})
	}

	if vm.Networks != nil {
		for _, net := range vm.Networks {
			data.Networks = append(data.Networks, map[string]interface{}{
				"nic_type":        net.NICType,
				"vlan":            net.VLAN,
				"bandwidth_limit": net.BandwidthLimit,
				"ips":             net.IPs,
			})
		}
	}

	if vm.StartScript != "" {
		script := base64.StdEncoding.EncodeToString([]byte(vm.StartScript))
		data.Script = &script
	}

	switch {
	case vm.SSHKey != "":
		data.SSH = &vm.SSHKey
	case vm.Password != "":
		data.Password = &vm.Password
	default:
		return nil, errors.New("missing required field: 'password' or 'ssh'")
	}

	// TODO: remove
	fmt.Printf("request: %+v\n", data)

	return data, nil
}

func (v *VirtualMachine) DecodeAPIResponse(ctx context.Context, data io.Reader) error {
	op, err := apiTypes.OperationFromContext(ctx)
	if err != nil {
		return err
	}

	// POST request returns provision progress information.
	if op == apiTypes.OperationCreate {
		resp := ProvisionProgress{}
		err := json.NewDecoder(data).Decode(&resp)
		if err != nil {
			return err
		}
		v.Progress = resp

		return nil
	}

	// TODO:
	// GET {id}/info
	// GET {id}/status
	// PUT * ??
	err = json.NewDecoder(data).Decode(v)
	if err != nil {
		return err
	}

	if op == apiTypes.OperationGet {
		loc := v1.Location{
			Code:        v.LocationCode,
			CountryCode: v.LocationCountry,
			Identifier:  v.LocationIdentifier,
			Name:        v.LocationName,
		}
		v.Location = loc
		v.LocationCode = ""
		v.LocationCountry = ""
		v.LocationIdentifier = ""
		v.LocationName = ""
	}

	return nil
}

func (v *VirtualMachine) FilterAPIResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	// TODO: remove
	buf := bytes.NewBuffer(nil)
	tee := io.TeeReader(res.Body, buf)
	body, _ := io.ReadAll(tee)
	fmt.Printf("response: %+v\n", string(body))
	res.Body = io.NopCloser(buf)
	return res, nil
}
