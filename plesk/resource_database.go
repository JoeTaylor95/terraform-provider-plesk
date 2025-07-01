package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDatabase() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceDatabaseCreate,
        ReadContext:   resourceDatabaseRead,
        UpdateContext: resourceDatabaseUpdate,
        DeleteContext: resourceDatabaseDelete,
        Schema: map[string]*schema.Schema{
            "name": {
                Type:         schema.TypeString,
                Required:     true,
                ForceNew:     true,
                ValidateFunc: validation.StringLenBetween(1, 64),
                Description:  "The name of the database.",
            },
            "type": {
                Type:         schema.TypeString,
                Optional:     true,
                Default:      "mysql",
                ValidateFunc: validation.StringInSlice([]string{"mysql", "pgsql"}, false),
                Description:  "The database type (mysql or pgsql).",
            },
            "server_id": {
                Type:        schema.TypeInt,
                Optional:    true,
                Description: "ID of the database server to use (optional).",
            },
        },
    }
}

func resourceDatabaseCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    payload := map[string]interface{}{
        "name": d.Get("name").(string),
        "type": d.Get("type").(string),
    }

    if v, ok := d.GetOk("server_id"); ok {
        payload["server_id"] = v.(int)
    }

    respBody, diags := client.Post(ctx, "/api/v2/databases", payload)
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
    return resourceDatabaseRead(ctx, d, m)
}

func resourceDatabaseRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, fmt.Sprintf("/api/v2/databases/%s", d.Id()))
    if diags.HasError() {
        return diags
    }

    var db struct {
        ID       string `json:"id"`
        Name     string `json:"name"`
        Type     string `json:"type"`
        ServerID int    `json:"server_id"`
    }
    if err := json.Unmarshal(respBody, &db); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse response",
            Detail:   err.Error(),
        }}
    }

    d.Set("name", db.Name)
    d.Set("type", db.Type)
    d.Set("server_id", db.ServerID)

    return nil
}

func resourceDatabaseUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    payload := map[string]interface{}{}

    if d.HasChange("name") {
        payload["name"] = d.Get("name").(string)
    }

    if d.HasChange("type") {
        payload["type"] = d.Get("type").(string)
    }

    if d.HasChange("server_id") {
        payload["server_id"] = d.Get("server_id").(int)
    }

    if len(payload) == 0 {
        return nil
    }

    _, diags := client.Post(ctx, fmt.Sprintf("/api/v2/databases/%s", d.Id()), payload)
    return diags
}

func resourceDatabaseDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    _, diags := client.Post(ctx, fmt.Sprintf("/api/v2/databases/%s", d.Id()), nil)
    return diags
}
