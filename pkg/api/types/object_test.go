package types

import (
	"context"
	"errors"
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

func (o *apiTestObject) GetIdentifier(context.Context) (string, error) {
	return o.Val, nil
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

var errTestFailingGetIdentifier = errors.New("failed to get identifier")

type apiTestObjectWithFailingGetIdentifier struct{}

func (apiTestObjectWithFailingGetIdentifier) EndpointURL(ctx context.Context) (*url.URL, error) {
	return url.Parse("/resource/v1")
}
func (apiTestObjectWithFailingGetIdentifier) GetIdentifier(context.Context) (string, error) {
	return "", errTestFailingGetIdentifier
}

var _ = Describe("GetObjectIdentifier function", func() {
	It("errors when an object is not set on singleObjectOperation", func() {
		sio := apiTestObject{}
		_, err := GetObjectIdentifier(&sio, true)
		Expect(err).To(MatchError(ErrUnidentifiedObject))
	})
	It("errors when object.GetIdentifier fails", func() {
		sio := apiTestObjectWithFailingGetIdentifier{}
		_, err := GetObjectIdentifier(&sio, true)
		Expect(err).To(MatchError(errTestFailingGetIdentifier))
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
