package sensitiveattribute

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &sensitiveAttributeResource{}
	_ resource.ResourceWithConfigure   = &sensitiveAttributeResource{}
	_ resource.ResourceWithImportState = &sensitiveAttributeResource{}
	_ resource.Resource                = &defaultSensitiveAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultSensitiveAttributeResource{}
	_ resource.ResourceWithImportState = &defaultSensitiveAttributeResource{}
)

// Create a Sensitive Attribute resource
func NewSensitiveAttributeResource() resource.Resource {
	return &sensitiveAttributeResource{}
}

func NewDefaultSensitiveAttributeResource() resource.Resource {
	return &defaultSensitiveAttributeResource{}
}

// sensitiveAttributeResource is the resource implementation.
type sensitiveAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSensitiveAttributeResource is the resource implementation.
type defaultSensitiveAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *sensitiveAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_attribute"
}

func (r *defaultSensitiveAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_sensitive_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *sensitiveAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultSensitiveAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type sensitiveAttributeResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	Description                                  types.String `tfsdk:"description"`
	AttributeType                                types.Set    `tfsdk:"attribute_type"`
	IncludeDefaultSensitiveOperationalAttributes types.Bool   `tfsdk:"include_default_sensitive_operational_attributes"`
	AllowInReturnedEntries                       types.String `tfsdk:"allow_in_returned_entries"`
	AllowInFilter                                types.String `tfsdk:"allow_in_filter"`
	AllowInAdd                                   types.String `tfsdk:"allow_in_add"`
	AllowInCompare                               types.String `tfsdk:"allow_in_compare"`
	AllowInModify                                types.String `tfsdk:"allow_in_modify"`
}

// GetSchema defines the schema for the resource.
func (r *sensitiveAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sensitiveAttributeSchema(ctx, req, resp, false)
}

func (r *defaultSensitiveAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	sensitiveAttributeSchema(ctx, req, resp, true)
}

func sensitiveAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Sensitive Attribute.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Sensitive Attribute",
				Optional:    true,
			},
			"attribute_type": schema.SetAttribute{
				Description: "The name(s) or OID(s) of the attribute types for attributes whose values may be considered sensitive.",
				Required:    true,
				ElementType: types.StringType,
			},
			"include_default_sensitive_operational_attributes": schema.BoolAttribute{
				Description: "Indicates whether to automatically include any server-generated operational attributes that may contain sensitive data.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_in_returned_entries": schema.StringAttribute{
				Description: "Indicates whether sensitive attributes should be included in entries returned to the client. This includes not only search result entries, but also other forms including in the values of controls like the pre-read, post-read, get authorization entry, and LDAP join response controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_in_filter": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in search filters. This also includes filters that may be used in other forms, including assertion and LDAP join request controls.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_in_add": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in add requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_in_compare": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with compare requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_in_modify": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with modify requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for sensitive-attribute sensitive-attribute
func addOptionalSensitiveAttributeFields(ctx context.Context, addRequest *client.AddSensitiveAttributeRequest, plan sensitiveAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeDefaultSensitiveOperationalAttributes) {
		addRequest.IncludeDefaultSensitiveOperationalAttributes = plan.IncludeDefaultSensitiveOperationalAttributes.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInReturnedEntries) {
		allowInReturnedEntries, err := client.NewEnumsensitiveAttributeAllowInReturnedEntriesPropFromValue(plan.AllowInReturnedEntries.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInReturnedEntries = allowInReturnedEntries
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInFilter) {
		allowInFilter, err := client.NewEnumsensitiveAttributeAllowInFilterPropFromValue(plan.AllowInFilter.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInFilter = allowInFilter
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInAdd) {
		allowInAdd, err := client.NewEnumsensitiveAttributeAllowInAddPropFromValue(plan.AllowInAdd.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInAdd = allowInAdd
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInCompare) {
		allowInCompare, err := client.NewEnumsensitiveAttributeAllowInComparePropFromValue(plan.AllowInCompare.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInCompare = allowInCompare
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowInModify) {
		allowInModify, err := client.NewEnumsensitiveAttributeAllowInModifyPropFromValue(plan.AllowInModify.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowInModify = allowInModify
	}
	return nil
}

// Read a SensitiveAttributeResponse object into the model struct
func readSensitiveAttributeResponse(ctx context.Context, r *client.SensitiveAttributeResponse, state *sensitiveAttributeResourceModel, expectedValues *sensitiveAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.IncludeDefaultSensitiveOperationalAttributes = internaltypes.BoolTypeOrNil(r.IncludeDefaultSensitiveOperationalAttributes)
	state.AllowInReturnedEntries = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInReturnedEntriesProp(r.AllowInReturnedEntries), internaltypes.IsEmptyString(expectedValues.AllowInReturnedEntries))
	state.AllowInFilter = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInFilterProp(r.AllowInFilter), internaltypes.IsEmptyString(expectedValues.AllowInFilter))
	state.AllowInAdd = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInAddProp(r.AllowInAdd), internaltypes.IsEmptyString(expectedValues.AllowInAdd))
	state.AllowInCompare = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInCompareProp(r.AllowInCompare), internaltypes.IsEmptyString(expectedValues.AllowInCompare))
	state.AllowInModify = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInModifyProp(r.AllowInModify), internaltypes.IsEmptyString(expectedValues.AllowInModify))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSensitiveAttributeOperations(plan sensitiveAttributeResourceModel, state sensitiveAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeDefaultSensitiveOperationalAttributes, state.IncludeDefaultSensitiveOperationalAttributes, "include-default-sensitive-operational-attributes")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInReturnedEntries, state.AllowInReturnedEntries, "allow-in-returned-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInFilter, state.AllowInFilter, "allow-in-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInAdd, state.AllowInAdd, "allow-in-add")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInCompare, state.AllowInCompare, "allow-in-compare")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowInModify, state.AllowInModify, "allow-in-modify")
	return ops
}

