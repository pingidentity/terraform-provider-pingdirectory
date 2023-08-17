package trustmanagerprovider

import (
	"context"
	"time"

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
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &trustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &trustManagerProviderResource{}
	_ resource.ResourceWithImportState = &trustManagerProviderResource{}
	_ resource.Resource                = &defaultTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &defaultTrustManagerProviderResource{}
)

// Create a Trust Manager Provider resource
func NewTrustManagerProviderResource() resource.Resource {
	return &trustManagerProviderResource{}
}

func NewDefaultTrustManagerProviderResource() resource.Resource {
	return &defaultTrustManagerProviderResource{}
}

// trustManagerProviderResource is the resource implementation.
type trustManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTrustManagerProviderResource is the resource implementation.
type defaultTrustManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *trustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_manager_provider"
}

func (r *defaultTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_trust_manager_provider"
}

// Configure adds the provider configured client to the resource.
func (r *trustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type trustManagerProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	TrustStoreFile                  types.String `tfsdk:"trust_store_file"`
	TrustStoreType                  types.String `tfsdk:"trust_store_type"`
	TrustStorePin                   types.String `tfsdk:"trust_store_pin"`
	TrustStorePinFile               types.String `tfsdk:"trust_store_pin_file"`
	TrustStorePinPassphraseProvider types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
	IncludeJVMDefaultIssuers        types.Bool   `tfsdk:"include_jvm_default_issuers"`
}

// GetSchema defines the schema for the resource.
func (r *trustManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	trustManagerProviderSchema(ctx, req, resp, false)
}

func (r *defaultTrustManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	trustManagerProviderSchema(ctx, req, resp, true)
}

func trustManagerProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Trust Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Trust Manager Provider resource. Options are ['blind', 'file-based', 'jvm-default', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"blind", "file-based", "jvm-default", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Trust Manager Provider.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Trust Manager Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_file": schema.StringAttribute{
				Description: "Specifies the path to the file containing the trust information. It can be an absolute path or a path that is relative to the Directory Server instance root.",
				Optional:    true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the trust store file.",
				Optional:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "Specifies the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Optional:    true,
				Sensitive:   true,
			},
			"trust_store_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Optional:    true,
			},
			"trust_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicate whether the Trust Manager Provider is enabled for use.",
				Required:    true,
			},
			"include_jvm_default_issuers": schema.BoolAttribute{
				Description: "Indicates whether certificates issued by an authority included in the JVM's set of default issuers should be automatically trusted, even if they would not otherwise be trusted by this provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsTrustManagerProvider() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_jvm_default_issuers"),
			path.MatchRoot("type"),
			[]string{"blind", "file-based", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_file"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_type"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin_file"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"file-based"},
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
	}
}

// Add config validators
func (r trustManagerProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsTrustManagerProvider()
}

// Add config validators
func (r defaultTrustManagerProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsTrustManagerProvider()
}

// Add optional fields to create request for blind trust-manager-provider
func addOptionalBlindTrustManagerProviderFields(ctx context.Context, addRequest *client.AddBlindTrustManagerProviderRequest, plan trustManagerProviderResourceModel) {
	if internaltypes.IsDefined(plan.IncludeJVMDefaultIssuers) {
		addRequest.IncludeJVMDefaultIssuers = plan.IncludeJVMDefaultIssuers.ValueBoolPointer()
	}
}

// Add optional fields to create request for file-based trust-manager-provider
func addOptionalFileBasedTrustManagerProviderFields(ctx context.Context, addRequest *client.AddFileBasedTrustManagerProviderRequest, plan trustManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		addRequest.TrustStoreType = plan.TrustStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		addRequest.TrustStorePin = plan.TrustStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePinFile) {
		addRequest.TrustStorePinFile = plan.TrustStorePinFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePinPassphraseProvider) {
		addRequest.TrustStorePinPassphraseProvider = plan.TrustStorePinPassphraseProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeJVMDefaultIssuers) {
		addRequest.IncludeJVMDefaultIssuers = plan.IncludeJVMDefaultIssuers.ValueBoolPointer()
	}
}

// Add optional fields to create request for jvm-default trust-manager-provider
func addOptionalJvmDefaultTrustManagerProviderFields(ctx context.Context, addRequest *client.AddJvmDefaultTrustManagerProviderRequest, plan trustManagerProviderResourceModel) {
}

