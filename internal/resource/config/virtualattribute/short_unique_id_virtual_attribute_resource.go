package virtualattribute

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
	_ resource.Resource                = &shortUniqueIdVirtualAttributeResource{}
	_ resource.ResourceWithConfigure   = &shortUniqueIdVirtualAttributeResource{}
	_ resource.ResourceWithImportState = &shortUniqueIdVirtualAttributeResource{}
)

// Create a Short Unique Id Virtual Attribute resource
func NewShortUniqueIdVirtualAttributeResource() resource.Resource {
	return &shortUniqueIdVirtualAttributeResource{}
}

// shortUniqueIdVirtualAttributeResource is the resource implementation.
type shortUniqueIdVirtualAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *shortUniqueIdVirtualAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_short_unique_id_virtual_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *shortUniqueIdVirtualAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type shortUniqueIdVirtualAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	SequenceNumberAttribute                      types.String `tfsdk:"sequence_number_attribute"`
	Description                                  types.String `tfsdk:"description"`
	Enabled                                      types.Bool   `tfsdk:"enabled"`
	AttributeType                                types.String `tfsdk:"attribute_type"`
	BaseDN                                       types.Set    `tfsdk:"base_dn"`
	GroupDN                                      types.Set    `tfsdk:"group_dn"`
	Filter                                       types.Set    `tfsdk:"filter"`
	ClientConnectionPolicy                       types.Set    `tfsdk:"client_connection_policy"`
	MultipleVirtualAttributeEvaluationOrderIndex types.Int64  `tfsdk:"multiple_virtual_attribute_evaluation_order_index"`
	MultipleVirtualAttributeMergeBehavior        types.String `tfsdk:"multiple_virtual_attribute_merge_behavior"`
	AllowIndexConflicts                          types.Bool   `tfsdk:"allow_index_conflicts"`
}

// GetSchema defines the schema for the resource.
func (r *shortUniqueIdVirtualAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Short Unique Id Virtual Attribute.",
		Attributes: map[string]schema.Attribute{
			"sequence_number_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute which contains the sequence number from which unique identifiers are generated. The attribute should have Integer syntax or a String syntax permitting integer values. If this property is modified then the filter property should be updated accordingly so that only entries containing the sequence number attribute are eligible to have a value generated for this virtual attribute.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Virtual Attribute",
				Optional:    true,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Virtual Attribute is enabled for use.",
				Optional:    true,
				Computed:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be dynamically assigned by the virtual attribute.",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs for the branches containing entries that are eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"group_dn": schema.SetAttribute{
				Description: "Specifies the DNs of the groups whose members can be eligible to use this virtual attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"filter": schema.SetAttribute{
				Description: "Specifies the search filters to be applied against entries to determine if the virtual attribute is to be generated for those entries.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"client_connection_policy": schema.SetAttribute{
				Description: "Specifies a set of client connection policies for which this Virtual Attribute should be generated. If this is undefined, then this Virtual Attribute will always be generated. If it is associated with one or more client connection policies, then this Virtual Attribute will be generated only for operations requested by clients assigned to one of those client connection policies.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"multiple_virtual_attribute_evaluation_order_index": schema.Int64Attribute{
				Description: "Specifies the order in which virtual attribute definitions for the same attribute type will be evaluated when generating values for an entry.",
				Optional:    true,
				Computed:    true,
			},
			"multiple_virtual_attribute_merge_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that will be exhibited for cases in which multiple virtual attribute definitions apply to the same multivalued attribute type. This will be ignored for single-valued attribute types.",
				Optional:    true,
				Computed:    true,
			},
			"allow_index_conflicts": schema.BoolAttribute{
				Description: "Indicates whether the server should allow creating or altering this virtual attribute definition even if it conflicts with one or more indexes defined in the server.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a ShortUniqueIdVirtualAttributeResponse object into the model struct
func readShortUniqueIdVirtualAttributeResponse(ctx context.Context, r *client.ShortUniqueIdVirtualAttributeResponse, state *shortUniqueIdVirtualAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SequenceNumberAttribute = types.StringValue(r.SequenceNumberAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.GroupDN = internaltypes.GetStringSet(r.GroupDN)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.ClientConnectionPolicy = internaltypes.GetStringSet(r.ClientConnectionPolicy)
	state.MultipleVirtualAttributeEvaluationOrderIndex = internaltypes.Int64TypeOrNil(r.MultipleVirtualAttributeEvaluationOrderIndex)
	state.MultipleVirtualAttributeMergeBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvirtualAttributeMultipleVirtualAttributeMergeBehaviorProp(r.MultipleVirtualAttributeMergeBehavior), true)
	state.AllowIndexConflicts = internaltypes.BoolTypeOrNil(r.AllowIndexConflicts)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createShortUniqueIdVirtualAttributeOperations(plan shortUniqueIdVirtualAttributeResourceModel, state shortUniqueIdVirtualAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.SequenceNumberAttribute, state.SequenceNumberAttribute, "sequence-number-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.GroupDN, state.GroupDN, "group-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ClientConnectionPolicy, state.ClientConnectionPolicy, "client-connection-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.MultipleVirtualAttributeEvaluationOrderIndex, state.MultipleVirtualAttributeEvaluationOrderIndex, "multiple-virtual-attribute-evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.MultipleVirtualAttributeMergeBehavior, state.MultipleVirtualAttributeMergeBehavior, "multiple-virtual-attribute-merge-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowIndexConflicts, state.AllowIndexConflicts, "allow-index-conflicts")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *shortUniqueIdVirtualAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan shortUniqueIdVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Short Unique Id Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state shortUniqueIdVirtualAttributeResourceModel
	readShortUniqueIdVirtualAttributeResponse(ctx, readResponse.ShortUniqueIdVirtualAttributeResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createShortUniqueIdVirtualAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Short Unique Id Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readShortUniqueIdVirtualAttributeResponse(ctx, updateResponse.ShortUniqueIdVirtualAttributeResponse, &state, &resp.Diagnostics)
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
func (r *shortUniqueIdVirtualAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state shortUniqueIdVirtualAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VirtualAttributeApi.GetVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Short Unique Id Virtual Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readShortUniqueIdVirtualAttributeResponse(ctx, readResponse.ShortUniqueIdVirtualAttributeResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *shortUniqueIdVirtualAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan shortUniqueIdVirtualAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state shortUniqueIdVirtualAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.VirtualAttributeApi.UpdateVirtualAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createShortUniqueIdVirtualAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VirtualAttributeApi.UpdateVirtualAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Short Unique Id Virtual Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readShortUniqueIdVirtualAttributeResponse(ctx, updateResponse.ShortUniqueIdVirtualAttributeResponse, &state, &resp.Diagnostics)
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
func (r *shortUniqueIdVirtualAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *shortUniqueIdVirtualAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
