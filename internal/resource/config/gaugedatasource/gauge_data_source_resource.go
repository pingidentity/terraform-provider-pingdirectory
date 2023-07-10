package gaugedatasource

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &gaugeDataSourceResource{}
	_ resource.ResourceWithConfigure   = &gaugeDataSourceResource{}
	_ resource.ResourceWithImportState = &gaugeDataSourceResource{}
	_ resource.Resource                = &defaultGaugeDataSourceResource{}
	_ resource.ResourceWithConfigure   = &defaultGaugeDataSourceResource{}
	_ resource.ResourceWithImportState = &defaultGaugeDataSourceResource{}
)

// Create a Gauge Data Source resource
func NewGaugeDataSourceResource() resource.Resource {
	return &gaugeDataSourceResource{}
}

func NewDefaultGaugeDataSourceResource() resource.Resource {
	return &defaultGaugeDataSourceResource{}
}

// gaugeDataSourceResource is the resource implementation.
type gaugeDataSourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGaugeDataSourceResource is the resource implementation.
type defaultGaugeDataSourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *gaugeDataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gauge_data_source"
}

func (r *defaultGaugeDataSourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_gauge_data_source"
}

// Configure adds the provider configured client to the resource.
func (r *gaugeDataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultGaugeDataSourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type gaugeDataSourceResourceModel struct {
	Id                            types.String  `tfsdk:"id"`
	LastUpdated                   types.String  `tfsdk:"last_updated"`
	Notifications                 types.Set     `tfsdk:"notifications"`
	RequiredActions               types.Set     `tfsdk:"required_actions"`
	Type                          types.String  `tfsdk:"type"`
	DataOrientation               types.String  `tfsdk:"data_orientation"`
	StatisticType                 types.String  `tfsdk:"statistic_type"`
	DivideValueBy                 types.Float64 `tfsdk:"divide_value_by"`
	DivideValueByAttribute        types.String  `tfsdk:"divide_value_by_attribute"`
	DivideValueByCounterAttribute types.String  `tfsdk:"divide_value_by_counter_attribute"`
	Description                   types.String  `tfsdk:"description"`
	AdditionalText                types.String  `tfsdk:"additional_text"`
	MonitorObjectclass            types.String  `tfsdk:"monitor_objectclass"`
	MonitorAttribute              types.String  `tfsdk:"monitor_attribute"`
	IncludeFilter                 types.String  `tfsdk:"include_filter"`
	ResourceAttribute             types.String  `tfsdk:"resource_attribute"`
	ResourceType                  types.String  `tfsdk:"resource_type"`
	MinimumUpdateInterval         types.String  `tfsdk:"minimum_update_interval"`
}

// GetSchema defines the schema for the resource.
func (r *gaugeDataSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	gaugeDataSourceSchema(ctx, req, resp, false)
}

func (r *defaultGaugeDataSourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	gaugeDataSourceSchema(ctx, req, resp, true)
}

func gaugeDataSourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Gauge Data Source.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Gauge Data Source resource. Options are ['indicator', 'numeric']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"indicator", "numeric"}...),
				},
			},
			"data_orientation": schema.StringAttribute{
				Description: "Indicates whether a higher or lower value is a more severe condition.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"statistic_type": schema.StringAttribute{
				Description: "Specifies the type of statistic to include in the output for the monitored attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"divide_value_by": schema.Float64Attribute{
				Description: "An optional floating point value that can be used to scale the resulting value.",
				Optional:    true,
			},
			"divide_value_by_attribute": schema.StringAttribute{
				Description: "An optional property that can scale the resulting value by another attribute in the monitored entry.",
				Optional:    true,
			},
			"divide_value_by_counter_attribute": schema.StringAttribute{
				Description: "An optional property that can scale the resulting value by another attribute whose value represents a counter in the monitored entry.",
				Optional:    true,
			},
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
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"indicator", "numeric"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *gaugeDataSourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeDataSourceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanGaugeDataSource(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model gaugeDataSourceResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.DataOrientation) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'data_orientation' not supported by pingdirectory_gauge_data_source resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'data_orientation', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.DivideValueBy) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'divide_value_by' not supported by pingdirectory_gauge_data_source resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'divide_value_by', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.StatisticType) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'statistic_type' not supported by pingdirectory_gauge_data_source resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'statistic_type', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.DivideValueByAttribute) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'divide_value_by_attribute' not supported by pingdirectory_gauge_data_source resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'divide_value_by_attribute', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.DivideValueByCounterAttribute) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'divide_value_by_counter_attribute' not supported by pingdirectory_gauge_data_source resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'divide_value_by_counter_attribute', the 'type' attribute must be one of ['numeric']")
	}
}

