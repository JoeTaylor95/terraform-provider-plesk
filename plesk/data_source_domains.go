package plesk

import (
    "context"
    "encoding/json"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDomains() *schema.Resource {
    return &schema.Resource{
        ReadContext: dataSourceDomainsRead,
        Schema: map[string]*schema.Schema{
            "domains": {
                Type:     schema.TypeList,
                Computed: true,
                Elem: &schema.Resource{
                    Schema: map[string]*schema.Schema{
                        "id": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                        "name": {
                            Type:     schema.TypeString,
                            Computed: true,
                        },
                    },
                },
            },
        },
    }
}

func dataSourceDomainsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, "/api/v2/domains")
    if diags.HasError() {
        return diags
    }

    var domainsResp struct {
        Domains []struct {
            ID   string `json:"id"`
            Name string `json:"name"`
        } `json:"domains"`
    }
    if err := json.Unmarshal(respBody, &domainsResp); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse domains response",
            Detail:   err.Error(),
        }}
    }

    domains := make([]map[string]interface{}, 0, len(domainsResp.Domains))
    for _, domain := range domainsResp.Domains {
        domains = append(domains, map[string]interface{}{
            "id":   domain.ID,
            "name": domain.Name,
        })
    }

    d.SetId("plesk-domains") // static ID for the data source instance
    d.Set("domains", domains)

    return nil
}
