package config

import (
	   "context"
	   "strings"
	   "time"

	   "github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	   internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

    "github.com/hashicorp/terraform-plugin-framework/diag"
	   "github.com/hashicorp/terraform-plugin-framework/path"
	   "github.com/hashicorp/terraform-plugin-framework/resource"
	   "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	   "github.com/hashicorp/terraform-plugin-framework/types"
	   "github.com/hashicorp/terraform-plugin-log/tflog"
	   client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	   _ resource.Resource                = &debugTargetResource{}
	   _ resource.ResourceWithConfigure   = &debugTargetResource{}
	   _ resource.ResourceWithImportState = &debugTargetResource{}
)

// Create a Debug Target resource
func NewDebugTargetResource() resource.Resource {
	   return &debugTargetResource{}
}

// debugTargetResource is the resource implementation.
type debugTargetResource struct {
	   providerConfig internaltypes.ProviderConfiguration
	   apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *debugTargetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	   resp.TypeName = req.ProviderTypeName + "_debug_target"
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

type debugTargetResourceModel struct {
    Id              types.String `tfsdk:"id"`
    LastUpdated     types.String `tfsdk:"last_updated"`
    Notifications   types.Set    `tfsdk:"notifications"`
    RequiredActions types.Set    `tfsdk:"required_actions"`
    LogPublisherName              types.String `tfsdk:"log_publisher_name"`
    DebugScope types.String `tfsdk:"debug_scope"`
    DebugLevel types.String `tfsdk:"debug_level"`
    DebugCategory types.Set `tfsdk:"debug_category"`
    OmitMethodEntryArguments types.Bool `tfsdk:"omit_method_entry_arguments"`
    OmitMethodReturnValue types.Bool `tfsdk:"omit_method_return_value"`
    IncludeThrowableCause types.Bool `tfsdk:"include_throwable_cause"`
    ThrowableStackFrames types.Int64 `tfsdk:"throwable_stack_frames"`
    Description types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *debugTargetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    schema := schema.Schema{
        Description: "Manages a Debug Target.",
        Attributes: map[string]schema.Attribute{
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
            },
            "debug_level": schema.StringAttribute{
                Description: "Specifies the lowest severity level of debug messages to log.",
                Required:    true,
            },
            "debug_category": schema.SetAttribute{
                Description: "Specifies the debug message categories to be logged.",
                Optional:    true,
                Computed:    true,
                ElementType:    types.StringType,
            },
            "omit_method_entry_arguments": schema.BoolAttribute{
                Description: "Specifies the property to indicate whether to include method arguments in debug messages.",
                Optional:    true,
                Computed:    true,
            },
            "omit_method_return_value": schema.BoolAttribute{
                Description: "Specifies the property to indicate whether to include the return value in debug messages.",
                Optional:    true,
                Computed:    true,
            },
            "include_throwable_cause": schema.BoolAttribute{
                Description: "Specifies the property to indicate whether to include the cause of exceptions in exception thrown and caught messages.",
                Optional:    true,
                Computed:    true,
            },
            "throwable_stack_frames": schema.Int64Attribute{
                Description: "Specifies the property to indicate the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
                Optional:    true,
                Computed:    true,
            },
            "description": schema.StringAttribute{
                Description: "A description for this Debug Target",
                Optional:    true,
            },
        },
    }
    AddCommonSchema(&schema, true)
    resp.Schema = schema
}

// Add optional fields to create request
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
        boolVal := plan.OmitMethodEntryArguments.ValueBool()
        addRequest.OmitMethodEntryArguments = &boolVal
    }
    if internaltypes.IsDefined(plan.OmitMethodReturnValue) {
        boolVal := plan.OmitMethodReturnValue.ValueBool()
        addRequest.OmitMethodReturnValue = &boolVal
    }
    if internaltypes.IsDefined(plan.IncludeThrowableCause) {
        boolVal := plan.IncludeThrowableCause.ValueBool()
        addRequest.IncludeThrowableCause = &boolVal
    }
    if internaltypes.IsDefined(plan.ThrowableStackFrames) {
        intVal := int32(plan.ThrowableStackFrames.ValueInt64())
        addRequest.ThrowableStackFrames = &intVal
    }
    // Empty strings are treated as equivalent to null
    if internaltypes.IsNonEmptyString(plan.Description) {
        stringVal := plan.Description.ValueString()
        addRequest.Description = &stringVal
    }
    return nil
}

