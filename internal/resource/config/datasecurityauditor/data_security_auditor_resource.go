// Copyright © 2025 Ping Identity Corporation

package datasecurityauditor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	_ resource.Resource                = &dataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &dataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &dataSecurityAuditorResource{}
	_ resource.Resource                = &defaultDataSecurityAuditorResource{}
	_ resource.ResourceWithConfigure   = &defaultDataSecurityAuditorResource{}
	_ resource.ResourceWithImportState = &defaultDataSecurityAuditorResource{}
)

// Create a Data Security Auditor resource
func NewDataSecurityAuditorResource() resource.Resource {
	return &dataSecurityAuditorResource{}
}

func NewDefaultDataSecurityAuditorResource() resource.Resource {
	return &defaultDataSecurityAuditorResource{}
}

// dataSecurityAuditorResource is the resource implementation.
type dataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDataSecurityAuditorResource is the resource implementation.
type defaultDataSecurityAuditorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *dataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_security_auditor"
}

func (r *defaultDataSecurityAuditorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_data_security_auditor"
}

// Configure adds the provider configured client to the resource.
func (r *dataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDataSecurityAuditorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type dataSecurityAuditorResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *dataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	dataSecurityAuditorSchema(ctx, req, resp, false)
}

func (r *defaultDataSecurityAuditorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	dataSecurityAuditorSchema(ctx, req, resp, true)
}

func dataSecurityAuditorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Data Security Auditor.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Data Security Auditor resource. Options are ['expired-password', 'idle-account', 'disabled-account', 'weakly-encoded-password', 'privilege', 'account-usability-issues', 'locked-account', 'filter', 'account-validity-window', 'multiple-password', 'deprecated-password-storage-scheme', 'nonexistent-password-policy', 'access-control', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"expired-password", "idle-account", "disabled-account", "weakly-encoded-password", "privilege", "account-usability-issues", "locked-account", "filter", "account-validity-window", "multiple-password", "deprecated-password-storage-scheme", "nonexistent-password-policy", "access-control", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Data Security Auditor.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Data Security Auditor. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"report_file": schema.StringAttribute{
				Description: "Specifies the name of the detailed report file.",
				Optional:    true,
				Computed:    true,
			},
			"filter": schema.SetAttribute{
				Description: "The filter to use to identify entries that should be reported. Multiple filters may be configured, and each reported entry will indicate which of these filter(s) matched that entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"account_expiration_warning_interval": schema.StringAttribute{
				Description: "If set, the auditor will report all users with account expiration times are in the future, but are within the specified length of time away from the current time.",
				Optional:    true,
			},
			"include_privilege": schema.SetAttribute{
				Description: "If defined, only entries with the specified privileges will be reported. By default, entries with any privilege assigned will be reported.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"maximum_idle_time": schema.StringAttribute{
				Description: "If set, users that have not authenticated for more than the specified time will be reported even if idle account lockout is not configured. Note that users may only be reported if the last login time tracking is enabled.",
				Optional:    true,
			},
			"weak_password_storage_scheme": schema.SetAttribute{
				Description: "The password storage schemes that are considered weak. Users with any of the specified password storage schemes will be included in the report.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"weak_crypt_encoding": schema.SetAttribute{
				Description: "Reporting on users with passwords encoded using the Crypt Password Storage scheme may be further limited by selecting one or more encoding mechanisms that are considered weak.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"idle_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for idle accounts. If the length of time since a user last authenticated is greater than the warning interval but less than the error interval (or if it is greater than the warning interval and no error interval is defined), then a warning will be generated for that account.",
				Optional:    true,
			},
			"idle_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for idle accounts. If the length of time since a user last authenticated is greater than the error interval, then an error will be generated for that account. If no error interval is defined, then only the warning interval will be used.",
				Optional:    true,
			},
			"never_logged_in_account_warning_interval": schema.StringAttribute{
				Description: "The length of time to use as the warning interval for accounts that do not appear to have authenticated. If this is not specified, then the idle account warning interval will be used.",
				Optional:    true,
			},
			"never_logged_in_account_error_interval": schema.StringAttribute{
				Description: "The length of time to use as the error interval for accounts that do not appear to have authenticated. If this is not specified, then the never-logged-in warning interval will be used. The idle account warning and error intervals will be used if no never-logged-in interval is configured.",
				Optional:    true,
			},
			"include_attribute": schema.SetAttribute{
				Description: "Specifies the attributes from the audited entries that should be included detailed reports. By default, no attributes are included.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"password_evaluation_age": schema.StringAttribute{
				Description: "If set, the auditor will report all users with passwords older than the specified value even if password expiration is not enabled.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Data Security Auditor is enabled for use.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"audit_backend": schema.SetAttribute{
				Description: "Specifies which backends the data security auditor may be applied to. By default, the data security auditors will audit entries in all backend types that support data auditing (Local DB, LDIF, and Config File Handler).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"audit_severity": schema.StringAttribute{
				Description: "Specifies the severity of events to include in the report.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("notice"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"error", "warning", "notice", "verbose"}...),
				},
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
func (r *dataSecurityAuditorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel dataSecurityAuditorResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for expired-password type
	if resourceType == "expired-password" {
		if !internaltypes.IsDefined(configModel.IncludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("pwdChangedTime"), types.StringValue("ds-pwp-warned-time"), types.StringValue("ds-pwp-last-login-time")})
			if !planModel.IncludeAttribute.Equal(defaultVal) {
				planModel.IncludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for idle-account type
	if resourceType == "idle-account" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("idle-accounts.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for disabled-account type
	if resourceType == "disabled-account" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("entries-with-disabled-accounts.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for weakly-encoded-password type
	if resourceType == "weakly-encoded-password" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("users-with-weakly-encoded-passwords.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WeakPasswordStorageScheme) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("3DES"), types.StringValue("AES"), types.StringValue("Base64"), types.StringValue("Blowfish"), types.StringValue("Clear"), types.StringValue("CRYPT"), types.StringValue("MD5"), types.StringValue("RC4"), types.StringValue("SHA-1"), types.StringValue("Salted MD5"), types.StringValue("Salted SHA-1")})
			if !planModel.WeakPasswordStorageScheme.Equal(defaultVal) {
				planModel.WeakPasswordStorageScheme = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WeakCryptEncoding) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("crypt"), types.StringValue("md5")})
			if !planModel.WeakCryptEncoding.Equal(defaultVal) {
				planModel.WeakCryptEncoding = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for privilege type
	if resourceType == "privilege" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("entries-with-privileges.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for account-usability-issues type
	if resourceType == "account-usability-issues" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("account-usability-issues.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for locked-account type
	if resourceType == "locked-account" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("locked-accounts.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("pwdAccountLockedTime"), types.StringValue("ds-pwp-auth-failure"), types.StringValue("ds-pwp-last-login-time")})
			if !planModel.IncludeAttribute.Equal(defaultVal) {
				planModel.IncludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for account-validity-window type
	if resourceType == "account-validity-window" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("account-validity-window.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-pwp-account-activation-time"), types.StringValue("ds-pwp-account-expiration-time")})
			if !planModel.IncludeAttribute.Equal(defaultVal) {
				planModel.IncludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for multiple-password type
	if resourceType == "multiple-password" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("users-with-multiple-passwords.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for deprecated-password-storage-scheme type
	if resourceType == "deprecated-password-storage-scheme" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("deprecated-password-storage-schemes.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for nonexistent-password-policy type
	if resourceType == "nonexistent-password-policy" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("nonexistent-password-policies.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-pwp-password-policy-dn")})
			if !planModel.IncludeAttribute.Equal(defaultVal) {
				planModel.IncludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for access-control type
	if resourceType == "access-control" {
		if !internaltypes.IsDefined(configModel.ReportFile) {
			defaultVal := types.StringValue("entries-with-acis.ldif")
			if !planModel.ReportFile.Equal(defaultVal) {
				planModel.ReportFile = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("aci")})
			if !planModel.IncludeAttribute.Equal(defaultVal) {
				planModel.IncludeAttribute = defaultVal
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

func (model *dataSecurityAuditorResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "expired-password" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "idle-account" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "disabled-account" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "weakly-encoded-password" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "privilege" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "account-usability-issues" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "locked-account" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "filter" {
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "account-validity-window" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "multiple-password" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "deprecated-password-storage-scheme" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "nonexistent-password-policy" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "access-control" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "third-party" {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsDataSecurityAuditor() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"idle-account"},
			configvalidators.Implies(
				path.MatchRoot("never_logged_in_account_error_interval"),
				path.MatchRoot("never_logged_in_account_warning_interval"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password_evaluation_age"),
			path.MatchRoot("type"),
			[]string{"expired-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("idle_account_warning_interval"),
			path.MatchRoot("type"),
			[]string{"idle-account"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("idle_account_error_interval"),
			path.MatchRoot("type"),
			[]string{"idle-account"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("never_logged_in_account_warning_interval"),
			path.MatchRoot("type"),
			[]string{"idle-account"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("never_logged_in_account_error_interval"),
			path.MatchRoot("type"),
			[]string{"idle-account"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("weak_password_storage_scheme"),
			path.MatchRoot("type"),
			[]string{"weakly-encoded-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("weak_crypt_encoding"),
			path.MatchRoot("type"),
			[]string{"weakly-encoded-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_privilege"),
			path.MatchRoot("type"),
			[]string{"privilege"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_idle_time"),
			path.MatchRoot("type"),
			[]string{"locked-account"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("filter"),
			path.MatchRoot("type"),
			[]string{"filter"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("account_expiration_warning_interval"),
			path.MatchRoot("type"),
			[]string{"account-validity-window"},
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
			"filter",
			[]path.Expression{path.MatchRoot("report_file"), path.MatchRoot("filter")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("report_file"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"idle-account",
			[]path.Expression{path.MatchRoot("idle_account_warning_interval")},
		),
	}
}

// Add config validators
func (r dataSecurityAuditorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsDataSecurityAuditor()
}

// Add config validators
func (r defaultDataSecurityAuditorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsDataSecurityAuditor()
}

// Add optional fields to create request for expired-password data-security-auditor
func addOptionalExpiredPasswordDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddExpiredPasswordDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordEvaluationAge) {
		addRequest.PasswordEvaluationAge = plan.PasswordEvaluationAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for idle-account data-security-auditor
func addOptionalIdleAccountDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddIdleAccountDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdleAccountErrorInterval) {
		addRequest.IdleAccountErrorInterval = plan.IdleAccountErrorInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NeverLoggedInAccountWarningInterval) {
		addRequest.NeverLoggedInAccountWarningInterval = plan.NeverLoggedInAccountWarningInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NeverLoggedInAccountErrorInterval) {
		addRequest.NeverLoggedInAccountErrorInterval = plan.NeverLoggedInAccountErrorInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for disabled-account data-security-auditor
func addOptionalDisabledAccountDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddDisabledAccountDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for weakly-encoded-password data-security-auditor
func addOptionalWeaklyEncodedPasswordDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddWeaklyEncodedPasswordDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.WeakPasswordStorageScheme) {
		var slice []string
		plan.WeakPasswordStorageScheme.ElementsAs(ctx, &slice, false)
		addRequest.WeakPasswordStorageScheme = slice
	}
	if internaltypes.IsDefined(plan.WeakCryptEncoding) {
		var slice []string
		plan.WeakCryptEncoding.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdataSecurityAuditorWeakCryptEncodingProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdataSecurityAuditorWeakCryptEncodingPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.WeakCryptEncoding = enumSlice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for privilege data-security-auditor
func addOptionalPrivilegeDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddPrivilegeDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludePrivilege) {
		var slice []string
		plan.IncludePrivilege.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdataSecurityAuditorIncludePrivilegeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdataSecurityAuditorIncludePrivilegePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IncludePrivilege = enumSlice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for account-usability-issues data-security-auditor
func addOptionalAccountUsabilityIssuesDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddAccountUsabilityIssuesDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for locked-account data-security-auditor
func addOptionalLockedAccountDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddLockedAccountDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumIdleTime) {
		addRequest.MaximumIdleTime = plan.MaximumIdleTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for filter data-security-auditor
func addOptionalFilterDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddFilterDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for account-validity-window data-security-auditor
func addOptionalAccountValidityWindowDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddAccountValidityWindowDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountExpirationWarningInterval) {
		addRequest.AccountExpirationWarningInterval = plan.AccountExpirationWarningInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for multiple-password data-security-auditor
func addOptionalMultiplePasswordDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddMultiplePasswordDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for deprecated-password-storage-scheme data-security-auditor
func addOptionalDeprecatedPasswordStorageSchemeDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddDeprecatedPasswordStorageSchemeDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for nonexistent-password-policy data-security-auditor
func addOptionalNonexistentPasswordPolicyDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddNonexistentPasswordPolicyDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for access-control data-security-auditor
func addOptionalAccessControlDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddAccessControlDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReportFile) {
		addRequest.ReportFile = plan.ReportFile.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Add optional fields to create request for third-party data-security-auditor
func addOptionalThirdPartyDataSecurityAuditorFields(ctx context.Context, addRequest *client.AddThirdPartyDataSecurityAuditorRequest, plan dataSecurityAuditorResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAttribute) {
		var slice []string
		plan.IncludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.AuditBackend) {
		var slice []string
		plan.AuditBackend.ElementsAs(ctx, &slice, false)
		addRequest.AuditBackend = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuditSeverity) {
		auditSeverity, err := client.NewEnumdataSecurityAuditorAuditSeverityPropFromValue(plan.AuditSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuditSeverity = auditSeverity
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateDataSecurityAuditorUnknownValues(model *dataSecurityAuditorResourceModel) {
	if model.WeakCryptEncoding.IsUnknown() || model.WeakCryptEncoding.IsNull() {
		model.WeakCryptEncoding, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludePrivilege.IsUnknown() || model.IncludePrivilege.IsNull() {
		model.IncludePrivilege, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Filter.IsUnknown() || model.Filter.IsNull() {
		model.Filter, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.WeakPasswordStorageScheme.IsUnknown() || model.WeakPasswordStorageScheme.IsNull() {
		model.WeakPasswordStorageScheme, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *dataSecurityAuditorResourceModel) populateAllComputedStringAttributes() {
	if model.IdleAccountWarningInterval.IsUnknown() || model.IdleAccountWarningInterval.IsNull() {
		model.IdleAccountWarningInterval = types.StringValue("")
	}
	if model.PasswordEvaluationAge.IsUnknown() || model.PasswordEvaluationAge.IsNull() {
		model.PasswordEvaluationAge = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.IdleAccountErrorInterval.IsUnknown() || model.IdleAccountErrorInterval.IsNull() {
		model.IdleAccountErrorInterval = types.StringValue("")
	}
	if model.ReportFile.IsUnknown() || model.ReportFile.IsNull() {
		model.ReportFile = types.StringValue("")
	}
	if model.NeverLoggedInAccountWarningInterval.IsUnknown() || model.NeverLoggedInAccountWarningInterval.IsNull() {
		model.NeverLoggedInAccountWarningInterval = types.StringValue("")
	}
	if model.MaximumIdleTime.IsUnknown() || model.MaximumIdleTime.IsNull() {
		model.MaximumIdleTime = types.StringValue("")
	}
	if model.AuditSeverity.IsUnknown() || model.AuditSeverity.IsNull() {
		model.AuditSeverity = types.StringValue("")
	}
	if model.NeverLoggedInAccountErrorInterval.IsUnknown() || model.NeverLoggedInAccountErrorInterval.IsNull() {
		model.NeverLoggedInAccountErrorInterval = types.StringValue("")
	}
	if model.AccountExpirationWarningInterval.IsUnknown() || model.AccountExpirationWarningInterval.IsNull() {
		model.AccountExpirationWarningInterval = types.StringValue("")
	}
}

// Read a ExpiredPasswordDataSecurityAuditorResponse object into the model struct
func readExpiredPasswordDataSecurityAuditorResponse(ctx context.Context, r *client.ExpiredPasswordDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("expired-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.PasswordEvaluationAge = internaltypes.StringTypeOrNil(r.PasswordEvaluationAge, internaltypes.IsEmptyString(expectedValues.PasswordEvaluationAge))
	config.CheckMismatchedPDFormattedAttributes("password_evaluation_age",
		expectedValues.PasswordEvaluationAge, state.PasswordEvaluationAge, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a IdleAccountDataSecurityAuditorResponse object into the model struct
func readIdleAccountDataSecurityAuditorResponse(ctx context.Context, r *client.IdleAccountDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("idle-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IdleAccountWarningInterval = types.StringValue(r.IdleAccountWarningInterval)
	config.CheckMismatchedPDFormattedAttributes("idle_account_warning_interval",
		expectedValues.IdleAccountWarningInterval, state.IdleAccountWarningInterval, diagnostics)
	state.IdleAccountErrorInterval = internaltypes.StringTypeOrNil(r.IdleAccountErrorInterval, internaltypes.IsEmptyString(expectedValues.IdleAccountErrorInterval))
	config.CheckMismatchedPDFormattedAttributes("idle_account_error_interval",
		expectedValues.IdleAccountErrorInterval, state.IdleAccountErrorInterval, diagnostics)
	state.NeverLoggedInAccountWarningInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountWarningInterval, internaltypes.IsEmptyString(expectedValues.NeverLoggedInAccountWarningInterval))
	config.CheckMismatchedPDFormattedAttributes("never_logged_in_account_warning_interval",
		expectedValues.NeverLoggedInAccountWarningInterval, state.NeverLoggedInAccountWarningInterval, diagnostics)
	state.NeverLoggedInAccountErrorInterval = internaltypes.StringTypeOrNil(r.NeverLoggedInAccountErrorInterval, internaltypes.IsEmptyString(expectedValues.NeverLoggedInAccountErrorInterval))
	config.CheckMismatchedPDFormattedAttributes("never_logged_in_account_error_interval",
		expectedValues.NeverLoggedInAccountErrorInterval, state.NeverLoggedInAccountErrorInterval, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a DisabledAccountDataSecurityAuditorResponse object into the model struct
func readDisabledAccountDataSecurityAuditorResponse(ctx context.Context, r *client.DisabledAccountDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disabled-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a WeaklyEncodedPasswordDataSecurityAuditorResponse object into the model struct
func readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx context.Context, r *client.WeaklyEncodedPasswordDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a PrivilegeDataSecurityAuditorResponse object into the model struct
func readPrivilegeDataSecurityAuditorResponse(ctx context.Context, r *client.PrivilegeDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a AccountUsabilityIssuesDataSecurityAuditorResponse object into the model struct
func readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx context.Context, r *client.AccountUsabilityIssuesDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("account-usability-issues")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a LockedAccountDataSecurityAuditorResponse object into the model struct
func readLockedAccountDataSecurityAuditorResponse(ctx context.Context, r *client.LockedAccountDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("locked-account")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.MaximumIdleTime = internaltypes.StringTypeOrNil(r.MaximumIdleTime, internaltypes.IsEmptyString(expectedValues.MaximumIdleTime))
	config.CheckMismatchedPDFormattedAttributes("maximum_idle_time",
		expectedValues.MaximumIdleTime, state.MaximumIdleTime, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a FilterDataSecurityAuditorResponse object into the model struct
func readFilterDataSecurityAuditorResponse(ctx context.Context, r *client.FilterDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("filter")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Filter = internaltypes.GetStringSet(r.Filter)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a AccountValidityWindowDataSecurityAuditorResponse object into the model struct
func readAccountValidityWindowDataSecurityAuditorResponse(ctx context.Context, r *client.AccountValidityWindowDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("account-validity-window")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AccountExpirationWarningInterval = internaltypes.StringTypeOrNil(r.AccountExpirationWarningInterval, internaltypes.IsEmptyString(expectedValues.AccountExpirationWarningInterval))
	config.CheckMismatchedPDFormattedAttributes("account_expiration_warning_interval",
		expectedValues.AccountExpirationWarningInterval, state.AccountExpirationWarningInterval, diagnostics)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a MultiplePasswordDataSecurityAuditorResponse object into the model struct
func readMultiplePasswordDataSecurityAuditorResponse(ctx context.Context, r *client.MultiplePasswordDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multiple-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse object into the model struct
func readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx context.Context, r *client.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deprecated-password-storage-scheme")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a NonexistentPasswordPolicyDataSecurityAuditorResponse object into the model struct
func readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx context.Context, r *client.NonexistentPasswordPolicyDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("nonexistent-password-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a AccessControlDataSecurityAuditorResponse object into the model struct
func readAccessControlDataSecurityAuditorResponse(ctx context.Context, r *client.AccessControlDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("access-control")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ReportFile = types.StringValue(r.ReportFile)
	state.IncludeAttribute = internaltypes.GetStringSet(r.IncludeAttribute)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AuditBackend = internaltypes.GetStringSet(r.AuditBackend)
	state.AuditSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Read a ThirdPartyDataSecurityAuditorResponse object into the model struct
func readThirdPartyDataSecurityAuditorResponse(ctx context.Context, r *client.ThirdPartyDataSecurityAuditorResponse, state *dataSecurityAuditorResourceModel, expectedValues *dataSecurityAuditorResourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumdataSecurityAuditorAuditSeverityProp(r.AuditSeverity), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDataSecurityAuditorUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createDataSecurityAuditorOperations(plan dataSecurityAuditorResourceModel, state dataSecurityAuditorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ReportFile, state.ReportFile, "report-file")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountExpirationWarningInterval, state.AccountExpirationWarningInterval, "account-expiration-warning-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludePrivilege, state.IncludePrivilege, "include-privilege")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumIdleTime, state.MaximumIdleTime, "maximum-idle-time")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WeakPasswordStorageScheme, state.WeakPasswordStorageScheme, "weak-password-storage-scheme")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WeakCryptEncoding, state.WeakCryptEncoding, "weak-crypt-encoding")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleAccountWarningInterval, state.IdleAccountWarningInterval, "idle-account-warning-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleAccountErrorInterval, state.IdleAccountErrorInterval, "idle-account-error-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.NeverLoggedInAccountWarningInterval, state.NeverLoggedInAccountWarningInterval, "never-logged-in-account-warning-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.NeverLoggedInAccountErrorInterval, state.NeverLoggedInAccountErrorInterval, "never-logged-in-account-error-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttribute, state.IncludeAttribute, "include-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordEvaluationAge, state.PasswordEvaluationAge, "password-evaluation-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuditBackend, state.AuditBackend, "audit-backend")
	operations.AddStringOperationIfNecessary(&ops, plan.AuditSeverity, state.AuditSeverity, "audit-severity")
	return ops
}

// Create a expired-password data-security-auditor
func (r *dataSecurityAuditorResource) CreateExpiredPasswordDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddExpiredPasswordDataSecurityAuditorRequest([]client.EnumexpiredPasswordDataSecurityAuditorSchemaUrn{client.ENUMEXPIREDPASSWORDDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITOREXPIRED_PASSWORD},
		plan.Name.ValueString())
	err := addOptionalExpiredPasswordDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddExpiredPasswordDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readExpiredPasswordDataSecurityAuditorResponse(ctx, addResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a idle-account data-security-auditor
func (r *dataSecurityAuditorResource) CreateIdleAccountDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddIdleAccountDataSecurityAuditorRequest([]client.EnumidleAccountDataSecurityAuditorSchemaUrn{client.ENUMIDLEACCOUNTDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORIDLE_ACCOUNT},
		plan.IdleAccountWarningInterval.ValueString(),
		plan.Name.ValueString())
	err := addOptionalIdleAccountDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddIdleAccountDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readIdleAccountDataSecurityAuditorResponse(ctx, addResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a disabled-account data-security-auditor
func (r *dataSecurityAuditorResource) CreateDisabledAccountDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddDisabledAccountDataSecurityAuditorRequest([]client.EnumdisabledAccountDataSecurityAuditorSchemaUrn{client.ENUMDISABLEDACCOUNTDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORDISABLED_ACCOUNT},
		plan.Name.ValueString())
	err := addOptionalDisabledAccountDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddDisabledAccountDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readDisabledAccountDataSecurityAuditorResponse(ctx, addResponse.DisabledAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a weakly-encoded-password data-security-auditor
func (r *dataSecurityAuditorResource) CreateWeaklyEncodedPasswordDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddWeaklyEncodedPasswordDataSecurityAuditorRequest([]client.EnumweaklyEncodedPasswordDataSecurityAuditorSchemaUrn{client.ENUMWEAKLYENCODEDPASSWORDDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORWEAKLY_ENCODED_PASSWORD},
		plan.Name.ValueString())
	err := addOptionalWeaklyEncodedPasswordDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddWeaklyEncodedPasswordDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, addResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a privilege data-security-auditor
func (r *dataSecurityAuditorResource) CreatePrivilegeDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddPrivilegeDataSecurityAuditorRequest([]client.EnumprivilegeDataSecurityAuditorSchemaUrn{client.ENUMPRIVILEGEDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORPRIVILEGE},
		plan.Name.ValueString())
	err := addOptionalPrivilegeDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddPrivilegeDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readPrivilegeDataSecurityAuditorResponse(ctx, addResponse.PrivilegeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a account-usability-issues data-security-auditor
func (r *dataSecurityAuditorResource) CreateAccountUsabilityIssuesDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddAccountUsabilityIssuesDataSecurityAuditorRequest([]client.EnumaccountUsabilityIssuesDataSecurityAuditorSchemaUrn{client.ENUMACCOUNTUSABILITYISSUESDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORACCOUNT_USABILITY_ISSUES},
		plan.Name.ValueString())
	err := addOptionalAccountUsabilityIssuesDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddAccountUsabilityIssuesDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx, addResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a locked-account data-security-auditor
func (r *dataSecurityAuditorResource) CreateLockedAccountDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddLockedAccountDataSecurityAuditorRequest([]client.EnumlockedAccountDataSecurityAuditorSchemaUrn{client.ENUMLOCKEDACCOUNTDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORLOCKED_ACCOUNT},
		plan.Name.ValueString())
	err := addOptionalLockedAccountDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddLockedAccountDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readLockedAccountDataSecurityAuditorResponse(ctx, addResponse.LockedAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a filter data-security-auditor
func (r *dataSecurityAuditorResource) CreateFilterDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	var FilterSlice []string
	plan.Filter.ElementsAs(ctx, &FilterSlice, false)
	addRequest := client.NewAddFilterDataSecurityAuditorRequest([]client.EnumfilterDataSecurityAuditorSchemaUrn{client.ENUMFILTERDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORFILTER},
		plan.ReportFile.ValueString(),
		FilterSlice,
		plan.Name.ValueString())
	err := addOptionalFilterDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddFilterDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readFilterDataSecurityAuditorResponse(ctx, addResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a account-validity-window data-security-auditor
func (r *dataSecurityAuditorResource) CreateAccountValidityWindowDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddAccountValidityWindowDataSecurityAuditorRequest([]client.EnumaccountValidityWindowDataSecurityAuditorSchemaUrn{client.ENUMACCOUNTVALIDITYWINDOWDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORACCOUNT_VALIDITY_WINDOW},
		plan.Name.ValueString())
	err := addOptionalAccountValidityWindowDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddAccountValidityWindowDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readAccountValidityWindowDataSecurityAuditorResponse(ctx, addResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a multiple-password data-security-auditor
func (r *dataSecurityAuditorResource) CreateMultiplePasswordDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddMultiplePasswordDataSecurityAuditorRequest([]client.EnummultiplePasswordDataSecurityAuditorSchemaUrn{client.ENUMMULTIPLEPASSWORDDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORMULTIPLE_PASSWORD},
		plan.Name.ValueString())
	err := addOptionalMultiplePasswordDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddMultiplePasswordDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readMultiplePasswordDataSecurityAuditorResponse(ctx, addResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a deprecated-password-storage-scheme data-security-auditor
func (r *dataSecurityAuditorResource) CreateDeprecatedPasswordStorageSchemeDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddDeprecatedPasswordStorageSchemeDataSecurityAuditorRequest([]client.EnumdeprecatedPasswordStorageSchemeDataSecurityAuditorSchemaUrn{client.ENUMDEPRECATEDPASSWORDSTORAGESCHEMEDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORDEPRECATED_PASSWORD_STORAGE_SCHEME},
		plan.Name.ValueString())
	err := addOptionalDeprecatedPasswordStorageSchemeDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddDeprecatedPasswordStorageSchemeDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx, addResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a nonexistent-password-policy data-security-auditor
func (r *dataSecurityAuditorResource) CreateNonexistentPasswordPolicyDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddNonexistentPasswordPolicyDataSecurityAuditorRequest([]client.EnumnonexistentPasswordPolicyDataSecurityAuditorSchemaUrn{client.ENUMNONEXISTENTPASSWORDPOLICYDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORNONEXISTENT_PASSWORD_POLICY},
		plan.Name.ValueString())
	err := addOptionalNonexistentPasswordPolicyDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddNonexistentPasswordPolicyDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx, addResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a access-control data-security-auditor
func (r *dataSecurityAuditorResource) CreateAccessControlDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddAccessControlDataSecurityAuditorRequest([]client.EnumaccessControlDataSecurityAuditorSchemaUrn{client.ENUMACCESSCONTROLDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORACCESS_CONTROL},
		plan.Name.ValueString())
	err := addOptionalAccessControlDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddAccessControlDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readAccessControlDataSecurityAuditorResponse(ctx, addResponse.AccessControlDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party data-security-auditor
func (r *dataSecurityAuditorResource) CreateThirdPartyDataSecurityAuditor(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan dataSecurityAuditorResourceModel) (*dataSecurityAuditorResourceModel, error) {
	addRequest := client.NewAddThirdPartyDataSecurityAuditorRequest([]client.EnumthirdPartyDataSecurityAuditorSchemaUrn{client.ENUMTHIRDPARTYDATASECURITYAUDITORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DATA_SECURITY_AUDITORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.ReportFile.ValueString(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyDataSecurityAuditorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Data Security Auditor", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddDataSecurityAuditorRequest(
		client.AddThirdPartyDataSecurityAuditorRequestAsAddDataSecurityAuditorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.AddDataSecurityAuditorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Data Security Auditor", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state dataSecurityAuditorResourceModel
	readThirdPartyDataSecurityAuditorResponse(ctx, addResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *dataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *dataSecurityAuditorResourceModel
	var err error
	if plan.Type.ValueString() == "expired-password" {
		state, err = r.CreateExpiredPasswordDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "idle-account" {
		state, err = r.CreateIdleAccountDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "disabled-account" {
		state, err = r.CreateDisabledAccountDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "weakly-encoded-password" {
		state, err = r.CreateWeaklyEncodedPasswordDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "privilege" {
		state, err = r.CreatePrivilegeDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "account-usability-issues" {
		state, err = r.CreateAccountUsabilityIssuesDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "locked-account" {
		state, err = r.CreateLockedAccountDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "filter" {
		state, err = r.CreateFilterDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "account-validity-window" {
		state, err = r.CreateAccountValidityWindowDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "multiple-password" {
		state, err = r.CreateMultiplePasswordDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "deprecated-password-storage-scheme" {
		state, err = r.CreateDeprecatedPasswordStorageSchemeDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "nonexistent-password-policy" {
		state, err = r.CreateNonexistentPasswordPolicyDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "access-control" {
		state, err = r.CreateAccessControlDataSecurityAuditor(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyDataSecurityAuditor(ctx, req, resp, plan)
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
func (r *defaultDataSecurityAuditorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan dataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Data Security Auditor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state dataSecurityAuditorResourceModel
	if readResponse.ExpiredPasswordDataSecurityAuditorResponse != nil {
		readExpiredPasswordDataSecurityAuditorResponse(ctx, readResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IdleAccountDataSecurityAuditorResponse != nil {
		readIdleAccountDataSecurityAuditorResponse(ctx, readResponse.IdleAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DisabledAccountDataSecurityAuditorResponse != nil {
		readDisabledAccountDataSecurityAuditorResponse(ctx, readResponse.DisabledAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
		readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PrivilegeDataSecurityAuditorResponse != nil {
		readPrivilegeDataSecurityAuditorResponse(ctx, readResponse.PrivilegeDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
		readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx, readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LockedAccountDataSecurityAuditorResponse != nil {
		readLockedAccountDataSecurityAuditorResponse(ctx, readResponse.LockedAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FilterDataSecurityAuditorResponse != nil {
		readFilterDataSecurityAuditorResponse(ctx, readResponse.FilterDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccountValidityWindowDataSecurityAuditorResponse != nil {
		readAccountValidityWindowDataSecurityAuditorResponse(ctx, readResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MultiplePasswordDataSecurityAuditorResponse != nil {
		readMultiplePasswordDataSecurityAuditorResponse(ctx, readResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
		readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx, readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
		readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx, readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccessControlDataSecurityAuditorResponse != nil {
		readAccessControlDataSecurityAuditorResponse(ctx, readResponse.AccessControlDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyDataSecurityAuditorResponse != nil {
		readThirdPartyDataSecurityAuditorResponse(ctx, readResponse.ThirdPartyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DataSecurityAuditorAPI.UpdateDataSecurityAuditor(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DataSecurityAuditorAPI.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ExpiredPasswordDataSecurityAuditorResponse != nil {
			readExpiredPasswordDataSecurityAuditorResponse(ctx, updateResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IdleAccountDataSecurityAuditorResponse != nil {
			readIdleAccountDataSecurityAuditorResponse(ctx, updateResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DisabledAccountDataSecurityAuditorResponse != nil {
			readDisabledAccountDataSecurityAuditorResponse(ctx, updateResponse.DisabledAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
			readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PrivilegeDataSecurityAuditorResponse != nil {
			readPrivilegeDataSecurityAuditorResponse(ctx, updateResponse.PrivilegeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
			readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx, updateResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LockedAccountDataSecurityAuditorResponse != nil {
			readLockedAccountDataSecurityAuditorResponse(ctx, updateResponse.LockedAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FilterDataSecurityAuditorResponse != nil {
			readFilterDataSecurityAuditorResponse(ctx, updateResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccountValidityWindowDataSecurityAuditorResponse != nil {
			readAccountValidityWindowDataSecurityAuditorResponse(ctx, updateResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MultiplePasswordDataSecurityAuditorResponse != nil {
			readMultiplePasswordDataSecurityAuditorResponse(ctx, updateResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
			readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx, updateResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
			readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx, updateResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccessControlDataSecurityAuditorResponse != nil {
			readAccessControlDataSecurityAuditorResponse(ctx, updateResponse.AccessControlDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyDataSecurityAuditorResponse != nil {
			readThirdPartyDataSecurityAuditorResponse(ctx, updateResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *dataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDataSecurityAuditorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDataSecurityAuditor(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state dataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DataSecurityAuditorAPI.GetDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Data Security Auditor", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Data Security Auditor", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ExpiredPasswordDataSecurityAuditorResponse != nil {
		readExpiredPasswordDataSecurityAuditorResponse(ctx, readResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.IdleAccountDataSecurityAuditorResponse != nil {
		readIdleAccountDataSecurityAuditorResponse(ctx, readResponse.IdleAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DisabledAccountDataSecurityAuditorResponse != nil {
		readDisabledAccountDataSecurityAuditorResponse(ctx, readResponse.DisabledAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
		readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, readResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PrivilegeDataSecurityAuditorResponse != nil {
		readPrivilegeDataSecurityAuditorResponse(ctx, readResponse.PrivilegeDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
		readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx, readResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LockedAccountDataSecurityAuditorResponse != nil {
		readLockedAccountDataSecurityAuditorResponse(ctx, readResponse.LockedAccountDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FilterDataSecurityAuditorResponse != nil {
		readFilterDataSecurityAuditorResponse(ctx, readResponse.FilterDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccountValidityWindowDataSecurityAuditorResponse != nil {
		readAccountValidityWindowDataSecurityAuditorResponse(ctx, readResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MultiplePasswordDataSecurityAuditorResponse != nil {
		readMultiplePasswordDataSecurityAuditorResponse(ctx, readResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
		readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx, readResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
		readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx, readResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AccessControlDataSecurityAuditorResponse != nil {
		readAccessControlDataSecurityAuditorResponse(ctx, readResponse.AccessControlDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyDataSecurityAuditorResponse != nil {
		readThirdPartyDataSecurityAuditorResponse(ctx, readResponse.ThirdPartyDataSecurityAuditorResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *dataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDataSecurityAuditorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDataSecurityAuditor(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDataSecurityAuditor(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan dataSecurityAuditorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state dataSecurityAuditorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DataSecurityAuditorAPI.UpdateDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createDataSecurityAuditorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DataSecurityAuditorAPI.UpdateDataSecurityAuditorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Data Security Auditor", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ExpiredPasswordDataSecurityAuditorResponse != nil {
			readExpiredPasswordDataSecurityAuditorResponse(ctx, updateResponse.ExpiredPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.IdleAccountDataSecurityAuditorResponse != nil {
			readIdleAccountDataSecurityAuditorResponse(ctx, updateResponse.IdleAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DisabledAccountDataSecurityAuditorResponse != nil {
			readDisabledAccountDataSecurityAuditorResponse(ctx, updateResponse.DisabledAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse != nil {
			readWeaklyEncodedPasswordDataSecurityAuditorResponse(ctx, updateResponse.WeaklyEncodedPasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PrivilegeDataSecurityAuditorResponse != nil {
			readPrivilegeDataSecurityAuditorResponse(ctx, updateResponse.PrivilegeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccountUsabilityIssuesDataSecurityAuditorResponse != nil {
			readAccountUsabilityIssuesDataSecurityAuditorResponse(ctx, updateResponse.AccountUsabilityIssuesDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LockedAccountDataSecurityAuditorResponse != nil {
			readLockedAccountDataSecurityAuditorResponse(ctx, updateResponse.LockedAccountDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FilterDataSecurityAuditorResponse != nil {
			readFilterDataSecurityAuditorResponse(ctx, updateResponse.FilterDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccountValidityWindowDataSecurityAuditorResponse != nil {
			readAccountValidityWindowDataSecurityAuditorResponse(ctx, updateResponse.AccountValidityWindowDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MultiplePasswordDataSecurityAuditorResponse != nil {
			readMultiplePasswordDataSecurityAuditorResponse(ctx, updateResponse.MultiplePasswordDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse != nil {
			readDeprecatedPasswordStorageSchemeDataSecurityAuditorResponse(ctx, updateResponse.DeprecatedPasswordStorageSchemeDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse != nil {
			readNonexistentPasswordPolicyDataSecurityAuditorResponse(ctx, updateResponse.NonexistentPasswordPolicyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AccessControlDataSecurityAuditorResponse != nil {
			readAccessControlDataSecurityAuditorResponse(ctx, updateResponse.AccessControlDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyDataSecurityAuditorResponse != nil {
			readThirdPartyDataSecurityAuditorResponse(ctx, updateResponse.ThirdPartyDataSecurityAuditorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *dataSecurityAuditorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state dataSecurityAuditorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DataSecurityAuditorAPI.DeleteDataSecurityAuditorExecute(r.apiClient.DataSecurityAuditorAPI.DeleteDataSecurityAuditor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Data Security Auditor", err, httpResp)
		return
	}
}

func (r *dataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDataSecurityAuditor(ctx, req, resp)
}

func (r *defaultDataSecurityAuditorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDataSecurityAuditor(ctx, req, resp)
}

func importDataSecurityAuditor(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
