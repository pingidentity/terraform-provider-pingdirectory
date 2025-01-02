package keypair

import (
	"context"

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
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultKeyPairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keyPairResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	Type                          types.String `tfsdk:"type"`
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
			"type": schema.StringAttribute{
				Description: "The type of Key Pair resource. Options are ['key-pair']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("key-pair"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"key-pair"}...),
				},
			},
			"key_algorithm": schema.StringAttribute{
				Description: "The algorithm name and the length in bits of the key, e.g. RSA_2048.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("RSA_2048"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"RSA_2048", "RSA_3072", "RSA_4096", "EC_256", "EC_384", "EC_521"}...),
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
				Default:     stringdefault.StaticString("cn=Directory Server,O=Ping Identity Key Pair"),
			},
			"certificate_chain": schema.StringAttribute{
				Description: "The PEM-encoded X.509 certificate chain.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"private_key": schema.StringAttribute{
				Description: "The base64-encoded private key that is encrypted using the preferred encryption settings definition.",
				Optional:    true,
				Sensitive:   true,
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		keyAlgorithmAttr := schemaDef.Attributes["key_algorithm"].(schema.StringAttribute)
		keyAlgorithmAttr.PlanModifiers = append(keyAlgorithmAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["key_algorithm"] = keyAlgorithmAttr
		selfSignedCertificateValidityAttr := schemaDef.Attributes["self_signed_certificate_validity"].(schema.StringAttribute)
		selfSignedCertificateValidityAttr.PlanModifiers = append(selfSignedCertificateValidityAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["self_signed_certificate_validity"] = selfSignedCertificateValidityAttr
		subjectDnAttr := schemaDef.Attributes["subject_dn"].(schema.StringAttribute)
		subjectDnAttr.PlanModifiers = append(subjectDnAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["subject_dn"] = subjectDnAttr
		privateKeyAttr := schemaDef.Attributes["private_key"].(schema.StringAttribute)
		privateKeyAttr.PlanModifiers = append(privateKeyAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["private_key"] = privateKeyAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
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

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *keyPairResourceModel) populateAllComputedStringAttributes() {
	if model.PrivateKey.IsUnknown() || model.PrivateKey.IsNull() {
		model.PrivateKey = types.StringValue("")
	}
	if model.SelfSignedCertificateValidity.IsUnknown() || model.SelfSignedCertificateValidity.IsNull() {
		model.SelfSignedCertificateValidity = types.StringValue("")
	}
	if model.CertificateChain.IsUnknown() || model.CertificateChain.IsNull() {
		model.CertificateChain = types.StringValue("")
	}
	if model.SubjectDN.IsUnknown() || model.SubjectDN.IsNull() {
		model.SubjectDN = types.StringValue("")
	}
	if model.KeyAlgorithm.IsUnknown() || model.KeyAlgorithm.IsNull() {
		model.KeyAlgorithm = types.StringValue("")
	}
}

// Read a KeyPairResponse object into the model struct
func readKeyPairResponse(ctx context.Context, r *client.KeyPairResponse, state *keyPairResourceModel, expectedValues *keyPairResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("key-pair")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyAlgorithm = types.StringValue(r.KeyAlgorithm.String())
	state.SelfSignedCertificateValidity = internaltypes.StringTypeOrNil(r.SelfSignedCertificateValidity, true)
	config.CheckMismatchedPDFormattedAttributes("self_signed_certificate_validity",
		expectedValues.SelfSignedCertificateValidity, state.SelfSignedCertificateValidity, diagnostics)
	state.SubjectDN = internaltypes.StringTypeOrNil(r.SubjectDN, true)
	state.CertificateChain = internaltypes.StringTypeOrNil(r.CertificateChain, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *keyPairResourceModel) setStateValuesNotReturnedByAPI(expectedValues *keyPairResourceModel) {
	if !expectedValues.PrivateKey.IsUnknown() {
		state.PrivateKey = expectedValues.PrivateKey
	}
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

// Create a key-pair key-pair
func (r *keyPairResource) CreateKeyPair(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan keyPairResourceModel) (*keyPairResourceModel, error) {
	addRequest := client.NewAddKeyPairRequest(plan.Name.ValueString())
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
	apiAddRequest := r.apiClient.KeyPairAPI.AddKeyPair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyPairRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.KeyPairAPI.AddKeyPairExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Key Pair", err, httpResp)
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
func (r *defaultKeyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyPairResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyPairAPI.GetKeyPair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state keyPairResourceModel
	readKeyPairResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.KeyPairAPI.UpdateKeyPair(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createKeyPairOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyPairAPI.UpdateKeyPairExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readKeyPairResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *keyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readKeyPair(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultKeyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readKeyPair(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readKeyPair(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state keyPairResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.KeyPairAPI.GetKeyPair(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Pair", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readKeyPairResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *keyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateKeyPair(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultKeyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateKeyPair(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateKeyPair(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
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
	updateRequest := apiClient.KeyPairAPI.UpdateKeyPair(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createKeyPairOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.KeyPairAPI.UpdateKeyPairExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Pair", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readKeyPairResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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

	httpResp, err := r.apiClient.KeyPairAPI.DeleteKeyPairExecute(r.apiClient.KeyPairAPI.DeleteKeyPair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Key Pair", err, httpResp)
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
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
