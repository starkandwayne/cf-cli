package resources

import (
	"fmt"
	"github.com/starkandwayne/cf-cli/cf/models"
)

type ServicePlanResource struct {
	Resource
	Entity ServicePlanEntity
}

type ServicePlanEntity struct {
	Name                string
	ServiceOfferingGuid string                  `json:"service_guid"`
	ServiceOffering     ServiceOfferingResource `json:"service"`
}

type ServicePlanDescription struct {
	ServiceLabel    string
	ServicePlanName string
	ServiceProvider string
}

func (resource ServicePlanResource) ToFields() (fields models.ServicePlanFields) {
	fields.Guid = resource.Metadata.Guid
	fields.Name = resource.Entity.Name
	return
}

func (planDesc ServicePlanDescription) String() string {
	if planDesc.ServiceProvider == "" {
		return fmt.Sprintf("%s %s", planDesc.ServiceLabel, planDesc.ServicePlanName) // v2 plan
	} else {
		return fmt.Sprintf("%s %s %s", planDesc.ServiceLabel, planDesc.ServiceProvider, planDesc.ServicePlanName) // v1 plan
	}
}

type ServiceMigrateV1ToV2Response struct {
	ChangedCount int `json:"changed_count"`
}
