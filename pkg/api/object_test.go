package api

import (
	"context"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type api_test_nonstruct_object bool

func (o *api_test_nonstruct_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_nonpointer_object bool

func (o api_test_nonpointer_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_noident_object struct {
	Value string `json:"value"`
}

func (o api_test_noident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_invalidident_object struct {
	Value int `json:"value" anxcloud:"identifier"`
}

func (o api_test_invalidident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_embeddedident_object struct {
	api_test_object
}

type api_test_ptrembeddedident_object struct {
	*api_test_object
}

type api_test_multiembeddedident_object struct {
	url.Values
	api_test_object
}

func (o api_test_multiembeddedident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type api_test_multiident_object struct {
	Identifier  string `json:"identifier" anxcloud:"identifier"`
	Identifier2 string `json:"identifier2" anxcloud:"identifier"`
}

func (o api_test_multiident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type api_test_embeddedmultiident_object struct {
	api_test_multiident_object
}

func (o api_test_embeddedmultiident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

type api_test_multiembeddedmultiident_object struct {
	api_test_object
	api_test_embeddedident_object
}

func (o api_test_multiembeddedmultiident_object) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

var _ = Describe("GetObjectIdentifier function", func() {
	It("errors out on invalid Object types", func() {
		nso := api_test_nonstruct_object(false)
		identifier, err := GetObjectIdentifier(&nso, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented as structs"))
		Expect(identifier).To(BeEmpty())

		npo := api_test_nonpointer_object(false)
		identifier, err = GetObjectIdentifier(npo, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented on a pointer to struct"))
		Expect(identifier).To(BeEmpty())

		nio := api_test_noident_object{"invalid"}
		identifier, err = GetObjectIdentifier(&nio, false)
		Expect(err).To(MatchError(ErrObjectWithoutIdentifier))
		Expect(identifier).To(BeEmpty())

		iio := api_test_invalidident_object{32}
		identifier, err = GetObjectIdentifier(&iio, false)
		Expect(err).To(MatchError(ErrObjectIdentifierTypeNotSupported))
		Expect(identifier).To(BeEmpty())

		mio := api_test_multiident_object{"identifier", "identifier2"}
		identifier, err = GetObjectIdentifier(&mio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())

		emio := api_test_embeddedmultiident_object{api_test_multiident_object{"identifier", "identifier2"}}
		identifier, err = GetObjectIdentifier(&emio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())

		memio := api_test_multiembeddedmultiident_object{api_test_object{"identifier"}, api_test_embeddedident_object{api_test_object{"another identifier"}}}
		identifier, err = GetObjectIdentifier(&memio, false)
		Expect(err).To(MatchError(ErrObjectWithMultipleIdentifier))
		Expect(identifier).To(BeEmpty())
	})

	It("accepts valid Object types", func() {
		sio := api_test_object{"identifier"}
		identifier, err := GetObjectIdentifier(&sio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))
	})

	It("accepts valid Object types where the identifier is in embedded fields", func() {
		eio := api_test_embeddedident_object{api_test_object{"identifier"}}
		identifier, err := GetObjectIdentifier(&eio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		peio := api_test_ptrembeddedident_object{&api_test_object{"identifier"}}
		identifier, err = GetObjectIdentifier(&peio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		meio := api_test_multiembeddedident_object{url.Values{}, api_test_object{"identifier"}}
		identifier, err = GetObjectIdentifier(&meio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))
	})

	Context("when doing an operation on a specific object", func() {
		It("errors out on valid Object type but empty identifier", func() {
			o := api_test_object{""}
			identifier, err := GetObjectIdentifier(&o, true)
			Expect(err).To(MatchError(ErrUnidentifiedObject))
			Expect(identifier).To(BeEmpty())
		})

		It("returns the correct identifier", func() {
			o := api_test_object{"test"}
			identifier, err := GetObjectIdentifier(&o, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(identifier).To(Equal("test"))
		})
	})
})
