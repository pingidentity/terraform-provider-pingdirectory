package delegatedadmincorrelatedrestresource

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
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delegatedAdminCorrelatedRestResourceResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminCorrelatedRestResourceResource{}
	_ resource.ResourceWithImportState = &delegatedAdminCorrelatedRestResourceResource{}
	_ resource.Resource                = &defaultDelegatedAdminCorrelatedRestResourceResource{}
	_ resource.ResourceWithConfigure   = &defaultDelegatedAdminCorrelatedRestResourceResource{}
	_ resource.ResourceWithImportState = &defaultDelegatedAdminCorrelatedRestResourceResource{}
)

// Create a Delegated Admin Correlated Rest Resource resource
func NewDelegatedAdminCorrelatedRestResourceResource() resource.Resource {
	return &delegatedAdminCorrelatedRestResourceResource{}
}

func NewDefaultDelegatedAdminCorrelatedRestResourceResource() resource.Resource {
	return &defaultDelegatedAdminCorrelatedRestResourceResource{}
}

// delegatedAdminCorrelatedRestResourceResource is the resource implementation.
type delegatedAdminCorrelatedRestResourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelegatedAdminCorrelatedRestResourceResource is the resource implementation.
type defaultDelegatedAdminCorrelatedRestResourceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminCorrelatedRestResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_correlated_rest_resource"
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delegated_admin_correlated_rest_resource"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminCorrelatedRestResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminCorrelatedRestResourceResourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	Name                                      types.String `tfsdk:"name"`
	Notifications                             types.Set    `tfsdk:"notifications"`
	RequiredActions                           types.Set    `tfsdk:"required_actions"`
	Type                                      types.String `tfsdk:"type"`
	RestResourceTypeName                      types.String `tfsdk:"rest_resource_type_name"`
	DisplayName                               types.String `tfsdk:"display_name"`
	CorrelatedRESTResource                    types.String `tfsdk:"correlated_rest_resource"`
	PrimaryRESTResourceCorrelationAttribute   types.String `tfsdk:"primary_rest_resource_correlation_attribute"`
	SecondaryRESTResourceCorrelationAttribute types.String `tfsdk:"secondary_rest_resource_correlation_attribute"`
	UseSecondaryValueForLinking               types.Bool   `tfsdk:"use_secondary_value_for_linking"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminCorrelatedRestResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminCorrelatedRestResourceSchema(ctx, req, resp, false)
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delegatedAdminCorrelatedRestResourceSchema(ctx, req, resp, true)
}

func delegatedAdminCorrelatedRestResourceSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Delegated Admin Correlated Rest Resource.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Correlated REST Resource resource. Options are ['delegated-admin-correlated-rest-resource']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("delegated-admin-correlated-rest-resource"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"delegated-admin-correlated-rest-resource"}...),
				},
			},
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Correlated REST Resource.",
				Required:    true,
			},
			"correlated_rest_resource": schema.StringAttribute{
				Description: "The REST Resource Type that will be linked to this REST Resource Type.",
				Required:    true,
			},
			"primary_rest_resource_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the parent REST Resource Type whose value will be used to match objects in the Delegated Admin Correlated REST Resource. This attribute must be writeable when use-secondary-value-for-linking is enabled.",
				Required:    true,
			},
			"secondary_rest_resource_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Delegated Admin Correlated REST Resource whose value will be matched with the primary-rest-resource-correlation-attribute. This attribute must be writeable when use-secondary-value-for-linking is disabled.",
				Required:    true,
			},
			"use_secondary_value_for_linking": schema.BoolAttribute{
				Description: "Indicates whether links should be created using the secondary correlation attribute value.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "rest_resource_type_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for delegated-admin-correlated-rest-resource delegated-admin-correlated-rest-resource
func addOptionalDelegatedAdminCorrelatedRestResourceFields(ctx context.Context, addRequest *client.AddDelegatedAdminCorrelatedRestResourceRequest, plan delegatedAdminCorrelatedRestResourceResourceModel) {
	if internaltypes.IsDefined(plan.UseSecondaryValueForLinking) {
		addRequest.UseSecondaryValueForLinking = plan.UseSecondaryValueForLinking.ValueBoolPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *delegatedAdminCorrelatedRestResourceResourceModel) populateAllComputedStringAttributes() {
	if model.CorrelatedRESTResource.IsUnknown() || model.CorrelatedRESTResource.IsNull() {
		model.CorrelatedRESTResource = types.StringValue("")
	}
	if model.DisplayName.IsUnknown() || model.DisplayName.IsNull() {
		model.DisplayName = types.StringValue("")
	}
	if model.PrimaryRESTResourceCorrelationAttribute.IsUnknown() || model.PrimaryRESTResourceCorrelationAttribute.IsNull() {
		model.PrimaryRESTResourceCorrelationAttribute = types.StringValue("")
	}
	if model.SecondaryRESTResourceCorrelationAttribute.IsUnknown() || model.SecondaryRESTResourceCorrelationAttribute.IsNull() {
		model.SecondaryRESTResourceCorrelationAttribute = types.StringValue("")
	}
}

// Read a DelegatedAdminCorrelatedRestResourceResponse object into the model struct
func readDelegatedAdminCorrelatedRestResourceResponse(ctx context.Context, r *client.DelegatedAdminCorrelatedRestResourceResponse, state *delegatedAdminCorrelatedRestResourceResourceModel, expectedValues *delegatedAdminCorrelatedRestResourceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("delegated-admin-correlated-rest-resource")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.CorrelatedRESTResource = types.StringValue(r.CorrelatedRESTResource)
	state.PrimaryRESTResourceCorrelationAttribute = types.StringValue(r.PrimaryRESTResourceCorrelationAttribute)
	state.SecondaryRESTResourceCorrelationAttribute = types.StringValue(r.SecondaryRESTResourceCorrelationAttribute)
	state.UseSecondaryValueForLinking = internaltypes.BoolTypeOrNil(r.UseSecondaryValueForLinking)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *delegatedAdminCorrelatedRestResourceResourceModel) setStateValuesNotReturnedByAPI(expectedValues *delegatedAdminCorrelatedRestResourceResourceModel) {
	if !expectedValues.RestResourceTypeName.IsUnknown() {
		state.RestResourceTypeName = expectedValues.RestResourceTypeName
	}
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminCorrelatedRestResourceOperations(plan delegatedAdminCorrelatedRestResourceResourceModel, state delegatedAdminCorrelatedRestResourceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelatedRESTResource, state.CorrelatedRESTResource, "correlated-rest-resource")
	operations.AddStringOperationIfNecessary(&ops, plan.PrimaryRESTResourceCorrelationAttribute, state.PrimaryRESTResourceCorrelationAttribute, "primary-rest-resource-correlation-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SecondaryRESTResourceCorrelationAttribute, state.SecondaryRESTResourceCorrelationAttribute, "secondary-rest-resource-correlation-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSecondaryValueForLinking, state.UseSecondaryValueForLinking, "use-secondary-value-for-linking")
	return ops
}

// Create a delegated-admin-correlated-rest-resource delegated-admin-correlated-rest-resource
func (r *delegatedAdminCorrelatedRestResourceResource) CreateDelegatedAdminCorrelatedRestResource(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan delegatedAdminCorrelatedRestResourceResourceModel) (*delegatedAdminCorrelatedRestResourceResourceModel, error) {
	addRequest := client.NewAddDelegatedAdminCorrelatedRestResourceRequest(plan.DisplayName.ValueString(),
		plan.CorrelatedRESTResource.ValueString(),
		plan.PrimaryRESTResourceCorrelationAttribute.ValueString(),
		plan.SecondaryRESTResourceCorrelationAttribute.ValueString(),
		plan.Name.ValueString())
	addOptionalDelegatedAdminCorrelatedRestResourceFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.AddDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminCorrelatedRestResourceRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.AddDelegatedAdminCorrelatedRestResourceExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delegated Admin Correlated Rest Resource", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delegatedAdminCorrelatedRestResourceResourceModel
	readDelegatedAdminCorrelatedRestResourceResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *delegatedAdminCorrelatedRestResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminCorrelatedRestResourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateDelegatedAdminCorrelatedRestResource(ctx, req, resp, plan)
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
func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminCorrelatedRestResourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.GetDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Correlated Rest Resource", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminCorrelatedRestResourceResourceModel
	readDelegatedAdminCorrelatedRestResourceResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.UpdateDelegatedAdminCorrelatedRestResource(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.RestResourceTypeName.ValueString())
	ops := createDelegatedAdminCorrelatedRestResourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.UpdateDelegatedAdminCorrelatedRestResourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Correlated Rest Resource", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminCorrelatedRestResourceResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *delegatedAdminCorrelatedRestResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminCorrelatedRestResource(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelegatedAdminCorrelatedRestResource(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readDelegatedAdminCorrelatedRestResource(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state delegatedAdminCorrelatedRestResourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminCorrelatedRestResourceAPI.GetDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Correlated Rest Resource", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Correlated Rest Resource", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminCorrelatedRestResourceResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *delegatedAdminCorrelatedRestResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminCorrelatedRestResource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelegatedAdminCorrelatedRestResource(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelegatedAdminCorrelatedRestResource(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delegatedAdminCorrelatedRestResourceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminCorrelatedRestResourceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminCorrelatedRestResourceAPI.UpdateDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.RestResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminCorrelatedRestResourceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminCorrelatedRestResourceAPI.UpdateDelegatedAdminCorrelatedRestResourceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Correlated Rest Resource", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminCorrelatedRestResourceResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDelegatedAdminCorrelatedRestResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminCorrelatedRestResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delegatedAdminCorrelatedRestResourceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.DeleteDelegatedAdminCorrelatedRestResourceExecute(r.apiClient.DelegatedAdminCorrelatedRestResourceAPI.DeleteDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.RestResourceTypeName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delegated Admin Correlated Rest Resource", err, httpResp)
		return
	}
}

func (r *delegatedAdminCorrelatedRestResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminCorrelatedRestResource(ctx, req, resp)
}

func (r *defaultDelegatedAdminCorrelatedRestResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelegatedAdminCorrelatedRestResource(ctx, req, resp)
}

func importDelegatedAdminCorrelatedRestResource(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [rest-resource-type-name]/[delegated-admin-correlated-rest-resource-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
