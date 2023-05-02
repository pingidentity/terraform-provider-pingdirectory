package azureauthenticationmethod

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &defaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultAzureAuthenticationMethodResource{}
	_ resource.Resource                = &defaultDefaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultDefaultAzureAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultDefaultAzureAuthenticationMethodResource{}
)

// Create a Default Azure Authentication Method resource
func NewDefaultAzureAuthenticationMethodResource() resource.Resource {
	return &defaultAzureAuthenticationMethodResource{}
}

func NewDefaultDefaultAzureAuthenticationMethodResource() resource.Resource {
	return &defaultDefaultAzureAuthenticationMethodResource{}
}

// defaultAzureAuthenticationMethodResource is the resource implementation.
type defaultAzureAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDefaultAzureAuthenticationMethodResource is the resource implementation.
type defaultDefaultAzureAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *defaultAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_azure_authentication_method"
}

func (r *defaultDefaultAzureAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_default_azure_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *defaultAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDefaultAzureAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type defaultAzureAuthenticationMethodResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	TenantID        types.String `tfsdk:"tenant_id"`
	ClientID        types.String `tfsdk:"client_id"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *defaultAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	defaultAzureAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultDefaultAzureAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	defaultAzureAuthenticationMethodSchema(ctx, req, resp, true)
}

func defaultAzureAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Default Azure Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Description: "The tenant ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_TENANT_ID environment variable.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "The client ID to use to authenticate. If this is not provided, then it will be obtained from the AZURE_CLIENT_ID",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Azure Authentication Method",
				Optional:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalDefaultAzureAuthenticationMethodFields(ctx context.Context, addRequest *client.AddDefaultAzureAuthenticationMethodRequest, plan defaultAzureAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TenantID) {
		addRequest.TenantID = plan.TenantID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientID) {
		addRequest.ClientID = plan.ClientID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DefaultAzureAuthenticationMethodResponse object into the model struct
func readDefaultAzureAuthenticationMethodResponse(ctx context.Context, r *client.DefaultAzureAuthenticationMethodResponse, state *defaultAzureAuthenticationMethodResourceModel, expectedValues *defaultAzureAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.TenantID = internaltypes.StringTypeOrNil(r.TenantID, internaltypes.IsEmptyString(expectedValues.TenantID))
	state.ClientID = internaltypes.StringTypeOrNil(r.ClientID, internaltypes.IsEmptyString(expectedValues.ClientID))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDefaultAzureAuthenticationMethodOperations(plan defaultAzureAuthenticationMethodResourceModel, state defaultAzureAuthenticationMethodResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.TenantID, state.TenantID, "tenant-id")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *defaultAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultAzureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDefaultAzureAuthenticationMethodRequest(plan.Id.ValueString(),
		[]client.EnumdefaultAzureAuthenticationMethodSchemaUrn{client.ENUMDEFAULTAZUREAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0AZURE_AUTHENTICATION_METHODDEFAULT})
	addOptionalDefaultAzureAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAzureAuthenticationMethodRequest(
		client.AddDefaultAzureAuthenticationMethodRequestAsAddAzureAuthenticationMethodRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.AddAzureAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Default Azure Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state defaultAzureAuthenticationMethodResourceModel
	readDefaultAzureAuthenticationMethodResponse(ctx, addResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultDefaultAzureAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultAzureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Default Azure Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultAzureAuthenticationMethodResourceModel
	readDefaultAzureAuthenticationMethodResponse(ctx, readResponse.DefaultAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDefaultAzureAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Default Azure Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDefaultAzureAuthenticationMethodResponse(ctx, updateResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *defaultAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDefaultAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDefaultAzureAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDefaultAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDefaultAzureAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state defaultAzureAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.GetAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Default Azure Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDefaultAzureAuthenticationMethodResponse(ctx, readResponse.DefaultAzureAuthenticationMethodResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *defaultAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDefaultAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDefaultAzureAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDefaultAzureAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDefaultAzureAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan defaultAzureAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultAzureAuthenticationMethodResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDefaultAzureAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AzureAuthenticationMethodApi.UpdateAzureAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Default Azure Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDefaultAzureAuthenticationMethodResponse(ctx, updateResponse.DefaultAzureAuthenticationMethodResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultDefaultAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *defaultAzureAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state defaultAzureAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethodExecute(r.apiClient.AzureAuthenticationMethodApi.DeleteAzureAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Default Azure Authentication Method", err, httpResp)
		return
	}
}

func (r *defaultAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDefaultAzureAuthenticationMethod(ctx, req, resp)
}

func (r *defaultDefaultAzureAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDefaultAzureAuthenticationMethod(ctx, req, resp)
}

func importDefaultAzureAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
