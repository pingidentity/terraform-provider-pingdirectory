package restresourcetype

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &groupRestResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &groupRestResourceTypeResource{}
	_ resource.ResourceWithImportState = &groupRestResourceTypeResource{}
	_ resource.Resource                = &defaultGroupRestResourceTypeResource{}
	_ resource.ResourceWithConfigure   = &defaultGroupRestResourceTypeResource{}
	_ resource.ResourceWithImportState = &defaultGroupRestResourceTypeResource{}
)

// Create a Group Rest Resource Type resource
func NewGroupRestResourceTypeResource() resource.Resource {
	return &groupRestResourceTypeResource{}
}

func NewDefaultGroupRestResourceTypeResource() resource.Resource {
	return &defaultGroupRestResourceTypeResource{}
}

// groupRestResourceTypeResource is the resource implementation.
type groupRestResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGroupRestResourceTypeResource is the resource implementation.
type defaultGroupRestResourceTypeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groupRestResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_rest_resource_type"
}

func (r *defaultGroupRestResourceTypeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_group_rest_resource_type"
}

// Configure adds the provider configured client to the resource.
func (r *groupRestResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGroupRestResourceTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type groupRestResourceTypeResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	MembersColumnName              types.String `tfsdk:"members_column_name"`
	NonmembersColumnName           types.String `tfsdk:"nonmembers_column_name"`
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
}

// GetSchema defines the schema for the resource.
func (r *groupRestResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groupRestResourceTypeSchema(ctx, req, resp, false)
}

func (r *defaultGroupRestResourceTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groupRestResourceTypeSchema(ctx, req, resp, true)
}

func groupRestResourceTypeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Group Rest Resource Type.",
		Attributes: map[string]schema.Attribute{
			"members_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group member column that will be displayed in the Delegated Admin UI",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"nonmembers_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group nonmember column that will be displayed in the Delegated Admin UI",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"create_rdn_attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type to be used as the RDN of new resources.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"post_create_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a new resource is created. The values are only set at creation time. Subsequent modifications to attributes in the constructed attribute value-pattern are not propagated here.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"update_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a resource is updated. The constructed values replace any existing values of the attribute.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"delegated_admin_report_size_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that may be included in a report.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add config validators
func (r groupRestResourceTypeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("parent_resource_type"),
			path.MatchRoot("parent_dn"),
		),
	}
}

// Add optional fields to create request
func addOptionalGroupRestResourceTypeFields(ctx context.Context, addRequest *client.AddGroupRestResourceTypeRequest, plan groupRestResourceTypeResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MembersColumnName) {
		addRequest.MembersColumnName = plan.MembersColumnName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.NonmembersColumnName) {
		addRequest.NonmembersColumnName = plan.NonmembersColumnName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
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
		addRequest.ParentDN = plan.ParentDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ParentResourceType) {
		addRequest.ParentResourceType = plan.ParentResourceType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RelativeDNFromParentResource) {
		addRequest.RelativeDNFromParentResource = plan.RelativeDNFromParentResource.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CreateRDNAttributeType) {
		addRequest.CreateRDNAttributeType = plan.CreateRDNAttributeType.ValueStringPointer()
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
		addRequest.DisplayName = plan.DisplayName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchFilterPattern) {
		addRequest.SearchFilterPattern = plan.SearchFilterPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrimaryDisplayAttributeType) {
		addRequest.PrimaryDisplayAttributeType = plan.PrimaryDisplayAttributeType.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DelegatedAdminSearchSizeLimit) {
		addRequest.DelegatedAdminSearchSizeLimit = plan.DelegatedAdminSearchSizeLimit.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DelegatedAdminReportSizeLimit) {
		addRequest.DelegatedAdminReportSizeLimit = plan.DelegatedAdminReportSizeLimit.ValueInt64Pointer()
	}
}

// Read a GroupRestResourceTypeResponse object into the model struct
func readGroupRestResourceTypeResponse(ctx context.Context, r *client.GroupRestResourceTypeResponse, state *groupRestResourceTypeResourceModel, expectedValues *groupRestResourceTypeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.MembersColumnName = internaltypes.StringTypeOrNil(r.MembersColumnName, internaltypes.IsEmptyString(expectedValues.MembersColumnName))
	state.NonmembersColumnName = internaltypes.StringTypeOrNil(r.NonmembersColumnName, internaltypes.IsEmptyString(expectedValues.NonmembersColumnName))
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
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroupRestResourceTypeOperations(plan groupRestResourceTypeResourceModel, state groupRestResourceTypeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.MembersColumnName, state.MembersColumnName, "members-column-name")
	operations.AddStringOperationIfNecessary(&ops, plan.NonmembersColumnName, state.NonmembersColumnName, "nonmembers-column-name")
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
	return ops
}

// Create a new resource
func (r *groupRestResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupRestResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGroupRestResourceTypeRequest(plan.Id.ValueString(),
		[]client.EnumgroupRestResourceTypeSchemaUrn{client.ENUMGROUPRESTRESOURCETYPESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0REST_RESOURCE_TYPEGROUP},
		plan.Enabled.ValueBool(),
		plan.ResourceEndpoint.ValueString(),
		plan.StructuralLDAPObjectclass.ValueString(),
		plan.SearchBaseDN.ValueString())
	addOptionalGroupRestResourceTypeFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.RestResourceTypeApi.AddRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddRestResourceTypeRequest(
		client.AddGroupRestResourceTypeRequestAsAddRestResourceTypeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.RestResourceTypeApi.AddRestResourceTypeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Group Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state groupRestResourceTypeResourceModel
	readGroupRestResourceTypeResponse(ctx, addResponse.GroupRestResourceTypeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGroupRestResourceTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupRestResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RestResourceTypeApi.GetRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Group Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groupRestResourceTypeResourceModel
	readGroupRestResourceTypeResponse(ctx, readResponse.GroupRestResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.RestResourceTypeApi.UpdateRestResourceType(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGroupRestResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.RestResourceTypeApi.UpdateRestResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Group Rest Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroupRestResourceTypeResponse(ctx, updateResponse.GroupRestResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *groupRestResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroupRestResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroupRestResourceTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroupRestResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGroupRestResourceType(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state groupRestResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.RestResourceTypeApi.GetRestResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Group Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGroupRestResourceTypeResponse(ctx, readResponse.GroupRestResourceTypeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *groupRestResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroupRestResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroupRestResourceTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroupRestResourceType(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGroupRestResourceType(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan groupRestResourceTypeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groupRestResourceTypeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.RestResourceTypeApi.UpdateRestResourceType(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGroupRestResourceTypeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.RestResourceTypeApi.UpdateRestResourceTypeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Group Rest Resource Type", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroupRestResourceTypeResponse(ctx, updateResponse.GroupRestResourceTypeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGroupRestResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groupRestResourceTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groupRestResourceTypeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.RestResourceTypeApi.DeleteRestResourceTypeExecute(r.apiClient.RestResourceTypeApi.DeleteRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Group Rest Resource Type", err, httpResp)
		return
	}
}

func (r *groupRestResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroupRestResourceType(ctx, req, resp)
}

func (r *defaultGroupRestResourceTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroupRestResourceType(ctx, req, resp)
}

func importGroupRestResourceType(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
