package restresourcetype

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &genericRestResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &genericRestResourceTypeResource{}
	_ resource.ResourceWithImportState = &genericRestResourceTypeResource{}
)

// Create a Generic Rest Resource Type resource
func NewGenericRestResourceTypeResource() resource.Resource {
	return &genericRestResourceTypeResource{}
}

// genericRestResourceTypeResource is the resource implementation.
type genericRestResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *genericRestResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_generic_rest_resource_type"
}

// Configure adds the provider configured client to the resource.
func (r *genericRestResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type genericRestResourceTypeResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	ResourceEndpoint               types.String `tfsdk:"resource_endpoint"`
	StructuralLDAPObjectclass      types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass       types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	SearchBaseDN                   types.String `tfsdk:"search_base_dn"`
	IncludeFilter                  types.Set    `tfsdk:"include_filter"`
	ParentDN                       types.String `tfsdk:"parent_dn"`
	ParentResourceType             types.String `tfsdk:"parent_resource_type"`
	RelativeDNFromParentResource   types.String `tfsdk:"relative_dn_from_parent_resource"`
	CreateRDNAttributeType         types.String `tfsdk:"create_rdn_attribute_type"`
	PostCreateConstructedAttribute types.Set    `tfsdk:"post_create_constructed_attribute"`
	UpdateConstructedAttribute     types.Set    `tfsdk:"update_constructed_attribute"`
	DisplayName                    types.String `tfsdk:"display_name"`
	SearchFilterPattern            types.String `tfsdk:"search_filter_pattern"`
	PrimaryDisplayAttributeType    types.String `tfsdk:"primary_display_attribute_type"`
	DelegatedAdminSearchSizeLimit  types.Int64  `tfsdk:"delegated_admin_search_size_limit"`
	DelegatedAdminReportSizeLimit  types.Int64  `tfsdk:"delegated_admin_report_size_limit"`
	MembersColumnName              types.String `tfsdk:"members_column_name"`
	NonmembersColumnName           types.String `tfsdk:"nonmembers_column_name"`
}

