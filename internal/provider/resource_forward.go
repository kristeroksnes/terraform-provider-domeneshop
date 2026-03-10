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

func resourceForward() *schema.Resource {
	return &schema.Resource{
		Description: `Use this resource to create and manage HTTP forwards ("WWW forwarding").`,

		CreateContext: resourceForwardCreate,
		ReadContext:   resourceForwardRead,
		UpdateContext: resourceForwardUpdate,
		DeleteContext: resourceForwardDelete,
		Importer: &schema.ResourceImporter{
			State: resourceForwardImport,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Description: "ID of the domain.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"host": {
				Description: "Subdomain of the forward, `@` for the root domain.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"url": {
				Description: "The URL to forward to. Must include scheme, e.g. `https://` or `ftp://`.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"frame": {
				Description: "Whether to enable frame forwarding using an iframe embed. NOT recommended for a variety of reasons.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func resourceForwardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	forward := domeneshop.HTTPForward{
		Host: d.Get("host").(string),
		Url:  d.Get("url").(string),
	}

	if v, ok := d.GetOk("frame"); ok {
		forward.Frame = domeneshop.PtrBool(v.(bool))
	}

	r, err := client.ForwardsApi.CreateForward(auth, domainID).HTTPForward(forward).Execute()
	if err.Error() != "" {
		return diag.Errorf("error creating HTTP forward: %s %s", err.Error(), err.Body())
	}

	// The Location header field contains the location of
	// the the newly created resource.
	// Location: https://api.domeneshop.no/v0/domains/{domainId}/dns/{host}
	v := strings.Split(r.Header.Get("Location"), "/")
	host := v[len(v)-1]

	d.SetId(host)

	return resourceForwardRead(ctx, d, meta)
}

func resourceForwardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	resp, r, err := client.ForwardsApi.GetForward(auth, domainID, d.Id()).Execute()
	if err.Error() != "" {
		if r.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error getting HTTP forward (Host: %s): %s %s", d.Id(), err.Error(), err.Body())
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

	return nil
}

func resourceForwardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	forward := domeneshop.HTTPForward{
		Host: d.Get("host").(string),
		Url:  d.Get("url").(string),
	}

	if v, ok := d.GetOk("frame"); ok {
		forward.Frame = domeneshop.PtrBool(v.(bool))
	}

	_, _, err := client.ForwardsApi.ModifyForward(auth, domainID, d.Id()).HTTPForward(forward).Execute()
	if err.Error() != "" {
		return diag.Errorf("error updating HTTP forward (Host: %s): %s %s", d.Id(), err.Error(), err.Body())
	}

	return resourceForwardRead(ctx, d, meta)
}

func resourceForwardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient).client
	auth := meta.(*apiClient).auth

	domainID := int32(d.Get("domain_id").(int))

	_, err := client.ForwardsApi.DeleteForward(auth, domainID, d.Id()).Execute()
	if err.Error() != "" {
		return diag.Errorf("error deleting HTTP forward (Host: %s): %s %s", d.Id(), err.Error(), err.Body())
	}

	return nil
}

func resourceForwardImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), "/")
	if len(s) != 2 {
		return nil, fmt.Errorf("error importing HTTP forward: Expected domain_id/host, got %s", d.Id())
	}

	domainID, err := strconv.Atoi(s[0])
	if err != nil {
		return nil, fmt.Errorf("error importing HTTP forward: Expected domain_id to be integer, got %s", s[0])
	}

	if err := d.Set("domain_id", domainID); err != nil {
		return nil, fmt.Errorf("error setting domain_id: %s", err)
	}
	if err := d.Set("host", s[1]); err != nil {
		return nil, fmt.Errorf("error setting host: %s", err)
	}
	d.SetId(s[1])

	return []*schema.ResourceData{d}, nil
}
