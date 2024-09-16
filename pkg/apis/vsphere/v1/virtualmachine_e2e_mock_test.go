package v1_test

import (
	"fmt"
	"go.anx.io/go-anxcloud/pkg/apis/vsphere/v1"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const (
	mockVMIdentifier                   = "VMmething1234567aaaaaaaabbbbbbbb"
	mockLocationIdentifier             = "LOCething1234567aaaaaaaabbbbbbbb"
	mockProvisioningLocationIdentifier = "PROVLOCng1234567aaaaaaaabbbbbbbb"
	mockVLANIdentifier                 = "VLANthing1234567aaaaaaaabbbbbbbb"
	mockTemplateIdentifier             = "TMPL1234-ffff-5678-aaaa-aaaaaaaabbbb"
	mockSSHKey                         = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAILHfMcIohDRvMkwBJ4teXJxOTOxEWXbv1gmQlPrWcuiC comment"
	mockIPAddress                      = "255.255.255.255"
)

var (
	mock *Server

	mockStatus = v1.StatusPoweredOn

	mockGetDeleted = false

	mockGetPowerdOn = false
)

func initMockServer() {
	mock = NewServer()
}

/* freshly created:
{
  "identifier": "fa2af2a60bb94df39bbb577fa62d770a",
  "name": "440021-test",
  "custom_name": null,
  "guest_os": null,
  "firmware": null,
  "status": "poweredOff",
  "ram": 1024,
  "cpu": 1,
  "cpu_clock_rate": null,
  "cpu_performance_type": "-",
  "vtpm_enabled": false,
  "cores": 1,
  "disks": 0,
  "disk_info": [],
  "network": [],
  "version_tools": null,
  "guest_tools_status": "Inactive",
  "location_code": "ANX04",
  "location_country": "AT",
  "location_identifier": "52b5f6b2fd3a4a7eaaedf1a7c019e9ea",
  "location_name": "ANX04 - AT, Vienna, Datasix",
  "provisioning_location_identifier": "",
  "template_id": null,
  "resource_salesperson": null
}

// during poweron
{
  "identifier": "fa2af2a60bb94df39bbb577fa62d770a",
  "name": "440021-test",
  "custom_name": null,
  "guest_os": "Ubuntu 24.04 LTS",
  "firmware": "BIOS",
  "status": "poweredOn",
  "ram": 1024,
  "cpu": 1,
  "cpu_clock_rate": 2015,
  "cpu_performance_type": "performance-intel",
  "vtpm_enabled": false,
  "cores": 1,
  "disks": 1,
  "disk_info": [
    {
      "disk_gb": 10,
      "disk_id": 2000,
      "disk_type": "STD4",
      "iops": 900,
      "latence": 30,
      "storage_type": "HDD",
      "bus_type": "VIRTIO",
      "bus_type_label": "VIRTIO(0:0) Hard disk 0"
    }
  ],
  "network": [
    {
      "nic": 5,
      "id": 4000,
      "bandwidth_limit": 1000,
      "vlan": "f4c5c4bb17ca4a81ab0b8ef7d69f3852",
      "ips_v4": [
        "185.x.x.250"
      ],
      "ips_v6": [],
      "mac_address": "00:11:22:33:44:55",
      "connected": true,
      "connected_at_start": true
    }
  ],
  "version_tools": "guestToolsCurrent",
  "guest_tools_status": "active",
  "location_code": "ANX04",
  "location_country": "AT",
  "location_identifier": "52b5f6b2fd3a4a7eaaedf1a7c019e9ea",
  "location_name": "ANX04 - AT, Vienna, Datasix",
  "provisioning_location_identifier": "b164595577114876af7662092da89f76",
  "template_id": null,
  "resource_salesperson": null
}
*/

func mockVMInfoResponseBody(name, desc string) map[string]interface{} {
	// example data taken from https://engine.anexia-it.com/api/vsphere/doc/#!/status/retrieveStatus_get
	return map[string]interface{}{
		"cores":                1,
		"cpu":                  1,
		"cpu_clock_rate":       2095,
		"cpu_performance_type": "performance",
		"custom_name":          desc,
		"disks":                2,
		"disk_info": []map[string]interface{}{
			{
				"bus_type":       "SCSI",
				"bus_type_label": "SCSI(0:0) Hard disk 1",
				"disk_gb":        10,
				"disk_id":        12343567,
				"disk_type":      "ENT6",
				"iops":           900,
				"latence":        30,
				"storage_type":   "HDD",
			},
			{
				"disk_gb": 10,
				"disk_id": 23424323,
			},
		},
		"firmware":            "UEFI",
		"guest_os":            "Debian GNU/Linux 10 (64-bit)",
		"guest_tools_status":  "Active",
		"identifier":          mockVMIdentifier,
		"location_code":       "ANX04",
		"location_country":    "AT",
		"location_identifier": mockLocationIdentifier,
		"location_name":       "ANX04 - AT, Vienna, Datasix",
		"name":                name,
		"network": []map[string]interface{}{
			{
				"nic":             5,
				"bandwidth_limit": 1000,
				"vlan":            mockVLANIdentifier,
				"id":              1235,
				"ips_v4": []string{
					mockIPAddress,
				},
				"ips_v6": []string{
					"0000:111:2222:aaaa:bbbb",
				},
				"mac_address": "00:11:22:33:44:55",
			},
		},
		"provisioning_location_identifier": mockProvisioningLocationIdentifier,
		"ram":                              2048,
		"status":                           mockStatus,
		"template_id":                      mockTemplateIdentifier,
		"version_tools":                    "guestToolsUnmanaged",
		"vtpm_enabled":                     false,
	}
}

func prepareGetInfo(name, desc string) {
	if isIntegrationTest {
		return
	}
	var response http.HandlerFunc

	switch {
	case mockGetDeleted:
		response = RespondWith(404, ``)
	default:
		body := mockVMInfoResponseBody(name, desc)
		response = RespondWithJSONEncoded(200, body)

		if mockGetPowerdOn {
			body["status"] = v1.StatusPoweredOn
		}
	}

	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/vsphere/v1/info.json/"+mockVMIdentifier+"/info"),
		response,
	))
}

