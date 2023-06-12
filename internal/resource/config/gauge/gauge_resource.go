package gauge

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &gaugeResource{}
	_ resource.ResourceWithConfigure   = &gaugeResource{}
	_ resource.ResourceWithImportState = &gaugeResource{}
	_ resource.Resource                = &defaultGaugeResource{}
	_ resource.ResourceWithConfigure   = &defaultGaugeResource{}
	_ resource.ResourceWithImportState = &defaultGaugeResource{}
)

// Create a Gauge resource
func NewGaugeResource() resource.Resource {
	return &gaugeResource{}
}

func NewDefaultGaugeResource() resource.Resource {
	return &defaultGaugeResource{}
}

// gaugeResource is the resource implementation.
type gaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGaugeResource is the resource implementation.
type defaultGaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *gaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gauge"
}

func (r *defaultGaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_gauge"
}

// Configure adds the provider configured client to the resource.
func (r *gaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type gaugeResourceModel struct {
	Id                             types.String  `tfsdk:"id"`
	LastUpdated                    types.String  `tfsdk:"last_updated"`
	Notifications                  types.Set     `tfsdk:"notifications"`
	RequiredActions                types.Set     `tfsdk:"required_actions"`
	Type                           types.String  `tfsdk:"type"`
	GaugeDataSource                types.String  `tfsdk:"gauge_data_source"`
	CriticalValue                  types.String  `tfsdk:"critical_value"`
	CriticalExitValue              types.Float64 `tfsdk:"critical_exit_value"`
	MajorValue                     types.String  `tfsdk:"major_value"`
	MajorExitValue                 types.Float64 `tfsdk:"major_exit_value"`
	MinorValue                     types.String  `tfsdk:"minor_value"`
	MinorExitValue                 types.Float64 `tfsdk:"minor_exit_value"`
	WarningValue                   types.String  `tfsdk:"warning_value"`
	WarningExitValue               types.Float64 `tfsdk:"warning_exit_value"`
	Description                    types.String  `tfsdk:"description"`
	Enabled                        types.Bool    `tfsdk:"enabled"`
	OverrideSeverity               types.String  `tfsdk:"override_severity"`
	AlertLevel                     types.String  `tfsdk:"alert_level"`
	UpdateInterval                 types.String  `tfsdk:"update_interval"`
	SamplesPerUpdateInterval       types.Int64   `tfsdk:"samples_per_update_interval"`
	IncludeResource                types.Set     `tfsdk:"include_resource"`
	ExcludeResource                types.Set     `tfsdk:"exclude_resource"`
	ServerUnavailableSeverityLevel types.String  `tfsdk:"server_unavailable_severity_level"`
	ServerDegradedSeverityLevel    types.String  `tfsdk:"server_degraded_severity_level"`
}

// GetSchema defines the schema for the resource.
func (r *gaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	gaugeSchema(ctx, req, resp, false)
}

func (r *defaultGaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	gaugeSchema(ctx, req, resp, true)
}

func gaugeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Gauge.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Gauge resource. Options are ['indicator', 'numeric']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"indicator", "numeric"}...),
				},
			},
			"gauge_data_source": schema.StringAttribute{
				Description: "Specifies the source of data to use in determining this Indicator Gauge's severity and status.",
				Required:    true,
			},
			"critical_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity should be critical.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"critical_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'critical'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"major_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'major'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"major_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'major'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"minor_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'minor'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"minor_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'minor'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"warning_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'warning'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warning_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'warning'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Gauge",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Gauge is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"override_severity": schema.StringAttribute{
				Description: "When defined, causes this Gauge to assume the specified severity, overriding its computed severity. This is useful for testing alarms generated by Gauges as well as suppressing alarms for known conditions.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_level": schema.StringAttribute{
				Description: "Specifies the level at which alerts are sent for alarms raised by this Gauge.",
				Optional:    true,
			},
			"update_interval": schema.StringAttribute{
				Description: "The frequency with which this Gauge is updated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"samples_per_update_interval": schema.Int64Attribute{
				Description: "Indicates the number of times the monitor data source value will be collected during the update interval.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"include_resource": schema.SetAttribute{
				Description: "Specifies set of resources to be monitored.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_resource": schema.SetAttribute{
				Description: "Specifies resources to exclude from being monitored.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server_unavailable_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_degraded_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered degraded.",
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
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *gaugeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanGauge(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model gaugeResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.MinorExitValue) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'minor_exit_value' not supported by pingdirectory_gauge resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'minor_exit_value', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.CriticalExitValue) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'critical_exit_value' not supported by pingdirectory_gauge resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'critical_exit_value', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.MajorExitValue) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'major_exit_value' not supported by pingdirectory_gauge resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'major_exit_value', the 'type' attribute must be one of ['numeric']")
	}
	if internaltypes.IsDefined(model.WarningExitValue) && model.Type.ValueString() != "numeric" {
		resp.Diagnostics.AddError("Attribute 'warning_exit_value' not supported by pingdirectory_gauge resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'warning_exit_value', the 'type' attribute must be one of ['numeric']")
	}
}

