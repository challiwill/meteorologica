package azure_test

import (
	. "github.com/challiwill/meteorologica/azure"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Usage", func() {
	Describe("Hash", func() {
		var (
			usage      Usage
			otherUsage Usage
		)

		BeforeEach(func() {
			usage = Usage{
				AccountOwnerId:         "some-account-owner-id",
				AccountName:            "some-account-name",
				ServiceAdministratorId: "some-admin-id",
				SubscriptionId:         "some-sub-id",
				SubscriptionGuid:       "some-sub-guid",
				SubscriptionName:       "some-sub-name",
				Date:                   "some-date",
				Month:                  2,
				Day:                    3,
				Year:                   1000,
				Product:                "some-product",
				MeterID:                "some-meter-id",
				MeterCategory:          "some-meter-category",
				MeterSubCategory:       "some-meter-subcategory",
				MeterRegion:            "some-meter-region",
				MeterName:              "some-meter-name",
				ConsumedQuantity:       1.0,
				ResourceRate:           2.0,
				ExtendedCost:           3.0,
				ResourceLocation:       "some-resource-location",
				ConsumedService:        "some-consumed-service",
				InstanceID:             "some-instance-id",
				ServiceInfo1:           "some-service-info",
				ServiceInfo2:           "some-service-info",
				AdditionalInfo:         "some-additional-info",
				Tags:                   "some-tags",
				StoreServiceIdentifier: "some-store-id",
				DepartmentName:         "some-department-name",
				CostCenter:             "some-cost-center",
				UnitOfMeasure:          "some-unit-of-measure",
				ResourceGroup:          "some-resource-group",
			}
			otherUsage = Usage{
				AccountOwnerId:         "some-account-owner-id",
				AccountName:            "some-account-name",
				ServiceAdministratorId: "some-admin-id",
				SubscriptionId:         "some-sub-id",
				SubscriptionGuid:       "some-other-sub-guid",
				SubscriptionName:       "some-sub-name",
				Date:                   "some-date",
				Month:                  2,
				Day:                    3,
				Year:                   1000,
				Product:                "some-product",
				MeterID:                "some-meter-id",
				MeterCategory:          "some-meter-category",
				MeterSubCategory:       "some-meter-subcategory",
				MeterRegion:            "some-meter-region",
				MeterName:              "some-meter-name",
				ConsumedQuantity:       1.0,
				ResourceRate:           2.0,
				ExtendedCost:           3.0,
				ResourceLocation:       "some-resource-location",
				ConsumedService:        "some-consumed-service",
				InstanceID:             "some-instance-id",
				ServiceInfo1:           "some-service-info",
				ServiceInfo2:           "some-service-info",
				AdditionalInfo:         "some-additional-info",
				Tags:                   "some-tags",
				StoreServiceIdentifier: "some-store-id",
				DepartmentName:         "some-department-name",
				CostCenter:             "some-cost-center",
				UnitOfMeasure:          "some-unit-of-measure",
				ResourceGroup:          "some-resource-group",
			}
		})

		It("does not return empty string", func() {
			Expect(usage.Hash()).NotTo(BeEmpty())
		})

		It("returns different hash for different structs", func() {
			Expect(usage.Hash()).NotTo(Equal(otherUsage.Hash()))
		})

		It("returns the same hash each time", func() {
			Expect(usage.Hash()).To(Equal(usage.Hash()))
		})
	})
})