func prepareCreate(name, desc string) {
	if isIntegrationTest {
		return
	}

	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("POST", "/api/vsphere/v1/provisioning/vm.json/"+mockLocationIdentifier+"/"+string(v1.TypeTemplate)+"/"+mockTemplateIdentifier),
		VerifyJSONRepresenting(map[string]interface{}{
			"additional_disks": []map[string]interface{}{
				{"gb": 10, "type": "STD4"},
			},
			"cpus":                 2,
			"cpu_performance_type": "performance-amd",
			"custom_name":          desc,
			"disk_gb":              10,
			"disk_type":            "ENT6",
			"hostname":             name,
			"memory_mb":            2048,
			"network": []map[string]interface{}{
				{"nic_type": "vmxnet3", "vlan": mockVLANIdentifier, "bandwidth_limit": 1000, "ips": []string{mockIPAddress}},
			},
			"script": "Iy9iaW4vc2gK",
			"ssh":    mockSSHKey,
		}),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"identifier": mockProgressIdentifier,
			"progress":   10,
			"queued":     false,
			"errors":     []string{},
		}),
	))
}

func prepareList(name, desc string) {
	if isIntegrationTest {
		return
	}

	mockVMs := []map[string]interface{}{
		{"identifier": "foo0", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo1", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo2", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo3", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo4", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo5", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo6", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo7", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo8", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		{"identifier": "foo9", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": mockLocationIdentifier},
		mockVMInfoResponseBody(name, desc),
	}

	pages := [][]map[string]interface{}{
		mockVMs[0:10],
		mockVMs[10:],
		{},
	}

	Expect(pages[0]).To(HaveLen(10))
	Expect(pages[1]).To(HaveLen(1))
	Expect(pages[len(pages)-1]).To(BeEmpty())

	for i, data := range pages {
		mock.AppendHandlers(CombineHandlers(
			VerifyRequest("GET", "/api/vsphere/v1/vmlist/list.json", fmt.Sprintf("page=%v&limit=10", i+1)),
			RespondWithJSONEncoded(200, map[string]interface{}{
				"page":        i + 1,
				"total_pages": len(pages),
				"total_items": len(mockVMs),
				"limit":       len(data),
				"data":        data,
			}),
		))
	}
}

func prepareDelete() {
	if isIntegrationTest {
		return
	}
	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("DELETE", "/api/vsphere/v1/provisioning/vm.json/"+mockVMIdentifier),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"identifier":                 mockVMIdentifier,
			"delete_will_be_executed_at": time.Now().Format(time.DateTime),
		}),
	))
}

func prepareEventuallyActive(name, desc string) {
	prepareGetInfo(name, desc)
	mockGetPowerdOn = true // Powered-on on second request.
}