// Create a sensitive-attribute sensitive-attribute
func (r *sensitiveAttributeResource) CreateSensitiveAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan sensitiveAttributeResourceModel) (*sensitiveAttributeResourceModel, error) {
	var AttributeTypeSlice []string
	plan.AttributeType.ElementsAs(ctx, &AttributeTypeSlice, false)
	addRequest := client.NewAddSensitiveAttributeRequest(plan.Id.ValueString(),
		AttributeTypeSlice)
	err := addOptionalSensitiveAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Sensitive Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SensitiveAttributeApi.AddSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSensitiveAttributeRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.SensitiveAttributeApi.AddSensitiveAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Sensitive Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state sensitiveAttributeResourceModel
	readSensitiveAttributeResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *sensitiveAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateSensitiveAttribute(ctx, req, resp, plan)
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
func (r *defaultSensitiveAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SensitiveAttributeApi.GetSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state sensitiveAttributeResourceModel
	readSensitiveAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SensitiveAttributeApi.UpdateSensitiveAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSensitiveAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SensitiveAttributeApi.UpdateSensitiveAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sensitive Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSensitiveAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *sensitiveAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSensitiveAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSensitiveAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state sensitiveAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SensitiveAttributeApi.GetSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSensitiveAttributeResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *sensitiveAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSensitiveAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSensitiveAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSensitiveAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan sensitiveAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state sensitiveAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SensitiveAttributeApi.UpdateSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSensitiveAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SensitiveAttributeApi.UpdateSensitiveAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Sensitive Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSensitiveAttributeResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSensitiveAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *sensitiveAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sensitiveAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SensitiveAttributeApi.DeleteSensitiveAttributeExecute(r.apiClient.SensitiveAttributeApi.DeleteSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Sensitive Attribute", err, httpResp)
		return
	}
}

func (r *sensitiveAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSensitiveAttribute(ctx, req, resp)
}

func (r *defaultSensitiveAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSensitiveAttribute(ctx, req, resp)
}

func importSensitiveAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
