package types

import (
	"context"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type RandomEmbeddedData struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

type apiTestObject struct {
	// mimicks a bug triggered by lbaas/v1, where all Objects embed HasState at top of the object
	RandomEmbeddedData

	Val string `json:"value" anxcloud:"identifier"`
}

func (o *apiTestObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type apiTestNonstructObject bool

func (o *apiTestNonstructObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type apiTestNonpointerObject bool

func (o apiTestNonpointerObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type apiTestNoidentObject struct {
	Value string `json:"value"`
}

func (o apiTestNoidentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type apiTestInvalididentObject struct {
	Value int `json:"value" anxcloud:"identifier"`
}

func (o apiTestInvalididentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type apiTestEmbeddedidentObject struct {
	apiTestObject
}

type apiTestPtrembeddedidentObject struct {
	*apiTestObject
}

type apiTestMultiembeddedidentObject struct {
	url.Values
	apiTestObject
}

func (o apiTestMultiembeddedidentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type apiTestMultiidentObject struct {
	Identifier  string `json:"identifier" anxcloud:"identifier"`
	Identifier2 string `json:"identifier2" anxcloud:"identifier"`
}

func (o apiTestMultiidentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type apiTestEmbeddedmultiidentObject struct {
	apiTestMultiidentObject
}

func (o apiTestEmbeddedmultiidentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type apiTestMultiembeddedmultiidentObject struct {
	apiTestObject
	apiTestEmbeddedidentObject
}

func (o apiTestMultiembeddedmultiidentObject) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

var _ = Describe("GetObjectIdentifier function", func() {
	It("errors out on invalid Object types", func() {
		nso := apiTestNonstructObject(false)
		identifier, err := GetObjectIdentifier(&nso, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented as structs"))
		Expect(identifier).To(BeEmpty())

		npo := apiTestNonpointerObject(false)
		identifier, err = GetObjectIdentifier(npo, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented on a pointer to struct"))
		Expect(identifier).To(BeEmpty())

		nio := apiTestNoidentObject{"invalid"}
		identifier, err = GetObjectIdentifier(&nio, false)
		Expect(err).To(MatchError(ErrObjectWithoutIdentifier))
		Expect(identifier).To(BeEmpty())

		iio := apiTestInvalididentObject{32}
		identifier, err = GetObjectIdentifier(&iio, false)
		Expect(err).To(MatchError(ErrObjectIdentifierTypeNotSupported))
		Expect(identifier).To(BeEmpty())

		mio := apiTestMultiidentObject{"identifier", "identifier2"}
		identifier, err = GetObjectIdentifier(&mio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())

		emio := apiTestEmbeddedmultiidentObject{apiTestMultiidentObject{"identifier", "identifier2"}}
		identifier, err = GetObjectIdentifier(&emio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())

		memio := apiTestMultiembeddedmultiidentObject{apiTestObject{RandomEmbeddedData{}, "identifier"}, apiTestEmbeddedidentObject{apiTestObject{RandomEmbeddedData{}, "another identifier"}}}
		identifier, err = GetObjectIdentifier(&memio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())
	})

	It("accepts valid Object types", func() {
		sio := apiTestObject{RandomEmbeddedData{}, "identifier"}
		identifier, err := GetObjectIdentifier(&sio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))
	})

	It("accepts valid Object types where the identifier is in embedded fields", func() {
		eio := apiTestEmbeddedidentObject{apiTestObject{RandomEmbeddedData{}, "identifier"}}
		identifier, err := GetObjectIdentifier(&eio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		peio := apiTestPtrembeddedidentObject{&apiTestObject{RandomEmbeddedData{}, "identifier"}}
		identifier, err = GetObjectIdentifier(&peio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		meio := apiTestMultiembeddedidentObject{url.Values{}, apiTestObject{RandomEmbeddedData{}, "identifier"}}
		identifier, err = GetObjectIdentifier(&meio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))
	})

	Context("when doing an operation on a specific object", func() {
		It("errors out on valid Object type but empty identifier", func() {
			o := apiTestObject{RandomEmbeddedData{}, ""}
			identifier, err := GetObjectIdentifier(&o, true)
			Expect(err).To(MatchError(ErrUnidentifiedObject))
			Expect(identifier).To(BeEmpty())
		})

		It("returns the correct identifier", func() {
			o := apiTestObject{RandomEmbeddedData{}, "test"}
			identifier, err := GetObjectIdentifier(&o, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(identifier).To(Equal("test"))
		})
	})
})
