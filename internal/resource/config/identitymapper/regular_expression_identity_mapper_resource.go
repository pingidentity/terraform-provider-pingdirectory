package identitymapper

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
	_ resource.Resource                = &regularExpressionIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &regularExpressionIdentityMapperResource{}
	_ resource.ResourceWithImportState = &regularExpressionIdentityMapperResource{}
)

// Create a Regular Expression Identity Mapper resource
func NewRegularExpressionIdentityMapperResource() resource.Resource {
	return &regularExpressionIdentityMapperResource{}
}

// regularExpressionIdentityMapperResource is the resource implementation.
type regularExpressionIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *regularExpressionIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regular_expression_identity_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *regularExpressionIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type regularExpressionIdentityMapperResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	MatchAttribute  types.Set    `tfsdk:"match_attribute"`
	MatchBaseDN     types.Set    `tfsdk:"match_base_dn"`
	MatchFilter     types.String `tfsdk:"match_filter"`
	MatchPattern    types.String `tfsdk:"match_pattern"`
	ReplacePattern  types.String `tfsdk:"replace_pattern"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *regularExpressionIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Regular Expression Identity Mapper.",
		Attributes: map[string]schema.Attribute{
			"match_attribute": schema.SetAttribute{
				Description: "Specifies the name or OID of the attribute whose value should match the provided identifier string after it has been processed by the associated regular expression.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"match_base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) that should be used when performing searches to map the provided ID string to a user entry. If multiple values are given, searches are performed below all the specified base DNs.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"match_filter": schema.StringAttribute{
				Description: "An optional filter that mapped users must match.",
				Optional:    true,
			},
			"match_pattern": schema.StringAttribute{
				Description: "Specifies the regular expression pattern that is used to identify portions of the ID string that will be replaced.",
				Required:    true,
			},
			"replace_pattern": schema.StringAttribute{
				Description: "Specifies the replacement pattern that should be used for substrings in the ID string that match the provided regular expression pattern.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Identity Mapper",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Identity Mapper is enabled for use.",
				Required:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalRegularExpressionIdentityMapperFields(ctx context.Context, addRequest *client.AddRegularExpressionIdentityMapperRequest, plan regularExpressionIdentityMapperResourceModel) {
	if internaltypes.IsDefined(plan.MatchAttribute) {
		var slice []string
		plan.MatchAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MatchAttribute = slice
	}
	if internaltypes.IsDefined(plan.MatchBaseDN) {
		var slice []string
		plan.MatchBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.MatchBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MatchFilter) {
		stringVal := plan.MatchFilter.ValueString()
		addRequest.MatchFilter = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReplacePattern) {
		stringVal := plan.ReplacePattern.ValueString()
		addRequest.ReplacePattern = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a RegularExpressionIdentityMapperResponse object into the model struct
func readRegularExpressionIdentityMapperResponse(ctx context.Context, r *client.RegularExpressionIdentityMapperResponse, state *regularExpressionIdentityMapperResourceModel, expectedValues *regularExpressionIdentityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.MatchBaseDN = internaltypes.GetStringSet(r.MatchBaseDN)
	state.MatchFilter = internaltypes.StringTypeOrNil(r.MatchFilter, internaltypes.IsEmptyString(expectedValues.MatchFilter))
	state.MatchPattern = types.StringValue(r.MatchPattern)
	state.ReplacePattern = internaltypes.StringTypeOrNil(r.ReplacePattern, internaltypes.IsEmptyString(expectedValues.ReplacePattern))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createRegularExpressionIdentityMapperOperations(plan regularExpressionIdentityMapperResourceModel, state regularExpressionIdentityMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MatchAttribute, state.MatchAttribute, "match-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MatchBaseDN, state.MatchBaseDN, "match-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchFilter, state.MatchFilter, "match-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchPattern, state.MatchPattern, "match-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplacePattern, state.ReplacePattern, "replace-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *regularExpressionIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan regularExpressionIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddRegularExpressionIdentityMapperRequest(plan.Id.ValueString(),
		[]client.EnumregularExpressionIdentityMapperSchemaUrn{client.ENUMREGULAREXPRESSIONIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERREGULAR_EXPRESSION},
		plan.MatchPattern.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalRegularExpressionIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddRegularExpressionIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Regular Expression Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state regularExpressionIdentityMapperResourceModel
	readRegularExpressionIdentityMapperResponse(ctx, addResponse.RegularExpressionIdentityMapperResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *regularExpressionIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state regularExpressionIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Regular Expression Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readRegularExpressionIdentityMapperResponse(ctx, readResponse.RegularExpressionIdentityMapperResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *regularExpressionIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan regularExpressionIdentityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state regularExpressionIdentityMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.IdentityMapperApi.UpdateIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createRegularExpressionIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Regular Expression Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readRegularExpressionIdentityMapperResponse(ctx, updateResponse.RegularExpressionIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *regularExpressionIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state regularExpressionIdentityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdentityMapperApi.DeleteIdentityMapperExecute(r.apiClient.IdentityMapperApi.DeleteIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Regular Expression Identity Mapper", err, httpResp)
		return
	}
}

func (r *regularExpressionIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
