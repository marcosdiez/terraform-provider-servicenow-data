package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/models"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ServiceCatalogOrderResourceName = "servicenow-data_servicecatalog_order"

func ServiceCatalogOrderResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: serviceCatalogOrderCreate,
		ReadContext:   serviceCatalogOrderRead,
		// UpdateContext:  serviceCatalogOrderUpdate, SNOW catalog orders are immutable
		DeleteContext: serviceCatalogOrderDelete,
		Schema: *models.MergeSchema(
			models.ServiceCatalogOrderColumns,
			*models.MergeSchema(
				models.ServiceCatalogOrderRequestColumns,
				models.ServiceCatalogOrderResponseColumns,
			),
		),
		SchemaVersion:  1,
		StateUpgraders: nil,
		CustomizeDiff:  nil,
		Importer: &schema.ResourceImporter{
			StateContext: nil,
		},
		DeprecationMessage: "",
		Timeouts: &schema.ResourceTimeout{
			Create:  schema.DefaultTimeout(5 * time.Second),
			Read:    schema.DefaultTimeout(5 * time.Second),
			Update:  schema.DefaultTimeout(5 * time.Second),
			Delete:  schema.DefaultTimeout(5 * time.Second),
			Default: schema.DefaultTimeout(5 * time.Second),
		},
		Description:   "A ServiceCatalog order",
		UseJSONNumber: false,
	}
}

func prepareCreatePayload(d *schema.ResourceData) map[string]interface{} {
	payload := map[string]interface{}{}
	for k, _ := range models.ServiceCatalogOrderRequestColumns {
		payload[k] = d.Get(k)
	}
	return payload
}

func sendCreatePayload(ctx context.Context, url string, payload map[string]interface{}, m interface{}) (map[string]any, diag.Diagnostics) {
	c := m.(*client.Client)
	rawData, err := c.SendRequest(http.MethodPost, url, payload, 200)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error on SendRequest: %s", err))
	}

	var data map[string]any
	err = json.Unmarshal(*rawData, &data)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("Error on json.Unmarsha: %s", err))
	}
	return data, nil
}

func saveCreatePayload(data_result map[string]any, d *schema.ResourceData) diag.Diagnostics {
	for k, v := range data_result {
		if k == "$$uiNotification" {
			continue
		}
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func serviceCatalogOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	url := fmt.Sprintf("/api/sn_sc/v1/servicecatalog/items/%s/order_now", d.Get("sc_cat_item_sys_id"))
	payload := prepareCreatePayload(d)
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: url=%s, params=%s", url, payload))

	data, diags := sendCreatePayload(ctx, url, payload, m)
	if diags != nil {
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: result=%s", data))
	data_result := data["result"].(map[string]any)

	diags = saveCreatePayload(data_result, d)
	if diags != nil {
		return diags
	}

	d.SetId(data_result["sys_id"].(string))
	return nil
}

func serviceCatalogOrderRead(_ context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	// c := m.(*client.Client)
	// Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
	return diag.FromErr(fmt.Errorf("serviceCatalogOrderRead has not been implemented yet"))

	// c := m.(*client.Client)
	// var tableID, sysID string
	// var err error
	// // Warning or errors can be collected in a slice type
	// var diags diag.Diagnostics
	// tableID, sysID, err = ExtractIDs(data.Id())
	// if err != nil {
	// 	return append(diags, diag.FromErr(err)...)
	// }

	// rowData, err := c.GetServiceCatalogOrder(tableID, map[string]interface{}{"sys_id": sysID})
	// if err != nil {
	// 	return diag.FromErr(err)
	// }
	// if len(rowData.SysData) == 0 {
	// 	data.SetId("")
	// } else {
	// 	parseResultDiag := ParsedResultToSchema(data, rowData)
	// 	diags = append(diags, parseResultDiag...)
	// 	if len(parseResultDiag) > 0 {
	// 		return diags
	// 	}
	// 	data.SetId(fmt.Sprintf("%s/%s", tableID, sysID))
	// }
	// return diags
}

func serviceCatalogOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderDelete: deleting is not possible. This is a NO OP by design. data=%s", d))
	return nil
}

// func ExtractIDs(ID string) (tableID, sysID string, err error) {
// 	ids := strings.Split(ID, `/`)
// 	if len(ids) != 2 {
// 		return "", "", fmt.Errorf("faulty id!%s", ID)
// 	}
// 	return ids[0], ids[1], nil
// }

// func ParsedResultToSchema(d *schema.ResourceData, result *models.ParsedResult) diag.Diagnostics {
// 	for k, v := range result.SysData {
// 		if err := d.Set(k, v); err != nil {
// 			return diag.FromErr(err)
// 		}
// 	}
// 	if err := d.Set("row_data", result.RowData); err != nil {
// 		return diag.FromErr(err)
// 	}
// 	return nil
// }
