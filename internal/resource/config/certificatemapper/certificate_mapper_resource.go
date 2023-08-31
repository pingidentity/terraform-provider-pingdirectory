package certificatemapper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &certificateMapperResource{}
	_ resource.ResourceWithConfigure   = &certificateMapperResource{}
	_ resource.ResourceWithImportState = &certificateMapperResource{}
	_ resource.Resource                = &defaultCertificateMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultCertificateMapperResource{}
	_ resource.ResourceWithImportState = &defaultCertificateMapperResource{}
)

// Create a Certificate Mapper resource
func NewCertificateMapperResource() resource.Resource {
	return &certificateMapperResource{}
}

func NewDefaultCertificateMapperResource() resource.Resource {
	return &defaultCertificateMapperResource{}
}

// certificateMapperResource is the resource implementation.
type certificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCertificateMapperResource is the resource implementation.
type defaultCertificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *certificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_mapper"
}

func (r *defaultCertificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_certificate_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *certificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificateMapperResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	Type                    types.String `tfsdk:"type"`
	ExtensionClass          types.String `tfsdk:"extension_class"`
	ExtensionArgument       types.Set    `tfsdk:"extension_argument"`
	FingerprintAttribute    types.String `tfsdk:"fingerprint_attribute"`
	FingerprintAlgorithm    types.String `tfsdk:"fingerprint_algorithm"`
	SubjectAttributeMapping types.Set    `tfsdk:"subject_attribute_mapping"`
	ScriptClass             types.String `tfsdk:"script_class"`
	ScriptArgument          types.Set    `tfsdk:"script_argument"`
	SubjectAttribute        types.String `tfsdk:"subject_attribute"`
	UserBaseDN              types.Set    `tfsdk:"user_base_dn"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *certificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificateMapperSchema(ctx, req, resp, false)
}

func (r *defaultCertificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificateMapperSchema(ctx, req, resp, true)
}

func certificateMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Certificate Mapper.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Certificate Mapper resource. Options are ['subject-equals-dn', 'subject-dn-to-user-attribute', 'groovy-scripted', 'subject-attribute-to-user-attribute', 'fingerprint', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"subject-equals-dn", "subject-dn-to-user-attribute", "groovy-scripted", "subject-attribute-to-user-attribute", "fingerprint", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Certificate Mapper.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"fingerprint_attribute": schema.StringAttribute{
				Description: "Specifies the attribute in which to look for the fingerprint.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"fingerprint_algorithm": schema.StringAttribute{
				Description: "Specifies the name of the digest algorithm to compute the fingerprint of client certificates.",
				Optional:    true,
			},
			"subject_attribute_mapping": schema.SetAttribute{
				Description: "Specifies a mapping between certificate attributes and user attributes.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Certificate Mapper.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"subject_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute whose value should exactly match the certificate subject DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_base_dn": schema.SetAttribute{
				Description:         "When the `type` attribute is set to  one of [`subject-dn-to-user-attribute`, `subject-attribute-to-user-attribute`]: Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry. When the `type` attribute is set to `fingerprint`: Specifies the set of base DNs below which to search for users.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`subject-dn-to-user-attribute`, `subject-attribute-to-user-attribute`]: Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry.\n  - `fingerprint`: Specifies the set of base DNs below which to search for users.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Certificate Mapper",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Certificate Mapper is enabled.",
				Required:    true,
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
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *certificateMapperResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model certificateMapperResourceModel
	req.Plan.Get(ctx, &model)
	resourceType := model.Type.ValueString()
	// Set defaults for subject-dn-to-user-attribute type
	if resourceType == "subject-dn-to-user-attribute" {
		if !internaltypes.IsDefined(model.SubjectAttribute) {
			model.SubjectAttribute = types.StringValue("ds-certificate-subject-dn")
		}
	}
	// Set defaults for fingerprint type
	if resourceType == "fingerprint" {
		if !internaltypes.IsDefined(model.FingerprintAttribute) {
			model.FingerprintAttribute = types.StringValue("ds-certificate-fingerprint")
		}
	}
	resp.Plan.Set(ctx, &model)
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsCertificateMapper() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("subject_attribute"),
			path.MatchRoot("type"),
			[]string{"subject-dn-to-user-attribute"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("user_base_dn"),
			path.MatchRoot("type"),
			[]string{"subject-dn-to-user-attribute", "subject-attribute-to-user-attribute", "fingerprint"},
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
			path.MatchRoot("subject_attribute_mapping"),
			path.MatchRoot("type"),
			[]string{"subject-attribute-to-user-attribute"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("fingerprint_attribute"),
			path.MatchRoot("type"),
			[]string{"fingerprint"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("fingerprint_algorithm"),
			path.MatchRoot("type"),
			[]string{"fingerprint"},
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
			"groovy-scripted",
			[]path.Expression{path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"subject-attribute-to-user-attribute",
			[]path.Expression{path.MatchRoot("subject_attribute_mapping")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"fingerprint",
			[]path.Expression{path.MatchRoot("fingerprint_algorithm")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r certificateMapperResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsCertificateMapper()
}

// Add config validators
func (r defaultCertificateMapperResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsCertificateMapper()
}

// Add optional fields to create request for subject-equals-dn certificate-mapper
func addOptionalSubjectEqualsDnCertificateMapperFields(ctx context.Context, addRequest *client.AddSubjectEqualsDnCertificateMapperRequest, plan certificateMapperResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for subject-dn-to-user-attribute certificate-mapper
func addOptionalSubjectDnToUserAttributeCertificateMapperFields(ctx context.Context, addRequest *client.AddSubjectDnToUserAttributeCertificateMapperRequest, plan certificateMapperResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectAttribute) {
		addRequest.SubjectAttribute = plan.SubjectAttribute.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.UserBaseDN) {
		var slice []string
		plan.UserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.UserBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted certificate-mapper
func addOptionalGroovyScriptedCertificateMapperFields(ctx context.Context, addRequest *client.AddGroovyScriptedCertificateMapperRequest, plan certificateMapperResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for subject-attribute-to-user-attribute certificate-mapper
func addOptionalSubjectAttributeToUserAttributeCertificateMapperFields(ctx context.Context, addRequest *client.AddSubjectAttributeToUserAttributeCertificateMapperRequest, plan certificateMapperResourceModel) {
	if internaltypes.IsDefined(plan.UserBaseDN) {
		var slice []string
		plan.UserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.UserBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for fingerprint certificate-mapper
func addOptionalFingerprintCertificateMapperFields(ctx context.Context, addRequest *client.AddFingerprintCertificateMapperRequest, plan certificateMapperResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FingerprintAttribute) {
		addRequest.FingerprintAttribute = plan.FingerprintAttribute.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.UserBaseDN) {
		var slice []string
		plan.UserBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.UserBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party certificate-mapper
func addOptionalThirdPartyCertificateMapperFields(ctx context.Context, addRequest *client.AddThirdPartyCertificateMapperRequest, plan certificateMapperResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateCertificateMapperUnknownValues(model *certificateMapperResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SubjectAttributeMapping.IsUnknown() || model.SubjectAttributeMapping.IsNull() {
		model.SubjectAttributeMapping, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.UserBaseDN.IsUnknown() || model.UserBaseDN.IsNull() {
		model.UserBaseDN, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SubjectAttribute.IsUnknown() || model.SubjectAttribute.IsNull() {
		model.SubjectAttribute = types.StringValue("")
	}
	if model.FingerprintAttribute.IsUnknown() || model.FingerprintAttribute.IsNull() {
		model.FingerprintAttribute = types.StringValue("")
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *certificateMapperResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.FingerprintAlgorithm.IsUnknown() || model.FingerprintAlgorithm.IsNull() {
		model.FingerprintAlgorithm = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
}

// Read a SubjectEqualsDnCertificateMapperResponse object into the model struct
func readSubjectEqualsDnCertificateMapperResponse(ctx context.Context, r *client.SubjectEqualsDnCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-equals-dn")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Read a SubjectDnToUserAttributeCertificateMapperResponse object into the model struct
func readSubjectDnToUserAttributeCertificateMapperResponse(ctx context.Context, r *client.SubjectDnToUserAttributeCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-dn-to-user-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubjectAttribute = types.StringValue(r.SubjectAttribute)
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Read a GroovyScriptedCertificateMapperResponse object into the model struct
func readGroovyScriptedCertificateMapperResponse(ctx context.Context, r *client.GroovyScriptedCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Read a SubjectAttributeToUserAttributeCertificateMapperResponse object into the model struct
func readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx context.Context, r *client.SubjectAttributeToUserAttributeCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-attribute-to-user-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubjectAttributeMapping = internaltypes.GetStringSet(r.SubjectAttributeMapping)
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Read a FingerprintCertificateMapperResponse object into the model struct
func readFingerprintCertificateMapperResponse(ctx context.Context, r *client.FingerprintCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fingerprint")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FingerprintAttribute = types.StringValue(r.FingerprintAttribute)
	state.FingerprintAlgorithm = types.StringValue(r.FingerprintAlgorithm.String())
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Read a ThirdPartyCertificateMapperResponse object into the model struct
func readThirdPartyCertificateMapperResponse(ctx context.Context, r *client.ThirdPartyCertificateMapperResponse, state *certificateMapperResourceModel, expectedValues *certificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateCertificateMapperUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createCertificateMapperOperations(plan certificateMapperResourceModel, state certificateMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.FingerprintAttribute, state.FingerprintAttribute, "fingerprint-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.FingerprintAlgorithm, state.FingerprintAlgorithm, "fingerprint-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SubjectAttributeMapping, state.SubjectAttributeMapping, "subject-attribute-mapping")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectAttribute, state.SubjectAttribute, "subject-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserBaseDN, state.UserBaseDN, "user-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a subject-equals-dn certificate-mapper
func (r *certificateMapperResource) CreateSubjectEqualsDnCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	addRequest := client.NewAddSubjectEqualsDnCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumsubjectEqualsDnCertificateMapperSchemaUrn{client.ENUMSUBJECTEQUALSDNCERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERSUBJECT_EQUALS_DN},
		plan.Enabled.ValueBool())
	addOptionalSubjectEqualsDnCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddSubjectEqualsDnCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readSubjectEqualsDnCertificateMapperResponse(ctx, addResponse.SubjectEqualsDnCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a subject-dn-to-user-attribute certificate-mapper
func (r *certificateMapperResource) CreateSubjectDnToUserAttributeCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	addRequest := client.NewAddSubjectDnToUserAttributeCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumsubjectDnToUserAttributeCertificateMapperSchemaUrn{client.ENUMSUBJECTDNTOUSERATTRIBUTECERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERSUBJECT_DN_TO_USER_ATTRIBUTE},
		plan.Enabled.ValueBool())
	addOptionalSubjectDnToUserAttributeCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddSubjectDnToUserAttributeCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readSubjectDnToUserAttributeCertificateMapperResponse(ctx, addResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted certificate-mapper
func (r *certificateMapperResource) CreateGroovyScriptedCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedCertificateMapperSchemaUrn{client.ENUMGROOVYSCRIPTEDCERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddGroovyScriptedCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readGroovyScriptedCertificateMapperResponse(ctx, addResponse.GroovyScriptedCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a subject-attribute-to-user-attribute certificate-mapper
func (r *certificateMapperResource) CreateSubjectAttributeToUserAttributeCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	var SubjectAttributeMappingSlice []string
	plan.SubjectAttributeMapping.ElementsAs(ctx, &SubjectAttributeMappingSlice, false)
	addRequest := client.NewAddSubjectAttributeToUserAttributeCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumsubjectAttributeToUserAttributeCertificateMapperSchemaUrn{client.ENUMSUBJECTATTRIBUTETOUSERATTRIBUTECERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERSUBJECT_ATTRIBUTE_TO_USER_ATTRIBUTE},
		SubjectAttributeMappingSlice,
		plan.Enabled.ValueBool())
	addOptionalSubjectAttributeToUserAttributeCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddSubjectAttributeToUserAttributeCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, addResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a fingerprint certificate-mapper
func (r *certificateMapperResource) CreateFingerprintCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	fingerprintAlgorithm, err := client.NewEnumcertificateMapperFingerprintAlgorithmPropFromValue(plan.FingerprintAlgorithm.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for FingerprintAlgorithm", err.Error())
		return nil, err
	}
	addRequest := client.NewAddFingerprintCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumfingerprintCertificateMapperSchemaUrn{client.ENUMFINGERPRINTCERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERFINGERPRINT},
		*fingerprintAlgorithm,
		plan.Enabled.ValueBool())
	addOptionalFingerprintCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddFingerprintCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readFingerprintCertificateMapperResponse(ctx, addResponse.FingerprintCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party certificate-mapper
func (r *certificateMapperResource) CreateThirdPartyCertificateMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan certificateMapperResourceModel) (*certificateMapperResourceModel, error) {
	addRequest := client.NewAddThirdPartyCertificateMapperRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyCertificateMapperSchemaUrn{client.ENUMTHIRDPARTYCERTIFICATEMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CERTIFICATE_MAPPERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyCertificateMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CertificateMapperApi.AddCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCertificateMapperRequest(
		client.AddThirdPartyCertificateMapperRequestAsAddCertificateMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CertificateMapperApi.AddCertificateMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateMapperResourceModel
	readThirdPartyCertificateMapperResponse(ctx, addResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *certificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *certificateMapperResourceModel
	var err error
	if plan.Type.ValueString() == "subject-equals-dn" {
		state, err = r.CreateSubjectEqualsDnCertificateMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "subject-dn-to-user-attribute" {
		state, err = r.CreateSubjectDnToUserAttributeCertificateMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedCertificateMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "subject-attribute-to-user-attribute" {
		state, err = r.CreateSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "fingerprint" {
		state, err = r.CreateFingerprintCertificateMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyCertificateMapper(ctx, req, resp, plan)
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
func (r *defaultCertificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state certificateMapperResourceModel
	if readResponse.SubjectEqualsDnCertificateMapperResponse != nil {
		readSubjectEqualsDnCertificateMapperResponse(ctx, readResponse.SubjectEqualsDnCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectDnToUserAttributeCertificateMapperResponse != nil {
		readSubjectDnToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedCertificateMapperResponse != nil {
		readGroovyScriptedCertificateMapperResponse(ctx, readResponse.GroovyScriptedCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
		readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FingerprintCertificateMapperResponse != nil {
		readFingerprintCertificateMapperResponse(ctx, readResponse.FingerprintCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyCertificateMapperResponse != nil {
		readThirdPartyCertificateMapperResponse(ctx, readResponse.ThirdPartyCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CertificateMapperApi.UpdateCertificateMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SubjectEqualsDnCertificateMapperResponse != nil {
			readSubjectEqualsDnCertificateMapperResponse(ctx, updateResponse.SubjectEqualsDnCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubjectDnToUserAttributeCertificateMapperResponse != nil {
			readSubjectDnToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedCertificateMapperResponse != nil {
			readGroovyScriptedCertificateMapperResponse(ctx, updateResponse.GroovyScriptedCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
			readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FingerprintCertificateMapperResponse != nil {
			readFingerprintCertificateMapperResponse(ctx, updateResponse.FingerprintCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyCertificateMapperResponse != nil {
			readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *certificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readCertificateMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state certificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Mapper", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Mapper", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SubjectEqualsDnCertificateMapperResponse != nil {
		readSubjectEqualsDnCertificateMapperResponse(ctx, readResponse.SubjectEqualsDnCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectDnToUserAttributeCertificateMapperResponse != nil {
		readSubjectDnToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedCertificateMapperResponse != nil {
		readGroovyScriptedCertificateMapperResponse(ctx, readResponse.GroovyScriptedCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
		readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FingerprintCertificateMapperResponse != nil {
		readFingerprintCertificateMapperResponse(ctx, readResponse.FingerprintCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyCertificateMapperResponse != nil {
		readThirdPartyCertificateMapperResponse(ctx, readResponse.ThirdPartyCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *certificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCertificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCertificateMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan certificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state certificateMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CertificateMapperApi.UpdateCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SubjectEqualsDnCertificateMapperResponse != nil {
			readSubjectEqualsDnCertificateMapperResponse(ctx, updateResponse.SubjectEqualsDnCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubjectDnToUserAttributeCertificateMapperResponse != nil {
			readSubjectDnToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedCertificateMapperResponse != nil {
			readGroovyScriptedCertificateMapperResponse(ctx, updateResponse.GroovyScriptedCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
			readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FingerprintCertificateMapperResponse != nil {
			readFingerprintCertificateMapperResponse(ctx, updateResponse.FingerprintCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyCertificateMapperResponse != nil {
			readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCertificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *certificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state certificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CertificateMapperApi.DeleteCertificateMapperExecute(r.apiClient.CertificateMapperApi.DeleteCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Certificate Mapper", err, httpResp)
		return
	}
}

func (r *certificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCertificateMapper(ctx, req, resp)
}

func (r *defaultCertificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCertificateMapper(ctx, req, resp)
}

func importCertificateMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
