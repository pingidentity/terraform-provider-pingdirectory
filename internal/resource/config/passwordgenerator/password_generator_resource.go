package passwordgenerator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &passwordGeneratorResource{}
	_ resource.ResourceWithConfigure   = &passwordGeneratorResource{}
	_ resource.ResourceWithImportState = &passwordGeneratorResource{}
	_ resource.Resource                = &defaultPasswordGeneratorResource{}
	_ resource.ResourceWithConfigure   = &defaultPasswordGeneratorResource{}
	_ resource.ResourceWithImportState = &defaultPasswordGeneratorResource{}
)

// Create a Password Generator resource
func NewPasswordGeneratorResource() resource.Resource {
	return &passwordGeneratorResource{}
}

func NewDefaultPasswordGeneratorResource() resource.Resource {
	return &defaultPasswordGeneratorResource{}
}

// passwordGeneratorResource is the resource implementation.
type passwordGeneratorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPasswordGeneratorResource is the resource implementation.
type defaultPasswordGeneratorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passwordGeneratorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_generator"
}

func (r *defaultPasswordGeneratorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_password_generator"
}

// Configure adds the provider configured client to the resource.
func (r *passwordGeneratorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultPasswordGeneratorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type passwordGeneratorResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	DictionaryFile            types.String `tfsdk:"dictionary_file"`
	MinimumPasswordCharacters types.Int64  `tfsdk:"minimum_password_characters"`
	MinimumPasswordWords      types.Int64  `tfsdk:"minimum_password_words"`
	CapitalizeWords           types.Bool   `tfsdk:"capitalize_words"`
	ScriptClass               types.String `tfsdk:"script_class"`
	ScriptArgument            types.Set    `tfsdk:"script_argument"`
	PasswordCharacterSet      types.Set    `tfsdk:"password_character_set"`
	PasswordFormat            types.String `tfsdk:"password_format"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *passwordGeneratorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordGeneratorSchema(ctx, req, resp, false)
}

func (r *defaultPasswordGeneratorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordGeneratorSchema(ctx, req, resp, true)
}

func passwordGeneratorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Password Generator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Generator resource. Options are ['random', 'groovy-scripted', 'passphrase', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"random", "groovy-scripted", "passphrase", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Password Generator.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Password Generator. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"dictionary_file": schema.StringAttribute{
				Description: "The path to the dictionary file that will be used to obtain the words for use in generated passwords.",
				Optional:    true,
			},
			"minimum_password_characters": schema.Int64Attribute{
				Description: "The minimum number of characters that generated passwords will be required to have.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"minimum_password_words": schema.Int64Attribute{
				Description: "The minimum number of words that must be concatenated in the course of generating a password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"capitalize_words": schema.BoolAttribute{
				Description: "Indicates whether to capitalize each word used in the generated password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Password Generator.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Password Generator. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"password_character_set": schema.SetAttribute{
				Description: "Specifies one or more named character sets.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"password_format": schema.StringAttribute{
				Description: "Specifies the format to use for the generated password.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Generator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Password Generator is enabled for use.",
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
func (r *passwordGeneratorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordGeneratorResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanPasswordGenerator(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model passwordGeneratorResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.DictionaryFile) && model.Type.ValueString() != "passphrase" {
		resp.Diagnostics.AddError("Attribute 'dictionary_file' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'dictionary_file', the 'type' attribute must be one of ['passphrase']")
	}
	if internaltypes.IsDefined(model.PasswordFormat) && model.Type.ValueString() != "random" {
		resp.Diagnostics.AddError("Attribute 'password_format' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_format', the 'type' attribute must be one of ['random']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_argument', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.MinimumPasswordCharacters) && model.Type.ValueString() != "passphrase" {
		resp.Diagnostics.AddError("Attribute 'minimum_password_characters' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'minimum_password_characters', the 'type' attribute must be one of ['passphrase']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.MinimumPasswordWords) && model.Type.ValueString() != "passphrase" {
		resp.Diagnostics.AddError("Attribute 'minimum_password_words' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'minimum_password_words', the 'type' attribute must be one of ['passphrase']")
	}
	if internaltypes.IsDefined(model.PasswordCharacterSet) && model.Type.ValueString() != "random" {
		resp.Diagnostics.AddError("Attribute 'password_character_set' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_character_set', the 'type' attribute must be one of ['random']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_class', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.CapitalizeWords) && model.Type.ValueString() != "passphrase" {
		resp.Diagnostics.AddError("Attribute 'capitalize_words' not supported by pingdirectory_password_generator resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'capitalize_words', the 'type' attribute must be one of ['passphrase']")
	}
}

// Add optional fields to create request for random password-generator
func addOptionalRandomPasswordGeneratorFields(ctx context.Context, addRequest *client.AddRandomPasswordGeneratorRequest, plan passwordGeneratorResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted password-generator
func addOptionalGroovyScriptedPasswordGeneratorFields(ctx context.Context, addRequest *client.AddGroovyScriptedPasswordGeneratorRequest, plan passwordGeneratorResourceModel) {
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

// Add optional fields to create request for passphrase password-generator
func addOptionalPassphrasePasswordGeneratorFields(ctx context.Context, addRequest *client.AddPassphrasePasswordGeneratorRequest, plan passwordGeneratorResourceModel) {
	if internaltypes.IsDefined(plan.MinimumPasswordCharacters) {
		addRequest.MinimumPasswordCharacters = plan.MinimumPasswordCharacters.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MinimumPasswordWords) {
		addRequest.MinimumPasswordWords = plan.MinimumPasswordWords.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.CapitalizeWords) {
		addRequest.CapitalizeWords = plan.CapitalizeWords.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party password-generator
func addOptionalThirdPartyPasswordGeneratorFields(ctx context.Context, addRequest *client.AddThirdPartyPasswordGeneratorRequest, plan passwordGeneratorResourceModel) {
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
func populatePasswordGeneratorUnknownValues(ctx context.Context, model *passwordGeneratorResourceModel) {
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.PasswordCharacterSet.ElementType(ctx) == nil {
		model.PasswordCharacterSet = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
}

// Read a RandomPasswordGeneratorResponse object into the model struct
func readRandomPasswordGeneratorResponse(ctx context.Context, r *client.RandomPasswordGeneratorResponse, state *passwordGeneratorResourceModel, expectedValues *passwordGeneratorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("random")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordCharacterSet = internaltypes.GetStringSet(r.PasswordCharacterSet)
	state.PasswordFormat = types.StringValue(r.PasswordFormat)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordGeneratorUnknownValues(ctx, state)
}

// Read a GroovyScriptedPasswordGeneratorResponse object into the model struct
func readGroovyScriptedPasswordGeneratorResponse(ctx context.Context, r *client.GroovyScriptedPasswordGeneratorResponse, state *passwordGeneratorResourceModel, expectedValues *passwordGeneratorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordGeneratorUnknownValues(ctx, state)
}

// Read a PassphrasePasswordGeneratorResponse object into the model struct
func readPassphrasePasswordGeneratorResponse(ctx context.Context, r *client.PassphrasePasswordGeneratorResponse, state *passwordGeneratorResourceModel, expectedValues *passwordGeneratorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("passphrase")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DictionaryFile = types.StringValue(r.DictionaryFile)
	state.MinimumPasswordCharacters = internaltypes.Int64TypeOrNil(r.MinimumPasswordCharacters)
	state.MinimumPasswordWords = internaltypes.Int64TypeOrNil(r.MinimumPasswordWords)
	state.CapitalizeWords = internaltypes.BoolTypeOrNil(r.CapitalizeWords)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordGeneratorUnknownValues(ctx, state)
}

// Read a ThirdPartyPasswordGeneratorResponse object into the model struct
func readThirdPartyPasswordGeneratorResponse(ctx context.Context, r *client.ThirdPartyPasswordGeneratorResponse, state *passwordGeneratorResourceModel, expectedValues *passwordGeneratorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordGeneratorUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createPasswordGeneratorOperations(plan passwordGeneratorResourceModel, state passwordGeneratorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.DictionaryFile, state.DictionaryFile, "dictionary-file")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumPasswordCharacters, state.MinimumPasswordCharacters, "minimum-password-characters")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinimumPasswordWords, state.MinimumPasswordWords, "minimum-password-words")
	operations.AddBoolOperationIfNecessary(&ops, plan.CapitalizeWords, state.CapitalizeWords, "capitalize-words")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PasswordCharacterSet, state.PasswordCharacterSet, "password-character-set")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordFormat, state.PasswordFormat, "password-format")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a random password-generator
func (r *passwordGeneratorResource) CreateRandomPasswordGenerator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordGeneratorResourceModel) (*passwordGeneratorResourceModel, error) {
	var PasswordCharacterSetSlice []string
	plan.PasswordCharacterSet.ElementsAs(ctx, &PasswordCharacterSetSlice, false)
	addRequest := client.NewAddRandomPasswordGeneratorRequest(plan.Name.ValueString(),
		[]client.EnumrandomPasswordGeneratorSchemaUrn{client.ENUMRANDOMPASSWORDGENERATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_GENERATORRANDOM},
		PasswordCharacterSetSlice,
		plan.PasswordFormat.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalRandomPasswordGeneratorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordGeneratorApi.AddPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordGeneratorRequest(
		client.AddRandomPasswordGeneratorRequestAsAddPasswordGeneratorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.AddPasswordGeneratorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Generator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordGeneratorResourceModel
	readRandomPasswordGeneratorResponse(ctx, addResponse.RandomPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted password-generator
func (r *passwordGeneratorResource) CreateGroovyScriptedPasswordGenerator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordGeneratorResourceModel) (*passwordGeneratorResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedPasswordGeneratorRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedPasswordGeneratorSchemaUrn{client.ENUMGROOVYSCRIPTEDPASSWORDGENERATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_GENERATORGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedPasswordGeneratorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordGeneratorApi.AddPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordGeneratorRequest(
		client.AddGroovyScriptedPasswordGeneratorRequestAsAddPasswordGeneratorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.AddPasswordGeneratorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Generator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordGeneratorResourceModel
	readGroovyScriptedPasswordGeneratorResponse(ctx, addResponse.GroovyScriptedPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a passphrase password-generator
func (r *passwordGeneratorResource) CreatePassphrasePasswordGenerator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordGeneratorResourceModel) (*passwordGeneratorResourceModel, error) {
	addRequest := client.NewAddPassphrasePasswordGeneratorRequest(plan.Name.ValueString(),
		[]client.EnumpassphrasePasswordGeneratorSchemaUrn{client.ENUMPASSPHRASEPASSWORDGENERATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_GENERATORPASSPHRASE},
		plan.DictionaryFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalPassphrasePasswordGeneratorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordGeneratorApi.AddPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordGeneratorRequest(
		client.AddPassphrasePasswordGeneratorRequestAsAddPasswordGeneratorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.AddPasswordGeneratorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Generator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordGeneratorResourceModel
	readPassphrasePasswordGeneratorResponse(ctx, addResponse.PassphrasePasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party password-generator
func (r *passwordGeneratorResource) CreateThirdPartyPasswordGenerator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordGeneratorResourceModel) (*passwordGeneratorResourceModel, error) {
	addRequest := client.NewAddThirdPartyPasswordGeneratorRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyPasswordGeneratorSchemaUrn{client.ENUMTHIRDPARTYPASSWORDGENERATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_GENERATORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyPasswordGeneratorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordGeneratorApi.AddPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordGeneratorRequest(
		client.AddThirdPartyPasswordGeneratorRequestAsAddPasswordGeneratorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.AddPasswordGeneratorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Generator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordGeneratorResourceModel
	readThirdPartyPasswordGeneratorResponse(ctx, addResponse.ThirdPartyPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passwordGeneratorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordGeneratorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passwordGeneratorResourceModel
	var err error
	if plan.Type.ValueString() == "random" {
		state, err = r.CreateRandomPasswordGenerator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedPasswordGenerator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "passphrase" {
		state, err = r.CreatePassphrasePasswordGenerator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPasswordGenerator(ctx, req, resp, plan)
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
func (r *defaultPasswordGeneratorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordGeneratorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.GetPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Generator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passwordGeneratorResourceModel
	if readResponse.RandomPasswordGeneratorResponse != nil {
		readRandomPasswordGeneratorResponse(ctx, readResponse.RandomPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordGeneratorResponse != nil {
		readGroovyScriptedPasswordGeneratorResponse(ctx, readResponse.GroovyScriptedPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PassphrasePasswordGeneratorResponse != nil {
		readPassphrasePasswordGeneratorResponse(ctx, readResponse.PassphrasePasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordGeneratorResponse != nil {
		readThirdPartyPasswordGeneratorResponse(ctx, readResponse.ThirdPartyPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PasswordGeneratorApi.UpdatePasswordGenerator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPasswordGeneratorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PasswordGeneratorApi.UpdatePasswordGeneratorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Generator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.RandomPasswordGeneratorResponse != nil {
			readRandomPasswordGeneratorResponse(ctx, updateResponse.RandomPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedPasswordGeneratorResponse != nil {
			readGroovyScriptedPasswordGeneratorResponse(ctx, updateResponse.GroovyScriptedPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PassphrasePasswordGeneratorResponse != nil {
			readPassphrasePasswordGeneratorResponse(ctx, updateResponse.PassphrasePasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordGeneratorResponse != nil {
			readThirdPartyPasswordGeneratorResponse(ctx, updateResponse.ThirdPartyPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *passwordGeneratorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordGeneratorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPasswordGenerator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state passwordGeneratorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PasswordGeneratorApi.GetPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Generator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.RandomPasswordGeneratorResponse != nil {
		readRandomPasswordGeneratorResponse(ctx, readResponse.RandomPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordGeneratorResponse != nil {
		readGroovyScriptedPasswordGeneratorResponse(ctx, readResponse.GroovyScriptedPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PassphrasePasswordGeneratorResponse != nil {
		readPassphrasePasswordGeneratorResponse(ctx, readResponse.PassphrasePasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordGeneratorResponse != nil {
		readThirdPartyPasswordGeneratorResponse(ctx, readResponse.ThirdPartyPasswordGeneratorResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passwordGeneratorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordGeneratorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordGenerator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePasswordGenerator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passwordGeneratorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passwordGeneratorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PasswordGeneratorApi.UpdatePasswordGenerator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPasswordGeneratorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PasswordGeneratorApi.UpdatePasswordGeneratorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Generator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.RandomPasswordGeneratorResponse != nil {
			readRandomPasswordGeneratorResponse(ctx, updateResponse.RandomPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedPasswordGeneratorResponse != nil {
			readGroovyScriptedPasswordGeneratorResponse(ctx, updateResponse.GroovyScriptedPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PassphrasePasswordGeneratorResponse != nil {
			readPassphrasePasswordGeneratorResponse(ctx, updateResponse.PassphrasePasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordGeneratorResponse != nil {
			readThirdPartyPasswordGeneratorResponse(ctx, updateResponse.ThirdPartyPasswordGeneratorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPasswordGeneratorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passwordGeneratorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passwordGeneratorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PasswordGeneratorApi.DeletePasswordGeneratorExecute(r.apiClient.PasswordGeneratorApi.DeletePasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Password Generator", err, httpResp)
		return
	}
}

func (r *passwordGeneratorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordGenerator(ctx, req, resp)
}

func (r *defaultPasswordGeneratorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordGenerator(ctx, req, resp)
}

func importPasswordGenerator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
