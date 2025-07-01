package plesk

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// ResourceDnsRecord defines the DNS record resource schema and CRUD operations.
func ResourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the domain this DNS record belongs to.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "DNS record type (e.g., A, AAAA, CNAME, MX, TXT, NS, SRV, PTR).",
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "PTR"}, false),
			},
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host or subdomain for the DNS record.",
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value for the DNS record (IP, hostname, etc.).",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     3600,
				Description: "Time To Live (TTL) for the DNS record in seconds.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Priority for MX or SRV records.",
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainID := d.Get("domain_id").(string)

	reqBody := map[string]interface{}{
		"type":  d.Get("type").(string),
		"host":  d.Get("host").(string),
		"value": d.Get("value").(string),
		"ttl":   d.Get("ttl").(int),
	}

	if v, ok := d.GetOk("priority"); ok {
		reqBody["priority"] = v.(int)
	}

	path := fmt.Sprintf("/api/v2/domains/%s/dns/records", domainID)
	respBody, diags := client.Post(ctx, path, reqBody)
	if diags.HasError() {
		return diags
	}

	// Expect response to include new record ID
	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return diag.Errorf("failed to parse DNS record create response: %s", err)
	}

	d.SetId(resp.ID)
	return resourceDnsRecordRead(ctx, d, m)
}

func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainID := d.Get("domain_id").(string)
	recordID := d.Id()

	path := fmt.Sprintf("/api/v2/domains/%s/dns/records/%s", domainID, recordID)
	respBody, diags := client.Get(ctx, path)
	if diags.HasError() {
		return diags
	}

	var resp struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		Host     string `json:"host"`
		Value    string `json:"value"`
		TTL      int    `json:"ttl"`
		Priority int    `json:"priority,omitempty"`
	}
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return diag.Errorf("failed to parse DNS record read response: %s", err)
	}

	d.Set("domain_id", domainID)
	d.Set("type", resp.Type)
	d.Set("host", resp.Host)
	d.Set("value", resp.Value)
	d.Set("ttl", resp.TTL)
	d.Set("priority", resp.Priority)

	return nil
}

func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainID := d.Get("domain_id").(string)
	recordID := d.Id()

	reqBody := map[string]interface{}{
		"type":  d.Get("type").(string),
		"host":  d.Get("host").(string),
		"value": d.Get("value").(string),
		"ttl":   d.Get("ttl").(int),
	}

	if v, ok := d.GetOk("priority"); ok {
		reqBody["priority"] = v.(int)
	}

	path := fmt.Sprintf("/api/v2/domains/%s/dns/records/%s", domainID, recordID)
	_, diags := client.Put(ctx, path, reqBody)
	if diags.HasError() {
		return diags
	}

	return resourceDnsRecordRead(ctx, d, m)
}

func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	domainID := d.Get("domain_id").(string)
	recordID := d.Id()

	path := fmt.Sprintf("/api/v2/domains/%s/dns/records/%s", domainID, recordID)
	return client.Delete(ctx, path)
}
