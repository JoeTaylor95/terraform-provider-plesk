package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceReseller() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceResellerCreate,
        ReadContext:   resourceResellerRead,
        DeleteContext: resourceResellerDelete,
        Schema: map[string]*schema.Schema{
            "login": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
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

func resourceResellerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

    respBody, diags := client.Post(ctx, "/api/v2/resellers", payload)
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
    } else if login, ok := respData["login"].(string); ok {
        d.SetId(login)
    } else {
        d.SetId(d.Get("login").(string))
    }

    return resourceResellerRead(ctx, d, m)
}

func resourceResellerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, "/api/v2/resellers")
    if diags.HasError() {
        return diags
    }

    var response struct {
        Resellers []struct {
            ID    string `json:"id"`
            Login string `json:"login"`
            Email string `json:"email,omitempty"`
        } `json:"resellers"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse resellers JSON response",
            Detail:   err.Error(),
        }}
    }

    for _, reseller := range response.Resellers {
        if reseller.ID == d.Id() || reseller.Login == d.Id() {
            d.Set("login", reseller.Login)
            d.Set("email", reseller.Email)
            return nil
        }
    }

    d.SetId("")
    return nil
}

func resourceResellerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    path := fmt.Sprintf("/api/v2/resellers/%s", d.Id())
    _, diags := client.Post(ctx, path, nil) // Use DELETE if supported
    if diags.HasError() {
        return diags
    }
    d.SetId("")
    return nil
}
