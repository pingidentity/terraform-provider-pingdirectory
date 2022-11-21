package ldap

import (
	"context"
	"strings"
	"terraform-provider-pingdirectory/internal/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/go-ldap/ldap/v3"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &usersResource{}
	_ resource.ResourceWithConfigure   = &usersResource{}
	_ resource.ResourceWithImportState = &usersResource{}
)

const baseDN = "ou=people,dc=example,dc=com"

// Create a Users resource
func NewUsersResource() resource.Resource {
	return &usersResource{}
}

// usersResource is the resource implementation.
type usersResource struct {
	providerConfig utils.ProviderConfiguration
}

// usersResourceModel maps the resource schema data.
type usersResourceModel struct {
	Uid         types.String `tfsdk:"uid"`
	Sn          types.String `tfsdk:"sn"`
	Cn          types.String `tfsdk:"cn"`
	GivenName   types.String `tfsdk:"given_name"`
	Mail        types.String `tfsdk:"mail"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// Helper function to get the DN for a user with a given ID
func userDn(uid string) string {
	return "uid=" + uid + "," + baseDN
}

// Get the common name for a user
func commonName(sn string, givenName string) string {
	return givenName + " " + sn
}

// Helper method to create an authenticated connection to the directory server
func bind(url string, username string, password string) (*ldap.Conn, error) {
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
			"sn": {
				Description: "Surname of the user.",
				Type:        types.StringType,
				Required:    true,
			},
			"cn": {
				Description: "Common name of the user.",
				Type:        types.StringType,
				Computed:    true,
			},
			"given_name": {
				Description: "Given name of the user.",
				Type:        types.StringType,
				Required:    true,
			},
			"mail": {
				Description: "Email address of the user.",
				Type:        types.StringType,
				Required:    true,
			},
			"description": {
				Description: "Description of the user.",
				Type:        types.StringType,
				Optional:    true,
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

	resourceConfig := req.ProviderData.(utils.ResourceConfiguration)
	r.providerConfig = resourceConfig.ProviderConfig
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

	l, err := bind(r.providerConfig.LdapHost, r.providerConfig.Username, r.providerConfig.Password)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	// NOTE: this does no input sanitization so it's probably HIGHLY insecure with the string concatenation and directly using input
	addRequest := ldap.NewAddRequest(userDn(plan.Uid.ValueString()), nil)
	addRequest.Attribute("objectClass", []string{"person", "organizationalPerson", "inetOrgPerson"})
	addRequest.Attribute("sn", []string{plan.Sn.ValueString()})
	addRequest.Attribute("cn", []string{commonName(plan.Sn.ValueString(), plan.GivenName.ValueString())})
	addRequest.Attribute("givenName", []string{plan.GivenName.ValueString()})
	addRequest.Attribute("uid", []string{plan.Uid.ValueString()})
	addRequest.Attribute("mail", []string{plan.Mail.ValueString()})
	// Just set a default password
	// I'm not sure if a provider can manage password like this - because the value saved in the state
	// will be an encrypted version of the password, and the directory server doesn't allow changing a password
	// to the same value as the current value. It's probably just an API that doesn't make sense to manage with Terraform.
	addRequest.Attribute("userPassword", []string{r.providerConfig.DefaultUserPassword})
	if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
		addRequest.Attribute("description", []string{plan.Description.ValueString()})
	}

	err = l.Add(addRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while adding the user", err.Error())
		return
	}

	// Populate Computed attribute values
	plan.Cn = types.StringValue(commonName(plan.Sn.ValueString(), plan.GivenName.ValueString()))
	plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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

	l, err := bind(r.providerConfig.LdapHost, r.providerConfig.Username, r.providerConfig.Password)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	//TODO appropriate parameters here - just following one of the package's examples
	// NOTE: this does no input sanitization so it's probably HIGHLY insecure
	searchRequest := ldap.NewSearchRequest(baseDN, ldap.ScopeSingleLevel, ldap.NeverDerefAliases, 0, 0, false,
		"(&(uid="+state.Uid.ValueString()+")(objectClass=person))", []string{"*"}, nil)
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
	state.Sn = types.StringValue(searchResult.Entries[0].GetAttributeValue("sn"))
	state.Cn = types.StringValue(searchResult.Entries[0].GetAttributeValue("cn"))
	state.GivenName = types.StringValue(searchResult.Entries[0].GetAttributeValue("givenName"))
	state.Mail = types.StringValue(searchResult.Entries[0].GetAttributeValue("mail"))
	description := searchResult.Entries[0].GetAttributeValue("description")
	if description != "" {
		state.Description = types.StringValue(searchResult.Entries[0].GetAttributeValue("description"))
	} else {
		state.Description = types.StringNull()
	}

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

	l, err := bind(r.providerConfig.LdapHost, r.providerConfig.Username, r.providerConfig.Password)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	// Need to get the current state of the entry
	//TODO use the state or search the directory? Could lead to issues if the directory state is updated and we form the modify request incorrectly.
	// Kind of a provider best practice question.
	//TODO use --permissiveModify request control?
	var state usersResourceModel
	req.State.Get(ctx, &state)

	// The UID can't change here because it's set to RequiresReplace above
	replaceCn := false
	modifyRequest := ldap.NewModifyRequest(userDn(plan.Uid.ValueString()), nil)
	modifyRequest.Controls = append(modifyRequest.Controls)
	if !state.Sn.Equal(plan.Sn) {
		tflog.Info(ctx, "Replacing sn: ("+state.Sn.ValueString()+" -> "+plan.Sn.ValueString()+")")
		modifyRequest.Replace("sn", []string{plan.Sn.ValueString()})
		replaceCn = true
	}
	if !state.GivenName.Equal(plan.GivenName) {
		tflog.Info(ctx, "Replacing givenName: ("+state.GivenName.ValueString()+" -> "+plan.GivenName.ValueString()+")")
		modifyRequest.Replace("givenName", []string{plan.GivenName.ValueString()})
		replaceCn = true
	}
	if replaceCn {
		tflog.Info(ctx, "Replacing cn")
		modifyRequest.Replace("cn", []string{commonName(plan.Sn.ValueString(), plan.GivenName.ValueString())})
	}
	if !state.Mail.Equal(plan.Mail) {
		tflog.Info(ctx, "Replacing mail: ("+state.Mail.ValueString()+" -> "+plan.Mail.ValueString()+")")
		modifyRequest.Replace("mail", []string{plan.Description.ValueString()})
	}
	if !plan.Description.Equal(state.Description) {
		if !plan.Description.IsUnknown() && !plan.Description.IsNull() {
			tflog.Info(ctx, "Replacing description: ("+state.Description.ValueString()+" -> "+plan.Description.ValueString()+")")
			modifyRequest.Replace("description", []string{plan.Description.ValueString()})
		} else {
			tflog.Info(ctx, "Deleting description")
			modifyRequest.Replace("description", []string{})
		}
	}

	err = l.Modify(modifyRequest)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while updating the entry", err.Error())
		return
	}

	// Update resource state with updated items and timestamp
	plan.Cn = types.StringValue(commonName(plan.Sn.ValueString(), plan.GivenName.ValueString()))
	plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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

	l, err := bind(r.providerConfig.LdapHost, r.providerConfig.Username, r.providerConfig.Password)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while authenticating to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	deleteRequest := ldap.NewDelRequest(userDn(state.Uid.ValueString()), nil)

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