// Add optional fields to create request for third-party trust-manager-provider
func addOptionalThirdPartyTrustManagerProviderFields(ctx context.Context, addRequest *client.AddThirdPartyTrustManagerProviderRequest, plan trustManagerProviderResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.IncludeJVMDefaultIssuers) {
		addRequest.IncludeJVMDefaultIssuers = plan.IncludeJVMDefaultIssuers.ValueBoolPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateTrustManagerProviderUnknownValues(model *trustManagerProviderResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TrustStorePin.IsUnknown() {
		model.TrustStorePin = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *trustManagerProviderResourceModel) populateAllComputedStringAttributes() {
	if model.TrustStorePinFile.IsUnknown() || model.TrustStorePinFile.IsNull() {
		model.TrustStorePinFile = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.TrustStorePinPassphraseProvider.IsUnknown() || model.TrustStorePinPassphraseProvider.IsNull() {
		model.TrustStorePinPassphraseProvider = types.StringValue("")
	}
	if model.TrustStoreFile.IsUnknown() || model.TrustStoreFile.IsNull() {
		model.TrustStoreFile = types.StringValue("")
	}
	if model.TrustStoreType.IsUnknown() || model.TrustStoreType.IsNull() {
		model.TrustStoreType = types.StringValue("")
	}
}

// Read a BlindTrustManagerProviderResponse object into the model struct
func readBlindTrustManagerProviderResponse(ctx context.Context, r *client.BlindTrustManagerProviderResponse, state *trustManagerProviderResourceModel, expectedValues *trustManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("blind")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = internaltypes.BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTrustManagerProviderUnknownValues(state)
}

// Read a FileBasedTrustManagerProviderResponse object into the model struct
func readFileBasedTrustManagerProviderResponse(ctx context.Context, r *client.FileBasedTrustManagerProviderResponse, state *trustManagerProviderResourceModel, expectedValues *trustManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TrustStoreFile = types.StringValue(r.TrustStoreFile)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	state.TrustStorePinFile = internaltypes.StringTypeOrNil(r.TrustStorePinFile, internaltypes.IsEmptyString(expectedValues.TrustStorePinFile))
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TrustStorePinPassphraseProvider))
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = internaltypes.BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTrustManagerProviderUnknownValues(state)
}

// Read a JvmDefaultTrustManagerProviderResponse object into the model struct
func readJvmDefaultTrustManagerProviderResponse(ctx context.Context, r *client.JvmDefaultTrustManagerProviderResponse, state *trustManagerProviderResourceModel, expectedValues *trustManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jvm-default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTrustManagerProviderUnknownValues(state)
}

// Read a ThirdPartyTrustManagerProviderResponse object into the model struct
func readThirdPartyTrustManagerProviderResponse(ctx context.Context, r *client.ThirdPartyTrustManagerProviderResponse, state *trustManagerProviderResourceModel, expectedValues *trustManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = internaltypes.BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTrustManagerProviderUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *trustManagerProviderResourceModel) setStateValuesNotReturnedByAPI(expectedValues *trustManagerProviderResourceModel) {
	if !expectedValues.TrustStorePin.IsUnknown() {
		state.TrustStorePin = expectedValues.TrustStorePin
	}
}

// Create any update operations necessary to make the state match the plan
func createTrustManagerProviderOperations(plan trustManagerProviderResourceModel, state trustManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinFile, state.TrustStorePinFile, "trust-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinPassphraseProvider, state.TrustStorePinPassphraseProvider, "trust-store-pin-passphrase-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeJVMDefaultIssuers, state.IncludeJVMDefaultIssuers, "include-jvm-default-issuers")
	return ops
}

