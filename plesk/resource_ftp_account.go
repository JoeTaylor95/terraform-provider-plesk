package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceFTPAccount() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceFTPAccountCreate,
        ReadContext:   resourceFTPAccountRead,
        DeleteContext: resourceFTPAccountDelete,
        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "password": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
            },
            "home_dir": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceFTPAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    payload := map[string]interface{}{
        "name":     d.Get("name"),
        "password": d.Get("password"),
    }

    if v, ok := d.GetOk("home_dir"); ok {
        payload["home_dir"] = v
    }

    respBody, diags := client.Post(ctx, "/api/v2/ftpusers", payload)
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

    if id, ok := respData["name"].(string); ok {
        d.SetId(id)
    } else {
        d.SetId(d.Get("name").(string))
    }

    return resourceFTPAccountRead(ctx, d, m)
}

func resourceFTPAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, "/api/v2/ftpusers")
    if diags.HasError() {
        return diags
    }

    var response struct {
        FTPUsers []struct {
            Name    string `json:"name"`
            HomeDir string `json:"home_dir,omitempty"`
        } `json:"ftpusers"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse ftpusers JSON response",
            Detail:   err.Error(),
        }}
    }

    for _, ftp := range response.FTPUsers {
        if ftp.Name == d.Id() {
            d.Set("name", ftp.Name)
            d.Set("home_dir", ftp.HomeDir)
            return nil
        }
    }

    d.SetId("")
    return nil
}

func resourceFTPAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    path := fmt.Sprintf("/api/v2/ftpusers/%s", d.Id())
    _, diags := client.Post(ctx, path, nil) // Use DELETE if supported
    if diags.HasError() {
        return diags
    }
    d.SetId("")
    return nil
}
