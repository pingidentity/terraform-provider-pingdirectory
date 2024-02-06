package jsonfieldconstraints

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/planmodifiers"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jsonFieldConstraintsResource{}
	_ resource.ResourceWithConfigure   = &jsonFieldConstraintsResource{}
	_ resource.ResourceWithImportState = &jsonFieldConstraintsResource{}
	_ resource.Resource                = &defaultJsonFieldConstraintsResource{}
	_ resource.ResourceWithConfigure   = &defaultJsonFieldConstraintsResource{}
	_ resource.ResourceWithImportState = &defaultJsonFieldConstraintsResource{}
)

// Create a Json Field Constraints resource
func NewJsonFieldConstraintsResource() resource.Resource {
	return &jsonFieldConstraintsResource{}
}

func NewDefaultJsonFieldConstraintsResource() resource.Resource {
	return &defaultJsonFieldConstraintsResource{}
}

// jsonFieldConstraintsResource is the resource implementation.
type jsonFieldConstraintsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJsonFieldConstraintsResource is the resource implementation.
type defaultJsonFieldConstraintsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jsonFieldConstraintsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_json_field_constraints"
}

func (r *defaultJsonFieldConstraintsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_json_field_constraints"
}

// Configure adds the provider configured client to the resource.
func (r *jsonFieldConstraintsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultJsonFieldConstraintsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type jsonFieldConstraintsResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	Type                          types.String `tfsdk:"type"`
	JsonAttributeConstraintsName  types.String `tfsdk:"json_attribute_constraints_name"`
	Description                   types.String `tfsdk:"description"`
	JsonField                     types.String `tfsdk:"json_field"`
	ValueType                     types.String `tfsdk:"value_type"`
	IsRequired                    types.Bool   `tfsdk:"is_required"`
	IsArray                       types.String `tfsdk:"is_array"`
	AllowNullValue                types.Bool   `tfsdk:"allow_null_value"`
	AllowEmptyObject              types.Bool   `tfsdk:"allow_empty_object"`
	IndexValues                   types.Bool   `tfsdk:"index_values"`
	IndexEntryLimit               types.Int64  `tfsdk:"index_entry_limit"`
	PrimeIndex                    types.Bool   `tfsdk:"prime_index"`
	CacheMode                     types.String `tfsdk:"cache_mode"`
	TokenizeValues                types.Bool   `tfsdk:"tokenize_values"`
	AllowedValue                  types.Set    `tfsdk:"allowed_value"`
	AllowedValueRegularExpression types.Set    `tfsdk:"allowed_value_regular_expression"`
	MinimumNumericValue           types.String `tfsdk:"minimum_numeric_value"`
	MaximumNumericValue           types.String `tfsdk:"maximum_numeric_value"`
	MinimumValueLength            types.Int64  `tfsdk:"minimum_value_length"`
	MaximumValueLength            types.Int64  `tfsdk:"maximum_value_length"`
	MinimumValueCount             types.Int64  `tfsdk:"minimum_value_count"`
	MaximumValueCount             types.Int64  `tfsdk:"maximum_value_count"`
}

// GetSchema defines the schema for the resource.
func (r *jsonFieldConstraintsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonFieldConstraintsSchema(ctx, req, resp, false)
}

func (r *defaultJsonFieldConstraintsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonFieldConstraintsSchema(ctx, req, resp, true)
}

func jsonFieldConstraintsSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Json Field Constraints.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of JSON Field Constraints resource. Options are ['json-field-constraints']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("json-field-constraints"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"json-field-constraints"}...),
				},
			},
			"json_attribute_constraints_name": schema.StringAttribute{
				Description: "Name of the parent JSON Attribute Constraints",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this JSON Field Constraints",
				Optional:    true,
			},
			"json_field": schema.StringAttribute{
				Description: "The full name of the JSON field to which these constraints apply.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"value_type": schema.StringAttribute{
				Description: "The data type that will be required for values of the target field.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"any", "boolean", "integer", "null", "number", "object", "string"}...),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"is_required": schema.BoolAttribute{
				Description: "Indicates whether the target field must be present in JSON objects stored as values of the associated attribute type.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"is_array": schema.StringAttribute{
				Description: "Indicates whether the value of the target field may be an array of values rather than a single value. If this property is set to \"required\" or \"optional\", then the constraints defined for this field will be applied to each element of the array.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("prohibited"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"required", "optional", "prohibited"}...),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"allow_null_value": schema.BoolAttribute{
				Description: "Indicates whether the target field may have a value that is the JSON null value as an alternative to a value (or array of values) of the specified value-type.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_empty_object": schema.BoolAttribute{
				Description: "Indicates whether the target field may have a value that is an empty JSON object (i.e., a JSON object with zero fields). This may only be set to true if value-type property is set to object.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"index_values": schema.BoolAttribute{
				Description: "Indicates whether backends that support JSON indexing should maintain an index for values of the target field.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that may contain a particular value for the target field before the server will stop maintaining the index for that value.",
				Optional:    true,
			},
			"prime_index": schema.BoolAttribute{
				Description: "Indicates whether backends that support database priming should load the contents of the associated JSON index into memory whenever the backend is opened.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit when caching data for the associated JSON index. This can be useful in environments in which the system does not have enough memory to fully cache the entire data set, as it makes it possible to prioritize which data is the most important to keep in memory.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"cache-keys-and-values", "cache-keys-only", "no-caching"}...),
				},
				PlanModifiers: []planmodifier.String{
					planmodifiers.ToLowercasePlanModifier(),
				},
			},
			"tokenize_values": schema.BoolAttribute{
				Description: "Indicates whether the backend should attempt to assign a compact token for each distinct value for the target field in an attempt to reduce the encoded size of the field in JSON objects. These tokens would be assigned prior to using any from the token set used for automatic compaction of some JSON string values.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allowed_value": schema.SetAttribute{
				Description: "Specifies an explicit set of string values that will be the only values permitted for the target field. If a set of allowed values is defined, then the server will reject any attempt to store a JSON object with a value for the target field that is not included in that set.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allowed_value_regular_expression": schema.SetAttribute{
				Description: "Specifies an explicit set of regular expressions that may be used to restrict the set of values that may be used for the target field. If a set of allowed value regular expressions is defined, then the server will reject any attempt to store a JSON object with a value for the target field that does not match at least one of those regular expressions.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"minimum_numeric_value": schema.StringAttribute{
				Description: "Specifies the smallest numeric value that may be used as the value for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is less than that minimum numeric value.",
				Optional:    true,
			},
			"maximum_numeric_value": schema.StringAttribute{
				Description: "Specifies the largest numeric value that may be used as the value for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is greater than that maximum numeric value.",
				Optional:    true,
			},
			"minimum_value_length": schema.Int64Attribute{
				Description: "Specifies the smallest number of characters that may be present in string values of the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is shorter than that minimum value length.",
				Optional:    true,
			},
			"maximum_value_length": schema.Int64Attribute{
				Description: "Specifies the largest number of characters that may be present in string values of the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is longer than that maximum value length.",
				Optional:    true,
			},
			"minimum_value_count": schema.Int64Attribute{
				Description: "Specifies the smallest number of elements that may be present in an array of values for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is an array with fewer than this number of elements.",
				Optional:    true,
			},
			"maximum_value_count": schema.Int64Attribute{
				Description: "Specifies the largest number of elements that may be present in an array of values for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is an array with more than this number of elements.",
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "json_field", "json_attribute_constraints_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		jsonFieldAttr := schemaDef.Attributes["json_field"].(schema.StringAttribute)
		jsonFieldAttr.PlanModifiers = append(jsonFieldAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["json_field"] = jsonFieldAttr
		valueTypeAttr := schemaDef.Attributes["value_type"].(schema.StringAttribute)
		valueTypeAttr.PlanModifiers = append(valueTypeAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["value_type"] = valueTypeAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for json-field-constraints json-field-constraints
func addOptionalJsonFieldConstraintsFields(ctx context.Context, addRequest *client.AddJsonFieldConstraintsRequest, plan jsonFieldConstraintsResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IsRequired) {
		addRequest.IsRequired = plan.IsRequired.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IsArray) {
		isArray, err := client.NewEnumjsonFieldConstraintsIsArrayPropFromValue(plan.IsArray.ValueString())
		if err != nil {
			return err
		}
		addRequest.IsArray = isArray
	}
	if internaltypes.IsDefined(plan.AllowNullValue) {
		addRequest.AllowNullValue = plan.AllowNullValue.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowEmptyObject) {
		addRequest.AllowEmptyObject = plan.AllowEmptyObject.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IndexValues) {
		addRequest.IndexValues = plan.IndexValues.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IndexEntryLimit) {
		addRequest.IndexEntryLimit = plan.IndexEntryLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.PrimeIndex) {
		addRequest.PrimeIndex = plan.PrimeIndex.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CacheMode) {
		cacheMode, err := client.NewEnumjsonFieldConstraintsCacheModePropFromValue(plan.CacheMode.ValueString())
		if err != nil {
			return err
		}
		addRequest.CacheMode = cacheMode
	}
	if internaltypes.IsDefined(plan.TokenizeValues) {
		addRequest.TokenizeValues = plan.TokenizeValues.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowedValue) {
		var slice []string
		plan.AllowedValue.ElementsAs(ctx, &slice, false)
		addRequest.AllowedValue = slice
	}
	if internaltypes.IsDefined(plan.AllowedValueRegularExpression) {
		var slice []string
		plan.AllowedValueRegularExpression.ElementsAs(ctx, &slice, false)
		addRequest.AllowedValueRegularExpression = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinimumNumericValue) {
		addRequest.MinimumNumericValue = plan.MinimumNumericValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumNumericValue) {
		addRequest.MaximumNumericValue = plan.MaximumNumericValue.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MinimumValueLength) {
		addRequest.MinimumValueLength = plan.MinimumValueLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumValueLength) {
		addRequest.MaximumValueLength = plan.MaximumValueLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MinimumValueCount) {
		addRequest.MinimumValueCount = plan.MinimumValueCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaximumValueCount) {
		addRequest.MaximumValueCount = plan.MaximumValueCount.ValueInt64Pointer()
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *jsonFieldConstraintsResourceModel) populateAllComputedStringAttributes() {
	if model.MinimumNumericValue.IsUnknown() || model.MinimumNumericValue.IsNull() {
		model.MinimumNumericValue = types.StringValue("")
	}
	if model.JsonField.IsUnknown() || model.JsonField.IsNull() {
		model.JsonField = types.StringValue("")
	}
	if model.IsArray.IsUnknown() || model.IsArray.IsNull() {
		model.IsArray = types.StringValue("")
	}
	if model.MaximumNumericValue.IsUnknown() || model.MaximumNumericValue.IsNull() {
		model.MaximumNumericValue = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ValueType.IsUnknown() || model.ValueType.IsNull() {
		model.ValueType = types.StringValue("")
	}
	if model.CacheMode.IsUnknown() || model.CacheMode.IsNull() {
		model.CacheMode = types.StringValue("")
	}
}

// Read a JsonFieldConstraintsResponse object into the model struct
func readJsonFieldConstraintsResponse(ctx context.Context, r *client.JsonFieldConstraintsResponse, state *jsonFieldConstraintsResourceModel, expectedValues *jsonFieldConstraintsResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-field-constraints")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.JsonField = types.StringValue(r.JsonField)
	state.ValueType = types.StringValue(r.ValueType.String())
	state.IsRequired = internaltypes.BoolTypeOrNil(r.IsRequired)
	state.IsArray = internaltypes.StringTypeOrNil(
		client.StringPointerEnumjsonFieldConstraintsIsArrayProp(r.IsArray), true)
	state.AllowNullValue = internaltypes.BoolTypeOrNil(r.AllowNullValue)
	state.AllowEmptyObject = internaltypes.BoolTypeOrNil(r.AllowEmptyObject)
	state.IndexValues = internaltypes.BoolTypeOrNil(r.IndexValues)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumjsonFieldConstraintsCacheModeProp(r.CacheMode), internaltypes.IsEmptyString(expectedValues.CacheMode))
	state.TokenizeValues = internaltypes.BoolTypeOrNil(r.TokenizeValues)
	state.AllowedValue = internaltypes.GetStringSet(r.AllowedValue)
	state.AllowedValueRegularExpression = internaltypes.GetStringSet(r.AllowedValueRegularExpression)
	state.MinimumNumericValue = internaltypes.StringTypeOrNil(r.MinimumNumericValue, internaltypes.IsEmptyString(expectedValues.MinimumNumericValue))
	state.MaximumNumericValue = internaltypes.StringTypeOrNil(r.MaximumNumericValue, internaltypes.IsEmptyString(expectedValues.MaximumNumericValue))
	state.MinimumValueLength = internaltypes.Int64TypeOrNil(r.MinimumValueLength)
	state.MaximumValueLength = internaltypes.Int64TypeOrNil(r.MaximumValueLength)
	state.MinimumValueCount = internaltypes.Int64TypeOrNil(r.MinimumValueCount)
	state.MaximumValueCount = internaltypes.Int64TypeOrNil(r.MaximumValueCount)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *jsonFieldConstraintsResourceModel) setStateValuesNotReturnedByAPI(expectedValues *jsonFieldConstraintsResourceModel) {
	if !expectedValues.JsonAttributeConstraintsName.IsUnknown() {
		state.JsonAttributeConstraintsName = expectedValues.JsonAttributeConstraintsName
	}
}

// Create any update operations necessary to make the state match the plan
func createJsonFieldConstraintsOperations(plan jsonFieldConstraintsResourceModel, state jsonFieldConstraintsResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.JsonField, state.JsonField, "json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.ValueType, state.ValueType, "value-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IsRequired, state.IsRequired, "is-required")
	operations.AddStringOperationIfNecessary(&ops, plan.IsArray, state.IsArray, "is-array")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowNullValue, state.AllowNullValue, "allow-null-value")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowEmptyObject, state.AllowEmptyObject, "allow-empty-object")
	operations.AddBoolOperationIfNecessary(&ops, plan.IndexValues, state.IndexValues, "index-values")
	operations.AddInt64OperationIfNecessary(&ops, plan.IndexEntryLimit, state.IndexEntryLimit, "index-entry-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.PrimeIndex, state.PrimeIndex, "prime-index")
	operations.AddStringOperationIfNecessary(&ops, plan.CacheMode, state.CacheMode, "cache-mode")
	operations.AddBoolOperationIfNecessary(&ops, plan.TokenizeValues, state.TokenizeValues, "tokenize-values")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedValue, state.AllowedValue, "allowed-value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedValueRegularExpression, state.AllowedValueRegularExpression, "allowed-value-regular-expression")
	operations.AddStringOperationIfNecessary(&ops, plan.MinimumNumericValue, state.MinimumNumericValue, "minimum-numeric-value")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumNumericValue, state.MaximumNumericValue, "maximum-numeric-value")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumValueLength, state.MinimumValueLength, "minimum-value-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumValueLength, state.MaximumValueLength, "maximum-value-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumValueCount, state.MinimumValueCount, "minimum-value-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumValueCount, state.MaximumValueCount, "maximum-value-count")
	return ops
}

// Create a json-field-constraints json-field-constraints
func (r *jsonFieldConstraintsResource) CreateJsonFieldConstraints(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan jsonFieldConstraintsResourceModel) (*jsonFieldConstraintsResourceModel, error) {
	valueType, err := client.NewEnumjsonFieldConstraintsValueTypePropFromValue(plan.ValueType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for ValueType", err.Error())
		return nil, err
	}
	addRequest := client.NewAddJsonFieldConstraintsRequest(plan.JsonField.ValueString(),
		*valueType)
	err = addOptionalJsonFieldConstraintsFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Json Field Constraints", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.JsonFieldConstraintsAPI.AddJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.JsonAttributeConstraintsName.ValueString())
	apiAddRequest = apiAddRequest.AddJsonFieldConstraintsRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.JsonFieldConstraintsAPI.AddJsonFieldConstraintsExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Json Field Constraints", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jsonFieldConstraintsResourceModel
	readJsonFieldConstraintsResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *jsonFieldConstraintsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonFieldConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateJsonFieldConstraints(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultJsonFieldConstraintsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonFieldConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.JsonFieldConstraintsAPI.GetJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.JsonField.ValueString(), plan.JsonAttributeConstraintsName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Field Constraints", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jsonFieldConstraintsResourceModel
	readJsonFieldConstraintsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.JsonFieldConstraintsAPI.UpdateJsonFieldConstraints(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.JsonField.ValueString(), plan.JsonAttributeConstraintsName.ValueString())
	ops := createJsonFieldConstraintsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.JsonFieldConstraintsAPI.UpdateJsonFieldConstraintsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Field Constraints", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonFieldConstraintsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jsonFieldConstraintsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonFieldConstraints(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultJsonFieldConstraintsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonFieldConstraints(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readJsonFieldConstraints(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state jsonFieldConstraintsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.JsonFieldConstraintsAPI.GetJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.JsonField.ValueString(), state.JsonAttributeConstraintsName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Json Field Constraints", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Field Constraints", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJsonFieldConstraintsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *jsonFieldConstraintsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonFieldConstraints(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJsonFieldConstraintsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonFieldConstraints(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJsonFieldConstraints(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jsonFieldConstraintsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jsonFieldConstraintsResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.JsonFieldConstraintsAPI.UpdateJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.JsonField.ValueString(), plan.JsonAttributeConstraintsName.ValueString())

	// Determine what update operations are necessary
	ops := createJsonFieldConstraintsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.JsonFieldConstraintsAPI.UpdateJsonFieldConstraintsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Field Constraints", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonFieldConstraintsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJsonFieldConstraintsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jsonFieldConstraintsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jsonFieldConstraintsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.JsonFieldConstraintsAPI.DeleteJsonFieldConstraintsExecute(r.apiClient.JsonFieldConstraintsAPI.DeleteJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.JsonField.ValueString(), state.JsonAttributeConstraintsName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Json Field Constraints", err, httpResp)
		return
	}
}

func (r *jsonFieldConstraintsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonFieldConstraints(ctx, req, resp)
}

func (r *defaultJsonFieldConstraintsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonFieldConstraints(ctx, req, resp)
}

func importJsonFieldConstraints(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [json-attribute-constraints-name]/[json-field-constraints-json-field]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("json_attribute_constraints_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("json_field"), split[1])...)
}
