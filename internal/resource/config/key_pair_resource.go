package config

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &keyPairResource{}
	_ resource.ResourceWithConfigure   = &keyPairResource{}
	_ resource.ResourceWithImportState = &keyPairResource{}
	_ resource.Resource                = &defaultKeyPairResource{}
	_ resource.ResourceWithConfigure   = &defaultKeyPairResource{}
	_ resource.ResourceWithImportState = &defaultKeyPairResource{}
)

// Create a Key Pair resource
func NewKeyPairResource() resource.Resource {
	return &keyPairResource{}
}

func NewDefaultKeyPairResource() resource.Resource {
	return &defaultKeyPairResource{}
}

// keyPairResource is the resource implementation.
type keyPairResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultKeyPairResource is the resource implementation.
type defaultKeyPairResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *keyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_pair"
}

func (r *defaultKeyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_key_pair"
}

// Configure adds the provider configured client to the resource.
func (r *keyPairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultKeyPairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type keyPairResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	KeyAlgorithm                  types.String `tfsdk:"key_algorithm"`
	SelfSignedCertificateValidity types.String `tfsdk:"self_signed_certificate_validity"`
	SubjectDN                     types.String `tfsdk:"subject_dn"`
	CertificateChain              types.String `tfsdk:"certificate_chain"`
	PrivateKey                    types.String `tfsdk:"private_key"`
}

type defaultKeyPairResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	KeyAlgorithm                  types.String `tfsdk:"key_algorithm"`
	SelfSignedCertificateValidity types.String `tfsdk:"self_signed_certificate_validity"`
	SubjectDN                     types.String `tfsdk:"subject_dn"`
	CertificateChain              types.String `tfsdk:"certificate_chain"`
	PrivateKey                    types.String `tfsdk:"private_key"`
}

// GetSchema defines the schema for the resource.
func (r *keyPairResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	keyPairSchema(ctx, req, resp, false)
}

func (r *defaultKeyPairResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	keyPairSchema(ctx, req, resp, true)
}

func keyPairSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Key Pair.",
		Attributes: map[string]schema.Attribute{
			"key_algorithm": schema.StringAttribute{
				Description: "The algorithm name and the length in bits of the key, e.g. RSA_2048.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"self_signed_certificate_validity": schema.StringAttribute{
				Description: "The validity period for a self-signed certificate. If not specified, the self-signed certificate will be valid for approximately 20 years. This is not used when importing an existing key-pair. The system will not automatically rotate expired certificates. It is up to the administrator to do that when that happens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject_dn": schema.StringAttribute{
				Description: "The DN that should be used as the subject for the self-signed certificate and certificate signing request. This is not used when importing an existing key-pair.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"certificate_chain": schema.StringAttribute{
				Description: "The PEM-encoded X.509 certificate chain.",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The base64-encoded private key that is encrypted using the preferred encryption settings definition.",
				Optional:    true,
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"key-pair"}...),
		}
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for key-pair key-pair
func addOptionalKeyPairFields(ctx context.Context, addRequest *client.AddKeyPairRequest, plan keyPairResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyAlgorithm) {
		keyAlgorithm, err := client.NewEnumkeyPairKeyAlgorithmPropFromValue(plan.KeyAlgorithm.ValueString())
		if err != nil {
			return err
		}
		addRequest.KeyAlgorithm = keyAlgorithm
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SelfSignedCertificateValidity) {
		addRequest.SelfSignedCertificateValidity = plan.SelfSignedCertificateValidity.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectDN) {
		addRequest.SubjectDN = plan.SubjectDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CertificateChain) {
		addRequest.CertificateChain = plan.CertificateChain.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKey) {
		addRequest.PrivateKey = plan.PrivateKey.ValueStringPointer()
	}
	return nil
}

