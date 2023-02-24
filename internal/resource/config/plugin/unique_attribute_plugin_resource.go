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
	_ resource.Resource                = &uniqueAttributePluginResource{}
	_ resource.ResourceWithConfigure   = &uniqueAttributePluginResource{}
	_ resource.ResourceWithImportState = &uniqueAttributePluginResource{}
)

// Create a Unique Attribute Plugin resource
func NewUniqueAttributePluginResource() resource.Resource {
	return &uniqueAttributePluginResource{}
}

// uniqueAttributePluginResource is the resource implementation.
type uniqueAttributePluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *uniqueAttributePluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unique_attribute_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *uniqueAttributePluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type uniqueAttributePluginResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	PluginType                             types.Set    `tfsdk:"plugin_type"`
	Type                                   types.Set    `tfsdk:"type"`
	MultipleAttributeBehavior              types.String `tfsdk:"multiple_attribute_behavior"`
	BaseDN                                 types.Set    `tfsdk:"base_dn"`
	PreventConflictsWithSoftDeletedEntries types.Bool   `tfsdk:"prevent_conflicts_with_soft_deleted_entries"`
	Filter                                 types.String `tfsdk:"filter"`
	Description                            types.String `tfsdk:"description"`
	Enabled                                types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations            types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *uniqueAttributePluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Unique Attribute Plugin.",
		Attributes: map[string]schema.Attribute{
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"type": schema.SetAttribute{
				Description: "Specifies the type of attributes to check for value uniqueness.",
				Required:    true,
				ElementType: types.StringType,
			},
			"multiple_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if multiple attribute types are specified.",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies a base DN within which the attribute must be unique.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"prevent_conflicts_with_soft_deleted_entries": schema.BoolAttribute{
				Description: "Indicates whether this Unique Attribute Plugin should reject a change that would result in one or more conflicts, even if those conflicts only exist in soft-deleted entries.",
				Optional:    true,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "Specifies the search filter to apply to determine if attribute uniqueness is enforced for the matching entries.",
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
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalUniqueAttributePluginFields(ctx context.Context, addRequest *client.AddUniqueAttributePluginRequest, plan uniqueAttributePluginResourceModel) error {
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
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleAttributeBehavior) {
		multipleAttributeBehavior, err := client.NewEnumpluginMultipleAttributeBehaviorPropFromValue(plan.MultipleAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleAttributeBehavior = multipleAttributeBehavior
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.PreventConflictsWithSoftDeletedEntries) {
		boolVal := plan.PreventConflictsWithSoftDeletedEntries.ValueBool()
		addRequest.PreventConflictsWithSoftDeletedEntries = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Filter) {
		stringVal := plan.Filter.ValueString()
		addRequest.Filter = &stringVal
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

// Read a UniqueAttributePluginResponse object into the model struct
func readUniqueAttributePluginResponse(ctx context.Context, r *client.UniqueAttributePluginResponse, state *uniqueAttributePluginResourceModel, expectedValues *uniqueAttributePluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.Type = internaltypes.GetStringSet(r.Type)
	state.MultipleAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleAttributeBehaviorProp(r.MultipleAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultipleAttributeBehavior))
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.PreventConflictsWithSoftDeletedEntries = internaltypes.BoolTypeOrNil(r.PreventConflictsWithSoftDeletedEntries)
	state.Filter = internaltypes.StringTypeOrNil(r.Filter, internaltypes.IsEmptyString(expectedValues.Filter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createUniqueAttributePluginOperations(plan uniqueAttributePluginResourceModel, state uniqueAttributePluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Type, state.Type, "type")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleAttributeBehavior, state.MultipleAttributeBehavior, "multiple-attribute-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventConflictsWithSoftDeletedEntries, state.PreventConflictsWithSoftDeletedEntries, "prevent-conflicts-with-soft-deleted-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *uniqueAttributePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan uniqueAttributePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var TypeSlice []string
	plan.Type.ElementsAs(ctx, &TypeSlice, false)
	addRequest := client.NewAddUniqueAttributePluginRequest(plan.Id.ValueString(),
		[]client.EnumuniqueAttributePluginSchemaUrn{client.ENUMUNIQUEATTRIBUTEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINUNIQUE_ATTRIBUTE},
		TypeSlice,
		plan.Enabled.ValueBool())
	err := addOptionalUniqueAttributePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Unique Attribute Plugin", err.Error())
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
		client.AddUniqueAttributePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Unique Attribute Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uniqueAttributePluginResourceModel
	readUniqueAttributePluginResponse(ctx, addResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *uniqueAttributePluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state uniqueAttributePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unique Attribute Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readUniqueAttributePluginResponse(ctx, readResponse.UniqueAttributePluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *uniqueAttributePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan uniqueAttributePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state uniqueAttributePluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createUniqueAttributePluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unique Attribute Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUniqueAttributePluginResponse(ctx, updateResponse.UniqueAttributePluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *uniqueAttributePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state uniqueAttributePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Unique Attribute Plugin", err, httpResp)
		return
	}
}

func (r *uniqueAttributePluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
