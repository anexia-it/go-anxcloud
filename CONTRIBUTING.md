# Guidance on how to contribute

> By submitting a pull request or filing a bug, issue, or feature request,
> you are agreeing to comply with this waiver of copyright interest.
> Details can be found in our [LICENSE](LICENSE).


There are two primary ways to help:
 - Using the issue tracker, and
 - Changing the code-base.


## Using the issue tracker

Use the issue tracker to suggest feature requests, report bugs, and ask questions.
This is also a great way to connect with the developers of the project as well
as others who are interested in this solution.

You can also use it to find a bug to fix or feature to implement. Mention in
the issue that you want to work on it (to prevent multiple people working on the same),
then follow the _Changing the code-base_ guidance below.


## Changing the code-base

Generally speaking, you should fork this repository, make changes in your
own fork, and then submit a pull request. People having commit access on this repository can
also push their branches in this repository instead of a fork, but still have to open pull
requests to have things merged into `main`.

All new code should have associated unit tests that validate implemented features
and the presence or lack of defects. For tests we use [`ginkgo`](https://onsi.github.io/ginkgo/)
and [`gomega`](https://onsi.github.io/gomega/), maybe take a look at
[`pkg/api/errors_test.go`](pkg/api/errors_test.go) for an example on how we use it.

Code in this project follows the standard go code style (`go fmt`), our CI system enforces it.
Please add new code in a sensible location, using the trimmed-down tree below as a guideline.

```plain
- pkg
+-- api         Generic client for our API
| +-- types     Things needed to implement Objects compatible with the generic client
|
+-- apis        Base for APIs usable with the generic client
| |
| +-- lbaas     Types and supporting code for multiple versions of LBaaS API
| | +-- v1      Types and supporting code for using LBaaS API v1 with the generic client
| |
| +-- clouddns  Types and supporting code for multiple versions of CloudDNS API
|   +-- v1      Types and supporting code for using CloudDNS API v1 with the generic client
|
+-- client      Wrapper around http.Client adding logging, authentication, ...
|
+-- utils       Utilities used in more than one package
  +-- test      Utilities for testing things
```

There are a lot more packages in this repository, many for the legacy clients. There is also the
`tests` directory in the repository root, containing the old integration tests not yet moved into
their packages. This tree is supposed to be an example and guideline, not a full list of everything
in the project.


### pre-commit hook

There is a `pre-commit` hook checking some things before you even make a commit, you can install it with
`make install-precommit-hook`. It might occasionally warn you about it being not up to date, run that command
again to update it.

It's not very clever yet and in some cases might hinder working, for example when rebasing. You can always just
uninstall it with `rm .git/hooks/pre-commit`.


## Testing

We use `Ginkgo` (v2) for our tests. Unit and integration tests are both directly in the package they test, integration
tests are disabled by default with build tags, like this:

```go
//go:build integration
// +build integration

package foo
```

Integration tests for bigger scopes (in a package and testing integration of sub-packages together) are placed in the
matching higher hierarchy package (e.g.: test basics in `pkg/vsphere/provisioning/disktype` and test whole "create VM
and do things with it" in `pkg/vsphere`).

Tests are executed with either `make test` (for unit tests) or `make func-test` (for e2e/integration tests). Note the
integration tests need an authentication token in environment variable `ANEXIA_TOKEN`.

To run unit tests for a single package or everything below a given package you can use either `go test ./pkg/lbaas/acl`
(for testing only acl package) or `go test ./pkg/lbaas/...` (for testing everything below ./pkg/lbaas, including the
`lbaas` package itself). For running integration tags, add the build tag like `go test -tags integration ./pkg/lbaas/...`.
You can also use the `ginkgo` CLI tool `ginkgo -tags integration pkg/lbaas/...`.


### Code generator

We use a code generator to automate some manual work. Currently it's only used to generate tests making sure
and `Object` (something you can use with the generic client) really implements hooks. It works by adding some
magic comments in the format `// anexia:something:something-else=foo,bar`.

Whenever you touch such a comment, something tagged with such a comment or the code generator itself you have
to update the generated code by running `make generate` and commit it together with your changes. The CI will
bark at you when you forget this (but we won't, promise).

More info about this can be found in [its documentation](docs/code-generator.md).
