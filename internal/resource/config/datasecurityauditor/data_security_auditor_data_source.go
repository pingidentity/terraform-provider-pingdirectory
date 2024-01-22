package datasecurityauditor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &dataSecurityAuditorDataSource{}
	_ datasource.DataSourceWithConfigure = &dataSecurityAuditorDataSource{}
)

// Create a Data Security Auditor data source
func NewDataSecurityAuditorDataSource() datasource.DataSource {
	return &dataSecurityAuditorDataSource{}
}

// dataSecurityAuditorDataSource is the datasource implementation.
type dataSecurityAuditorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *dataSecurityAuditorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_security_auditor"
}

// Configure adds the provider configured client to the data source.
func (r *dataSecurityAuditorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type dataSecurityAuditorDataSourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	Type                                types.String `tfsdk:"type"`
	ExtensionClass                      types.String `tfsdk:"extension_class"`
	ExtensionArgument                   types.Set    `tfsdk:"extension_argument"`
	ReportFile                          types.String `tfsdk:"report_file"`
	Filter                              types.Set    `tfsdk:"filter"`
	AccountExpirationWarningInterval    types.String `tfsdk:"account_expiration_warning_interval"`
	IncludePrivilege                    types.Set    `tfsdk:"include_privilege"`
	MaximumIdleTime                     types.String `tfsdk:"maximum_idle_time"`
	WeakPasswordStorageScheme           types.Set    `tfsdk:"weak_password_storage_scheme"`
	WeakCryptEncoding                   types.Set    `tfsdk:"weak_crypt_encoding"`
	IdleAccountWarningInterval          types.String `tfsdk:"idle_account_warning_interval"`
	IdleAccountErrorInterval            types.String `tfsdk:"idle_account_error_interval"`
	NeverLoggedInAccountWarningInterval types.String `tfsdk:"never_logged_in_account_warning_interval"`
	NeverLoggedInAccountErrorInterval   types.String `tfsdk:"never_logged_in_account_error_interval"`
	IncludeAttribute                    types.Set    `tfsdk:"include_attribute"`
	PasswordEvaluationAge               types.String `tfsdk:"password_evaluation_age"`
	Enabled                             types.Bool   `tfsdk:"enabled"`
	AuditBackend                        types.Set    `tfsdk:"audit_backend"`
	AuditSeverity                       types.String `tfsdk:"audit_severity"`
}

