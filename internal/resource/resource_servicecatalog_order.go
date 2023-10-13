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

func sendCreatePayload(ctx context.Context, url string, payload map[string]interface{}, client *client.Client) (map[string]any, diag.Diagnostics) {
	rawData, err := client.SendRequest(http.MethodPost, url, payload, 200)
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

func saveCreatePayload(data_result map[string]any, data *schema.ResourceData) diag.Diagnostics {
	for k, v := range data_result {
		if k == "$$uiNotification" {
			continue
		}
		if err := data.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func serviceCatalogOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	url := fmt.Sprintf("/api/sn_sc/v1/servicecatalog/items/%s/order_now", d.Get("sc_cat_item_sys_id"))
	payload := prepareCreatePayload(d)
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: url=%s, params=%s", url, payload))

	client := m.(*client.Client)
	data, diags := sendCreatePayload(ctx, url, payload, client)
	if diags != nil {
		return diags
	}

	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: result=%s", data))
	data_result := data["result"].(map[string]any)

	diags = saveCreatePayload(data_result, d)
	if diags != nil {
		return diags
	}

	sys_id := data_result["sys_id"].(string)

	sc_req_item_raw_data, err := serviceCatalogOrderReadTableData(sys_id, client)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderCreate: result=%s", sc_req_item_raw_data.RowData))

	if err := d.Set("sc_req_item_raw_data", sc_req_item_raw_data.RowData); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(sys_id)
	return nil
}

func serviceCatalogOrderReadTableData(sys_id string, client *client.Client) (*models.ParsedResult, error) {
	query := map[string]interface{}{"request": sys_id}
	rowData, err := client.GetTableRow(models.ServiceCatalogTableName, query)
	if err != nil {
		return nil, err
	}
	return rowData, nil
}

func serviceCatalogOrderRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {

	sys_id := data.Id()
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderRead: sys_id=%s", sys_id))
	client := m.(*client.Client)

	sc_req_item_raw_data, err := serviceCatalogOrderReadTableData(sys_id, client)
	if err != nil {
		return diag.FromErr(err)
	}
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderRead: result=%s", sc_req_item_raw_data.RowData))

	if err := data.Set("sc_req_item_raw_data", sc_req_item_raw_data.RowData); err != nil {
		return diag.FromErr(err)
	}
	var diags diag.Diagnostics
	return diags
}

func serviceCatalogOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("serviceCatalogOrderDelete: deleting is not possible. This is a NO OP by design. data=%s", d))
	return nil
}
