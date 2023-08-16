package logfieldsyntax

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logFieldSyntaxResource{}
	_ resource.ResourceWithConfigure   = &logFieldSyntaxResource{}
	_ resource.ResourceWithImportState = &logFieldSyntaxResource{}
)

// Create a Log Field Syntax resource
func NewLogFieldSyntaxResource() resource.Resource {
	return &logFieldSyntaxResource{}
}

// logFieldSyntaxResource is the resource implementation.
type logFieldSyntaxResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logFieldSyntaxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_field_syntax"
}

// Configure adds the provider configured client to the resource.
func (r *logFieldSyntaxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type logFieldSyntaxResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	Type                       types.String `tfsdk:"type"`
	IncludedSensitiveAttribute types.Set    `tfsdk:"included_sensitive_attribute"`
	ExcludedSensitiveAttribute types.Set    `tfsdk:"excluded_sensitive_attribute"`
	IncludedSensitiveField     types.Set    `tfsdk:"included_sensitive_field"`
	ExcludedSensitiveField     types.Set    `tfsdk:"excluded_sensitive_field"`
	Description                types.String `tfsdk:"description"`
	DefaultBehavior            types.String `tfsdk:"default_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *logFieldSyntaxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Field Syntax.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Syntax resource. Options are ['json', 'attribute-based', 'generic']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"json", "attribute-based", "generic"}...),
				},
			},
			"included_sensitive_attribute": schema.SetAttribute{
				Description: "The set of attribute types that will be considered sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_sensitive_attribute": schema.SetAttribute{
				Description: "The set of attribute types that will not be considered sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"included_sensitive_field": schema.SetAttribute{
				Description: "The names of the JSON fields that will be considered sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_sensitive_field": schema.SetAttribute{
				Description: "The names of the JSON fields that will not be considered sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Syntax",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_behavior": schema.StringAttribute{
				Description: "The default behavior that the server should exhibit when logging fields with this syntax. This may be overridden on a per-field basis.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators
func (r logFieldSyntaxResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_sensitive_attribute"),
			path.MatchRoot("type"),
			[]string{"attribute-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("excluded_sensitive_field"),
			path.MatchRoot("type"),
			[]string{"json"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_sensitive_attribute"),
			path.MatchRoot("type"),
			[]string{"attribute-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("included_sensitive_field"),
			path.MatchRoot("type"),
			[]string{"json"},
		),
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateLogFieldSyntaxUnknownValues(ctx context.Context, model *logFieldSyntaxResourceModel) {
	if model.IncludedSensitiveAttribute.IsUnknown() || model.IncludedSensitiveAttribute.IsNull() {
		model.IncludedSensitiveAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedSensitiveAttribute.IsUnknown() || model.ExcludedSensitiveAttribute.IsNull() {
		model.ExcludedSensitiveAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludedSensitiveField.IsUnknown() || model.ExcludedSensitiveField.IsNull() {
		model.ExcludedSensitiveField, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludedSensitiveField.IsUnknown() || model.IncludedSensitiveField.IsNull() {
		model.IncludedSensitiveField, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Read a JsonLogFieldSyntaxResponse object into the model struct
func readJsonLogFieldSyntaxResponse(ctx context.Context, r *client.JsonLogFieldSyntaxResponse, state *logFieldSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IncludedSensitiveField = internaltypes.GetStringSet(r.IncludedSensitiveField)
	state.ExcludedSensitiveField = internaltypes.GetStringSet(r.ExcludedSensitiveField)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFieldSyntaxUnknownValues(ctx, state)
}

// Read a AttributeBasedLogFieldSyntaxResponse object into the model struct
func readAttributeBasedLogFieldSyntaxResponse(ctx context.Context, r *client.AttributeBasedLogFieldSyntaxResponse, state *logFieldSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IncludedSensitiveAttribute = internaltypes.GetStringSet(r.IncludedSensitiveAttribute)
	state.ExcludedSensitiveAttribute = internaltypes.GetStringSet(r.ExcludedSensitiveAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFieldSyntaxUnknownValues(ctx, state)
}

// Read a GenericLogFieldSyntaxResponse object into the model struct
func readGenericLogFieldSyntaxResponse(ctx context.Context, r *client.GenericLogFieldSyntaxResponse, state *logFieldSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogFieldSyntaxUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createLogFieldSyntaxOperations(plan logFieldSyntaxResourceModel, state logFieldSyntaxResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedSensitiveAttribute, state.IncludedSensitiveAttribute, "included-sensitive-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedSensitiveAttribute, state.ExcludedSensitiveAttribute, "excluded-sensitive-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedSensitiveField, state.IncludedSensitiveField, "included-sensitive-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedSensitiveField, state.ExcludedSensitiveField, "excluded-sensitive-field")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultBehavior, state.DefaultBehavior, "default-behavior")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *logFieldSyntaxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFieldSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldSyntaxApi.GetLogFieldSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logFieldSyntaxResourceModel
	if readResponse.JsonLogFieldSyntaxResponse != nil {
		readJsonLogFieldSyntaxResponse(ctx, readResponse.JsonLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeBasedLogFieldSyntaxResponse != nil {
		readAttributeBasedLogFieldSyntaxResponse(ctx, readResponse.AttributeBasedLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericLogFieldSyntaxResponse != nil {
		readGenericLogFieldSyntaxResponse(ctx, readResponse.GenericLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFieldSyntaxApi.UpdateLogFieldSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogFieldSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFieldSyntaxApi.UpdateLogFieldSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.JsonLogFieldSyntaxResponse != nil {
			readJsonLogFieldSyntaxResponse(ctx, updateResponse.JsonLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.AttributeBasedLogFieldSyntaxResponse != nil {
			readAttributeBasedLogFieldSyntaxResponse(ctx, updateResponse.AttributeBasedLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericLogFieldSyntaxResponse != nil {
			readGenericLogFieldSyntaxResponse(ctx, updateResponse.GenericLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
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
func (r *logFieldSyntaxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state logFieldSyntaxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldSyntaxApi.GetLogFieldSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.JsonLogFieldSyntaxResponse != nil {
		readJsonLogFieldSyntaxResponse(ctx, readResponse.JsonLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeBasedLogFieldSyntaxResponse != nil {
		readAttributeBasedLogFieldSyntaxResponse(ctx, readResponse.AttributeBasedLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericLogFieldSyntaxResponse != nil {
		readGenericLogFieldSyntaxResponse(ctx, readResponse.GenericLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logFieldSyntaxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan logFieldSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logFieldSyntaxResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LogFieldSyntaxApi.UpdateLogFieldSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogFieldSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFieldSyntaxApi.UpdateLogFieldSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.JsonLogFieldSyntaxResponse != nil {
			readJsonLogFieldSyntaxResponse(ctx, updateResponse.JsonLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.AttributeBasedLogFieldSyntaxResponse != nil {
			readAttributeBasedLogFieldSyntaxResponse(ctx, updateResponse.AttributeBasedLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericLogFieldSyntaxResponse != nil {
			readGenericLogFieldSyntaxResponse(ctx, updateResponse.GenericLogFieldSyntaxResponse, &state, &resp.Diagnostics)
		}
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
func (r *logFieldSyntaxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logFieldSyntaxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
