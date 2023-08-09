package passwordvalidator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passwordValidatorDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordValidatorDataSource{}
)

// Create a Password Validator data source
func NewPasswordValidatorDataSource() datasource.DataSource {
	return &passwordValidatorDataSource{}
}

// passwordValidatorDataSource is the datasource implementation.
type passwordValidatorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordValidatorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_validator"
}

// Configure adds the provider configured client to the data source.
func (r *passwordValidatorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type passwordValidatorDataSourceModel struct {
	Id                                             types.String `tfsdk:"id"`
	Name                                           types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *passwordValidatorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Password Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Validator resource. Options are ['character-set', 'similarity-based', 'attribute-value', 'custom', 'repeated-characters', 'dictionary', 'haystack', 'utf-8', 'groovy-scripted', 'pwned-passwords', 'disallowed-characters', 'length-based', 'regular-expression', 'unique-characters', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Password Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Password Validator. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"min_unique_characters": schema.Int64Attribute{
				Description: "Specifies the minimum number of unique characters that a password will be allowed to contain.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"match_pattern": schema.StringAttribute{
				Description: "The regular expression to use for this password validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"match_behavior": schema.StringAttribute{
				Description: "The behavior to exhibit if a user's proposed password matches the regular expression defined in the match-pattern property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_password_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of characters that can be included in a proposed password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_password_length": schema.Int64Attribute{
				Description: "Specifies the minimum number of characters that must be included in a proposed password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disallowed_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed anywhere in a password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disallowed_leading_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed as the first character of the password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disallowed_trailing_characters": schema.StringAttribute{
				Description: "A set of characters that will not be allowed as the last character of the password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pwned_passwords_base_url": schema.StringAttribute{
				Description: "The base URL for requests used to interact with the Pwned Passwords service. The first five characters of the hexadecimal representation of the unsalted SHA-1 digest of a proposed password will be appended to this base URL to construct the HTTP GET request used to obtain information about potential matches.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Pwned Passwords service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invoke_for_add": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided in LDAP add requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invoke_for_self_change": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided by an end user in the course of changing their own password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"invoke_for_admin_reset": schema.BoolAttribute{
				Description: "Indicates whether this password validator should be used to validate clear-text passwords provided by administrators when changing the password for another user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"accept_password_on_service_error": schema.BoolAttribute{
				Description: "Indicates whether to accept the proposed password if an error occurs while attempting to interact with the Pwned Passwords service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies which key manager provider should be used to obtain a client certificate to present to the validation server when performing HTTPS communication. This may be left undefined if communication will not be secured with HTTPS, or if there is no need to present a client certificate to the validation service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies which trust manager provider should be used to determine whether to trust the certificate presented by the server when performing HTTPS communication. This may be left undefined if HTTPS communication is not needed, or if the validation service presents a certificate that is trusted by the default JVM configuration (which should be the case for the Pwned Password servers).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Password Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Password Validator. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allow_non_ascii_characters": schema.BoolAttribute{
				Description: "Indicates whether passwords will be allowed to include characters from outside the ASCII character set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_unknown_characters": schema.BoolAttribute{
				Description: "Indicates whether passwords will be allowed to include characters that are not recognized by the JVM's Unicode support.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_character_type": schema.SetAttribute{
				Description: "Specifies the set of character types that are allowed to be present in passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"assumed_password_guesses_per_second": schema.StringAttribute{
				Description: "The number of password guesses per second that a potential attacker may be expected to make.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_acceptable_time_to_exhaust_search_space": schema.StringAttribute{
				Description: "The minimum length of time (using the configured number of password guesses per second) required to exhaust the entire search space for a proposed password in order for that password to be considered acceptable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"dictionary_file": schema.StringAttribute{
				Description: "Specifies the path to the file containing a list of words that cannot be used as passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_consecutive_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of times that any character can appear consecutively in a password value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"case_sensitive_validation": schema.BoolAttribute{
				Description: " When the `type` attribute is set to one of [`repeated-characters`, `unique-characters`]: Indicates whether this password validator should treat password characters in a case-sensitive manner. When the `type` attribute is set to `dictionary`: Indicates whether this password validator is to treat password characters in a case-sensitive manner.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ignore_leading_non_alphabetic_characters": schema.BoolAttribute{
				Description: "Indicates whether to ignore any digits, symbols, or other non-alphabetic characters that may appear at the beginning of a proposed password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ignore_trailing_non_alphabetic_characters": schema.BoolAttribute{
				Description: "Indicates whether to ignore any digits, symbols, or other non-alphabetic characters that may appear at the end of a proposed password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"strip_diacritical_marks": schema.BoolAttribute{
				Description: "Indicates whether to strip characters of any diacritical marks (like accents, cedillas, circumflexes, diaereses, tildes, and umlauts) they may contain. Any characters with a diacritical mark would be replaced with a base version",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alternative_password_character_mapping": schema.SetAttribute{
				Description: "Provides a set of character substitutions that can be applied to the proposed password when checking to see if it is in the provided dictionary. Each mapping should consist of a single character followed by a colon and a list of the alternative characters that may be used in place of that character.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"maximum_allowed_percent_of_password": schema.Int64Attribute{
				Description: "The maximum allowed percent of a proposed password that any single dictionary word is allowed to comprise. A value of 100 indicates that a proposed password will only be rejected if the dictionary contains the entire proposed password (after any configured transformations have been applied).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"match_attribute": schema.SetAttribute{
				Description: "Specifies the name(s) of the attribute(s) whose values should be checked to determine whether they match the provided password. If no values are provided, then the server checks if the proposed password matches the value of any user attribute in the target user's entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"test_password_substring_of_attribute_value": schema.BoolAttribute{
				Description: "Indicates whether to reject any proposed password that is a substring of a value in one of the match attributes in the target user's entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"test_attribute_value_substring_of_password": schema.BoolAttribute{
				Description: "Indicates whether to reject any proposed password in which a value in one of the match attributes in the target user's entry is a substring of that password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_attribute_value_length_for_substring_matches": schema.Int64Attribute{
				Description: "The minimum length that an attribute value must have for it to be considered when rejecting passwords that contain the value of another attribute as a substring.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"test_reversed_password": schema.BoolAttribute{
				Description: " When the `type` attribute is set to `attribute-value`: Indicates whether to perform matching against the reversed value of the provided password in addition to the order in which it was given. When the `type` attribute is set to `dictionary`: Indicates whether this password validator is to test the reversed value of the provided password as well as the order in which it was given.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_password_difference": schema.Int64Attribute{
				Description: "Specifies the minimum difference of new and old password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"character_set": schema.SetAttribute{
				Description: " When the `type` attribute is set to `character-set`: Specifies a character set containing characters that a password may contain and a value indicating the minimum number of characters required from that set. When the `type` attribute is set to `repeated-characters`: Specifies a set of characters that should be considered equivalent for the purpose of this password validator. This can be used, for example, to ensure that passwords contain no more than three consecutive digits.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allow_unclassified_characters": schema.BoolAttribute{
				Description: "Indicates whether this password validator allows passwords to contain characters outside of any of the user-defined character sets.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_required_character_sets": schema.Int64Attribute{
				Description: "Specifies the minimum number of character sets that must be represented in a proposed password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Validator",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the password validator is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"validator_requirement_description": schema.StringAttribute{
				Description: "Specifies a message that can be used to describe the requirements imposed by this password validator to end users. If a value is provided for this property, then it will override any description that may have otherwise been generated by the validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"validator_failure_message": schema.StringAttribute{
				Description: "Specifies a message that may be provided to the end user in the event that a proposed password is rejected by this validator. If a value is provided for this property, then it will override any failure message that may have otherwise been generated by the validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a CharacterSetPasswordValidatorResponse object into the model struct
func readCharacterSetPasswordValidatorResponseDataSource(ctx context.Context, r *client.CharacterSetPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("character-set")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CharacterSet = internaltypes.GetStringSet(r.CharacterSet)
	state.AllowUnclassifiedCharacters = types.BoolValue(r.AllowUnclassifiedCharacters)
	state.MinimumRequiredCharacterSets = internaltypes.Int64TypeOrNil(r.MinimumRequiredCharacterSets)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a SimilarityBasedPasswordValidatorResponse object into the model struct
func readSimilarityBasedPasswordValidatorResponseDataSource(ctx context.Context, r *client.SimilarityBasedPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("similarity-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MinPasswordDifference = types.Int64Value(r.MinPasswordDifference)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a AttributeValuePasswordValidatorResponse object into the model struct
func readAttributeValuePasswordValidatorResponseDataSource(ctx context.Context, r *client.AttributeValuePasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-value")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.TestPasswordSubstringOfAttributeValue = internaltypes.BoolTypeOrNil(r.TestPasswordSubstringOfAttributeValue)
	state.TestAttributeValueSubstringOfPassword = internaltypes.BoolTypeOrNil(r.TestAttributeValueSubstringOfPassword)
	state.MinimumAttributeValueLengthForSubstringMatches = internaltypes.Int64TypeOrNil(r.MinimumAttributeValueLengthForSubstringMatches)
	state.TestReversedPassword = types.BoolValue(r.TestReversedPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a CustomPasswordValidatorResponse object into the model struct
func readCustomPasswordValidatorResponseDataSource(ctx context.Context, r *client.CustomPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a RepeatedCharactersPasswordValidatorResponse object into the model struct
func readRepeatedCharactersPasswordValidatorResponseDataSource(ctx context.Context, r *client.RepeatedCharactersPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("repeated-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxConsecutiveLength = types.Int64Value(r.MaxConsecutiveLength)
	state.CaseSensitiveValidation = types.BoolValue(r.CaseSensitiveValidation)
	state.CharacterSet = internaltypes.GetStringSet(r.CharacterSet)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a DictionaryPasswordValidatorResponse object into the model struct
func readDictionaryPasswordValidatorResponseDataSource(ctx context.Context, r *client.DictionaryPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a HaystackPasswordValidatorResponse object into the model struct
func readHaystackPasswordValidatorResponseDataSource(ctx context.Context, r *client.HaystackPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("haystack")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AssumedPasswordGuessesPerSecond = types.StringValue(r.AssumedPasswordGuessesPerSecond)
	state.MinimumAcceptableTimeToExhaustSearchSpace = types.StringValue(r.MinimumAcceptableTimeToExhaustSearchSpace)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a Utf8PasswordValidatorResponse object into the model struct
func readUtf8PasswordValidatorResponseDataSource(ctx context.Context, r *client.Utf8PasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("utf-8")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowNonAsciiCharacters = internaltypes.BoolTypeOrNil(r.AllowNonAsciiCharacters)
	state.AllowUnknownCharacters = internaltypes.BoolTypeOrNil(r.AllowUnknownCharacters)
	state.AllowedCharacterType = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordValidatorAllowedCharacterTypeProp(r.AllowedCharacterType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a GroovyScriptedPasswordValidatorResponse object into the model struct
func readGroovyScriptedPasswordValidatorResponseDataSource(ctx context.Context, r *client.GroovyScriptedPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a PwnedPasswordsPasswordValidatorResponse object into the model struct
func readPwnedPasswordsPasswordValidatorResponseDataSource(ctx context.Context, r *client.PwnedPasswordsPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pwned-passwords")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PwnedPasswordsBaseURL = types.StringValue(r.PwnedPasswordsBaseURL)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.InvokeForAdd = types.BoolValue(r.InvokeForAdd)
	state.InvokeForSelfChange = types.BoolValue(r.InvokeForSelfChange)
	state.InvokeForAdminReset = types.BoolValue(r.InvokeForAdminReset)
	state.AcceptPasswordOnServiceError = types.BoolValue(r.AcceptPasswordOnServiceError)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a DisallowedCharactersPasswordValidatorResponse object into the model struct
func readDisallowedCharactersPasswordValidatorResponseDataSource(ctx context.Context, r *client.DisallowedCharactersPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disallowed-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DisallowedCharacters = internaltypes.StringTypeOrNil(r.DisallowedCharacters, false)
	state.DisallowedLeadingCharacters = internaltypes.StringTypeOrNil(r.DisallowedLeadingCharacters, false)
	state.DisallowedTrailingCharacters = internaltypes.StringTypeOrNil(r.DisallowedTrailingCharacters, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a LengthBasedPasswordValidatorResponse object into the model struct
func readLengthBasedPasswordValidatorResponseDataSource(ctx context.Context, r *client.LengthBasedPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("length-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.MinPasswordLength = internaltypes.Int64TypeOrNil(r.MinPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a RegularExpressionPasswordValidatorResponse object into the model struct
func readRegularExpressionPasswordValidatorResponseDataSource(ctx context.Context, r *client.RegularExpressionPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("regular-expression")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchPattern = types.StringValue(r.MatchPattern)
	state.MatchBehavior = types.StringValue(r.MatchBehavior.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a UniqueCharactersPasswordValidatorResponse object into the model struct
func readUniqueCharactersPasswordValidatorResponseDataSource(ctx context.Context, r *client.UniqueCharactersPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("unique-characters")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MinUniqueCharacters = types.Int64Value(r.MinUniqueCharacters)
	state.CaseSensitiveValidation = types.BoolValue(r.CaseSensitiveValidation)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read a ThirdPartyPasswordValidatorResponse object into the model struct
func readThirdPartyPasswordValidatorResponseDataSource(ctx context.Context, r *client.ThirdPartyPasswordValidatorResponse, state *passwordValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ValidatorRequirementDescription = internaltypes.StringTypeOrNil(r.ValidatorRequirementDescription, false)
	state.ValidatorFailureMessage = internaltypes.StringTypeOrNil(r.ValidatorFailureMessage, false)
}

// Read resource information
func (r *passwordValidatorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordValidatorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordValidatorApi.GetPasswordValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.CharacterSetPasswordValidatorResponse != nil {
		readCharacterSetPasswordValidatorResponseDataSource(ctx, readResponse.CharacterSetPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SimilarityBasedPasswordValidatorResponse != nil {
		readSimilarityBasedPasswordValidatorResponseDataSource(ctx, readResponse.SimilarityBasedPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeValuePasswordValidatorResponse != nil {
		readAttributeValuePasswordValidatorResponseDataSource(ctx, readResponse.AttributeValuePasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomPasswordValidatorResponse != nil {
		readCustomPasswordValidatorResponseDataSource(ctx, readResponse.CustomPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.RepeatedCharactersPasswordValidatorResponse != nil {
		readRepeatedCharactersPasswordValidatorResponseDataSource(ctx, readResponse.RepeatedCharactersPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DictionaryPasswordValidatorResponse != nil {
		readDictionaryPasswordValidatorResponseDataSource(ctx, readResponse.DictionaryPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HaystackPasswordValidatorResponse != nil {
		readHaystackPasswordValidatorResponseDataSource(ctx, readResponse.HaystackPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Utf8PasswordValidatorResponse != nil {
		readUtf8PasswordValidatorResponseDataSource(ctx, readResponse.Utf8PasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordValidatorResponse != nil {
		readGroovyScriptedPasswordValidatorResponseDataSource(ctx, readResponse.GroovyScriptedPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PwnedPasswordsPasswordValidatorResponse != nil {
		readPwnedPasswordsPasswordValidatorResponseDataSource(ctx, readResponse.PwnedPasswordsPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DisallowedCharactersPasswordValidatorResponse != nil {
		readDisallowedCharactersPasswordValidatorResponseDataSource(ctx, readResponse.DisallowedCharactersPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LengthBasedPasswordValidatorResponse != nil {
		readLengthBasedPasswordValidatorResponseDataSource(ctx, readResponse.LengthBasedPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionPasswordValidatorResponse != nil {
		readRegularExpressionPasswordValidatorResponseDataSource(ctx, readResponse.RegularExpressionPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UniqueCharactersPasswordValidatorResponse != nil {
		readUniqueCharactersPasswordValidatorResponseDataSource(ctx, readResponse.UniqueCharactersPasswordValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordValidatorResponse != nil {
		readThirdPartyPasswordValidatorResponseDataSource(ctx, readResponse.ThirdPartyPasswordValidatorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
