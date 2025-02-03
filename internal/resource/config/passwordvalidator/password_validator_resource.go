// Copyright © 2025 Ping Identity Corporation

package passwordvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordValidatorResource{}
	_ resource.ResourceWithConfigure   = &passwordValidatorResource{}
	_ resource.ResourceWithImportState = &passwordValidatorResource{}
	_ resource.Resource                = &defaultPasswordValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultPasswordValidatorResource{}
	_ resource.ResourceWithImportState = &defaultPasswordValidatorResource{}
)

// Create a Password Validator resource
func NewPasswordValidatorResource() resource.Resource {
	return &passwordValidatorResource{}
}

func NewDefaultPasswordValidatorResource() resource.Resource {
	return &defaultPasswordValidatorResource{}
}

// passwordValidatorResource is the resource implementation.
type passwordValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPasswordValidatorResource is the resource implementation.
type defaultPasswordValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passwordValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_validator"
}

func (r *defaultPasswordValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_password_validator"
}

// Configure adds the provider configured client to the resource.
func (r *passwordValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPasswordValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordValidatorResourceModel struct {
	Id                                             types.String `tfsdk:"id"`
	Name                                           types.String `tfsdk:"name"`
	Notifications                                  types.Set    `tfsdk:"notifications"`
	RequiredActions                                types.Set    `tfsdk:"required_actions"`
	Type                                           types.String `tfsdk:"type"`
	ExtensionClass                                 types.String `tfsdk:"extension_class"`
	ExtensionArgument                              types.Set    `tfsdk:"extension_argument"`
	MinUniqueCharacters                            types.Int64  `tfsdk:"min_unique_characters"`
	MatchPattern                                   types.String `tfsdk:"match_pattern"`
	MatchBehavior                                  types.String `tfsdk:"match_behavior"`
	MaxPasswordLength                              types.Int64  `tfsdk:"max_password_length"`
	MinPasswordLength                              types.Int64  `tfsdk:"min_password_length"`
	DisallowedCharacters                           types.String `tfsdk:"disallowed_characters"`
	DisallowedLeadingCharacters                    types.String `tfsdk:"disallowed_leading_characters"`
	DisallowedTrailingCharacters                   types.String `tfsdk:"disallowed_trailing_characters"`
	PwnedPasswordsBaseURL                          types.String `tfsdk:"pwned_passwords_base_url"`
	HttpProxyExternalServer                        types.String `tfsdk:"http_proxy_external_server"`
	HttpConnectTimeout                             types.String `tfsdk:"http_connect_timeout"`
	HttpResponseTimeout                            types.String `tfsdk:"http_response_timeout"`
	InvokeForAdd                                   types.Bool   `tfsdk:"invoke_for_add"`
	InvokeForSelfChange                            types.Bool   `tfsdk:"invoke_for_self_change"`
	InvokeForAdminReset                            types.Bool   `tfsdk:"invoke_for_admin_reset"`
	AcceptPasswordOnServiceError                   types.Bool   `tfsdk:"accept_password_on_service_error"`
	KeyManagerProvider                             types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                           types.String `tfsdk:"trust_manager_provider"`
	ScriptClass                                    types.String `tfsdk:"script_class"`
	ScriptArgument                                 types.Set    `tfsdk:"script_argument"`
	AllowNonAsciiCharacters                        types.Bool   `tfsdk:"allow_non_ascii_characters"`
	AllowUnknownCharacters                         types.Bool   `tfsdk:"allow_unknown_characters"`
	AllowedCharacterType                           types.Set    `tfsdk:"allowed_character_type"`
	AssumedPasswordGuessesPerSecond                types.String `tfsdk:"assumed_password_guesses_per_second"`
	MinimumAcceptableTimeToExhaustSearchSpace      types.String `tfsdk:"minimum_acceptable_time_to_exhaust_search_space"`
	DictionaryFile                                 types.String `tfsdk:"dictionary_file"`
	MaxConsecutiveLength                           types.Int64  `tfsdk:"max_consecutive_length"`
	CaseSensitiveValidation                        types.Bool   `tfsdk:"case_sensitive_validation"`
	IgnoreLeadingNonAlphabeticCharacters           types.Bool   `tfsdk:"ignore_leading_non_alphabetic_characters"`
	IgnoreTrailingNonAlphabeticCharacters          types.Bool   `tfsdk:"ignore_trailing_non_alphabetic_characters"`
	StripDiacriticalMarks                          types.Bool   `tfsdk:"strip_diacritical_marks"`
	AlternativePasswordCharacterMapping            types.Set    `tfsdk:"alternative_password_character_mapping"`
	MaximumAllowedPercentOfPassword                types.Int64  `tfsdk:"maximum_allowed_percent_of_password"`
	MatchAttribute                                 types.Set    `tfsdk:"match_attribute"`
	TestPasswordSubstringOfAttributeValue          types.Bool   `tfsdk:"test_password_substring_of_attribute_value"`
	TestAttributeValueSubstringOfPassword          types.Bool   `tfsdk:"test_attribute_value_substring_of_password"`
	MinimumAttributeValueLengthForSubstringMatches types.Int64  `tfsdk:"minimum_attribute_value_length_for_substring_matches"`
	TestReversedPassword                           types.Bool   `tfsdk:"test_reversed_password"`
	MinPasswordDifference                          types.Int64  `tfsdk:"min_password_difference"`
	CharacterSet                                   types.Set    `tfsdk:"character_set"`
	AllowUnclassifiedCharacters                    types.Bool   `tfsdk:"allow_unclassified_characters"`
	MinimumRequiredCharacterSets                   types.Int64  `tfsdk:"minimum_required_character_sets"`
	Description                                    types.String `tfsdk:"description"`
	Enabled                                        types.Bool   `tfsdk:"enabled"`
	ValidatorRequirementDescription                types.String `tfsdk:"validator_requirement_description"`
	ValidatorFailureMessage                        types.String `tfsdk:"validator_failure_message"`
}

// GetSchema defines the schema for the resource.
func (r *passwordValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordValidatorSchema(ctx, req, resp, false)
}

func (r *defaultPasswordValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordValidatorSchema(ctx, req, resp, true)
}

func passwordValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Password Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Validator resource. Options are ['character-set', 'similarity-based', 'attribute-value', 'custom', 'repeated-characters', 'dictionary', 'haystack', 'utf-8', 'groovy-scripted', 'pwned-passwords', 'disallowed-characters', 'length-based', 'regular-expression', 'unique-characters', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"character-set", "similarity-based", "attribute-value", "custom", "repeated-characters", "dictionary", "haystack", "utf-8", "groovy-scripted", "pwned-passwords", "disallowed-characters", "length-based", "regular-expression", "unique-characters", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Password Validator.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Password Validator. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"min_unique_characters": schema.Int64Attribute{
				Description: "Specifies the minimum number of unique characters that a password will be allowed to contain.",
				Optional:    true,
			},
			"match_pattern": schema.StringAttribute{
				Description: "The regular expression to use for this password validator.",
				Optional:    true,
			},
			"match_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if a user's proposed password matches the regular expression defined in the match-pattern property.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"require-match", "reject-match"}...),
				},
			},
			"max_password_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of characters that can be included in a proposed password.",
				Optional:    true,
				Computed:    true,
			},
			"min_password_length": schema.Int64Attribute{
				Description: "Specifies the minimum number of characters that must be included in a proposed password.",
				Optional:    true,
				Computed:    true,
			},
			"disallowed_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed anywhere in a password.",
				Optional:    true,
			},
			"disallowed_leading_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed as the first character of the password.",
				Optional:    true,
			},
			"disallowed_trailing_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed as the last character of the password.",
				Optional:    true,
			},
			"pwned_passwords_base_url": schema.StringAttribute{
				Description: "The base URL for requests used to interact with the Pwned Passwords service. The first five characters of the hexadecimal representation of the unsalted SHA-1 digest of a proposed password will be appended to this base URL to construct the HTTP GET request used to obtain information about potential matches.",
				Optional:    true,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Pwned Passwords service.",
				Optional:    true,
			},
			"http_connect_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait to obtain an HTTP connection.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_response_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. The maximum length of time to wait for a response to an HTTP request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"invoke_for_add": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided in LDAP add requests.",
				Optional:    true,
				Computed:    true,
			},
			"invoke_for_self_change": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided by an end user in the course of changing their own password.",
				Optional:    true,
				Computed:    true,
			},
			"invoke_for_admin_reset": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided by administrators when changing the password for another user.",
				Optional:    true,
				Computed:    true,
			},
			"accept_password_on_service_error": schema.BoolAttribute{
				Description: "Indicates whether to accept the proposed password if an error occurs while attempting to interact with the Pwned Passwords service.",
				Optional:    true,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies which key manager provider should be used to obtain a client certificate to present to the validation server when performing HTTPS communication. This may be left undefined if communication will not be secured with HTTPS, or if there is no need to present a client certificate to the validation service.",
				Optional:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies which trust manager provider should be used to determine whether to trust the certificate presented by the server when performing HTTPS communication. This may be left undefined if HTTPS communication is not needed, or if the validation service presents a certificate that is trusted by the default JVM configuration (which should be the case for the Pwned Password servers).",
				Optional:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Password Validator.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Password Validator. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allow_non_ascii_characters": schema.BoolAttribute{
				Description: "Indicates whether passwords will be allowed to include characters from outside the ASCII character set.",
				Optional:    true,
				Computed:    true,
			},
			"allow_unknown_characters": schema.BoolAttribute{
				Description: "Indicates whether passwords will be allowed to include characters that are not recognized by the JVM's Unicode support.",
				Optional:    true,
				Computed:    true,
			},
			"allowed_character_type": schema.SetAttribute{
				Description: "Specifies the set of character types that are allowed to be present in passwords.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"assumed_password_guesses_per_second": schema.StringAttribute{
				Description: "The number of password guesses per second that a potential attacker may be expected to make.",
				Optional:    true,
				Computed:    true,
			},
			"minimum_acceptable_time_to_exhaust_search_space": schema.StringAttribute{
				Description: "The minimum length of time (using the configured number of password guesses per second) required to exhaust the entire search space for a proposed password in order for that password to be considered acceptable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dictionary_file": schema.StringAttribute{
				Description: "Specifies the path to the file containing a list of words that cannot be used as passwords.",
				Optional:    true,
			},
			"max_consecutive_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of times that any character can appear consecutively in a password value.",
				Optional:    true,
			},
			"case_sensitive_validation": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`repeated-characters`, `unique-characters`]: Indicates whether this password validator should treat password characters in a case-sensitive manner. When the `type` attribute is set to `dictionary`: Indicates whether this password validator is to treat password characters in a case-sensitive manner.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`repeated-characters`, `unique-characters`]: Indicates whether this password validator should treat password characters in a case-sensitive manner.\n  - `dictionary`: Indicates whether this password validator is to treat password characters in a case-sensitive manner.",
				Optional:            true,
				Computed:            true,
			},
			"ignore_leading_non_alphabetic_characters": schema.BoolAttribute{
				Description: "Indicates whether to ignore any digits, symbols, or other non-alphabetic characters that may appear at the beginning of a proposed password.",
				Optional:    true,
				Computed:    true,
			},
			"ignore_trailing_non_alphabetic_characters": schema.BoolAttribute{
				Description: "Indicates whether to ignore any digits, symbols, or other non-alphabetic characters that may appear at the end of a proposed password.",
				Optional:    true,
				Computed:    true,
			},
			"strip_diacritical_marks": schema.BoolAttribute{
				Description: "Indicates whether to strip characters of any diacritical marks (like accents, cedillas, circumflexes, diaereses, tildes, and umlauts) they may contain. Any characters with a diacritical mark would be replaced with a base version",
				Optional:    true,
				Computed:    true,
			},
			"alternative_password_character_mapping": schema.SetAttribute{
				Description: "Provides a set of character substitutions that can be applied to the proposed password when checking to see if it is in the provided dictionary. Each mapping should consist of a single character followed by a colon and a list of the alternative characters that may be used in place of that character.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"maximum_allowed_percent_of_password": schema.Int64Attribute{
				Description: "The maximum allowed percent of a proposed password that any single dictionary word is allowed to comprise. A value of 100 indicates that a proposed password will only be rejected if the dictionary contains the entire proposed password (after any configured transformations have been applied).",
				Optional:    true,
				Computed:    true,
			},
			"match_attribute": schema.SetAttribute{
				Description: "Specifies the name(s) of the attribute(s) whose values should be checked to determine whether they match the provided password. If no values are provided, then the server checks if the proposed password matches the value of any user attribute in the target user's entry.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"test_password_substring_of_attribute_value": schema.BoolAttribute{
				Description: "Indicates whether to reject any proposed password that is a substring of a value in one of the match attributes in the target user's entry.",
				Optional:    true,
				Computed:    true,
			},
			"test_attribute_value_substring_of_password": schema.BoolAttribute{
				Description: "Indicates whether to reject any proposed password in which a value in one of the match attributes in the target user's entry is a substring of that password.",
				Optional:    true,
				Computed:    true,
			},
			"minimum_attribute_value_length_for_substring_matches": schema.Int64Attribute{
				Description: "The minimum length that an attribute value must have for it to be considered when rejecting passwords that contain the value of another attribute as a substring.",
				Optional:    true,
				Computed:    true,
			},
			"test_reversed_password": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `attribute-value`: Indicates whether to perform matching against the reversed value of the provided password in addition to the order in which it was given. When the `type` attribute is set to `dictionary`: Indicates whether this password validator is to test the reversed value of the provided password as well as the order in which it was given.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `attribute-value`: Indicates whether to perform matching against the reversed value of the provided password in addition to the order in which it was given.\n  - `dictionary`: Indicates whether this password validator is to test the reversed value of the provided password as well as the order in which it was given.",
				Optional:            true,
				Computed:            true,
			},
			"min_password_difference": schema.Int64Attribute{
				Description: "Specifies the minimum difference of new and old password.",
				Optional:    true,
			},
			"character_set": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `character-set`: Specifies a character set containing characters that a password may contain and a value indicating the minimum number of characters required from that set. When the `type` attribute is set to `repeated-characters`: Specifies a set of characters that should be considered equivalent for the purpose of this password validator. This can be used, for example, to ensure that passwords contain no more than three consecutive digits.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `character-set`: Specifies a character set containing characters that a password may contain and a value indicating the minimum number of characters required from that set.\n  - `repeated-characters`: Specifies a set of characters that should be considered equivalent for the purpose of this password validator. This can be used, for example, to ensure that passwords contain no more than three consecutive digits.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_unclassified_characters": schema.BoolAttribute{
				Description: "Indicates whether this password validator allows passwords to contain characters outside of any of the user-defined character sets.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"minimum_required_character_sets": schema.Int64Attribute{
				Description: "Specifies the minimum number of character sets that must be represented in a proposed password.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the password validator is enabled for use.",
				Required:    true,
			},
			"validator_requirement_description": schema.StringAttribute{
				Description: "Specifies a message that can be used to describe the requirements imposed by this password validator to end users. If a value is provided for this property, then it will override any description that may have otherwise been generated by the validator.",
				Optional:    true,
			},
			"validator_failure_message": schema.StringAttribute{
				Description: "Specifies a message that may be provided to the end user in the event that a proposed password is rejected by this validator. If a value is provided for this property, then it will override any failure message that may have otherwise been generated by the validator.",
				Optional:    true,
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
func (r *passwordValidatorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_password_validator")
	var planModel, configModel passwordValidatorResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for character-set type
	if resourceType == "character-set" {
		if !internaltypes.IsDefined(configModel.MinimumRequiredCharacterSets) {
			defaultVal := types.Int64Value(1)
			if !planModel.MinimumRequiredCharacterSets.Equal(defaultVal) {
				planModel.MinimumRequiredCharacterSets = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for attribute-value type
	if resourceType == "attribute-value" {
		if !internaltypes.IsDefined(configModel.TestPasswordSubstringOfAttributeValue) {
			defaultVal := types.BoolValue(false)
			if !planModel.TestPasswordSubstringOfAttributeValue.Equal(defaultVal) {
				planModel.TestPasswordSubstringOfAttributeValue = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TestAttributeValueSubstringOfPassword) {
			defaultVal := types.BoolValue(false)
			if !planModel.TestAttributeValueSubstringOfPassword.Equal(defaultVal) {
				planModel.TestAttributeValueSubstringOfPassword = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MinimumAttributeValueLengthForSubstringMatches) {
			defaultVal := types.Int64Value(4)
			if !planModel.MinimumAttributeValueLengthForSubstringMatches.Equal(defaultVal) {
				planModel.MinimumAttributeValueLengthForSubstringMatches = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for dictionary type
	if resourceType == "dictionary" {
		if !internaltypes.IsDefined(configModel.CaseSensitiveValidation) {
			defaultVal := types.BoolValue(false)
			if !planModel.CaseSensitiveValidation.Equal(defaultVal) {
				planModel.CaseSensitiveValidation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TestReversedPassword) {
			defaultVal := types.BoolValue(true)
			if !planModel.TestReversedPassword.Equal(defaultVal) {
				planModel.TestReversedPassword = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IgnoreLeadingNonAlphabeticCharacters) {
			defaultVal := types.BoolValue(false)
			if !planModel.IgnoreLeadingNonAlphabeticCharacters.Equal(defaultVal) {
				planModel.IgnoreLeadingNonAlphabeticCharacters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IgnoreTrailingNonAlphabeticCharacters) {
			defaultVal := types.BoolValue(false)
			if !planModel.IgnoreTrailingNonAlphabeticCharacters.Equal(defaultVal) {
				planModel.IgnoreTrailingNonAlphabeticCharacters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.StripDiacriticalMarks) {
			defaultVal := types.BoolValue(false)
			if !planModel.StripDiacriticalMarks.Equal(defaultVal) {
				planModel.StripDiacriticalMarks = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaximumAllowedPercentOfPassword) {
			defaultVal := types.Int64Value(100)
			if !planModel.MaximumAllowedPercentOfPassword.Equal(defaultVal) {
				planModel.MaximumAllowedPercentOfPassword = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for haystack type
	if resourceType == "haystack" {
		if !internaltypes.IsDefined(configModel.AssumedPasswordGuessesPerSecond) {
			defaultVal := types.StringValue("100,000,000,000")
			if !planModel.AssumedPasswordGuessesPerSecond.Equal(defaultVal) {
				planModel.AssumedPasswordGuessesPerSecond = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for utf-8 type
	if resourceType == "utf-8" {
		if !internaltypes.IsDefined(configModel.AllowNonAsciiCharacters) {
			defaultVal := types.BoolValue(true)
			if !planModel.AllowNonAsciiCharacters.Equal(defaultVal) {
				planModel.AllowNonAsciiCharacters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AllowUnknownCharacters) {
			defaultVal := types.BoolValue(false)
			if !planModel.AllowUnknownCharacters.Equal(defaultVal) {
				planModel.AllowUnknownCharacters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AllowedCharacterType) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("letters"), types.StringValue("numbers"), types.StringValue("punctuation"), types.StringValue("symbols"), types.StringValue("spaces"), types.StringValue("marks")})
			if !planModel.AllowedCharacterType.Equal(defaultVal) {
				planModel.AllowedCharacterType = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for pwned-passwords type
	if resourceType == "pwned-passwords" {
		if !internaltypes.IsDefined(configModel.PwnedPasswordsBaseURL) {
			defaultVal := types.StringValue("https://api.pwnedpasswords.com/range/")
			if !planModel.PwnedPasswordsBaseURL.Equal(defaultVal) {
				planModel.PwnedPasswordsBaseURL = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InvokeForAdd) {
			defaultVal := types.BoolValue(true)
			if !planModel.InvokeForAdd.Equal(defaultVal) {
				planModel.InvokeForAdd = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InvokeForSelfChange) {
			defaultVal := types.BoolValue(true)
			if !planModel.InvokeForSelfChange.Equal(defaultVal) {
				planModel.InvokeForSelfChange = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.InvokeForAdminReset) {
			defaultVal := types.BoolValue(true)
			if !planModel.InvokeForAdminReset.Equal(defaultVal) {
				planModel.InvokeForAdminReset = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AcceptPasswordOnServiceError) {
			defaultVal := types.BoolValue(true)
			if !planModel.AcceptPasswordOnServiceError.Equal(defaultVal) {
				planModel.AcceptPasswordOnServiceError = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for length-based type
	if resourceType == "length-based" {
		if !internaltypes.IsDefined(configModel.MaxPasswordLength) {
			defaultVal := types.Int64Value(0)
			if !planModel.MaxPasswordLength.Equal(defaultVal) {
				planModel.MaxPasswordLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MinPasswordLength) {
			defaultVal := types.Int64Value(6)
			if !planModel.MinPasswordLength.Equal(defaultVal) {
				planModel.MinPasswordLength = defaultVal
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

func (r *defaultPasswordValidatorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig, "pingdirectory_default_password_validator")
}

func modifyPlanPasswordValidator(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, resourceName string) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model passwordValidatorResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpConnectTimeout) {
		resp.Diagnostics.AddError("Attribute 'http_connect_timeout' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.HttpResponseTimeout) {
		resp.Diagnostics.AddError("Attribute 'http_response_timeout' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	compare, err = version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "utf-8" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9300,
			providerConfig.ProductVersion, resourceName+" with type \"utf_8\"")
	}
	if internaltypes.IsDefined(model.Type) && model.Type.ValueString() == "disallowed-characters" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9300,
			providerConfig.ProductVersion, resourceName+" with type \"disallowed_characters\"")
	}
	compare, err = version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

func (model *passwordValidatorResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "character-set" {
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "similarity-based" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "attribute-value" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "repeated-characters" {
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "dictionary" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "haystack" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "utf-8" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.StripDiacriticalMarks = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
	}
	if resourceType == "groovy-scripted" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "pwned-passwords" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxPasswordLength = types.Int64Null()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.CaseSensitiveValidation = types.BoolNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "disallowed-characters" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "length-based" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.InvokeForSelfChange = types.BoolNull()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "regular-expression" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "unique-characters" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
	if resourceType == "third-party" {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
		model.PwnedPasswordsBaseURL = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.InvokeForSelfChange = types.BoolNull()
		model.MinPasswordLength = types.Int64Null()
		model.IgnoreLeadingNonAlphabeticCharacters = types.BoolNull()
		model.InvokeForAdd = types.BoolNull()
		model.AcceptPasswordOnServiceError = types.BoolNull()
		model.TestAttributeValueSubstringOfPassword = types.BoolNull()
		model.MinimumRequiredCharacterSets = types.Int64Null()
		model.HttpResponseTimeout = types.StringNull()
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
		model.StripDiacriticalMarks = types.BoolNull()
		model.AllowNonAsciiCharacters = types.BoolNull()
		model.MinimumAttributeValueLengthForSubstringMatches = types.Int64Null()
		model.InvokeForAdminReset = types.BoolNull()
		model.CaseSensitiveValidation = types.BoolNull()
		model.HttpConnectTimeout = types.StringNull()
		model.AllowUnclassifiedCharacters = types.BoolNull()
		model.TestPasswordSubstringOfAttributeValue = types.BoolNull()
		model.MaximumAllowedPercentOfPassword = types.Int64Null()
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringNull()
		model.IgnoreTrailingNonAlphabeticCharacters = types.BoolNull()
		model.AssumedPasswordGuessesPerSecond = types.StringNull()
		model.TestReversedPassword = types.BoolNull()
		model.AllowUnknownCharacters = types.BoolNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsPasswordValidator() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"disallowed-characters"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("disallowed_characters"),
				path.MatchRoot("disallowed_leading_characters"),
				path.MatchRoot("disallowed_trailing_characters"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("character_set"),
			path.MatchRoot("type"),
			[]string{"character-set", "repeated-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_unclassified_characters"),
			path.MatchRoot("type"),
			[]string{"character-set"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("minimum_required_character_sets"),
			path.MatchRoot("type"),
			[]string{"character-set"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_password_difference"),
			path.MatchRoot("type"),
			[]string{"similarity-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_attribute"),
			path.MatchRoot("type"),
			[]string{"attribute-value"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("test_password_substring_of_attribute_value"),
			path.MatchRoot("type"),
			[]string{"attribute-value"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("test_attribute_value_substring_of_password"),
			path.MatchRoot("type"),
			[]string{"attribute-value"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("minimum_attribute_value_length_for_substring_matches"),
			path.MatchRoot("type"),
			[]string{"attribute-value"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("test_reversed_password"),
			path.MatchRoot("type"),
			[]string{"attribute-value", "dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_consecutive_length"),
			path.MatchRoot("type"),
			[]string{"repeated-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("case_sensitive_validation"),
			path.MatchRoot("type"),
			[]string{"repeated-characters", "dictionary", "unique-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("dictionary_file"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ignore_leading_non_alphabetic_characters"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ignore_trailing_non_alphabetic_characters"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("strip_diacritical_marks"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("alternative_password_character_mapping"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_allowed_percent_of_password"),
			path.MatchRoot("type"),
			[]string{"dictionary"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("assumed_password_guesses_per_second"),
			path.MatchRoot("type"),
			[]string{"haystack"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("minimum_acceptable_time_to_exhaust_search_space"),
			path.MatchRoot("type"),
			[]string{"haystack"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_non_ascii_characters"),
			path.MatchRoot("type"),
			[]string{"utf-8"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_unknown_characters"),
			path.MatchRoot("type"),
			[]string{"utf-8"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_character_type"),
			path.MatchRoot("type"),
			[]string{"utf-8"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("pwned_passwords_base_url"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_connect_timeout"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_response_timeout"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("invoke_for_add"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("invoke_for_self_change"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("invoke_for_admin_reset"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("accept_password_on_service_error"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_manager_provider"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_manager_provider"),
			path.MatchRoot("type"),
			[]string{"pwned-passwords"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disallowed_characters"),
			path.MatchRoot("type"),
			[]string{"disallowed-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disallowed_leading_characters"),
			path.MatchRoot("type"),
			[]string{"disallowed-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("disallowed_trailing_characters"),
			path.MatchRoot("type"),
			[]string{"disallowed-characters"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_password_length"),
			path.MatchRoot("type"),
			[]string{"length-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_password_length"),
			path.MatchRoot("type"),
			[]string{"length-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_pattern"),
			path.MatchRoot("type"),
			[]string{"regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_behavior"),
			path.MatchRoot("type"),
			[]string{"regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_unique_characters"),
			path.MatchRoot("type"),
			[]string{"unique-characters"},
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
			"character-set",
			[]path.Expression{path.MatchRoot("character_set"), path.MatchRoot("allow_unclassified_characters"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"similarity-based",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("min_password_difference")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"attribute-value",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("test_reversed_password")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"repeated-characters",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("max_consecutive_length"), path.MatchRoot("case_sensitive_validation")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"dictionary",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("dictionary_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"haystack",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"utf-8",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"pwned-passwords",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"disallowed-characters",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"length-based",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"regular-expression",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("match_pattern"), path.MatchRoot("match_behavior")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"unique-characters",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("case_sensitive_validation"), path.MatchRoot("min_unique_characters")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r passwordValidatorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPasswordValidator()
}

// Add config validators
func (r defaultPasswordValidatorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPasswordValidator()
}

// Add optional fields to create request for character-set password-validator
func addOptionalCharacterSetPasswordValidatorFields(ctx context.Context, addRequest *client.AddCharacterSetPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.MinimumRequiredCharacterSets) {
		addRequest.MinimumRequiredCharacterSets = plan.MinimumRequiredCharacterSets.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for similarity-based password-validator
func addOptionalSimilarityBasedPasswordValidatorFields(ctx context.Context, addRequest *client.AddSimilarityBasedPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for attribute-value password-validator
func addOptionalAttributeValuePasswordValidatorFields(ctx context.Context, addRequest *client.AddAttributeValuePasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.MatchAttribute) {
		var slice []string
		plan.MatchAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MatchAttribute = slice
	}
	if internaltypes.IsDefined(plan.TestPasswordSubstringOfAttributeValue) {
		addRequest.TestPasswordSubstringOfAttributeValue = plan.TestPasswordSubstringOfAttributeValue.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.TestAttributeValueSubstringOfPassword) {
		addRequest.TestAttributeValueSubstringOfPassword = plan.TestAttributeValueSubstringOfPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MinimumAttributeValueLengthForSubstringMatches) {
		addRequest.MinimumAttributeValueLengthForSubstringMatches = plan.MinimumAttributeValueLengthForSubstringMatches.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for repeated-characters password-validator
func addOptionalRepeatedCharactersPasswordValidatorFields(ctx context.Context, addRequest *client.AddRepeatedCharactersPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.CharacterSet) {
		var slice []string
		plan.CharacterSet.ElementsAs(ctx, &slice, false)
		addRequest.CharacterSet = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for dictionary password-validator
func addOptionalDictionaryPasswordValidatorFields(ctx context.Context, addRequest *client.AddDictionaryPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.CaseSensitiveValidation) {
		addRequest.CaseSensitiveValidation = plan.CaseSensitiveValidation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.TestReversedPassword) {
		addRequest.TestReversedPassword = plan.TestReversedPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IgnoreLeadingNonAlphabeticCharacters) {
		addRequest.IgnoreLeadingNonAlphabeticCharacters = plan.IgnoreLeadingNonAlphabeticCharacters.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IgnoreTrailingNonAlphabeticCharacters) {
		addRequest.IgnoreTrailingNonAlphabeticCharacters = plan.IgnoreTrailingNonAlphabeticCharacters.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.StripDiacriticalMarks) {
		addRequest.StripDiacriticalMarks = plan.StripDiacriticalMarks.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AlternativePasswordCharacterMapping) {
		var slice []string
		plan.AlternativePasswordCharacterMapping.ElementsAs(ctx, &slice, false)
		addRequest.AlternativePasswordCharacterMapping = slice
	}
	if internaltypes.IsDefined(plan.MaximumAllowedPercentOfPassword) {
		addRequest.MaximumAllowedPercentOfPassword = plan.MaximumAllowedPercentOfPassword.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for haystack password-validator
func addOptionalHaystackPasswordValidatorFields(ctx context.Context, addRequest *client.AddHaystackPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AssumedPasswordGuessesPerSecond) {
		addRequest.AssumedPasswordGuessesPerSecond = plan.AssumedPasswordGuessesPerSecond.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinimumAcceptableTimeToExhaustSearchSpace) {
		addRequest.MinimumAcceptableTimeToExhaustSearchSpace = plan.MinimumAcceptableTimeToExhaustSearchSpace.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for utf-8 password-validator
func addOptionalUtf8PasswordValidatorFields(ctx context.Context, addRequest *client.AddUtf8PasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.AllowNonAsciiCharacters) {
		addRequest.AllowNonAsciiCharacters = plan.AllowNonAsciiCharacters.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowUnknownCharacters) {
		addRequest.AllowUnknownCharacters = plan.AllowUnknownCharacters.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowedCharacterType) {
		var slice []string
		plan.AllowedCharacterType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpasswordValidatorAllowedCharacterTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpasswordValidatorAllowedCharacterTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedCharacterType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for groovy-scripted password-validator
func addOptionalGroovyScriptedPasswordValidatorFields(ctx context.Context, addRequest *client.AddGroovyScriptedPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for pwned-passwords password-validator
func addOptionalPwnedPasswordsPasswordValidatorFields(ctx context.Context, addRequest *client.AddPwnedPasswordsPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PwnedPasswordsBaseURL) {
		addRequest.PwnedPasswordsBaseURL = plan.PwnedPasswordsBaseURL.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpConnectTimeout) {
		addRequest.HttpConnectTimeout = plan.HttpConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpResponseTimeout) {
		addRequest.HttpResponseTimeout = plan.HttpResponseTimeout.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForAdd) {
		addRequest.InvokeForAdd = plan.InvokeForAdd.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForSelfChange) {
		addRequest.InvokeForSelfChange = plan.InvokeForSelfChange.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.InvokeForAdminReset) {
		addRequest.InvokeForAdminReset = plan.InvokeForAdminReset.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AcceptPasswordOnServiceError) {
		addRequest.AcceptPasswordOnServiceError = plan.AcceptPasswordOnServiceError.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for disallowed-characters password-validator
func addOptionalDisallowedCharactersPasswordValidatorFields(ctx context.Context, addRequest *client.AddDisallowedCharactersPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisallowedCharacters) {
		addRequest.DisallowedCharacters = plan.DisallowedCharacters.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisallowedLeadingCharacters) {
		addRequest.DisallowedLeadingCharacters = plan.DisallowedLeadingCharacters.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisallowedTrailingCharacters) {
		addRequest.DisallowedTrailingCharacters = plan.DisallowedTrailingCharacters.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for length-based password-validator
func addOptionalLengthBasedPasswordValidatorFields(ctx context.Context, addRequest *client.AddLengthBasedPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.MaxPasswordLength) {
		addRequest.MaxPasswordLength = plan.MaxPasswordLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MinPasswordLength) {
		addRequest.MinPasswordLength = plan.MinPasswordLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for regular-expression password-validator
func addOptionalRegularExpressionPasswordValidatorFields(ctx context.Context, addRequest *client.AddRegularExpressionPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for unique-characters password-validator
func addOptionalUniqueCharactersPasswordValidatorFields(ctx context.Context, addRequest *client.AddUniqueCharactersPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party password-validator
func addOptionalThirdPartyPasswordValidatorFields(ctx context.Context, addRequest *client.AddThirdPartyPasswordValidatorRequest, plan passwordValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorRequirementDescription) {
		addRequest.ValidatorRequirementDescription = plan.ValidatorRequirementDescription.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ValidatorFailureMessage) {
		addRequest.ValidatorFailureMessage = plan.ValidatorFailureMessage.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populatePasswordValidatorUnknownValues(model *passwordValidatorResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AlternativePasswordCharacterMapping.IsUnknown() || model.AlternativePasswordCharacterMapping.IsNull() {
		model.AlternativePasswordCharacterMapping, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedCharacterType.IsUnknown() || model.AllowedCharacterType.IsNull() {
		model.AllowedCharacterType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.CharacterSet.IsUnknown() || model.CharacterSet.IsNull() {
		model.CharacterSet, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MatchAttribute.IsUnknown() || model.MatchAttribute.IsNull() {
		model.MatchAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *passwordValidatorResourceModel) populateAllComputedStringAttributes() {
	if model.TrustManagerProvider.IsUnknown() || model.TrustManagerProvider.IsNull() {
		model.TrustManagerProvider = types.StringValue("")
	}
	if model.MatchBehavior.IsUnknown() || model.MatchBehavior.IsNull() {
		model.MatchBehavior = types.StringValue("")
	}
	if model.DisallowedTrailingCharacters.IsUnknown() || model.DisallowedTrailingCharacters.IsNull() {
		model.DisallowedTrailingCharacters = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.DictionaryFile.IsUnknown() || model.DictionaryFile.IsNull() {
		model.DictionaryFile = types.StringValue("")
	}
	if model.DisallowedLeadingCharacters.IsUnknown() || model.DisallowedLeadingCharacters.IsNull() {
		model.DisallowedLeadingCharacters = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.ValidatorFailureMessage.IsUnknown() || model.ValidatorFailureMessage.IsNull() {
		model.ValidatorFailureMessage = types.StringValue("")
	}
	if model.MatchPattern.IsUnknown() || model.MatchPattern.IsNull() {
		model.MatchPattern = types.StringValue("")
	}
	if model.AssumedPasswordGuessesPerSecond.IsUnknown() || model.AssumedPasswordGuessesPerSecond.IsNull() {
		model.AssumedPasswordGuessesPerSecond = types.StringValue("")
	}
	if model.MinimumAcceptableTimeToExhaustSearchSpace.IsUnknown() || model.MinimumAcceptableTimeToExhaustSearchSpace.IsNull() {
		model.MinimumAcceptableTimeToExhaustSearchSpace = types.StringValue("")
	}
	if model.DisallowedCharacters.IsUnknown() || model.DisallowedCharacters.IsNull() {
		model.DisallowedCharacters = types.StringValue("")
	}
	if model.HttpConnectTimeout.IsUnknown() || model.HttpConnectTimeout.IsNull() {
		model.HttpConnectTimeout = types.StringValue("")
	}
	if model.KeyManagerProvider.IsUnknown() || model.KeyManagerProvider.IsNull() {
		model.KeyManagerProvider = types.StringValue("")
	}
	if model.PwnedPasswordsBaseURL.IsUnknown() || model.PwnedPasswordsBaseURL.IsNull() {
		model.PwnedPasswordsBaseURL = types.StringValue("")
	}
	if model.HttpResponseTimeout.IsUnknown() || model.HttpResponseTimeout.IsNull() {
		model.HttpResponseTimeout = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
	if model.ValidatorRequirementDescription.IsUnknown() || model.ValidatorRequirementDescription.IsNull() {
		model.ValidatorRequirementDescription = types.StringValue("")
	}
}

// Read a CharacterSetPasswordValidatorResponse object into the model struct
func readCharacterSetPasswordValidatorResponse(ctx context.Context, r *client.CharacterSetPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("character-set")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CharacterSet = internaltypes.GetStringSet(r.CharacterSet)
	state.AllowUnclassifiedCharacters = types.BoolValue(r.AllowUnclassifiedCharacters)
	state.MinimumRequiredCharacterSets = internaltypes.Int64TypeOrNil(r.MinimumRequiredCharacterSets)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a SimilarityBasedPasswordValidatorResponse object into the model struct
func readSimilarityBasedPasswordValidatorResponse(ctx context.Context, r *client.SimilarityBasedPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("similarity-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MinPasswordDifference = types.Int64Value(r.MinPasswordDifference)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a AttributeValuePasswordValidatorResponse object into the model struct
func readAttributeValuePasswordValidatorResponse(ctx context.Context, r *client.AttributeValuePasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-value")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.TestPasswordSubstringOfAttributeValue = internaltypes.BoolTypeOrNil(r.TestPasswordSubstringOfAttributeValue)
	state.TestAttributeValueSubstringOfPassword = internaltypes.BoolTypeOrNil(r.TestAttributeValueSubstringOfPassword)
	state.MinimumAttributeValueLengthForSubstringMatches = internaltypes.Int64TypeOrNil(r.MinimumAttributeValueLengthForSubstringMatches)
	state.TestReversedPassword = types.BoolValue(r.TestReversedPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a CustomPasswordValidatorResponse object into the model struct
func readCustomPasswordValidatorResponse(ctx context.Context, r *client.CustomPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a RepeatedCharactersPasswordValidatorResponse object into the model struct
func readRepeatedCharactersPasswordValidatorResponse(ctx context.Context, r *client.RepeatedCharactersPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("repeated-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxConsecutiveLength = types.Int64Value(r.MaxConsecutiveLength)
	state.CaseSensitiveValidation = types.BoolValue(r.CaseSensitiveValidation)
	state.CharacterSet = internaltypes.GetStringSet(r.CharacterSet)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a DictionaryPasswordValidatorResponse object into the model struct
func readDictionaryPasswordValidatorResponse(ctx context.Context, r *client.DictionaryPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("dictionary")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DictionaryFile = types.StringValue(r.DictionaryFile)
	state.CaseSensitiveValidation = types.BoolValue(r.CaseSensitiveValidation)
	state.TestReversedPassword = types.BoolValue(r.TestReversedPassword)
	state.IgnoreLeadingNonAlphabeticCharacters = internaltypes.BoolTypeOrNil(r.IgnoreLeadingNonAlphabeticCharacters)
	state.IgnoreTrailingNonAlphabeticCharacters = internaltypes.BoolTypeOrNil(r.IgnoreTrailingNonAlphabeticCharacters)
	state.StripDiacriticalMarks = internaltypes.BoolTypeOrNil(r.StripDiacriticalMarks)
	state.AlternativePasswordCharacterMapping = internaltypes.GetStringSet(r.AlternativePasswordCharacterMapping)
	state.MaximumAllowedPercentOfPassword = internaltypes.Int64TypeOrNil(r.MaximumAllowedPercentOfPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a HaystackPasswordValidatorResponse object into the model struct
func readHaystackPasswordValidatorResponse(ctx context.Context, r *client.HaystackPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("haystack")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AssumedPasswordGuessesPerSecond = types.StringValue(r.AssumedPasswordGuessesPerSecond)
	state.MinimumAcceptableTimeToExhaustSearchSpace = types.StringValue(r.MinimumAcceptableTimeToExhaustSearchSpace)
	config.CheckMismatchedPDFormattedAttributes("minimum_acceptable_time_to_exhaust_search_space",
		expectedValues.MinimumAcceptableTimeToExhaustSearchSpace, state.MinimumAcceptableTimeToExhaustSearchSpace, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a Utf8PasswordValidatorResponse object into the model struct
func readUtf8PasswordValidatorResponse(ctx context.Context, r *client.Utf8PasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("utf-8")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowNonAsciiCharacters = internaltypes.BoolTypeOrNil(r.AllowNonAsciiCharacters)
	state.AllowUnknownCharacters = internaltypes.BoolTypeOrNil(r.AllowUnknownCharacters)
	state.AllowedCharacterType = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordValidatorAllowedCharacterTypeProp(r.AllowedCharacterType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a GroovyScriptedPasswordValidatorResponse object into the model struct
func readGroovyScriptedPasswordValidatorResponse(ctx context.Context, r *client.GroovyScriptedPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a PwnedPasswordsPasswordValidatorResponse object into the model struct
func readPwnedPasswordsPasswordValidatorResponse(ctx context.Context, r *client.PwnedPasswordsPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pwned-passwords")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PwnedPasswordsBaseURL = types.StringValue(r.PwnedPasswordsBaseURL)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.HttpConnectTimeout = internaltypes.StringTypeOrNil(r.HttpConnectTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_connect_timeout",
		expectedValues.HttpConnectTimeout, state.HttpConnectTimeout, diagnostics)
	state.HttpResponseTimeout = internaltypes.StringTypeOrNil(r.HttpResponseTimeout, true)
	config.CheckMismatchedPDFormattedAttributes("http_response_timeout",
		expectedValues.HttpResponseTimeout, state.HttpResponseTimeout, diagnostics)
	state.InvokeForAdd = types.BoolValue(r.InvokeForAdd)
	state.InvokeForSelfChange = types.BoolValue(r.InvokeForSelfChange)
	state.InvokeForAdminReset = types.BoolValue(r.InvokeForAdminReset)
	state.AcceptPasswordOnServiceError = types.BoolValue(r.AcceptPasswordOnServiceError)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a DisallowedCharactersPasswordValidatorResponse object into the model struct
func readDisallowedCharactersPasswordValidatorResponse(ctx context.Context, r *client.DisallowedCharactersPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disallowed-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DisallowedCharacters = internaltypes.StringTypeOrNil(r.DisallowedCharacters, internaltypes.IsEmptyString(expectedValues.DisallowedCharacters))
	state.DisallowedLeadingCharacters = internaltypes.StringTypeOrNil(r.DisallowedLeadingCharacters, internaltypes.IsEmptyString(expectedValues.DisallowedLeadingCharacters))
	state.DisallowedTrailingCharacters = internaltypes.StringTypeOrNil(r.DisallowedTrailingCharacters, internaltypes.IsEmptyString(expectedValues.DisallowedTrailingCharacters))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a LengthBasedPasswordValidatorResponse object into the model struct
func readLengthBasedPasswordValidatorResponse(ctx context.Context, r *client.LengthBasedPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("length-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.MinPasswordLength = internaltypes.Int64TypeOrNil(r.MinPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a RegularExpressionPasswordValidatorResponse object into the model struct
func readRegularExpressionPasswordValidatorResponse(ctx context.Context, r *client.RegularExpressionPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("regular-expression")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchPattern = types.StringValue(r.MatchPattern)
	state.MatchBehavior = types.StringValue(r.MatchBehavior.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a UniqueCharactersPasswordValidatorResponse object into the model struct
func readUniqueCharactersPasswordValidatorResponse(ctx context.Context, r *client.UniqueCharactersPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unique-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MinUniqueCharacters = types.Int64Value(r.MinUniqueCharacters)
	state.CaseSensitiveValidation = types.BoolValue(r.CaseSensitiveValidation)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Read a ThirdPartyPasswordValidatorResponse object into the model struct
func readThirdPartyPasswordValidatorResponse(ctx context.Context, r *client.ThirdPartyPasswordValidatorResponse, state *passwordValidatorResourceModel, expectedValues *passwordValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, internaltypes.IsEmptyString(expectedValues.ValidatorRequirementDescription))
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, internaltypes.IsEmptyString(expectedValues.ValidatorFailureMessage))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordValidatorUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createPasswordValidatorOperations(plan passwordValidatorResourceModel, state passwordValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinUniqueCharacters, state.MinUniqueCharacters, "min-unique-characters")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchPattern, state.MatchPattern, "match-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchBehavior, state.MatchBehavior, "match-behavior")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxPasswordLength, state.MaxPasswordLength, "max-password-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinPasswordLength, state.MinPasswordLength, "min-password-length")
	operations.AddStringOperationIfNecessary(&ops, plan.DisallowedCharacters, state.DisallowedCharacters, "disallowed-characters")
	operations.AddStringOperationIfNecessary(&ops, plan.DisallowedLeadingCharacters, state.DisallowedLeadingCharacters, "disallowed-leading-characters")
	operations.AddStringOperationIfNecessary(&ops, plan.DisallowedTrailingCharacters, state.DisallowedTrailingCharacters, "disallowed-trailing-characters")
	operations.AddStringOperationIfNecessary(&ops, plan.PwnedPasswordsBaseURL, state.PwnedPasswordsBaseURL, "pwned-passwords-base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpConnectTimeout, state.HttpConnectTimeout, "http-connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpResponseTimeout, state.HttpResponseTimeout, "http-response-timeout")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForAdd, state.InvokeForAdd, "invoke-for-add")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForSelfChange, state.InvokeForSelfChange, "invoke-for-self-change")
	operations.AddBoolOperationIfNecessary(&ops, plan.InvokeForAdminReset, state.InvokeForAdminReset, "invoke-for-admin-reset")
	operations.AddBoolOperationIfNecessary(&ops, plan.AcceptPasswordOnServiceError, state.AcceptPasswordOnServiceError, "accept-password-on-service-error")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowNonAsciiCharacters, state.AllowNonAsciiCharacters, "allow-non-ascii-characters")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUnknownCharacters, state.AllowUnknownCharacters, "allow-unknown-characters")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedCharacterType, state.AllowedCharacterType, "allowed-character-type")
	operations.AddStringOperationIfNecessary(&ops, plan.AssumedPasswordGuessesPerSecond, state.AssumedPasswordGuessesPerSecond, "assumed-password-guesses-per-second")
	operations.AddStringOperationIfNecessary(&ops, plan.MinimumAcceptableTimeToExhaustSearchSpace, state.MinimumAcceptableTimeToExhaustSearchSpace, "minimum-acceptable-time-to-exhaust-search-space")
	operations.AddStringOperationIfNecessary(&ops, plan.DictionaryFile, state.DictionaryFile, "dictionary-file")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConsecutiveLength, state.MaxConsecutiveLength, "max-consecutive-length")
	operations.AddBoolOperationIfNecessary(&ops, plan.CaseSensitiveValidation, state.CaseSensitiveValidation, "case-sensitive-validation")
	operations.AddBoolOperationIfNecessary(&ops, plan.IgnoreLeadingNonAlphabeticCharacters, state.IgnoreLeadingNonAlphabeticCharacters, "ignore-leading-non-alphabetic-characters")
	operations.AddBoolOperationIfNecessary(&ops, plan.IgnoreTrailingNonAlphabeticCharacters, state.IgnoreTrailingNonAlphabeticCharacters, "ignore-trailing-non-alphabetic-characters")
	operations.AddBoolOperationIfNecessary(&ops, plan.StripDiacriticalMarks, state.StripDiacriticalMarks, "strip-diacritical-marks")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AlternativePasswordCharacterMapping, state.AlternativePasswordCharacterMapping, "alternative-password-character-mapping")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumAllowedPercentOfPassword, state.MaximumAllowedPercentOfPassword, "maximum-allowed-percent-of-password")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MatchAttribute, state.MatchAttribute, "match-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.TestPasswordSubstringOfAttributeValue, state.TestPasswordSubstringOfAttributeValue, "test-password-substring-of-attribute-value")
	operations.AddBoolOperationIfNecessary(&ops, plan.TestAttributeValueSubstringOfPassword, state.TestAttributeValueSubstringOfPassword, "test-attribute-value-substring-of-password")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumAttributeValueLengthForSubstringMatches, state.MinimumAttributeValueLengthForSubstringMatches, "minimum-attribute-value-length-for-substring-matches")
	operations.AddBoolOperationIfNecessary(&ops, plan.TestReversedPassword, state.TestReversedPassword, "test-reversed-password")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinPasswordDifference, state.MinPasswordDifference, "min-password-difference")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CharacterSet, state.CharacterSet, "character-set")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUnclassifiedCharacters, state.AllowUnclassifiedCharacters, "allow-unclassified-characters")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumRequiredCharacterSets, state.MinimumRequiredCharacterSets, "minimum-required-character-sets")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidatorRequirementDescription, state.ValidatorRequirementDescription, "validator-requirement-description")
	operations.AddStringOperationIfNecessary(&ops, plan.ValidatorFailureMessage, state.ValidatorFailureMessage, "validator-failure-message")
	return ops
}

// Create a character-set password-validator
func (r *passwordValidatorResource) CreateCharacterSetPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	var CharacterSetSlice []string
	plan.CharacterSet.ElementsAs(ctx, &CharacterSetSlice, false)
	addRequest := client.NewAddCharacterSetPasswordValidatorRequest([]client.EnumcharacterSetPasswordValidatorSchemaUrn{client.ENUMCHARACTERSETPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORCHARACTER_SET},
		CharacterSetSlice,
		plan.AllowUnclassifiedCharacters.ValueBool(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalCharacterSetPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddCharacterSetPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readCharacterSetPasswordValidatorResponse(ctx, addResponse.CharacterSetPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a similarity-based password-validator
func (r *passwordValidatorResource) CreateSimilarityBasedPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddSimilarityBasedPasswordValidatorRequest([]client.EnumsimilarityBasedPasswordValidatorSchemaUrn{client.ENUMSIMILARITYBASEDPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORSIMILARITY_BASED},
		plan.MinPasswordDifference.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalSimilarityBasedPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddSimilarityBasedPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readSimilarityBasedPasswordValidatorResponse(ctx, addResponse.SimilarityBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a attribute-value password-validator
func (r *passwordValidatorResource) CreateAttributeValuePasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddAttributeValuePasswordValidatorRequest([]client.EnumattributeValuePasswordValidatorSchemaUrn{client.ENUMATTRIBUTEVALUEPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORATTRIBUTE_VALUE},
		plan.TestReversedPassword.ValueBool(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalAttributeValuePasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddAttributeValuePasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readAttributeValuePasswordValidatorResponse(ctx, addResponse.AttributeValuePasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a repeated-characters password-validator
func (r *passwordValidatorResource) CreateRepeatedCharactersPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddRepeatedCharactersPasswordValidatorRequest([]client.EnumrepeatedCharactersPasswordValidatorSchemaUrn{client.ENUMREPEATEDCHARACTERSPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORREPEATED_CHARACTERS},
		plan.MaxConsecutiveLength.ValueInt64(),
		plan.CaseSensitiveValidation.ValueBool(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalRepeatedCharactersPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddRepeatedCharactersPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readRepeatedCharactersPasswordValidatorResponse(ctx, addResponse.RepeatedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a dictionary password-validator
func (r *passwordValidatorResource) CreateDictionaryPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddDictionaryPasswordValidatorRequest([]client.EnumdictionaryPasswordValidatorSchemaUrn{client.ENUMDICTIONARYPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORDICTIONARY},
		plan.DictionaryFile.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalDictionaryPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddDictionaryPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readDictionaryPasswordValidatorResponse(ctx, addResponse.DictionaryPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a haystack password-validator
func (r *passwordValidatorResource) CreateHaystackPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddHaystackPasswordValidatorRequest([]client.EnumhaystackPasswordValidatorSchemaUrn{client.ENUMHAYSTACKPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORHAYSTACK},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalHaystackPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddHaystackPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readHaystackPasswordValidatorResponse(ctx, addResponse.HaystackPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a utf-8 password-validator
func (r *passwordValidatorResource) CreateUtf8PasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddUtf8PasswordValidatorRequest([]client.Enumutf8PasswordValidatorSchemaUrn{client.ENUMUTF8PASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORUTF_8},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalUtf8PasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddUtf8PasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readUtf8PasswordValidatorResponse(ctx, addResponse.Utf8PasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted password-validator
func (r *passwordValidatorResource) CreateGroovyScriptedPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedPasswordValidatorRequest([]client.EnumgroovyScriptedPasswordValidatorSchemaUrn{client.ENUMGROOVYSCRIPTEDPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalGroovyScriptedPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddGroovyScriptedPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readGroovyScriptedPasswordValidatorResponse(ctx, addResponse.GroovyScriptedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pwned-passwords password-validator
func (r *passwordValidatorResource) CreatePwnedPasswordsPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddPwnedPasswordsPasswordValidatorRequest([]client.EnumpwnedPasswordsPasswordValidatorSchemaUrn{client.ENUMPWNEDPASSWORDSPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORPWNED_PASSWORDS},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalPwnedPasswordsPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddPwnedPasswordsPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readPwnedPasswordsPasswordValidatorResponse(ctx, addResponse.PwnedPasswordsPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a disallowed-characters password-validator
func (r *passwordValidatorResource) CreateDisallowedCharactersPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddDisallowedCharactersPasswordValidatorRequest([]client.EnumdisallowedCharactersPasswordValidatorSchemaUrn{client.ENUMDISALLOWEDCHARACTERSPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORDISALLOWED_CHARACTERS},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalDisallowedCharactersPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddDisallowedCharactersPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readDisallowedCharactersPasswordValidatorResponse(ctx, addResponse.DisallowedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a length-based password-validator
func (r *passwordValidatorResource) CreateLengthBasedPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddLengthBasedPasswordValidatorRequest([]client.EnumlengthBasedPasswordValidatorSchemaUrn{client.ENUMLENGTHBASEDPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORLENGTH_BASED},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalLengthBasedPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddLengthBasedPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readLengthBasedPasswordValidatorResponse(ctx, addResponse.LengthBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a regular-expression password-validator
func (r *passwordValidatorResource) CreateRegularExpressionPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	matchBehavior, err := client.NewEnumpasswordValidatorMatchBehaviorPropFromValue(plan.MatchBehavior.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for MatchBehavior", err.Error())
		return nil, err
	}
	addRequest := client.NewAddRegularExpressionPasswordValidatorRequest([]client.EnumregularExpressionPasswordValidatorSchemaUrn{client.ENUMREGULAREXPRESSIONPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORREGULAR_EXPRESSION},
		plan.MatchPattern.ValueString(),
		*matchBehavior,
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err = addOptionalRegularExpressionPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddRegularExpressionPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readRegularExpressionPasswordValidatorResponse(ctx, addResponse.RegularExpressionPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a unique-characters password-validator
func (r *passwordValidatorResource) CreateUniqueCharactersPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddUniqueCharactersPasswordValidatorRequest([]client.EnumuniqueCharactersPasswordValidatorSchemaUrn{client.ENUMUNIQUECHARACTERSPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORUNIQUE_CHARACTERS},
		plan.MinUniqueCharacters.ValueInt64(),
		plan.CaseSensitiveValidation.ValueBool(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalUniqueCharactersPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddUniqueCharactersPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readUniqueCharactersPasswordValidatorResponse(ctx, addResponse.UniqueCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party password-validator
func (r *passwordValidatorResource) CreateThirdPartyPasswordValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordValidatorResourceModel) (*passwordValidatorResourceModel, error) {
	addRequest := client.NewAddThirdPartyPasswordValidatorRequest([]client.EnumthirdPartyPasswordValidatorSchemaUrn{client.ENUMTHIRDPARTYPASSWORDVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_VALIDATORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyPasswordValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordValidatorAPI.AddPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordValidatorRequest(
		client.AddThirdPartyPasswordValidatorRequestAsAddPasswordValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.AddPasswordValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordValidatorResourceModel
	readThirdPartyPasswordValidatorResponse(ctx, addResponse.ThirdPartyPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passwordValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passwordValidatorResourceModel
	var err error
	if plan.Type.ValueString() == "character-set" {
		state, err = r.CreateCharacterSetPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "similarity-based" {
		state, err = r.CreateSimilarityBasedPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "attribute-value" {
		state, err = r.CreateAttributeValuePasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "repeated-characters" {
		state, err = r.CreateRepeatedCharactersPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "dictionary" {
		state, err = r.CreateDictionaryPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "haystack" {
		state, err = r.CreateHaystackPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "utf-8" {
		state, err = r.CreateUtf8PasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "pwned-passwords" {
		state, err = r.CreatePwnedPasswordsPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "disallowed-characters" {
		state, err = r.CreateDisallowedCharactersPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "length-based" {
		state, err = r.CreateLengthBasedPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "regular-expression" {
		state, err = r.CreateRegularExpressionPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "unique-characters" {
		state, err = r.CreateUniqueCharactersPasswordValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPasswordValidator(ctx, req, resp, plan)
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
func (r *defaultPasswordValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.GetPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passwordValidatorResourceModel
	if readResponse.CharacterSetPasswordValidatorResponse != nil {
		readCharacterSetPasswordValidatorResponse(ctx, readResponse.CharacterSetPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimilarityBasedPasswordValidatorResponse != nil {
		readSimilarityBasedPasswordValidatorResponse(ctx, readResponse.SimilarityBasedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeValuePasswordValidatorResponse != nil {
		readAttributeValuePasswordValidatorResponse(ctx, readResponse.AttributeValuePasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomPasswordValidatorResponse != nil {
		readCustomPasswordValidatorResponse(ctx, readResponse.CustomPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RepeatedCharactersPasswordValidatorResponse != nil {
		readRepeatedCharactersPasswordValidatorResponse(ctx, readResponse.RepeatedCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DictionaryPasswordValidatorResponse != nil {
		readDictionaryPasswordValidatorResponse(ctx, readResponse.DictionaryPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HaystackPasswordValidatorResponse != nil {
		readHaystackPasswordValidatorResponse(ctx, readResponse.HaystackPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Utf8PasswordValidatorResponse != nil {
		readUtf8PasswordValidatorResponse(ctx, readResponse.Utf8PasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordValidatorResponse != nil {
		readGroovyScriptedPasswordValidatorResponse(ctx, readResponse.GroovyScriptedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PwnedPasswordsPasswordValidatorResponse != nil {
		readPwnedPasswordsPasswordValidatorResponse(ctx, readResponse.PwnedPasswordsPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DisallowedCharactersPasswordValidatorResponse != nil {
		readDisallowedCharactersPasswordValidatorResponse(ctx, readResponse.DisallowedCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LengthBasedPasswordValidatorResponse != nil {
		readLengthBasedPasswordValidatorResponse(ctx, readResponse.LengthBasedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionPasswordValidatorResponse != nil {
		readRegularExpressionPasswordValidatorResponse(ctx, readResponse.RegularExpressionPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UniqueCharactersPasswordValidatorResponse != nil {
		readUniqueCharactersPasswordValidatorResponse(ctx, readResponse.UniqueCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordValidatorResponse != nil {
		readThirdPartyPasswordValidatorResponse(ctx, readResponse.ThirdPartyPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PasswordValidatorAPI.UpdatePasswordValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPasswordValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PasswordValidatorAPI.UpdatePasswordValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CharacterSetPasswordValidatorResponse != nil {
			readCharacterSetPasswordValidatorResponse(ctx, updateResponse.CharacterSetPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimilarityBasedPasswordValidatorResponse != nil {
			readSimilarityBasedPasswordValidatorResponse(ctx, updateResponse.SimilarityBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AttributeValuePasswordValidatorResponse != nil {
			readAttributeValuePasswordValidatorResponse(ctx, updateResponse.AttributeValuePasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomPasswordValidatorResponse != nil {
			readCustomPasswordValidatorResponse(ctx, updateResponse.CustomPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RepeatedCharactersPasswordValidatorResponse != nil {
			readRepeatedCharactersPasswordValidatorResponse(ctx, updateResponse.RepeatedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DictionaryPasswordValidatorResponse != nil {
			readDictionaryPasswordValidatorResponse(ctx, updateResponse.DictionaryPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HaystackPasswordValidatorResponse != nil {
			readHaystackPasswordValidatorResponse(ctx, updateResponse.HaystackPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Utf8PasswordValidatorResponse != nil {
			readUtf8PasswordValidatorResponse(ctx, updateResponse.Utf8PasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedPasswordValidatorResponse != nil {
			readGroovyScriptedPasswordValidatorResponse(ctx, updateResponse.GroovyScriptedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PwnedPasswordsPasswordValidatorResponse != nil {
			readPwnedPasswordsPasswordValidatorResponse(ctx, updateResponse.PwnedPasswordsPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DisallowedCharactersPasswordValidatorResponse != nil {
			readDisallowedCharactersPasswordValidatorResponse(ctx, updateResponse.DisallowedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LengthBasedPasswordValidatorResponse != nil {
			readLengthBasedPasswordValidatorResponse(ctx, updateResponse.LengthBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RegularExpressionPasswordValidatorResponse != nil {
			readRegularExpressionPasswordValidatorResponse(ctx, updateResponse.RegularExpressionPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UniqueCharactersPasswordValidatorResponse != nil {
			readUniqueCharactersPasswordValidatorResponse(ctx, updateResponse.UniqueCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordValidatorResponse != nil {
			readThirdPartyPasswordValidatorResponse(ctx, updateResponse.ThirdPartyPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *passwordValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPasswordValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPasswordValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state passwordValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PasswordValidatorAPI.GetPasswordValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Password Validator", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Validator", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.CharacterSetPasswordValidatorResponse != nil {
		readCharacterSetPasswordValidatorResponse(ctx, readResponse.CharacterSetPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimilarityBasedPasswordValidatorResponse != nil {
		readSimilarityBasedPasswordValidatorResponse(ctx, readResponse.SimilarityBasedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeValuePasswordValidatorResponse != nil {
		readAttributeValuePasswordValidatorResponse(ctx, readResponse.AttributeValuePasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomPasswordValidatorResponse != nil {
		readCustomPasswordValidatorResponse(ctx, readResponse.CustomPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RepeatedCharactersPasswordValidatorResponse != nil {
		readRepeatedCharactersPasswordValidatorResponse(ctx, readResponse.RepeatedCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DictionaryPasswordValidatorResponse != nil {
		readDictionaryPasswordValidatorResponse(ctx, readResponse.DictionaryPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HaystackPasswordValidatorResponse != nil {
		readHaystackPasswordValidatorResponse(ctx, readResponse.HaystackPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Utf8PasswordValidatorResponse != nil {
		readUtf8PasswordValidatorResponse(ctx, readResponse.Utf8PasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordValidatorResponse != nil {
		readGroovyScriptedPasswordValidatorResponse(ctx, readResponse.GroovyScriptedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PwnedPasswordsPasswordValidatorResponse != nil {
		readPwnedPasswordsPasswordValidatorResponse(ctx, readResponse.PwnedPasswordsPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DisallowedCharactersPasswordValidatorResponse != nil {
		readDisallowedCharactersPasswordValidatorResponse(ctx, readResponse.DisallowedCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LengthBasedPasswordValidatorResponse != nil {
		readLengthBasedPasswordValidatorResponse(ctx, readResponse.LengthBasedPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionPasswordValidatorResponse != nil {
		readRegularExpressionPasswordValidatorResponse(ctx, readResponse.RegularExpressionPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.UniqueCharactersPasswordValidatorResponse != nil {
		readUniqueCharactersPasswordValidatorResponse(ctx, readResponse.UniqueCharactersPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordValidatorResponse != nil {
		readThirdPartyPasswordValidatorResponse(ctx, readResponse.ThirdPartyPasswordValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passwordValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePasswordValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passwordValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passwordValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PasswordValidatorAPI.UpdatePasswordValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPasswordValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PasswordValidatorAPI.UpdatePasswordValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CharacterSetPasswordValidatorResponse != nil {
			readCharacterSetPasswordValidatorResponse(ctx, updateResponse.CharacterSetPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimilarityBasedPasswordValidatorResponse != nil {
			readSimilarityBasedPasswordValidatorResponse(ctx, updateResponse.SimilarityBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AttributeValuePasswordValidatorResponse != nil {
			readAttributeValuePasswordValidatorResponse(ctx, updateResponse.AttributeValuePasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomPasswordValidatorResponse != nil {
			readCustomPasswordValidatorResponse(ctx, updateResponse.CustomPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RepeatedCharactersPasswordValidatorResponse != nil {
			readRepeatedCharactersPasswordValidatorResponse(ctx, updateResponse.RepeatedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DictionaryPasswordValidatorResponse != nil {
			readDictionaryPasswordValidatorResponse(ctx, updateResponse.DictionaryPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HaystackPasswordValidatorResponse != nil {
			readHaystackPasswordValidatorResponse(ctx, updateResponse.HaystackPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Utf8PasswordValidatorResponse != nil {
			readUtf8PasswordValidatorResponse(ctx, updateResponse.Utf8PasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedPasswordValidatorResponse != nil {
			readGroovyScriptedPasswordValidatorResponse(ctx, updateResponse.GroovyScriptedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PwnedPasswordsPasswordValidatorResponse != nil {
			readPwnedPasswordsPasswordValidatorResponse(ctx, updateResponse.PwnedPasswordsPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DisallowedCharactersPasswordValidatorResponse != nil {
			readDisallowedCharactersPasswordValidatorResponse(ctx, updateResponse.DisallowedCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LengthBasedPasswordValidatorResponse != nil {
			readLengthBasedPasswordValidatorResponse(ctx, updateResponse.LengthBasedPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RegularExpressionPasswordValidatorResponse != nil {
			readRegularExpressionPasswordValidatorResponse(ctx, updateResponse.RegularExpressionPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.UniqueCharactersPasswordValidatorResponse != nil {
			readUniqueCharactersPasswordValidatorResponse(ctx, updateResponse.UniqueCharactersPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordValidatorResponse != nil {
			readThirdPartyPasswordValidatorResponse(ctx, updateResponse.ThirdPartyPasswordValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPasswordValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passwordValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passwordValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PasswordValidatorAPI.DeletePasswordValidatorExecute(r.apiClient.PasswordValidatorAPI.DeletePasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Password Validator", err, httpResp)
		return
	}
}

func (r *passwordValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordValidator(ctx, req, resp)
}

func (r *defaultPasswordValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordValidator(ctx, req, resp)
}

func importPasswordValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