// Create a blind trust-manager-provider
func (r *trustManagerProviderResource) CreateBlindTrustManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan trustManagerProviderResourceModel) (*trustManagerProviderResourceModel, error) {
	addRequest := client.NewAddBlindTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumblindTrustManagerProviderSchemaUrn{client.ENUMBLINDTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERBLIND},
		plan.Enabled.ValueBool())
	addOptionalBlindTrustManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddBlindTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state trustManagerProviderResourceModel
	readBlindTrustManagerProviderResponse(ctx, addResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based trust-manager-provider
func (r *trustManagerProviderResource) CreateFileBasedTrustManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan trustManagerProviderResourceModel) (*trustManagerProviderResourceModel, error) {
	addRequest := client.NewAddFileBasedTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedTrustManagerProviderSchemaUrn{client.ENUMFILEBASEDTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERFILE_BASED},
		plan.TrustStoreFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedTrustManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddFileBasedTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state trustManagerProviderResourceModel
	readFileBasedTrustManagerProviderResponse(ctx, addResponse.FileBasedTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jvm-default trust-manager-provider
func (r *trustManagerProviderResource) CreateJvmDefaultTrustManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan trustManagerProviderResourceModel) (*trustManagerProviderResourceModel, error) {
	addRequest := client.NewAddJvmDefaultTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumjvmDefaultTrustManagerProviderSchemaUrn{client.ENUMJVMDEFAULTTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERJVM_DEFAULT},
		plan.Enabled.ValueBool())
	addOptionalJvmDefaultTrustManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddJvmDefaultTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state trustManagerProviderResourceModel
	readJvmDefaultTrustManagerProviderResponse(ctx, addResponse.JvmDefaultTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party trust-manager-provider
func (r *trustManagerProviderResource) CreateThirdPartyTrustManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan trustManagerProviderResourceModel) (*trustManagerProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyTrustManagerProviderSchemaUrn{client.ENUMTHIRDPARTYTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyTrustManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddThirdPartyTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state trustManagerProviderResourceModel
	readThirdPartyTrustManagerProviderResponse(ctx, addResponse.ThirdPartyTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *trustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan trustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *trustManagerProviderResourceModel
	var err error
	if plan.Type.ValueString() == "blind" {
		state, err = r.CreateBlindTrustManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based" {
		state, err = r.CreateFileBasedTrustManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jvm-default" {
		state, err = r.CreateJvmDefaultTrustManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyTrustManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan trustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state trustManagerProviderResourceModel
	if readResponse.BlindTrustManagerProviderResponse != nil {
		readBlindTrustManagerProviderResponse(ctx, readResponse.BlindTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedTrustManagerProviderResponse != nil {
		readFileBasedTrustManagerProviderResponse(ctx, readResponse.FileBasedTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JvmDefaultTrustManagerProviderResponse != nil {
		readJvmDefaultTrustManagerProviderResponse(ctx, readResponse.JvmDefaultTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyTrustManagerProviderResponse != nil {
		readThirdPartyTrustManagerProviderResponse(ctx, readResponse.ThirdPartyTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Trust Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.BlindTrustManagerProviderResponse != nil {
			readBlindTrustManagerProviderResponse(ctx, updateResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedTrustManagerProviderResponse != nil {
			readFileBasedTrustManagerProviderResponse(ctx, updateResponse.FileBasedTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JvmDefaultTrustManagerProviderResponse != nil {
			readJvmDefaultTrustManagerProviderResponse(ctx, updateResponse.JvmDefaultTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyTrustManagerProviderResponse != nil {
			readThirdPartyTrustManagerProviderResponse(ctx, updateResponse.ThirdPartyTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *trustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readTrustManagerProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state trustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.BlindTrustManagerProviderResponse != nil {
		readBlindTrustManagerProviderResponse(ctx, readResponse.BlindTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedTrustManagerProviderResponse != nil {
		readFileBasedTrustManagerProviderResponse(ctx, readResponse.FileBasedTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JvmDefaultTrustManagerProviderResponse != nil {
		readJvmDefaultTrustManagerProviderResponse(ctx, readResponse.JvmDefaultTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyTrustManagerProviderResponse != nil {
		readThirdPartyTrustManagerProviderResponse(ctx, readResponse.ThirdPartyTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *trustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTrustManagerProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan trustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state trustManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Trust Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.BlindTrustManagerProviderResponse != nil {
			readBlindTrustManagerProviderResponse(ctx, updateResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedTrustManagerProviderResponse != nil {
			readFileBasedTrustManagerProviderResponse(ctx, updateResponse.FileBasedTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JvmDefaultTrustManagerProviderResponse != nil {
			readJvmDefaultTrustManagerProviderResponse(ctx, updateResponse.JvmDefaultTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyTrustManagerProviderResponse != nil {
			readThirdPartyTrustManagerProviderResponse(ctx, updateResponse.ThirdPartyTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *defaultTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *trustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state trustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Trust Manager Provider", err, httpResp)
		return
	}
}

func (r *trustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTrustManagerProvider(ctx, req, resp)
}

func (r *defaultTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTrustManagerProvider(ctx, req, resp)
}

func importTrustManagerProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
