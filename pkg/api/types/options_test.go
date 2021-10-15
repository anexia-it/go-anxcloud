package types

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Options type", func() {
	var options []Options

	BeforeEach(func() {
		options = []Options{
			&GetOptions{},
			&ListOptions{},
			&CreateOptions{},
			&DestroyOptions{},
			&UpdateOptions{},
		}
	})

	It("can set additional options from all operation specific option types", func() {
		for i, opts := range options {
			err := opts.Set("test", i, false)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("barks when setting additional keys twice without overwrite set", func() {
		for i, opts := range options {
			err := opts.Set("test", i, false)
			Expect(err).NotTo(HaveOccurred())

			err = opts.Set("test", i, false)
			Expect(err).To(MatchError(ErrKeyAlreadySet))
		}
	})

	It("does not bark when setting additional keys twice with overwrite set", func() {
		for i, opts := range options {
			err := opts.Set("test", i, false)
			Expect(err).NotTo(HaveOccurred())

			err = opts.Set("test", i, true)
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("barks when retrieving a not-set key", func() {
		for _, opts := range options {
			val, err := opts.Get("test")
			Expect(err).To(MatchError(ErrKeyNotSet))
			Expect(val).To(BeNil())
		}
	})

	It("does not bark when retrieving an additional key that was set before", func() {
		for i, opts := range options {
			err := opts.Set("test", i, false)
			Expect(err).NotTo(HaveOccurred())

			val, err := opts.Get("test")
			Expect(err).NotTo(HaveOccurred())
			Expect(val).To(Equal(i))
		}
	})
})

func TestAPIUnits(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "test suite for pkg/api/types")
}
