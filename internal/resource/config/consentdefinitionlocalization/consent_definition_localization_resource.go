// Copyright © 2025 Ping Identity Corporation

package consentdefinitionlocalization

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consentDefinitionLocalizationResource{}
	_ resource.ResourceWithConfigure   = &consentDefinitionLocalizationResource{}
	_ resource.ResourceWithImportState = &consentDefinitionLocalizationResource{}
	_ resource.Resource                = &defaultConsentDefinitionLocalizationResource{}
	_ resource.ResourceWithConfigure   = &defaultConsentDefinitionLocalizationResource{}
	_ resource.ResourceWithImportState = &defaultConsentDefinitionLocalizationResource{}
)

// Create a Consent Definition Localization resource
func NewConsentDefinitionLocalizationResource() resource.Resource {
	return &consentDefinitionLocalizationResource{}
}

func NewDefaultConsentDefinitionLocalizationResource() resource.Resource {
	return &defaultConsentDefinitionLocalizationResource{}
}

// consentDefinitionLocalizationResource is the resource implementation.
type consentDefinitionLocalizationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConsentDefinitionLocalizationResource is the resource implementation.
type defaultConsentDefinitionLocalizationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consentDefinitionLocalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_definition_localization"
}

func (r *defaultConsentDefinitionLocalizationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_consent_definition_localization"
}

// Configure adds the provider configured client to the resource.
func (r *consentDefinitionLocalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultConsentDefinitionLocalizationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentDefinitionLocalizationResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	Type                  types.String `tfsdk:"type"`
	ConsentDefinitionName types.String `tfsdk:"consent_definition_name"`
	Locale                types.String `tfsdk:"locale"`
	Version               types.String `tfsdk:"version"`
	TitleText             types.String `tfsdk:"title_text"`
	DataText              types.String `tfsdk:"data_text"`
	PurposeText           types.String `tfsdk:"purpose_text"`
}

// GetSchema defines the schema for the resource.
func (r *consentDefinitionLocalizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consentDefinitionLocalizationSchema(ctx, req, resp, false)
}

func (r *defaultConsentDefinitionLocalizationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consentDefinitionLocalizationSchema(ctx, req, resp, true)
}

func consentDefinitionLocalizationSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Consent Definition Localization.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Consent Definition Localization resource. Options are ['consent-definition-localization']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("consent-definition-localization"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"consent-definition-localization"}...),
				},
			},
			"consent_definition_name": schema.StringAttribute{
				Description: "Name of the parent Consent Definition",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"locale": schema.StringAttribute{
				Description: "The locale of this Consent Definition Localization.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"version": schema.StringAttribute{
				Description: "The version of this Consent Definition Localization, using the format MAJOR.MINOR.",
				Required:    true,
			},
			"title_text": schema.StringAttribute{
				Description: "Localized text that may be used to provide a title or summary for a consent request or a granted consent.",
				Optional:    true,
			},
			"data_text": schema.StringAttribute{
				Description: "Localized text describing the data to be shared.",
				Required:    true,
			},
			"purpose_text": schema.StringAttribute{
				Description: "Localized text describing how the data is to be used.",
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "locale", "consent_definition_name"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		localeAttr := schemaDef.Attributes["locale"].(schema.StringAttribute)
		localeAttr.PlanModifiers = append(localeAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["locale"] = localeAttr
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Add optional fields to create request for consent-definition-localization consent-definition-localization
func addOptionalConsentDefinitionLocalizationFields(ctx context.Context, addRequest *client.AddConsentDefinitionLocalizationRequest, plan consentDefinitionLocalizationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TitleText) {
		addRequest.TitleText = plan.TitleText.ValueStringPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *consentDefinitionLocalizationResourceModel) populateAllComputedStringAttributes() {
	if model.Locale.IsUnknown() || model.Locale.IsNull() {
		model.Locale = types.StringValue("")
	}
	if model.TitleText.IsUnknown() || model.TitleText.IsNull() {
		model.TitleText = types.StringValue("")
	}
	if model.Version.IsUnknown() || model.Version.IsNull() {
		model.Version = types.StringValue("")
	}
	if model.PurposeText.IsUnknown() || model.PurposeText.IsNull() {
		model.PurposeText = types.StringValue("")
	}
	if model.DataText.IsUnknown() || model.DataText.IsNull() {
		model.DataText = types.StringValue("")
	}
}

// Read a ConsentDefinitionLocalizationResponse object into the model struct
func readConsentDefinitionLocalizationResponse(ctx context.Context, r *client.ConsentDefinitionLocalizationResponse, state *consentDefinitionLocalizationResourceModel, expectedValues *consentDefinitionLocalizationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent-definition-localization")
	state.Id = types.StringValue(r.Id)
	state.Locale = types.StringValue(r.Locale)
	state.Version = types.StringValue(r.Version)
	state.TitleText = internaltypes.StringTypeOrNil(r.TitleText, internaltypes.IsEmptyString(expectedValues.TitleText))
	state.DataText = types.StringValue(r.DataText)
	state.PurposeText = types.StringValue(r.PurposeText)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *consentDefinitionLocalizationResourceModel) setStateValuesNotReturnedByAPI(expectedValues *consentDefinitionLocalizationResourceModel) {
	if !expectedValues.ConsentDefinitionName.IsUnknown() {
		state.ConsentDefinitionName = expectedValues.ConsentDefinitionName
	}
}

// Create any update operations necessary to make the state match the plan
func createConsentDefinitionLocalizationOperations(plan consentDefinitionLocalizationResourceModel, state consentDefinitionLocalizationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Locale, state.Locale, "locale")
	operations.AddStringOperationIfNecessary(&ops, plan.Version, state.Version, "version")
	operations.AddStringOperationIfNecessary(&ops, plan.TitleText, state.TitleText, "title-text")
	operations.AddStringOperationIfNecessary(&ops, plan.DataText, state.DataText, "data-text")
	operations.AddStringOperationIfNecessary(&ops, plan.PurposeText, state.PurposeText, "purpose-text")
	return ops
}

// Create a consent-definition-localization consent-definition-localization
func (r *consentDefinitionLocalizationResource) CreateConsentDefinitionLocalization(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan consentDefinitionLocalizationResourceModel) (*consentDefinitionLocalizationResourceModel, error) {
	addRequest := client.NewAddConsentDefinitionLocalizationRequest(plan.Locale.ValueString(),
		plan.Version.ValueString(),
		plan.DataText.ValueString(),
		plan.PurposeText.ValueString(),
		plan.Locale.ValueString())
	addOptionalConsentDefinitionLocalizationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConsentDefinitionLocalizationAPI.AddConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.ConsentDefinitionName.ValueString())
	apiAddRequest = apiAddRequest.AddConsentDefinitionLocalizationRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationAPI.AddConsentDefinitionLocalizationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Consent Definition Localization", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state consentDefinitionLocalizationResourceModel
	readConsentDefinitionLocalizationResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *consentDefinitionLocalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionLocalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateConsentDefinitionLocalization(ctx, req, resp, plan)
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
func (r *defaultConsentDefinitionLocalizationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentDefinitionLocalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationAPI.GetConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Locale.ValueString(), plan.ConsentDefinitionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition Localization", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state consentDefinitionLocalizationResourceModel
	readConsentDefinitionLocalizationResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConsentDefinitionLocalizationAPI.UpdateConsentDefinitionLocalization(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Locale.ValueString(), plan.ConsentDefinitionName.ValueString())
	ops := createConsentDefinitionLocalizationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentDefinitionLocalizationAPI.UpdateConsentDefinitionLocalizationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition Localization", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionLocalizationResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *consentDefinitionLocalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsentDefinitionLocalization(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultConsentDefinitionLocalizationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsentDefinitionLocalization(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readConsentDefinitionLocalization(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state consentDefinitionLocalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConsentDefinitionLocalizationAPI.GetConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Locale.ValueString(), state.ConsentDefinitionName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition Localization", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Definition Localization", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentDefinitionLocalizationResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *consentDefinitionLocalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsentDefinitionLocalization(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConsentDefinitionLocalizationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsentDefinitionLocalization(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConsentDefinitionLocalization(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan consentDefinitionLocalizationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consentDefinitionLocalizationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConsentDefinitionLocalizationAPI.UpdateConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Locale.ValueString(), plan.ConsentDefinitionName.ValueString())

	// Determine what update operations are necessary
	ops := createConsentDefinitionLocalizationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConsentDefinitionLocalizationAPI.UpdateConsentDefinitionLocalizationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Definition Localization", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentDefinitionLocalizationResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConsentDefinitionLocalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *consentDefinitionLocalizationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state consentDefinitionLocalizationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConsentDefinitionLocalizationAPI.DeleteConsentDefinitionLocalizationExecute(r.apiClient.ConsentDefinitionLocalizationAPI.DeleteConsentDefinitionLocalization(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Locale.ValueString(), state.ConsentDefinitionName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Consent Definition Localization", err, httpResp)
		return
	}
}

func (r *consentDefinitionLocalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsentDefinitionLocalization(ctx, req, resp)
}

func (r *defaultConsentDefinitionLocalizationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsentDefinitionLocalization(ctx, req, resp)
}

func importConsentDefinitionLocalization(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [consent-definition-name]/[consent-definition-localization-locale]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("consent_definition_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("locale"), split[1])...)
}
