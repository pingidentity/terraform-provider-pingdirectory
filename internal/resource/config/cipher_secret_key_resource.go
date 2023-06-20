package config

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cipherSecretKeyResource{}
	_ resource.ResourceWithConfigure   = &cipherSecretKeyResource{}
	_ resource.ResourceWithImportState = &cipherSecretKeyResource{}
)

// Create a Cipher Secret Key resource
func NewCipherSecretKeyResource() resource.Resource {
	return &cipherSecretKeyResource{}
}

// cipherSecretKeyResource is the resource implementation.
type cipherSecretKeyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *cipherSecretKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_cipher_secret_key"
}

// Configure adds the provider configured client to the resource.
func (r *cipherSecretKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type cipherSecretKeyResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	ServerInstanceName             types.String `tfsdk:"server_instance_name"`
	CipherTransformationName       types.String `tfsdk:"cipher_transformation_name"`
	InitializationVectorLengthBits types.Int64  `tfsdk:"initialization_vector_length_bits"`
	KeyID                          types.String `tfsdk:"key_id"`
	IsCompromised                  types.Bool   `tfsdk:"is_compromised"`
	SymmetricKey                   types.Set    `tfsdk:"symmetric_key"`
	KeyLengthBits                  types.Int64  `tfsdk:"key_length_bits"`
}

// GetSchema defines the schema for the resource.
func (r *cipherSecretKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Cipher Secret Key.",
		Attributes: map[string]schema.Attribute{
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"cipher_transformation_name": schema.StringAttribute{
				Description: "The algorithm name used to produce this cipher, e.g. AES/CBC/PKCS5Padding.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"initialization_vector_length_bits": schema.Int64Attribute{
				Description: "The initialization vector length of the cipher in bits.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"key_id": schema.StringAttribute{
				Description: "The unique system-generated identifier for the Secret Key.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_compromised": schema.BoolAttribute{
				Description: "If the key is compromised, an administrator may set this flag to immediately trigger the creation of a new secret key. After the new key is generated, the value of this property will be reset to false.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"symmetric_key": schema.SetAttribute{
				Description: "The symmetric key that is used for both encryption of plain text and decryption of cipher text. This stores the secret key for each server instance encrypted with that server's inter-server certificate.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"key_length_bits": schema.Int64Attribute{
				Description: "The length of the key in bits.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a CipherSecretKeyResponse object into the model struct
func readCipherSecretKeyResponse(ctx context.Context, r *client.CipherSecretKeyResponse, state *cipherSecretKeyResourceModel, expectedValues *cipherSecretKeyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceName = expectedValues.ServerInstanceName
	state.CipherTransformationName = internaltypes.StringTypeOrNil(r.CipherTransformationName, true)
	state.InitializationVectorLengthBits = internaltypes.Int64TypeOrNil(r.InitializationVectorLengthBits)
	state.KeyID = types.StringValue(r.KeyID)
	state.IsCompromised = internaltypes.BoolTypeOrNil(r.IsCompromised)
	state.SymmetricKey = internaltypes.GetStringSet(r.SymmetricKey)
	state.KeyLengthBits = types.Int64Value(r.KeyLengthBits)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCipherSecretKeyOperations(plan cipherSecretKeyResourceModel, state cipherSecretKeyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.CipherTransformationName, state.CipherTransformationName, "cipher-transformation-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitializationVectorLengthBits, state.InitializationVectorLengthBits, "initialization-vector-length-bits")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyID, state.KeyID, "key-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IsCompromised, state.IsCompromised, "is-compromised")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SymmetricKey, state.SymmetricKey, "symmetric-key")
	operations.AddInt64OperationIfNecessary(&ops, plan.KeyLengthBits, state.KeyLengthBits, "key-length-bits")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *cipherSecretKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cipherSecretKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherSecretKeyApi.GetCipherSecretKey(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Secret Key", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state cipherSecretKeyResourceModel
	readCipherSecretKeyResponse(ctx, readResponse, &state, &plan, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherSecretKeyApi.UpdateCipherSecretKey(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString())
	ops := createCipherSecretKeyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherSecretKeyApi.UpdateCipherSecretKeyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Cipher Secret Key", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCipherSecretKeyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *cipherSecretKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cipherSecretKeyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherSecretKeyApi.GetCipherSecretKey(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Secret Key", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCipherSecretKeyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *cipherSecretKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cipherSecretKeyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state cipherSecretKeyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.CipherSecretKeyApi.UpdateCipherSecretKey(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := createCipherSecretKeyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherSecretKeyApi.UpdateCipherSecretKeyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Cipher Secret Key", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCipherSecretKeyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *cipherSecretKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *cipherSecretKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [server-instance-name]/[cipher-secret-key-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_instance_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
