//go:build integration

package v1_test

import (
	"go.anx.io/go-anxcloud/pkg/api"
	ipamv1 "go.anx.io/go-anxcloud/pkg/apis/ipam/v1"
)

func preparePrefixCreate(a api.API, p ipamv1.Prefix)                        {}
func preparePrefixGet(a api.API, p ipamv1.Prefix)                           {}
func preparePrefixDelete(a api.API)                                         {}
func preparePrefixUpdate(a api.API, p ipamv1.Prefix, newDescription string) {}
func preparePrefixList(a api.API, p ipamv1.Prefix)                          {}
