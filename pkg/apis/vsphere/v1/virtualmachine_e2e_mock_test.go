package v1_test

import (
	"fmt"
	"go.anx.io/go-anxcloud/pkg/apis/vsphere/v1"
	"net/http"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/ghttp"
)

const (
	mockADCIdentifier  = "something1234567aaaaaaaabbbbbbbb"
	mockVLANIdentifier = "something1234567aaaaaaaabbbbbbbb"
)

var (
	mock *Server

	mockStatus = v1.StatusPoweredOn

	mockGetDeleted = false

	templateType = "templates"

	templateIdentifier = "6a681fb4-0bc9-4b68-bada-b76a8ec31e0f"
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

func mockADCInfoResponseBody(name, desc string) map[string]interface{} {

	// example data taken from https://engine.anexia-it.com/api/vsphere/doc/#!/status/retrieveStatus_get
	return map[string]interface{}{
		"identifier":           mockADCIdentifier,
		"name":                 name,
		"custom_name":          desc,
		"guest_os":             "Debian GNU/Linux 10 (64-bit)",
		"firmware":             "UEFI",
		"status":               mockStatus,
		"ram":                  1024,
		"cpu":                  1,
		"cpu_clock_rate":       2095,
		"cpu_performance_type": "performance",
		"vtpm_enabled":         false,
		"cores":                1,
		"disks":                1,
		"disk_info": []map[string]interface{}{
			{
				"bus_type":       "SCSI",
				"bus_type_label": "SCSI(0:0) Hard disk 1",
				"disk_gb":        4,
				"disk_id":        12343567,
				"disk_type":      "STD4",
				"iops":           900,
				"latency":        30,
				"storage_type":   "HDD",
			},
		},
		"network": []map[string]interface{}{
			{
				"nic":             5,
				"bandwidth_limit": 1000,
				"vlan":            mockVLANIdentifier,
				"id":              1235,
				"ips_v4": []string{
					"1.2.3.4",
				},
				"ips_v6": []string{
					"0000:111:2222:aaaa:bbbb",
				},
				"mac_address": "00:11:22:33:44:55",
			},
		},
		"version_tools":                    "guestToolsUnmanaged",
		"guest_tools_status":               "Active",
		"location_code":                    "ANX04",
		"location_country":                 "AT",
		"location_identifier":              locationIdentifier,
		"location_name":                    "ANX04 - AT, Vienna, Datasix",
		"provisioning_location_identifier": provisioningLocationIdentifier,
		"template_id":                      templateIdentifier,
		"resource_salesperson":             "AW",
	}
}

func prepareGet(name, desc string) {
	var response http.HandlerFunc

	if mockGetDeleted {
		response = RespondWith(404, ``)
	} else {
		body := mockADCInfoResponseBody(name, desc)
		response = RespondWithJSONEncoded(200, body)
	}

	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("GET", "/api/vsphere/v1/info.json/"+mockADCIdentifier+"/info"),
		response,
	))
}

func prepareCreate(name, desc string) {
	mock.AppendHandlers(CombineHandlers(
		VerifyRequest("POST", "/api/provisioning/v1/vm.json/"+locationIdentifier+"/"+templateType+"/"+templateIdentifier),
		VerifyJSONRepresenting(map[string]interface{}{
			"hostname":    name,
			"custom_name": desc,
			"memory_mb":   1024,
			"cpus":        1,
			"disk_gb":     4,
		}),
		RespondWithJSONEncoded(200, map[string]interface{}{
			"identifier": mockADCIdentifier,
			"progress":   10,
			"queued":     false,
			"errors":     []string{},
		}),
	))
}

func prepareList(name, desc string) {
	mockADCs := []map[string]interface{}{
		{"identifier": "foo0", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo1", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo2", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo3", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo4", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo5", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo6", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo7", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo8", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		{"identifier": "foo9", "name": "foo", "status": "poweredOff", "cores": 1, "cpu": 1, "ram": 1024, "location_code": "ANX04", "location_identifier": locationIdentifier},
		mockADCInfoResponseBody(name, desc),
	}

	pages := [][]map[string]interface{}{
		mockADCs[0:10],
		mockADCs[10:],
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
				"total_items": len(mockADCs),
				"limit":       len(data),
				"data":        data,
			}),
		))
	}
}
