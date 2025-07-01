package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceExtension() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceExtensionInstall,
        ReadContext:   resourceExtensionRead,
        DeleteContext: resourceExtensionUninstall,
        Schema: map[string]*schema.Schema{
            "id": {
                Type:        schema.TypeString,
                Required:    true,
                ForceNew:    true,
                Description: "Extension identifier.",
            },
            "enabled": {
                Type:        schema.TypeBool,
                Optional:    true,
                Default:     true,
                Description: "Whether the extension is enabled.",
            },
        },
    }
}

func resourceExtensionInstall(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    extensionID := d.Get("id").(string)

    payload := map[string]interface{}{
        "id": extensionID,
    }

    _, diags := client.Post(ctx, "/api/v2/extensions", payload)
    if diags.HasError() {
        return diags
    }

    d.SetId(extensionID)

    // Enable if requested
    if d.Get("enabled").(bool) {
        _, diags = client.Post(ctx, fmt.Sprintf("/api/v2/extensions/%s/enable", extensionID), nil)
        if diags.HasError() {
            return diags
        }
    }

    return resourceExtensionRead(ctx, d, m)
}

func resourceExtensionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    extensionID := d.Id()

    respBody, diags := client.Get(ctx, fmt.Sprintf("/api/v2/extensions/%s", extensionID))
    if diags.HasError() {
        return diags
    }

    var ext struct {
        ID      string `json:"id"`
        Enabled bool   `json:"enabled"`
    }
    if err := json.Unmarshal(respBody, &ext); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse extension response",
            Detail:   err.Error(),
        }}
    }

    d.Set("id", ext.ID)
    d.Set("enabled", ext.Enabled)

    return nil
}

func resourceExtensionUninstall(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    extensionID := d.Id()

    _, diags := client.Post(ctx, fmt.Sprintf("/api/v2/extensions/%s", extensionID), nil) // Assuming DELETE via POST or change to DELETE
    if diags.HasError() {
        return diags
    }

    d.SetId("")
    return nil
}
