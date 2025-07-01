package plesk

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ProviderConfigure initializes and returns the Plesk API client based on
// the configuration from the Terraform provider block.
func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
    var diags diag.Diagnostics

    host := d.Get("host").(string)
    port := d.Get("port").(string)
    token := d.Get("token").(string)

    // Basic validation
    if host == "" {
        diags = append(diags, diag.Diagnostic{
            Severity: diag.Error,
            Summary:  "Missing Host",
            Detail:   "The Plesk host must be provided.",
        })
        return nil, diags
    }

    if token == "" {
        diags = append(diags, diag.Diagnostic{
            Severity: diag.Error,
            Summary:  "Missing API Token",
            Detail:   "The Plesk API token must be provided.",
        })
        return nil, diags
    }

    client := &Client{
        Host:  host,
        Port:  port,
        Token: token,
        // Use default HTTP client; could add timeout, TLS config etc.
        Client: nil,
    }

    return client, diags
}
