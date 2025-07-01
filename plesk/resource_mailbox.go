package plesk

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func ResourceMailbox() *schema.Resource {
    return &schema.Resource{
        CreateContext: resourceMailboxCreate,
        ReadContext:   resourceMailboxRead,
        DeleteContext: resourceMailboxDelete,
        Schema: map[string]*schema.Schema{
            "email": {
                Type:     schema.TypeString,
                Required: true,
                ForceNew: true,
            },
            "password": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
            },
        },
    }
}

func resourceMailboxCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    payload := map[string]interface{}{
        "email":    d.Get("email"),
        "password": d.Get("password"),
    }

    respBody, diags := client.Post(ctx, "/api/v2/mail", payload)
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
    } else if email, ok := respData["email"].(string); ok {
        d.SetId(email)
    } else {
        d.SetId(d.Get("email").(string))
    }

    return resourceMailboxRead(ctx, d, m)
}

func resourceMailboxRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)

    respBody, diags := client.Get(ctx, "/api/v2/mail")
    if diags.HasError() {
        return diags
    }

    var response struct {
        Mailboxes []struct {
            ID    string `json:"id"`
            Email string `json:"email"`
        } `json:"mailboxes"`
    }

    if err := json.Unmarshal(respBody, &response); err != nil {
        return diag.Diagnostics{{
            Severity: diag.Error,
            Summary:  "Failed to parse mailboxes JSON response",
            Detail:   err.Error(),
        }}
    }

    for _, mailbox := range response.Mailboxes {
        if mailbox.ID == d.Id() || mailbox.Email == d.Id() {
            d.Set("email", mailbox.Email)
            return nil
        }
    }

    d.SetId("")
    return nil
}

func resourceMailboxDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*Client)
    path := fmt.Sprintf("/api/v2/mail/%s", d.Id())
    _, diags := client.Post(ctx, path, nil) // Use DELETE if supported
    if diags.HasError() {
        return diags
    }
    d.SetId("")
    return nil
}
