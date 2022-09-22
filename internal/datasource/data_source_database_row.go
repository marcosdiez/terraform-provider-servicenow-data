package datasource

import (
	"context"
	"fmt"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/client"
	"github.com/BESTSELLER/terraform-provider-servicenow-data/internal/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TODO
func DatabaseRowDatasource() *schema.Resource {
	return &schema.Resource{
		Schema:        resource.RowSchema,
		ReadContext:   dataSourceDatabaseRowRead,
		Description:   "A row in a SN table",
		UseJSONNumber: false,
	}
}

func dataSourceDatabaseRowRead(ctx context.Context, data *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	tableID := data.Get("table_id").(string)
	if tableID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "tableID is mandatory",
		})
		return diags
	}
	rowID := data.Get("row_data").(map[string]string)["sys_id"]

	if rowID == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "sys_id is mandatory",
		})
		return diags
	}

	rowData, err := c.GetTableRow(tableID, rowID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := data.Set("row_data", &rowData); err != nil {
		return diag.FromErr(err)
	}

	data.SetId(fmt.Sprintf("%s/%s", tableID, rowID))

	return diags
}
