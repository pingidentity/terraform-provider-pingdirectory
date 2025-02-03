// Copyright © 2025 Ping Identity Corporation

package consentservice

import (
	"context"

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
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consentServiceResource{}
	_ resource.ResourceWithConfigure   = &consentServiceResource{}
	_ resource.ResourceWithImportState = &consentServiceResource{}
)

// Create a Consent Service resource
func NewConsentServiceResource() resource.Resource {
	return &consentServiceResource{}
}

// consentServiceResource is the resource implementation.
type consentServiceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consentServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_consent_service"
}

// Configure adds the provider configured client to the resource.
func (r *consentServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consentServiceResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	Type                        types.String `tfsdk:"type"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	BaseDN                      types.String `tfsdk:"base_dn"`
	BindDN                      types.String `tfsdk:"bind_dn"`
	SearchSizeLimit             types.Int64  `tfsdk:"search_size_limit"`
	ConsentRecordIdentityMapper types.Set    `tfsdk:"consent_record_identity_mapper"`
	ServiceAccountDN            types.Set    `tfsdk:"service_account_dn"`
	UnprivilegedConsentScope    types.String `tfsdk:"unprivileged_consent_scope"`
	PrivilegedConsentScope      types.String `tfsdk:"privileged_consent_scope"`
	Audience                    types.String `tfsdk:"audience"`
}

// GetSchema defines the schema for the resource.
func (r *consentServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Consent Service.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Consent Service resource. Options are ['consent-service']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"consent-service"}...),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Consent Service is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.StringAttribute{
				Description: "The base DN under which consent records are stored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_dn": schema.StringAttribute{
				Description: "The DN of an internal service account used by the Consent Service to make internal LDAP requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"search_size_limit": schema.Int64Attribute{
				Description: "The maximum number of consent resources that may be returned from a search request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"consent_record_identity_mapper": schema.SetAttribute{
				Description: "If specified, the Identity Mapper(s) that may be used to map consent record subject and actor values to DNs. This is typically only needed if privileged API clients will be used.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"service_account_dn": schema.SetAttribute{
				Description: "The set of account DNs that the Consent Service will consider to be privileged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"unprivileged_consent_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Consent Service for unprivileged clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"privileged_consent_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Consent Service if the client is to be considered privileged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"audience": schema.StringAttribute{
				Description: "A string or URI that identifies the Consent Service in the context of OAuth2 authorization.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ConsentServiceResponse object into the model struct
func readConsentServiceResponse(ctx context.Context, r *client.ConsentServiceResponse, state *consentServiceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent-service")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.StringTypeOrNil(r.BaseDN, true)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, true)
	state.SearchSizeLimit = internaltypes.Int64TypeOrNil(r.SearchSizeLimit)
	state.ConsentRecordIdentityMapper = internaltypes.GetStringSet(r.ConsentRecordIdentityMapper)
	state.ServiceAccountDN = internaltypes.GetStringSet(r.ServiceAccountDN)
	state.UnprivilegedConsentScope = internaltypes.StringTypeOrNil(r.UnprivilegedConsentScope, true)
	state.PrivilegedConsentScope = internaltypes.StringTypeOrNil(r.PrivilegedConsentScope, true)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConsentServiceOperations(plan consentServiceResourceModel, state consentServiceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDN, state.BindDN, "bind-dn")
	operations.AddInt64OperationIfNecessary(&ops, plan.SearchSizeLimit, state.SearchSizeLimit, "search-size-limit")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ConsentRecordIdentityMapper, state.ConsentRecordIdentityMapper, "consent-record-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ServiceAccountDN, state.ServiceAccountDN, "service-account-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.UnprivilegedConsentScope, state.UnprivilegedConsentScope, "unprivileged-consent-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivilegedConsentScope, state.PrivilegedConsentScope, "privileged-consent-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Audience, state.Audience, "audience")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *consentServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consentServiceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentServiceAPI.GetConsentService(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Service", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state consentServiceResourceModel
	readConsentServiceResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConsentServiceAPI.UpdateConsentService(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createConsentServiceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentServiceAPI.UpdateConsentServiceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Service", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentServiceResponse(ctx, updateResponse, &state, &resp.Diagnostics)
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *consentServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state consentServiceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentServiceAPI.GetConsentService(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Service", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentServiceResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *consentServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan consentServiceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consentServiceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ConsentServiceAPI.UpdateConsentService(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createConsentServiceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConsentServiceAPI.UpdateConsentServiceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Consent Service", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsentServiceResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *consentServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *consentServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
