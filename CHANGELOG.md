## 0.3.25

ENHANCEMENTS

* lbaas - add load balancer as a service endpoints

## 0.3.24

ENHANCEMENTS

* vsphere/progress Fix a bug where 404 from progress API results in a loop until timeout

## 0.3.23

ENHANCEMENTS

* vmlist - Added VMList API client

## 0.3.22

ENHANCEMENTS

* clouddns - Added CloudDNS API client

## 0.3.21

ENHANCEMENTS

* client - Add client option to provide an `io.Writer` to dump http request/response for debugging

## 0.3.20

BUGFIXES

* vsphere/info - `DiskGB` variable type changed from float32 to float64

## 0.3.19

ENHANCEMENTS

* vlan - added `vm_provisioning` flag to `UpdateDefinition`
* vsphere/info - `DiskGB` attribute changed to floating point type
* ipam/address - random sleep added to ip reservation as workaround

## 0.3.18

BUGFIXES

* Changed location identifier tag name

## 0.3.17

BUGFIXES  

* Added multi VLAN support to ip-prefix. (#19)
