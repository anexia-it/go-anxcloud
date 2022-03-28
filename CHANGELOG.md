# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
* generic client: ErrorFromResponse and NewHTTPError to allow easier mocking (#118, @LittleFox94)
* lbaas/v1: add state getters (StateSuccess, StateProgressing, StateFailure) (#116, @LittleFox94)
* utilities for comparing, searching and reconciling Objects by a list of (nested) attribute names (#117, @LittleFox94)
  - `pkg/utils/object/compare`

### Changed
* **breaking**: renamed corev1.Info to corev1.Resource (#113, @LittleFox94)

<!--
Please add your release notes under the correct category (Added, Changed, ...) and use the following format as a
guideline:

* <api-scope>[/<sub-api-scope] - <a short description of what changed>

e.g.

* vsphere/provisioning - added progress identifier
-->

## [0.4.1] - 2022-03-04

### Added
* new APIs supported with generic client:
  - core/location
  - vlan

### Changed
* github.com/satori/go.uuid updated to latest master
  - fixes a security issue for random UUIDs being not that random (CVE-2021-3538, https://github.com/satori/go.uuid/issues/73)
    + although this is a real security issue, it's not relevant for go-anxcloud as it only consumes UUIDs in production code
    + UUIDs were generated in test code, but those not being as random as they should be isn't a real problem
  - the library seems dead and all usages not part of a public interface have been removed from go-anxcloud to
    prepare removing it completely when the legacy CloudDNS client is gone

## [0.4.0] - 2022-01-24

### Added
* new client, unifying features across APIs and reducing code duplication (PR #56)
  - already supported:
    + core/resources
    + lbaas
    + clouddns/zone
* client: new option `BaseURL` (PR #58)
* client: interface to retrieve metrics, such as requests currently in-flight or request duration (PR #66)
* old-style clients:
  - lbaas: pagination support (PR #45)
  - core/location: add getters for locations by identifier and code (PR #49)
  - vsphere/info: add CPU performance type and CPU clock rate attributes (PR #51)

### Changed
* import path is now go.anx.io/go-anxcloud, please change your code accordingly
* client: now uses [`logr`](https://github.com/go-logr/logr) for logging (PR #50)
* package is now tested against Go versions 1.16 and 1.17

### Deprecated
* the old-style clients are deprecated and will be removed in the minor version following the version with everything supported by the generic client we have old-style clients for
  - write code against the generic client instead, create issues for APIs you need to help us prioritize them
* client: the `LogWriter` option for dumping requests and responses is replaced by the `Logger` option (PR #50)

### Removed
* client: DefaultBaseURL was exported by mistake and is not exposed anymore, use the BaseURL method on the client instance instead (PR #58)

### Fixed
* connections could end up dangling around (PR #48)

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