// Add optional fields to create request for indicator gauge
func addOptionalIndicatorGaugeFields(ctx context.Context, addRequest *client.AddIndicatorGaugeRequest, plan gaugeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CriticalValue) {
		addRequest.CriticalValue = plan.CriticalValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MajorValue) {
		addRequest.MajorValue = plan.MajorValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinorValue) {
		addRequest.MinorValue = plan.MinorValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WarningValue) {
		addRequest.WarningValue = plan.WarningValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OverrideSeverity) {
		overrideSeverity, err := client.NewEnumgaugeOverrideSeverityPropFromValue(plan.OverrideSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.OverrideSeverity = overrideSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AlertLevel) {
		alertLevel, err := client.NewEnumgaugeAlertLevelPropFromValue(plan.AlertLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.AlertLevel = alertLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateInterval) {
		addRequest.UpdateInterval = plan.UpdateInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SamplesPerUpdateInterval) {
		addRequest.SamplesPerUpdateInterval = plan.SamplesPerUpdateInterval.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeResource) {
		var slice []string
		plan.IncludeResource.ElementsAs(ctx, &slice, false)
		addRequest.IncludeResource = slice
	}
	if internaltypes.IsDefined(plan.ExcludeResource) {
		var slice []string
		plan.ExcludeResource.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeResource = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerUnavailableSeverityLevel) {
		serverUnavailableSeverityLevel, err := client.NewEnumgaugeServerUnavailableSeverityLevelPropFromValue(plan.ServerUnavailableSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerUnavailableSeverityLevel = serverUnavailableSeverityLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerDegradedSeverityLevel) {
		serverDegradedSeverityLevel, err := client.NewEnumgaugeServerDegradedSeverityLevelPropFromValue(plan.ServerDegradedSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerDegradedSeverityLevel = serverDegradedSeverityLevel
	}
	return nil
}

// Add optional fields to create request for numeric gauge
func addOptionalNumericGaugeFields(ctx context.Context, addRequest *client.AddNumericGaugeRequest, plan gaugeResourceModel) error {
	if internaltypes.IsDefined(plan.CriticalValue) {
		floatVal, err := strconv.ParseFloat(plan.CriticalValue.ValueString(), 64)
		if err != nil {
			return err
		}
		addRequest.CriticalValue = &floatVal
	}
	if internaltypes.IsDefined(plan.CriticalExitValue) {
		addRequest.CriticalExitValue = plan.CriticalExitValue.ValueFloat64Pointer()
	}
	if internaltypes.IsDefined(plan.MajorValue) {
		floatVal, err := strconv.ParseFloat(plan.MajorValue.ValueString(), 64)
		if err != nil {
			return err
		}
		addRequest.MajorValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MajorExitValue) {
		addRequest.MajorExitValue = plan.MajorExitValue.ValueFloat64Pointer()
	}
	if internaltypes.IsDefined(plan.MinorValue) {
		floatVal, err := strconv.ParseFloat(plan.MinorValue.ValueString(), 64)
		if err != nil {
			return err
		}
		addRequest.MinorValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MinorExitValue) {
		addRequest.MinorExitValue = plan.MinorExitValue.ValueFloat64Pointer()
	}
	if internaltypes.IsDefined(plan.WarningValue) {
		floatVal, err := strconv.ParseFloat(plan.WarningValue.ValueString(), 64)
		if err != nil {
			return err
		}
		addRequest.WarningValue = &floatVal
	}
	if internaltypes.IsDefined(plan.WarningExitValue) {
		addRequest.WarningExitValue = plan.WarningExitValue.ValueFloat64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OverrideSeverity) {
		overrideSeverity, err := client.NewEnumgaugeOverrideSeverityPropFromValue(plan.OverrideSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.OverrideSeverity = overrideSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AlertLevel) {
		alertLevel, err := client.NewEnumgaugeAlertLevelPropFromValue(plan.AlertLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.AlertLevel = alertLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateInterval) {
		addRequest.UpdateInterval = plan.UpdateInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SamplesPerUpdateInterval) {
		addRequest.SamplesPerUpdateInterval = plan.SamplesPerUpdateInterval.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeResource) {
		var slice []string
		plan.IncludeResource.ElementsAs(ctx, &slice, false)
		addRequest.IncludeResource = slice
	}
	if internaltypes.IsDefined(plan.ExcludeResource) {
		var slice []string
		plan.ExcludeResource.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeResource = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerUnavailableSeverityLevel) {
		serverUnavailableSeverityLevel, err := client.NewEnumgaugeServerUnavailableSeverityLevelPropFromValue(plan.ServerUnavailableSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerUnavailableSeverityLevel = serverUnavailableSeverityLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerDegradedSeverityLevel) {
		serverDegradedSeverityLevel, err := client.NewEnumgaugeServerDegradedSeverityLevelPropFromValue(plan.ServerDegradedSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerDegradedSeverityLevel = serverDegradedSeverityLevel
	}
	return nil
}

// Read a IndicatorGaugeResponse object into the model struct
func readIndicatorGaugeResponse(ctx context.Context, r *client.IndicatorGaugeResponse, state *gaugeResourceModel, expectedValues *gaugeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("indicator")
	state.Id = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	state.CriticalValue = internaltypes.StringTypeOrNil(r.CriticalValue, internaltypes.IsEmptyString(expectedValues.CriticalValue))
	state.MajorValue = internaltypes.StringTypeOrNil(r.MajorValue, internaltypes.IsEmptyString(expectedValues.MajorValue))
	state.MinorValue = internaltypes.StringTypeOrNil(r.MinorValue, internaltypes.IsEmptyString(expectedValues.MinorValue))
	state.WarningValue = internaltypes.StringTypeOrNil(r.WarningValue, internaltypes.IsEmptyString(expectedValues.WarningValue))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.OverrideSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeOverrideSeverityProp(r.OverrideSeverity), internaltypes.IsEmptyString(expectedValues.OverrideSeverity))
	state.AlertLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeAlertLevelProp(r.AlertLevel), internaltypes.IsEmptyString(expectedValues.AlertLevel))
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, internaltypes.IsEmptyString(expectedValues.UpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("update_interval",
		expectedValues.UpdateInterval, state.UpdateInterval, diagnostics)
	state.SamplesPerUpdateInterval = internaltypes.Int64TypeOrNil(r.SamplesPerUpdateInterval)
	state.IncludeResource = internaltypes.GetStringSet(r.IncludeResource)
	state.ExcludeResource = internaltypes.GetStringSet(r.ExcludeResource)
	state.ServerUnavailableSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerUnavailableSeverityLevelProp(r.ServerUnavailableSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerUnavailableSeverityLevel))
	state.ServerDegradedSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerDegradedSeverityLevelProp(r.ServerDegradedSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerDegradedSeverityLevel))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a NumericGaugeResponse object into the model struct
func readNumericGaugeResponse(ctx context.Context, r *client.NumericGaugeResponse, state *gaugeResourceModel, expectedValues *gaugeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("numeric")
	state.Id = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	if r.CriticalValue == nil {
		state.CriticalValue = types.StringNull()
	} else {
		state.CriticalValue = types.StringValue(strconv.FormatFloat(*r.CriticalValue, 'f', -1, 64))
	}
	state.CriticalExitValue = internaltypes.Float64TypeOrNil(r.CriticalExitValue)
	if r.MajorValue == nil {
		state.MajorValue = types.StringNull()
	} else {
		state.MajorValue = types.StringValue(strconv.FormatFloat(*r.MajorValue, 'f', -1, 64))
	}
	state.MajorExitValue = internaltypes.Float64TypeOrNil(r.MajorExitValue)
	if r.MinorValue == nil {
		state.MinorValue = types.StringNull()
	} else {
		state.MinorValue = types.StringValue(strconv.FormatFloat(*r.MinorValue, 'f', -1, 64))
	}
	state.MinorExitValue = internaltypes.Float64TypeOrNil(r.MinorExitValue)
	if r.WarningValue == nil {
		state.WarningValue = types.StringNull()
	} else {
		state.WarningValue = types.StringValue(strconv.FormatFloat(*r.WarningValue, 'f', -1, 64))
	}
	state.WarningExitValue = internaltypes.Float64TypeOrNil(r.WarningExitValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.OverrideSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeOverrideSeverityProp(r.OverrideSeverity), internaltypes.IsEmptyString(expectedValues.OverrideSeverity))
	state.AlertLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeAlertLevelProp(r.AlertLevel), internaltypes.IsEmptyString(expectedValues.AlertLevel))
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, internaltypes.IsEmptyString(expectedValues.UpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("update_interval",
		expectedValues.UpdateInterval, state.UpdateInterval, diagnostics)
	state.SamplesPerUpdateInterval = internaltypes.Int64TypeOrNil(r.SamplesPerUpdateInterval)
	state.IncludeResource = internaltypes.GetStringSet(r.IncludeResource)
	state.ExcludeResource = internaltypes.GetStringSet(r.ExcludeResource)
	state.ServerUnavailableSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerUnavailableSeverityLevelProp(r.ServerUnavailableSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerUnavailableSeverityLevel))
	state.ServerDegradedSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerDegradedSeverityLevelProp(r.ServerDegradedSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerDegradedSeverityLevel))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGaugeOperations(plan gaugeResourceModel, state gaugeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeDataSource, state.GaugeDataSource, "gauge-data-source")
	operations.AddStringOperationIfNecessary(&ops, plan.CriticalValue, state.CriticalValue, "critical-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.CriticalExitValue, state.CriticalExitValue, "critical-exit-value")
	operations.AddStringOperationIfNecessary(&ops, plan.MajorValue, state.MajorValue, "major-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MajorExitValue, state.MajorExitValue, "major-exit-value")
	operations.AddStringOperationIfNecessary(&ops, plan.MinorValue, state.MinorValue, "minor-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MinorExitValue, state.MinorExitValue, "minor-exit-value")
	operations.AddStringOperationIfNecessary(&ops, plan.WarningValue, state.WarningValue, "warning-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.WarningExitValue, state.WarningExitValue, "warning-exit-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.AlertLevel, state.AlertLevel, "alert-level")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateInterval, state.UpdateInterval, "update-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.SamplesPerUpdateInterval, state.SamplesPerUpdateInterval, "samples-per-update-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeResource, state.IncludeResource, "include-resource")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeResource, state.ExcludeResource, "exclude-resource")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerUnavailableSeverityLevel, state.ServerUnavailableSeverityLevel, "server-unavailable-severity-level")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerDegradedSeverityLevel, state.ServerDegradedSeverityLevel, "server-degraded-severity-level")
	return ops
}

// Create a indicator gauge
func (r *gaugeResource) CreateIndicatorGauge(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan gaugeResourceModel) (*gaugeResourceModel, error) {
	addRequest := client.NewAddIndicatorGaugeRequest(plan.Id.ValueString(),
		[]client.EnumindicatorGaugeSchemaUrn{client.ENUMINDICATORGAUGESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGEINDICATOR},
		plan.GaugeDataSource.ValueString())
	err := addOptionalIndicatorGaugeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Gauge", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.GaugeApi.AddGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGaugeRequest(
		client.AddIndicatorGaugeRequestAsAddGaugeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeApi.AddGaugeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Gauge", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state gaugeResourceModel
	readIndicatorGaugeResponse(ctx, addResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a numeric gauge
func (r *gaugeResource) CreateNumericGauge(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan gaugeResourceModel) (*gaugeResourceModel, error) {
	addRequest := client.NewAddNumericGaugeRequest(plan.Id.ValueString(),
		[]client.EnumnumericGaugeSchemaUrn{client.ENUMNUMERICGAUGESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGENUMERIC},
		plan.GaugeDataSource.ValueString())
	err := addOptionalNumericGaugeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Gauge", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.GaugeApi.AddGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGaugeRequest(
		client.AddNumericGaugeRequestAsAddGaugeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeApi.AddGaugeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Gauge", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state gaugeResourceModel
	readNumericGaugeResponse(ctx, addResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *gaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *gaugeResourceModel
	var err error
	if plan.Type.ValueString() == "indicator" {
		state, err = r.CreateIndicatorGauge(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "numeric" {
		state, err = r.CreateNumericGauge(ctx, req, resp, plan)
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
func (r *defaultGaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state gaugeResourceModel
	if plan.Type.ValueString() == "indicator" {
		readIndicatorGaugeResponse(ctx, readResponse.IndicatorGaugeResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "numeric" {
		readNumericGaugeResponse(ctx, readResponse.NumericGaugeResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GaugeApi.UpdateGauge(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "indicator" {
			readIndicatorGaugeResponse(ctx, updateResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "numeric" {
			readNumericGaugeResponse(ctx, updateResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *gaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGauge(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state gaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.IndicatorGaugeResponse != nil {
		readIndicatorGaugeResponse(ctx, readResponse.IndicatorGaugeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NumericGaugeResponse != nil {
		readNumericGaugeResponse(ctx, readResponse.NumericGaugeResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *gaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGauge(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan gaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state gaugeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.GaugeApi.UpdateGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "indicator" {
			readIndicatorGaugeResponse(ctx, updateResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "numeric" {
			readNumericGaugeResponse(ctx, updateResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *gaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state gaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GaugeApi.DeleteGaugeExecute(r.apiClient.GaugeApi.DeleteGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Gauge", err, httpResp)
		return
	}
}

func (r *gaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGauge(ctx, req, resp)
}

func (r *defaultGaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGauge(ctx, req, resp)
}

func importGauge(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
