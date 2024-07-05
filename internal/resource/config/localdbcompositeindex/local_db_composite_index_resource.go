package localdbcompositeindex

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localDbCompositeIndexResource{}
	_ resource.ResourceWithConfigure   = &localDbCompositeIndexResource{}
	_ resource.ResourceWithImportState = &localDbCompositeIndexResource{}
	_ resource.Resource                = &defaultLocalDbCompositeIndexResource{}
	_ resource.ResourceWithConfigure   = &defaultLocalDbCompositeIndexResource{}
	_ resource.ResourceWithImportState = &defaultLocalDbCompositeIndexResource{}
)

// Create a Local Db Composite Index resource
func NewLocalDbCompositeIndexResource() resource.Resource {
	return &localDbCompositeIndexResource{}
}

func NewDefaultLocalDbCompositeIndexResource() resource.Resource {
	return &defaultLocalDbCompositeIndexResource{}
}

// localDbCompositeIndexResource is the resource implementation.
type localDbCompositeIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLocalDbCompositeIndexResource is the resource implementation.
type defaultLocalDbCompositeIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *localDbCompositeIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_composite_index"
}

func (r *defaultLocalDbCompositeIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_local_db_composite_index"
}

// Configure adds the provider configured client to the resource.
func (r *localDbCompositeIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLocalDbCompositeIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type localDbCompositeIndexResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	Notifications          types.Set    `tfsdk:"notifications"`
	RequiredActions        types.Set    `tfsdk:"required_actions"`
	Type                   types.String `tfsdk:"type"`
	BackendName            types.String `tfsdk:"backend_name"`
	Description            types.String `tfsdk:"description"`
	IndexFilterPattern     types.String `tfsdk:"index_filter_pattern"`
	IndexBaseDNPattern     types.String `tfsdk:"index_base_dn_pattern"`
	IndexEntryLimit        types.Int64  `tfsdk:"index_entry_limit"`
	PrimeIndex             types.Bool   `tfsdk:"prime_index"`
	PrimeInternalNodesOnly types.Bool   `tfsdk:"prime_internal_nodes_only"`
	CacheMode              types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the resource.
func (r *localDbCompositeIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbCompositeIndexSchema(ctx, req, resp, false)
}

func (r *defaultLocalDbCompositeIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbCompositeIndexSchema(ctx, req, resp, true)
}

func localDbCompositeIndexSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Local Db Composite Index.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Local DB Composite Index resource. Options are ['local-db-composite-index']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("local-db-composite-index"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"local-db-composite-index"}...),
				},
			},
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Local DB Composite Index",
				Optional:    true,
			},
			"index_filter_pattern": schema.StringAttribute{
				Description: "A filter pattern that identifies which entries to include in the index.",
				Required:    true,
			},
			"index_base_dn_pattern": schema.StringAttribute{
				Description: "An optional base DN pattern that identifies portions of the DIT in which entries to index may exist.",
				Optional:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that any single index key will be allowed to match before the server stops maintaining the ID set for that index key.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"prime_index": schema.BoolAttribute{
				Description: "Indicates whether the server should load the contents of this index into memory when the backend is being opened.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"prime_internal_nodes_only": schema.BoolAttribute{
				Description: "Indicates whether to only prime the internal nodes of the index database, rather than priming both internal and leaf nodes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"cache_mode": schema.StringAttribute{
				Description: "The behavior that the server should exhibit when storing information from this index in the database cache.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"cache-keys-and-values", "cache-keys-only", "no-caching"}...),
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "backend_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		indexFilterPatternAttr := schemaDef.Attributes["index_filter_pattern"].(schema.StringAttribute)
		indexFilterPatternAttr.PlanModifiers = append(indexFilterPatternAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["index_filter_pattern"] = indexFilterPatternAttr
		indexBaseDnPatternAttr := schemaDef.Attributes["index_base_dn_pattern"].(schema.StringAttribute)
		indexBaseDnPatternAttr.PlanModifiers = append(indexBaseDnPatternAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["index_base_dn_pattern"] = indexBaseDnPatternAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for local-db-composite-index local-db-composite-index
func addOptionalLocalDbCompositeIndexFields(ctx context.Context, addRequest *client.AddLocalDbCompositeIndexRequest, plan localDbCompositeIndexResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IndexBaseDNPattern) {
		addRequest.IndexBaseDNPattern = plan.IndexBaseDNPattern.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IndexEntryLimit) {
		addRequest.IndexEntryLimit = plan.IndexEntryLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.PrimeIndex) {
		addRequest.PrimeIndex = plan.PrimeIndex.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.PrimeInternalNodesOnly) {
		addRequest.PrimeInternalNodesOnly = plan.PrimeInternalNodesOnly.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CacheMode) {
		cacheMode, err := client.NewEnumlocalDbCompositeIndexCacheModePropFromValue(plan.CacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.CacheMode = cacheMode
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *localDbCompositeIndexResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.IndexBaseDNPattern.IsUnknown() || model.IndexBaseDNPattern.IsNull() {
		model.IndexBaseDNPattern = types.StringValue("")
	}
	if model.IndexFilterPattern.IsUnknown() || model.IndexFilterPattern.IsNull() {
		model.IndexFilterPattern = types.StringValue("")
	}
	if model.CacheMode.IsUnknown() || model.CacheMode.IsNull() {
		model.CacheMode = types.StringValue("")
	}
}

// Read a LocalDbCompositeIndexResponse object into the model struct
func readLocalDbCompositeIndexResponse(ctx context.Context, r *client.LocalDbCompositeIndexResponse, state *localDbCompositeIndexResourceModel, expectedValues *localDbCompositeIndexResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("local-db-composite-index")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.IndexFilterPattern = types.StringValue(r.IndexFilterPattern)
	state.IndexBaseDNPattern = internaltypes.StringTypeOrNil(r.IndexBaseDNPattern, internaltypes.IsEmptyString(expectedValues.IndexBaseDNPattern))
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.PrimeInternalNodesOnly = internaltypes.BoolTypeOrNil(r.PrimeInternalNodesOnly)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbCompositeIndexCacheModeProp(r.CacheMode), internaltypes.IsEmptyString(expectedValues.CacheMode))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *localDbCompositeIndexResourceModel) setStateValuesNotReturnedByAPI(expectedValues *localDbCompositeIndexResourceModel) {
	if !expectedValues.BackendName.IsUnknown() {
		state.BackendName = expectedValues.BackendName
	}
}

// Create any update operations necessary to make the state match the plan
func createLocalDbCompositeIndexOperations(plan localDbCompositeIndexResourceModel, state localDbCompositeIndexResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.IndexFilterPattern, state.IndexFilterPattern, "index-filter-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.IndexBaseDNPattern, state.IndexBaseDNPattern, "index-base-dn-pattern")
	operations.AddInt64OperationIfNecessary(&ops, plan.IndexEntryLimit, state.IndexEntryLimit, "index-entry-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeIndex, state.PrimeIndex, "prime-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeInternalNodesOnly, state.PrimeInternalNodesOnly, "prime-internal-nodes-only")
	operations.AddStringOperationIfNecessary(&ops, plan.CacheMode, state.CacheMode, "cache-mode")
	return ops
}

// Create a local-db-composite-index local-db-composite-index
func (r *localDbCompositeIndexResource) CreateLocalDbCompositeIndex(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan localDbCompositeIndexResourceModel) (*localDbCompositeIndexResourceModel, error) {
	addRequest := client.NewAddLocalDbCompositeIndexRequest(plan.IndexFilterPattern.ValueString(),
		plan.Name.ValueString())
	err := addOptionalLocalDbCompositeIndexFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Db Composite Index", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LocalDbCompositeIndexAPI.AddLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendName.ValueString())
	apiAddRequest = apiAddRequest.AddLocalDbCompositeIndexRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.LocalDbCompositeIndexAPI.AddLocalDbCompositeIndexExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Local Db Composite Index", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state localDbCompositeIndexResourceModel
	readLocalDbCompositeIndexResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *localDbCompositeIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbCompositeIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateLocalDbCompositeIndex(ctx, req, resp, plan)
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
func (r *defaultLocalDbCompositeIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbCompositeIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbCompositeIndexAPI.GetLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Composite Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state localDbCompositeIndexResourceModel
	readLocalDbCompositeIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LocalDbCompositeIndexAPI.UpdateLocalDbCompositeIndex(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())
	ops := createLocalDbCompositeIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LocalDbCompositeIndexAPI.UpdateLocalDbCompositeIndexExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Composite Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbCompositeIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *localDbCompositeIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbCompositeIndex(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLocalDbCompositeIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbCompositeIndex(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLocalDbCompositeIndex(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state localDbCompositeIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LocalDbCompositeIndexAPI.GetLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Composite Index", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Composite Index", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbCompositeIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *localDbCompositeIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbCompositeIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbCompositeIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbCompositeIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLocalDbCompositeIndex(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan localDbCompositeIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state localDbCompositeIndexResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LocalDbCompositeIndexAPI.UpdateLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())

	// Determine what update operations are necessary
	ops := createLocalDbCompositeIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LocalDbCompositeIndexAPI.UpdateLocalDbCompositeIndexExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Composite Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbCompositeIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLocalDbCompositeIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *localDbCompositeIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localDbCompositeIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LocalDbCompositeIndexAPI.DeleteLocalDbCompositeIndexExecute(r.apiClient.LocalDbCompositeIndexAPI.DeleteLocalDbCompositeIndex(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.BackendName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Local Db Composite Index", err, httpResp)
		return
	}
}

func (r *localDbCompositeIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbCompositeIndex(ctx, req, resp)
}

func (r *defaultLocalDbCompositeIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbCompositeIndex(ctx, req, resp)
}

func importLocalDbCompositeIndex(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [backend-name]/[local-db-composite-index-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("backend_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
