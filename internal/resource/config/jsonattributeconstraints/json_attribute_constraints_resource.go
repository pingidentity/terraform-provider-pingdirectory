package jsonattributeconstraints

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &jsonAttributeConstraintsResource{}
	_ resource.ResourceWithConfigure   = &jsonAttributeConstraintsResource{}
	_ resource.ResourceWithImportState = &jsonAttributeConstraintsResource{}
	_ resource.Resource                = &defaultJsonAttributeConstraintsResource{}
	_ resource.ResourceWithConfigure   = &defaultJsonAttributeConstraintsResource{}
	_ resource.ResourceWithImportState = &defaultJsonAttributeConstraintsResource{}
)

// Create a Json Attribute Constraints resource
func NewJsonAttributeConstraintsResource() resource.Resource {
	return &jsonAttributeConstraintsResource{}
}

func NewDefaultJsonAttributeConstraintsResource() resource.Resource {
	return &defaultJsonAttributeConstraintsResource{}
}

// jsonAttributeConstraintsResource is the resource implementation.
type jsonAttributeConstraintsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJsonAttributeConstraintsResource is the resource implementation.
type defaultJsonAttributeConstraintsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jsonAttributeConstraintsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_json_attribute_constraints"
}

func (r *defaultJsonAttributeConstraintsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_json_attribute_constraints"
}

// Configure adds the provider configured client to the resource.
func (r *jsonAttributeConstraintsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultJsonAttributeConstraintsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type jsonAttributeConstraintsResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	Type               types.String `tfsdk:"type"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	AttributeType      types.String `tfsdk:"attribute_type"`
	AllowUnnamedFields types.Bool   `tfsdk:"allow_unnamed_fields"`
}

// GetSchema defines the schema for the resource.
func (r *jsonAttributeConstraintsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonAttributeConstraintsSchema(ctx, req, resp, false)
}

func (r *defaultJsonAttributeConstraintsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonAttributeConstraintsSchema(ctx, req, resp, true)
}

func jsonAttributeConstraintsSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Json Attribute Constraints.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of JSON Attribute Constraints resource. Options are ['json-attribute-constraints']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("json-attribute-constraints"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"json-attribute-constraints"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this JSON Attribute Constraints",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this JSON Attribute Constraints is enabled.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"attribute_type": schema.StringAttribute{
				Description: "The name or OID of the LDAP attribute type whose values will be subject to the associated field constraints. This attribute type must be defined in the server schema, and it must have a \"JSON object\" syntax.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allow_unnamed_fields": schema.BoolAttribute{
				Description: "Indicates whether JSON objects stored as values of attributes with the associated attribute-type will be permitted to include fields for which there is no subordinate json-field-constraints definition. If unnamed fields are allowed, then no constraints will be imposed on the values of those fields. However, if unnamed fields are not allowed, then the server will reject any attempt to store a JSON object with a field for which there is no corresponding json-fields-constraints definition.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "attribute_type"})
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for json-attribute-constraints json-attribute-constraints
func addOptionalJsonAttributeConstraintsFields(ctx context.Context, addRequest *client.AddJsonAttributeConstraintsRequest, plan jsonAttributeConstraintsResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowUnnamedFields) {
		addRequest.AllowUnnamedFields = plan.AllowUnnamedFields.ValueBoolPointer()
	}
}

// Read a JsonAttributeConstraintsResponse object into the model struct
func readJsonAttributeConstraintsResponse(ctx context.Context, r *client.JsonAttributeConstraintsResponse, state *jsonAttributeConstraintsResourceModel, expectedValues *jsonAttributeConstraintsResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-attribute-constraints")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.AllowUnnamedFields = internaltypes.BoolTypeOrNil(r.AllowUnnamedFields)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createJsonAttributeConstraintsOperations(plan jsonAttributeConstraintsResourceModel, state jsonAttributeConstraintsResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUnnamedFields, state.AllowUnnamedFields, "allow-unnamed-fields")
	return ops
}

// Create a json-attribute-constraints json-attribute-constraints
func (r *jsonAttributeConstraintsResource) CreateJsonAttributeConstraints(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan jsonAttributeConstraintsResourceModel) (*jsonAttributeConstraintsResourceModel, error) {
	addRequest := client.NewAddJsonAttributeConstraintsRequest(plan.AttributeType.ValueString())
	addOptionalJsonAttributeConstraintsFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.JsonAttributeConstraintsApi.AddJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddJsonAttributeConstraintsRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.JsonAttributeConstraintsApi.AddJsonAttributeConstraintsExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Json Attribute Constraints", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jsonAttributeConstraintsResourceModel
	readJsonAttributeConstraintsResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *jsonAttributeConstraintsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonAttributeConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateJsonAttributeConstraints(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultJsonAttributeConstraintsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonAttributeConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.JsonAttributeConstraintsApi.GetJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Attribute Constraints", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jsonAttributeConstraintsResourceModel
	readJsonAttributeConstraintsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.JsonAttributeConstraintsApi.UpdateJsonAttributeConstraints(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString())
	ops := createJsonAttributeConstraintsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.JsonAttributeConstraintsApi.UpdateJsonAttributeConstraintsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Attribute Constraints", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonAttributeConstraintsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jsonAttributeConstraintsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonAttributeConstraints(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultJsonAttributeConstraintsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonAttributeConstraints(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readJsonAttributeConstraints(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state jsonAttributeConstraintsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.JsonAttributeConstraintsApi.GetJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.AttributeType.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Json Attribute Constraints", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Attribute Constraints", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJsonAttributeConstraintsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *jsonAttributeConstraintsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonAttributeConstraints(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJsonAttributeConstraintsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonAttributeConstraints(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJsonAttributeConstraints(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jsonAttributeConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jsonAttributeConstraintsResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.JsonAttributeConstraintsApi.UpdateJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.AttributeType.ValueString())

	// Determine what update operations are necessary
	ops := createJsonAttributeConstraintsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.JsonAttributeConstraintsApi.UpdateJsonAttributeConstraintsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Attribute Constraints", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonAttributeConstraintsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJsonAttributeConstraintsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jsonAttributeConstraintsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jsonAttributeConstraintsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.JsonAttributeConstraintsApi.DeleteJsonAttributeConstraintsExecute(r.apiClient.JsonAttributeConstraintsApi.DeleteJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AttributeType.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Json Attribute Constraints", err, httpResp)
		return
	}
}

func (r *jsonAttributeConstraintsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonAttributeConstraints(ctx, req, resp)
}

func (r *defaultJsonAttributeConstraintsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonAttributeConstraints(ctx, req, resp)
}

func importJsonAttributeConstraints(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to attribute_type attribute
	resource.ImportStatePassthroughID(ctx, path.Root("attribute_type"), req, resp)
}
