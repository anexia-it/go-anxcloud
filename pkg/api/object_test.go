package api

import (
	"context"
	"net/url"

	uuid "github.com/satori/go.uuid"

	"github.com/anexia-it/go-anxcloud/pkg/api/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type api_test_nonstruct_object bool

func (o *api_test_nonstruct_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_nonpointer_object bool

func (o api_test_nonpointer_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_noident_object struct {
	Value string `json:"value"`
}

func (o api_test_noident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_invalidident_object struct {
	Value int `json:"value" anxcloud:"identifier"`
}

func (o api_test_invalidident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_uuidident_object struct {
	Identifier uuid.UUID `json:"identifier" anxcloud:"identifier"`
}

func (o api_test_uuidident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/invalid_anyway")
}

type api_test_embeddedident_object struct {
	api_test_object
}

type api_test_ptrembeddedident_object struct {
	*api_test_object
}

type api_test_multiembeddedident_object struct {
	uuid.UUID
	api_test_object
}

type api_test_multiident_object struct {
	Identifier  string `json:"identifier" anxcloud:"identifier"`
	Identifier2 string `json:"identifier2" anxcloud:"identifier"`
}

func (o api_test_multiident_object) EndpointURL(ctx context.Context, op types.Operation, opts types.Options) (*url.URL, error) {
	return url.Parse("/resource/v1")
}

var _ = Describe("getObjectIdentifier function", func() {
	It("errors out on invalid Object types", func() {
		nso := api_test_nonstruct_object(false)
		identifier, err := getObjectIdentifier(&nso, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented as structs"))
		Expect(identifier).To(BeEmpty())

		npo := api_test_nonpointer_object(false)
		identifier, err = getObjectIdentifier(npo, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("must be implemented on a pointer to struct"))
		Expect(identifier).To(BeEmpty())

		nio := api_test_noident_object{"invalid"}
		identifier, err = getObjectIdentifier(&nio, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("lacks identifier field"))
		Expect(identifier).To(BeEmpty())

		iio := api_test_invalidident_object{32}
		identifier, err = getObjectIdentifier(&iio, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("identifier field has an unsupported type"))
		Expect(identifier).To(BeEmpty())

		mio := api_test_multiident_object{"identifier", "identifier2"}
		identifier, err = getObjectIdentifier(&mio, false)
		Expect(err).To(MatchError(ErrTypeNotSupported))
		Expect(err.Error()).To(ContainSubstring("api_test_multiident_object has multiple fields tagged as identifier"))
		Expect(identifier).To(BeEmpty())
	})

	It("accepts valid Object types", func() {
		sio := api_test_object{"identifier"}
		identifier, err := getObjectIdentifier(&sio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		uio := api_test_uuidident_object{uuid.FromStringOrNil("6010622e-3e14-11ec-a5c3-0f457821b3ba")}
		identifier, err = getObjectIdentifier(&uio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("6010622e-3e14-11ec-a5c3-0f457821b3ba"))
	})

	It("accepts valid Object types where the identifier is in embedded fields", func() {
		eio := api_test_embeddedident_object{api_test_object{"identifier"}}
		identifier, err := getObjectIdentifier(&eio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		peio := api_test_ptrembeddedident_object{&api_test_object{"identifier"}}
		identifier, err = getObjectIdentifier(&peio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))

		meio := api_test_multiembeddedident_object{uuid.NewV4(), api_test_object{"identifier"}}
		identifier, err = getObjectIdentifier(&meio, true)
		Expect(err).NotTo(HaveOccurred())
		Expect(identifier).To(Equal("identifier"))
	})

	Context("when doing an operation on a specific object", func() {
		It("errors out on valid Object type but empty identifier", func() {
			o := api_test_object{""}
			identifier, err := getObjectIdentifier(&o, true)
			Expect(err).To(MatchError(ErrUnidentifiedObject))
			Expect(identifier).To(BeEmpty())
		})

		It("returns the correct identifier", func() {
			o := api_test_object{"test"}
			identifier, err := getObjectIdentifier(&o, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(identifier).To(Equal("test"))
		})
	})
})
