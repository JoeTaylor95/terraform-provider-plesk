package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDatabaseUser() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceDatabaseUserCreate,
        ReadContext:   resourceDatabaseUserRead,
        UpdateContext: resourceDatabaseUserUpdate,
        DeleteContext: resourceDatabaseUserDelete,
        Schema: map[string]*schema.Schema{
            "username": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                ValidateFunc: validation.StringLenBetween(1, 32),
                Description:  "Database user name.",
            },
            "password": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
                Description: "Database user password.",
            },
            "database_id": {
                Type:     schema.TypeString,
                Required: true,
                Description: "ID of the database the user belongs to.",
            },
        },
    }
}

func resourceDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    payload := map[string]interface{}{
        "username":    d.Get("username").(string),
        "password":    d.Get("password").(string),
        "database_id": d.Get("database_id").(string),
    }

    respBody, diags := client.Post(ctx, "/api/v2/dbusers", payload)
    if diags.HasError() {
        return diags
    }

    var resp struct {
        ID string `json:"id"`
    }
    if err := json.Unmarshal(respBody, &resp); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse response",
            Detail:   err.Error(),
        }}
    }

    d.SetId(resp.ID)
    return resourceDatabaseUserRead(ctx, d, m)
}

func resourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, fmt.Sprintf("/api/v2/dbusers/%s", d.Id()))
    if diags.HasError() {
        return diags
    }

    var user struct {
        ID         string `json:"id"`
        Username   string `json:"username"`
        DatabaseID string `json:"database_id"`
    }
    if err := json.Unmarshal(respBody, &user); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse response",
            Detail:   err.Error(),
        }}
    }

    d.Set("username", user.Username)
    d.Set("database_id", user.DatabaseID)
    // Note: password is sensitive, generally not retrievable, so do not set

    return nil
}

func resourceDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    payload := map[string]interface{}{}

    if d.HasChange("password") {
        payload["password"] = d.Get("password").(string)
    }

    if len(payload) == 0 {
        return nil
    }

    _, diags := client.Post(ctx, fmt.Sprintf("/api/v2/dbusers/%s", d.Id()), payload)
    return diags
}

func resourceDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    _, diags := client.Post(ctx, fmt.Sprintf("/api/v2/dbusers/%s", d.Id()), nil)
    return diags
}
