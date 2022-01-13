package pagination

import (
	"context"
	"errors"

	"github.com/anexia-it/go-anxcloud/pkg/utils/param"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type testPageablePage struct {
	entries []string

	limit   int
	total   int
	page    int
	options []param.Parameter
}

func (p testPageablePage) Num() int {
	return p.page
}

func (p testPageablePage) Size() int {
	return len(p.entries)
}

func (p testPageablePage) Total() int {
	return p.total
}

func (p testPageablePage) Options() []param.Parameter {
	return p.options
}

func (p testPageablePage) Content() interface{} {
	return p.entries
}

type testPageable struct {
	entries []string
}

func (p testPageable) GetPage(ctx context.Context, page, limit int, opts ...param.Parameter) (Page, error) {
	startIdx := (page - 1) * limit
	endIdx := (page) * limit

	if startIdx > len(p.entries) {
		return nil, errors.New("Page out of range")
	}

	if endIdx > len(p.entries) {
		endIdx = len(p.entries) - 1
	}

	return testPageablePage{
		entries: p.entries[startIdx:endIdx],
		total:   len(p.entries),
		limit:   limit,
		page:    page,
		options: opts,
	}, nil
}

func (p testPageable) NextPage(ctx context.Context, page Page) (Page, error) {
	return p.GetPage(ctx, page.Num()+1, page.(testPageablePage).limit, page.(testPageablePage).options...)
}

var _ = Describe("AsChan function", func() {
	var testStrings = []string{
		"Hello world",
		"foo", "bar", "baz",
		"some random test strings",
		"black lives matter",
		"trans rights are human rights",
		"still only random strings, but why not these?",
	}

	var pageable Pageable

	BeforeEach(func() {
		pageable = testPageable{
			entries: testStrings,
		}
	})

	It("iterates pages via channel", func() {
		asChan, cancelFunc := AsChan(context.TODO(), pageable)
		defer cancelFunc()

		counter := 1
		for elem := range asChan {
			Expect(elem).To(Equal(testStrings[counter-1]))
			counter++
		}
		Expect(counter).To(Equal(len(testStrings)))
	})
})

var _ = Describe("LoopUntil function", func() {
	var testStrings = []string{
		"Hello world",
		"foo", "bar", "baz",
		"some random test strings",
		"black lives matter",
		"trans rights are human rights",
		"still only random strings, but why not these?",
	}

	var pageable Pageable

	BeforeEach(func() {
		pageable = testPageable{
			entries: testStrings,
		}
	})

	Context("stop condition never met", func() {
		It("loops through all elements and returns the expected error", func() {
			counter := 1
			err := LoopUntil(context.TODO(), pageable, func(elem interface{}) (bool, error) {
				Expect(elem).To(Equal(testStrings[counter-1]))
				counter++
				return false, nil
			})
			Expect(counter).To(Equal(len(testStrings)))
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ErrConditionNeverMet))
		})
	})

	Context("stop condition at last element", func() {
		It("loops through all elements and does not return an error", func() {
			counter := 1
			err := LoopUntil(context.TODO(), pageable, func(elem interface{}) (bool, error) {
				Expect(elem).To(Equal(testStrings[counter-1]))
				counter++

				if counter == len(testStrings) {
					return true, nil
				}

				return false, nil
			})
			Expect(counter).To(Equal(len(testStrings)))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("second element returns an error", func() {
		It("loops through first two elements and returns the error", func() {
			testError := errors.New("test error")

			counter := 1
			err := LoopUntil(context.TODO(), pageable, func(elem interface{}) (bool, error) {
				Expect(elem).To(Equal(testStrings[counter-1]))
				counter++

				if counter == 3 {
					return false, testError
				}

				return false, nil
			})
			Expect(counter).To(Equal(3))
			Expect(err).To(MatchError(testError))
		})
	})
})
