package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceForward() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about a HTTP forwarding.",

		ReadContext: dataSourceForwardRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"host": {
				Description: "The subdomain this forward applies to, without the domain part.",
				Type:        schema.TypeString,
				Required:    true,
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
	}
}

func dataSourceForwardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	host := d.Get("host").(string)

	resp, _, err := client.ForwardsApi.GetForward(auth, domainID, host).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting HTTP forward (Host: %s): %s %s", host, err.Error(), err.Body())
	}

	if err := d.Set("host", resp.GetHost()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", resp.GetUrl()); err != nil {
		return diag.FromErr(err)
	}

	if resp.HasFrame() {
		if err := d.Set("frame", resp.GetFrame()); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(resp.GetHost())

	return nil
}
