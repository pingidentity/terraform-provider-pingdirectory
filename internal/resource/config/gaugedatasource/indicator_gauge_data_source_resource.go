package gaugedatasource

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &indicatorGaugeDataSourceResource{}
	_ resource.ResourceWithConfigure   = &indicatorGaugeDataSourceResource{}
	_ resource.ResourceWithImportState = &indicatorGaugeDataSourceResource{}
	_ resource.Resource                = &defaultIndicatorGaugeDataSourceResource{}
	_ resource.ResourceWithConfigure   = &defaultIndicatorGaugeDataSourceResource{}
	_ resource.ResourceWithImportState = &defaultIndicatorGaugeDataSourceResource{}
)

// Create a Indicator Gauge Data Source resource
func NewIndicatorGaugeDataSourceResource() resource.Resource {
	return &indicatorGaugeDataSourceResource{}
}

func NewDefaultIndicatorGaugeDataSourceResource() resource.Resource {
	return &defaultIndicatorGaugeDataSourceResource{}
}

// indicatorGaugeDataSourceResource is the resource implementation.
type indicatorGaugeDataSourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultIndicatorGaugeDataSourceResource is the resource implementation.
type defaultIndicatorGaugeDataSourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *indicatorGaugeDataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indicator_gauge_data_source"
}

func (r *defaultIndicatorGaugeDataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_indicator_gauge_data_source"
}

// Configure adds the provider configured client to the resource.
func (r *indicatorGaugeDataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultIndicatorGaugeDataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type indicatorGaugeDataSourceResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	Description           types.String `tfsdk:"description"`
	AdditionalText        types.String `tfsdk:"additional_text"`
	MonitorObjectclass    types.String `tfsdk:"monitor_objectclass"`
	MonitorAttribute      types.String `tfsdk:"monitor_attribute"`
	IncludeFilter         types.String `tfsdk:"include_filter"`
	ResourceAttribute     types.String `tfsdk:"resource_attribute"`
	ResourceType          types.String `tfsdk:"resource_type"`
	MinimumUpdateInterval types.String `tfsdk:"minimum_update_interval"`
}

// GetSchema defines the schema for the resource.
func (r *indicatorGaugeDataSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	indicatorGaugeDataSourceSchema(ctx, req, resp, false)
}

func (r *defaultIndicatorGaugeDataSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	indicatorGaugeDataSourceSchema(ctx, req, resp, true)
}

func indicatorGaugeDataSourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Indicator Gauge Data Source.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Gauge Data Source",
				Optional:    true,
			},
			"additional_text": schema.StringAttribute{
				Description: "Additional information about the source of this data that is added to alerts sent as a result of gauges that use this Gauge Data Source.",
				Optional:    true,
			},
			"monitor_objectclass": schema.StringAttribute{
				Description: "The object class name of the monitor entries to examine for generating gauge data.",
				Required:    true,
			},
			"monitor_attribute": schema.StringAttribute{
				Description: "Specifies the attribute on the monitor entries from which to derive the current gauge value.",
				Required:    true,
			},
			"include_filter": schema.StringAttribute{
				Description: "An optional LDAP filter that can be used restrict which monitor entries are used to compute output.",
				Optional:    true,
			},
			"resource_attribute": schema.StringAttribute{
				Description: "Specifies the attribute whose value is used to identify the specific resource being monitored (e.g. device name).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_type": schema.StringAttribute{
				Description: "A string indicating the type of resource being monitored.",
				Optional:    true,
			},
			"minimum_update_interval": schema.StringAttribute{
				Description: "The minimum frequency with which gauges using this Gauge Data Source can be configured for update. In order to prevent undesirable side effects, some Gauge Data Sources may use this property to impose a higher bound on the update frequency of gauges.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalIndicatorGaugeDataSourceFields(ctx context.Context, addRequest *client.AddIndicatorGaugeDataSourceRequest, plan indicatorGaugeDataSourceResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdditionalText) {
		addRequest.AdditionalText = plan.AdditionalText.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludeFilter) {
		addRequest.IncludeFilter = plan.IncludeFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResourceAttribute) {
		addRequest.ResourceAttribute = plan.ResourceAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResourceType) {
		addRequest.ResourceType = plan.ResourceType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinimumUpdateInterval) {
		addRequest.MinimumUpdateInterval = plan.MinimumUpdateInterval.ValueStringPointer()
	}
}

