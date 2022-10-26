package pingdirectory

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/go-ldap/ldap/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &usersResource{}
	_ resource.ResourceWithConfigure   = &usersResource{}
	_ resource.ResourceWithImportState = &usersResource{}
)

const baseDN = "ou=people,dc=example,dc=com"

// NewUsersResource is a helper function to simplify the provider implementation.
func NewUsersResource() resource.Resource {
	return &usersResource{}
}

// usersResource is the resource implementation.
type usersResource struct {
	ldapConnectionInfo pingdirectoryProviderModel
}

// usersResourceModel maps the resource schema data.
type usersResourceModel struct {
	Uid         types.String `tfsdk:"uid"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Helper function to get the DN for a user with a given ID
func UserDn(uid string) string {
	return "uid=" + uid + "," + baseDN
}

// Helper method to create an authenticated connection to the directory server
func Bind(url string, username string, password string) (*ldap.Conn, error) {
	l, err := ldap.DialURL(url)
	if err != nil {
		return nil, err
	}

	err = l.Bind(username, password)
	if err != nil {
		l.Close()
		return nil, err
	}

	return l, nil
}

// Metadata returns the resource type name.
func (r *usersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// GetSchema defines the schema for the resource.
func (r *usersResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a user.",
		Attributes: map[string]tfsdk.Attribute{
			"uid": {
				Description: "User ID of the user.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"description": {
				Description: "Description of the user.",
				Type:        types.StringType,
				Required:    true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the user.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *usersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.ldapConnectionInfo = req.ProviderData.(pingdirectoryProviderModel)
}

// Create a new resource
func (r *usersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan usersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	l, err := Bind(r.ldapConnectionInfo.Host.Value, r.ldapConnectionInfo.Username.Value, r.ldapConnectionInfo.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	// NOTE: this does no input sanitization so it's probably HIGHLY insecure
	//TODO let these other attribute values come from the resource definition
	addRequest := ldap.NewAddRequest(UserDn(plan.Uid.Value), nil)
	addRequest.Attribute("description", []string{plan.Description.Value})
	addRequest.Attribute("objectClass", []string{"person", "organizationalPerson", "inetOrgPerson"})
	addRequest.Attribute("sn", []string{"Mahomes"})
	addRequest.Attribute("cn", []string{"Patrick Mahomes"})
	addRequest.Attribute("givenName", []string{"Patrick"})
	addRequest.Attribute("uid", []string{plan.Uid.Value})
	addRequest.Attribute("mail", []string{plan.Uid.Value + "@example.com"})
	addRequest.Attribute("userPassword", []string{"2FederateM0re"})

	err = l.Add(addRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while adding the user", err.Error())
		return
	}

	// Populate Computed attribute values
	plan.LastUpdated = types.String{Value: string(time.Now().Format(time.RFC850))}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *usersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	//TODO how to gracefully handle if the user on the remote server gets deleted (outside of terraform's control)?
	// Kind of a general question for best practice for a provider.

	// Get current state
	var state usersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	l, err := Bind(r.ldapConnectionInfo.Host.Value, r.ldapConnectionInfo.Username.Value, r.ldapConnectionInfo.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	//TODO appropriate parameters here - just following one of the package's examples
	// NOTE: this does no input sanitization so it's probably HIGHLY insecure
	//TODO add more attributes to the search if more are supported in the schema
	searchRequest := ldap.NewSearchRequest(baseDN, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(&(uid="+state.Uid.Value+")(objectClass=person))", []string{"dn", "description"}, nil)
	searchResult, err := l.Search(searchRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while searching the PingDirectory server", err.Error())
		return
	}

	if len(searchResult.Entries) == 0 {
		resp.Diagnostics.AddError("No entries found matching the search in the PingDirectory Server", "")
		return
	}
	if len(searchResult.Entries) > 1 {
		var sb strings.Builder
		for _, entry := range searchResult.Entries {
			sb.WriteString(entry.DN + " ")
		}
		resp.Diagnostics.AddError("More than one entry found matching the search in the PingDirectory Server", "Entries found: "+sb.String())
		return
	}

	// Update state from the search
	//TODO if more attributes are supported add here
	state.Description = types.String{Value: searchResult.Entries[0].GetAttributeValue("description")}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *usersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan usersResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	l, err := Bind(r.ldapConnectionInfo.Host.Value, r.ldapConnectionInfo.Username.Value, r.ldapConnectionInfo.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	// The UID can't change here because it's set to RequiresReplace above
	modifyRequest := ldap.NewModifyRequest(UserDn(plan.Uid.Value), nil)
	modifyRequest.Replace("description", []string{plan.Description.Value})

	err = l.Modify(modifyRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while updating the entry", err.Error())
		return
	}

	// Update resource state with updated items and timestamp
	//TODO update any future computed values here
	plan.LastUpdated = types.String{Value: string(time.Now().Format(time.RFC850))}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *usersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state usersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	l, err := Bind(r.ldapConnectionInfo.Host.Value, r.ldapConnectionInfo.Username.Value, r.ldapConnectionInfo.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	deleteRequest := ldap.NewDelRequest(UserDn(state.Uid.Value), nil)

	err = l.Del(deleteRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while deleting the user in the PingDirectory server", err.Error())
		return
	}
}

func (r *usersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Uid attribute
	resource.ImportStatePassthroughID(ctx, path.Root("uid"), req, resp)
}
