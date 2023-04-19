package certificatemapper

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &thirdPartyCertificateMapperResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyCertificateMapperResource{}
	_ resource.ResourceWithImportState = &thirdPartyCertificateMapperResource{}
	_ resource.Resource                = &defaultThirdPartyCertificateMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyCertificateMapperResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyCertificateMapperResource{}
)

// Create a Third Party Certificate Mapper resource
func NewThirdPartyCertificateMapperResource() resource.Resource {
	return &thirdPartyCertificateMapperResource{}
}

func NewDefaultThirdPartyCertificateMapperResource() resource.Resource {
	return &defaultThirdPartyCertificateMapperResource{}
}

// thirdPartyCertificateMapperResource is the resource implementation.
type thirdPartyCertificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyCertificateMapperResource is the resource implementation.
type defaultThirdPartyCertificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyCertificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_certificate_mapper"
}

func (r *defaultThirdPartyCertificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_certificate_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyCertificateMapperResourceModel struct {
	Id                types.String `tfsdk:"id"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyCertificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyCertificateMapperSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyCertificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyCertificateMapperSchema(ctx, req, resp, true)
}

func thirdPartyCertificateMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Certificate Mapper.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Certificate Mapper.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalThirdPartyCertificateMapperFields(ctx context.Context, addRequest *client.AddThirdPartyCertificateMapperRequest, plan thirdPartyCertificateMapperResourceModel) {
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

// Read a ThirdPartyCertificateMapperResponse object into the model struct
func readThirdPartyCertificateMapperResponse(ctx context.Context, r *client.ThirdPartyCertificateMapperResponse, state *thirdPartyCertificateMapperResourceModel, expectedValues *thirdPartyCertificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyCertificateMapperOperations(plan thirdPartyCertificateMapperResourceModel, state thirdPartyCertificateMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *thirdPartyCertificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyCertificateMapperRequest(plan.Id.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyCertificateMapperResourceModel
	readThirdPartyCertificateMapperResponse(ctx, addResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultThirdPartyCertificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyCertificateMapperResourceModel
	readThirdPartyCertificateMapperResponse(ctx, readResponse.ThirdPartyCertificateMapperResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CertificateMapperApi.UpdateCertificateMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyCertificateMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyCertificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyCertificateMapperResponse(ctx, readResponse.ThirdPartyCertificateMapperResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyCertificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyCertificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyCertificateMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyCertificateMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CertificateMapperApi.UpdateCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyCertificateMapperResponse(ctx, updateResponse.ThirdPartyCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyCertificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyCertificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyCertificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CertificateMapperApi.DeleteCertificateMapperExecute(r.apiClient.CertificateMapperApi.DeleteCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Certificate Mapper", err, httpResp)
		return
	}
}

func (r *thirdPartyCertificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyCertificateMapper(ctx, req, resp)
}

func (r *defaultThirdPartyCertificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyCertificateMapper(ctx, req, resp)
}

func importThirdPartyCertificateMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