// Read a IndicatorGaugeDataSourceResponse object into the model struct
func readIndicatorGaugeDataSourceResponse(ctx context.Context, r *client.IndicatorGaugeDataSourceResponse, state *indicatorGaugeDataSourceResourceModel, expectedValues *indicatorGaugeDataSourceResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AdditionalText = internaltypes.StringTypeOrNil(r.AdditionalText, internaltypes.IsEmptyString(expectedValues.AdditionalText))
	state.MonitorObjectclass = types.StringValue(r.MonitorObjectclass)
	state.MonitorAttribute = types.StringValue(r.MonitorAttribute)
	state.IncludeFilter = internaltypes.StringTypeOrNil(r.IncludeFilter, internaltypes.IsEmptyString(expectedValues.IncludeFilter))
	state.ResourceAttribute = internaltypes.StringTypeOrNil(r.ResourceAttribute, internaltypes.IsEmptyString(expectedValues.ResourceAttribute))
	state.ResourceType = internaltypes.StringTypeOrNil(r.ResourceType, internaltypes.IsEmptyString(expectedValues.ResourceType))
	state.MinimumUpdateInterval = internaltypes.StringTypeOrNil(r.MinimumUpdateInterval, internaltypes.IsEmptyString(expectedValues.MinimumUpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("minimum_update_interval",
		expectedValues.MinimumUpdateInterval, state.MinimumUpdateInterval, diagnostics)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createIndicatorGaugeDataSourceOperations(plan indicatorGaugeDataSourceResourceModel, state indicatorGaugeDataSourceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalText, state.AdditionalText, "additional-text")
	operations.AddStringOperationIfNecessary(&ops, plan.MonitorObjectclass, state.MonitorObjectclass, "monitor-objectclass")
	operations.AddStringOperationIfNecessary(&ops, plan.MonitorAttribute, state.MonitorAttribute, "monitor-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceAttribute, state.ResourceAttribute, "resource-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceType, state.ResourceType, "resource-type")
	operations.AddStringOperationIfNecessary(&ops, plan.MinimumUpdateInterval, state.MinimumUpdateInterval, "minimum-update-interval")
	return ops
}

// Create a new resource
func (r *indicatorGaugeDataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indicatorGaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddIndicatorGaugeDataSourceRequest(plan.Id.ValueString(),
		[]client.EnumindicatorGaugeDataSourceSchemaUrn{client.ENUMINDICATORGAUGEDATASOURCESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGE_DATA_SOURCEINDICATOR},
		plan.MonitorObjectclass.ValueString(),
		plan.MonitorAttribute.ValueString())
	addOptionalIndicatorGaugeDataSourceFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.GaugeDataSourceApi.AddGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGaugeDataSourceRequest(
		client.AddIndicatorGaugeDataSourceRequestAsAddGaugeDataSourceRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.AddGaugeDataSourceExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Indicator Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state indicatorGaugeDataSourceResourceModel
	readIndicatorGaugeDataSourceResponse(ctx, addResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultIndicatorGaugeDataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indicatorGaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.GetGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Indicator Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state indicatorGaugeDataSourceResourceModel
	readIndicatorGaugeDataSourceResponse(ctx, readResponse.IndicatorGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GaugeDataSourceApi.UpdateGaugeDataSource(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createIndicatorGaugeDataSourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.UpdateGaugeDataSourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Indicator Gauge Data Source", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIndicatorGaugeDataSourceResponse(ctx, updateResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
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
func (r *indicatorGaugeDataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIndicatorGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIndicatorGaugeDataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIndicatorGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIndicatorGaugeDataSource(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state indicatorGaugeDataSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.GaugeDataSourceApi.GetGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Indicator Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIndicatorGaugeDataSourceResponse(ctx, readResponse.IndicatorGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *indicatorGaugeDataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIndicatorGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIndicatorGaugeDataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIndicatorGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIndicatorGaugeDataSource(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan indicatorGaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state indicatorGaugeDataSourceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.GaugeDataSourceApi.UpdateGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createIndicatorGaugeDataSourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.GaugeDataSourceApi.UpdateGaugeDataSourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Indicator Gauge Data Source", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIndicatorGaugeDataSourceResponse(ctx, updateResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultIndicatorGaugeDataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *indicatorGaugeDataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state indicatorGaugeDataSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GaugeDataSourceApi.DeleteGaugeDataSourceExecute(r.apiClient.GaugeDataSourceApi.DeleteGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Indicator Gauge Data Source", err, httpResp)
		return
	}
}

func (r *indicatorGaugeDataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIndicatorGaugeDataSource(ctx, req, resp)
}

func (r *defaultIndicatorGaugeDataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIndicatorGaugeDataSource(ctx, req, resp)
}

func importIndicatorGaugeDataSource(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
