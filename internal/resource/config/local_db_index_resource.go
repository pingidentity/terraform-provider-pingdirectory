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
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &localDbIndexResource{}
	_ resource.ResourceWithConfigure   = &localDbIndexResource{}
	_ resource.ResourceWithImportState = &localDbIndexResource{}
	_ resource.Resource                = &defaultLocalDbIndexResource{}
	_ resource.ResourceWithConfigure   = &defaultLocalDbIndexResource{}
	_ resource.ResourceWithImportState = &defaultLocalDbIndexResource{}
)

// Create a Local Db Index resource
func NewLocalDbIndexResource() resource.Resource {
	return &localDbIndexResource{}
}

func NewDefaultLocalDbIndexResource() resource.Resource {
	return &defaultLocalDbIndexResource{}
}

// localDbIndexResource is the resource implementation.
type localDbIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLocalDbIndexResource is the resource implementation.
type defaultLocalDbIndexResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *localDbIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_db_index"
}

func (r *defaultLocalDbIndexResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_local_db_index"
}

// Configure adds the provider configured client to the resource.
func (r *localDbIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLocalDbIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type localDbIndexResourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	LastUpdated                                  types.String `tfsdk:"last_updated"`
	Notifications                                types.Set    `tfsdk:"notifications"`
	RequiredActions                              types.Set    `tfsdk:"required_actions"`
	BackendName                                  types.String `tfsdk:"backend_name"`
	Attribute                                    types.String `tfsdk:"attribute"`
	IndexEntryLimit                              types.Int64  `tfsdk:"index_entry_limit"`
	SubstringIndexEntryLimit                     types.Int64  `tfsdk:"substring_index_entry_limit"`
	MaintainMatchCountForKeysExceedingEntryLimit types.Bool   `tfsdk:"maintain_match_count_for_keys_exceeding_entry_limit"`
	IndexType                                    types.Set    `tfsdk:"index_type"`
	SubstringLength                              types.Int64  `tfsdk:"substring_length"`
	PrimeIndex                                   types.Bool   `tfsdk:"prime_index"`
	PrimeInternalNodesOnly                       types.Bool   `tfsdk:"prime_internal_nodes_only"`
	EqualityIndexFilter                          types.Set    `tfsdk:"equality_index_filter"`
	MaintainEqualityIndexWithoutFilter           types.Bool   `tfsdk:"maintain_equality_index_without_filter"`
	CacheMode                                    types.String `tfsdk:"cache_mode"`
}

// GetSchema defines the schema for the resource.
func (r *localDbIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbIndexSchema(ctx, req, resp, false)
}

func (r *defaultLocalDbIndexResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	localDbIndexSchema(ctx, req, resp, true)
}

func localDbIndexSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Local Db Index.",
		Attributes: map[string]schema.Attribute{
			"backend_name": schema.StringAttribute{
				Description: "Name of the parent Backend",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"attribute": schema.StringAttribute{
				Description: "Specifies the name of the attribute for which the index is to be maintained.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained.",
				Optional:    true,
				Computed:    true,
			},
			"substring_index_entry_limit": schema.Int64Attribute{
				Description: "Specifies, for substring indexes, the maximum number of entries that are allowed to match a given index key before that particular index key is no longer maintained. Setting a large limit can dramatically increase the database size on disk and have a big impact on server performance if the indexed attribute is modified frequently. When a very large limit is required, creating a dedicated composite index with an index-filter-pattern of (attr=*?*) will give the best balance between search and update performance.",
				Optional:    true,
				Computed:    true,
			},
			"maintain_match_count_for_keys_exceeding_entry_limit": schema.BoolAttribute{
				Description: "Indicates whether to continue to maintain a count of the number of matching entries for an index key even after that count exceeds the index entry limit.",
				Optional:    true,
				Computed:    true,
			},
			"index_type": schema.SetAttribute{
				Description: "Specifies the type(s) of indexing that should be performed for the associated attribute.",
				Required:    true,
				ElementType: types.StringType,
			},
			"substring_length": schema.Int64Attribute{
				Description: "The length of substrings in a substring index.",
				Optional:    true,
				Computed:    true,
			},
			"prime_index": schema.BoolAttribute{
				Description: "If this option is enabled and this index's backend is configured to prime indexes, then this index will be loaded at startup.",
				Optional:    true,
				Computed:    true,
			},
			"prime_internal_nodes_only": schema.BoolAttribute{
				Description: "If this option is enabled and this index's backend is configured to prime indexes using the preload method, then only the internal database nodes (i.e., the database keys but not values) should be primed when the backend is initialized.",
				Optional:    true,
				Computed:    true,
			},
			"equality_index_filter": schema.SetAttribute{
				Description: "A search filter that may be used in conjunction with an equality component for the associated attribute type. If an equality index filter is defined, then an additional equality index will be maintained for the associated attribute, but only for entries which match the provided filter. Further, the index will be used only for searches containing an equality component with the associated attribute type ANDed with this filter.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maintain_equality_index_without_filter": schema.BoolAttribute{
				Description: "Indicates whether to maintain a separate equality index for the associated attribute without any filter, in addition to maintaining an index for each equality index filter that is defined. If this is false, then the attribute will not be indexed for equality by itself but only in conjunction with the defined equality index filters.",
				Optional:    true,
				Computed:    true,
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the database for this index. This controls how much database cache memory can be consumed by this index.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	AddCommonSchema(&schema, false)
	if setOptionalToComputed {
		SetAllAttributesToOptionalAndComputed(&schema, []string{"attribute", "backend_name"})
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalLocalDbIndexFields(ctx context.Context, addRequest *client.AddLocalDbIndexRequest, plan localDbIndexResourceModel) error {
	if internaltypes.IsDefined(plan.IndexEntryLimit) {
		intVal := int32(plan.IndexEntryLimit.ValueInt64())
		addRequest.IndexEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.SubstringIndexEntryLimit) {
		intVal := int32(plan.SubstringIndexEntryLimit.ValueInt64())
		addRequest.SubstringIndexEntryLimit = &intVal
	}
	if internaltypes.IsDefined(plan.MaintainMatchCountForKeysExceedingEntryLimit) {
		boolVal := plan.MaintainMatchCountForKeysExceedingEntryLimit.ValueBool()
		addRequest.MaintainMatchCountForKeysExceedingEntryLimit = &boolVal
	}
	if internaltypes.IsDefined(plan.SubstringLength) {
		intVal := int32(plan.SubstringLength.ValueInt64())
		addRequest.SubstringLength = &intVal
	}
	if internaltypes.IsDefined(plan.PrimeIndex) {
		boolVal := plan.PrimeIndex.ValueBool()
		addRequest.PrimeIndex = &boolVal
	}
	if internaltypes.IsDefined(plan.PrimeInternalNodesOnly) {
		boolVal := plan.PrimeInternalNodesOnly.ValueBool()
		addRequest.PrimeInternalNodesOnly = &boolVal
	}
	if internaltypes.IsDefined(plan.EqualityIndexFilter) {
		var slice []string
		plan.EqualityIndexFilter.ElementsAs(ctx, &slice, false)
		addRequest.EqualityIndexFilter = slice
	}
	if internaltypes.IsDefined(plan.MaintainEqualityIndexWithoutFilter) {
		boolVal := plan.MaintainEqualityIndexWithoutFilter.ValueBool()
		addRequest.MaintainEqualityIndexWithoutFilter = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CacheMode) {
		cacheMode, err := client.NewEnumlocalDbIndexCacheModePropFromValue(plan.CacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.CacheMode = cacheMode
	}
	return nil
}

// Read a LocalDbIndexResponse object into the model struct
func readLocalDbIndexResponse(ctx context.Context, r *client.LocalDbIndexResponse, state *localDbIndexResourceModel, expectedValues *localDbIndexResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BackendName = expectedValues.BackendName
	state.Attribute = types.StringValue(r.Attribute)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.SubstringIndexEntryLimit = internaltypes.Int64TypeOrNil(r.SubstringIndexEntryLimit)
	state.MaintainMatchCountForKeysExceedingEntryLimit = internaltypes.BoolTypeOrNil(r.MaintainMatchCountForKeysExceedingEntryLimit)
	state.IndexType = internaltypes.GetStringSet(
		client.StringSliceEnumlocalDbIndexIndexTypeProp(r.IndexType))
	state.SubstringLength = internaltypes.Int64TypeOrNil(r.SubstringLength)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.PrimeInternalNodesOnly = internaltypes.BoolTypeOrNil(r.PrimeInternalNodesOnly)
	state.EqualityIndexFilter = internaltypes.GetStringSet(r.EqualityIndexFilter)
	state.MaintainEqualityIndexWithoutFilter = internaltypes.BoolTypeOrNil(r.MaintainEqualityIndexWithoutFilter)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbIndexCacheModeProp(r.CacheMode), internaltypes.IsEmptyString(expectedValues.CacheMode))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLocalDbIndexOperations(plan localDbIndexResourceModel, state localDbIndexResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Attribute, state.Attribute, "attribute")
	operations.AddInt64OperationIfNecessary(&ops, plan.IndexEntryLimit, state.IndexEntryLimit, "index-entry-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.SubstringIndexEntryLimit, state.SubstringIndexEntryLimit, "substring-index-entry-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.MaintainMatchCountForKeysExceedingEntryLimit, state.MaintainMatchCountForKeysExceedingEntryLimit, "maintain-match-count-for-keys-exceeding-entry-limit")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexType, state.IndexType, "index-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.SubstringLength, state.SubstringLength, "substring-length")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeIndex, state.PrimeIndex, "prime-index")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeInternalNodesOnly, state.PrimeInternalNodesOnly, "prime-internal-nodes-only")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EqualityIndexFilter, state.EqualityIndexFilter, "equality-index-filter")
	operations.AddBoolOperationIfNecessary(&ops, plan.MaintainEqualityIndexWithoutFilter, state.MaintainEqualityIndexWithoutFilter, "maintain-equality-index-without-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.CacheMode, state.CacheMode, "cache-mode")
	return ops
}

// Create a new resource
func (r *localDbIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var IndexTypeSlice []client.EnumlocalDbIndexIndexTypeProp
	plan.IndexType.ElementsAs(ctx, &IndexTypeSlice, false)
	addRequest := client.NewAddLocalDbIndexRequest(plan.Attribute.ValueString(),
		plan.Attribute.ValueString(),
		IndexTypeSlice)
	err := addOptionalLocalDbIndexFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Db Index", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LocalDbIndexApi.AddLocalDbIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendName.ValueString())
	apiAddRequest = apiAddRequest.AddLocalDbIndexRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.LocalDbIndexApi.AddLocalDbIndexExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Local Db Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state localDbIndexResourceModel
	readLocalDbIndexResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLocalDbIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbIndexApi.GetLocalDbIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Attribute.ValueString(), plan.BackendName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state localDbIndexResourceModel
	readLocalDbIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LocalDbIndexApi.UpdateLocalDbIndex(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Attribute.ValueString(), plan.BackendName.ValueString())
	ops := createLocalDbIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LocalDbIndexApi.UpdateLocalDbIndexExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *localDbIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLocalDbIndex(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state localDbIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LocalDbIndexApi.GetLocalDbIndex(
		ProviderBasicAuthContext(ctx, providerConfig), state.Attribute.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *localDbIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbIndexResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLocalDbIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLocalDbIndex(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan localDbIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state localDbIndexResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LocalDbIndexApi.UpdateLocalDbIndex(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Attribute.ValueString(), plan.BackendName.ValueString())

	// Determine what update operations are necessary
	ops := createLocalDbIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LocalDbIndexApi.UpdateLocalDbIndexExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLocalDbIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *localDbIndexResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state localDbIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LocalDbIndexApi.DeleteLocalDbIndexExecute(r.apiClient.LocalDbIndexApi.DeleteLocalDbIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Attribute.ValueString(), state.BackendName.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Local Db Index", err, httpResp)
		return
	}
}

func (r *localDbIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbIndex(ctx, req, resp)
}

func (r *defaultLocalDbIndexResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLocalDbIndex(ctx, req, resp)
}

func importLocalDbIndex(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [backend-name]/[local-db-index-attribute]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("backend_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute"), split[1])...)
}
