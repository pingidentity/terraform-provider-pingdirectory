package config

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &rootDnResource{}
	_ resource.ResourceWithConfigure   = &rootDnResource{}
	_ resource.ResourceWithImportState = &rootDnResource{}
)

// Create a Root DN resource
func NewRootDnResource() resource.Resource {
	return &rootDnResource{}
}

// rootDnResource is the resource implementation.
type rootDnResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// rootDnResourceModel maps the resource schema data.
type rootDnResourceModel struct {
	// Id field required for acceptance testing framework
	DefaultRootPrivilegeName types.Set    `tfsdk:"default_root_privilege_name"`
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
}

// all are optional and computed
// DefaultRootPrivilegeName will be the patch / update target

// Metadata returns the resource type name.
func (r *rootDnResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_root_dn"
}

// GetSchema defines the schema for the resource.
func (r *rootDnResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages the Root DN of PingDirectory.",
		// All are considered computed, since we are importing the existing Root DN
		// from a server, rather than "creating" the Root DN object
		// like a typical Terraform resource.
		Attributes: map[string]schema.Attribute{
			"default_root_privilege_name": schema.SetAttribute{
				Description: "Specifies the names of the privileges that root users will be granted by default.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Configure adds the provider configured client to the resource.
func (r *rootDnResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Create a new resource
// For Root DN, create doesn't actually "create" anything - it "adopts" the existing
// server Root DN settings into being managed by Terraform. This method reads the existing Root DN config
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *rootDnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan rootDnResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.RootDnApi.GetRootDn(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Root DN", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read existing Root DN
	var state rootDnResourceModel
	readRootDnResponse(ctx, getResp, &state)

	// Determine what changes need to be made to match the plan
	updateRDnRequest := r.apiClient.RootDnApi.UpdateRootDn(ProviderBasicAuthContext(ctx, r.providerConfig))

	ops := createRootDnOperations(plan, state)

	if len(ops) > 0 {
		updateRDnRequest = updateRDnRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)
		globalResp, httpResp, err := r.apiClient.RootDnApi.UpdateRootDnExecute(updateRDnRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Root DN", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := globalResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readRootDnResponse(ctx, globalResp, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		// Just put the initial read into the plan
		plan = state
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *rootDnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state rootDnResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.RootDnApi.GetRootDn(ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Root DN", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRootDnResponse(ctx, getResp, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a RootDNResponse object into the model struct
func readRootDnResponse(ctx context.Context, r *client.RootDnResponse, state *rootDnResourceModel) {
	// Placeholder Id value for acceptance test framework
	// Id not in RootDN, can be set to anything
	state.Id = types.StringValue("id")
	state.DefaultRootPrivilegeName = internaltypes.GetEnumSet(r.DefaultRootPrivilegeName)

	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20)
}

// Create any update operations necessary to make the state match the plan
func createRootDnOperations(plan rootDnResourceModel, state rootDnResourceModel) []client.Operation {
	var ops []client.Operation

	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultRootPrivilegeName, state.DefaultRootPrivilegeName, "default-root-privilege-name")

	return ops
}

// Update the Root DN Permissions - similar to the Create method since the permissions are adopted
func (r *rootDnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan rootDnResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state
	var state rootDnResourceModel
	req.State.Get(ctx, &state)
	updateRDnRequest := r.apiClient.RootDnApi.UpdateRootDn(ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createRootDnOperations(plan, state)
	if len(ops) > 0 {
		updateRDnRequest = updateRDnRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		globalResp, httpResp, err := r.apiClient.RootDnApi.UpdateRootDnExecute(updateRDnRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Root DN", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := globalResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readRootDnResponse(ctx, globalResp, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No Root DN operations created for update")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// Terraform can't actually delete the Root DN, so this method does nothing.
// Terraform will just "forget" about the Root DN resource, and it can be managed elsewhere.
func (r *rootDnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *rootDnResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set an arbitrary state value to appease terraform - the placeholder will immediately be
	// replaced with the actual instance name when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
