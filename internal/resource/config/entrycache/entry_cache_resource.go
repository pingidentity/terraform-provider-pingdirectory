package entrycache

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &entryCacheResource{}
	_ resource.ResourceWithConfigure   = &entryCacheResource{}
	_ resource.ResourceWithImportState = &entryCacheResource{}
	_ resource.Resource                = &defaultEntryCacheResource{}
	_ resource.ResourceWithConfigure   = &defaultEntryCacheResource{}
	_ resource.ResourceWithImportState = &defaultEntryCacheResource{}
)

// Create a Entry Cache resource
func NewEntryCacheResource() resource.Resource {
	return &entryCacheResource{}
}

func NewDefaultEntryCacheResource() resource.Resource {
	return &defaultEntryCacheResource{}
}

// entryCacheResource is the resource implementation.
type entryCacheResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultEntryCacheResource is the resource implementation.
type defaultEntryCacheResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *entryCacheResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_entry_cache"
}

func (r *defaultEntryCacheResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_entry_cache"
}

// Configure adds the provider configured client to the resource.
func (r *entryCacheResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultEntryCacheResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type entryCacheResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	MaxMemoryPercent            types.Int64  `tfsdk:"max_memory_percent"`
	MaxEntries                  types.Int64  `tfsdk:"max_entries"`
	OnlyCacheFrequentlyAccessed types.Bool   `tfsdk:"only_cache_frequently_accessed"`
	IncludeFilter               types.Set    `tfsdk:"include_filter"`
	ExcludeFilter               types.Set    `tfsdk:"exclude_filter"`
	MinCacheEntryValueCount     types.Int64  `tfsdk:"min_cache_entry_value_count"`
	MinCacheEntryAttribute      types.Set    `tfsdk:"min_cache_entry_attribute"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	CacheLevel                  types.Int64  `tfsdk:"cache_level"`
	CacheUnindexedSearchResults types.Bool   `tfsdk:"cache_unindexed_search_results"`
}

// GetSchema defines the schema for the resource.
func (r *entryCacheResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	entryCacheSchema(ctx, req, resp, false)
}

func (r *defaultEntryCacheResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	entryCacheSchema(ctx, req, resp, true)
}

func entryCacheSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Entry Cache.",
		Attributes: map[string]schema.Attribute{
			"max_memory_percent": schema.Int64Attribute{
				Description: "Specifies the maximum amount of memory, as a percentage of the total maximum JVM heap size, that this cache should occupy when full. If the amount of memory the cache is using is greater than this amount, then an attempt to put a new entry in the cache will be ignored and will cause the oldest entry to be purged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_entries": schema.Int64Attribute{
				Description: "Specifies the maximum number of entries that will be allowed in the cache. Once the cache reaches this size, then adding new entries will cause existing entries to be purged, starting with the oldest.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"only_cache_frequently_accessed": schema.BoolAttribute{
				Description: "Specifies that the cache should only store entries which are accessed much more frequently than the average entry. The cache will observe attempts to place entries in the cache and compare an entry's accesses to the average entry's.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of filters that define the entries that should be included in the cache.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_filter": schema.SetAttribute{
				Description: "The set of filters that define the entries that should be excluded from the cache.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"min_cache_entry_value_count": schema.Int64Attribute{
				Description: "Specifies the minimum number of attribute values (optionally across a specified subset of attributes as defined in the min-cache-entry-attributes property) for entries that should be held in the cache. Entries with fewer than this number of attribute values will be excluded from the cache.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_cache_entry_attribute": schema.SetAttribute{
				Description: "Specifies the names of the attribute types for which the min-cache-entry-value-count property should apply. If no attribute types are specified, then all user attributes will be examined.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Entry Cache",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Entry Cache is enabled.",
				Required:    true,
			},
			"cache_level": schema.Int64Attribute{
				Description: "Specifies the cache level in the cache order if more than one instance of the cache is configured.",
				Required:    true,
			},
			"cache_unindexed_search_results": schema.BoolAttribute{
				Description: "Indicates whether the entry cache should be updated with entries that have been returned to the client during the course of processing an unindexed search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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

// Add optional fields to create request for fifo entry-cache
func addOptionalFifoEntryCacheFields(ctx context.Context, addRequest *client.AddFifoEntryCacheRequest, plan entryCacheResourceModel) {
	if internaltypes.IsDefined(plan.MaxMemoryPercent) {
		addRequest.MaxMemoryPercent = plan.MaxMemoryPercent.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxEntries) {
		addRequest.MaxEntries = plan.MaxEntries.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.OnlyCacheFrequentlyAccessed) {
		addRequest.OnlyCacheFrequentlyAccessed = plan.OnlyCacheFrequentlyAccessed.ValueBoolPointer()
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
	if internaltypes.IsDefined(plan.MinCacheEntryValueCount) {
		addRequest.MinCacheEntryValueCount = plan.MinCacheEntryValueCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MinCacheEntryAttribute) {
		var slice []string
		plan.MinCacheEntryAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MinCacheEntryAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CacheUnindexedSearchResults) {
		addRequest.CacheUnindexedSearchResults = plan.CacheUnindexedSearchResults.ValueBoolPointer()
	}
}

// Read a FifoEntryCacheResponse object into the model struct
func readFifoEntryCacheResponse(ctx context.Context, r *client.FifoEntryCacheResponse, state *entryCacheResourceModel, expectedValues *entryCacheResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.MaxMemoryPercent = internaltypes.Int64TypeOrNil(r.MaxMemoryPercent)
	state.MaxEntries = internaltypes.Int64TypeOrNil(r.MaxEntries)
	state.OnlyCacheFrequentlyAccessed = internaltypes.BoolTypeOrNil(r.OnlyCacheFrequentlyAccessed)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ExcludeFilter = internaltypes.GetStringSet(r.ExcludeFilter)
	state.MinCacheEntryValueCount = internaltypes.Int64TypeOrNil(r.MinCacheEntryValueCount)
	state.MinCacheEntryAttribute = internaltypes.GetStringSet(r.MinCacheEntryAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.CacheLevel = types.Int64Value(r.CacheLevel)
	state.CacheUnindexedSearchResults = internaltypes.BoolTypeOrNil(r.CacheUnindexedSearchResults)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createEntryCacheOperations(plan entryCacheResourceModel, state entryCacheResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxMemoryPercent, state.MaxMemoryPercent, "max-memory-percent")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxEntries, state.MaxEntries, "max-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.OnlyCacheFrequentlyAccessed, state.OnlyCacheFrequentlyAccessed, "only-cache-frequently-accessed")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeFilter, state.ExcludeFilter, "exclude-filter")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinCacheEntryValueCount, state.MinCacheEntryValueCount, "min-cache-entry-value-count")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MinCacheEntryAttribute, state.MinCacheEntryAttribute, "min-cache-entry-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddInt64OperationIfNecessary(&ops, plan.CacheLevel, state.CacheLevel, "cache-level")
	operations.AddBoolOperationIfNecessary(&ops, plan.CacheUnindexedSearchResults, state.CacheUnindexedSearchResults, "cache-unindexed-search-results")
	return ops
}

// Create a fifo entry-cache
func (r *entryCacheResource) CreateFifoEntryCache(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan entryCacheResourceModel) (*entryCacheResourceModel, error) {
	addRequest := client.NewAddFifoEntryCacheRequest(plan.Id.ValueString(),
		[]client.EnumfifoEntryCacheSchemaUrn{client.ENUMFIFOENTRYCACHESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ENTRY_CACHEFIFO},
		plan.Enabled.ValueBool(),
		plan.CacheLevel.ValueInt64())
	addOptionalFifoEntryCacheFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.EntryCacheApi.AddEntryCache(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFifoEntryCacheRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.EntryCacheApi.AddEntryCacheExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Entry Cache", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state entryCacheResourceModel
	readFifoEntryCacheResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *entryCacheResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan entryCacheResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateFifoEntryCache(ctx, req, resp, plan)
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
func (r *defaultEntryCacheResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan entryCacheResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.EntryCacheApi.GetEntryCache(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Entry Cache", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state entryCacheResourceModel
	readFifoEntryCacheResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.EntryCacheApi.UpdateEntryCache(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createEntryCacheOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.EntryCacheApi.UpdateEntryCacheExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Entry Cache", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFifoEntryCacheResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *entryCacheResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readEntryCache(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultEntryCacheResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readEntryCache(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readEntryCache(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state entryCacheResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.EntryCacheApi.GetEntryCache(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Entry Cache", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFifoEntryCacheResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *entryCacheResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateEntryCache(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultEntryCacheResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateEntryCache(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateEntryCache(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan entryCacheResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state entryCacheResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.EntryCacheApi.UpdateEntryCache(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createEntryCacheOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.EntryCacheApi.UpdateEntryCacheExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Entry Cache", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFifoEntryCacheResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultEntryCacheResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *entryCacheResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state entryCacheResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.EntryCacheApi.DeleteEntryCacheExecute(r.apiClient.EntryCacheApi.DeleteEntryCache(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Entry Cache", err, httpResp)
		return
	}
}

func (r *entryCacheResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importEntryCache(ctx, req, resp)
}

func (r *defaultEntryCacheResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importEntryCache(ctx, req, resp)
}

func importEntryCache(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
