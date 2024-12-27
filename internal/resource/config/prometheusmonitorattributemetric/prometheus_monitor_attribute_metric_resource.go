package prometheusmonitorattributemetric

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &prometheusMonitorAttributeMetricResource{}
	_ resource.ResourceWithConfigure   = &prometheusMonitorAttributeMetricResource{}
	_ resource.ResourceWithImportState = &prometheusMonitorAttributeMetricResource{}
	_ resource.Resource                = &defaultPrometheusMonitorAttributeMetricResource{}
	_ resource.ResourceWithConfigure   = &defaultPrometheusMonitorAttributeMetricResource{}
	_ resource.ResourceWithImportState = &defaultPrometheusMonitorAttributeMetricResource{}
)

// Create a Prometheus Monitor Attribute Metric resource
func NewPrometheusMonitorAttributeMetricResource() resource.Resource {
	return &prometheusMonitorAttributeMetricResource{}
}

func NewDefaultPrometheusMonitorAttributeMetricResource() resource.Resource {
	return &defaultPrometheusMonitorAttributeMetricResource{}
}

// prometheusMonitorAttributeMetricResource is the resource implementation.
type prometheusMonitorAttributeMetricResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPrometheusMonitorAttributeMetricResource is the resource implementation.
type defaultPrometheusMonitorAttributeMetricResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *prometheusMonitorAttributeMetricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_monitor_attribute_metric"
}

func (r *defaultPrometheusMonitorAttributeMetricResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_prometheus_monitor_attribute_metric"
}

// Configure adds the provider configured client to the resource.
func (r *prometheusMonitorAttributeMetricResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPrometheusMonitorAttributeMetricResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type prometheusMonitorAttributeMetricResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	MetricName               types.String `tfsdk:"metric_name"`
	MonitorAttributeName     types.String `tfsdk:"monitor_attribute_name"`
	MonitorObjectClassName   types.String `tfsdk:"monitor_object_class_name"`
	MetricType               types.String `tfsdk:"metric_type"`
	Filter                   types.String `tfsdk:"filter"`
	MetricDescription        types.String `tfsdk:"metric_description"`
	LabelNameValuePair       types.Set    `tfsdk:"label_name_value_pair"`
}

// GetSchema defines the schema for the resource.
func (r *prometheusMonitorAttributeMetricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	prometheusMonitorAttributeMetricSchema(ctx, req, resp, false)
}

func (r *defaultPrometheusMonitorAttributeMetricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	prometheusMonitorAttributeMetricSchema(ctx, req, resp, true)
}

func prometheusMonitorAttributeMetricSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Prometheus Monitor Attribute Metric. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Prometheus Monitor Attribute Metric resource. Options are ['prometheus-monitor-attribute-metric']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("prometheus-monitor-attribute-metric"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"prometheus-monitor-attribute-metric"}...),
				},
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"metric_name": schema.StringAttribute{
				Description: "The name that will be used in the metric to be consumed by Prometheus.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"monitor_attribute_name": schema.StringAttribute{
				Description: "The name of the monitor attribute that contains the numeric value to be published.",
				Required:    true,
			},
			"monitor_object_class_name": schema.StringAttribute{
				Description: "The name of the object class for monitor entries that contain the monitor attribute.",
				Required:    true,
			},
			"metric_type": schema.StringAttribute{
				Description: "The metric type that should be used for the value of the specified monitor attribute.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"counter", "gauge", "boolean"}...),
				},
			},
			"filter": schema.StringAttribute{
				Description: "A filter that may be used to restrict the set of monitor entries for which the metric should be generated.",
				Optional:    true,
			},
			"metric_description": schema.StringAttribute{
				Description: "A human-readable description that should be published as part of the metric definition.",
				Optional:    true,
			},
			"label_name_value_pair": schema.SetAttribute{
				Description: "A set of name-value pairs for labels that should be included in the published metric for the target attribute.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "metric_name", "http_servlet_extension_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		monitorAttributeNameAttr := schemaDef.Attributes["monitor_attribute_name"].(schema.StringAttribute)
		monitorAttributeNameAttr.PlanModifiers = append(monitorAttributeNameAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["monitor_attribute_name"] = monitorAttributeNameAttr
		monitorObjectClassNameAttr := schemaDef.Attributes["monitor_object_class_name"].(schema.StringAttribute)
		monitorObjectClassNameAttr.PlanModifiers = append(monitorObjectClassNameAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["monitor_object_class_name"] = monitorObjectClassNameAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *prometheusMonitorAttributeMetricResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_prometheus_monitor_attribute_metric")
}

func (r *defaultPrometheusMonitorAttributeMetricResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_prometheus_monitor_attribute_metric")
}

func modifyPlanPrometheusMonitorAttributeMetric(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		providerConfig.ProductVersion, resourceName)
}

// Add optional fields to create request for prometheus-monitor-attribute-metric prometheus-monitor-attribute-metric
func addOptionalPrometheusMonitorAttributeMetricFields(ctx context.Context, addRequest *client.AddPrometheusMonitorAttributeMetricRequest, plan prometheusMonitorAttributeMetricResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Filter) {
		addRequest.Filter = plan.Filter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MetricDescription) {
		addRequest.MetricDescription = plan.MetricDescription.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LabelNameValuePair) {
		var slice []string
		plan.LabelNameValuePair.ElementsAs(ctx, &slice, false)
		addRequest.LabelNameValuePair = slice
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *prometheusMonitorAttributeMetricResourceModel) populateAllComputedStringAttributes() {
	if model.MetricName.IsUnknown() || model.MetricName.IsNull() {
		model.MetricName = types.StringValue("")
	}
	if model.Filter.IsUnknown() || model.Filter.IsNull() {
		model.Filter = types.StringValue("")
	}
	if model.MetricType.IsUnknown() || model.MetricType.IsNull() {
		model.MetricType = types.StringValue("")
	}
	if model.MonitorAttributeName.IsUnknown() || model.MonitorAttributeName.IsNull() {
		model.MonitorAttributeName = types.StringValue("")
	}
	if model.MetricDescription.IsUnknown() || model.MetricDescription.IsNull() {
		model.MetricDescription = types.StringValue("")
	}
	if model.MonitorObjectClassName.IsUnknown() || model.MonitorObjectClassName.IsNull() {
		model.MonitorObjectClassName = types.StringValue("")
	}
}

// Read a PrometheusMonitorAttributeMetricResponse object into the model struct
func readPrometheusMonitorAttributeMetricResponse(ctx context.Context, r *client.PrometheusMonitorAttributeMetricResponse, state *prometheusMonitorAttributeMetricResourceModel, expectedValues *prometheusMonitorAttributeMetricResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("prometheus-monitor-attribute-metric")
	state.Id = types.StringValue(r.Id)
	state.MetricName = types.StringValue(r.MetricName)
	state.MonitorAttributeName = types.StringValue(r.MonitorAttributeName)
	state.MonitorObjectClassName = types.StringValue(r.MonitorObjectClassName)
	state.MetricType = types.StringValue(r.MetricType.String())
	state.Filter = internaltypes.StringTypeOrNil(r.Filter, internaltypes.IsEmptyString(expectedValues.Filter))
	state.MetricDescription = internaltypes.StringTypeOrNil(r.MetricDescription, internaltypes.IsEmptyString(expectedValues.MetricDescription))
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *prometheusMonitorAttributeMetricResourceModel) setStateValuesNotReturnedByAPI(expectedValues *prometheusMonitorAttributeMetricResourceModel) {
	if !expectedValues.HttpServletExtensionName.IsUnknown() {
		state.HttpServletExtensionName = expectedValues.HttpServletExtensionName
	}
}

// Create any update operations necessary to make the state match the plan
func createPrometheusMonitorAttributeMetricOperations(plan prometheusMonitorAttributeMetricResourceModel, state prometheusMonitorAttributeMetricResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.MetricName, state.MetricName, "metric-name")
	operations.AddStringOperationIfNecessary(&ops, plan.MonitorAttributeName, state.MonitorAttributeName, "monitor-attribute-name")
	operations.AddStringOperationIfNecessary(&ops, plan.MonitorObjectClassName, state.MonitorObjectClassName, "monitor-object-class-name")
	operations.AddStringOperationIfNecessary(&ops, plan.MetricType, state.MetricType, "metric-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringOperationIfNecessary(&ops, plan.MetricDescription, state.MetricDescription, "metric-description")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LabelNameValuePair, state.LabelNameValuePair, "label-name-value-pair")
	return ops
}

// Create a prometheus-monitor-attribute-metric prometheus-monitor-attribute-metric
func (r *prometheusMonitorAttributeMetricResource) CreatePrometheusMonitorAttributeMetric(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan prometheusMonitorAttributeMetricResourceModel) (*prometheusMonitorAttributeMetricResourceModel, error) {
	metricType, err := client.NewEnumprometheusMonitorAttributeMetricMetricTypePropFromValue(plan.MetricType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for MetricType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddPrometheusMonitorAttributeMetricRequest(plan.MetricName.ValueString(),
		plan.MonitorAttributeName.ValueString(),
		plan.MonitorObjectClassName.ValueString(),
		*metricType)
	addOptionalPrometheusMonitorAttributeMetricFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PrometheusMonitorAttributeMetricAPI.AddPrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.HttpServletExtensionName.ValueString())
	apiAddRequest = apiAddRequest.AddPrometheusMonitorAttributeMetricRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.AddPrometheusMonitorAttributeMetricExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Prometheus Monitor Attribute Metric", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state prometheusMonitorAttributeMetricResourceModel
	readPrometheusMonitorAttributeMetricResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *prometheusMonitorAttributeMetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan prometheusMonitorAttributeMetricResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreatePrometheusMonitorAttributeMetric(ctx, req, resp, plan)
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
func (r *defaultPrometheusMonitorAttributeMetricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan prometheusMonitorAttributeMetricResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.GetPrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.MetricName.ValueString(), plan.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitor Attribute Metric", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state prometheusMonitorAttributeMetricResourceModel
	readPrometheusMonitorAttributeMetricResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PrometheusMonitorAttributeMetricAPI.UpdatePrometheusMonitorAttributeMetric(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.MetricName.ValueString(), plan.HttpServletExtensionName.ValueString())
	ops := createPrometheusMonitorAttributeMetricOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.UpdatePrometheusMonitorAttributeMetricExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Prometheus Monitor Attribute Metric", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPrometheusMonitorAttributeMetricResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *prometheusMonitorAttributeMetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPrometheusMonitorAttributeMetricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPrometheusMonitorAttributeMetric(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state prometheusMonitorAttributeMetricResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PrometheusMonitorAttributeMetricAPI.GetPrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.MetricName.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitor Attribute Metric", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitor Attribute Metric", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPrometheusMonitorAttributeMetricResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *prometheusMonitorAttributeMetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPrometheusMonitorAttributeMetricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePrometheusMonitorAttributeMetric(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePrometheusMonitorAttributeMetric(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan prometheusMonitorAttributeMetricResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state prometheusMonitorAttributeMetricResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PrometheusMonitorAttributeMetricAPI.UpdatePrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.MetricName.ValueString(), plan.HttpServletExtensionName.ValueString())

	// Determine what update operations are necessary
	ops := createPrometheusMonitorAttributeMetricOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PrometheusMonitorAttributeMetricAPI.UpdatePrometheusMonitorAttributeMetricExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Prometheus Monitor Attribute Metric", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPrometheusMonitorAttributeMetricResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPrometheusMonitorAttributeMetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *prometheusMonitorAttributeMetricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state prometheusMonitorAttributeMetricResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.DeletePrometheusMonitorAttributeMetricExecute(r.apiClient.PrometheusMonitorAttributeMetricAPI.DeletePrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.MetricName.ValueString(), state.HttpServletExtensionName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Prometheus Monitor Attribute Metric", err, httpResp)
		return
	}
}

func (r *prometheusMonitorAttributeMetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPrometheusMonitorAttributeMetric(ctx, req, resp)
}

func (r *defaultPrometheusMonitorAttributeMetricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPrometheusMonitorAttributeMetric(ctx, req, resp)
}

func importPrometheusMonitorAttributeMetric(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [http-servlet-extension-name]/[prometheus-monitor-attribute-metric-metric-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("http_servlet_extension_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metric_name"), split[1])...)
}
