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
	_ resource.Resource                = &subjectAttributeToUserAttributeCertificateMapperResource{}
	_ resource.ResourceWithConfigure   = &subjectAttributeToUserAttributeCertificateMapperResource{}
	_ resource.ResourceWithImportState = &subjectAttributeToUserAttributeCertificateMapperResource{}
	_ resource.Resource                = &defaultSubjectAttributeToUserAttributeCertificateMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultSubjectAttributeToUserAttributeCertificateMapperResource{}
	_ resource.ResourceWithImportState = &defaultSubjectAttributeToUserAttributeCertificateMapperResource{}
)

// Create a Subject Attribute To User Attribute Certificate Mapper resource
func NewSubjectAttributeToUserAttributeCertificateMapperResource() resource.Resource {
	return &subjectAttributeToUserAttributeCertificateMapperResource{}
}

func NewDefaultSubjectAttributeToUserAttributeCertificateMapperResource() resource.Resource {
	return &defaultSubjectAttributeToUserAttributeCertificateMapperResource{}
}

// subjectAttributeToUserAttributeCertificateMapperResource is the resource implementation.
type subjectAttributeToUserAttributeCertificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSubjectAttributeToUserAttributeCertificateMapperResource is the resource implementation.
type defaultSubjectAttributeToUserAttributeCertificateMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subject_attribute_to_user_attribute_certificate_mapper"
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_subject_attribute_to_user_attribute_certificate_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type subjectAttributeToUserAttributeCertificateMapperResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	SubjectAttributeMapping types.Set    `tfsdk:"subject_attribute_mapping"`
	UserBaseDN              types.Set    `tfsdk:"user_base_dn"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	subjectAttributeToUserAttributeCertificateMapperSchema(ctx, req, resp, false)
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	subjectAttributeToUserAttributeCertificateMapperSchema(ctx, req, resp, true)
}

func subjectAttributeToUserAttributeCertificateMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Subject Attribute To User Attribute Certificate Mapper.",
		Attributes: map[string]schema.Attribute{
			"subject_attribute_mapping": schema.SetAttribute{
				Description: "Specifies a mapping between certificate attributes and user attributes.",
				Required:    true,
				ElementType: types.StringType,
			},
			"user_base_dn": schema.SetAttribute{
				Description: "Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry.",
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
func addOptionalSubjectAttributeToUserAttributeCertificateMapperFields(ctx context.Context, addRequest *client.AddSubjectAttributeToUserAttributeCertificateMapperRequest, plan subjectAttributeToUserAttributeCertificateMapperResourceModel) {
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

// Read a SubjectAttributeToUserAttributeCertificateMapperResponse object into the model struct
func readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx context.Context, r *client.SubjectAttributeToUserAttributeCertificateMapperResponse, state *subjectAttributeToUserAttributeCertificateMapperResourceModel, expectedValues *subjectAttributeToUserAttributeCertificateMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SubjectAttributeMapping = internaltypes.GetStringSet(r.SubjectAttributeMapping)
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSubjectAttributeToUserAttributeCertificateMapperOperations(plan subjectAttributeToUserAttributeCertificateMapperResourceModel, state subjectAttributeToUserAttributeCertificateMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SubjectAttributeMapping, state.SubjectAttributeMapping, "subject-attribute-mapping")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UserBaseDN, state.UserBaseDN, "user-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subjectAttributeToUserAttributeCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var SubjectAttributeMappingSlice []string
	plan.SubjectAttributeMapping.ElementsAs(ctx, &SubjectAttributeMappingSlice, false)
	addRequest := client.NewAddSubjectAttributeToUserAttributeCertificateMapperRequest(plan.Id.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state subjectAttributeToUserAttributeCertificateMapperResourceModel
	readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, addResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan subjectAttributeToUserAttributeCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state subjectAttributeToUserAttributeCertificateMapperResourceModel
	readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CertificateMapperApi.UpdateCertificateMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSubjectAttributeToUserAttributeCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSubjectAttributeToUserAttributeCertificateMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state subjectAttributeToUserAttributeCertificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CertificateMapperApi.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *subjectAttributeToUserAttributeCertificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSubjectAttributeToUserAttributeCertificateMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan subjectAttributeToUserAttributeCertificateMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state subjectAttributeToUserAttributeCertificateMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CertificateMapperApi.UpdateCertificateMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSubjectAttributeToUserAttributeCertificateMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CertificateMapperApi.UpdateCertificateMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSubjectAttributeToUserAttributeCertificateMapperResponse(ctx, updateResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *subjectAttributeToUserAttributeCertificateMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state subjectAttributeToUserAttributeCertificateMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CertificateMapperApi.DeleteCertificateMapperExecute(r.apiClient.CertificateMapperApi.DeleteCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Subject Attribute To User Attribute Certificate Mapper", err, httpResp)
		return
	}
}

func (r *subjectAttributeToUserAttributeCertificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp)
}

func (r *defaultSubjectAttributeToUserAttributeCertificateMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSubjectAttributeToUserAttributeCertificateMapper(ctx, req, resp)
}

func importSubjectAttributeToUserAttributeCertificateMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
