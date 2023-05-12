package extendedoperationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &replaceCertificateExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &replaceCertificateExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &replaceCertificateExtendedOperationHandlerResource{}
	_ resource.Resource                = &defaultReplaceCertificateExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultReplaceCertificateExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultReplaceCertificateExtendedOperationHandlerResource{}
)

// Create a Replace Certificate Extended Operation Handler resource
func NewReplaceCertificateExtendedOperationHandlerResource() resource.Resource {
	return &replaceCertificateExtendedOperationHandlerResource{}
}

func NewDefaultReplaceCertificateExtendedOperationHandlerResource() resource.Resource {
	return &defaultReplaceCertificateExtendedOperationHandlerResource{}
}

// replaceCertificateExtendedOperationHandlerResource is the resource implementation.
type replaceCertificateExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultReplaceCertificateExtendedOperationHandlerResource is the resource implementation.
type defaultReplaceCertificateExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *replaceCertificateExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replace_certificate_extended_operation_handler"
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_replace_certificate_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *replaceCertificateExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type replaceCertificateExtendedOperationHandlerResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	AllowRemotelyProvidedCertificates types.Bool   `tfsdk:"allow_remotely_provided_certificates"`
	AllowedOperation                  types.Set    `tfsdk:"allowed_operation"`
	ConnectionCriteria                types.String `tfsdk:"connection_criteria"`
	RequestCriteria                   types.String `tfsdk:"request_criteria"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *replaceCertificateExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replaceCertificateExtendedOperationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replaceCertificateExtendedOperationHandlerSchema(ctx, req, resp, true)
}

func replaceCertificateExtendedOperationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Replace Certificate Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"allow_remotely_provided_certificates": schema.BoolAttribute{
				Description: "Indicates whether clients should be allowed to directly provide a new listener or inter-server certificate chain in the extended request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_operation": schema.SetAttribute{
				Description: "The types of replace certificate operations that clients will be allowed to request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A set of criteria that client connections must satisfy before they will be allowed to request the associated extended operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"request_criteria": schema.StringAttribute{
				Description: "A set of criteria that the extended requests must satisfy before they will be processed by the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Extended Operation Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).",
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
func addOptionalReplaceCertificateExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddReplaceCertificateExtendedOperationHandlerRequest, plan replaceCertificateExtendedOperationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.AllowRemotelyProvidedCertificates) {
		addRequest.AllowRemotelyProvidedCertificates = plan.AllowRemotelyProvidedCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowedOperation) {
		var slice []string
		plan.AllowedOperation.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumextendedOperationHandlerAllowedOperationProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumextendedOperationHandlerAllowedOperationPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedOperation = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Read a ReplaceCertificateExtendedOperationHandlerResponse object into the model struct
func readReplaceCertificateExtendedOperationHandlerResponse(ctx context.Context, r *client.ReplaceCertificateExtendedOperationHandlerResponse, state *replaceCertificateExtendedOperationHandlerResourceModel, expectedValues *replaceCertificateExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AllowRemotelyProvidedCertificates = internaltypes.BoolTypeOrNil(r.AllowRemotelyProvidedCertificates)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumextendedOperationHandlerAllowedOperationProp(r.AllowedOperation))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createReplaceCertificateExtendedOperationHandlerOperations(plan replaceCertificateExtendedOperationHandlerResourceModel, state replaceCertificateExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowRemotelyProvidedCertificates, state.AllowRemotelyProvidedCertificates, "allow-remotely-provided-certificates")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedOperation, state.AllowedOperation, "allowed-operation")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *replaceCertificateExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replaceCertificateExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddReplaceCertificateExtendedOperationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumreplaceCertificateExtendedOperationHandlerSchemaUrn{client.ENUMREPLACECERTIFICATEEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERREPLACE_CERTIFICATE},
		plan.Enabled.ValueBool())
	err := addOptionalReplaceCertificateExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Replace Certificate Extended Operation Handler", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddReplaceCertificateExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Replace Certificate Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state replaceCertificateExtendedOperationHandlerResourceModel
	readReplaceCertificateExtendedOperationHandlerResponse(ctx, addResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan replaceCertificateExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replace Certificate Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state replaceCertificateExtendedOperationHandlerResourceModel
	readReplaceCertificateExtendedOperationHandlerResponse(ctx, readResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createReplaceCertificateExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replace Certificate Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplaceCertificateExtendedOperationHandlerResponse(ctx, updateResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *replaceCertificateExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplaceCertificateExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readReplaceCertificateExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readReplaceCertificateExtendedOperationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state replaceCertificateExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Replace Certificate Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readReplaceCertificateExtendedOperationHandlerResponse(ctx, readResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *replaceCertificateExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplaceCertificateExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateReplaceCertificateExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateReplaceCertificateExtendedOperationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan replaceCertificateExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state replaceCertificateExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createReplaceCertificateExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Replace Certificate Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readReplaceCertificateExtendedOperationHandlerResponse(ctx, updateResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultReplaceCertificateExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *replaceCertificateExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state replaceCertificateExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandlerExecute(r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Replace Certificate Extended Operation Handler", err, httpResp)
		return
	}
}

func (r *replaceCertificateExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplaceCertificateExtendedOperationHandler(ctx, req, resp)
}

func (r *defaultReplaceCertificateExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importReplaceCertificateExtendedOperationHandler(ctx, req, resp)
}

func importReplaceCertificateExtendedOperationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
