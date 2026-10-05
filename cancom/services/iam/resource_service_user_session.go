package iam

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/cancom/terraform-provider-cancom/client"
	client_iam "github.com/cancom/terraform-provider-cancom/client/services/iam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceServiceUserSession() *schema.Resource {
	return &schema.Resource{
		Description:   "IAM --- Service user session token with automatic re-roll support.",
		CreateContext: resourceServiceUserSessionCreate,
		ReadContext:   resourceServiceUserSessionRead,
		UpdateContext: resourceServiceUserSessionUpdate,
		DeleteContext: resourceServiceUserSessionDelete,
		CustomizeDiff: customdiff.All(
			func(ctx context.Context, diff *schema.ResourceDiff, m interface{}) error {
				if diff.Id() == "" {
					return nil
				}
				exp := int64(diff.Get("expires_at").(int))
				if !sessionNeedsReroll(exp, diff.Get("reroll_days").(int), time.Now()) {
					return nil
				}
				// ForceNew requires a pending change, so mark jwt as unknown first.
				if err := diff.SetNewComputed("jwt"); err != nil {
					return err
				}
				return diff.ForceNew("jwt")
			},
		),
		Schema: map[string]*schema.Schema{
			"service_user_crn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Principal CRN of the service user (e.g. crn:cancom::iam:serviceuser:testuser2).",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment / description for the session.",
			},
			"reroll_days": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Number of days before expiration to trigger a recreation (reroll) of the session token.",
			},
			"jwt": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The JWT session token.",
			},
			"session_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The session ID extracted from the token / API.",
			},
			"expires_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration timestamp (UNIX epoch in seconds) of the session token.",
			},
			"ttl": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "TTL timestamp of the session.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the session is currently enabled.",
			},
		},
	}
}

func sessionNeedsReroll(expiresAt int64, rerollDays int, now time.Time) bool {
	if expiresAt <= 0 {
		return false
	}
	return expiresAt <= now.Add(time.Duration(rerollDays)*24*time.Hour).Unix()
}

func resourceServiceUserSessionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	sessionCreateRequest := client_iam.SessionCreateRequest{
		ServiceUser: d.Get("service_user_crn").(string),
		Comment:     d.Get("comment").(string),
	}

	resp, err := (*client_iam.Client)(c).CreateSession(&sessionCreateRequest)
	if err != nil {
		return diag.Errorf("Error creating service user session: %s", err)
	}

	d.SetId(resp.SessionID)
	d.Set("jwt", resp.Jwt)
	d.Set("session_id", resp.SessionID)
	d.Set("expires_at", int(resp.ExpiresAt))

	return resourceServiceUserSessionRead(ctx, d, m)
}

func resourceServiceUserSessionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	serviceUser := d.Get("service_user_crn").(string)
	sessionID := d.Id()

	if serviceUser == "" || sessionID == "" {
		d.SetId("")
		return diags
	}

	session, err := (*client_iam.Client)(c).GetSession(serviceUser, sessionID)
	if err != nil {
		var httpErr *client.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.Errorf("Error reading service user session: %s", err)
	}

	if !session.Enabled {
		d.SetId("")
		return diags
	}

	if ttlUnix, err := strconv.ParseInt(session.TTL, 10, 64); err == nil && ttlUnix > 0 {
		d.Set("expires_at", int(ttlUnix))
	}

	d.Set("service_user_crn", session.Principal)
	d.Set("session_id", session.SessionID)
	d.Set("comment", session.Comment)
	d.Set("enabled", session.Enabled)
	d.Set("ttl", session.TTL)

	return diags
}

func resourceServiceUserSessionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	serviceUser := d.Get("service_user_crn").(string)
	sessionID := d.Id()

	sessionUpdateRequest := client_iam.SessionUpdateRequest{
		Comment: d.Get("comment").(string),
	}

	err = (*client_iam.Client)(c).UpdateSession(serviceUser, sessionID, &sessionUpdateRequest)
	if err != nil {
		return diag.Errorf("Error updating service user session: %s", err)
	}

	return resourceServiceUserSessionRead(ctx, d, m)
}

func resourceServiceUserSessionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c, err := m.(*client.CcpClient).GetService("iam")
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics

	serviceUser := d.Get("service_user_crn").(string)
	sessionID := d.Id()

	err = (*client_iam.Client)(c).DeleteSession(serviceUser, sessionID)
	if err != nil {
		var httpErr *client.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			d.SetId("")
			return diags
		}
		return diag.Errorf("Error deleting service user session: %s", err)
	}

	d.SetId("")

	return diags
}
