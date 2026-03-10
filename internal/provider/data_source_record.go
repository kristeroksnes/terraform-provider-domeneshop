package provider

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about a DNS record.",

		ReadContext: dataSourceRecordRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"record_id": {
				Description: "ID of DNS the record.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"host": {
				Description: "The host/subdomain the DNS record applies to.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ttl": {
				Description: "TTL of DNS record in seconds.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"type": {
				Description: "The type of the record.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"data": {
				Description: "The value of the record.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"priority": {
				Description: "MX/SRV record priority, also known as preference. Lower values are usually preferred first, but this is not guaranteed.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"weight": {
				Description: "SRV record weight. Relevant if multiple records have same preference.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"port": {
				Description: "SRV record port. The port where the service is found.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"flags": {
				Description: "CAA record flags.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"tag": {
				Description: "CAA/DS record tag.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"alg": {
				Description: "DS record algorithm.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"digest": {
				Description: "DS record digest type.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"usage": {
				Description: "TLSA record certificate usage.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"selector": {
				Description: "TLSA record selector.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"dtype": {
				Description: "TLSA record matching type.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	recordID := int32(d.Get("record_id").(int))

	resp, _, err := client.DnsApi.GetRecord(auth, domainID, recordID).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting DNS record (ID: %d): %s %s", recordID, err.Error(), err.Body())
	}

	if err := d.Set("host", resp.GetHost()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", resp.GetType()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("data", resp.GetData()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ttl", resp.GetTtl()); err != nil {
		return diag.FromErr(err)
	}

	if v, ok := resp.GetPriorityOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("priority", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetWeightOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("weight", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetPortOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("port", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetFlagsOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("flags", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetTagOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("tag", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetAlgOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("alg", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetDigestOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("digest", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetUsageOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("usage", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetSelectorOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("selector", i); err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := resp.GetDtypeOk(); ok {
		i, _ := strconv.Atoi(*v)
		if err := d.Set("dtype", i); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(strconv.Itoa(int(recordID)))

	return nil
}
