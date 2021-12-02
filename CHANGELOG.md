# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.3.28] - 2021-10-05
### Added
* ipam/prefix - allow to create empty prefixes
* ipam/address - allow to retrieve addresses with given filters

## [0.3.27] - 2021-09-27
### Added
*  ipam/prefix - Prefix Type is now available in Info

## [0.3.26] - 2021-09-23
### Added
*  vsphere - Deprovision now returns progress identifier

## [0.3.25] - 2021-09-22
### Added
* lbaas - add load balancer as a service endpoints
* lbaas - add acl endpoints
* client - allow to specify a user agent option

## [0.3.24] - 2021-08-23
### Added
* vsphere/progress Fix a bug where 404 from progress API results in a loop until timeout

## [0.3.23] - 2021-08-12
### Added
* vmlist - Added VMList API client

## [0.3.22] - 2021-07-20
### Added
* clouddns - Added CloudDNS API client

## [0.3.21] - 2021-07-17
### Added
* client - Add client option to provide an `io.Writer` to dump http request/response for debugging

## [0.3.20] - 2021-04-07
### Fixed
* vsphere/info - `DiskGB` variable type changed from float32 to float64

## [0.3.19] - 2021-03-29
### Added
* vlan - added `vm_provisioning` flag to `UpdateDefinition`
* vsphere/info - `DiskGB` attribute changed to floating point type
* ipam/address - random sleep added to ip reservation as workaround

## [0.3.18] - 2021-03-22
### Fixed
* Changed location identifier tag name

## [0.3.17] - 2021-03-19
### Fixed
* Added multi VLAN support to ip-prefix. (#19)
