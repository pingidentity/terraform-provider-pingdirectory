package config

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
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
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLocalDbVlvIndexResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type localDbVlvIndexResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
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

func localDbVlvIndexSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Local Db Vlv Index.",
		Attributes: map[string]schema.Attribute{
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the cache mode that should be used when accessing the records in the database for this index.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		SetAllAttributesToOptionalAndComputed(&schema, []string{"name", "backend_name"})
	}
	AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Add optional fields to create request
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

// Read a LocalDbVlvIndexResponse object into the model struct
func readLocalDbVlvIndexResponse(ctx context.Context, r *client.LocalDbVlvIndexResponse, state *localDbVlvIndexResourceModel, expectedValues *localDbVlvIndexResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BackendName = expectedValues.BackendName
	state.BaseDN = types.StringValue(r.BaseDN)
	state.Scope = types.StringValue(r.Scope.String())
	state.Filter = types.StringValue(r.Filter)
	state.SortOrder = types.StringValue(r.SortOrder)
	state.Name = types.StringValue(r.Name)
	state.MaxBlockSize = internaltypes.Int64TypeOrNil(r.MaxBlockSize)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlocalDbVlvIndexCacheModeProp(r.CacheMode), internaltypes.IsEmptyString(expectedValues.CacheMode))
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
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

// Create a new resource
func (r *localDbVlvIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbVlvIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope, err := client.NewEnumlocalDbVlvIndexScopePropFromValue(plan.Scope.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for Scope", err.Error())
		return
	}
	addRequest := client.NewAddLocalDbVlvIndexRequest(plan.Name.ValueString(),
		plan.BaseDN.ValueString(),
		*scope,
		plan.Filter.ValueString(),
		plan.SortOrder.ValueString(),
		plan.Name.ValueString())
	err = addOptionalLocalDbVlvIndexFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Local Db Vlv Index", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LocalDbVlvIndexApi.AddLocalDbVlvIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendName.ValueString())
	apiAddRequest = apiAddRequest.AddLocalDbVlvIndexRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.LocalDbVlvIndexApi.AddLocalDbVlvIndexExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Local Db Vlv Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state localDbVlvIndexResourceModel
	readLocalDbVlvIndexResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLocalDbVlvIndexResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan localDbVlvIndexResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LocalDbVlvIndexApi.GetLocalDbVlvIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
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
	updateRequest := r.apiClient.LocalDbVlvIndexApi.UpdateLocalDbVlvIndex(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())
	ops := createLocalDbVlvIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LocalDbVlvIndexApi.UpdateLocalDbVlvIndexExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Vlv Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbVlvIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *localDbVlvIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLocalDbVlvIndexResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLocalDbVlvIndex(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLocalDbVlvIndex(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state localDbVlvIndexResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LocalDbVlvIndexApi.GetLocalDbVlvIndex(
		ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.BackendName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Local Db Vlv Index", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLocalDbVlvIndexResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
	updateRequest := apiClient.LocalDbVlvIndexApi.UpdateLocalDbVlvIndex(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.BackendName.ValueString())

	// Determine what update operations are necessary
	ops := createLocalDbVlvIndexOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LocalDbVlvIndexApi.UpdateLocalDbVlvIndexExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Local Db Vlv Index", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLocalDbVlvIndexResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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

	httpResp, err := r.apiClient.LocalDbVlvIndexApi.DeleteLocalDbVlvIndexExecute(r.apiClient.LocalDbVlvIndexApi.DeleteLocalDbVlvIndex(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.BackendName.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Local Db Vlv Index", err, httpResp)
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