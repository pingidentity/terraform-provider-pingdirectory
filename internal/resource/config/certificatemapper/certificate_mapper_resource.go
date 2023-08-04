package certificatemapper

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type certificateMapperResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	LastUpdated             types.String `tfsdk:"last_updated"`
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
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
		typeAttr.PlanModifiers = []planmodifier.String{}
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *certificateMapperResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCertificateMapperResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanCertificateMapper(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model certificateMapperResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.SubjectAttributeMapping) && model.Type.ValueString() != "subject-attribute-to-user-attribute" {
		resp.Diagnostics.AddError("Attribute 'subject_attribute_mapping' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'subject_attribute_mapping', the 'type' attribute must be one of ['subject-attribute-to-user-attribute']")
	}
	if internaltypes.IsDefined(model.UserBaseDN) && model.Type.ValueString() != "subject-dn-to-user-attribute" && model.Type.ValueString() != "subject-attribute-to-user-attribute" && model.Type.ValueString() != "fingerprint" {
		resp.Diagnostics.AddError("Attribute 'user_base_dn' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'user_base_dn', the 'type' attribute must be one of ['subject-dn-to-user-attribute', 'subject-attribute-to-user-attribute', 'fingerprint']")
	}
	if internaltypes.IsDefined(model.FingerprintAlgorithm) && model.Type.ValueString() != "fingerprint" {
		resp.Diagnostics.AddError("Attribute 'fingerprint_algorithm' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'fingerprint_algorithm', the 'type' attribute must be one of ['fingerprint']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.SubjectAttribute) && model.Type.ValueString() != "subject-dn-to-user-attribute" {
		resp.Diagnostics.AddError("Attribute 'subject_attribute' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'subject_attribute', the 'type' attribute must be one of ['subject-dn-to-user-attribute']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_argument', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.FingerprintAttribute) && model.Type.ValueString() != "fingerprint" {
		resp.Diagnostics.AddError("Attribute 'fingerprint_attribute' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'fingerprint_attribute', the 'type' attribute must be one of ['fingerprint']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_certificate_mapper resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_class', the 'type' attribute must be one of ['groovy-scripted']")
	}
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
func populateCertificateMapperUnknownValues(ctx context.Context, model *certificateMapperResourceModel) {
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.SubjectAttributeMapping.ElementType(ctx) == nil {
		model.SubjectAttributeMapping = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.UserBaseDN.ElementType(ctx) == nil {
		model.UserBaseDN = types.SetNull(types.StringType)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	populateCertificateMapperUnknownValues(ctx, state)
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
	if plan.Type.ValueString() == "subject-equals-dn" {
		readSubjectEqualsDnCertificateMapperResponse(ctx, readResponse.SubjectEqualsDnCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "subject-dn-to-user-attribute" {
		readSubjectDnToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		readGroovyScriptedCertificateMapperResponse(ctx, readResponse.GroovyScriptedCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "subject-attribute-to-user-attribute" {
		readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "fingerprint" {
		readFingerprintCertificateMapperResponse(ctx, readResponse.FingerprintCertificateMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
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
		if plan.Type.ValueString() == "subject-equals-dn" {
			readSubjectEqualsDnCertificateMapperResponse(ctx, updateResponse.SubjectEqualsDnCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "subject-dn-to-user-attribute" {
			readSubjectDnToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedCertificateMapperResponse(ctx, updateResponse.GroovyScriptedCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "subject-attribute-to-user-attribute" {
			readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "fingerprint" {
			readFingerprintCertificateMapperResponse(ctx, updateResponse.FingerprintCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *certificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCertificateMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Mapper", err, httpResp)
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
		if plan.Type.ValueString() == "subject-equals-dn" {
			readSubjectEqualsDnCertificateMapperResponse(ctx, updateResponse.SubjectEqualsDnCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "subject-dn-to-user-attribute" {
			readSubjectDnToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedCertificateMapperResponse(ctx, updateResponse.GroovyScriptedCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "subject-attribute-to-user-attribute" {
			readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "fingerprint" {
			readFingerprintCertificateMapperResponse(ctx, updateResponse.FingerprintCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
	if err != nil {
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
