package plugin

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &internalSearchRatePluginResource{}
	_ resource.ResourceWithConfigure   = &internalSearchRatePluginResource{}
	_ resource.ResourceWithImportState = &internalSearchRatePluginResource{}
	_ resource.Resource                = &defaultInternalSearchRatePluginResource{}
	_ resource.ResourceWithConfigure   = &defaultInternalSearchRatePluginResource{}
	_ resource.ResourceWithImportState = &defaultInternalSearchRatePluginResource{}
)

// Create a Internal Search Rate Plugin resource
func NewInternalSearchRatePluginResource() resource.Resource {
	return &internalSearchRatePluginResource{}
}

func NewDefaultInternalSearchRatePluginResource() resource.Resource {
	return &defaultInternalSearchRatePluginResource{}
}

// internalSearchRatePluginResource is the resource implementation.
type internalSearchRatePluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultInternalSearchRatePluginResource is the resource implementation.
type defaultInternalSearchRatePluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *internalSearchRatePluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_internal_search_rate_plugin"
}

func (r *defaultInternalSearchRatePluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_internal_search_rate_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *internalSearchRatePluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultInternalSearchRatePluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type internalSearchRatePluginResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	PluginType                  types.Set    `tfsdk:"plugin_type"`
	NumThreads                  types.Int64  `tfsdk:"num_threads"`
	BaseDN                      types.String `tfsdk:"base_dn"`
	LowerBound                  types.Int64  `tfsdk:"lower_bound"`
	UpperBound                  types.Int64  `tfsdk:"upper_bound"`
	FilterPrefix                types.String `tfsdk:"filter_prefix"`
	FilterSuffix                types.String `tfsdk:"filter_suffix"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *internalSearchRatePluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	internalSearchRatePluginSchema(ctx, req, resp, false)
}

func (r *defaultInternalSearchRatePluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	internalSearchRatePluginSchema(ctx, req, resp, true)
}

func internalSearchRatePluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Internal Search Rate Plugin.",
		Attributes: map[string]schema.Attribute{
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"num_threads": schema.Int64Attribute{
				Description: "Specifies the number of concurrent threads that should be used to process the search operations.",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "Specifies the base DN to use for the searches to perform.",
				Required:    true,
			},
			"lower_bound": schema.Int64Attribute{
				Description: "Specifies the lower bound for the numeric value which will be inserted into the search filter.",
				Optional:    true,
				Computed:    true,
			},
			"upper_bound": schema.Int64Attribute{
				Description: "Specifies the upper bound for the numeric value which will be inserted into the search filter.",
				Optional:    true,
			},
			"filter_prefix": schema.StringAttribute{
				Description: "Specifies a prefix which will be used in front of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should contain the entire filter string.",
				Required:    true,
			},
			"filter_suffix": schema.StringAttribute{
				Description: "Specifies a suffix which will be used after of the randomly-selected numeric value in all search filters used. If no upper bound is defined, then this should be omitted.",
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
			"invoke_for_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether the plug-in should be invoked for internal operations.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalInternalSearchRatePluginFields(ctx context.Context, addRequest *client.AddInternalSearchRatePluginRequest, plan internalSearchRatePluginResourceModel) error {
	if internaltypes.IsDefined(plan.PluginType) {
		var slice []string
		plan.PluginType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpluginPluginTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpluginPluginTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PluginType = enumSlice
	}
	if internaltypes.IsDefined(plan.NumThreads) {
		intVal := int32(plan.NumThreads.ValueInt64())
		addRequest.NumThreads = &intVal
	}
	if internaltypes.IsDefined(plan.LowerBound) {
		intVal := int32(plan.LowerBound.ValueInt64())
		addRequest.LowerBound = &intVal
	}
	if internaltypes.IsDefined(plan.UpperBound) {
		intVal := int32(plan.UpperBound.ValueInt64())
		addRequest.UpperBound = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FilterSuffix) {
		stringVal := plan.FilterSuffix.ValueString()
		addRequest.FilterSuffix = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.InvokeForInternalOperations) {
		boolVal := plan.InvokeForInternalOperations.ValueBool()
		addRequest.InvokeForInternalOperations = &boolVal
	}
	return nil
}

// Read a InternalSearchRatePluginResponse object into the model struct
func readInternalSearchRatePluginResponse(ctx context.Context, r *client.InternalSearchRatePluginResponse, state *internalSearchRatePluginResourceModel, expectedValues *internalSearchRatePluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.NumThreads = types.Int64Value(int64(r.NumThreads))
	state.BaseDN = types.StringValue(r.BaseDN)
	state.LowerBound = internaltypes.Int64TypeOrNil(r.LowerBound)
	state.UpperBound = internaltypes.Int64TypeOrNil(r.UpperBound)
	state.FilterPrefix = types.StringValue(r.FilterPrefix)
	state.FilterSuffix = internaltypes.StringTypeOrNil(r.FilterSuffix, internaltypes.IsEmptyString(expectedValues.FilterSuffix))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createInternalSearchRatePluginOperations(plan internalSearchRatePluginResourceModel, state internalSearchRatePluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumThreads, state.NumThreads, "num-threads")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowerBound, state.LowerBound, "lower-bound")
	operations.AddInt64OperationIfNecessary(&ops, plan.UpperBound, state.UpperBound, "upper-bound")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterPrefix, state.FilterPrefix, "filter-prefix")
	operations.AddStringOperationIfNecessary(&ops, plan.FilterSuffix, state.FilterSuffix, "filter-suffix")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *internalSearchRatePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan internalSearchRatePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddInternalSearchRatePluginRequest(plan.Id.ValueString(),
		[]client.EnuminternalSearchRatePluginSchemaUrn{client.ENUMINTERNALSEARCHRATEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGININTERNAL_SEARCH_RATE},
		plan.BaseDN.ValueString(),
		plan.FilterPrefix.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalInternalSearchRatePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Internal Search Rate Plugin", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddInternalSearchRatePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Internal Search Rate Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state internalSearchRatePluginResourceModel
	readInternalSearchRatePluginResponse(ctx, addResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultInternalSearchRatePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan internalSearchRatePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Internal Search Rate Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state internalSearchRatePluginResourceModel
	readInternalSearchRatePluginResponse(ctx, readResponse.InternalSearchRatePluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createInternalSearchRatePluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Internal Search Rate Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readInternalSearchRatePluginResponse(ctx, updateResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *internalSearchRatePluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readInternalSearchRatePlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultInternalSearchRatePluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readInternalSearchRatePlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readInternalSearchRatePlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state internalSearchRatePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Internal Search Rate Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readInternalSearchRatePluginResponse(ctx, readResponse.InternalSearchRatePluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *internalSearchRatePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateInternalSearchRatePlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultInternalSearchRatePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateInternalSearchRatePlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateInternalSearchRatePlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan internalSearchRatePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state internalSearchRatePluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createInternalSearchRatePluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Internal Search Rate Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readInternalSearchRatePluginResponse(ctx, updateResponse.InternalSearchRatePluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultInternalSearchRatePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *internalSearchRatePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state internalSearchRatePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Internal Search Rate Plugin", err, httpResp)
		return
	}
}

func (r *internalSearchRatePluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importInternalSearchRatePlugin(ctx, req, resp)
}

func (r *defaultInternalSearchRatePluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importInternalSearchRatePlugin(ctx, req, resp)
}

func importInternalSearchRatePlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
