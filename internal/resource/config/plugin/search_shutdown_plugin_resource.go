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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &searchShutdownPluginResource{}
	_ resource.ResourceWithConfigure   = &searchShutdownPluginResource{}
	_ resource.ResourceWithImportState = &searchShutdownPluginResource{}
	_ resource.Resource                = &defaultSearchShutdownPluginResource{}
	_ resource.ResourceWithConfigure   = &defaultSearchShutdownPluginResource{}
	_ resource.ResourceWithImportState = &defaultSearchShutdownPluginResource{}
)

// Create a Search Shutdown Plugin resource
func NewSearchShutdownPluginResource() resource.Resource {
	return &searchShutdownPluginResource{}
}

func NewDefaultSearchShutdownPluginResource() resource.Resource {
	return &defaultSearchShutdownPluginResource{}
}

// searchShutdownPluginResource is the resource implementation.
type searchShutdownPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSearchShutdownPluginResource is the resource implementation.
type defaultSearchShutdownPluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *searchShutdownPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_search_shutdown_plugin"
}

func (r *defaultSearchShutdownPluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_search_shutdown_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *searchShutdownPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSearchShutdownPluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type searchShutdownPluginResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	BaseDN                types.String `tfsdk:"base_dn"`
	Scope                 types.String `tfsdk:"scope"`
	Filter                types.String `tfsdk:"filter"`
	IncludeAttribute      types.Set    `tfsdk:"include_attribute"`
	OutputFile            types.String `tfsdk:"output_file"`
	PreviousFileExtension types.String `tfsdk:"previous_file_extension"`
	Description           types.String `tfsdk:"description"`
	Enabled               types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *searchShutdownPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchShutdownPluginSchema(ctx, req, resp, false)
}

func (r *defaultSearchShutdownPluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	searchShutdownPluginSchema(ctx, req, resp, true)
}

func searchShutdownPluginSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Search Shutdown Plugin.",
		Attributes: map[string]schema.Attribute{
			"base_dn": schema.StringAttribute{
				Description: "The base DN to use for the search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.StringAttribute{
				Description: "The scope to use for the search.",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "The filter to use for the search.",
				Required:    true,
			},
			"include_attribute": schema.SetAttribute{
				Description: "The name of an attribute that should be included in the results. This may include any token which is allowed as a requested attribute in search requests, including the name of an attribute, an asterisk (to indicate all user attributes), a plus sign (to indicate all operational attributes), an object class name preceded with an at symbol (to indicate all attributes associated with that object class), an attribute name preceded by a caret (to indicate that attribute should be excluded), or an object class name preceded by a caret and an at symbol (to indicate that all attributes associated with that object class should be excluded).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"output_file": schema.StringAttribute{
				Description: "The path of an LDIF file that should be created with the results of the search.",
				Required:    true,
			},
			"previous_file_extension": schema.StringAttribute{
				Description: "An extension that should be appended to the name of an existing output file rather than deleting it. If a file already exists with the full previous file name, then it will be deleted before the current file is renamed to become the previous file.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalSearchShutdownPluginFields(ctx context.Context, addRequest *client.AddSearchShutdownPluginRequest, plan searchShutdownPluginResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseDN) {
		stringVal := plan.BaseDN.ValueString()
		addRequest.BaseDN = &stringVal
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PreviousFileExtension) {
		stringVal := plan.PreviousFileExtension.ValueString()
		addRequest.PreviousFileExtension = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a SearchShutdownPluginResponse object into the model struct
func readSearchShutdownPluginResponse(ctx context.Context, r *client.SearchShutdownPluginResponse, state *searchShutdownPluginResourceModel, expectedValues *searchShutdownPluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.StringTypeOrNil(r.BaseDN, internaltypes.IsEmptyString(expectedValues.BaseDN))
	state.Scope = types.StringValue(r.Scope.String())
	state.Filter = types.StringValue(r.Filter)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.OutputFile = types.StringValue(r.OutputFile)
	state.PreviousFileExtension = internaltypes.StringTypeOrNil(r.PreviousFileExtension, internaltypes.IsEmptyString(expectedValues.PreviousFileExtension))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSearchShutdownPluginOperations(plan searchShutdownPluginResourceModel, state searchShutdownPluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Scope, state.Scope, "scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputFile, state.OutputFile, "output-file")
	operations.AddStringOperationIfNecessary(&ops, plan.PreviousFileExtension, state.PreviousFileExtension, "previous-file-extension")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *searchShutdownPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchShutdownPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope, err := client.NewEnumpluginScopePropFromValue(plan.Scope.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for Scope", err.Error())
		return
	}
	addRequest := client.NewAddSearchShutdownPluginRequest(plan.Id.ValueString(),
		[]client.EnumsearchShutdownPluginSchemaUrn{client.ENUMSEARCHSHUTDOWNPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINSEARCH_SHUTDOWN},
		*scope,
		plan.Filter.ValueString(),
		plan.OutputFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalSearchShutdownPluginFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PluginApi.AddPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPluginRequest(
		client.AddSearchShutdownPluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Search Shutdown Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state searchShutdownPluginResourceModel
	readSearchShutdownPluginResponse(ctx, addResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSearchShutdownPluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan searchShutdownPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Shutdown Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state searchShutdownPluginResourceModel
	readSearchShutdownPluginResponse(ctx, readResponse.SearchShutdownPluginResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSearchShutdownPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Shutdown Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSearchShutdownPluginResponse(ctx, updateResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *searchShutdownPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchShutdownPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSearchShutdownPluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSearchShutdownPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSearchShutdownPlugin(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state searchShutdownPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Search Shutdown Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSearchShutdownPluginResponse(ctx, readResponse.SearchShutdownPluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *searchShutdownPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchShutdownPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSearchShutdownPluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSearchShutdownPlugin(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSearchShutdownPlugin(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan searchShutdownPluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state searchShutdownPluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSearchShutdownPluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Search Shutdown Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSearchShutdownPluginResponse(ctx, updateResponse.SearchShutdownPluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSearchShutdownPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *searchShutdownPluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state searchShutdownPluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Search Shutdown Plugin", err, httpResp)
		return
	}
}

func (r *searchShutdownPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchShutdownPlugin(ctx, req, resp)
}

func (r *defaultSearchShutdownPluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSearchShutdownPlugin(ctx, req, resp)
}

func importSearchShutdownPlugin(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
