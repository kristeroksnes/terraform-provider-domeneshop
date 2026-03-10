package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceForwards() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve a list of HTTP forwards.",

		ReadContext: dataSourceForwardsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"forwards": {
				Description: "List of HTTP forwards.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Description: "The subdomain this forward applies to, without the domain part.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"url": {
							Description: "The URL to forward to. Must include scheme, e.g. `https://` or `ftp://`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"frame": {
							Description: "Whether to enable frame forwarding using an iframe embed. NOT recommended for a variety of reasons.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceForwardsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	resp, _, err := client.ForwardsApi.GetForwards(auth, domainID).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting HTTP forwards: %s %s", err.Error(), err.Body())
	}

	var forwards []interface{}
	for _, v := range resp {
		forwards = append(forwards, map[string]interface{}{
			"host":  v.GetHost(),
			"url":   v.GetUrl(),
			"frame": v.GetFrame(),
		})
	}

	if err := d.Set("forwards", forwards); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%v", forwards)
	d.SetId(fmt.Sprintf("%d", schema.HashString(id)))

	return nil
}
