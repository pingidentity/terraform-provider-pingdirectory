package trustmanagerprovider

import (
	"context"
	"terraform-provider-pingdirectory/internal/utils"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fileBasedTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &fileBasedTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &fileBasedTrustManagerProviderResource{}
)

// Create a File Based Trust Manager Provider resource
func NewFileBasedTrustManagerProviderResource() resource.Resource {
	return &fileBasedTrustManagerProviderResource{}
}

// fileBasedTrustManagerProviderResource is the resource implementation.
type fileBasedTrustManagerProviderResource struct {
	providerConfig utils.ProviderConfiguration
	apiClient      *client.APIClient
}

// fileBasedTrustManagerProviderResourceModel maps the resource schema data.
type fileBasedTrustManagerProviderResourceModel struct {
	Name                            types.String `tfsdk:"name"`
	TrustStoreFile                  types.String `tfsdk:"trust_store_file"`
	TrustStoreType                  types.String `tfsdk:"trust_store_type"`
	TrustStorePin                   types.String `tfsdk:"trust_store_pin"`
	TrustStorePinFile               types.String `tfsdk:"trust_store_pin_file"`
	TrustStorePinPassphraseProvider types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
	IncludeJVMDefaultIssuers        types.Bool   `tfsdk:"include_jvm_default_issuers"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
}

// Metadata returns the resource type name.
func (r *fileBasedTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_based_trust_manager_provider"
}

// GetSchema defines the schema for the resource.
func (r *fileBasedTrustManagerProviderResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a File Based Trust Manager Provider.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the Trust Manager Provider.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"trust_store_file": {
				Description: "Specifies the path to the file containing the trust information. It can be an absolute path or a path that is relative to the Directory Server instance root.",
				Type:        types.StringType,
				Required:    true,
			},
			"trust_store_type": {
				Description: "Specifies the format for the data in the trust store file.",
				Type:        types.StringType,
				Optional:    true,
			},
			"trust_store_pin": {
				Description: "Specifies the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
			},
			"trust_store_pin_file": {
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Type:        types.StringType,
				Optional:    true,
			},
			"trust_store_pin_passphrase_provider": {
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Trust Manager Provider.",
				Type:        types.StringType,
				Optional:    true,
			},
			"enabled": {
				Description: "Indicate whether the Trust Manager Provider is enabled for use.",
				Type:        types.BoolType,
				Required:    true,
			},
			// Optional boolean fields must be Computed because PD gives them a default value
			"include_jvm_default_issuers": {
				Description: "Indicates whether certificates issued by an authority included in the JVM's set of default issuers should be automatically trusted, even if they would not otherwise be trusted by this provider.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the Trust Manager Provider.",
				Type:        types.StringType,
				Computed:    true,
				Required:    false,
				Optional:    false,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *fileBasedTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(utils.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Add optional fields to create request
func addOptionalFileBasedTrustManagerProviderFields(addRequest *client.AddFileBasedTrustManagerProviderRequest, plan fileBasedTrustManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if utils.IsNonEmptyString(plan.TrustStoreType) {
		stringVal := plan.TrustStoreType.ValueString()
		addRequest.TrustStoreType = &stringVal
	}
	if utils.IsNonEmptyString(plan.TrustStorePin) {
		stringVal := plan.TrustStorePin.ValueString()
		addRequest.TrustStorePin = &stringVal
	}
	if utils.IsNonEmptyString(plan.TrustStorePinFile) {
		stringVal := plan.TrustStorePinFile.ValueString()
		addRequest.TrustStorePinFile = &stringVal
	}
	if utils.IsNonEmptyString(plan.TrustStorePinPassphraseProvider) {
		stringVal := plan.TrustStorePinPassphraseProvider.ValueString()
		addRequest.TrustStorePinPassphraseProvider = &stringVal
	}
	// Non string values just have to be defined
	if utils.IsDefinedBool(plan.IncludeJVMDefaultIssuers) {
		boolVal := plan.IncludeJVMDefaultIssuers.ValueBool()
		addRequest.IncludeJVMDefaultIssuers = &boolVal
	}
}

// Create a new resource
func (r *fileBasedTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileBasedTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedTrustManagerProviderSchemaUrn{client.ENUMFILEBASEDTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERFILE_BASED},
		plan.TrustStoreFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedTrustManagerProviderFields(addRequest, plan)
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(utils.BasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddFileBasedTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return
	}

	// Read the response into the state
	var state fileBasedTrustManagerProviderResourceModel
	readFileBasedTrustManagerProviderResponse(trustManagerResponse.FileBasedTrustManagerProviderResponse, &state, &plan)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a FileBasedTrustManagerProviderResponse object into the model struct
func readFileBasedTrustManagerProviderResponse(r *client.FileBasedTrustManagerProviderResponse,
	state *fileBasedTrustManagerProviderResourceModel, expectedValues *fileBasedTrustManagerProviderResourceModel) {
	state.Name = types.StringValue(r.Id)
	state.TrustStoreFile = types.StringValue(r.TrustStoreFile)
	// If a plan was provided and is using an empty string, use that for a nil string in the response.
	// To PingDirectory, nil and empty string is equivalent, but to Terraform they are distinct. So we
	// just want to match whatever is in the plan here.
	state.TrustStoreType = utils.StringTypeOrNil(r.TrustStoreType, utils.IsEmptyString(expectedValues.TrustStoreType))
	state.TrustStorePin = utils.StringTypeOrNil(r.TrustStorePin, utils.IsEmptyString(expectedValues.TrustStorePin))
	state.TrustStorePinFile = utils.StringTypeOrNil(r.TrustStorePinFile, utils.IsEmptyString(expectedValues.TrustStorePinFile))
	state.TrustStorePinPassphraseProvider = utils.StringTypeOrNil(r.TrustStorePinPassphraseProvider, utils.IsEmptyString(expectedValues.TrustStorePinPassphraseProvider))

	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = utils.BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
}

// Read resource information
func (r *fileBasedTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state fileBasedTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		utils.BasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
		return
	}

	// Read the response into the state
	readFileBasedTrustManagerProviderResponse(trustManagerResponse.FileBasedTrustManagerProviderResponse, &state, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func createFileBasedTrustManagerProviderOperations(plan fileBasedTrustManagerProviderResourceModel, state fileBasedTrustManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	utils.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	utils.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	utils.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	utils.AddStringOperationIfNecessary(&ops, plan.TrustStorePinFile, state.TrustStorePinFile, "trust-store-pin-file")
	utils.AddStringOperationIfNecessary(&ops, plan.TrustStorePinPassphraseProvider, state.TrustStorePinPassphraseProvider, "trust-store-pin-passphrase-provider")
	utils.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	utils.AddBoolOperationIfNecessary(&ops, plan.IncludeJVMDefaultIssuers, state.IncludeJVMDefaultIssuers, "include-jvm-default-issuers")
	return ops
}

// Update a resource
func (r *fileBasedTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan fileBasedTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileBasedTrustManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(utils.BasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createFileBasedTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))

		trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Trust Manager Provider", err, httpResp)
			return
		}

		// Read the response
		readFileBasedTrustManagerProviderResponse(trustManagerResponse.FileBasedTrustManagerProviderResponse, &state, &plan)
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
func (r *fileBasedTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileBasedTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(
		r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProvider(utils.BasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		utils.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Trust Manager Provider", err, httpResp)
		return
	}
}

func (r *fileBasedTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
