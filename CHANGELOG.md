# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

<!--
Please add your changelog entry under this comment in the correct category (Security, Fixed, Added, Changed, Deprecated, Removed - in this order).

Changelog entries are best in the following format, where scope is something like "generic client" or "lbaas/v1"
(for LBaaS API bindings). If the change isn't user-facing but still relevant enough for a changelog entry, add
"(internal)" before the scope.

* (internal)? scope: short description (pull request, author)

Some examples, more below in the actual changelog (newer entries are more likely to be good entries):
* generic client: List resources with a channel (#42, @LittleFox94)
* core/v1: added helper methods to tag resources (#122, @marioreggiori)
* (internal) generic client: add hook FilterRequestURLHook (#123, @marioreggiori)

-->
### Added
* api: Add RateLimitError (#438, @nachtjasmin)
* Replace all occurrences of `vmxnet3` with `virtio` (#444, @nachtjasmin)

## [0.7.7] -- 2025-01-09

### Added

* api/mock/mock_api_implementation: panic when trying to fake `ResourceWithTag` (#419, @drpsychick)

### Changed

* go-anxcloud is now tested with Go versions 1.22 and 1.23 (#428, @drpsychick)

### Fixed

* CloudDNS now consistently uses `zone_name` during requests, fixing integration tests (#432,  @drpsychick)
* golang/x/net: update to 0.33.0 due to CVE-2024-45338 (#431, @drpsychick)

## [0.7.6] -- 2024-11-04

### Added

* vsphere/info/network: add missing `bandwidth_limit` field in `info.Network` (#417, @89q12) 

## [0.7.5] -- 2024-10-22

### Added

* vsphere/provisioning/progress: fix `Status`, it's a string (#412, @drpsychick)

## [0.7.4] -- 2024-10-21

### Added

* vsphere/provisioning/vm: add `BandwidthLimit` option (#408 @ProbstenHias)
* vsphere/provisioning/progress: add `Status` (#410, @drpsychick)
* kubernetes/v1: Add `CniPlugin`, `ApiServerAllowlist` and `ApiServerAllowlist` options (#409, @ProbstenHias)

## [0.7.3] -- 2024-08-06

### Added

* pkg/api/mock: add pre-update hook (#392, @nachtjasmin)

## [0.7.2] - 2024-06-19

### Fixed
* pkg/api: use logger of client if `WithLogger` was not called explicitly (#381, @nachtjasmin)

## [0.7.1] - 2024-06-06

### Added
* ipam/address: Add ip-version and prefix filter attributes to `ReserveRandom` request (#378, @anx-mschaefer)

### Fixed
* ipam: Fix several error messages to contain the correct resource (#378, @anx-mschaefer)

## [0.7.0] - 2024-05-24

* go-anxcloud is now tested with Go versions 1.21 and 1.22 and we might start using features only available since 1.21

### Added
* `WithRequestOptions` API option to configure default request options (#361, @anx-mschaefer)
* kubernetes/v1: Add `autoscaling` attribute via `EnableAutoscaling` field (#369, @nachtjasmin)

### Changed
* (internal) add "error-return" to request option interfaces (#361, @anx-mschaefer)

## [0.6.4] - 2024-03-15

### Fixed
* fix renamed inputs for crazy-max/ghaction-import-gpg action in release workflow

## [0.6.3] - 2024-03-15

### Changed
* e5e/v1: rename storage backend 'zip' to 'archive' and add new `deployment_state` attribute (#347, @anx-mschaefer)

## [0.6.2] - 2023-12-12

### Added
* generic mock client: introduce pre-create hook (#322, @anx-mschaefer)

## [0.6.1] - 2023-11-17

### Added
* e5e/v1: added new `Application` & `Function` resources (#317, @anx-mschaefer)

## [0.6.0] - 2023-11-10

* go-anxcloud is now tested with Go versions 1.20 and 1.21 and we might start using features only available since 1.20

### Added
* core/v1: add `created_at` and `updated_at` fields to Resource type
* generic client: make common.PartialResource filterable
* lbaas/v2: added new `Cluster`, `Node` & `LoadBalancer` resources (#309, @anx-mschaefer)
* frontier/v1: added new `API`, `Action`, `Deployment` & `Endpoint` resources (#311, @anx-mschaefer)

## [0.5.3] - 2023-06-16

### Added
* vsphere/provisioning/vm: add AdditionalDisks field to definition (#275, @marioreggiori)

## [0.5.2] - 2023-06-06

### Fixed
* (internal) core/v1: Added retry on conflict(409) engine error while tagging resources (#266, @11tuvork28)

### Changed
* go-anxcloud is now tested with Go versions 1.19 and 1.20 and we might start using features only available since 1.19
* github workflows refactored to run integration tests only when labeling a pull request with `integration_tests`

## [0.5.1] - 2023-02-17

### Changed
* trace logging in pkg/client now goes to a logr.Logger attached to the request context, falling back to the logger configured on the client (#212, @LittleFox94)

### Fixed
* trace logging in pkg/client now really includes the request/response bodies (#211, @LittleFox94)
* kubernetes/v1: GetKubeConfig helper waits until kubeconfig is available (#221, @marioreggiori)

### Added
* kubernetes/v1: configurable cluster prefixes (#208, @marioreggiori)
* generic client: common resource types (#208, @marioreggiori)

## [0.5.0] - 2022-11-30

### Added
* kubernetes/v1: add Cluster and NodePool bindings (#151, @marioreggiori)
* kubernetes/v1: add kubeconfig helper (#175, @marioreggiori)

### Changed
* lbaas/v1: changed state retriever interface `StateSuccess`, `StateProgressing` & `StateFailure` to `StateOK`, `StatePending` & `StateError` (#185, @marioreggiori)
* (internal) generic client: common Generic Service interfaces and utils have been extracted into a new shared package (#185, @marioreggiori)

## [0.4.6] - 2022-10-04

### Added
* generic client: AutoTag create option added (#167, @marioreggiori)

### Changed
* upgraded to target Go 1.18, so now testing with Go versions 1.18 and 1.19 (#161, @LittleFox94)
* clouddns: prevent creation of DNS records with empty names due to API change -> use "@" to target domain root instead (#168, @marioreggiori)
* (internal) generic client apis: change pkg of tests to `v1_test` to mitigate import cycle issues (#169, @marioreggiori)

## [0.4.5] - 2022-08-08

### Added
* vsphere/v1: Template bindings and FindNamedTemplate helper added to retrieve templates by name (#148, @marioreggiori)
* generic client: mock implementation (#139, @marioreggiori)
* (internal) object-generator: generate GetIdentifier method in `runtime` mode (#150, @marioreggiori)

### Changed
* generic client: GetIdentifier method added to types.Object interface (#150, @marioreggiori)
* (internal) generic client: uses types.Object.GetIdentifier method (#150, @marioreggiori)

## [0.4.4] - 2022-06-10

### Fixed
* pkg/utils/object/compare.Reconcile now accepts arrays of `*struct` and `types.Object` as target/existing input (#145, @LittleFox94)

### Added
* core/v1: helper methods Tag, Untag and ListTags (#122, @marioreggiori)
* lbaas/v1: ACL and Rule API bindings added (#142, @toothstone & @marioreggiori)

### Changed
* moved pkg/api.GetObjectIdentifier and related errors to pkg/api/types (#144, @LittleFox94)
  - the previous locations are still available, but marked as deprecated

## [0.4.3] - 2022-05-04

### Fixed
* clouddns/v1: creating a Record didn't retrieve its Identifier (#120, @LittleFox94)
* lbaas/v1: fix some attributes not being sent to the Engine when creating Backends (#135, @LittleFox94)

### Added
* (internal) generic client: FilterRequestURLHook for modifying request URLs (#123, @marioreggiori)

### Changed
* (internal) core/v1: ResourceWithTag uses RequestBodyHook and FilterRequestURLHook instead of RequestFilterHook (#123, @marioreggiori)

## [0.4.2] - 2022-03-29

### Added
* generic client: ErrorFromResponse and NewHTTPError to allow easier mocking (#118, @LittleFox94)
* lbaas/v1: add state getters (StateSuccess, StateProgressing, StateFailure) (#116, @LittleFox94)
  - **breaking**: changes the `State` attribute of all resources to be added via embedded `HasState`
* utilities for comparing, searching and reconciling Objects by a list of (nested) attribute names (#117, @LittleFox94)
  - `pkg/utils/object/compare`

### Changed
* **breaking**: renamed corev1.Info to corev1.Resource (#113, @LittleFox94)

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