// Add optional fields to create request for indicator gauge-data-source
func addOptionalIndicatorGaugeDataSourceFields(ctx context.Context, addRequest *client.AddIndicatorGaugeDataSourceRequest, plan gaugeDataSourceResourceModel) error {
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
	return nil
}

// Add optional fields to create request for numeric gauge-data-source
func addOptionalNumericGaugeDataSourceFields(ctx context.Context, addRequest *client.AddNumericGaugeDataSourceRequest, plan gaugeDataSourceResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DataOrientation) {
		dataOrientation, err := client.NewEnumgaugeDataSourceDataOrientationPropFromValue(plan.DataOrientation.ValueString())
		if err != nil {
			return err
		}
		addRequest.DataOrientation = dataOrientation
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.StatisticType) {
		statisticType, err := client.NewEnumgaugeDataSourceStatisticTypePropFromValue(plan.StatisticType.ValueString())
		if err != nil {
			return err
		}
		addRequest.StatisticType = statisticType
	}
	if internaltypes.IsDefined(plan.DivideValueBy) {
		addRequest.DivideValueBy = plan.DivideValueBy.ValueFloat64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DivideValueByAttribute) {
		addRequest.DivideValueByAttribute = plan.DivideValueByAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DivideValueByCounterAttribute) {
		addRequest.DivideValueByCounterAttribute = plan.DivideValueByCounterAttribute.ValueStringPointer()
	}
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
	return nil
}