// GetSchema defines the schema for the datasource.
func (r *dataSecurityAuditorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Data Security Auditor.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Data Security Auditor resource. Options are ['expired-password', 'idle-account', 'disabled-account', 'weakly-encoded-password', 'privilege', 'account-usability-issues', 'locked-account', 'filter', 'account-validity-window', 'multiple-password', 'deprecated-password-storage-scheme', 'nonexistent-password-policy', 'access-control', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Data Security Auditor.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Data Security Auditor. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.SetAttribute{
				Description: "The filter to use to identify entries that should be reported. Multiple filters may be configured, and each reported entry will indicate which of these filter(s) matched that entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"account_expiration_warning_interval": schema.StringAttribute{
				Description: "If set, the auditor will report all users with account expiration times are in the future, but are within the specified length of time away from the current time.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_privilege": schema.SetAttribute{
				Description: "If defined, only entries with the specified privileges will be reported. By default, entries with any privilege assigned will be reported.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maximum_idle_time": schema.StringAttribute{
				Description: "If set, users that have not authenticated for more than the specified time will be reported even if idle account lockout is not configured. Note that users may only be reported if the last login time tracking is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"weak_password_storage_scheme": schema.SetAttribute{
				Description: "The password storage schemes that are considered weak. Users with any of the specified password storage schemes will be included in the report.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"weak_crypt_encoding": schema.SetAttribute{
				Description: "Reporting on users with passwords encoded using the Crypt Password Storage scheme may be further limited by selecting one or more encoding mechanisms that are considered weak.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"idle_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for idle accounts. If the length of time since a user last authenticated is greater than the warning interval but less than the error interval (or if it is greater than the warning interval and no error interval is defined), then a warning will be generated for that account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"idle_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for idle accounts. If the length of time since a user last authenticated is greater than the error interval, then an error will be generated for that account. If no error interval is defined, then only the warning interval will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"never_logged_in_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for accounts that do not appear to have authenticated. If this is not specified, then the idle account warning interval will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"never_logged_in_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for accounts that do not appear to have authenticated. If this is not specified, then the never-logged-in warning interval will be used. The idle account warning and error intervals will be used if no never-logged-in interval is configured.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_attribute": schema.SetAttribute{
				Description: "Specifies the attributes from the audited entries that should be included detailed reports. By default, no attributes are included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"password_evaluation_age": schema.StringAttribute{
				Description: "If set, the auditor will report all users with passwords older than the specified value even if password expiration is not enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Data Security Auditor is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"audit_backend": schema.SetAttribute{
				Description: "Specifies which backends the data security auditor may be applied to. By default, the data security auditors will audit entries in all backend types that support data auditing (Local DB, LDIF, and Config File Handler).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"audit_severity": schema.StringAttribute{
				Description: "Specifies the severity of events to include in the report.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ExpiredPasswordDataSecurityAuditorResponse object into the model struct
func readExpiredPasswordDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.ExpiredPasswordDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("expired-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.PasswordEvaluationAge = internaltypes.StringTypeOrNil(r.PasswordEvaluationAge, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a IdleAccountDataSecurityAuditorResponse object into the model struct
func readIdleAccountDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.IdleAccountDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("idle-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IdleAccountWarningInterval = types.StringValue(r.IdleAccountWarningInterval)
	state.IdleAccountErrorInterval = internaltypes.StringTypeOrNil(r.IdleAccountErrorInterval, false)
	state.NeverLoggedInAccountWarningInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountWarningInterval, false)
	state.NeverLoggedInAccountErrorInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountErrorInterval, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a DisabledAccountDataSecurityAuditorResponse object into the model struct
func readDisabledAccountDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.DisabledAccountDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disabled-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a WeaklyEncodedPasswordDataSecurityAuditorResponse object into the model struct
func readWeaklyEncodedPasswordDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.WeaklyEncodedPasswordDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("weakly-encoded-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.WeakPasswordStorageScheme = internaltypes.GetStringSet(r.WeakPasswordStorageScheme)
	state.WeakCryptEncoding = internaltypes.GetStringSet(
		client.StringSliceEnumdataSecurityAuditorWeakCryptEncodingProp(r.WeakCryptEncoding))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a PrivilegeDataSecurityAuditorResponse object into the model struct
func readPrivilegeDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.PrivilegeDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("privilege")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludePrivilege = internaltypes.GetStringSet(
		client.StringSliceEnumdataSecurityAuditorIncludePrivilegeProp(r.IncludePrivilege))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a AccountUsabilityIssuesDataSecurityAuditorResponse object into the model struct
func readAccountUsabilityIssuesDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.AccountUsabilityIssuesDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("account-usability-issues")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a LockedAccountDataSecurityAuditorResponse object into the model struct
func readLockedAccountDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.LockedAccountDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("locked-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.MaximumIdleTime = internaltypes.StringTypeOrNil(r.MaximumIdleTime, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a FilterDataSecurityAuditorResponse object into the model struct
func readFilterDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.FilterDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("filter")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a AccountValidityWindowDataSecurityAuditorResponse object into the model struct
func readAccountValidityWindowDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.AccountValidityWindowDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("account-validity-window")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AccountExpirationWarningInterval = internaltypes.StringTypeOrNil(r.AccountExpirationWarningInterval, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a MultiplePasswordDataSecurityAuditorResponse object into the model struct
func readMultiplePasswordDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.MultiplePasswordDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multiple-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse object into the model struct
func readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deprecated-password-storage-scheme")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a NonexistentPasswordPolicyDataSecurityAuditorResponse object into the model struct
func readNonexistentPasswordPolicyDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.NonexistentPasswordPolicyDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nonexistent-password-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a AccessControlDataSecurityAuditorResponse object into the model struct
func readAccessControlDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.AccessControlDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("access-control")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read a ThirdPartyDataSecurityAuditorResponse object into the model struct
func readThirdPartyDataSecurityAuditorResponseDataSource(ctx context.Context, r *client.ThirdPartyDataSecurityAuditorResponse, state *dataSecurityAuditorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), false)
}

// Read resource information
func (r *dataSecurityAuditorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state dataSecurityAuditorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ExpiredPasswordDataSecurityAuditorResponse != nil {
		readExpiredPasswordDataSecurityAuditorResponseDataSource(ctx, readResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IdleAccountDataSecurityAuditorResponse != nil {
		readIdleAccountDataSecurityAuditorResponseDataSource(ctx, readResponse.IdleAccountDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DisabledAccountDataSecurityAuditorResponse != nil {
		readDisabledAccountDataSecurityAuditorResponseDataSource(ctx, readResponse.DisabledAccountDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
		readWeaklyEncodedPasswordDataSecurityAuditorResponseDataSource(ctx, readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PrivilegeDataSecurityAuditorResponse != nil {
		readPrivilegeDataSecurityAuditorResponseDataSource(ctx, readResponse.PrivilegeDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
		readAccountUsabilityIssuesDataSecurityAuditorResponseDataSource(ctx, readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LockedAccountDataSecurityAuditorResponse != nil {
		readLockedAccountDataSecurityAuditorResponseDataSource(ctx, readResponse.LockedAccountDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FilterDataSecurityAuditorResponse != nil {
		readFilterDataSecurityAuditorResponseDataSource(ctx, readResponse.FilterDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AccountValidityWindowDataSecurityAuditorResponse != nil {
		readAccountValidityWindowDataSecurityAuditorResponseDataSource(ctx, readResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MultiplePasswordDataSecurityAuditorResponse != nil {
		readMultiplePasswordDataSecurityAuditorResponseDataSource(ctx, readResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
		readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponseDataSource(ctx, readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
		readNonexistentPasswordPolicyDataSecurityAuditorResponseDataSource(ctx, readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AccessControlDataSecurityAuditorResponse != nil {
		readAccessControlDataSecurityAuditorResponseDataSource(ctx, readResponse.AccessControlDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyDataSecurityAuditorResponse != nil {
		readThirdPartyDataSecurityAuditorResponseDataSource(ctx, readResponse.ThirdPartyDataSecurityAuditorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
