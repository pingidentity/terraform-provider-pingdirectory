package httpservletextension

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &prometheusMonitoringHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &prometheusMonitoringHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &prometheusMonitoringHttpServletExtensionResource{}
	_ resource.Resource                = &defaultPrometheusMonitoringHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultPrometheusMonitoringHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &defaultPrometheusMonitoringHttpServletExtensionResource{}
)

// Create a Prometheus Monitoring Http Servlet Extension resource
func NewPrometheusMonitoringHttpServletExtensionResource() resource.Resource {
	return &prometheusMonitoringHttpServletExtensionResource{}
}

func NewDefaultPrometheusMonitoringHttpServletExtensionResource() resource.Resource {
	return &defaultPrometheusMonitoringHttpServletExtensionResource{}
}

// prometheusMonitoringHttpServletExtensionResource is the resource implementation.
type prometheusMonitoringHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPrometheusMonitoringHttpServletExtensionResource is the resource implementation.
type defaultPrometheusMonitoringHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *prometheusMonitoringHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_monitoring_http_servlet_extension"
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_prometheus_monitoring_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *prometheusMonitoringHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type prometheusMonitoringHttpServletExtensionResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	BaseContextPath                    types.String `tfsdk:"base_context_path"`
	IncludeInstanceNameLabel           types.Bool   `tfsdk:"include_instance_name_label"`
	IncludeProductNameLabel            types.Bool   `tfsdk:"include_product_name_label"`
	IncludeLocationNameLabel           types.Bool   `tfsdk:"include_location_name_label"`
	AlwaysIncludeMonitorEntryNameLabel types.Bool   `tfsdk:"always_include_monitor_entry_name_label"`
	IncludeMonitorObjectClassNameLabel types.Bool   `tfsdk:"include_monitor_object_class_name_label"`
	IncludeMonitorAttributeNameLabel   types.Bool   `tfsdk:"include_monitor_attribute_name_label"`
	LabelNameValuePair                 types.Set    `tfsdk:"label_name_value_pair"`
	Description                        types.String `tfsdk:"description"`
	CrossOriginPolicy                  types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader                     types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader        types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *prometheusMonitoringHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	prometheusMonitoringHttpServletExtensionSchema(ctx, req, resp, false)
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	prometheusMonitoringHttpServletExtensionSchema(ctx, req, resp, true)
}

func prometheusMonitoringHttpServletExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Prometheus Monitoring Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_instance_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include an \"instance\" label whose value is the instance name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_product_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"product\" label whose value is the product name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_location_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"location\" label whose value is the location name for this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"always_include_monitor_entry_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should always include a \"monitor_entry\" label whose value is the name of the monitor entry from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_monitor_object_class_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_object_class\" label whose value is the name of the object class for the monitor entry from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_monitor_attribute_name_label": schema.BoolAttribute{
				Description: "Indicates whether generated metrics should include a \"monitor_attribute\" label whose value is the name of the monitor attribute from which the metric was obtained.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"label_name_value_pair": schema.SetAttribute{
				Description: "A set of name-value pairs for labels that should be included in all metrics exposed by this Directory Server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
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

// Validate that this resource is being used with a compatible PingDirectory version
func (r *prometheusMonitoringHttpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		r.providerConfig.ServerVersion, "pingdirectory_prometheus_monitoring_http_servlet_extension")
}

// Validate that this resource is being used with a compatible PingDirectory version
func (r *defaultPrometheusMonitoringHttpServletExtensionResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
		r.providerConfig.ServerVersion, "pingdirectory_default_prometheus_monitoring_http_servlet_extension")
}

// Add optional fields to create request
func addOptionalPrometheusMonitoringHttpServletExtensionFields(ctx context.Context, addRequest *client.AddPrometheusMonitoringHttpServletExtensionRequest, plan prometheusMonitoringHttpServletExtensionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BaseContextPath) {
		stringVal := plan.BaseContextPath.ValueString()
		addRequest.BaseContextPath = &stringVal
	}
	if internaltypes.IsDefined(plan.IncludeInstanceNameLabel) {
		boolVal := plan.IncludeInstanceNameLabel.ValueBool()
		addRequest.IncludeInstanceNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeProductNameLabel) {
		boolVal := plan.IncludeProductNameLabel.ValueBool()
		addRequest.IncludeProductNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeLocationNameLabel) {
		boolVal := plan.IncludeLocationNameLabel.ValueBool()
		addRequest.IncludeLocationNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.AlwaysIncludeMonitorEntryNameLabel) {
		boolVal := plan.AlwaysIncludeMonitorEntryNameLabel.ValueBool()
		addRequest.AlwaysIncludeMonitorEntryNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeMonitorObjectClassNameLabel) {
		boolVal := plan.IncludeMonitorObjectClassNameLabel.ValueBool()
		addRequest.IncludeMonitorObjectClassNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeMonitorAttributeNameLabel) {
		boolVal := plan.IncludeMonitorAttributeNameLabel.ValueBool()
		addRequest.IncludeMonitorAttributeNameLabel = &boolVal
	}
	if internaltypes.IsDefined(plan.LabelNameValuePair) {
		var slice []string
		plan.LabelNameValuePair.ElementsAs(ctx, &slice, false)
		addRequest.LabelNameValuePair = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		stringVal := plan.CrossOriginPolicy.ValueString()
		addRequest.CrossOriginPolicy = &stringVal
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		stringVal := plan.CorrelationIDResponseHeader.ValueString()
		addRequest.CorrelationIDResponseHeader = &stringVal
	}
}

