package provider

import (
	"context"
	"net/http"
	"time"

	"github.com/JoeTaylor95/terraform-provider-plesk/plesk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Plesk server hostname or IP",
				DefaultFunc: schema.EnvDefaultFunc("PLESK_HOST", nil),
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Plesk API port",
				Default:     "8443",
				DefaultFunc: schema.EnvDefaultFunc("PLESK_PORT", "8443"),
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Bearer token to authenticate to the Plesk API",
				DefaultFunc: schema.EnvDefaultFunc("PLESK_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"plesk_site":          plesk.ResourceSite(),
			"plesk_ftp_account":   plesk.ResourceFTPAccount(),
			"plesk_user":          plesk.ResourceUser(),
			"plesk_reseller":      plesk.ResourceReseller(),
			"plesk_mailbox":       plesk.ResourceMailbox(),
			"plesk_database":      plesk.ResourceDatabase(),
			"plesk_database_user": plesk.ResourceDatabaseUser(),
			"plesk_extension":     plesk.ResourceExtension(),
			"plesk_dns_record":    plesk.ResourceDnsRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"plesk_domains": plesk.DataSourceDomains(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	host := d.Get("host").(string)
	port := d.Get("port").(string)
	token := d.Get("token").(string)

	if host == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Plesk host must be provided",
		})
		return nil, diags
	}

	if token == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Plesk API token must be provided",
		})
		return nil, diags
	}

	client := &plesk.Client{
		Host:  host,
		Port:  port,
		Token: token,
		Client: &http.Client{ // Capital C here!
			Timeout: 15 * time.Second,
		},
	}

	// Test API connectivity and authentication
	respBody, errDiags := client.Get(ctx, "/api/v2/server")
	if errDiags.HasError() {
		diags = append(diags, errDiags...)
		return nil, diags
	}

	if len(respBody) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Empty response from Plesk server API",
		})
		return nil, diags
	}

	return client, diags
}
