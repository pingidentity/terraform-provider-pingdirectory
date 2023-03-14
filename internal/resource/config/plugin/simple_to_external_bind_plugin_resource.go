package plugin

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
	_ resource.Resource                = &simpleToExternalBindPluginResource{}
	_ resource.ResourceWithConfigure   = &simpleToExternalBindPluginResource{}
	_ resource.ResourceWithImportState = &simpleToExternalBindPluginResource{}
	_ resource.Resource                = &defaultSimpleToExternalBindPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultSimpleToExternalBindPluginResource{}
	_ resource.ResourceWithImportState = &defaultSimpleToExternalBindPluginResource{}
)

// Create a Simple To External Bind Plugin resource
func NewSimpleToExternalBindPluginResource() resource.Resource {
	return &simpleToExternalBindPluginResource{}
}

func NewDefaultSimpleToExternalBindPluginResource() resource.Resource {
	return &defaultSimpleToExternalBindPluginResource{}
}

// simpleToExternalBindPluginResource is the resource implementation.
type simpleToExternalBindPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSimpleToExternalBindPluginResource is the resource implementation.
type defaultSimpleToExternalBindPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *simpleToExternalBindPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_simple_to_external_bind_plugin"
}

func (r *defaultSimpleToExternalBindPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_simple_to_external_bind_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *simpleToExternalBindPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSimpleToExternalBindPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type simpleToExternalBindPluginResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	ConnectionCriteria types.String `tfsdk:"connection_criteria"`
	RequestCriteria    types.String `tfsdk:"request_criteria"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *simpleToExternalBindPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleToExternalBindPluginSchema(ctx, req, resp, false)
}

func (r *defaultSimpleToExternalBindPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	simpleToExternalBindPluginSchema(ctx, req, resp, true)
}

func simpleToExternalBindPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Simple To External Bind Plugin.",
		Attributes: map[string]schema.Attribute{
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a connection criteria object that may be used to indicate the set of clients for which this plugin should be used. If a value is provided, then this plugin will only be used for requests from client connections matching this criteria.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a request criteria object that may be used to indicate the set of requests for which this plugin should be used. If a value is provided, then this plugin will only be used for bind requests matching this criteria.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Plugin",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the plug-in is enabled for use.",
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
func addOptionalSimpleToExternalBindPluginFields(ctx context.Context, addRequest *client.AddSimpleToExternalBindPluginRequest, plan simpleToExternalBindPluginResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		stringVal := plan.ConnectionCriteria.ValueString()
		addRequest.ConnectionCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		stringVal := plan.RequestCriteria.ValueString()
		addRequest.RequestCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a SimpleToExternalBindPluginResponse object into the model struct
func readSimpleToExternalBindPluginResponse(ctx context.Context, r *client.SimpleToExternalBindPluginResponse, state *simpleToExternalBindPluginResourceModel, expectedValues *simpleToExternalBindPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSimpleToExternalBindPluginOperations(plan simpleToExternalBindPluginResourceModel, state simpleToExternalBindPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *simpleToExternalBindPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleToExternalBindPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSimpleToExternalBindPluginRequest(plan.Id.ValueString(),
		[]client.EnumsimpleToExternalBindPluginSchemaUrn{client.ENUMSIMPLETOEXTERNALBINDPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSIMPLE_TO_EXTERNAL_BIND},
		plan.Enabled.ValueBool())
	addOptionalSimpleToExternalBindPluginFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSimpleToExternalBindPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Simple To External Bind Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state simpleToExternalBindPluginResourceModel
	readSimpleToExternalBindPluginResponse(ctx, addResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSimpleToExternalBindPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan simpleToExternalBindPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple To External Bind Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state simpleToExternalBindPluginResourceModel
	readSimpleToExternalBindPluginResponse(ctx, readResponse.SimpleToExternalBindPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSimpleToExternalBindPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple To External Bind Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleToExternalBindPluginResponse(ctx, updateResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *simpleToExternalBindPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleToExternalBindPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleToExternalBindPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSimpleToExternalBindPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSimpleToExternalBindPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state simpleToExternalBindPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Simple To External Bind Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSimpleToExternalBindPluginResponse(ctx, readResponse.SimpleToExternalBindPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *simpleToExternalBindPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleToExternalBindPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSimpleToExternalBindPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSimpleToExternalBindPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSimpleToExternalBindPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan simpleToExternalBindPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state simpleToExternalBindPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSimpleToExternalBindPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Simple To External Bind Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSimpleToExternalBindPluginResponse(ctx, updateResponse.SimpleToExternalBindPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSimpleToExternalBindPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *simpleToExternalBindPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state simpleToExternalBindPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Simple To External Bind Plugin", err, httpResp)
		return
	}
}

func (r *simpleToExternalBindPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleToExternalBindPlugin(ctx, req, resp)
}

func (r *defaultSimpleToExternalBindPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSimpleToExternalBindPlugin(ctx, req, resp)
}

func importSimpleToExternalBindPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
