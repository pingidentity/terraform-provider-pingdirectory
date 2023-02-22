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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &composedAttributePluginResource{}
	_ resource.ResourceWithConfigure   = &composedAttributePluginResource{}
	_ resource.ResourceWithImportState = &composedAttributePluginResource{}
)

// Create a Composed Attribute Plugin resource
func NewComposedAttributePluginResource() resource.Resource {
	return &composedAttributePluginResource{}
}

// composedAttributePluginResource is the resource implementation.
type composedAttributePluginResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *composedAttributePluginResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_composed_attribute_plugin"
}

// Configure adds the provider configured client to the resource.
func (r *composedAttributePluginResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type composedAttributePluginResourceModel struct {
	Id                                                   types.String `tfsdk:"id"`
	LastUpdated                                          types.String `tfsdk:"last_updated"`
	Notifications                                        types.Set    `tfsdk:"notifications"`
	RequiredActions                                      types.Set    `tfsdk:"required_actions"`
	PluginType                                           types.Set    `tfsdk:"plugin_type"`
	AttributeType                                        types.String `tfsdk:"attribute_type"`
	ValuePattern                                         types.Set    `tfsdk:"value_pattern"`
	MultipleValuePatternBehavior                         types.String `tfsdk:"multiple_value_pattern_behavior"`
	MultiValuedAttributeBehavior                         types.String `tfsdk:"multi_valued_attribute_behavior"`
	TargetAttributeExistsDuringInitialPopulationBehavior types.String `tfsdk:"target_attribute_exists_during_initial_population_behavior"`
	UpdateSourceAttributeBehavior                        types.String `tfsdk:"update_source_attribute_behavior"`
	SourceAttributeRemovalBehavior                       types.String `tfsdk:"source_attribute_removal_behavior"`
	UpdateTargetAttributeBehavior                        types.String `tfsdk:"update_target_attribute_behavior"`
	IncludeBaseDN                                        types.Set    `tfsdk:"include_base_dn"`
	ExcludeBaseDN                                        types.Set    `tfsdk:"exclude_base_dn"`
	IncludeFilter                                        types.Set    `tfsdk:"include_filter"`
	ExcludeFilter                                        types.Set    `tfsdk:"exclude_filter"`
	UpdatedEntryNewlyMatchesCriteriaBehavior             types.String `tfsdk:"updated_entry_newly_matches_criteria_behavior"`
	UpdatedEntryNoLongerMatchesCriteriaBehavior          types.String `tfsdk:"updated_entry_no_longer_matches_criteria_behavior"`
	Description                                          types.String `tfsdk:"description"`
	Enabled                                              types.Bool   `tfsdk:"enabled"`
	InvokeForInternalOperations                          types.Bool   `tfsdk:"invoke_for_internal_operations"`
}

// GetSchema defines the schema for the resource.
func (r *composedAttributePluginResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Composed Attribute Plugin.",
		Attributes: map[string]schema.Attribute{
			"plugin_type": schema.SetAttribute{
				Description: "Specifies the set of plug-in types for the plug-in, which specifies the times at which the plug-in is invoked.",
				Required:    true,
				ElementType: types.StringType,
			},
			"attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute type for which values are to be generated.",
				Required:    true,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the values to use for the target attribute type.",
				Required:    true,
				ElementType: types.StringType,
			},
			"multiple_value_pattern_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the plugin is configured with multiple value patterns.",
				Optional:    true,
				Computed:    true,
			},
			"multi_valued_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for source attributes that have multiple values.",
				Optional:    true,
				Computed:    true,
			},
			"target_attribute_exists_during_initial_population_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if the target attribute exists when initially populating the entry with composed values (whether during an LDIF import, an add operation, or an invocation of the populate composed attribute values task).",
				Optional:    true,
				Computed:    true,
			},
			"update_source_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update one or more of the source attributes used in any of the value patterns.",
				Optional:    true,
				Computed:    true,
			},
			"source_attribute_removal_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that update an entry to remove source attributes in such a way that this plugin would no longer generate any composed values for that entry.",
				Optional:    true,
				Computed:    true,
			},
			"update_target_attribute_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify and modify DN operations that attempt to update the set of values for the target attribute.",
				Optional:    true,
				Computed:    true,
			},
			"include_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values may be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_base_dn": schema.SetAttribute{
				Description: "The set of base DNs below which composed values will not be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values may be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_filter": schema.SetAttribute{
				Description: "The set of search filters that identify entries for which composed values will not be generated.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"updated_entry_newly_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously did not satisfy either the base DN or filter criteria, but now do satisfy that criteria.",
				Optional:    true,
				Computed:    true,
			},
			"updated_entry_no_longer_matches_criteria_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit for modify or modify DN operations that update an entry that previously satisfied the base DN and filter criteria, but now no longer satisfies that criteria.",
				Optional:    true,
				Computed:    true,
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
func addOptionalComposedAttributePluginFields(ctx context.Context, addRequest *client.AddComposedAttributePluginRequest, plan composedAttributePluginResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleValuePatternBehavior) {
		multipleValuePatternBehavior, err := client.NewEnumpluginMultipleValuePatternBehaviorPropFromValue(plan.MultipleValuePatternBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleValuePatternBehavior = multipleValuePatternBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultiValuedAttributeBehavior) {
		multiValuedAttributeBehavior, err := client.NewEnumpluginMultiValuedAttributeBehaviorPropFromValue(plan.MultiValuedAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultiValuedAttributeBehavior = multiValuedAttributeBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TargetAttributeExistsDuringInitialPopulationBehavior) {
		targetAttributeExistsDuringInitialPopulationBehavior, err := client.NewEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorPropFromValue(plan.TargetAttributeExistsDuringInitialPopulationBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.TargetAttributeExistsDuringInitialPopulationBehavior = targetAttributeExistsDuringInitialPopulationBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateSourceAttributeBehavior) {
		updateSourceAttributeBehavior, err := client.NewEnumpluginUpdateSourceAttributeBehaviorPropFromValue(plan.UpdateSourceAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdateSourceAttributeBehavior = updateSourceAttributeBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SourceAttributeRemovalBehavior) {
		sourceAttributeRemovalBehavior, err := client.NewEnumpluginSourceAttributeRemovalBehaviorPropFromValue(plan.SourceAttributeRemovalBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SourceAttributeRemovalBehavior = sourceAttributeRemovalBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateTargetAttributeBehavior) {
		updateTargetAttributeBehavior, err := client.NewEnumpluginUpdateTargetAttributeBehaviorPropFromValue(plan.UpdateTargetAttributeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdateTargetAttributeBehavior = updateTargetAttributeBehavior
	}
	if internaltypes.IsDefined(plan.IncludeBaseDN) {
		var slice []string
		plan.IncludeBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.IncludeBaseDN = slice
	}
	if internaltypes.IsDefined(plan.ExcludeBaseDN) {
		var slice []string
		plan.ExcludeBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeBaseDN = slice
	}
	if internaltypes.IsDefined(plan.IncludeFilter) {
		var slice []string
		plan.IncludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludeFilter = slice
	}
	if internaltypes.IsDefined(plan.ExcludeFilter) {
		var slice []string
		plan.ExcludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeFilter = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdatedEntryNewlyMatchesCriteriaBehavior) {
		updatedEntryNewlyMatchesCriteriaBehavior, err := client.NewEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorPropFromValue(plan.UpdatedEntryNewlyMatchesCriteriaBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdatedEntryNewlyMatchesCriteriaBehavior = updatedEntryNewlyMatchesCriteriaBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdatedEntryNoLongerMatchesCriteriaBehavior) {
		updatedEntryNoLongerMatchesCriteriaBehavior, err := client.NewEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorPropFromValue(plan.UpdatedEntryNoLongerMatchesCriteriaBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.UpdatedEntryNoLongerMatchesCriteriaBehavior = updatedEntryNoLongerMatchesCriteriaBehavior
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

// Read a ComposedAttributePluginResponse object into the model struct
func readComposedAttributePluginResponse(ctx context.Context, r *client.ComposedAttributePluginResponse, state *composedAttributePluginResourceModel, expectedValues *composedAttributePluginResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PluginType = internaltypes.GetStringSet(
		client.StringSliceEnumpluginPluginTypeProp(r.PluginType))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.MultipleValuePatternBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultipleValuePatternBehaviorProp(r.MultipleValuePatternBehavior), internaltypes.IsEmptyString(expectedValues.MultipleValuePatternBehavior))
	state.MultiValuedAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginMultiValuedAttributeBehaviorProp(r.MultiValuedAttributeBehavior), internaltypes.IsEmptyString(expectedValues.MultiValuedAttributeBehavior))
	state.TargetAttributeExistsDuringInitialPopulationBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginTargetAttributeExistsDuringInitialPopulationBehaviorProp(r.TargetAttributeExistsDuringInitialPopulationBehavior), internaltypes.IsEmptyString(expectedValues.TargetAttributeExistsDuringInitialPopulationBehavior))
	state.UpdateSourceAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateSourceAttributeBehaviorProp(r.UpdateSourceAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateSourceAttributeBehavior))
	state.SourceAttributeRemovalBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginSourceAttributeRemovalBehaviorProp(r.SourceAttributeRemovalBehavior), internaltypes.IsEmptyString(expectedValues.SourceAttributeRemovalBehavior))
	state.UpdateTargetAttributeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdateTargetAttributeBehaviorProp(r.UpdateTargetAttributeBehavior), internaltypes.IsEmptyString(expectedValues.UpdateTargetAttributeBehavior))
	state.IncludeBaseDN = internaltypes.GetStringSet(r.IncludeBaseDN)
	state.ExcludeBaseDN = internaltypes.GetStringSet(r.ExcludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.UpdatedEntryNewlyMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNewlyMatchesCriteriaBehaviorProp(r.UpdatedEntryNewlyMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNewlyMatchesCriteriaBehavior))
	state.UpdatedEntryNoLongerMatchesCriteriaBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpluginUpdatedEntryNoLongerMatchesCriteriaBehaviorProp(r.UpdatedEntryNoLongerMatchesCriteriaBehavior), internaltypes.IsEmptyString(expectedValues.UpdatedEntryNoLongerMatchesCriteriaBehavior))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.InvokeForInternalOperations = internaltypes.BoolTypeOrNil(r.InvokeForInternalOperations)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createComposedAttributePluginOperations(plan composedAttributePluginResourceModel, state composedAttributePluginResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PluginType, state.PluginType, "plugin-type")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleValuePatternBehavior, state.MultipleValuePatternBehavior, "multiple-value-pattern-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.MultiValuedAttributeBehavior, state.MultiValuedAttributeBehavior, "multi-valued-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetAttributeExistsDuringInitialPopulationBehavior, state.TargetAttributeExistsDuringInitialPopulationBehavior, "target-attribute-exists-during-initial-population-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateSourceAttributeBehavior, state.UpdateSourceAttributeBehavior, "update-source-attribute-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.SourceAttributeRemovalBehavior, state.SourceAttributeRemovalBehavior, "source-attribute-removal-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateTargetAttributeBehavior, state.UpdateTargetAttributeBehavior, "update-target-attribute-behavior")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeBaseDN, state.IncludeBaseDN, "include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeBaseDN, state.ExcludeBaseDN, "exclude-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeFilter, state.ExcludeFilter, "exclude-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNewlyMatchesCriteriaBehavior, state.UpdatedEntryNewlyMatchesCriteriaBehavior, "updated-entry-newly-matches-criteria-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdatedEntryNoLongerMatchesCriteriaBehavior, state.UpdatedEntryNoLongerMatchesCriteriaBehavior, "updated-entry-no-longer-matches-criteria-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForInternalOperations, state.InvokeForInternalOperations, "invoke-for-internal-operations")
	return ops
}

// Create a new resource
func (r *composedAttributePluginResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan composedAttributePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var PluginTypeSlice []client.EnumpluginPluginTypeProp
	plan.PluginType.ElementsAs(ctx, &PluginTypeSlice, false)
	var ValuePatternSlice []string
	plan.ValuePattern.ElementsAs(ctx, &ValuePatternSlice, false)
	addRequest := client.NewAddComposedAttributePluginRequest(plan.Id.ValueString(),
		[]client.EnumcomposedAttributePluginSchemaUrn{client.ENUMCOMPOSEDATTRIBUTEPLUGINSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PLUGINCOMPOSED_ATTRIBUTE},
		PluginTypeSlice,
		plan.AttributeType.ValueString(),
		ValuePatternSlice,
		plan.Enabled.ValueBool())
	err := addOptionalComposedAttributePluginFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Composed Attribute Plugin", err.Error())
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
		client.AddComposedAttributePluginRequestAsAddPluginRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PluginApi.AddPluginExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Composed Attribute Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state composedAttributePluginResourceModel
	readComposedAttributePluginResponse(ctx, addResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)

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
func (r *composedAttributePluginResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state composedAttributePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PluginApi.GetPlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Composed Attribute Plugin", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readComposedAttributePluginResponse(ctx, readResponse.ComposedAttributePluginResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *composedAttributePluginResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan composedAttributePluginResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state composedAttributePluginResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.PluginApi.UpdatePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createComposedAttributePluginOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PluginApi.UpdatePluginExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Composed Attribute Plugin", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readComposedAttributePluginResponse(ctx, updateResponse.ComposedAttributePluginResponse, &state, &plan, &resp.Diagnostics)
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
func (r *composedAttributePluginResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state composedAttributePluginResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PluginApi.DeletePluginExecute(r.apiClient.PluginApi.DeletePlugin(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Composed Attribute Plugin", err, httpResp)
		return
	}
}

func (r *composedAttributePluginResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
