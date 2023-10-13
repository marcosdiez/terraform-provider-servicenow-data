package models

import (
	// "encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ServiceCatalogTableName = "sc_req_item"

var ServiceCatalogOrderColumns = map[string]*schema.Schema{
	// https://developer.servicenow.com/dev.do#!/reference/api/utah/rest/c_ServiceCatalogAPI#servicecat-PUT-items-submit_guide?navFilter=servicecatalog/items
	"sc_cat_item_sys_id": {
		Type:         schema.TypeString,
		ValidateFunc: validation.StringLenBetween(32, 32),
		Required:     true,
		ForceNew:     true,
	},
	"sc_req_item_raw_data": {
		Description: "Data from the sc_rec_item table",
		Computed:    true,
		Type:        schema.TypeMap,
		Elem: &schema.Schema{
			Type: schema.TypeString},
	},
}

var ServiceCatalogOrderRequestColumns = map[string]*schema.Schema{
	// REQUEST
	"sysparm_also_request_for": {
		Description: "Comma-separated string of user sys_ids of other users for which to order the specified item. User sys_ids are located in the User [sys_user] table.",
		Type:        schema.TypeString,
		Optional:    true,
		ForceNew:    true,
	},
	"sysparm_quantity": {
		Description: "Quantity of the item. Cannot be a negative number. Data type: Number",
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "1",
		ForceNew:    true,
	},
	"sysparm_requested_for": {
		Description:  "Sys_id of the user for whom to order the specified item. Located in the User [sys_user] table.",
		Type:         schema.TypeString,
		ValidateFunc: validation.StringLenBetween(32, 32),
		Optional:     true,
		ForceNew:     true,
	},
	"variables": {
		Description: "Name-value pairs of all mandatory cart item variables. Mandatory variables are defined on the associated form.",
		Required:    true,
		Type:        schema.TypeMap,
		ForceNew:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString},
	},
	"get_portal_messages": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "true",
		ForceNew: true,
	},
	"sysparm_no_validation": {
		Type:     schema.TypeString,
		Optional: true,
		Default:  "true",
		ForceNew: true,
	},
}

var ServiceCatalogOrderResponseColumns = map[string]*schema.Schema{
	// RESPONSE
	"sys_id": {
		Description: "Sys_id of the order.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"number": {
		Description: "Number of the generated request",
		Type:        schema.TypeString,
		Computed:    true,
	},
	// "$$uiNotification": {
	// 	// Computed: true, # a bug in the terraform library does not allow default values for TypeList, when computer
	// 	Type:     schema.TypeSet,
	// 	Computed: true,
	// 	Elem:     &schema.Schema{Type: schema.TypeString},
	// 	Optional: true,
	// },
	"request_number": {
		Description: "Request number.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"parent_id": {
		Description: "If available, the sys_id of the parent record from which the request is created",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"request_id": {
		Description: "Sys_id of the order request.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"parent_table": {
		Description: "If available, the name of the parent table from which the request is created.",
		Type:        schema.TypeString,
		Computed:    true,
	},
	"table": {
		Description: "Table name of the request.",
		Type:        schema.TypeString,
		Computed:    true,
	},
}