// GetSchema defines the schema for the resource.
func (r *genericRestResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Generic Rest Resource Type.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this REST Resource Type",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the REST Resource Type is enabled.",
				Required:    true,
			},
			"resource_endpoint": schema.StringAttribute{
				Description: "The HTTP addressable endpoint of this REST Resource Type relative to a REST API base URL. Do not include a leading '/'.",
				Required:    true,
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this REST Resource Type.",
				Required:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this REST Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory where resources of this type are located.",
				Required:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this REST Resource Type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"parent_dn": schema.StringAttribute{
				Description: "Specifies the DN of the parent entry for new resources of this type, when a parent resource is not provided by the app. The parent DN must be at or below the search base of this resource type.",
				Optional:    true,
			},
			"parent_resource_type": schema.StringAttribute{
				Description: "Specifies the name of another resource type which may be a parent of new resources of this type. The search base DN of the parent resource type must be at or above the search base DN of this resource type.",
				Optional:    true,
			},
			"relative_dn_from_parent_resource": schema.StringAttribute{
				Description: "Specifies a template for a relative DN from the parent resource which identifies the parent entry for a new resource of this type. If this property is not specified then new resources are created immediately below the parent resource or parent DN.",
				Optional:    true,
			},
			"create_rdn_attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type to be used as the RDN of new resources.",
				Optional:    true,
			},
			"post_create_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a new resource is created. The values are only set at creation time. Subsequent modifications to attributes in the constructed attribute value-pattern are not propagated here.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"update_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a resource is updated. The constructed values replace any existing values of the attribute.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this REST Resource Type.",
				Optional:    true,
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "Specifies the LDAP filter that should be used when searching for resources matching provided search text. All attribute types in the filter pattern referencing the search text must have a Delegated Admin Attribute definition.",
				Optional:    true,
			},
			"primary_display_attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type which is the primary display attribute. This attribute type must be in the search filter pattern and must have a Delegated Admin Attribute definition.",
				Optional:    true,
			},
			"delegated_admin_search_size_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that may be returned from a search request.",
				Optional:    true,
				Computed:    true,
			},
			"delegated_admin_report_size_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that may be included in a report.",
				Optional:    true,
				Computed:    true,
			},
			"members_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group member column that will be displayed in the Delegated Admin UI",
				Optional:    true,
				Computed:    true,
			},
			"nonmembers_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group nonmember column that will be displayed in the Delegated Admin UI",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalGenericRestResourceTypeFields(ctx context.Context, addRequest *client.AddGenericRestResourceTypeRequest, plan genericRestResourceTypeResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
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
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ParentDN) {
		stringVal := plan.ParentDN.ValueString()
		addRequest.ParentDN = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ParentResourceType) {
		stringVal := plan.ParentResourceType.ValueString()
		addRequest.ParentResourceType = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RelativeDNFromParentResource) {
		stringVal := plan.RelativeDNFromParentResource.ValueString()
		addRequest.RelativeDNFromParentResource = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CreateRDNAttributeType) {
		stringVal := plan.CreateRDNAttributeType.ValueString()
		addRequest.CreateRDNAttributeType = &stringVal
	}
	if internaltypes.IsDefined(plan.PostCreateConstructedAttribute) {
		var slice []string
		plan.PostCreateConstructedAttribute.ElementsAs(ctx, &slice, false)
		addRequest.PostCreateConstructedAttribute = slice
	}
	if internaltypes.IsDefined(plan.UpdateConstructedAttribute) {
		var slice []string
		plan.UpdateConstructedAttribute.ElementsAs(ctx, &slice, false)
		addRequest.UpdateConstructedAttribute = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DisplayName) {
		stringVal := plan.DisplayName.ValueString()
		addRequest.DisplayName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchFilterPattern) {
		stringVal := plan.SearchFilterPattern.ValueString()
		addRequest.SearchFilterPattern = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrimaryDisplayAttributeType) {
		stringVal := plan.PrimaryDisplayAttributeType.ValueString()
		addRequest.PrimaryDisplayAttributeType = &stringVal
	}
	if internaltypes.IsDefined(plan.DelegatedAdminSearchSizeLimit) {
		intVal := int32(plan.DelegatedAdminSearchSizeLimit.ValueInt64())
		addRequest.DelegatedAdminSearchSizeLimit = &intVal
	}
	if internaltypes.IsDefined(plan.DelegatedAdminReportSizeLimit) {
		intVal := int32(plan.DelegatedAdminReportSizeLimit.ValueInt64())
		addRequest.DelegatedAdminReportSizeLimit = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MembersColumnName) {
		stringVal := plan.MembersColumnName.ValueString()
		addRequest.MembersColumnName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NonmembersColumnName) {
		stringVal := plan.NonmembersColumnName.ValueString()
		addRequest.NonmembersColumnName = &stringVal
	}
}