// Read a IndicatorGaugeDataSourceResponse object into the model struct
func readIndicatorGaugeDataSourceResponse(ctx context.Context, r *client.IndicatorGaugeDataSourceResponse, state *gaugeDataSourceResourceModel, expectedValues *gaugeDataSourceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("indicator")
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

// Read a NumericGaugeDataSourceResponse object into the model struct
func readNumericGaugeDataSourceResponse(ctx context.Context, r *client.NumericGaugeDataSourceResponse, state *gaugeDataSourceResourceModel, expectedValues *gaugeDataSourceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("numeric")
	state.Id = types.StringValue(r.Id)
	state.DataOrientation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeDataSourceDataOrientationProp(r.DataOrientation), internaltypes.IsEmptyString(expectedValues.DataOrientation))
	state.StatisticType = types.StringValue(r.StatisticType.String())
	state.DivideValueBy = internaltypes.Float64TypeOrNil(r.DivideValueBy)
	state.DivideValueByAttribute = internaltypes.StringTypeOrNil(r.DivideValueByAttribute, internaltypes.IsEmptyString(expectedValues.DivideValueByAttribute))
	state.DivideValueByCounterAttribute = internaltypes.StringTypeOrNil(r.DivideValueByCounterAttribute, internaltypes.IsEmptyString(expectedValues.DivideValueByCounterAttribute))
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
func createGaugeDataSourceOperations(plan gaugeDataSourceResourceModel, state gaugeDataSourceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DataOrientation, state.DataOrientation, "data-orientation")
	operations.AddStringOperationIfNecessary(&ops, plan.StatisticType, state.StatisticType, "statistic-type")
	operations.AddFloat64OperationIfNecessary(&ops, plan.DivideValueBy, state.DivideValueBy, "divide-value-by")
	operations.AddStringOperationIfNecessary(&ops, plan.DivideValueByAttribute, state.DivideValueByAttribute, "divide-value-by-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DivideValueByCounterAttribute, state.DivideValueByCounterAttribute, "divide-value-by-counter-attribute")
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

// Create a indicator gauge-data-source
func (r *gaugeDataSourceResource) CreateIndicatorGaugeDataSource(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan gaugeDataSourceResourceModel) (*gaugeDataSourceResourceModel, error) {
	addRequest := client.NewAddIndicatorGaugeDataSourceRequest(plan.Id.ValueString(),
		[]client.EnumindicatorGaugeDataSourceSchemaUrn{client.ENUMINDICATORGAUGEDATASOURCESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGE_DATA_SOURCEINDICATOR},
		plan.MonitorObjectclass.ValueString(),
		plan.MonitorAttribute.ValueString())
	err := addOptionalIndicatorGaugeDataSourceFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Gauge Data Source", err.Error())
		return nil, err
	}
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Gauge Data Source", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state gaugeDataSourceResourceModel
	readIndicatorGaugeDataSourceResponse(ctx, addResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a numeric gauge-data-source
func (r *gaugeDataSourceResource) CreateNumericGaugeDataSource(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan gaugeDataSourceResourceModel) (*gaugeDataSourceResourceModel, error) {
	addRequest := client.NewAddNumericGaugeDataSourceRequest(plan.Id.ValueString(),
		[]client.EnumnumericGaugeDataSourceSchemaUrn{client.ENUMNUMERICGAUGEDATASOURCESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGE_DATA_SOURCENUMERIC},
		plan.MonitorObjectclass.ValueString(),
		plan.MonitorAttribute.ValueString())
	err := addOptionalNumericGaugeDataSourceFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Gauge Data Source", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.GaugeDataSourceApi.AddGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGaugeDataSourceRequest(
		client.AddNumericGaugeDataSourceRequestAsAddGaugeDataSourceRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.AddGaugeDataSourceExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Gauge Data Source", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state gaugeDataSourceResourceModel
	readNumericGaugeDataSourceResponse(ctx, addResponse.NumericGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *gaugeDataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *gaugeDataSourceResourceModel
	var err error
	if plan.Type.ValueString() == "indicator" {
		state, err = r.CreateIndicatorGaugeDataSource(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "numeric" {
		state, err = r.CreateNumericGaugeDataSource(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultGaugeDataSourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.GetGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state gaugeDataSourceResourceModel
	if plan.Type.ValueString() == "indicator" {
		readIndicatorGaugeDataSourceResponse(ctx, readResponse.IndicatorGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "numeric" {
		readNumericGaugeDataSourceResponse(ctx, readResponse.NumericGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GaugeDataSourceApi.UpdateGaugeDataSource(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGaugeDataSourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GaugeDataSourceApi.UpdateGaugeDataSourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gauge Data Source", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "indicator" {
			readIndicatorGaugeDataSourceResponse(ctx, updateResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "numeric" {
			readNumericGaugeDataSourceResponse(ctx, updateResponse.NumericGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
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
func (r *gaugeDataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeDataSourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGaugeDataSource(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state gaugeDataSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.GaugeDataSourceApi.GetGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge Data Source", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.IndicatorGaugeDataSourceResponse != nil {
		readIndicatorGaugeDataSourceResponse(ctx, readResponse.IndicatorGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NumericGaugeDataSourceResponse != nil {
		readNumericGaugeDataSourceResponse(ctx, readResponse.NumericGaugeDataSourceResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *gaugeDataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeDataSourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGaugeDataSource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGaugeDataSource(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan gaugeDataSourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state gaugeDataSourceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.GaugeDataSourceApi.UpdateGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGaugeDataSourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.GaugeDataSourceApi.UpdateGaugeDataSourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gauge Data Source", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "indicator" {
			readIndicatorGaugeDataSourceResponse(ctx, updateResponse.IndicatorGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "numeric" {
			readNumericGaugeDataSourceResponse(ctx, updateResponse.NumericGaugeDataSourceResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGaugeDataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *gaugeDataSourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state gaugeDataSourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GaugeDataSourceApi.DeleteGaugeDataSourceExecute(r.apiClient.GaugeDataSourceApi.DeleteGaugeDataSource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Gauge Data Source", err, httpResp)
		return
	}
}

func (r *gaugeDataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGaugeDataSource(ctx, req, resp)
}

func (r *defaultGaugeDataSourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGaugeDataSource(ctx, req, resp)
}

func importGaugeDataSource(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
