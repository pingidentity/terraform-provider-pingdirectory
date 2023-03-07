package backend

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &changelogBackendResource{}
	_ resource.ResourceWithConfigure   = &changelogBackendResource{}
	_ resource.ResourceWithImportState = &changelogBackendResource{}
)

// Create a Changelog Backend resource
func NewChangelogBackendResource() resource.Resource {
	return &changelogBackendResource{}
}

// changelogBackendResource is the resource implementation.
type changelogBackendResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *changelogBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_changelog_backend"
}

// Configure adds the provider configured client to the resource.
func (r *changelogBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type changelogBackendResourceModel struct {
	Id                                          types.String `tfsdk:"id"`
	LastUpdated                                 types.String `tfsdk:"last_updated"`
	Notifications                               types.Set    `tfsdk:"notifications"`
	RequiredActions                             types.Set    `tfsdk:"required_actions"`
	BaseDN                                      types.Set    `tfsdk:"base_dn"`
	DbDirectory                                 types.String `tfsdk:"db_directory"`
	DbDirectoryPermissions                      types.String `tfsdk:"db_directory_permissions"`
	DbCachePercent                              types.Int64  `tfsdk:"db_cache_percent"`
	JeProperty                                  types.Set    `tfsdk:"je_property"`
	ChangelogWriteBatchSize                     types.Int64  `tfsdk:"changelog_write_batch_size"`
	ChangelogPurgeBatchSize                     types.Int64  `tfsdk:"changelog_purge_batch_size"`
	ChangelogWriteQueueCapacity                 types.Int64  `tfsdk:"changelog_write_queue_capacity"`
	IndexIncludeAttribute                       types.Set    `tfsdk:"index_include_attribute"`
	IndexExcludeAttribute                       types.Set    `tfsdk:"index_exclude_attribute"`
	ChangelogMaximumAge                         types.String `tfsdk:"changelog_maximum_age"`
	TargetDatabaseSize                          types.String `tfsdk:"target_database_size"`
	ChangelogEntryIncludeBaseDN                 types.Set    `tfsdk:"changelog_entry_include_base_dn"`
	ChangelogEntryExcludeBaseDN                 types.Set    `tfsdk:"changelog_entry_exclude_base_dn"`
	ChangelogEntryIncludeFilter                 types.Set    `tfsdk:"changelog_entry_include_filter"`
	ChangelogEntryExcludeFilter                 types.Set    `tfsdk:"changelog_entry_exclude_filter"`
	ChangelogIncludeAttribute                   types.Set    `tfsdk:"changelog_include_attribute"`
	ChangelogExcludeAttribute                   types.Set    `tfsdk:"changelog_exclude_attribute"`
	ChangelogDeletedEntryIncludeAttribute       types.Set    `tfsdk:"changelog_deleted_entry_include_attribute"`
	ChangelogDeletedEntryExcludeAttribute       types.Set    `tfsdk:"changelog_deleted_entry_exclude_attribute"`
	ChangelogIncludeKeyAttribute                types.Set    `tfsdk:"changelog_include_key_attribute"`
	ChangelogMaxBeforeAfterValues               types.Int64  `tfsdk:"changelog_max_before_after_values"`
	WriteLastmodAttributes                      types.Bool   `tfsdk:"write_lastmod_attributes"`
	UseReversibleForm                           types.Bool   `tfsdk:"use_reversible_form"`
	IncludeVirtualAttributes                    types.Set    `tfsdk:"include_virtual_attributes"`
	ApplyAccessControlsToChangelogEntryContents types.Bool   `tfsdk:"apply_access_controls_to_changelog_entry_contents"`
	ReportExcludedChangelogAttributes           types.String `tfsdk:"report_excluded_changelog_attributes"`
	SoftDeleteEntryIncludedOperation            types.Set    `tfsdk:"soft_delete_entry_included_operation"`
	BackendID                                   types.String `tfsdk:"backend_id"`
	Description                                 types.String `tfsdk:"description"`
	Enabled                                     types.Bool   `tfsdk:"enabled"`
	SetDegradedAlertWhenDisabled                types.Bool   `tfsdk:"set_degraded_alert_when_disabled"`
	ReturnUnavailableWhenDisabled               types.Bool   `tfsdk:"return_unavailable_when_disabled"`
	NotificationManager                         types.String `tfsdk:"notification_manager"`
}

// GetSchema defines the schema for the resource.
func (r *changelogBackendResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Changelog Backend.",
		Attributes: map[string]schema.Attribute{
			"base_dn": schema.SetAttribute{
				Description: "Specifies the base DN(s) for the data that the backend handles.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"db_directory": schema.StringAttribute{
				Description: "Specifies the path to the filesystem directory that is used to hold the Berkeley DB Java Edition database files containing the data for this backend. The files for this backend are stored in a sub-directory named after the backend-id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_directory_permissions": schema.StringAttribute{
				Description: "Specifies the permissions that should be applied to the directory containing the backend database files and to directories and files created during backup of the backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"db_cache_percent": schema.Int64Attribute{
				Description: "Specifies the percentage of JVM memory to allocate to the changelog database cache.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"je_property": schema.SetAttribute{
				Description: "Specifies the database and environment properties for the Berkeley DB Java Edition database for this changelog backend.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_write_batch_size": schema.Int64Attribute{
				Description: "Specifies the number of changelog entries written in a single database transaction.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"changelog_purge_batch_size": schema.Int64Attribute{
				Description: "Specifies the number of changelog entries purged in a single database transaction.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"changelog_write_queue_capacity": schema.Int64Attribute{
				Description: "Specifies the capacity of the changelog write queue in number of changes.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"index_include_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types are to be specifically included in the set of attribute indexes maintained on the changelog. If this property does not have any values then no attribute types are indexed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"index_exclude_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types are to be specifically excluded from the set of attribute indexes maintained on the changelog. This property is useful when the index-include-attribute property contains one of the special values \"*\" and \"+\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_maximum_age": schema.StringAttribute{
				Description: "Changes are guaranteed to be maintained in the changelog database for at least this duration. Setting target-database-size can allow additional changes to be maintained up to the configured size on disk.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"target_database_size": schema.StringAttribute{
				Description: "The changelog database is allowed to grow up to this size on disk even if changes are older than the configured changelog-maximum-age.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"changelog_entry_include_base_dn": schema.SetAttribute{
				Description: "The base DNs for branches in the data for which to record changes in the changelog.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_entry_exclude_base_dn": schema.SetAttribute{
				Description: "The base DNs for branches in the data for which no changelog records should be generated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_entry_include_filter": schema.SetAttribute{
				Description: "A filter that indicates which changelog entries should actually be stored in the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_entry_exclude_filter": schema.SetAttribute{
				Description: "A filter that indicates which changelog entries should be excluded from the changelog. Note that this filter is evaluated against the changelog entry itself and not against the entry that was the target of the change referenced by the changelog entry. This filter may target any attributes that appear in changelog entries with the exception of the changeNumber and entry-size-bytes attributes, since they will not be known at the time of the filter evaluation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_include_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types will be included in a changelog entry for ADD and MODIFY operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_exclude_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be excluded in a changelog entry for ADD and MODIFY operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_deleted_entry_include_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be included in a changelog entry for DELETE operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_deleted_entry_exclude_attribute": schema.SetAttribute{
				Description: "Specifies a set of attribute types that should be excluded from a changelog entry for DELETE operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_include_key_attribute": schema.SetAttribute{
				Description: "Specifies which attribute types will be included in a changelog entry on every change.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"changelog_max_before_after_values": schema.Int64Attribute{
				Description: "This controls whether all attribute values for a modified attribute (even those values that have not changed) will be included in the changelog entry. If the number of attribute values does not exceed this limit, then all values for the modified attribute will be included in the changelog entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"write_lastmod_attributes": schema.BoolAttribute{
				Description: "Specifies whether values of creatorsName, createTimestamp, modifiersName and modifyTimestamp attributes will be written to changelog entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Specifies whether the changelog should provide enough information to be able to revert the changes if desired.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_virtual_attributes": schema.SetAttribute{
				Description: "Specifies the changelog entry elements (if any) in which virtual attributes should be included.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"apply_access_controls_to_changelog_entry_contents": schema.BoolAttribute{
				Description: "Indicates whether the contents of changelog entries should be subject to access control and sensitive attribute evaluation such that the contents of attributes like changes, deletedEntryAttrs, ds-changelog-entry-key-attr-values, ds-changelog-before-values, and ds-changelog-after-values may be altered based on attributes the user can see in the target entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"report_excluded_changelog_attributes": schema.StringAttribute{
				Description: "Indicates whether changelog entries that have been altered by applying access controls should include additional information about any attributes that may have been removed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_entry_included_operation": schema.SetAttribute{
				Description: "Specifies which operations performed on soft-deleted entries will appear in the changelog.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"backend_id": schema.StringAttribute{
				Description: "Specifies a name to identify the associated backend.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Backend",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the backend is enabled in the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"set_degraded_alert_when_disabled": schema.BoolAttribute{
				Description: "Determines whether the Directory Server enters a DEGRADED state (and sends a corresponding alert) when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"return_unavailable_when_disabled": schema.BoolAttribute{
				Description: "Determines whether any LDAP operation that would use this Backend is to return UNAVAILABLE when this Backend is disabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"notification_manager": schema.StringAttribute{
				Description: "Specifies a notification manager for changes resulting from operations processed through this Backend",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Read a ChangelogBackendResponse object into the model struct
func readChangelogBackendResponse(ctx context.Context, r *client.ChangelogBackendResponse, state *changelogBackendResourceModel, expectedValues *changelogBackendResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.DbDirectory = internaltypes.StringTypeOrNil(r.DbDirectory, true)
	state.DbDirectoryPermissions = internaltypes.StringTypeOrNil(r.DbDirectoryPermissions, true)
	state.DbCachePercent = internaltypes.Int64TypeOrNil(r.DbCachePercent)
	state.JeProperty = internaltypes.GetStringSet(r.JeProperty)
	state.ChangelogWriteBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogWriteBatchSize)
	state.ChangelogPurgeBatchSize = internaltypes.Int64TypeOrNil(r.ChangelogPurgeBatchSize)
	state.ChangelogWriteQueueCapacity = internaltypes.Int64TypeOrNil(r.ChangelogWriteQueueCapacity)
	state.IndexIncludeAttribute = internaltypes.GetStringSet(r.IndexIncludeAttribute)
	state.IndexExcludeAttribute = internaltypes.GetStringSet(r.IndexExcludeAttribute)
	state.ChangelogMaximumAge = types.StringValue(r.ChangelogMaximumAge)
	config.CheckMismatchedPDFormattedAttributes("changelog_maximum_age",
		expectedValues.ChangelogMaximumAge, state.ChangelogMaximumAge, diagnostics)
	state.TargetDatabaseSize = internaltypes.StringTypeOrNil(r.TargetDatabaseSize, true)
	config.CheckMismatchedPDFormattedAttributes("target_database_size",
		expectedValues.TargetDatabaseSize, state.TargetDatabaseSize, diagnostics)
	state.ChangelogEntryIncludeBaseDN = internaltypes.GetStringSet(r.ChangelogEntryIncludeBaseDN)
	state.ChangelogEntryExcludeBaseDN = internaltypes.GetStringSet(r.ChangelogEntryExcludeBaseDN)
	state.ChangelogEntryIncludeFilter = internaltypes.GetStringSet(r.ChangelogEntryIncludeFilter)
	state.ChangelogEntryExcludeFilter = internaltypes.GetStringSet(r.ChangelogEntryExcludeFilter)
	state.ChangelogIncludeAttribute = internaltypes.GetStringSet(r.ChangelogIncludeAttribute)
	state.ChangelogExcludeAttribute = internaltypes.GetStringSet(r.ChangelogExcludeAttribute)
	state.ChangelogDeletedEntryIncludeAttribute = internaltypes.GetStringSet(r.ChangelogDeletedEntryIncludeAttribute)
	state.ChangelogDeletedEntryExcludeAttribute = internaltypes.GetStringSet(r.ChangelogDeletedEntryExcludeAttribute)
	state.ChangelogIncludeKeyAttribute = internaltypes.GetStringSet(r.ChangelogIncludeKeyAttribute)
	state.ChangelogMaxBeforeAfterValues = internaltypes.Int64TypeOrNil(r.ChangelogMaxBeforeAfterValues)
	state.WriteLastmodAttributes = internaltypes.BoolTypeOrNil(r.WriteLastmodAttributes)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.IncludeVirtualAttributes = internaltypes.GetStringSet(
		client.StringSliceEnumbackendIncludeVirtualAttributesProp(r.IncludeVirtualAttributes))
	state.ApplyAccessControlsToChangelogEntryContents = internaltypes.BoolTypeOrNil(r.ApplyAccessControlsToChangelogEntryContents)
	state.ReportExcludedChangelogAttributes = internaltypes.StringTypeOrNil(
		client.StringPointerEnumbackendReportExcludedChangelogAttributesProp(r.ReportExcludedChangelogAttributes), true)
	state.SoftDeleteEntryIncludedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumbackendSoftDeleteEntryIncludedOperationProp(r.SoftDeleteEntryIncludedOperation))
	state.BackendID = types.StringValue(r.BackendID)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SetDegradedAlertWhenDisabled = internaltypes.BoolTypeOrNil(r.SetDegradedAlertWhenDisabled)
	state.ReturnUnavailableWhenDisabled = internaltypes.BoolTypeOrNil(r.ReturnUnavailableWhenDisabled)
	state.NotificationManager = internaltypes.StringTypeOrNil(r.NotificationManager, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createChangelogBackendOperations(plan changelogBackendResourceModel, state changelogBackendResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectory, state.DbDirectory, "db-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DbDirectoryPermissions, state.DbDirectoryPermissions, "db-directory-permissions")
	operations.AddInt64OperationIfNecessary(&ops, plan.DbCachePercent, state.DbCachePercent, "db-cache-percent")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.JeProperty, state.JeProperty, "je-property")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogWriteBatchSize, state.ChangelogWriteBatchSize, "changelog-write-batch-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogPurgeBatchSize, state.ChangelogPurgeBatchSize, "changelog-purge-batch-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogWriteQueueCapacity, state.ChangelogWriteQueueCapacity, "changelog-write-queue-capacity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexIncludeAttribute, state.IndexIncludeAttribute, "index-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IndexExcludeAttribute, state.IndexExcludeAttribute, "index-exclude-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.ChangelogMaximumAge, state.ChangelogMaximumAge, "changelog-maximum-age")
	operations.AddStringOperationIfNecessary(&ops, plan.TargetDatabaseSize, state.TargetDatabaseSize, "target-database-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryIncludeBaseDN, state.ChangelogEntryIncludeBaseDN, "changelog-entry-include-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryExcludeBaseDN, state.ChangelogEntryExcludeBaseDN, "changelog-entry-exclude-base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryIncludeFilter, state.ChangelogEntryIncludeFilter, "changelog-entry-include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogEntryExcludeFilter, state.ChangelogEntryExcludeFilter, "changelog-entry-exclude-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogIncludeAttribute, state.ChangelogIncludeAttribute, "changelog-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogExcludeAttribute, state.ChangelogExcludeAttribute, "changelog-exclude-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogDeletedEntryIncludeAttribute, state.ChangelogDeletedEntryIncludeAttribute, "changelog-deleted-entry-include-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogDeletedEntryExcludeAttribute, state.ChangelogDeletedEntryExcludeAttribute, "changelog-deleted-entry-exclude-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangelogIncludeKeyAttribute, state.ChangelogIncludeKeyAttribute, "changelog-include-key-attribute")
	operations.AddInt64OperationIfNecessary(&ops, plan.ChangelogMaxBeforeAfterValues, state.ChangelogMaxBeforeAfterValues, "changelog-max-before-after-values")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteLastmodAttributes, state.WriteLastmodAttributes, "write-lastmod-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseReversibleForm, state.UseReversibleForm, "use-reversible-form")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeVirtualAttributes, state.IncludeVirtualAttributes, "include-virtual-attributes")
	operations.AddBoolOperationIfNecessary(&ops, plan.ApplyAccessControlsToChangelogEntryContents, state.ApplyAccessControlsToChangelogEntryContents, "apply-access-controls-to-changelog-entry-contents")
	operations.AddStringOperationIfNecessary(&ops, plan.ReportExcludedChangelogAttributes, state.ReportExcludedChangelogAttributes, "report-excluded-changelog-attributes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SoftDeleteEntryIncludedOperation, state.SoftDeleteEntryIncludedOperation, "soft-delete-entry-included-operation")
	operations.AddStringOperationIfNecessary(&ops, plan.BackendID, state.BackendID, "backend-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.SetDegradedAlertWhenDisabled, state.SetDegradedAlertWhenDisabled, "set-degraded-alert-when-disabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.ReturnUnavailableWhenDisabled, state.ReturnUnavailableWhenDisabled, "return-unavailable-when-disabled")
	operations.AddStringOperationIfNecessary(&ops, plan.NotificationManager, state.NotificationManager, "notification-manager")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *changelogBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan changelogBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Changelog Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state changelogBackendResourceModel
	readChangelogBackendResponse(ctx, readResponse.ChangelogBackendResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.BackendApi.UpdateBackend(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())
	ops := createChangelogBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Changelog Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readChangelogBackendResponse(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *changelogBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state changelogBackendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.BackendApi.GetBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.BackendID.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Changelog Backend", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readChangelogBackendResponse(ctx, readResponse.ChangelogBackendResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *changelogBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan changelogBackendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state changelogBackendResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.BackendApi.UpdateBackend(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.BackendID.ValueString())

	// Determine what update operations are necessary
	ops := createChangelogBackendOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.BackendApi.UpdateBackendExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Changelog Backend", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readChangelogBackendResponse(ctx, updateResponse.ChangelogBackendResponse, &state, &plan, &resp.Diagnostics)
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
func (r *changelogBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *changelogBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to backend_id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("backend_id"), req, resp)
}
