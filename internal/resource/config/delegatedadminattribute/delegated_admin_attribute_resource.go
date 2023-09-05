package delegatedadminattribute

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	_ resource.Resource                = &delegatedAdminAttributeResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminAttributeResource{}
	_ resource.ResourceWithImportState = &delegatedAdminAttributeResource{}
	_ resource.Resource                = &defaultDelegatedAdminAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultDelegatedAdminAttributeResource{}
	_ resource.ResourceWithImportState = &defaultDelegatedAdminAttributeResource{}
)

// Create a Delegated Admin Attribute resource
func NewDelegatedAdminAttributeResource() resource.Resource {
	return &delegatedAdminAttributeResource{}
}

func NewDefaultDelegatedAdminAttributeResource() resource.Resource {
	return &defaultDelegatedAdminAttributeResource{}
}

// delegatedAdminAttributeResource is the resource implementation.
type delegatedAdminAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelegatedAdminAttributeResource is the resource implementation.
type defaultDelegatedAdminAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_attribute"
}

func (r *defaultDelegatedAdminAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delegated_admin_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDelegatedAdminAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminAttributeResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	Type                  types.String `tfsdk:"type"`
	RestResourceTypeName  types.String `tfsdk:"rest_resource_type_name"`
	AllowedMIMEType       types.Set    `tfsdk:"allowed_mime_type"`
	Description           types.String `tfsdk:"description"`
	AttributeType         types.String `tfsdk:"attribute_type"`
	DisplayName           types.String `tfsdk:"display_name"`
	Mutability            types.String `tfsdk:"mutability"`
	IncludeInSummary      types.Bool   `tfsdk:"include_in_summary"`
	MultiValued           types.Bool   `tfsdk:"multi_valued"`
	AttributeCategory     types.String `tfsdk:"attribute_category"`
	DisplayOrderIndex     types.Int64  `tfsdk:"display_order_index"`
	ReferenceResourceType types.String `tfsdk:"reference_resource_type"`
	AttributePresentation types.String `tfsdk:"attribute_presentation"`
	DateTimeFormat        types.String `tfsdk:"date_time_format"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminAttributeSchema(ctx, req, resp, false)
}

func (r *defaultDelegatedAdminAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminAttributeSchema(ctx, req, resp, true)
}

func delegatedAdminAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Delegated Admin Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Attribute resource. Options are ['certificate', 'photo', 'generic']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"certificate", "photo", "generic"}...),
				},
			},
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allowed_mime_type": schema.SetAttribute{
				Description: "The list of file types allowed to be uploaded. If no types are specified, then all types will be allowed.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Attribute",
				Optional:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Attribute.",
				Required:    true,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the attribute can be written.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("read-write"),
			},
			"include_in_summary": schema.BoolAttribute{
				Description: "Indicates whether this Delegated Admin Attribute is to be included in the summary display for a resource.",
				Optional:    true,
				Computed:    true,
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Indicates whether this Delegated Admin Attribute may have multiple values.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"attribute_category": schema.StringAttribute{
				Description: "Specifies which attribute category this attribute belongs to.",
				Optional:    true,
			},
			"display_order_index": schema.Int64Attribute{
				Description: "This property determines a display order for attributes within a given attribute category. Attributes are ordered within their category based on this index from least to greatest.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"reference_resource_type": schema.StringAttribute{
				Description: "For LDAP attributes with DN syntax, specifies what kind of resource is referenced.",
				Optional:    true,
			},
			"attribute_presentation": schema.StringAttribute{
				Description: "Indicates how the attribute is presented to the user of the app.",
				Optional:    true,
			},
			"date_time_format": schema.StringAttribute{
				Description: "Specifies the format string that is used to present a date and/or time value to the user of the app. This property only applies to LDAP attribute types whose LDAP syntax is GeneralizedTime and is ignored if the attribute type has any other syntax.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("yyyy-MM-dd"),
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "attribute_type", "rest_resource_type_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		attributeTypeAttr := schemaDef.Attributes["attribute_type"].(schema.StringAttribute)
		attributeTypeAttr.PlanModifiers = append(attributeTypeAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["attribute_type"] = attributeTypeAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *delegatedAdminAttributeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel delegatedAdminAttributeResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for generic type
	if resourceType == "generic" {
		if !internaltypes.IsDefined(configModel.IncludeInSummary) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInSummary.Equal(defaultVal) {
				planModel.IncludeInSummary = defaultVal
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

func (model *delegatedAdminAttributeResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "certificate" {
		model.IncludeInSummary = types.BoolNull()
	}
	if resourceType == "photo" {
		model.IncludeInSummary = types.BoolNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsDelegatedAdminAttribute() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_mime_type"),
			path.MatchRoot("type"),
			[]string{"certificate", "photo"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_in_summary"),
			path.MatchRoot("type"),
			[]string{"generic"},
		),
	}
}

// Add config validators
func (r delegatedAdminAttributeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsDelegatedAdminAttribute()
}

// Add config validators
func (r defaultDelegatedAdminAttributeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsDelegatedAdminAttribute()
}

// Add optional fields to create request for certificate delegated-admin-attribute
func addOptionalCertificateDelegatedAdminAttributeFields(ctx context.Context, addRequest *client.AddCertificateDelegatedAdminAttributeRequest, plan delegatedAdminAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedMIMEType) {
		var slice []string
		plan.AllowedMIMEType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdelegatedAdminAttributeCertificateAllowedMIMETypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdelegatedAdminAttributeCertificateAllowedMIMETypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedMIMEType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Mutability) {
		mutability, err := client.NewEnumdelegatedAdminAttributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	if internaltypes.IsDefined(plan.MultiValued) {
		addRequest.MultiValued = plan.MultiValued.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeCategory) {
		addRequest.AttributeCategory = plan.AttributeCategory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DisplayOrderIndex) {
		addRequest.DisplayOrderIndex = plan.DisplayOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReferenceResourceType) {
		addRequest.ReferenceResourceType = plan.ReferenceResourceType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributePresentation) {
		attributePresentation, err := client.NewEnumdelegatedAdminAttributeAttributePresentationPropFromValue(plan.AttributePresentation.ValueString())
		if err != nil {
			return err
		}
		addRequest.AttributePresentation = attributePresentation
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DateTimeFormat) {
		addRequest.DateTimeFormat = plan.DateTimeFormat.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for photo delegated-admin-attribute
func addOptionalPhotoDelegatedAdminAttributeFields(ctx context.Context, addRequest *client.AddPhotoDelegatedAdminAttributeRequest, plan delegatedAdminAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedMIMEType) {
		var slice []string
		plan.AllowedMIMEType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdelegatedAdminAttributePhotoAllowedMIMETypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdelegatedAdminAttributePhotoAllowedMIMETypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedMIMEType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Mutability) {
		mutability, err := client.NewEnumdelegatedAdminAttributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	if internaltypes.IsDefined(plan.MultiValued) {
		addRequest.MultiValued = plan.MultiValued.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeCategory) {
		addRequest.AttributeCategory = plan.AttributeCategory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DisplayOrderIndex) {
		addRequest.DisplayOrderIndex = plan.DisplayOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReferenceResourceType) {
		addRequest.ReferenceResourceType = plan.ReferenceResourceType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributePresentation) {
		attributePresentation, err := client.NewEnumdelegatedAdminAttributeAttributePresentationPropFromValue(plan.AttributePresentation.ValueString())
		if err != nil {
			return err
		}
		addRequest.AttributePresentation = attributePresentation
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DateTimeFormat) {
		addRequest.DateTimeFormat = plan.DateTimeFormat.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for generic delegated-admin-attribute
func addOptionalGenericDelegatedAdminAttributeFields(ctx context.Context, addRequest *client.AddGenericDelegatedAdminAttributeRequest, plan delegatedAdminAttributeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Mutability) {
		mutability, err := client.NewEnumdelegatedAdminAttributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	if internaltypes.IsDefined(plan.MultiValued) {
		addRequest.MultiValued = plan.MultiValued.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInSummary) {
		addRequest.IncludeInSummary = plan.IncludeInSummary.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeCategory) {
		addRequest.AttributeCategory = plan.AttributeCategory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DisplayOrderIndex) {
		addRequest.DisplayOrderIndex = plan.DisplayOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReferenceResourceType) {
		addRequest.ReferenceResourceType = plan.ReferenceResourceType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributePresentation) {
		attributePresentation, err := client.NewEnumdelegatedAdminAttributeAttributePresentationPropFromValue(plan.AttributePresentation.ValueString())
		if err != nil {
			return err
		}
		addRequest.AttributePresentation = attributePresentation
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DateTimeFormat) {
		addRequest.DateTimeFormat = plan.DateTimeFormat.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateDelegatedAdminAttributeUnknownValues(model *delegatedAdminAttributeResourceModel) {
	if model.AllowedMIMEType.IsUnknown() || model.AllowedMIMEType.IsNull() {
		model.AllowedMIMEType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *delegatedAdminAttributeResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.AttributeType.IsUnknown() || model.AttributeType.IsNull() {
		model.AttributeType = types.StringValue("")
	}
	if model.DisplayName.IsUnknown() || model.DisplayName.IsNull() {
		model.DisplayName = types.StringValue("")
	}
	if model.AttributeCategory.IsUnknown() || model.AttributeCategory.IsNull() {
		model.AttributeCategory = types.StringValue("")
	}
	if model.ReferenceResourceType.IsUnknown() || model.ReferenceResourceType.IsNull() {
		model.ReferenceResourceType = types.StringValue("")
	}
	if model.Mutability.IsUnknown() || model.Mutability.IsNull() {
		model.Mutability = types.StringValue("")
	}
	if model.AttributePresentation.IsUnknown() || model.AttributePresentation.IsNull() {
		model.AttributePresentation = types.StringValue("")
	}
	if model.DateTimeFormat.IsUnknown() || model.DateTimeFormat.IsNull() {
		model.DateTimeFormat = types.StringValue("")
	}
}

// Read a CertificateDelegatedAdminAttributeResponse object into the model struct
func readCertificateDelegatedAdminAttributeResponse(ctx context.Context, r *client.CertificateDelegatedAdminAttributeResponse, state *delegatedAdminAttributeResourceModel, expectedValues *delegatedAdminAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("certificate")
	state.Id = types.StringValue(r.Id)
	state.AllowedMIMEType = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminAttributeCertificateAllowedMIMETypeProp(r.AllowedMIMEType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, internaltypes.IsEmptyString(expectedValues.AttributeCategory))
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, internaltypes.IsEmptyString(expectedValues.ReferenceResourceType))
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), internaltypes.IsEmptyString(expectedValues.AttributePresentation))
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDelegatedAdminAttributeUnknownValues(state)
}

// Read a PhotoDelegatedAdminAttributeResponse object into the model struct
func readPhotoDelegatedAdminAttributeResponse(ctx context.Context, r *client.PhotoDelegatedAdminAttributeResponse, state *delegatedAdminAttributeResourceModel, expectedValues *delegatedAdminAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("photo")
	state.Id = types.StringValue(r.Id)
	state.AllowedMIMEType = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminAttributePhotoAllowedMIMETypeProp(r.AllowedMIMEType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, internaltypes.IsEmptyString(expectedValues.AttributeCategory))
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, internaltypes.IsEmptyString(expectedValues.ReferenceResourceType))
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), internaltypes.IsEmptyString(expectedValues.AttributePresentation))
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDelegatedAdminAttributeUnknownValues(state)
}

// Read a GenericDelegatedAdminAttributeResponse object into the model struct
func readGenericDelegatedAdminAttributeResponse(ctx context.Context, r *client.GenericDelegatedAdminAttributeResponse, state *delegatedAdminAttributeResourceModel, expectedValues *delegatedAdminAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.IncludeInSummary = types.BoolValue(r.IncludeInSummary)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, internaltypes.IsEmptyString(expectedValues.AttributeCategory))
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, internaltypes.IsEmptyString(expectedValues.ReferenceResourceType))
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), internaltypes.IsEmptyString(expectedValues.AttributePresentation))
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateDelegatedAdminAttributeUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *delegatedAdminAttributeResourceModel) setStateValuesNotReturnedByAPI(expectedValues *delegatedAdminAttributeResourceModel) {
	if !expectedValues.RestResourceTypeName.IsUnknown() {
		state.RestResourceTypeName = expectedValues.RestResourceTypeName
	}
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminAttributeOperations(plan delegatedAdminAttributeResourceModel, state delegatedAdminAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedMIMEType, state.AllowedMIMEType, "allowed-mime-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Mutability, state.Mutability, "mutability")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInSummary, state.IncludeInSummary, "include-in-summary")
	operations.AddBoolOperationIfNecessary(&ops, plan.MultiValued, state.MultiValued, "multi-valued")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeCategory, state.AttributeCategory, "attribute-category")
	operations.AddInt64OperationIfNecessary(&ops, plan.DisplayOrderIndex, state.DisplayOrderIndex, "display-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.ReferenceResourceType, state.ReferenceResourceType, "reference-resource-type")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributePresentation, state.AttributePresentation, "attribute-presentation")
	operations.AddStringOperationIfNecessary(&ops, plan.DateTimeFormat, state.DateTimeFormat, "date-time-format")
	return ops
}

// Create a certificate delegated-admin-attribute
func (r *delegatedAdminAttributeResource) CreateCertificateDelegatedAdminAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminAttributeResourceModel) (*delegatedAdminAttributeResourceModel, error) {
	addRequest := client.NewAddCertificateDelegatedAdminAttributeRequest(plan.AttributeType.ValueString(),
		[]client.EnumcertificateDelegatedAdminAttributeSchemaUrn{client.ENUMCERTIFICATEDELEGATEDADMINATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DELEGATED_ADMIN_ATTRIBUTECERTIFICATE},
		plan.DisplayName.ValueString())
	err := addOptionalCertificateDelegatedAdminAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Delegated Admin Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminAttributeRequest(
		client.AddCertificateDelegatedAdminAttributeRequestAsAddDelegatedAdminAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminAttributeResourceModel
	readCertificateDelegatedAdminAttributeResponse(ctx, addResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a photo delegated-admin-attribute
func (r *delegatedAdminAttributeResource) CreatePhotoDelegatedAdminAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminAttributeResourceModel) (*delegatedAdminAttributeResourceModel, error) {
	addRequest := client.NewAddPhotoDelegatedAdminAttributeRequest(plan.AttributeType.ValueString(),
		[]client.EnumphotoDelegatedAdminAttributeSchemaUrn{client.ENUMPHOTODELEGATEDADMINATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DELEGATED_ADMIN_ATTRIBUTEPHOTO},
		plan.DisplayName.ValueString())
	err := addOptionalPhotoDelegatedAdminAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Delegated Admin Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminAttributeRequest(
		client.AddPhotoDelegatedAdminAttributeRequestAsAddDelegatedAdminAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminAttributeResourceModel
	readPhotoDelegatedAdminAttributeResponse(ctx, addResponse.PhotoDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a generic delegated-admin-attribute
func (r *delegatedAdminAttributeResource) CreateGenericDelegatedAdminAttribute(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminAttributeResourceModel) (*delegatedAdminAttributeResourceModel, error) {
	addRequest := client.NewAddGenericDelegatedAdminAttributeRequest(plan.AttributeType.ValueString(),
		[]client.EnumgenericDelegatedAdminAttributeSchemaUrn{client.ENUMGENERICDELEGATEDADMINATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DELEGATED_ADMIN_ATTRIBUTEGENERIC},
		plan.DisplayName.ValueString())
	err := addOptionalGenericDelegatedAdminAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Delegated Admin Attribute", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminAttributeRequest(
		client.AddGenericDelegatedAdminAttributeRequestAsAddDelegatedAdminAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Attribute", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminAttributeResourceModel
	readGenericDelegatedAdminAttributeResponse(ctx, addResponse.GenericDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *delegatedAdminAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *delegatedAdminAttributeResourceModel
	var err error
	if plan.Type.ValueString() == "certificate" {
		state, err = r.CreateCertificateDelegatedAdminAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "photo" {
		state, err = r.CreatePhotoDelegatedAdminAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "generic" {
		state, err = r.CreateGenericDelegatedAdminAttribute(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultDelegatedAdminAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminAttributeResourceModel
	if readResponse.CertificateDelegatedAdminAttributeResponse != nil {
		readCertificateDelegatedAdminAttributeResponse(ctx, readResponse.CertificateDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PhotoDelegatedAdminAttributeResponse != nil {
		readPhotoDelegatedAdminAttributeResponse(ctx, readResponse.PhotoDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GenericDelegatedAdminAttributeResponse != nil {
		readGenericDelegatedAdminAttributeResponse(ctx, readResponse.GenericDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString())
	ops := createDelegatedAdminAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CertificateDelegatedAdminAttributeResponse != nil {
			readCertificateDelegatedAdminAttributeResponse(ctx, updateResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PhotoDelegatedAdminAttributeResponse != nil {
			readPhotoDelegatedAdminAttributeResponse(ctx, updateResponse.PhotoDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GenericDelegatedAdminAttributeResponse != nil {
			readGenericDelegatedAdminAttributeResponse(ctx, updateResponse.GenericDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *delegatedAdminAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDelegatedAdminAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDelegatedAdminAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state delegatedAdminAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.AttributeType.ValueString(), state.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.CertificateDelegatedAdminAttributeResponse != nil {
		readCertificateDelegatedAdminAttributeResponse(ctx, readResponse.CertificateDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PhotoDelegatedAdminAttributeResponse != nil {
		readPhotoDelegatedAdminAttributeResponse(ctx, readResponse.PhotoDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GenericDelegatedAdminAttributeResponse != nil {
		readGenericDelegatedAdminAttributeResponse(ctx, readResponse.GenericDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *delegatedAdminAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelegatedAdminAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CertificateDelegatedAdminAttributeResponse != nil {
			readCertificateDelegatedAdminAttributeResponse(ctx, updateResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PhotoDelegatedAdminAttributeResponse != nil {
			readPhotoDelegatedAdminAttributeResponse(ctx, updateResponse.PhotoDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GenericDelegatedAdminAttributeResponse != nil {
			readGenericDelegatedAdminAttributeResponse(ctx, updateResponse.GenericDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultDelegatedAdminAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delegatedAdminAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminAttributeApi.DeleteDelegatedAdminAttributeExecute(r.apiClient.DelegatedAdminAttributeApi.DeleteDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AttributeType.ValueString(), state.RestResourceTypeName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delegated Admin Attribute", err, httpResp)
		return
	}
}

func (r *delegatedAdminAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminAttribute(ctx, req, resp)
}

func (r *defaultDelegatedAdminAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminAttribute(ctx, req, resp)
}

func importDelegatedAdminAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [rest-resource-type-name]/[delegated-admin-attribute-attribute-type]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_type"), split[1])...)
}
