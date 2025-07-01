package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceSite() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceSiteCreate,
        ReadContext:   resourceSiteRead,
        DeleteContext: resourceSiteDelete,
        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "hosting_type": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "ftp_login": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "ftp_password": {
                Type:      schema.TypeString,
                Optional:  true,
                Sensitive: true,
            },
        },
    }
}

func resourceSiteCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    payload := map[string]interface{}{
        "name": d.Get("name"),
    }
    if v, ok := d.GetOk("hosting_type"); ok {
        payload["hosting_type"] = v
    }
    if v, ok := d.GetOk("ftp_login"); ok {
        payload["ftp_login"] = v
    }
    if v, ok := d.GetOk("ftp_password"); ok {
        payload["ftp_password"] = v
    }

    respBody, diags := client.Post(ctx, "/api/v2/domains", payload)
    if diags.HasError() {
        return diags
    }

    var respData map[string]interface{}
    if err := json.Unmarshal(respBody, &respData); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse response JSON",
            Detail:   err.Error(),
        }}
    }

    // ID can be in "id" field or fallback to name
    if id, ok := respData["id"].(string); ok {
        d.SetId(id)
    } else {
        d.SetId(d.Get("name").(string))
    }

    return resourceSiteRead(ctx, d, m)
}

func resourceSiteRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, "/api/v2/domains")
    if diags.HasError() {
        return diags
    }

    var response struct {
        Domains []struct {
            ID          string `json:"id"`
            Name        string `json:"name"`
            HostingType string `json:"hosting_type,omitempty"`
            FTPLogin    string `json:"ftp_login,omitempty"`
        } `json:"domains"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse domains JSON response",
            Detail:   err.Error(),
        }}
    }

    for _, domain := range response.Domains {
        if domain.ID == d.Id() {
            d.Set("name", domain.Name)
            d.Set("hosting_type", domain.HostingType)
            d.Set("ftp_login", domain.FTPLogin)
            return nil
        }
    }

    d.SetId("")
    return nil
}

func resourceSiteDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    path := fmt.Sprintf("/api/v2/domains/%s", d.Id())
    _, diags := client.Post(ctx, path, nil) // Change to DELETE if API supports it
    if diags.HasError() {
        return diags
    }
    d.SetId("")
    return nil
}
