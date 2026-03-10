package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDomain() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about a domain.",

		ReadContext: dataSourceDomainRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
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
	}
}

func dataSourceDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	resp, _, err := client.DomainsApi.GetDomain(auth, domainID).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting domain (ID: %d): %s %s", domainID, err.Error(), err.Body())
	}

	if err := d.Set("domain", resp.GetDomain()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("expiry_date", resp.GetExpiryDate()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("registered_date", resp.GetRegisteredDate()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("renew", resp.GetRenew()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("registrant", resp.GetRegistrant()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("status", resp.GetStatus()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nameservers", resp.GetNameservers()); err != nil {
		return diag.FromErr(err)
	}

	var services []interface{}
	if resp.HasServices() {
		services = append(services, map[string]interface{}{
			"registrar": resp.Services.GetRegistrar(),
			"dns":       resp.Services.GetDns(),
			"email":     resp.Services.GetEmail(),
			"webhotel":  resp.Services.GetWebhotel(),
		})
	}
	if err := d.Set("services", services); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(int(resp.GetId())))

	return nil
}
