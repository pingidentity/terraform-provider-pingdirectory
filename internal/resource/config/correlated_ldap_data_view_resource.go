package config

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &correlatedLdapDataViewResource{}
	_ resource.ResourceWithConfigure   = &correlatedLdapDataViewResource{}
	_ resource.ResourceWithImportState = &correlatedLdapDataViewResource{}
	_ resource.Resource                = &defaultCorrelatedLdapDataViewResource{}
	_ resource.ResourceWithConfigure   = &defaultCorrelatedLdapDataViewResource{}
	_ resource.ResourceWithImportState = &defaultCorrelatedLdapDataViewResource{}
)

// Create a Correlated Ldap Data View resource
func NewCorrelatedLdapDataViewResource() resource.Resource {
	return &correlatedLdapDataViewResource{}
}

func NewDefaultCorrelatedLdapDataViewResource() resource.Resource {
	return &defaultCorrelatedLdapDataViewResource{}
}

// correlatedLdapDataViewResource is the resource implementation.
type correlatedLdapDataViewResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCorrelatedLdapDataViewResource is the resource implementation.
type defaultCorrelatedLdapDataViewResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *correlatedLdapDataViewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_correlated_ldap_data_view"
}

func (r *defaultCorrelatedLdapDataViewResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_correlated_ldap_data_view"
}

// Configure adds the provider configured client to the resource.
func (r *correlatedLdapDataViewResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultCorrelatedLdapDataViewResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type correlatedLdapDataViewResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	ScimResourceTypeName          types.String `tfsdk:"scim_resource_type_name"`
	StructuralLDAPObjectclass     types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass      types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	IncludeBaseDN                 types.String `tfsdk:"include_base_dn"`
	IncludeFilter                 types.Set    `tfsdk:"include_filter"`
	IncludeOperationalAttribute   types.Set    `tfsdk:"include_operational_attribute"`
	CreateDNPattern               types.String `tfsdk:"create_dn_pattern"`
	PrimaryCorrelationAttribute   types.String `tfsdk:"primary_correlation_attribute"`
	SecondaryCorrelationAttribute types.String `tfsdk:"secondary_correlation_attribute"`
}

// GetSchema defines the schema for the resource.
func (r *correlatedLdapDataViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	correlatedLdapDataViewSchema(ctx, req, resp, false)
}

func (r *defaultCorrelatedLdapDataViewResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	correlatedLdapDataViewSchema(ctx, req, resp, true)
}

func correlatedLdapDataViewSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Correlated Ldap Data View.",
		Attributes: map[string]schema.Attribute{
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this Correlated LDAP Data View.",
				Required:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this Correlated LDAP Data View.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory that can be accessed by this Correlated LDAP Data View.",
				Required:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this Correlated LDAP Data View.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_operational_attribute": schema.SetAttribute{
				Description: "Specifies the set of operational LDAP attributes to be provided by this Correlated LDAP Data View.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"create_dn_pattern": schema.StringAttribute{
				Description: "Specifies the template to use for the DN when creating new entries.",
				Optional:    true,
			},
			"primary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the parent SCIM Resource Type whose value will be used to match objects in the Correlated LDAP Data View. If multiple correlation attributes are required they may be created using additional correlation-attribute-pairs.",
				Required:    true,
			},
			"secondary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Correlated LDAP Data View whose value will be matched with the primary-correlation-attribute. If multiple correlation attributes are required they may be specified by creating additional correlation-attribute-pairs.",
				Required:    true,
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id", "scim_resource_type_name"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for correlated-ldap-data-view correlated-ldap-data-view
func addOptionalCorrelatedLdapDataViewFields(ctx context.Context, addRequest *client.AddCorrelatedLdapDataViewRequest, plan correlatedLdapDataViewResourceModel) {
	if internaltypes.IsDefined(plan.AuxiliaryLDAPObjectclass) {
		var slice []string
		plan.AuxiliaryLDAPObjectclass.ElementsAs(ctx, &slice, false)
		addRequest.AuxiliaryLDAPObjectclass = slice
	}
	if internaltypes.IsDefined(plan.IncludeFilter) {
		var slice []string
		plan.IncludeFilter.ElementsAs(ctx, &slice, false)
		addRequest.IncludeFilter = slice
	}
	if internaltypes.IsDefined(plan.IncludeOperationalAttribute) {
		var slice []string
		plan.IncludeOperationalAttribute.ElementsAs(ctx, &slice, false)
		addRequest.IncludeOperationalAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CreateDNPattern) {
		addRequest.CreateDNPattern = plan.CreateDNPattern.ValueStringPointer()
	}
}

