package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRecords() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve a list of DNS records.",

		ReadContext: dataSourceRecordsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "Only return domains whose `domain` field includes this string.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"host": {
				Description: "Only return records whose `host` field matches this string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description: "Only return records whose `type` field matches this string.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"records": {
				Description: "List of records.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "ID of the DNS record.",
							Type:        schema.TypeInt,
							Computed:    true,
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
				},
			},
		},
	}
}

func dataSourceRecordsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	host := d.Get("host").(string)
	recordType := d.Get("type").(string)

	resp, _, err := client.DnsApi.GetRecords(auth, domainID).Host(host).Type_(recordType).Execute()
	if err.Error() != "" {
		return diag.Errorf("error getting domains: %s %s", err.Error(), err.Body())
	}

	var records []interface{}
	for _, v := range resp {
		priority, _ := strconv.Atoi(v.GetPriority())
		weight, _ := strconv.Atoi(v.GetWeight())
		port, _ := strconv.Atoi(v.GetPort())
		flags, _ := strconv.Atoi(v.GetFlags())
		alg, _ := strconv.Atoi(v.GetAlg())
		digest, _ := strconv.Atoi(v.GetDigest())
		usage, _ := strconv.Atoi(v.GetUsage())
		selector, _ := strconv.Atoi(v.GetSelector())
		dtype, _ := strconv.Atoi(v.GetDtype())
		records = append(records, map[string]interface{}{
			"id":       v.GetId(),
			"host":     v.GetHost(),
			"ttl":      v.GetTtl(),
			"type":     v.GetType(),
			"data":     v.GetData(),
			"priority": priority,
			"weight":   weight,
			"port":     port,
			"flags":    flags,
			"alg":      alg,
			"digest":   digest,
			"usage":    usage,
			"selector": selector,
			"dtype":    dtype,
		})
	}

	if err := d.Set("records", records); err != nil {
		return diag.FromErr(err)
	}

	id := fmt.Sprintf("%v", records)
	d.SetId(fmt.Sprintf("%d", schema.HashString(id)))

	return nil
}
