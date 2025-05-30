// Copyright © 2025 Ping Identity Corporation

package monitorprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &monitorProviderResource{}
	_ resource.ResourceWithConfigure   = &monitorProviderResource{}
	_ resource.ResourceWithImportState = &monitorProviderResource{}
	_ resource.Resource                = &defaultMonitorProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultMonitorProviderResource{}
	_ resource.ResourceWithImportState = &defaultMonitorProviderResource{}
)

// Create a Monitor Provider resource
func NewMonitorProviderResource() resource.Resource {
	return &monitorProviderResource{}
}

func NewDefaultMonitorProviderResource() resource.Resource {
	return &defaultMonitorProviderResource{}
}

// monitorProviderResource is the resource implementation.
type monitorProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultMonitorProviderResource is the resource implementation.
type defaultMonitorProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *monitorProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_provider"
}

func (r *defaultMonitorProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_monitor_provider"
}

// Configure adds the provider configured client to the resource.
func (r *monitorProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultMonitorProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type monitorProviderResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	Type                    types.String `tfsdk:"type"`
	ExtensionClass          types.String `tfsdk:"extension_class"`
	ExtensionArgument       types.Set    `tfsdk:"extension_argument"`
	CheckFrequency          types.String `tfsdk:"check_frequency"`
	ProlongedOutageDuration types.String `tfsdk:"prolonged_outage_duration"`
	ProlongedOutageBehavior types.String `tfsdk:"prolonged_outage_behavior"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

type defaultMonitorProviderResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Notifications                        types.Set    `tfsdk:"notifications"`
	RequiredActions                      types.Set    `tfsdk:"required_actions"`
	Type                                 types.String `tfsdk:"type"`
	ExtensionClass                       types.String `tfsdk:"extension_class"`
	ExtensionArgument                    types.Set    `tfsdk:"extension_argument"`
	LowSpaceWarningSizeThreshold         types.String `tfsdk:"low_space_warning_size_threshold"`
	LowSpaceWarningPercentThreshold      types.Int64  `tfsdk:"low_space_warning_percent_threshold"`
	LowSpaceErrorSizeThreshold           types.String `tfsdk:"low_space_error_size_threshold"`
	LowSpaceErrorPercentThreshold        types.Int64  `tfsdk:"low_space_error_percent_threshold"`
	OutOfSpaceErrorSizeThreshold         types.String `tfsdk:"out_of_space_error_size_threshold"`
	OutOfSpaceErrorPercentThreshold      types.Int64  `tfsdk:"out_of_space_error_percent_threshold"`
	AlertFrequency                       types.String `tfsdk:"alert_frequency"`
	CheckFrequency                       types.String `tfsdk:"check_frequency"`
	DiskDevices                          types.Set    `tfsdk:"disk_devices"`
	NetworkDevices                       types.Set    `tfsdk:"network_devices"`
	SystemUtilizationMonitorLogDirectory types.String `tfsdk:"system_utilization_monitor_log_directory"`
	ProlongedOutageDuration              types.String `tfsdk:"prolonged_outage_duration"`
	ProlongedOutageBehavior              types.String `tfsdk:"prolonged_outage_behavior"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *monitorProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	monitorProviderSchema(ctx, req, resp, false)
}

func (r *defaultMonitorProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	monitorProviderSchema(ctx, req, resp, true)
}

func monitorProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Monitor Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Monitor Provider resource. Options are ['memory-usage', 'stack-trace', 'encryption-settings-database-accessibility', 'custom', 'active-operations', 'ssl-context', 'version', 'host-system', 'general', 'disk-space-usage', 'system-info', 'client-connection', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"encryption-settings-database-accessibility", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Monitor Provider.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Monitor Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"check_frequency": schema.StringAttribute{
				Description: "The frequency with which this monitor provider should confirm the ability to access the server's encryption settings database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prolonged_outage_duration": schema.StringAttribute{
				Description: "The minimum length of time that an outage should persist before it is considered a prolonged outage. If an outage lasts at least as long as this duration, then the server will take the action indicated by the prolonged-outage-behavior property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"prolonged_outage_behavior": schema.StringAttribute{
				Description: "The behavior that the server should exhibit after a prolonged period of time when the encryption settings database remains unreadable.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"none", "issue-alert", "enter-lockdown-mode", "shut-down-server"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Monitor Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`memory-usage`, `stack-trace`, `encryption-settings-database-accessibility`, `custom`, `active-operations`, `ssl-context`, `version`, `general`, `disk-space-usage`, `system-info`, `client-connection`, `third-party`]: Indicates whether the Monitor Provider is enabled for use. When the `type` attribute is set to `host-system`: Indicates whether the Host System Monitor Provider is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`memory-usage`, `stack-trace`, `encryption-settings-database-accessibility`, `custom`, `active-operations`, `ssl-context`, `version`, `general`, `disk-space-usage`, `system-info`, `client-connection`, `third-party`]: Indicates whether the Monitor Provider is enabled for use.\n  - `host-system`: Indicates whether the Host System Monitor Provider is enabled for use.",
				Required:            true,
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
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"memory-usage", "stack-trace", "encryption-settings-database-accessibility", "custom", "active-operations", "ssl-context", "version", "host-system", "general", "disk-space-usage", "system-info", "client-connection", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["low_space_warning_size_threshold"] = schema.StringAttribute{
			Description: "Specifies the low space warning threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
		}
		schemaDef.Attributes["low_space_warning_percent_threshold"] = schema.Int64Attribute{
			Description: "Specifies the low space warning threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
		}
		schemaDef.Attributes["low_space_error_size_threshold"] = schema.StringAttribute{
			Description: "Specifies the low space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
		}
		schemaDef.Attributes["low_space_error_percent_threshold"] = schema.Int64Attribute{
			Description: "Specifies the low space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
		}
		schemaDef.Attributes["out_of_space_error_size_threshold"] = schema.StringAttribute{
			Description: "Specifies the out of space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
		}
		schemaDef.Attributes["out_of_space_error_percent_threshold"] = schema.Int64Attribute{
			Description: "Specifies the out of space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
		}
		schemaDef.Attributes["alert_frequency"] = schema.StringAttribute{
			Description: "Specifies the length of time between administrative alerts generated in response to lack of usable disk space. Administrative alerts will be generated whenever the amount of usable space drops below any threshold, and they will also be generated at regular intervals as long as the amount of usable space remains below the threshold value. A value of zero indicates that alerts should only be generated when the amount of usable space drops below a configured threshold.",
		}
		schemaDef.Attributes["disk_devices"] = schema.SetAttribute{
			Description: "Specifies which disk devices to monitor for I/O activity. Should be the device name as displayed by iostat -d.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["network_devices"] = schema.SetAttribute{
			Description: "Specifies which network interfaces to monitor for I/O activity. Should be the device name as displayed by netstat -i.",
			ElementType: types.StringType,
		}
		schemaDef.Attributes["system_utilization_monitor_log_directory"] = schema.StringAttribute{
			Description: "Specifies a relative or absolute path to the directory on the local filesystem containing the log files used by the system utilization monitor. The path must exist, and it must be a writable directory by the server process.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *monitorProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel monitorProviderResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for encryption-settings-database-accessibility type
	if resourceType == "encryption-settings-database-accessibility" {
		if !internaltypes.IsDefined(configModel.ProlongedOutageBehavior) {
			defaultVal := types.StringValue("none")
			if !planModel.ProlongedOutageBehavior.Equal(defaultVal) {
				planModel.ProlongedOutageBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (model *monitorProviderResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "third-party" {
		model.ProlongedOutageDuration = types.StringNull()
		model.ProlongedOutageBehavior = types.StringNull()
		model.CheckFrequency = types.StringNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsMonitorProvider() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("check_frequency"),
			path.MatchRoot("type"),
			[]string{"encryption-settings-database-accessibility"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prolonged_outage_duration"),
			path.MatchRoot("type"),
			[]string{"encryption-settings-database-accessibility"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prolonged_outage_behavior"),
			path.MatchRoot("type"),
			[]string{"encryption-settings-database-accessibility"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"encryption-settings-database-accessibility",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r monitorProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsMonitorProvider()
}

// Add config validators
func (r defaultMonitorProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disk_devices"),
			path.MatchRoot("type"),
			[]string{"host-system"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("network_devices"),
			path.MatchRoot("type"),
			[]string{"host-system"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("system_utilization_monitor_log_directory"),
			path.MatchRoot("type"),
			[]string{"host-system"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_space_warning_size_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_space_warning_percent_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_space_error_size_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_space_error_percent_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("out_of_space_error_size_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("out_of_space_error_percent_threshold"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("alert_frequency"),
			path.MatchRoot("type"),
			[]string{"disk-space-usage"},
		),
	}
	return append(configValidatorsMonitorProvider(), validators...)
}

// Add optional fields to create request for encryption-settings-database-accessibility monitor-provider
func addOptionalEncryptionSettingsDatabaseAccessibilityMonitorProviderFields(ctx context.Context, addRequest *client.AddEncryptionSettingsDatabaseAccessibilityMonitorProviderRequest, plan monitorProviderResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CheckFrequency) {
		addRequest.CheckFrequency = plan.CheckFrequency.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ProlongedOutageDuration) {
		addRequest.ProlongedOutageDuration = plan.ProlongedOutageDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ProlongedOutageBehavior) {
		prolongedOutageBehavior, err := client.NewEnummonitorProviderProlongedOutageBehaviorPropFromValue(plan.ProlongedOutageBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.ProlongedOutageBehavior = prolongedOutageBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party monitor-provider
func addOptionalThirdPartyMonitorProviderFields(ctx context.Context, addRequest *client.AddThirdPartyMonitorProviderRequest, plan monitorProviderResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateMonitorProviderUnknownValues(model *monitorProviderResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateMonitorProviderUnknownValuesDefault(model *defaultMonitorProviderResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DiskDevices.IsUnknown() || model.DiskDevices.IsNull() {
		model.DiskDevices, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.NetworkDevices.IsUnknown() || model.NetworkDevices.IsNull() {
		model.NetworkDevices, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Read a MemoryUsageMonitorProviderResponse object into the model struct
func readMemoryUsageMonitorProviderResponseDefault(ctx context.Context, r *client.MemoryUsageMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("memory-usage")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a StackTraceMonitorProviderResponse object into the model struct
func readStackTraceMonitorProviderResponseDefault(ctx context.Context, r *client.StackTraceMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("stack-trace")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse object into the model struct
func readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponse(ctx context.Context, r *client.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, state *monitorProviderResourceModel, expectedValues *monitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings-database-accessibility")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CheckFrequency = types.StringValue(r.CheckFrequency)
	config.CheckMismatchedPDFormattedAttributes("check_frequency",
		expectedValues.CheckFrequency, state.CheckFrequency, diagnostics)
	state.ProlongedOutageDuration = internaltypes.StringTypeOrNil(r.ProlongedOutageDuration, true)
	config.CheckMismatchedPDFormattedAttributes("prolonged_outage_duration",
		expectedValues.ProlongedOutageDuration, state.ProlongedOutageDuration, diagnostics)
	state.ProlongedOutageBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnummonitorProviderProlongedOutageBehaviorProp(r.ProlongedOutageBehavior), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValues(state)
}

// Read a EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse object into the model struct
func readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDefault(ctx context.Context, r *client.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings-database-accessibility")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CheckFrequency = types.StringValue(r.CheckFrequency)
	config.CheckMismatchedPDFormattedAttributes("check_frequency",
		expectedValues.CheckFrequency, state.CheckFrequency, diagnostics)
	state.ProlongedOutageDuration = internaltypes.StringTypeOrNil(r.ProlongedOutageDuration, true)
	config.CheckMismatchedPDFormattedAttributes("prolonged_outage_duration",
		expectedValues.ProlongedOutageDuration, state.ProlongedOutageDuration, diagnostics)
	state.ProlongedOutageBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnummonitorProviderProlongedOutageBehaviorProp(r.ProlongedOutageBehavior), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a CustomMonitorProviderResponse object into the model struct
func readCustomMonitorProviderResponseDefault(ctx context.Context, r *client.CustomMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a ActiveOperationsMonitorProviderResponse object into the model struct
func readActiveOperationsMonitorProviderResponseDefault(ctx context.Context, r *client.ActiveOperationsMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("active-operations")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a SslContextMonitorProviderResponse object into the model struct
func readSslContextMonitorProviderResponseDefault(ctx context.Context, r *client.SslContextMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ssl-context")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a VersionMonitorProviderResponse object into the model struct
func readVersionMonitorProviderResponseDefault(ctx context.Context, r *client.VersionMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("version")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a HostSystemMonitorProviderResponse object into the model struct
func readHostSystemMonitorProviderResponseDefault(ctx context.Context, r *client.HostSystemMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("host-system")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.DiskDevices = internaltypes.GetStringSet(r.DiskDevices)
	state.NetworkDevices = internaltypes.GetStringSet(r.NetworkDevices)
	state.SystemUtilizationMonitorLogDirectory = types.StringValue(r.SystemUtilizationMonitorLogDirectory)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a GeneralMonitorProviderResponse object into the model struct
func readGeneralMonitorProviderResponseDefault(ctx context.Context, r *client.GeneralMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("general")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a DiskSpaceUsageMonitorProviderResponse object into the model struct
func readDiskSpaceUsageMonitorProviderResponseDefault(ctx context.Context, r *client.DiskSpaceUsageMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disk-space-usage")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LowSpaceWarningSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceWarningSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("low_space_warning_size_threshold",
		expectedValues.LowSpaceWarningSizeThreshold, state.LowSpaceWarningSizeThreshold, diagnostics)
	state.LowSpaceWarningPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceWarningPercentThreshold)
	state.LowSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceErrorSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("low_space_error_size_threshold",
		expectedValues.LowSpaceErrorSizeThreshold, state.LowSpaceErrorSizeThreshold, diagnostics)
	state.LowSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceErrorPercentThreshold)
	state.OutOfSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.OutOfSpaceErrorSizeThreshold, true)
	config.CheckMismatchedPDFormattedAttributes("out_of_space_error_size_threshold",
		expectedValues.OutOfSpaceErrorSizeThreshold, state.OutOfSpaceErrorSizeThreshold, diagnostics)
	state.OutOfSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.OutOfSpaceErrorPercentThreshold)
	state.AlertFrequency = types.StringValue(r.AlertFrequency)
	config.CheckMismatchedPDFormattedAttributes("alert_frequency",
		expectedValues.AlertFrequency, state.AlertFrequency, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a SystemInfoMonitorProviderResponse object into the model struct
func readSystemInfoMonitorProviderResponseDefault(ctx context.Context, r *client.SystemInfoMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("system-info")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a ClientConnectionMonitorProviderResponse object into the model struct
func readClientConnectionMonitorProviderResponseDefault(ctx context.Context, r *client.ClientConnectionMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-connection")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Read a ThirdPartyMonitorProviderResponse object into the model struct
func readThirdPartyMonitorProviderResponse(ctx context.Context, r *client.ThirdPartyMonitorProviderResponse, state *monitorProviderResourceModel, expectedValues *monitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValues(state)
}

// Read a ThirdPartyMonitorProviderResponse object into the model struct
func readThirdPartyMonitorProviderResponseDefault(ctx context.Context, r *client.ThirdPartyMonitorProviderResponse, state *defaultMonitorProviderResourceModel, expectedValues *defaultMonitorProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateMonitorProviderUnknownValuesDefault(state)
}

// Create any update operations necessary to make the state match the plan
func createMonitorProviderOperations(plan monitorProviderResourceModel, state monitorProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.CheckFrequency, state.CheckFrequency, "check-frequency")
	operations.AddStringOperationIfNecessary(&ops, plan.ProlongedOutageDuration, state.ProlongedOutageDuration, "prolonged-outage-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.ProlongedOutageBehavior, state.ProlongedOutageBehavior, "prolonged-outage-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createMonitorProviderOperationsDefault(plan defaultMonitorProviderResourceModel, state defaultMonitorProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.LowSpaceWarningSizeThreshold, state.LowSpaceWarningSizeThreshold, "low-space-warning-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowSpaceWarningPercentThreshold, state.LowSpaceWarningPercentThreshold, "low-space-warning-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.LowSpaceErrorSizeThreshold, state.LowSpaceErrorSizeThreshold, "low-space-error-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowSpaceErrorPercentThreshold, state.LowSpaceErrorPercentThreshold, "low-space-error-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.OutOfSpaceErrorSizeThreshold, state.OutOfSpaceErrorSizeThreshold, "out-of-space-error-size-threshold")
	operations.AddInt64OperationIfNecessary(&ops, plan.OutOfSpaceErrorPercentThreshold, state.OutOfSpaceErrorPercentThreshold, "out-of-space-error-percent-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.AlertFrequency, state.AlertFrequency, "alert-frequency")
	operations.AddStringOperationIfNecessary(&ops, plan.CheckFrequency, state.CheckFrequency, "check-frequency")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DiskDevices, state.DiskDevices, "disk-devices")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.NetworkDevices, state.NetworkDevices, "network-devices")
	operations.AddStringOperationIfNecessary(&ops, plan.SystemUtilizationMonitorLogDirectory, state.SystemUtilizationMonitorLogDirectory, "system-utilization-monitor-log-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.ProlongedOutageDuration, state.ProlongedOutageDuration, "prolonged-outage-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.ProlongedOutageBehavior, state.ProlongedOutageBehavior, "prolonged-outage-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a encryption-settings-database-accessibility monitor-provider
func (r *monitorProviderResource) CreateEncryptionSettingsDatabaseAccessibilityMonitorProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan monitorProviderResourceModel) (*monitorProviderResourceModel, error) {
	addRequest := client.NewAddEncryptionSettingsDatabaseAccessibilityMonitorProviderRequest([]client.EnumencryptionSettingsDatabaseAccessibilityMonitorProviderSchemaUrn{client.ENUMENCRYPTIONSETTINGSDATABASEACCESSIBILITYMONITORPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0MONITOR_PROVIDERENCRYPTION_SETTINGS_DATABASE_ACCESSIBILITY},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalEncryptionSettingsDatabaseAccessibilityMonitorProviderFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Monitor Provider", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.MonitorProviderAPI.AddMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddMonitorProviderRequest(
		client.AddEncryptionSettingsDatabaseAccessibilityMonitorProviderRequestAsAddMonitorProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.MonitorProviderAPI.AddMonitorProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Monitor Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state monitorProviderResourceModel
	readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponse(ctx, addResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party monitor-provider
func (r *monitorProviderResource) CreateThirdPartyMonitorProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan monitorProviderResourceModel) (*monitorProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyMonitorProviderRequest([]client.EnumthirdPartyMonitorProviderSchemaUrn{client.ENUMTHIRDPARTYMONITORPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0MONITOR_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyMonitorProviderFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Monitor Provider", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.MonitorProviderAPI.AddMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddMonitorProviderRequest(
		client.AddThirdPartyMonitorProviderRequestAsAddMonitorProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.MonitorProviderAPI.AddMonitorProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Monitor Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state monitorProviderResourceModel
	readThirdPartyMonitorProviderResponse(ctx, addResponse.ThirdPartyMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *monitorProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan monitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *monitorProviderResourceModel
	var err error
	if plan.Type.ValueString() == "encryption-settings-database-accessibility" {
		state, err = r.CreateEncryptionSettingsDatabaseAccessibilityMonitorProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyMonitorProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

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
func (r *defaultMonitorProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultMonitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderAPI.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitor Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultMonitorProviderResourceModel
	if readResponse.MemoryUsageMonitorProviderResponse != nil {
		readMemoryUsageMonitorProviderResponseDefault(ctx, readResponse.MemoryUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StackTraceMonitorProviderResponse != nil {
		readStackTraceMonitorProviderResponseDefault(ctx, readResponse.StackTraceMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
		readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDefault(ctx, readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomMonitorProviderResponse != nil {
		readCustomMonitorProviderResponseDefault(ctx, readResponse.CustomMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveOperationsMonitorProviderResponse != nil {
		readActiveOperationsMonitorProviderResponseDefault(ctx, readResponse.ActiveOperationsMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SslContextMonitorProviderResponse != nil {
		readSslContextMonitorProviderResponseDefault(ctx, readResponse.SslContextMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VersionMonitorProviderResponse != nil {
		readVersionMonitorProviderResponseDefault(ctx, readResponse.VersionMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HostSystemMonitorProviderResponse != nil {
		readHostSystemMonitorProviderResponseDefault(ctx, readResponse.HostSystemMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralMonitorProviderResponse != nil {
		readGeneralMonitorProviderResponseDefault(ctx, readResponse.GeneralMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DiskSpaceUsageMonitorProviderResponse != nil {
		readDiskSpaceUsageMonitorProviderResponseDefault(ctx, readResponse.DiskSpaceUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SystemInfoMonitorProviderResponse != nil {
		readSystemInfoMonitorProviderResponseDefault(ctx, readResponse.SystemInfoMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClientConnectionMonitorProviderResponse != nil {
		readClientConnectionMonitorProviderResponseDefault(ctx, readResponse.ClientConnectionMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyMonitorProviderResponse != nil {
		readThirdPartyMonitorProviderResponseDefault(ctx, readResponse.ThirdPartyMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.MonitorProviderAPI.UpdateMonitorProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createMonitorProviderOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitorProviderAPI.UpdateMonitorProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Monitor Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.MemoryUsageMonitorProviderResponse != nil {
			readMemoryUsageMonitorProviderResponseDefault(ctx, updateResponse.MemoryUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StackTraceMonitorProviderResponse != nil {
			readStackTraceMonitorProviderResponseDefault(ctx, updateResponse.StackTraceMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
			readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDefault(ctx, updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomMonitorProviderResponse != nil {
			readCustomMonitorProviderResponseDefault(ctx, updateResponse.CustomMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ActiveOperationsMonitorProviderResponse != nil {
			readActiveOperationsMonitorProviderResponseDefault(ctx, updateResponse.ActiveOperationsMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SslContextMonitorProviderResponse != nil {
			readSslContextMonitorProviderResponseDefault(ctx, updateResponse.SslContextMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VersionMonitorProviderResponse != nil {
			readVersionMonitorProviderResponseDefault(ctx, updateResponse.VersionMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HostSystemMonitorProviderResponse != nil {
			readHostSystemMonitorProviderResponseDefault(ctx, updateResponse.HostSystemMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GeneralMonitorProviderResponse != nil {
			readGeneralMonitorProviderResponseDefault(ctx, updateResponse.GeneralMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DiskSpaceUsageMonitorProviderResponse != nil {
			readDiskSpaceUsageMonitorProviderResponseDefault(ctx, updateResponse.DiskSpaceUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SystemInfoMonitorProviderResponse != nil {
			readSystemInfoMonitorProviderResponseDefault(ctx, updateResponse.SystemInfoMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClientConnectionMonitorProviderResponse != nil {
			readClientConnectionMonitorProviderResponseDefault(ctx, updateResponse.ClientConnectionMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyMonitorProviderResponse != nil {
			readThirdPartyMonitorProviderResponseDefault(ctx, updateResponse.ThirdPartyMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *monitorProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state monitorProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderAPI.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Monitor Provider", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitor Provider", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
		readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponse(ctx, readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyMonitorProviderResponse != nil {
		readThirdPartyMonitorProviderResponse(ctx, readResponse.ThirdPartyMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultMonitorProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultMonitorProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderAPI.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitor Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.MemoryUsageMonitorProviderResponse != nil {
		readMemoryUsageMonitorProviderResponseDefault(ctx, readResponse.MemoryUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StackTraceMonitorProviderResponse != nil {
		readStackTraceMonitorProviderResponseDefault(ctx, readResponse.StackTraceMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomMonitorProviderResponse != nil {
		readCustomMonitorProviderResponseDefault(ctx, readResponse.CustomMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveOperationsMonitorProviderResponse != nil {
		readActiveOperationsMonitorProviderResponseDefault(ctx, readResponse.ActiveOperationsMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SslContextMonitorProviderResponse != nil {
		readSslContextMonitorProviderResponseDefault(ctx, readResponse.SslContextMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VersionMonitorProviderResponse != nil {
		readVersionMonitorProviderResponseDefault(ctx, readResponse.VersionMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HostSystemMonitorProviderResponse != nil {
		readHostSystemMonitorProviderResponseDefault(ctx, readResponse.HostSystemMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralMonitorProviderResponse != nil {
		readGeneralMonitorProviderResponseDefault(ctx, readResponse.GeneralMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DiskSpaceUsageMonitorProviderResponse != nil {
		readDiskSpaceUsageMonitorProviderResponseDefault(ctx, readResponse.DiskSpaceUsageMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SystemInfoMonitorProviderResponse != nil {
		readSystemInfoMonitorProviderResponseDefault(ctx, readResponse.SystemInfoMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClientConnectionMonitorProviderResponse != nil {
		readClientConnectionMonitorProviderResponseDefault(ctx, readResponse.ClientConnectionMonitorProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *monitorProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan monitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state monitorProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.MonitorProviderAPI.UpdateMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createMonitorProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitorProviderAPI.UpdateMonitorProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Monitor Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
			readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponse(ctx, updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyMonitorProviderResponse != nil {
			readThirdPartyMonitorProviderResponse(ctx, updateResponse.ThirdPartyMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultMonitorProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultMonitorProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultMonitorProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.MonitorProviderAPI.UpdateMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createMonitorProviderOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.MonitorProviderAPI.UpdateMonitorProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Monitor Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.MemoryUsageMonitorProviderResponse != nil {
			readMemoryUsageMonitorProviderResponseDefault(ctx, updateResponse.MemoryUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StackTraceMonitorProviderResponse != nil {
			readStackTraceMonitorProviderResponseDefault(ctx, updateResponse.StackTraceMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
			readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDefault(ctx, updateResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomMonitorProviderResponse != nil {
			readCustomMonitorProviderResponseDefault(ctx, updateResponse.CustomMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ActiveOperationsMonitorProviderResponse != nil {
			readActiveOperationsMonitorProviderResponseDefault(ctx, updateResponse.ActiveOperationsMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SslContextMonitorProviderResponse != nil {
			readSslContextMonitorProviderResponseDefault(ctx, updateResponse.SslContextMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VersionMonitorProviderResponse != nil {
			readVersionMonitorProviderResponseDefault(ctx, updateResponse.VersionMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HostSystemMonitorProviderResponse != nil {
			readHostSystemMonitorProviderResponseDefault(ctx, updateResponse.HostSystemMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GeneralMonitorProviderResponse != nil {
			readGeneralMonitorProviderResponseDefault(ctx, updateResponse.GeneralMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DiskSpaceUsageMonitorProviderResponse != nil {
			readDiskSpaceUsageMonitorProviderResponseDefault(ctx, updateResponse.DiskSpaceUsageMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SystemInfoMonitorProviderResponse != nil {
			readSystemInfoMonitorProviderResponseDefault(ctx, updateResponse.SystemInfoMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClientConnectionMonitorProviderResponse != nil {
			readClientConnectionMonitorProviderResponseDefault(ctx, updateResponse.ClientConnectionMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyMonitorProviderResponse != nil {
			readThirdPartyMonitorProviderResponseDefault(ctx, updateResponse.ThirdPartyMonitorProviderResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultMonitorProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *monitorProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state monitorProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.MonitorProviderAPI.DeleteMonitorProviderExecute(r.apiClient.MonitorProviderAPI.DeleteMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Monitor Provider", err, httpResp)
		return
	}
}

func (r *monitorProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMonitorProvider(ctx, req, resp)
}

func (r *defaultMonitorProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMonitorProvider(ctx, req, resp)
}

func importMonitorProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
