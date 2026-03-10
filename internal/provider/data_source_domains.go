package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomains() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve a list of domains.",

		ReadContext: dataSourceDomainsRead,

		Schema: map[string]*schema.Schema{
			"domain": {
				Description: "Only return domains whose `domain` field includes this string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"domains": {
				Description: "List of domains.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the domain.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"domain": {
							Description: "Name of the domain.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"expiry_date": {
							Description: "Expiry date.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"registered_date": {
							Description: "Registered date.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"renew": {
							Description: "Whether the domain should be renewed.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"registrant": {
							Description: "Name of the registrant.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"status": {
							Description: "Domain status.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"nameservers": {
							Description: "List of nameservers.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"services": {
							Description: "Domain services.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"registrar": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"dns": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"email": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"webhotel": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDomainsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domain := d.Get("domain").(string)

	resp, _, err := client.DomainsApi.GetDomains(auth).Domain(domain).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting domains: %s %s", err.Error(), err.Body())
	}

	var domains []interface{}
	for _, v := range resp {
		var services []interface{}
		if v.HasServices() {
			services = append(services, map[string]interface{}{
				"registrar": v.Services.GetRegistrar(),
				"dns":       v.Services.GetDns(),
				"email":     v.Services.GetEmail(),
				"webhotel":  v.Services.GetWebhotel(),
			})
		}
		domains = append(domains, map[string]interface{}{
			"id":              v.GetId(),
			"domain":          v.GetDomain(),
			"expiry_date":     v.GetExpiryDate(),
			"registered_date": v.GetRegisteredDate(),
			"renew":           v.GetRenew(),
			"registrant":      v.GetRegistrant(),
			"status":          v.GetStatus(),
			"nameservers":     v.GetNameservers(),
			"services":        services,
		})
	}

	if err := d.Set("domains", domains); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%v", domains)
	d.SetId(fmt.Sprintf("%d", schema.HashString(id)))

	return nil
}
