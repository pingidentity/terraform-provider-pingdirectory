package pingdirectory

import (
	"context"
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

// NewUsersResource is a helper function to simplify the provider implementation.
func NewUsersResource() resource.Resource {
	return &usersResource{}
}

// usersResource is the resource implementation.
type usersResource struct {
	//TODO find appropriate client type
	client pingdirectoryProviderModel
}

// usersResourceModel maps the resource schema data.
type usersResourceModel struct {
	Uid         types.String `tfsdk:"uid"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
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

	r.client = req.ProviderData.(pingdirectoryProviderModel)
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

	l, err := ldap.DialURL(r.client.Host.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while connecting to the PingDirectory server", err.Error())
		return
	}
	defer l.Close()

	err = l.Bind(r.client.Username.Value, r.client.Password.Value)
	if err != nil {
		resp.Diagnostics.AddError("An error occurred while binding to the PingDirectory server", err.Error())
		return
	}

	// NOTE: this does no input sanitization so it's probably HIGHLY insecure
	//TODO let these other attribute values come from the resource definition
	addRequest := ldap.NewAddRequest("uid="+plan.Uid.Value+",ou=people,dc=example,dc=com", nil)
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
	// Get current state
	var state usersResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//TODO read current user from PD
	// state.Uid
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading PingDirectory user",
	// 		"Could not read PingDirectory user with ID "+state.Uid.Value+": "+err.Error(),
	// 	)
	// 	return
	// }

	// Overwrite items with refreshed state
	// state.SomeComputedValue = types.String{Value: "computed"}

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

	// TODO update existing user entry
	// plan.Uid
	// plan.Description
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error Updating PingDirectory user",
	//		"Could not update user, unexpected error: "+err.Error(),
	//	)
	//	return
	//}

	// TODO fetch updated user (if necessary?)

	// Update resource state with updated items and timestamp
	// plan.SomeComputedValue = types.String{Value: "computed"}
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

	// TODO delete entry from PD
	// state.Uid
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error Deleting PingDirectory user",
	//		"Could not delete user, unexpected error: "+err.Error(),
	//	)
	//	return
	//}
}

func (r *usersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Uid attribute
	//TODO verify this works
	resource.ImportStatePassthroughID(ctx, path.Root("uid"), req, resp)
}
