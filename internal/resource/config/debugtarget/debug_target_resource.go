package debugtarget

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &debugTargetResource{}
	_ resource.ResourceWithConfigure   = &debugTargetResource{}
	_ resource.ResourceWithImportState = &debugTargetResource{}
	_ resource.Resource                = &defaultDebugTargetResource{}
	_ resource.ResourceWithConfigure   = &defaultDebugTargetResource{}
	_ resource.ResourceWithImportState = &defaultDebugTargetResource{}
)

// Create a Debug Target resource
func NewDebugTargetResource() resource.Resource {
	return &debugTargetResource{}
}

func NewDefaultDebugTargetResource() resource.Resource {
	return &defaultDebugTargetResource{}
}

// debugTargetResource is the resource implementation.
type debugTargetResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDebugTargetResource is the resource implementation.
type defaultDebugTargetResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *debugTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_debug_target"
}

func (r *defaultDebugTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_debug_target"
}

// Configure adds the provider configured client to the resource.
func (r *debugTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDebugTargetResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type debugTargetResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	LogPublisherName         types.String `tfsdk:"log_publisher_name"`
	DebugScope               types.String `tfsdk:"debug_scope"`
	DebugLevel               types.String `tfsdk:"debug_level"`
	DebugCategory            types.Set    `tfsdk:"debug_category"`
	OmitMethodEntryArguments types.Bool   `tfsdk:"omit_method_entry_arguments"`
	OmitMethodReturnValue    types.Bool   `tfsdk:"omit_method_return_value"`
	IncludeThrowableCause    types.Bool   `tfsdk:"include_throwable_cause"`
	ThrowableStackFrames     types.Int64  `tfsdk:"throwable_stack_frames"`
	Description              types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *debugTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	debugTargetSchema(ctx, req, resp, false)
}

func (r *defaultDebugTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	debugTargetSchema(ctx, req, resp, true)
}

func debugTargetSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Debug Target.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Debug Target resource. Options are ['debug-target']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("debug-target"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"debug-target"}...),
				},
			},
			"log_publisher_name": schema.StringAttribute{
				Description: "Name of the parent Log Publisher",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"debug_scope": schema.StringAttribute{
				Description: "Specifies the fully-qualified Java package, class, or method affected by the settings in this target definition. Use the number character (#) to separate the class name and the method name (that is, com.unboundid.directory.server.core.DirectoryServer#startUp).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"debug_level": schema.StringAttribute{
				Description: "Specifies the lowest severity level of debug messages to log.",
				Required:    true,
			},
			"debug_category": schema.SetAttribute{
				Description: "Specifies the debug message categories to be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"omit_method_entry_arguments": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include method arguments in debug messages.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"omit_method_return_value": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include the return value in debug messages.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"include_throwable_cause": schema.BoolAttribute{
				Description: "Specifies the property to indicate whether to include the cause of exceptions in exception thrown and caught messages.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"throwable_stack_frames": schema.Int64Attribute{
				Description: "Specifies the property to indicate the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"description": schema.StringAttribute{
				Description: "A description for this Debug Target",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "debug_scope", "log_publisher_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for debug-target debug-target
func addOptionalDebugTargetFields(ctx context.Context, addRequest *client.AddDebugTargetRequest, plan debugTargetResourceModel) error {
	if internaltypes.IsDefined(plan.DebugCategory) {
		var slice []string
		plan.DebugCategory.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdebugTargetDebugCategoryProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdebugTargetDebugCategoryPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DebugCategory = enumSlice
	}
	if internaltypes.IsDefined(plan.OmitMethodEntryArguments) {
		addRequest.OmitMethodEntryArguments = plan.OmitMethodEntryArguments.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.OmitMethodReturnValue) {
		addRequest.OmitMethodReturnValue = plan.OmitMethodReturnValue.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThrowableCause) {
		addRequest.IncludeThrowableCause = plan.IncludeThrowableCause.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ThrowableStackFrames) {
		addRequest.ThrowableStackFrames = plan.ThrowableStackFrames.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *debugTargetResourceModel) populateAllComputedStringAttributes() {
	if model.DebugScope.IsUnknown() || model.DebugScope.IsNull() {
		model.DebugScope = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.DebugLevel.IsUnknown() || model.DebugLevel.IsNull() {
		model.DebugLevel = types.StringValue("")
	}
}

// Read a DebugTargetResponse object into the model struct
func readDebugTargetResponse(ctx context.Context, r *client.DebugTargetResponse, state *debugTargetResourceModel, expectedValues *debugTargetResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("debug-target")
	state.Id = types.StringValue(r.Id)
	state.DebugScope = types.StringValue(r.DebugScope)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugCategory = internaltypes.GetStringSet(
		client.StringSliceEnumdebugTargetDebugCategoryProp(r.DebugCategory))
	state.OmitMethodEntryArguments = internaltypes.BoolTypeOrNil(r.OmitMethodEntryArguments)
	state.OmitMethodReturnValue = internaltypes.BoolTypeOrNil(r.OmitMethodReturnValue)
	state.IncludeThrowableCause = internaltypes.BoolTypeOrNil(r.IncludeThrowableCause)
	state.ThrowableStackFrames = internaltypes.Int64TypeOrNil(r.ThrowableStackFrames)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *debugTargetResourceModel) setStateValuesNotReturnedByAPI(expectedValues *debugTargetResourceModel) {
	if !expectedValues.LogPublisherName.IsUnknown() {
		state.LogPublisherName = expectedValues.LogPublisherName
	}
}

// Create any update operations necessary to make the state match the plan
func createDebugTargetOperations(plan debugTargetResourceModel, state debugTargetResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DebugScope, state.DebugScope, "debug-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.DebugLevel, state.DebugLevel, "debug-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DebugCategory, state.DebugCategory, "debug-category")
	operations.AddBoolOperationIfNecessary(&ops, plan.OmitMethodEntryArguments, state.OmitMethodEntryArguments, "omit-method-entry-arguments")
	operations.AddBoolOperationIfNecessary(&ops, plan.OmitMethodReturnValue, state.OmitMethodReturnValue, "omit-method-return-value")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThrowableCause, state.IncludeThrowableCause, "include-throwable-cause")
	operations.AddInt64OperationIfNecessary(&ops, plan.ThrowableStackFrames, state.ThrowableStackFrames, "throwable-stack-frames")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a debug-target debug-target
func (r *debugTargetResource) CreateDebugTarget(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan debugTargetResourceModel) (*debugTargetResourceModel, error) {
	debugLevel, err := client.NewEnumdebugTargetDebugLevelPropFromValue(plan.DebugLevel.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for DebugLevel", err.Error())
		return nil, err
	}
	addRequest := client.NewAddDebugTargetRequest(plan.DebugScope.ValueString(),
		plan.DebugScope.ValueString(),
		*debugLevel)
	err = addOptionalDebugTargetFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Debug Target", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DebugTargetApi.AddDebugTarget(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.LogPublisherName.ValueString())
	apiAddRequest = apiAddRequest.AddDebugTargetRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.DebugTargetApi.AddDebugTargetExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Debug Target", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state debugTargetResourceModel
	readDebugTargetResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *debugTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan debugTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateDebugTarget(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	state.setStateValuesNotReturnedByAPI(&plan)
	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultDebugTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan debugTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DebugTargetApi.GetDebugTarget(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.DebugScope.ValueString(), plan.LogPublisherName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Debug Target", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state debugTargetResourceModel
	readDebugTargetResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DebugTargetApi.UpdateDebugTarget(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.DebugScope.ValueString(), plan.LogPublisherName.ValueString())
	ops := createDebugTargetOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DebugTargetApi.UpdateDebugTargetExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Debug Target", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDebugTargetResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *debugTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDebugTarget(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDebugTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDebugTarget(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDebugTarget(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state debugTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DebugTargetApi.GetDebugTarget(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.DebugScope.ValueString(), state.LogPublisherName.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Debug Target", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Debug Target", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDebugTargetResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *debugTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDebugTarget(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDebugTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDebugTarget(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDebugTarget(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan debugTargetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state debugTargetResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DebugTargetApi.UpdateDebugTarget(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.DebugScope.ValueString(), plan.LogPublisherName.ValueString())

	// Determine what update operations are necessary
	ops := createDebugTargetOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DebugTargetApi.UpdateDebugTargetExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Debug Target", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDebugTargetResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultDebugTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *debugTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state debugTargetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DebugTargetApi.DeleteDebugTargetExecute(r.apiClient.DebugTargetApi.DeleteDebugTarget(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.DebugScope.ValueString(), state.LogPublisherName.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Debug Target", err, httpResp)
		return
	}
}

func (r *debugTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDebugTarget(ctx, req, resp)
}

func (r *defaultDebugTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDebugTarget(ctx, req, resp)
}

func importDebugTarget(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [log-publisher-name]/[debug-target-debug-scope]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("log_publisher_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("debug_scope"), split[1])...)
}