// Read a PrometheusMonitoringHttpServletExtensionResponse object into the model struct
func readPrometheusMonitoringHttpServletExtensionResponse(ctx context.Context, r *client.PrometheusMonitoringHttpServletExtensionResponse, state *prometheusMonitoringHttpServletExtensionResourceModel, expectedValues *prometheusMonitoringHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.IncludeInstanceNameLabel = internaltypes.BoolTypeOrNil(r.IncludeInstanceNameLabel)
	state.IncludeProductNameLabel = internaltypes.BoolTypeOrNil(r.IncludeProductNameLabel)
	state.IncludeLocationNameLabel = internaltypes.BoolTypeOrNil(r.IncludeLocationNameLabel)
	state.AlwaysIncludeMonitorEntryNameLabel = internaltypes.BoolTypeOrNil(r.AlwaysIncludeMonitorEntryNameLabel)
	state.IncludeMonitorObjectClassNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorObjectClassNameLabel)
	state.IncludeMonitorAttributeNameLabel = internaltypes.BoolTypeOrNil(r.IncludeMonitorAttributeNameLabel)
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPrometheusMonitoringHttpServletExtensionOperations(plan prometheusMonitoringHttpServletExtensionResourceModel, state prometheusMonitoringHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceNameLabel, state.IncludeInstanceNameLabel, "include-instance-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductNameLabel, state.IncludeProductNameLabel, "include-product-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeLocationNameLabel, state.IncludeLocationNameLabel, "include-location-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.AlwaysIncludeMonitorEntryNameLabel, state.AlwaysIncludeMonitorEntryNameLabel, "always-include-monitor-entry-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorObjectClassNameLabel, state.IncludeMonitorObjectClassNameLabel, "include-monitor-object-class-name-label")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeMonitorAttributeNameLabel, state.IncludeMonitorAttributeNameLabel, "include-monitor-attribute-name-label")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LabelNameValuePair, state.LabelNameValuePair, "label-name-value-pair")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
func (r *prometheusMonitoringHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan prometheusMonitoringHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPrometheusMonitoringHttpServletExtensionRequest(plan.Id.ValueString(),
		[]client.EnumprometheusMonitoringHttpServletExtensionSchemaUrn{client.ENUMPROMETHEUSMONITORINGHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONPROMETHEUS_MONITORING})
	addOptionalPrometheusMonitoringHttpServletExtensionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddPrometheusMonitoringHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Prometheus Monitoring Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state prometheusMonitoringHttpServletExtensionResourceModel
	readPrometheusMonitoringHttpServletExtensionResponse(ctx, addResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan prometheusMonitoringHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitoring Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state prometheusMonitoringHttpServletExtensionResourceModel
	readPrometheusMonitoringHttpServletExtensionResponse(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPrometheusMonitoringHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Prometheus Monitoring Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPrometheusMonitoringHttpServletExtensionResponse(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *prometheusMonitoringHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPrometheusMonitoringHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPrometheusMonitoringHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPrometheusMonitoringHttpServletExtension(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state prometheusMonitoringHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitoring Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPrometheusMonitoringHttpServletExtensionResponse(ctx, readResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *prometheusMonitoringHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePrometheusMonitoringHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePrometheusMonitoringHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePrometheusMonitoringHttpServletExtension(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan prometheusMonitoringHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state prometheusMonitoringHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPrometheusMonitoringHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Prometheus Monitoring Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPrometheusMonitoringHttpServletExtensionResponse(ctx, updateResponse.PrometheusMonitoringHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPrometheusMonitoringHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *prometheusMonitoringHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state prometheusMonitoringHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtensionExecute(r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Prometheus Monitoring Http Servlet Extension", err, httpResp)
		return
	}
}

func (r *prometheusMonitoringHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPrometheusMonitoringHttpServletExtension(ctx, req, resp)
}

func (r *defaultPrometheusMonitoringHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPrometheusMonitoringHttpServletExtension(ctx, req, resp)
}

func importPrometheusMonitoringHttpServletExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