// Read a DebugTargetResponse object into the model struct
func readDebugTargetResponse(ctx context.Context, r *client.DebugTargetResponse, state *debugTargetResourceModel, expectedValues *debugTargetResourceModel, diagnostics *diag.Diagnostics) {
    state.Id = types.StringValue(r.Id)
    state.LogPublisherName = expectedValues.LogPublisherName
    state.DebugScope = types.StringValue(r.DebugScope)
    state.DebugLevel = types.StringValue(r.DebugLevel.String())
    state.DebugCategory = internaltypes.GetStringSet(
        client.StringSliceEnumdebugTargetDebugCategoryProp(r.DebugCategory))
    state.OmitMethodEntryArguments = internaltypes.BoolTypeOrNil(r.OmitMethodEntryArguments)
    state.OmitMethodReturnValue = internaltypes.BoolTypeOrNil(r.OmitMethodReturnValue)
    state.IncludeThrowableCause = internaltypes.BoolTypeOrNil(r.IncludeThrowableCause)
    state.ThrowableStackFrames = internaltypes.Int64TypeOrNil(r.ThrowableStackFrames)
    state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
    state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
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

// Create a new resource
func (r *debugTargetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    // Retrieve values from plan
    var plan debugTargetResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    debugLevel, err := client.NewEnumdebugTargetDebugLevelPropFromValue(plan.DebugLevel.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Failed to parse enum value for DebugLevel", err.Error())
        return
    }
    addRequest := client.NewAddDebugTargetRequest(plan.Id.ValueString(),
        plan.DebugScope.ValueString(),
        *debugLevel)
    err = addOptionalDebugTargetFields(ctx, addRequest, plan)
    if err != nil {
        resp.Diagnostics.AddError("Failed to add optional properties to add request for Debug Target", err.Error())
        return
    }
    // Log request JSON
    requestJson, err := addRequest.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add request: "+string(requestJson))
    }
    apiAddRequest := r.apiClient.DebugTargetApi.AddDebugTarget(
        ProviderBasicAuthContext(ctx, r.providerConfig), plan.LogPublisherName.ValueString())
    apiAddRequest = apiAddRequest.AddDebugTargetRequest(*addRequest)

    addResponse, httpResp, err := r.apiClient.DebugTargetApi.AddDebugTargetExecute(apiAddRequest)
    if err != nil {
        ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Debug Target", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := addResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Add response: "+string(responseJson))
    }

    // Read the response into the state
    var state debugTargetResourceModel
    readDebugTargetResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

    // Populate Computed attribute values
    state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

    // Set state to fully populated data
    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Read resource information
func (r *debugTargetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    // Get current state
  var state debugTargetResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    readResponse, httpResp, err := r.apiClient.DebugTargetApi.GetDebugTarget(
        ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.LogPublisherName.ValueString()).Execute()
    if err != nil {
        ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Debug Target", err, httpResp)
        return
    }

    // Log response JSON
    responseJson, err := readResponse.MarshalJSON()
    if err == nil {
        tflog.Debug(ctx, "Read response: "+string(responseJson))
    }

    // Read the response into the state
    readDebugTargetResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

    // Set refreshed state
    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

// Update a resource
func (r *debugTargetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	   updateRequest := r.apiClient.DebugTargetApi.UpdateDebugTarget(
        ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.LogPublisherName.ValueString())

	   // Determine what update operations are necessary
	   ops := createDebugTargetOperations(plan, state)
	   if len(ops) > 0 {
	   	   updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
	   	   // Log operations
	   	   operations.LogUpdateOperations(ctx, ops)

        updateResponse, httpResp, err := r.apiClient.DebugTargetApi.UpdateDebugTargetExecute(updateRequest)
        if err != nil {
            ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Debug Target", err, httpResp)
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

    diags = resp.State.Set(ctx, state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
    	   return
    }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *debugTargetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    // Retrieve values from state
    var state debugTargetResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
    	   return
    }

    httpResp, err := r.apiClient.DebugTargetApi.DeleteDebugTargetExecute(r.apiClient.DebugTargetApi.DeleteDebugTarget(
        ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.LogPublisherName.ValueString()))
    if err != nil {
        ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Debug Target", err, httpResp)
        return
    }
}

func (r *debugTargetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    split := strings.Split(req.ID, "/")
    if len(split) != 2 {
      resp.Diagnostics.AddError("Invalid import id for resource", "Expected [log-publisher-name]/[debug-target-name]. Got: "+req.ID)
      return
    }
    // Set the required attributes to read the resource
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("log_publisher_name"), split[0])...)
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}

