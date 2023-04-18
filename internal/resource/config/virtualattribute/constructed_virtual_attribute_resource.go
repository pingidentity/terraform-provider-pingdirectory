package virtualattribute

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &constructedVirtualAttributeResource{}
	_ resource.ResourceWithConfigure   = &constructedVirtualAttributeResource{}
	_ resource.ResourceWithImportState = &constructedVirtualAttributeResource{}
	_ resource.Resource                = &defaultConstructedVirtualAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultConstructedVirtualAttributeResource{}
	_ resource.ResourceWithImportState = &defaultConstructedVirtualAttributeResource{}
)

// Create a Constructed Virtual Attribute resource
func NewConstructedVirtualAttributeResource() resource.Resource {
	return &constructedVirtualAttributeResource{}
}

func NewDefaultConstructedVirtualAttributeResource() resource.Resource {
	return &defaultConstructedVirtualAttributeResource{}
}

// constructedVirtualAttributeResource is the resource implementation.
type constructedVirtualAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConstructedVirtualAttributeResource is the resource implementation.
type defaultConstructedVirtualAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *constructedVirtualAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_constructed_virtual_attribute"
}

func (r *defaultConstructedVirtualAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_constructed_virtual_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *constructedVirtualAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultConstructedVirtualAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type constructedVirtualAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	ValuePattern                                 types.Set    `tfsdk:"value_pattern"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
	AttributeType                                types.String `tfsdk:"attribute_type"`
	BaseDN                                       types.Set    `tfsdk:"base_dn"`
	GroupDN                                      types.Set    `tfsdk:"group_dn"`
	Filter                                       types.Set    `tfsdk:"filter"`
	ClientConnectionPolicy                       types.Set    `tfsdk:"client_connection_policy"`
	ConflictBehavior                             types.String `tfsdk:"conflict_behavior"`
	RequireExplicitRequestByName                 types.Bool   `tfsdk:"require_explicit_request_by_name"`
	MultipleVirtualAttributeEvaluationOrderIndex types.Int64  `tfsdk:"multiple_virtual_attribute_evaluation_order_index"`
	MultipleVirtualAttributeMergeBehavior        types.String `tfsdk:"multiple_virtual_attribute_merge_behavior"`
	AllowIndexConflicts                          types.Bool   `tfsdk:"allow_index_conflicts"`
}

// GetSchema defines the schema for the resource.
func (r *constructedVirtualAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	constructedVirtualAttributeSchema(ctx, req, resp, false)
}

func (r *defaultConstructedVirtualAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	constructedVirtualAttributeSchema(ctx, req, resp, true)
}

func constructedVirtualAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Constructed Virtual Attribute.",
		Attributes: map[string]schema.Attribute{
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the virtual attribute value using fixed text and attribute values from the entry.",
				Required:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Virtual Attribute",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Virtual Attribute is enabled for use.",
				Required:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.",
				Required:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"group_dn": schema.SetAttribute{
				Description: "Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"filter": schema.SetAttribute{
				Description: "Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"client_connection_policy": schema.SetAttribute{
				Description: "Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"conflict_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server is to exhibit for entries that already contain one or more real values for the associated attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"require_explicit_request_by_name": schema.BoolAttribute{
				Description: "Indicates whether attributes of this type must be explicitly included by name in the list of requested attributes. Note that this will only apply to virtual attributes which are associated with an attribute type that is operational. It will be ignored for virtual attributes associated with a non-operational attribute type.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"multiple_virtual_attribute_evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.",
				Optional:    true,
			},
			"multiple_virtual_attribute_merge_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_index_conflicts": schema.BoolAttribute{
				Description: "Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalConstructedVirtualAttributeFields(ctx context.Context, addRequest *client.AddConstructedVirtualAttributeRequest, plan constructedVirtualAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BaseDN) {
		var slice []string
		plan.BaseDN.ElementsAs(ctx, &slice, false)
		addRequest.BaseDN = slice
	}
	if internaltypes.IsDefined(plan.GroupDN) {
		var slice []string
		plan.GroupDN.ElementsAs(ctx, &slice, false)
		addRequest.GroupDN = slice
	}
	if internaltypes.IsDefined(plan.Filter) {
		var slice []string
		plan.Filter.ElementsAs(ctx, &slice, false)
		addRequest.Filter = slice
	}
	if internaltypes.IsDefined(plan.ClientConnectionPolicy) {
		var slice []string
		plan.ClientConnectionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.ClientConnectionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConflictBehavior) {
		conflictBehavior, err := client.NewEnumvirtualAttributeConflictBehaviorPropFromValue(plan.ConflictBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConflictBehavior = conflictBehavior
	}
	if internaltypes.IsDefined(plan.RequireExplicitRequestByName) {
		addRequest.RequireExplicitRequestByName = plan.RequireExplicitRequestByName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MultipleVirtualAttributeEvaluationOrderIndex) {
		addRequest.MultipleVirtualAttributeEvaluationOrderIndex = plan.MultipleVirtualAttributeEvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MultipleVirtualAttributeMergeBehavior) {
		multipleVirtualAttributeMergeBehavior, err := client.NewEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorPropFromValue(plan.MultipleVirtualAttributeMergeBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.MultipleVirtualAttributeMergeBehavior = multipleVirtualAttributeMergeBehavior
	}
	if internaltypes.IsDefined(plan.AllowIndexConflicts) {
		addRequest.AllowIndexConflicts = plan.AllowIndexConflicts.ValueBoolPointer()
	}
	return nil
}

// Read a ConstructedVirtualAttributeResponse object into the model struct
func readConstructedVirtualAttributeResponse(ctx context.Context, r *client.ConstructedVirtualAttributeResponse, state *constructedVirtualAttributeResourceModel, expectedValues *constructedVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.ConflictBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeConflictBehaviorProp(r.ConflictBehavior), internaltypes.IsEmptyString(expectedValues.ConflictBehavior))
	state.RequireExplicitRequestByName = internaltypes.BoolTypeOrNil(r.RequireExplicitRequestByName)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), internaltypes.IsEmptyString(expectedValues.MultipleVirtualAttributeMergeBehavior))
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConstructedVirtualAttributeOperations(plan constructedVirtualAttributeResourceModel, state constructedVirtualAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ValuePattern, state.ValuePattern, "value-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GroupDN, state.GroupDN, "group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ClientConnectionPolicy, state.ClientConnectionPolicy, "client-connection-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.ConflictBehavior, state.ConflictBehavior, "conflict-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireExplicitRequestByName, state.RequireExplicitRequestByName, "require-explicit-request-by-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.MultipleVirtualAttributeEvaluationOrderIndex, state.MultipleVirtualAttributeEvaluationOrderIndex, "multiple-virtual-attribute-evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleVirtualAttributeMergeBehavior, state.MultipleVirtualAttributeMergeBehavior, "multiple-virtual-attribute-merge-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowIndexConflicts, state.AllowIndexConflicts, "allow-index-conflicts")
	return ops
}

// Create a new resource
func (r *constructedVirtualAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan constructedVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var ValuePatternSlice []string
	plan.ValuePattern.ElementsAs(ctx, &ValuePatternSlice, false)
	addRequest := client.NewAddConstructedVirtualAttributeRequest(plan.Id.ValueString(),
		[]client.EnumconstructedVirtualAttributeSchemaUrn{client.ENUMCONSTRUCTEDVIRTUALATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VIRTUAL_ATTRIBUTECONSTRUCTED},
		ValuePatternSlice,
		plan.Enabled.ValueBool(),
		plan.AttributeType.ValueString())
	err := addOptionalConstructedVirtualAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Constructed Virtual Attribute", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VirtualAttributeApi.AddVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddVirtualAttributeRequest(
		client.AddConstructedVirtualAttributeRequestAsAddVirtualAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VirtualAttributeApi.AddVirtualAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Constructed Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state constructedVirtualAttributeResourceModel
	readConstructedVirtualAttributeResponse(ctx, addResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultConstructedVirtualAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan constructedVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Constructed Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state constructedVirtualAttributeResourceModel
	readConstructedVirtualAttributeResponse(ctx, readResponse.ConstructedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConstructedVirtualAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Constructed Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConstructedVirtualAttributeResponse(ctx, updateResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *constructedVirtualAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConstructedVirtualAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConstructedVirtualAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConstructedVirtualAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConstructedVirtualAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state constructedVirtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Constructed Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConstructedVirtualAttributeResponse(ctx, readResponse.ConstructedVirtualAttributeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *constructedVirtualAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConstructedVirtualAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConstructedVirtualAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConstructedVirtualAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConstructedVirtualAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan constructedVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state constructedVirtualAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.VirtualAttributeApi.UpdateVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConstructedVirtualAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Constructed Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConstructedVirtualAttributeResponse(ctx, updateResponse.ConstructedVirtualAttributeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConstructedVirtualAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *constructedVirtualAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state constructedVirtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VirtualAttributeApi.DeleteVirtualAttributeExecute(r.apiClient.VirtualAttributeApi.DeleteVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Constructed Virtual Attribute", err, httpResp)
		return
	}
}

func (r *constructedVirtualAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConstructedVirtualAttribute(ctx, req, resp)
}

func (r *defaultConstructedVirtualAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConstructedVirtualAttribute(ctx, req, resp)
}

func importConstructedVirtualAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
