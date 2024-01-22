package localdbvlvindex

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localDbVlvIndexResource{}
	_ resource.ResourceWithConfigure   = &localDbVlvIndexResource{}
	_ resource.ResourceWithImportState = &localDbVlvIndexResource{}
	_ resource.Resource                = &defaultLocalDbVlvIndexResource{}
	_ resource.ResourceWithConfigure   = &defaultLocalDbVlvIndexResource{}
	_ resource.ResourceWithImportState = &defaultLocalDbVlvIndexResource{}
)

// Create a Local Db Vlv Index resource
func NewLocalDbVlvIndexResource() resource.Resource {
	return &localDbVlvIndexResource{}
}

func NewDefaultLocalDbVlvIndexResource() resource.Resource {
	return &defaultLocalDbVlvIndexResource{}
}

// localDbVlvIndexResource is the resource implementation.
type localDbVlvIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLocalDbVlvIndexResource is the resource implementation.
type defaultLocalDbVlvIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *localDbVlvIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_vlv_index"
}

func (r *defaultLocalDbVlvIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_local_db_vlv_index"
}

// Configure adds the provider configured client to the resource.
func (r *localDbVlvIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLocalDbVlvIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type localDbVlvIndexResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	BackendName     types.String `tfsdk:"backend_name"`
	BaseDN          types.String `tfsdk:"base_dn"`
	Scope           types.String `tfsdk:"scope"`
	Filter          types.String `tfsdk:"filter"`
	SortOrder       types.String `tfsdk:"sort_order"`
	Name            types.String `tfsdk:"name"`
	MaxBlockSize    types.Int64  `tfsdk:"max_block_size"`
	CacheMode       types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the resource.
func (r *localDbVlvIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbVlvIndexSchema(ctx, req, resp, false)
}

func (r *defaultLocalDbVlvIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbVlvIndexSchema(ctx, req, resp, true)
}

func localDbVlvIndexSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Local Db Vlv Index.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Local DB VLV Index resource. Options are ['local-db-vlv-index']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("local-db-vlv-index"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"local-db-vlv-index"}...),
				},
			},
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"base_dn": schema.StringAttribute{
				Description: "Specifies the base DN used in the search query that is being indexed.",
				Required:    true,
			},
			"scope": schema.StringAttribute{
				Description: "Specifies the LDAP scope of the query that is being indexed.",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "Specifies the LDAP filter used in the query that is being indexed.",
				Required:    true,
			},
			"sort_order": schema.StringAttribute{
				Description: "Specifies the names of the attributes that are used to sort the entries for the query being indexed.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Specifies a unique name for this VLV index.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_block_size": schema.Int64Attribute{
				Description: "Specifies the number of entry IDs to store in a single sorted set before it must be split.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(4000),
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the database for this index.",
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "name", "backend_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		nameAttr := schemaDef.Attributes["name"].(schema.StringAttribute)
		nameAttr.PlanModifiers = append(nameAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["name"] = nameAttr
		maxBlockSizeAttr := schemaDef.Attributes["max_block_size"].(schema.Int64Attribute)
		maxBlockSizeAttr.PlanModifiers = append(maxBlockSizeAttr.PlanModifiers, int64planmodifier.RequiresReplace())
		schemaDef.Attributes["max_block_size"] = maxBlockSizeAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for local-db-vlv-index local-db-vlv-index
func addOptionalLocalDbVlvIndexFields(ctx context.Context, addRequest *client.AddLocalDbVlvIndexRequest, plan localDbVlvIndexResourceModel) error {
	if internaltypes.IsDefined(plan.MaxBlockSize) {
		addRequest.MaxBlockSize = plan.MaxBlockSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CacheMode) {
		cacheMode, err := client.NewEnumlocalDbVlvIndexCacheModePropFromValue(plan.CacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.CacheMode = cacheMode
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *localDbVlvIndexResourceModel) populateAllComputedStringAttributes() {
	if model.Scope.IsUnknown() || model.Scope.IsNull() {
		model.Scope = types.StringValue("")
	}
	if model.Filter.IsUnknown() || model.Filter.IsNull() {
		model.Filter = types.StringValue("")
	}
	if model.BaseDN.IsUnknown() || model.BaseDN.IsNull() {
		model.BaseDN = types.StringValue("")
	}
	if model.SortOrder.IsUnknown() || model.SortOrder.IsNull() {
		model.SortOrder = types.StringValue("")
	}
	if model.CacheMode.IsUnknown() || model.CacheMode.IsNull() {
		model.CacheMode = types.StringValue("")
	}
	if model.Name.IsUnknown() || model.Name.IsNull() {
		model.Name = types.StringValue("")
	}
}

// Read a LocalDbVlvIndexResponse object into the model struct
func readLocalDbVlvIndexResponse(ctx context.Context, r *client.LocalDbVlvIndexResponse, state *localDbVlvIndexResourceModel, expectedValues *localDbVlvIndexResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db-vlv-index")
	state.Id = types.StringValue(r.Id)
	state.BaseDN = types.StringValue(r.BaseDN)
	state.Scope = types.StringValue(r.Scope.String())
	state.Filter = types.StringValue(r.Filter)
	state.SortOrder = types.StringValue(r.SortOrder)
	state.Name = types.StringValue(r.Name)
	state.MaxBlockSize = internaltypes.Int64TypeOrNil(r.MaxBlockSize)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbVlvIndexCacheModeProp(r.CacheMode), internaltypes.IsEmptyString(expectedValues.CacheMode))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *localDbVlvIndexResourceModel) setStateValuesNotReturnedByAPI(expectedValues *localDbVlvIndexResourceModel) {
	if !expectedValues.BackendName.IsUnknown() {
		state.BackendName = expectedValues.BackendName
	}
}

// Create any update operations necessary to make the state match the plan
func createLocalDbVlvIndexOperations(plan localDbVlvIndexResourceModel, state localDbVlvIndexResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Scope, state.Scope, "scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringOperationIfNecessary(&ops, plan.SortOrder, state.SortOrder, "sort-order")
	operations.AddStringOperationIfNecessary(&ops, plan.Name, state.Name, "name")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxBlockSize, state.MaxBlockSize, "max-block-size")
	operations.AddStringOperationIfNecessary(&ops, plan.CacheMode, state.CacheMode, "cache-mode")
	return ops
}

// Create a local-db-vlv-index local-db-vlv-index
func (r *localDbVlvIndexResource) CreateLocalDbVlvIndex(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan localDbVlvIndexResourceModel) (*localDbVlvIndexResourceModel, error) {
	scope, err := client.NewEnumlocalDbVlvIndexScopePropFromValue(plan.Scope.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for Scope", err.Error())
		return nil, err
	}
	addRequest := client.NewAddLocalDbVlvIndexRequest(plan.BaseDN.ValueString(),
		*scope,
		plan.Filter.ValueString(),
		plan.SortOrder.ValueString(),
		plan.Name.ValueString(),
		plan.Name.ValueString())
	err = addOptionalLocalDbVlvIndexFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Db Vlv Index", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LocalDbVlvIndexAPI.AddLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendName.ValueString())
	apiAddRequest = apiAddRequest.AddLocalDbVlvIndexRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.LocalDbVlvIndexAPI.AddLocalDbVlvIndexExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Local Db Vlv Index", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state localDbVlvIndexResourceModel
	readLocalDbVlvIndexResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *localDbVlvIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbVlvIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateLocalDbVlvIndex(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
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
func (r *defaultLocalDbVlvIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbVlvIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbVlvIndexAPI.GetLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state localDbVlvIndexResourceModel
	readLocalDbVlvIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LocalDbVlvIndexAPI.UpdateLocalDbVlvIndex(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())
	ops := createLocalDbVlvIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LocalDbVlvIndexAPI.UpdateLocalDbVlvIndexExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Vlv Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbVlvIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *localDbVlvIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLocalDbVlvIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLocalDbVlvIndex(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state localDbVlvIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LocalDbVlvIndexAPI.GetLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbVlvIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *localDbVlvIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbVlvIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLocalDbVlvIndex(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan localDbVlvIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state localDbVlvIndexResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LocalDbVlvIndexAPI.UpdateLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())

	// Determine what update operations are necessary
	ops := createLocalDbVlvIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LocalDbVlvIndexAPI.UpdateLocalDbVlvIndexExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Vlv Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbVlvIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLocalDbVlvIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *localDbVlvIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localDbVlvIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LocalDbVlvIndexAPI.DeleteLocalDbVlvIndexExecute(r.apiClient.LocalDbVlvIndexAPI.DeleteLocalDbVlvIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.BackendName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Local Db Vlv Index", err, httpResp)
		return
	}
}

func (r *localDbVlvIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbVlvIndex(ctx, req, resp)
}

func (r *defaultLocalDbVlvIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbVlvIndex(ctx, req, resp)
}

func importLocalDbVlvIndex(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [backend-name]/[local-db-vlv-index-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("backend_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
