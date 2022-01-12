# Dev workflows

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


## Code generator

Our build tools include a small code generator, currently only used to help with some tasks related to `Object`s,
the interface needed to be compatible with the generic client. It's code lives in `/tools/object_generator.go`.

This was built to deal with go's interfaces being implicit. The generic client allows `Object`s to implement
additional interfaces to hook into parts of the request processing and we want to make sure `Object`s always
implement the interfaces they think they do. The code generator is used to generate some tests (using `ginkgo` and
`gomega`) to make sure the interfaces specified in the magic comment are really implemented. These tests will only
be run when there is a spec runner test file for the package already - but you should have that anyway.

Only files with names not starting with `.`, ending with `.go` and not ending with `_test.go` are parsed, which
translates to every non-hidden non-test go file.


### Workflow & CI

When changing anything on a magic comment, a type marked with a magic comment or the code generator itself, you
have to run `make generate` to re-generate files. The CI will check if `make generate` changes anything, failing
if something was changed influencing the generated code without re-generating it.


### Magic comment format

Comments always start with `// anxcloud:`, that's what the generator is looking for. Single space required between the
comment-starting `//` and `anxcloud:`. To process the comment, this common prefix is stripped from it.

The payload of the comment (what is left after stripping the prefix) is then split by `:` to get some `specs`,
which are then split by `=` to have a spec name and value. Not all specs have a value, if a spec has multiple
values, it has them separated by `,`.

Some examples:

```go
// anxcloud:object
// -> this is parsed to having the spec called 'object' with no value

// anxcloud:object:hooks=ResponseBodyHook
// -> this is parsed to having the spec called 'object' with no value
//    and the spec 'hooks' with value 'ResponseBodyHook'.

// anxcloud:object:hooks=ResponseBodyHook,ResponseDecodeHook
// -> this is parsed to having the spec called 'object' with no value
//    and the spec 'hooks' with value 'ResponseBodyHook,ResponseDecodeHook'.
//    Code handling the 'hooks' spec will split the value by ',' to decode
//    single elements.
```

Currently all the specs have to be given in the same comment line, placed above the `type` keyword they apply to
**with one blank line in between**. The blank line ensures the magic comment isn't written into the documentation.

Example how it looks in real world:

```go
// anxcloud:object:hooks=ResponseBodyHook

// LoadBalancer describes a single LoadBalancer in Anexia LBaaS API.
type LoadBalancer struct {
    // [...]
}
```

This example makes sure the type `LoadBalancer` does everything necessary to be usable with the generic client
(`object` spec) and always implements the interface for the `ResponseBodyHook` correctly.


#### Known specs and their usage

* `object`

    | Usable on | Value  |
    |-----------|--------|
    | types     | (none) |

    Specifies the type is an `Object`, something usable with the generic API client.


* `hooks`

    | Usable on | Value  |
    |-----------|--------|
    | types     | names of hook interfaces from `pkg/api/types` |

    Explicitly specifies the type implements the given interfaces.

## pre-commit hook

There is a `pre-commit` hook checking some things before you even make a commit, you can install it with
`make install-precommit-hook`. It might occasionally warn you about it being not up to date, run that command
again to update it.

It's not very clever yet and in some cases might hinder working, for example when rebasing. You can always just
uninstall it with `rm .git/hooks/pre-commit`.
