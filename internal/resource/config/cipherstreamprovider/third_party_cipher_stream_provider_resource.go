package cipherstreamprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &thirdPartyCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &thirdPartyCipherStreamProviderResource{}
	_ resource.Resource                = &defaultThirdPartyCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyCipherStreamProviderResource{}
)

// Create a Third Party Cipher Stream Provider resource
func NewThirdPartyCipherStreamProviderResource() resource.Resource {
	return &thirdPartyCipherStreamProviderResource{}
}

func NewDefaultThirdPartyCipherStreamProviderResource() resource.Resource {
	return &defaultThirdPartyCipherStreamProviderResource{}
}

// thirdPartyCipherStreamProviderResource is the resource implementation.
type thirdPartyCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyCipherStreamProviderResource is the resource implementation.
type defaultThirdPartyCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_cipher_stream_provider"
}

func (r *defaultThirdPartyCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyCipherStreamProviderResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyCipherStreamProviderSchema(ctx, req, resp, true)
}

func thirdPartyCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Cipher Stream Provider.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Cipher Stream Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Cipher Stream Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.",
				Required:    true,
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
func addOptionalThirdPartyCipherStreamProviderFields(ctx context.Context, addRequest *client.AddThirdPartyCipherStreamProviderRequest, plan thirdPartyCipherStreamProviderResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ThirdPartyCipherStreamProviderResponse object into the model struct
func readThirdPartyCipherStreamProviderResponse(ctx context.Context, r *client.ThirdPartyCipherStreamProviderResponse, state *thirdPartyCipherStreamProviderResourceModel, expectedValues *thirdPartyCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyCipherStreamProviderOperations(plan thirdPartyCipherStreamProviderResourceModel, state thirdPartyCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *thirdPartyCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyCipherStreamProviderSchemaUrn{client.ENUMTHIRDPARTYCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddThirdPartyCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyCipherStreamProviderResourceModel
	readThirdPartyCipherStreamProviderResponse(ctx, addResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultThirdPartyCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyCipherStreamProviderResourceModel
	readThirdPartyCipherStreamProviderResponse(ctx, readResponse.ThirdPartyCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyCipherStreamProviderResponse(ctx, updateResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyCipherStreamProviderResponse(ctx, readResponse.ThirdPartyCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyCipherStreamProviderResponse(ctx, updateResponse.ThirdPartyCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *thirdPartyCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyCipherStreamProvider(ctx, req, resp)
}

func (r *defaultThirdPartyCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyCipherStreamProvider(ctx, req, resp)
}

func importThirdPartyCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
