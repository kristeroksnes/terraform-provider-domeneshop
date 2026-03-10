package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/innovationnorway/go-domeneshop/api/v0/domeneshop"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Description: "Use this resource to create and manage DNS records.",

		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRecordImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"host": {
				Description: "The host/subdomain the DNS record applies to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"type": {
				Description: "The type of the record. Possible values are: `A`, `AAAA`, `ANAME`, `CNAME`, `DS`, `MX`, `NS`, `SRV`, `TXT`, `TLSA`.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"data": {
				Description: "The value of the record.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ttl": {
				Description: "TTL of DNS record in seconds.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"priority": {
				Description: "MX/SRV record priority, also known as preference. Lower values are usually preferred first, but this is not guaranteed",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"weight": {
				Description: "SRV record weight. Relevant if multiple records have same preference.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"port": {
				Description: "SRV record port. The port where the service is found.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"flags": {
				Description: "CAA record flags.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"tag": {
				Description: "CAA/DS record tag.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"alg": {
				Description: "DS record algorithm.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"digest": {
				Description: "DS record digest type.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"usage": {
				Description: "TLSA record certificate usage.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"selector": {
				Description: "TLSA record selector.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"dtype": {
				Description: "TLSA record matching type.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	record := domeneshop.DNSRecord{
		Host: d.Get("host").(string),
		Type: d.Get("type").(string),
		Data: d.Get("data").(string),
	}

	if v, ok := d.GetOk("ttl"); ok {
		record.SetTtl(int32(v.(int)))
	}

	if v, ok := d.GetOk("priority"); ok {
		record.SetPriority(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("weight"); ok {
		record.SetWeight(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("port"); ok {
		record.SetPort(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("flags"); ok {
		record.SetFlags(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("tag"); ok {
		record.SetTag(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("alg"); ok {
		record.SetAlg(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("digest"); ok {
		record.SetDigest(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("usage"); ok {
		record.SetUsage(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("selector"); ok {
		record.SetSelector(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("dtype"); ok {
		record.SetDtype(strconv.Itoa(v.(int)))
	}

	resp, _, err := client.DnsApi.CreateRecord(auth, domainID).DNSRecord(record).Execute()
	if err.Error() != "" {
		return diag.Errorf("error creating DNS record: %s %s", err.Error(), err.Body())
	}

	d.SetId(strconv.Itoa(int(resp.GetId())))

	return resourceRecordRead(ctx, d, meta)
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	recordID, _ := strconv.Atoi(d.Id())

	resp, r, err := client.DnsApi.GetRecord(auth, domainID, int32(recordID)).Execute()
	if err.Error() != "" {
		if r.StatusCode == 404 {
			d.SetId("")
			return nil
		}

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

	return nil
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	recordID, _ := strconv.Atoi(d.Id())

	record := domeneshop.DNSRecord{
		Host: d.Get("host").(string),
		Type: d.Get("type").(string),
		Data: d.Get("data").(string),
	}

	if v, ok := d.GetOk("ttl"); ok {
		record.SetTtl(int32(v.(int)))
	}

	if v, ok := d.GetOk("priority"); ok {
		record.SetPriority(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("weight"); ok {
		record.SetWeight(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("port"); ok {
		record.SetPort(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("flags"); ok {
		record.SetFlags(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("tag"); ok {
		record.SetTag(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("alg"); ok {
		record.SetAlg(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("digest"); ok {
		record.SetDigest(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("usage"); ok {
		record.SetUsage(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("selector"); ok {
		record.SetSelector(strconv.Itoa(v.(int)))
	}

	if v, ok := d.GetOk("dtype"); ok {
		record.SetDtype(strconv.Itoa(v.(int)))
	}

	_, err := client.DnsApi.ModifyRecord(auth, domainID, int32(recordID)).DNSRecord(record).Execute()
	if err.Error() != "" {
		return diag.Errorf("error updating DNS record (ID: %d): %s %s", recordID, err.Error(), err.Body())
	}

	return resourceRecordRead(ctx, d, meta)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))
	recordID, _ := strconv.Atoi(d.Id())

	_, err := client.DnsApi.DeleteRecord(auth, domainID, int32(recordID)).Execute()
	if err.Error() != "" {
		return diag.Errorf("error deleting DNS record (ID: %d): %s %s", recordID, err.Error(), err.Body())
	}

	return nil
}

func resourceRecordImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("error importing DNS record: Expected domain_id/record_id, got %s", d.Id())
	}

	domainID, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, fmt.Errorf("error importing DNS record: Expected domain_id to be integer, got %s", s[0])
	}

	if err := d.Set("domain_id", domainID); err != nil {
		return nil, fmt.Errorf("error setting domain_id: %s", err)
	}
	d.SetId(s[1])

	return []*schema.ResourceData{d}, nil
}
