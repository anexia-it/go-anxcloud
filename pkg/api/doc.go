// Package api implements a generic API client for our engine, reducing the amount of duplicated code in
// this library.
//
// It is heavily inspired by the kubernetes client-go, but also allows handling quirks in the engine API
// in a graceful way, without making quirky APIs incompatible with the generic code in this package or
// letting the user of this library handle the quirks.
package api
