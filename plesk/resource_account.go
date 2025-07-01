package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceAccount() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceAccountCreate,
        ReadContext:   resourceAccountRead,
        DeleteContext: resourceAccountDelete,
        Schema: map[string]*schema.Schema{
            "login": {
                Type:     schema.TypeString,
                Required: true,
            },
            "password": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
            },
            "company": {
                Type:     schema.TypeString,
                Optional: true,
            },
            "email": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceAccountCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    payload := map[string]interface{}{
        "login":    d.Get("login"),
        "password": d.Get("password"),
    }

    if v, ok := d.GetOk("company"); ok {
        payload["company"] = v
    }
    if v, ok := d.GetOk("email"); ok {
        payload["email"] = v
    }

    respBody, diags := client.Post(ctx, "/api/v2/accounts", payload)
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

    if id, ok := respData["id"].(string); ok {
        d.SetId(id)
    } else {
        d.SetId(d.Get("login").(string))
    }

    return resourceAccountRead(ctx, d, m)
}

func resourceAccountRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    return nil
}

func resourceAccountDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    path := fmt.Sprintf("/api/v2/accounts/%s", d.Id())
    _, diags := client.Post(ctx, path, nil)
    if diags.HasError() {
        return diags
    }
    d.SetId("")
    return nil
}