// Read a CorrelatedLdapDataViewResponse object into the model struct
func readCorrelatedLdapDataViewResponse(ctx context.Context, r *client.CorrelatedLdapDataViewResponse, state *correlatedLdapDataViewResourceModel, expectedValues *correlatedLdapDataViewResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScimResourceTypeName = expectedValues.ScimResourceTypeName
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.IncludeBaseDN = types.StringValue(r.IncludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.IncludeOperationalAttribute = internaltypes.GetStringSet(r.IncludeOperationalAttribute)
	state.CreateDNPattern = internaltypes.StringTypeOrNil(r.CreateDNPattern, internaltypes.IsEmptyString(expectedValues.CreateDNPattern))
	state.PrimaryCorrelationAttribute = types.StringValue(r.PrimaryCorrelationAttribute)
	state.SecondaryCorrelationAttribute = types.StringValue(r.SecondaryCorrelationAttribute)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCorrelatedLdapDataViewOperations(plan correlatedLdapDataViewResourceModel, state correlatedLdapDataViewResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.StructuralLDAPObjectclass, state.StructuralLDAPObjectclass, "structural-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuxiliaryLDAPObjectclass, state.AuxiliaryLDAPObjectclass, "auxiliary-ldap-objectclass")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeBaseDN, state.IncludeBaseDN, "include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeOperationalAttribute, state.IncludeOperationalAttribute, "include-operational-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.CreateDNPattern, state.CreateDNPattern, "create-dn-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.PrimaryCorrelationAttribute, state.PrimaryCorrelationAttribute, "primary-correlation-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.SecondaryCorrelationAttribute, state.SecondaryCorrelationAttribute, "secondary-correlation-attribute")
	return ops
}

// Create a correlated-ldap-data-view correlated-ldap-data-view
func (r *correlatedLdapDataViewResource) CreateCorrelatedLdapDataView(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan correlatedLdapDataViewResourceModel) (*correlatedLdapDataViewResourceModel, error) {
	addRequest := client.NewAddCorrelatedLdapDataViewRequest(plan.Id.ValueString(),
		plan.StructuralLDAPObjectclass.ValueString(),
		plan.IncludeBaseDN.ValueString(),
		plan.PrimaryCorrelationAttribute.ValueString(),
		plan.SecondaryCorrelationAttribute.ValueString())
	addOptionalCorrelatedLdapDataViewFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CorrelatedLdapDataViewApi.AddCorrelatedLdapDataView(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.ScimResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddCorrelatedLdapDataViewRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.CorrelatedLdapDataViewApi.AddCorrelatedLdapDataViewExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Correlated Ldap Data View", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state correlatedLdapDataViewResourceModel
	readCorrelatedLdapDataViewResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *correlatedLdapDataViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan correlatedLdapDataViewResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateCorrelatedLdapDataView(ctx, req, resp, plan)
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
func (r *defaultCorrelatedLdapDataViewResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan correlatedLdapDataViewResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CorrelatedLdapDataViewApi.GetCorrelatedLdapDataView(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Correlated Ldap Data View", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state correlatedLdapDataViewResourceModel
	readCorrelatedLdapDataViewResponse(ctx, readResponse, &state, &plan, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CorrelatedLdapDataViewApi.UpdateCorrelatedLdapDataView(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ScimResourceTypeName.ValueString())
	ops := createCorrelatedLdapDataViewOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CorrelatedLdapDataViewApi.UpdateCorrelatedLdapDataViewExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Correlated Ldap Data View", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCorrelatedLdapDataViewResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *correlatedLdapDataViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCorrelatedLdapDataView(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCorrelatedLdapDataViewResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCorrelatedLdapDataView(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCorrelatedLdapDataView(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state correlatedLdapDataViewResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CorrelatedLdapDataViewApi.GetCorrelatedLdapDataView(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Correlated Ldap Data View", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCorrelatedLdapDataViewResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *correlatedLdapDataViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCorrelatedLdapDataView(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCorrelatedLdapDataViewResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCorrelatedLdapDataView(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCorrelatedLdapDataView(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan correlatedLdapDataViewResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state correlatedLdapDataViewResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CorrelatedLdapDataViewApi.UpdateCorrelatedLdapDataView(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString(), plan.ScimResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createCorrelatedLdapDataViewOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CorrelatedLdapDataViewApi.UpdateCorrelatedLdapDataViewExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Correlated Ldap Data View", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCorrelatedLdapDataViewResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCorrelatedLdapDataViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *correlatedLdapDataViewResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state correlatedLdapDataViewResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CorrelatedLdapDataViewApi.DeleteCorrelatedLdapDataViewExecute(r.apiClient.CorrelatedLdapDataViewApi.DeleteCorrelatedLdapDataView(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ScimResourceTypeName.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Correlated Ldap Data View", err, httpResp)
		return
	}
}

func (r *correlatedLdapDataViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCorrelatedLdapDataView(ctx, req, resp)
}

func (r *defaultCorrelatedLdapDataViewResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCorrelatedLdapDataView(ctx, req, resp)
}

func importCorrelatedLdapDataView(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [scim-resource-type-name]/[correlated-ldap-data-view-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scim_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
