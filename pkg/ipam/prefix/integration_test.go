//go:build integration
// +build integration

package prefix

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.anx.io/go-anxcloud/pkg/client"
	"go.anx.io/go-anxcloud/pkg/ipam/address"
)

const (
	locationID = "52b5f6b2fd3a4a7eaaedf1a7c019e9ea"
	vlanID     = "166fa87362c8498f8c4aa6d1c5b9042c"
)

var _ = Describe("ipam/prefix client", func() {
	var api API
	var cli client.Client

	BeforeEach(func() {
		c, err := client.New(client.AuthFromEnv(false))
		Expect(err).ToNot(HaveOccurred())
		cli = c

		api = NewAPI(cli)
	})

	checkEmpty := func(prefixID *string, createEmpty bool) {
		It("created the expected amount of addresses", func() {
			a := address.NewAPI(cli)

			filtered, err := a.GetFiltered(context.TODO(), 1, 1000, address.PrefixFilter(*prefixID))
			Expect(err).ToNot(HaveOccurred())

			if createEmpty {
				Expect(filtered).To(HaveLen(3)) // only network, broadcast and router addresses are created
			} else {
				Expect(filtered).To(HaveLen(8)) // we expect all IPs to be already created
			}
		})
	}

	prefixTest := func(createEmpty bool) {
		Context("with a prefix created for testing", Ordered, func() {
			var prefix string
			BeforeAll(func() {
				create := NewCreate(
					locationID,
					vlanID,
					4,
					TypePrivate,
					29,
				)
				create.CreateEmpty = createEmpty

				p, err := api.Create(context.TODO(), create)
				Expect(err).NotTo(HaveOccurred())

				DeferCleanup(func() {
					err := api.Delete(context.TODO(), p.ID)
					Expect(err).NotTo(HaveOccurred())
				})

				prefix = p.ID
			})

			It("lists prefixes including our test prefix", func() {
				found := false
				page := 1

				for !found {
					ps, err := api.List(context.TODO(), page, 20)
					Expect(err).NotTo(HaveOccurred())
					Expect(ps).NotTo(BeEmpty())

					for _, p := range ps {
						if p.ID == prefix {
							found = true
							break
						}
					}

					page++
				}
			})

			It("eventually retrieves test prefix with expected data and being Active", func() {
				poll := func(g Gomega) {
					info, err := api.Get(context.TODO(), prefix)
					g.Expect(err).NotTo(HaveOccurred())

					g.Expect(info.Vlans).To(HaveLen(1))
					g.Expect(info.Vlans[0].ID).To(Equal(vlanID))

					g.Expect(info.PrefixType).To(BeEquivalentTo(TypePrivate))

					g.Expect(info.Status).To(Equal("Active"))
				}

				Eventually(poll, 5*time.Minute, 10*time.Second).Should(Succeed())
			})

			checkEmpty(&prefix, createEmpty)

			It("updates the test prefix with changed data", func() {
				p, err := api.Update(context.TODO(), prefix, Update{CustomerDescription: "something else"})
				Expect(err).NotTo(HaveOccurred())

				Expect(p.ID).To(Equal(prefix))
				Expect(p.CustomerDescription).To(Equal("something else"))
			})
		})
	}

	Context("with createEmpty set to false", func() {
		prefixTest(false)
	})

	Context("with createEmpty set to true", func() {
		prefixTest(true)
	})
})