// Read a GenericRestResourceTypeResponse object into the model struct
func readGenericRestResourceTypeResponse(ctx context.Context, r *client.GenericRestResourceTypeResponse, state *genericRestResourceTypeResourceModel, expectedValues *genericRestResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ResourceEndpoint = types.StringValue(r.ResourceEndpoint)
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.SearchBaseDN = types.StringValue(r.SearchBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ParentDN = internaltypes.StringTypeOrNil(r.ParentDN, internaltypes.IsEmptyString(expectedValues.ParentDN))
	state.ParentResourceType = internaltypes.StringTypeOrNil(r.ParentResourceType, internaltypes.IsEmptyString(expectedValues.ParentResourceType))
	state.RelativeDNFromParentResource = internaltypes.StringTypeOrNil(r.RelativeDNFromParentResource, internaltypes.IsEmptyString(expectedValues.RelativeDNFromParentResource))
	state.CreateRDNAttributeType = internaltypes.StringTypeOrNil(r.CreateRDNAttributeType, internaltypes.IsEmptyString(expectedValues.CreateRDNAttributeType))
	state.PostCreateConstructedAttribute = internaltypes.GetStringSet(r.PostCreateConstructedAttribute)
	state.UpdateConstructedAttribute = internaltypes.GetStringSet(r.UpdateConstructedAttribute)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, internaltypes.IsEmptyString(expectedValues.DisplayName))
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, internaltypes.IsEmptyString(expectedValues.SearchFilterPattern))
	state.PrimaryDisplayAttributeType = internaltypes.StringTypeOrNil(r.PrimaryDisplayAttributeType, internaltypes.IsEmptyString(expectedValues.PrimaryDisplayAttributeType))
	state.DelegatedAdminSearchSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminSearchSizeLimit)
	state.DelegatedAdminReportSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminReportSizeLimit)
	state.MembersColumnName = internaltypes.StringTypeOrNil(r.MembersColumnName, internaltypes.IsEmptyString(expectedValues.MembersColumnName))
	state.NonmembersColumnName = internaltypes.StringTypeOrNil(r.NonmembersColumnName, internaltypes.IsEmptyString(expectedValues.NonmembersColumnName))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGenericRestResourceTypeOperations(plan genericRestResourceTypeResourceModel, state genericRestResourceTypeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.ResourceEndpoint, state.ResourceEndpoint, "resource-endpoint")
	operations.AddStringOperationIfNecessary(&ops, plan.StructuralLDAPObjectclass, state.StructuralLDAPObjectclass, "structural-ldap-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AuxiliaryLDAPObjectclass, state.AuxiliaryLDAPObjectclass, "auxiliary-ldap-objectclass")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchBaseDN, state.SearchBaseDN, "search-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.ParentDN, state.ParentDN, "parent-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.ParentResourceType, state.ParentResourceType, "parent-resource-type")
	operations.AddStringOperationIfNecessary(&ops, plan.RelativeDNFromParentResource, state.RelativeDNFromParentResource, "relative-dn-from-parent-resource")
	operations.AddStringOperationIfNecessary(&ops, plan.CreateRDNAttributeType, state.CreateRDNAttributeType, "create-rdn-attribute-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PostCreateConstructedAttribute, state.PostCreateConstructedAttribute, "post-create-constructed-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.UpdateConstructedAttribute, state.UpdateConstructedAttribute, "update-constructed-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchFilterPattern, state.SearchFilterPattern, "search-filter-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.PrimaryDisplayAttributeType, state.PrimaryDisplayAttributeType, "primary-display-attribute-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.DelegatedAdminSearchSizeLimit, state.DelegatedAdminSearchSizeLimit, "delegated-admin-search-size-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.DelegatedAdminReportSizeLimit, state.DelegatedAdminReportSizeLimit, "delegated-admin-report-size-limit")
	operations.AddStringOperationIfNecessary(&ops, plan.MembersColumnName, state.MembersColumnName, "members-column-name")
	operations.AddStringOperationIfNecessary(&ops, plan.NonmembersColumnName, state.NonmembersColumnName, "nonmembers-column-name")
	return ops
}

// Create a new resource
func (r *genericRestResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan genericRestResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGenericRestResourceTypeRequest(plan.Id.ValueString(),
		[]client.EnumgenericRestResourceTypeSchemaUrn{client.ENUMGENERICRESTRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REST_RESOURCE_TYPEGENERIC},
		plan.Enabled.ValueBool(),
		plan.ResourceEndpoint.ValueString(),
		plan.StructuralLDAPObjectclass.ValueString(),
		plan.SearchBaseDN.ValueString())
	addOptionalGenericRestResourceTypeFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RestResourceTypeApi.AddRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRestResourceTypeRequest(
		client.AddGenericRestResourceTypeRequestAsAddRestResourceTypeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RestResourceTypeApi.AddRestResourceTypeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Generic Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state genericRestResourceTypeResourceModel
	readGenericRestResourceTypeResponse(ctx, addResponse.GenericRestResourceTypeResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *genericRestResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state genericRestResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RestResourceTypeApi.GetRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generic Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGenericRestResourceTypeResponse(ctx, readResponse.GenericRestResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *genericRestResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan genericRestResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state genericRestResourceTypeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.RestResourceTypeApi.UpdateRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGenericRestResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RestResourceTypeApi.UpdateRestResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generic Rest Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGenericRestResourceTypeResponse(ctx, updateResponse.GenericRestResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *genericRestResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state genericRestResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RestResourceTypeApi.DeleteRestResourceTypeExecute(r.apiClient.RestResourceTypeApi.DeleteRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Generic Rest Resource Type", err, httpResp)
		return
	}
}

func (r *genericRestResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