// Read a KeyPairResponse object into the model struct
func readKeyPairResponse(ctx context.Context, r *client.KeyPairResponse, state *keyPairResourceModel, expectedValues *keyPairResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.KeyAlgorithm = types.StringValue(r.KeyAlgorithm.String())
	state.SelfSignedCertificateValidity = internaltypes.StringTypeOrNil(r.SelfSignedCertificateValidity, internaltypes.IsEmptyString(expectedValues.SelfSignedCertificateValidity))
	CheckMismatchedPDFormattedAttributes("self_signed_certificate_validity",
		expectedValues.SelfSignedCertificateValidity, state.SelfSignedCertificateValidity, diagnostics)
	state.SubjectDN = internaltypes.StringTypeOrNil(r.SubjectDN, internaltypes.IsEmptyString(expectedValues.SubjectDN))
	state.CertificateChain = internaltypes.StringTypeOrNil(r.CertificateChain, internaltypes.IsEmptyString(expectedValues.CertificateChain))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.PrivateKey = expectedValues.PrivateKey
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a KeyPairResponse object into the model struct
func readKeyPairResponseDefault(ctx context.Context, r *client.KeyPairResponse, state *defaultKeyPairResourceModel, expectedValues *defaultKeyPairResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.KeyAlgorithm = types.StringValue(r.KeyAlgorithm.String())
	state.SelfSignedCertificateValidity = internaltypes.StringTypeOrNil(r.SelfSignedCertificateValidity, internaltypes.IsEmptyString(expectedValues.SelfSignedCertificateValidity))
	CheckMismatchedPDFormattedAttributes("self_signed_certificate_validity",
		expectedValues.SelfSignedCertificateValidity, state.SelfSignedCertificateValidity, diagnostics)
	state.SubjectDN = internaltypes.StringTypeOrNil(r.SubjectDN, internaltypes.IsEmptyString(expectedValues.SubjectDN))
	state.CertificateChain = internaltypes.StringTypeOrNil(r.CertificateChain, internaltypes.IsEmptyString(expectedValues.CertificateChain))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.PrivateKey = expectedValues.PrivateKey
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createKeyPairOperations(plan keyPairResourceModel, state keyPairResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.KeyAlgorithm, state.KeyAlgorithm, "key-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.SelfSignedCertificateValidity, state.SelfSignedCertificateValidity, "self-signed-certificate-validity")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectDN, state.SubjectDN, "subject-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.CertificateChain, state.CertificateChain, "certificate-chain")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKey, state.PrivateKey, "private-key")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createKeyPairOperationsDefault(plan defaultKeyPairResourceModel, state defaultKeyPairResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.KeyAlgorithm, state.KeyAlgorithm, "key-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.SelfSignedCertificateValidity, state.SelfSignedCertificateValidity, "self-signed-certificate-validity")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectDN, state.SubjectDN, "subject-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.CertificateChain, state.CertificateChain, "certificate-chain")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKey, state.PrivateKey, "private-key")
	return ops
}

// Create a key-pair key-pair
func (r *keyPairResource) CreateKeyPair(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan keyPairResourceModel) (*keyPairResourceModel, error) {
	addRequest := client.NewAddKeyPairRequest(plan.Id.ValueString())
	err := addOptionalKeyPairFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Key Pair", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyPairApi.AddKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyPairRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.KeyPairApi.AddKeyPairExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Key Pair", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state keyPairResourceModel
	readKeyPairResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *keyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateKeyPair(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultKeyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultKeyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyPairApi.GetKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultKeyPairResourceModel
	readKeyPairResponseDefault(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.KeyPairApi.UpdateKeyPair(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createKeyPairOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyPairApi.UpdateKeyPairExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readKeyPairResponseDefault(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *keyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state keyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyPairApi.GetKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readKeyPairResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *defaultKeyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultKeyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyPairApi.GetKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readKeyPairResponseDefault(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *keyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan keyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state keyPairResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.KeyPairApi.UpdateKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createKeyPairOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyPairApi.UpdateKeyPairExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readKeyPairResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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

func (r *defaultKeyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultKeyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultKeyPairResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.KeyPairApi.UpdateKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createKeyPairOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyPairApi.UpdateKeyPairExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readKeyPairResponseDefault(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultKeyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *keyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state keyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.KeyPairApi.DeleteKeyPairExecute(r.apiClient.KeyPairApi.DeleteKeyPair(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Key Pair", err, httpResp)
		return
	}
}

func (r *keyPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKeyPair(ctx, req, resp)
}

func (r *defaultKeyPairResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKeyPair(ctx, req, resp)
}

func importKeyPair(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
