// Copyright © 2025 Ping Identity Corporation

package logfieldmapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logFieldMappingResource{}
	_ resource.ResourceWithConfigure   = &logFieldMappingResource{}
	_ resource.ResourceWithImportState = &logFieldMappingResource{}
	_ resource.Resource                = &defaultLogFieldMappingResource{}
	_ resource.ResourceWithConfigure   = &defaultLogFieldMappingResource{}
	_ resource.ResourceWithImportState = &defaultLogFieldMappingResource{}
)

// Create a Log Field Mapping resource
func NewLogFieldMappingResource() resource.Resource {
	return &logFieldMappingResource{}
}

func NewDefaultLogFieldMappingResource() resource.Resource {
	return &defaultLogFieldMappingResource{}
}

// logFieldMappingResource is the resource implementation.
type logFieldMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogFieldMappingResource is the resource implementation.
type defaultLogFieldMappingResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logFieldMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_field_mapping"
}

func (r *defaultLogFieldMappingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_field_mapping"
}

// Configure adds the provider configured client to the resource.
func (r *logFieldMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogFieldMappingResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFieldMappingResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
	Type                                types.String `tfsdk:"type"`
	LogFieldTimestamp                   types.String `tfsdk:"log_field_timestamp"`
	LogFieldConnectionID                types.String `tfsdk:"log_field_connection_id"`
	LogFieldStartupid                   types.String `tfsdk:"log_field_startupid"`
	LogFieldProductName                 types.String `tfsdk:"log_field_product_name"`
	LogFieldCategory                    types.String `tfsdk:"log_field_category"`
	LogFieldSeverity                    types.String `tfsdk:"log_field_severity"`
	LogFieldInstanceName                types.String `tfsdk:"log_field_instance_name"`
	LogFieldOperationID                 types.String `tfsdk:"log_field_operation_id"`
	LogFieldMessageType                 types.String `tfsdk:"log_field_message_type"`
	LogFieldOperationType               types.String `tfsdk:"log_field_operation_type"`
	LogFieldMessageID                   types.String `tfsdk:"log_field_message_id"`
	LogFieldResultCode                  types.String `tfsdk:"log_field_result_code"`
	LogFieldMessage                     types.String `tfsdk:"log_field_message"`
	LogFieldOrigin                      types.String `tfsdk:"log_field_origin"`
	LogFieldRequesterDN                 types.String `tfsdk:"log_field_requester_dn"`
	LogFieldDisconnectReason            types.String `tfsdk:"log_field_disconnect_reason"`
	LogFieldDeleteOldRDN                types.String `tfsdk:"log_field_delete_old_rdn"`
	LogFieldAuthenticatedUserDN         types.String `tfsdk:"log_field_authenticated_user_dn"`
	LogFieldProcessingTime              types.String `tfsdk:"log_field_processing_time"`
	LogFieldRequestedAttributes         types.String `tfsdk:"log_field_requested_attributes"`
	LogFieldSASLMechanismName           types.String `tfsdk:"log_field_sasl_mechanism_name"`
	LogFieldNewRDN                      types.String `tfsdk:"log_field_new_rdn"`
	LogFieldBaseDN                      types.String `tfsdk:"log_field_base_dn"`
	LogFieldBindDN                      types.String `tfsdk:"log_field_bind_dn"`
	LogFieldMatchedDN                   types.String `tfsdk:"log_field_matched_dn"`
	LogFieldRequesterIPAddress          types.String `tfsdk:"log_field_requester_ip_address"`
	LogFieldAuthenticationType          types.String `tfsdk:"log_field_authentication_type"`
	LogFieldNewSuperiorDN               types.String `tfsdk:"log_field_new_superior_dn"`
	LogFieldFilter                      types.String `tfsdk:"log_field_filter"`
	LogFieldAlternateAuthorizationDN    types.String `tfsdk:"log_field_alternate_authorization_dn"`
	LogFieldEntryDN                     types.String `tfsdk:"log_field_entry_dn"`
	LogFieldEntriesReturned             types.String `tfsdk:"log_field_entries_returned"`
	LogFieldAuthenticationFailureID     types.String `tfsdk:"log_field_authentication_failure_id"`
	LogFieldRequestOID                  types.String `tfsdk:"log_field_request_oid"`
	LogFieldResponseOID                 types.String `tfsdk:"log_field_response_oid"`
	LogFieldTargetProtocol              types.String `tfsdk:"log_field_target_protocol"`
	LogFieldTargetPort                  types.String `tfsdk:"log_field_target_port"`
	LogFieldTargetAddress               types.String `tfsdk:"log_field_target_address"`
	LogFieldTargetAttribute             types.String `tfsdk:"log_field_target_attribute"`
	LogFieldTargetHost                  types.String `tfsdk:"log_field_target_host"`
	LogFieldProtocolVersion             types.String `tfsdk:"log_field_protocol_version"`
	LogFieldProtocolName                types.String `tfsdk:"log_field_protocol_name"`
	LogFieldAuthenticationFailureReason types.String `tfsdk:"log_field_authentication_failure_reason"`
	LogFieldAdditionalInformation       types.String `tfsdk:"log_field_additional_information"`
	LogFieldUnindexed                   types.String `tfsdk:"log_field_unindexed"`
	LogFieldScope                       types.String `tfsdk:"log_field_scope"`
	LogFieldReferralUrls                types.String `tfsdk:"log_field_referral_urls"`
	LogFieldSourceAddress               types.String `tfsdk:"log_field_source_address"`
	LogFieldMessageIDToAbandon          types.String `tfsdk:"log_field_message_id_to_abandon"`
	LogFieldResponseControls            types.String `tfsdk:"log_field_response_controls"`
	LogFieldRequestControls             types.String `tfsdk:"log_field_request_controls"`
	LogFieldIntermediateClientResult    types.String `tfsdk:"log_field_intermediate_client_result"`
	LogFieldIntermediateClientRequest   types.String `tfsdk:"log_field_intermediate_client_request"`
	LogFieldReplicationChangeID         types.String `tfsdk:"log_field_replication_change_id"`
	Description                         types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *logFieldMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFieldMappingSchema(ctx, req, resp, false)
}

func (r *defaultLogFieldMappingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFieldMappingSchema(ctx, req, resp, true)
}

func logFieldMappingSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Field Mapping.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Mapping resource. Options are ['access', 'error']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"access", "error"}...),
				},
			},
			"log_field_timestamp": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The time that the operation was processed. When the `type` attribute is set to `error`: The time that the log message was generated.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The time that the operation was processed.\n  - `error`: The time that the log message was generated.",
				Optional:            true,
			},
			"log_field_connection_id": schema.StringAttribute{
				Description: "The connection ID assigned to the client connection.",
				Optional:    true,
			},
			"log_field_startupid": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The startup ID for the Directory Server. A different value will be generated each time the server is started, and it may be used to distinguish between operations with the same connection ID and operation ID across server restarts. When the `type` attribute is set to `error`: The startup ID for the Directory Server. A different value will be generated each time the server is started.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The startup ID for the Directory Server. A different value will be generated each time the server is started, and it may be used to distinguish between operations with the same connection ID and operation ID across server restarts.\n  - `error`: The startup ID for the Directory Server. A different value will be generated each time the server is started.",
				Optional:            true,
			},
			"log_field_product_name": schema.StringAttribute{
				Description: "The name for this Directory Server product, which may be used to identify which product was used to log the message if multiple products log to the same database table.",
				Optional:    true,
			},
			"log_field_category": schema.StringAttribute{
				Description: "The category for the log message.",
				Optional:    true,
			},
			"log_field_severity": schema.StringAttribute{
				Description: "The severity for the log message.",
				Optional:    true,
			},
			"log_field_instance_name": schema.StringAttribute{
				Description: "A name that uniquely identifies this Directory Server instance, which may be used to identify which instance was used to log the message if multiple server instances log to the same database table.",
				Optional:    true,
			},
			"log_field_operation_id": schema.StringAttribute{
				Description: "The operation ID for the operation processed by the server.",
				Optional:    true,
			},
			"log_field_message_type": schema.StringAttribute{
				Description: "The type of log message. Message types may include \"CONNECT\", \"DISCONNECT\", \"FORWARD\", \"RESULT\", \"ENTRY\", or \"REFERENCE\".",
				Optional:    true,
			},
			"log_field_operation_type": schema.StringAttribute{
				Description: "The type of operation that was processed.",
				Optional:    true,
			},
			"log_field_message_id": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The message ID included in the client request. When the `type` attribute is set to `error`: The numeric value which uniquely identifies the type of message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The message ID included in the client request.\n  - `error`: The numeric value which uniquely identifies the type of message.",
				Optional:            true,
			},
			"log_field_result_code": schema.StringAttribute{
				Description: "The numeric result code for the operation.",
				Optional:    true,
			},
			"log_field_message": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `access`: The diagnostic message for the operation. When the `type` attribute is set to `error`: The text of the log message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `access`: The diagnostic message for the operation.\n  - `error`: The text of the log message.",
				Optional:            true,
			},
			"log_field_origin": schema.StringAttribute{
				Description: "The origin for the operation. Values may include \"replication\" (if the operation was received via replication), \"internal\" (if it was an internal operation processed by a server component), or \"external\" (if it was a request from a client).",
				Optional:    true,
			},
			"log_field_requester_dn": schema.StringAttribute{
				Description: "The DN of the user that requested the operation.",
				Optional:    true,
			},
			"log_field_disconnect_reason": schema.StringAttribute{
				Description: "The reason that the client connection was closed.",
				Optional:    true,
			},
			"log_field_delete_old_rdn": schema.StringAttribute{
				Description: "Indicates whether the old RDN values should be removed from an entry while processing a modify DN operation.",
				Optional:    true,
			},
			"log_field_authenticated_user_dn": schema.StringAttribute{
				Description: "The DN of the user that authenticated to the server.",
				Optional:    true,
			},
			"log_field_processing_time": schema.StringAttribute{
				Description: "The length of time (in milliseconds with microsecond accuracy) required to process the operation.",
				Optional:    true,
			},
			"log_field_requested_attributes": schema.StringAttribute{
				Description: "The set of requested attributes for the search operation.",
				Optional:    true,
			},
			"log_field_sasl_mechanism_name": schema.StringAttribute{
				Description: "The name of the SASL mechanism used to authenticate.",
				Optional:    true,
			},
			"log_field_new_rdn": schema.StringAttribute{
				Description: "The new RDN to use for a modify DN operation.",
				Optional:    true,
			},
			"log_field_base_dn": schema.StringAttribute{
				Description: "The base DN for a search operation.",
				Optional:    true,
			},
			"log_field_bind_dn": schema.StringAttribute{
				Description: "The bind DN for a bind operation.",
				Optional:    true,
			},
			"log_field_matched_dn": schema.StringAttribute{
				Description: "The DN of the superior entry closest to the DN specified by the client.",
				Optional:    true,
			},
			"log_field_requester_ip_address": schema.StringAttribute{
				Description: "The IP address of the client that requested the operation.",
				Optional:    true,
			},
			"log_field_authentication_type": schema.StringAttribute{
				Description: "The type of authentication requested by the client.",
				Optional:    true,
			},
			"log_field_new_superior_dn": schema.StringAttribute{
				Description: "The new superior DN from a modify DN operation.",
				Optional:    true,
			},
			"log_field_filter": schema.StringAttribute{
				Description: "The filter from a search operation.",
				Optional:    true,
			},
			"log_field_alternate_authorization_dn": schema.StringAttribute{
				Description: "The DN of the alternate authorization identity used when processing the operation.",
				Optional:    true,
			},
			"log_field_entry_dn": schema.StringAttribute{
				Description: "The DN of the entry targeted by the operation.",
				Optional:    true,
			},
			"log_field_entries_returned": schema.StringAttribute{
				Description: "The number of search result entries returned to the client.",
				Optional:    true,
			},
			"log_field_authentication_failure_id": schema.StringAttribute{
				Description: "The numeric identifier for the authentication failure reason.",
				Optional:    true,
			},
			"log_field_request_oid": schema.StringAttribute{
				Description: "The OID of an extended request.",
				Optional:    true,
			},
			"log_field_response_oid": schema.StringAttribute{
				Description: "The OID of an extended response.",
				Optional:    true,
			},
			"log_field_target_protocol": schema.StringAttribute{
				Description: "The protocol used when forwarding the request to a backend server.",
				Optional:    true,
			},
			"log_field_target_port": schema.StringAttribute{
				Description: "The network port of the Directory Server to which the client connection has been established, or of the backend server to which the request has been forwarded.",
				Optional:    true,
			},
			"log_field_target_address": schema.StringAttribute{
				Description: "The network address of the Directory Server to which the client connection has been established.",
				Optional:    true,
			},
			"log_field_target_attribute": schema.StringAttribute{
				Description: "The name of the attribute targeted by a compare operation.",
				Optional:    true,
			},
			"log_field_target_host": schema.StringAttribute{
				Description: "The address of the server to which the request has been forwarded.",
				Optional:    true,
			},
			"log_field_protocol_version": schema.StringAttribute{
				Description: "The protocol version used by the client when communicating with the Directory Server.",
				Optional:    true,
			},
			"log_field_protocol_name": schema.StringAttribute{
				Description: "The name of the protocol the client is using to communicate with the Directory Server.",
				Optional:    true,
			},
			"log_field_authentication_failure_reason": schema.StringAttribute{
				Description: "A message explaining the reason that the authentication attempt failed.",
				Optional:    true,
			},
			"log_field_additional_information": schema.StringAttribute{
				Description: "Additional information about the operation that was processed which was not returned to the client.",
				Optional:    true,
			},
			"log_field_unindexed": schema.StringAttribute{
				Description: "Indicates whether the requested search operation was unindexed.",
				Optional:    true,
			},
			"log_field_scope": schema.StringAttribute{
				Description: "The scope for the search operation.",
				Optional:    true,
			},
			"log_field_referral_urls": schema.StringAttribute{
				Description: "The referral URLs returned to the client.",
				Optional:    true,
			},
			"log_field_source_address": schema.StringAttribute{
				Description: "The address of the client from which the connection was established.",
				Optional:    true,
			},
			"log_field_message_id_to_abandon": schema.StringAttribute{
				Description: "The message ID of the operation to be abandoned.",
				Optional:    true,
			},
			"log_field_response_controls": schema.StringAttribute{
				Description: "The OIDs of the response controls returned to the client.",
				Optional:    true,
			},
			"log_field_request_controls": schema.StringAttribute{
				Description: "The OIDs of the request controls returned to the client.",
				Optional:    true,
			},
			"log_field_intermediate_client_result": schema.StringAttribute{
				Description: "The contents of the intermediate client response control returned to the client.",
				Optional:    true,
			},
			"log_field_intermediate_client_request": schema.StringAttribute{
				Description: "The contents of the intermediate client request control provided by the client.",
				Optional:    true,
			},
			"log_field_replication_change_id": schema.StringAttribute{
				Description: "The replication change ID.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Mapping",
				Optional:    true,
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
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsLogFieldMapping() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_connection_id"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_operation_id"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_message_type"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_operation_type"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_result_code"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_origin"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_requester_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_disconnect_reason"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_delete_old_rdn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_authenticated_user_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_processing_time"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_requested_attributes"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_sasl_mechanism_name"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_new_rdn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_base_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_bind_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_matched_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_requester_ip_address"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_authentication_type"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_new_superior_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_filter"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_alternate_authorization_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_entry_dn"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_entries_returned"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_authentication_failure_id"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_request_oid"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_response_oid"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_target_protocol"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_target_port"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_target_address"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_target_attribute"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_target_host"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_protocol_version"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_protocol_name"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_authentication_failure_reason"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_additional_information"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_unindexed"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_scope"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_referral_urls"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_source_address"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_message_id_to_abandon"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_response_controls"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_request_controls"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_intermediate_client_result"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_intermediate_client_request"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_replication_change_id"),
			path.MatchRoot("type"),
			[]string{"access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_category"),
			path.MatchRoot("type"),
			[]string{"error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_severity"),
			path.MatchRoot("type"),
			[]string{"error"},
		),
	}
}

// Add config validators
func (r logFieldMappingResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogFieldMapping()
}

// Add config validators
func (r defaultLogFieldMappingResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogFieldMapping()
}

// Add optional fields to create request for access log-field-mapping
func addOptionalAccessLogFieldMappingFields(ctx context.Context, addRequest *client.AddAccessLogFieldMappingRequest, plan logFieldMappingResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTimestamp) {
		addRequest.LogFieldTimestamp = plan.LogFieldTimestamp.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldConnectionID) {
		addRequest.LogFieldConnectionID = plan.LogFieldConnectionID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldStartupid) {
		addRequest.LogFieldStartupid = plan.LogFieldStartupid.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProductName) {
		addRequest.LogFieldProductName = plan.LogFieldProductName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldInstanceName) {
		addRequest.LogFieldInstanceName = plan.LogFieldInstanceName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldOperationID) {
		addRequest.LogFieldOperationID = plan.LogFieldOperationID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessageType) {
		addRequest.LogFieldMessageType = plan.LogFieldMessageType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldOperationType) {
		addRequest.LogFieldOperationType = plan.LogFieldOperationType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessageID) {
		addRequest.LogFieldMessageID = plan.LogFieldMessageID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldResultCode) {
		addRequest.LogFieldResultCode = plan.LogFieldResultCode.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessage) {
		addRequest.LogFieldMessage = plan.LogFieldMessage.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldOrigin) {
		addRequest.LogFieldOrigin = plan.LogFieldOrigin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldRequesterDN) {
		addRequest.LogFieldRequesterDN = plan.LogFieldRequesterDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldDisconnectReason) {
		addRequest.LogFieldDisconnectReason = plan.LogFieldDisconnectReason.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldDeleteOldRDN) {
		addRequest.LogFieldDeleteOldRDN = plan.LogFieldDeleteOldRDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAuthenticatedUserDN) {
		addRequest.LogFieldAuthenticatedUserDN = plan.LogFieldAuthenticatedUserDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProcessingTime) {
		addRequest.LogFieldProcessingTime = plan.LogFieldProcessingTime.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldRequestedAttributes) {
		addRequest.LogFieldRequestedAttributes = plan.LogFieldRequestedAttributes.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldSASLMechanismName) {
		addRequest.LogFieldSASLMechanismName = plan.LogFieldSASLMechanismName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldNewRDN) {
		addRequest.LogFieldNewRDN = plan.LogFieldNewRDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBaseDN) {
		addRequest.LogFieldBaseDN = plan.LogFieldBaseDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBindDN) {
		addRequest.LogFieldBindDN = plan.LogFieldBindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMatchedDN) {
		addRequest.LogFieldMatchedDN = plan.LogFieldMatchedDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldRequesterIPAddress) {
		addRequest.LogFieldRequesterIPAddress = plan.LogFieldRequesterIPAddress.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAuthenticationType) {
		addRequest.LogFieldAuthenticationType = plan.LogFieldAuthenticationType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldNewSuperiorDN) {
		addRequest.LogFieldNewSuperiorDN = plan.LogFieldNewSuperiorDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldFilter) {
		addRequest.LogFieldFilter = plan.LogFieldFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAlternateAuthorizationDN) {
		addRequest.LogFieldAlternateAuthorizationDN = plan.LogFieldAlternateAuthorizationDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldEntryDN) {
		addRequest.LogFieldEntryDN = plan.LogFieldEntryDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldEntriesReturned) {
		addRequest.LogFieldEntriesReturned = plan.LogFieldEntriesReturned.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAuthenticationFailureID) {
		addRequest.LogFieldAuthenticationFailureID = plan.LogFieldAuthenticationFailureID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldRequestOID) {
		addRequest.LogFieldRequestOID = plan.LogFieldRequestOID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldResponseOID) {
		addRequest.LogFieldResponseOID = plan.LogFieldResponseOID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTargetProtocol) {
		addRequest.LogFieldTargetProtocol = plan.LogFieldTargetProtocol.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTargetPort) {
		addRequest.LogFieldTargetPort = plan.LogFieldTargetPort.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTargetAddress) {
		addRequest.LogFieldTargetAddress = plan.LogFieldTargetAddress.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTargetAttribute) {
		addRequest.LogFieldTargetAttribute = plan.LogFieldTargetAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTargetHost) {
		addRequest.LogFieldTargetHost = plan.LogFieldTargetHost.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProtocolVersion) {
		addRequest.LogFieldProtocolVersion = plan.LogFieldProtocolVersion.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProtocolName) {
		addRequest.LogFieldProtocolName = plan.LogFieldProtocolName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAuthenticationFailureReason) {
		addRequest.LogFieldAuthenticationFailureReason = plan.LogFieldAuthenticationFailureReason.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldAdditionalInformation) {
		addRequest.LogFieldAdditionalInformation = plan.LogFieldAdditionalInformation.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldUnindexed) {
		addRequest.LogFieldUnindexed = plan.LogFieldUnindexed.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldScope) {
		addRequest.LogFieldScope = plan.LogFieldScope.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldReferralUrls) {
		addRequest.LogFieldReferralUrls = plan.LogFieldReferralUrls.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldSourceAddress) {
		addRequest.LogFieldSourceAddress = plan.LogFieldSourceAddress.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessageIDToAbandon) {
		addRequest.LogFieldMessageIDToAbandon = plan.LogFieldMessageIDToAbandon.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldResponseControls) {
		addRequest.LogFieldResponseControls = plan.LogFieldResponseControls.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldRequestControls) {
		addRequest.LogFieldRequestControls = plan.LogFieldRequestControls.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldIntermediateClientResult) {
		addRequest.LogFieldIntermediateClientResult = plan.LogFieldIntermediateClientResult.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldIntermediateClientRequest) {
		addRequest.LogFieldIntermediateClientRequest = plan.LogFieldIntermediateClientRequest.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldReplicationChangeID) {
		addRequest.LogFieldReplicationChangeID = plan.LogFieldReplicationChangeID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for error log-field-mapping
func addOptionalErrorLogFieldMappingFields(ctx context.Context, addRequest *client.AddErrorLogFieldMappingRequest, plan logFieldMappingResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldTimestamp) {
		addRequest.LogFieldTimestamp = plan.LogFieldTimestamp.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldProductName) {
		addRequest.LogFieldProductName = plan.LogFieldProductName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldInstanceName) {
		addRequest.LogFieldInstanceName = plan.LogFieldInstanceName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldStartupid) {
		addRequest.LogFieldStartupid = plan.LogFieldStartupid.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldCategory) {
		addRequest.LogFieldCategory = plan.LogFieldCategory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldSeverity) {
		addRequest.LogFieldSeverity = plan.LogFieldSeverity.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessageID) {
		addRequest.LogFieldMessageID = plan.LogFieldMessageID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldMessage) {
		addRequest.LogFieldMessage = plan.LogFieldMessage.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logFieldMappingResourceModel) populateAllComputedStringAttributes() {
	if model.LogFieldOrigin.IsUnknown() || model.LogFieldOrigin.IsNull() {
		model.LogFieldOrigin = types.StringValue("")
	}
	if model.LogFieldMessageIDToAbandon.IsUnknown() || model.LogFieldMessageIDToAbandon.IsNull() {
		model.LogFieldMessageIDToAbandon = types.StringValue("")
	}
	if model.LogFieldNewRDN.IsUnknown() || model.LogFieldNewRDN.IsNull() {
		model.LogFieldNewRDN = types.StringValue("")
	}
	if model.LogFieldSourceAddress.IsUnknown() || model.LogFieldSourceAddress.IsNull() {
		model.LogFieldSourceAddress = types.StringValue("")
	}
	if model.LogFieldRequestedAttributes.IsUnknown() || model.LogFieldRequestedAttributes.IsNull() {
		model.LogFieldRequestedAttributes = types.StringValue("")
	}
	if model.LogFieldAlternateAuthorizationDN.IsUnknown() || model.LogFieldAlternateAuthorizationDN.IsNull() {
		model.LogFieldAlternateAuthorizationDN = types.StringValue("")
	}
	if model.LogFieldAuthenticationFailureID.IsUnknown() || model.LogFieldAuthenticationFailureID.IsNull() {
		model.LogFieldAuthenticationFailureID = types.StringValue("")
	}
	if model.LogFieldProtocolName.IsUnknown() || model.LogFieldProtocolName.IsNull() {
		model.LogFieldProtocolName = types.StringValue("")
	}
	if model.LogFieldBindDN.IsUnknown() || model.LogFieldBindDN.IsNull() {
		model.LogFieldBindDN = types.StringValue("")
	}
	if model.LogFieldRequestOID.IsUnknown() || model.LogFieldRequestOID.IsNull() {
		model.LogFieldRequestOID = types.StringValue("")
	}
	if model.LogFieldReplicationChangeID.IsUnknown() || model.LogFieldReplicationChangeID.IsNull() {
		model.LogFieldReplicationChangeID = types.StringValue("")
	}
	if model.LogFieldBaseDN.IsUnknown() || model.LogFieldBaseDN.IsNull() {
		model.LogFieldBaseDN = types.StringValue("")
	}
	if model.LogFieldNewSuperiorDN.IsUnknown() || model.LogFieldNewSuperiorDN.IsNull() {
		model.LogFieldNewSuperiorDN = types.StringValue("")
	}
	if model.LogFieldInstanceName.IsUnknown() || model.LogFieldInstanceName.IsNull() {
		model.LogFieldInstanceName = types.StringValue("")
	}
	if model.LogFieldSeverity.IsUnknown() || model.LogFieldSeverity.IsNull() {
		model.LogFieldSeverity = types.StringValue("")
	}
	if model.LogFieldProcessingTime.IsUnknown() || model.LogFieldProcessingTime.IsNull() {
		model.LogFieldProcessingTime = types.StringValue("")
	}
	if model.LogFieldTargetProtocol.IsUnknown() || model.LogFieldTargetProtocol.IsNull() {
		model.LogFieldTargetProtocol = types.StringValue("")
	}
	if model.LogFieldOperationType.IsUnknown() || model.LogFieldOperationType.IsNull() {
		model.LogFieldOperationType = types.StringValue("")
	}
	if model.LogFieldMessage.IsUnknown() || model.LogFieldMessage.IsNull() {
		model.LogFieldMessage = types.StringValue("")
	}
	if model.LogFieldProtocolVersion.IsUnknown() || model.LogFieldProtocolVersion.IsNull() {
		model.LogFieldProtocolVersion = types.StringValue("")
	}
	if model.LogFieldIntermediateClientResult.IsUnknown() || model.LogFieldIntermediateClientResult.IsNull() {
		model.LogFieldIntermediateClientResult = types.StringValue("")
	}
	if model.LogFieldAuthenticationType.IsUnknown() || model.LogFieldAuthenticationType.IsNull() {
		model.LogFieldAuthenticationType = types.StringValue("")
	}
	if model.LogFieldMatchedDN.IsUnknown() || model.LogFieldMatchedDN.IsNull() {
		model.LogFieldMatchedDN = types.StringValue("")
	}
	if model.LogFieldRequesterIPAddress.IsUnknown() || model.LogFieldRequesterIPAddress.IsNull() {
		model.LogFieldRequesterIPAddress = types.StringValue("")
	}
	if model.LogFieldProductName.IsUnknown() || model.LogFieldProductName.IsNull() {
		model.LogFieldProductName = types.StringValue("")
	}
	if model.LogFieldResponseControls.IsUnknown() || model.LogFieldResponseControls.IsNull() {
		model.LogFieldResponseControls = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.LogFieldUnindexed.IsUnknown() || model.LogFieldUnindexed.IsNull() {
		model.LogFieldUnindexed = types.StringValue("")
	}
	if model.LogFieldRequesterDN.IsUnknown() || model.LogFieldRequesterDN.IsNull() {
		model.LogFieldRequesterDN = types.StringValue("")
	}
	if model.LogFieldResultCode.IsUnknown() || model.LogFieldResultCode.IsNull() {
		model.LogFieldResultCode = types.StringValue("")
	}
	if model.LogFieldTargetHost.IsUnknown() || model.LogFieldTargetHost.IsNull() {
		model.LogFieldTargetHost = types.StringValue("")
	}
	if model.LogFieldOperationID.IsUnknown() || model.LogFieldOperationID.IsNull() {
		model.LogFieldOperationID = types.StringValue("")
	}
	if model.LogFieldScope.IsUnknown() || model.LogFieldScope.IsNull() {
		model.LogFieldScope = types.StringValue("")
	}
	if model.LogFieldMessageType.IsUnknown() || model.LogFieldMessageType.IsNull() {
		model.LogFieldMessageType = types.StringValue("")
	}
	if model.LogFieldTimestamp.IsUnknown() || model.LogFieldTimestamp.IsNull() {
		model.LogFieldTimestamp = types.StringValue("")
	}
	if model.LogFieldAuthenticationFailureReason.IsUnknown() || model.LogFieldAuthenticationFailureReason.IsNull() {
		model.LogFieldAuthenticationFailureReason = types.StringValue("")
	}
	if model.LogFieldDisconnectReason.IsUnknown() || model.LogFieldDisconnectReason.IsNull() {
		model.LogFieldDisconnectReason = types.StringValue("")
	}
	if model.LogFieldDeleteOldRDN.IsUnknown() || model.LogFieldDeleteOldRDN.IsNull() {
		model.LogFieldDeleteOldRDN = types.StringValue("")
	}
	if model.LogFieldTargetAttribute.IsUnknown() || model.LogFieldTargetAttribute.IsNull() {
		model.LogFieldTargetAttribute = types.StringValue("")
	}
	if model.LogFieldIntermediateClientRequest.IsUnknown() || model.LogFieldIntermediateClientRequest.IsNull() {
		model.LogFieldIntermediateClientRequest = types.StringValue("")
	}
	if model.LogFieldRequestControls.IsUnknown() || model.LogFieldRequestControls.IsNull() {
		model.LogFieldRequestControls = types.StringValue("")
	}
	if model.LogFieldEntryDN.IsUnknown() || model.LogFieldEntryDN.IsNull() {
		model.LogFieldEntryDN = types.StringValue("")
	}
	if model.LogFieldFilter.IsUnknown() || model.LogFieldFilter.IsNull() {
		model.LogFieldFilter = types.StringValue("")
	}
	if model.LogFieldAuthenticatedUserDN.IsUnknown() || model.LogFieldAuthenticatedUserDN.IsNull() {
		model.LogFieldAuthenticatedUserDN = types.StringValue("")
	}
	if model.LogFieldConnectionID.IsUnknown() || model.LogFieldConnectionID.IsNull() {
		model.LogFieldConnectionID = types.StringValue("")
	}
	if model.LogFieldTargetPort.IsUnknown() || model.LogFieldTargetPort.IsNull() {
		model.LogFieldTargetPort = types.StringValue("")
	}
	if model.LogFieldAdditionalInformation.IsUnknown() || model.LogFieldAdditionalInformation.IsNull() {
		model.LogFieldAdditionalInformation = types.StringValue("")
	}
	if model.LogFieldTargetAddress.IsUnknown() || model.LogFieldTargetAddress.IsNull() {
		model.LogFieldTargetAddress = types.StringValue("")
	}
	if model.LogFieldReferralUrls.IsUnknown() || model.LogFieldReferralUrls.IsNull() {
		model.LogFieldReferralUrls = types.StringValue("")
	}
	if model.LogFieldStartupid.IsUnknown() || model.LogFieldStartupid.IsNull() {
		model.LogFieldStartupid = types.StringValue("")
	}
	if model.LogFieldResponseOID.IsUnknown() || model.LogFieldResponseOID.IsNull() {
		model.LogFieldResponseOID = types.StringValue("")
	}
	if model.LogFieldCategory.IsUnknown() || model.LogFieldCategory.IsNull() {
		model.LogFieldCategory = types.StringValue("")
	}
	if model.LogFieldMessageID.IsUnknown() || model.LogFieldMessageID.IsNull() {
		model.LogFieldMessageID = types.StringValue("")
	}
	if model.LogFieldEntriesReturned.IsUnknown() || model.LogFieldEntriesReturned.IsNull() {
		model.LogFieldEntriesReturned = types.StringValue("")
	}
	if model.LogFieldSASLMechanismName.IsUnknown() || model.LogFieldSASLMechanismName.IsNull() {
		model.LogFieldSASLMechanismName = types.StringValue("")
	}
}

// Read a AccessLogFieldMappingResponse object into the model struct
func readAccessLogFieldMappingResponse(ctx context.Context, r *client.AccessLogFieldMappingResponse, state *logFieldMappingResourceModel, expectedValues *logFieldMappingResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFieldTimestamp = internaltypes.StringTypeOrNil(r.LogFieldTimestamp, internaltypes.IsEmptyString(expectedValues.LogFieldTimestamp))
	state.LogFieldConnectionID = internaltypes.StringTypeOrNil(r.LogFieldConnectionID, internaltypes.IsEmptyString(expectedValues.LogFieldConnectionID))
	state.LogFieldStartupid = internaltypes.StringTypeOrNil(r.LogFieldStartupid, internaltypes.IsEmptyString(expectedValues.LogFieldStartupid))
	state.LogFieldProductName = internaltypes.StringTypeOrNil(r.LogFieldProductName, internaltypes.IsEmptyString(expectedValues.LogFieldProductName))
	state.LogFieldInstanceName = internaltypes.StringTypeOrNil(r.LogFieldInstanceName, internaltypes.IsEmptyString(expectedValues.LogFieldInstanceName))
	state.LogFieldOperationID = internaltypes.StringTypeOrNil(r.LogFieldOperationID, internaltypes.IsEmptyString(expectedValues.LogFieldOperationID))
	state.LogFieldMessageType = internaltypes.StringTypeOrNil(r.LogFieldMessageType, internaltypes.IsEmptyString(expectedValues.LogFieldMessageType))
	state.LogFieldOperationType = internaltypes.StringTypeOrNil(r.LogFieldOperationType, internaltypes.IsEmptyString(expectedValues.LogFieldOperationType))
	state.LogFieldMessageID = internaltypes.StringTypeOrNil(r.LogFieldMessageID, internaltypes.IsEmptyString(expectedValues.LogFieldMessageID))
	state.LogFieldResultCode = internaltypes.StringTypeOrNil(r.LogFieldResultCode, internaltypes.IsEmptyString(expectedValues.LogFieldResultCode))
	state.LogFieldMessage = internaltypes.StringTypeOrNil(r.LogFieldMessage, internaltypes.IsEmptyString(expectedValues.LogFieldMessage))
	state.LogFieldOrigin = internaltypes.StringTypeOrNil(r.LogFieldOrigin, internaltypes.IsEmptyString(expectedValues.LogFieldOrigin))
	state.LogFieldRequesterDN = internaltypes.StringTypeOrNil(r.LogFieldRequesterDN, internaltypes.IsEmptyString(expectedValues.LogFieldRequesterDN))
	state.LogFieldDisconnectReason = internaltypes.StringTypeOrNil(r.LogFieldDisconnectReason, internaltypes.IsEmptyString(expectedValues.LogFieldDisconnectReason))
	state.LogFieldDeleteOldRDN = internaltypes.StringTypeOrNil(r.LogFieldDeleteOldRDN, internaltypes.IsEmptyString(expectedValues.LogFieldDeleteOldRDN))
	state.LogFieldAuthenticatedUserDN = internaltypes.StringTypeOrNil(r.LogFieldAuthenticatedUserDN, internaltypes.IsEmptyString(expectedValues.LogFieldAuthenticatedUserDN))
	state.LogFieldProcessingTime = internaltypes.StringTypeOrNil(r.LogFieldProcessingTime, internaltypes.IsEmptyString(expectedValues.LogFieldProcessingTime))
	state.LogFieldRequestedAttributes = internaltypes.StringTypeOrNil(r.LogFieldRequestedAttributes, internaltypes.IsEmptyString(expectedValues.LogFieldRequestedAttributes))
	state.LogFieldSASLMechanismName = internaltypes.StringTypeOrNil(r.LogFieldSASLMechanismName, internaltypes.IsEmptyString(expectedValues.LogFieldSASLMechanismName))
	state.LogFieldNewRDN = internaltypes.StringTypeOrNil(r.LogFieldNewRDN, internaltypes.IsEmptyString(expectedValues.LogFieldNewRDN))
	state.LogFieldBaseDN = internaltypes.StringTypeOrNil(r.LogFieldBaseDN, internaltypes.IsEmptyString(expectedValues.LogFieldBaseDN))
	state.LogFieldBindDN = internaltypes.StringTypeOrNil(r.LogFieldBindDN, internaltypes.IsEmptyString(expectedValues.LogFieldBindDN))
	state.LogFieldMatchedDN = internaltypes.StringTypeOrNil(r.LogFieldMatchedDN, internaltypes.IsEmptyString(expectedValues.LogFieldMatchedDN))
	state.LogFieldRequesterIPAddress = internaltypes.StringTypeOrNil(r.LogFieldRequesterIPAddress, internaltypes.IsEmptyString(expectedValues.LogFieldRequesterIPAddress))
	state.LogFieldAuthenticationType = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationType, internaltypes.IsEmptyString(expectedValues.LogFieldAuthenticationType))
	state.LogFieldNewSuperiorDN = internaltypes.StringTypeOrNil(r.LogFieldNewSuperiorDN, internaltypes.IsEmptyString(expectedValues.LogFieldNewSuperiorDN))
	state.LogFieldFilter = internaltypes.StringTypeOrNil(r.LogFieldFilter, internaltypes.IsEmptyString(expectedValues.LogFieldFilter))
	state.LogFieldAlternateAuthorizationDN = internaltypes.StringTypeOrNil(r.LogFieldAlternateAuthorizationDN, internaltypes.IsEmptyString(expectedValues.LogFieldAlternateAuthorizationDN))
	state.LogFieldEntryDN = internaltypes.StringTypeOrNil(r.LogFieldEntryDN, internaltypes.IsEmptyString(expectedValues.LogFieldEntryDN))
	state.LogFieldEntriesReturned = internaltypes.StringTypeOrNil(r.LogFieldEntriesReturned, internaltypes.IsEmptyString(expectedValues.LogFieldEntriesReturned))
	state.LogFieldAuthenticationFailureID = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationFailureID, internaltypes.IsEmptyString(expectedValues.LogFieldAuthenticationFailureID))
	state.LogFieldRequestOID = internaltypes.StringTypeOrNil(r.LogFieldRequestOID, internaltypes.IsEmptyString(expectedValues.LogFieldRequestOID))
	state.LogFieldResponseOID = internaltypes.StringTypeOrNil(r.LogFieldResponseOID, internaltypes.IsEmptyString(expectedValues.LogFieldResponseOID))
	state.LogFieldTargetProtocol = internaltypes.StringTypeOrNil(r.LogFieldTargetProtocol, internaltypes.IsEmptyString(expectedValues.LogFieldTargetProtocol))
	state.LogFieldTargetPort = internaltypes.StringTypeOrNil(r.LogFieldTargetPort, internaltypes.IsEmptyString(expectedValues.LogFieldTargetPort))
	state.LogFieldTargetAddress = internaltypes.StringTypeOrNil(r.LogFieldTargetAddress, internaltypes.IsEmptyString(expectedValues.LogFieldTargetAddress))
	state.LogFieldTargetAttribute = internaltypes.StringTypeOrNil(r.LogFieldTargetAttribute, internaltypes.IsEmptyString(expectedValues.LogFieldTargetAttribute))
	state.LogFieldTargetHost = internaltypes.StringTypeOrNil(r.LogFieldTargetHost, internaltypes.IsEmptyString(expectedValues.LogFieldTargetHost))
	state.LogFieldProtocolVersion = internaltypes.StringTypeOrNil(r.LogFieldProtocolVersion, internaltypes.IsEmptyString(expectedValues.LogFieldProtocolVersion))
	state.LogFieldProtocolName = internaltypes.StringTypeOrNil(r.LogFieldProtocolName, internaltypes.IsEmptyString(expectedValues.LogFieldProtocolName))
	state.LogFieldAuthenticationFailureReason = internaltypes.StringTypeOrNil(r.LogFieldAuthenticationFailureReason, internaltypes.IsEmptyString(expectedValues.LogFieldAuthenticationFailureReason))
	state.LogFieldAdditionalInformation = internaltypes.StringTypeOrNil(r.LogFieldAdditionalInformation, internaltypes.IsEmptyString(expectedValues.LogFieldAdditionalInformation))
	state.LogFieldUnindexed = internaltypes.StringTypeOrNil(r.LogFieldUnindexed, internaltypes.IsEmptyString(expectedValues.LogFieldUnindexed))
	state.LogFieldScope = internaltypes.StringTypeOrNil(r.LogFieldScope, internaltypes.IsEmptyString(expectedValues.LogFieldScope))
	state.LogFieldReferralUrls = internaltypes.StringTypeOrNil(r.LogFieldReferralUrls, internaltypes.IsEmptyString(expectedValues.LogFieldReferralUrls))
	state.LogFieldSourceAddress = internaltypes.StringTypeOrNil(r.LogFieldSourceAddress, internaltypes.IsEmptyString(expectedValues.LogFieldSourceAddress))
	state.LogFieldMessageIDToAbandon = internaltypes.StringTypeOrNil(r.LogFieldMessageIDToAbandon, internaltypes.IsEmptyString(expectedValues.LogFieldMessageIDToAbandon))
	state.LogFieldResponseControls = internaltypes.StringTypeOrNil(r.LogFieldResponseControls, internaltypes.IsEmptyString(expectedValues.LogFieldResponseControls))
	state.LogFieldRequestControls = internaltypes.StringTypeOrNil(r.LogFieldRequestControls, internaltypes.IsEmptyString(expectedValues.LogFieldRequestControls))
	state.LogFieldIntermediateClientResult = internaltypes.StringTypeOrNil(r.LogFieldIntermediateClientResult, internaltypes.IsEmptyString(expectedValues.LogFieldIntermediateClientResult))
	state.LogFieldIntermediateClientRequest = internaltypes.StringTypeOrNil(r.LogFieldIntermediateClientRequest, internaltypes.IsEmptyString(expectedValues.LogFieldIntermediateClientRequest))
	state.LogFieldReplicationChangeID = internaltypes.StringTypeOrNil(r.LogFieldReplicationChangeID, internaltypes.IsEmptyString(expectedValues.LogFieldReplicationChangeID))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a ErrorLogFieldMappingResponse object into the model struct
func readErrorLogFieldMappingResponse(ctx context.Context, r *client.ErrorLogFieldMappingResponse, state *logFieldMappingResourceModel, expectedValues *logFieldMappingResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFieldTimestamp = internaltypes.StringTypeOrNil(r.LogFieldTimestamp, internaltypes.IsEmptyString(expectedValues.LogFieldTimestamp))
	state.LogFieldProductName = internaltypes.StringTypeOrNil(r.LogFieldProductName, internaltypes.IsEmptyString(expectedValues.LogFieldProductName))
	state.LogFieldInstanceName = internaltypes.StringTypeOrNil(r.LogFieldInstanceName, internaltypes.IsEmptyString(expectedValues.LogFieldInstanceName))
	state.LogFieldStartupid = internaltypes.StringTypeOrNil(r.LogFieldStartupid, internaltypes.IsEmptyString(expectedValues.LogFieldStartupid))
	state.LogFieldCategory = internaltypes.StringTypeOrNil(r.LogFieldCategory, internaltypes.IsEmptyString(expectedValues.LogFieldCategory))
	state.LogFieldSeverity = internaltypes.StringTypeOrNil(r.LogFieldSeverity, internaltypes.IsEmptyString(expectedValues.LogFieldSeverity))
	state.LogFieldMessageID = internaltypes.StringTypeOrNil(r.LogFieldMessageID, internaltypes.IsEmptyString(expectedValues.LogFieldMessageID))
	state.LogFieldMessage = internaltypes.StringTypeOrNil(r.LogFieldMessage, internaltypes.IsEmptyString(expectedValues.LogFieldMessage))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLogFieldMappingOperations(plan logFieldMappingResourceModel, state logFieldMappingResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTimestamp, state.LogFieldTimestamp, "log-field-timestamp")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldConnectionID, state.LogFieldConnectionID, "log-field-connection-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldStartupid, state.LogFieldStartupid, "log-field-startupid")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldProductName, state.LogFieldProductName, "log-field-product-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldCategory, state.LogFieldCategory, "log-field-category")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldSeverity, state.LogFieldSeverity, "log-field-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldInstanceName, state.LogFieldInstanceName, "log-field-instance-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldOperationID, state.LogFieldOperationID, "log-field-operation-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessageType, state.LogFieldMessageType, "log-field-message-type")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldOperationType, state.LogFieldOperationType, "log-field-operation-type")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessageID, state.LogFieldMessageID, "log-field-message-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldResultCode, state.LogFieldResultCode, "log-field-result-code")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessage, state.LogFieldMessage, "log-field-message")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldOrigin, state.LogFieldOrigin, "log-field-origin")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldRequesterDN, state.LogFieldRequesterDN, "log-field-requester-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldDisconnectReason, state.LogFieldDisconnectReason, "log-field-disconnect-reason")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldDeleteOldRDN, state.LogFieldDeleteOldRDN, "log-field-delete-old-rdn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAuthenticatedUserDN, state.LogFieldAuthenticatedUserDN, "log-field-authenticated-user-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldProcessingTime, state.LogFieldProcessingTime, "log-field-processing-time")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldRequestedAttributes, state.LogFieldRequestedAttributes, "log-field-requested-attributes")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldSASLMechanismName, state.LogFieldSASLMechanismName, "log-field-sasl-mechanism-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldNewRDN, state.LogFieldNewRDN, "log-field-new-rdn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldBaseDN, state.LogFieldBaseDN, "log-field-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldBindDN, state.LogFieldBindDN, "log-field-bind-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMatchedDN, state.LogFieldMatchedDN, "log-field-matched-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldRequesterIPAddress, state.LogFieldRequesterIPAddress, "log-field-requester-ip-address")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAuthenticationType, state.LogFieldAuthenticationType, "log-field-authentication-type")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldNewSuperiorDN, state.LogFieldNewSuperiorDN, "log-field-new-superior-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldFilter, state.LogFieldFilter, "log-field-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAlternateAuthorizationDN, state.LogFieldAlternateAuthorizationDN, "log-field-alternate-authorization-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldEntryDN, state.LogFieldEntryDN, "log-field-entry-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldEntriesReturned, state.LogFieldEntriesReturned, "log-field-entries-returned")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAuthenticationFailureID, state.LogFieldAuthenticationFailureID, "log-field-authentication-failure-id")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldRequestOID, state.LogFieldRequestOID, "log-field-request-oid")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldResponseOID, state.LogFieldResponseOID, "log-field-response-oid")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTargetProtocol, state.LogFieldTargetProtocol, "log-field-target-protocol")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTargetPort, state.LogFieldTargetPort, "log-field-target-port")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTargetAddress, state.LogFieldTargetAddress, "log-field-target-address")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTargetAttribute, state.LogFieldTargetAttribute, "log-field-target-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldTargetHost, state.LogFieldTargetHost, "log-field-target-host")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldProtocolVersion, state.LogFieldProtocolVersion, "log-field-protocol-version")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldProtocolName, state.LogFieldProtocolName, "log-field-protocol-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAuthenticationFailureReason, state.LogFieldAuthenticationFailureReason, "log-field-authentication-failure-reason")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldAdditionalInformation, state.LogFieldAdditionalInformation, "log-field-additional-information")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldUnindexed, state.LogFieldUnindexed, "log-field-unindexed")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldScope, state.LogFieldScope, "log-field-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldReferralUrls, state.LogFieldReferralUrls, "log-field-referral-urls")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldSourceAddress, state.LogFieldSourceAddress, "log-field-source-address")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMessageIDToAbandon, state.LogFieldMessageIDToAbandon, "log-field-message-id-to-abandon")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldResponseControls, state.LogFieldResponseControls, "log-field-response-controls")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldRequestControls, state.LogFieldRequestControls, "log-field-request-controls")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldIntermediateClientResult, state.LogFieldIntermediateClientResult, "log-field-intermediate-client-result")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldIntermediateClientRequest, state.LogFieldIntermediateClientRequest, "log-field-intermediate-client-request")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldReplicationChangeID, state.LogFieldReplicationChangeID, "log-field-replication-change-id")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a access log-field-mapping
func (r *logFieldMappingResource) CreateAccessLogFieldMapping(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFieldMappingResourceModel) (*logFieldMappingResourceModel, error) {
	addRequest := client.NewAddAccessLogFieldMappingRequest([]client.EnumaccessLogFieldMappingSchemaUrn{client.ENUMACCESSLOGFIELDMAPPINGSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FIELD_MAPPINGACCESS},
		plan.Name.ValueString())
	addOptionalAccessLogFieldMappingFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFieldMappingAPI.AddLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFieldMappingRequest(
		client.AddAccessLogFieldMappingRequestAsAddLogFieldMappingRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFieldMappingAPI.AddLogFieldMappingExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Field Mapping", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFieldMappingResourceModel
	readAccessLogFieldMappingResponse(ctx, addResponse.AccessLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a error log-field-mapping
func (r *logFieldMappingResource) CreateErrorLogFieldMapping(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFieldMappingResourceModel) (*logFieldMappingResourceModel, error) {
	addRequest := client.NewAddErrorLogFieldMappingRequest([]client.EnumerrorLogFieldMappingSchemaUrn{client.ENUMERRORLOGFIELDMAPPINGSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FIELD_MAPPINGERROR},
		plan.Name.ValueString())
	addOptionalErrorLogFieldMappingFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFieldMappingAPI.AddLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFieldMappingRequest(
		client.AddErrorLogFieldMappingRequestAsAddLogFieldMappingRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFieldMappingAPI.AddLogFieldMappingExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Field Mapping", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFieldMappingResourceModel
	readErrorLogFieldMappingResponse(ctx, addResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logFieldMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logFieldMappingResourceModel
	var err error
	if plan.Type.ValueString() == "access" {
		state, err = r.CreateAccessLogFieldMapping(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "error" {
		state, err = r.CreateErrorLogFieldMapping(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

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
func (r *defaultLogFieldMappingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldMappingAPI.GetLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logFieldMappingResourceModel
	if readResponse.AccessLogFieldMappingResponse != nil {
		readAccessLogFieldMappingResponse(ctx, readResponse.AccessLogFieldMappingResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogFieldMappingResponse != nil {
		readErrorLogFieldMappingResponse(ctx, readResponse.ErrorLogFieldMappingResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFieldMappingAPI.UpdateLogFieldMapping(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogFieldMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFieldMappingAPI.UpdateLogFieldMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.AccessLogFieldMappingResponse != nil {
			readAccessLogFieldMappingResponse(ctx, updateResponse.AccessLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ErrorLogFieldMappingResponse != nil {
			readErrorLogFieldMappingResponse(ctx, updateResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *logFieldMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogFieldMappingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogFieldMapping(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logFieldMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFieldMappingAPI.GetLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Mapping", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Mapping", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.AccessLogFieldMappingResponse != nil {
		readAccessLogFieldMappingResponse(ctx, readResponse.AccessLogFieldMappingResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogFieldMappingResponse != nil {
		readErrorLogFieldMappingResponse(ctx, readResponse.ErrorLogFieldMappingResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logFieldMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogFieldMappingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFieldMapping(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogFieldMapping(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logFieldMappingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logFieldMappingResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFieldMappingAPI.UpdateLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogFieldMappingOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFieldMappingAPI.UpdateLogFieldMappingExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Mapping", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.AccessLogFieldMappingResponse != nil {
			readAccessLogFieldMappingResponse(ctx, updateResponse.AccessLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ErrorLogFieldMappingResponse != nil {
			readErrorLogFieldMappingResponse(ctx, updateResponse.ErrorLogFieldMappingResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultLogFieldMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logFieldMappingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logFieldMappingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFieldMappingAPI.DeleteLogFieldMappingExecute(r.apiClient.LogFieldMappingAPI.DeleteLogFieldMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log Field Mapping", err, httpResp)
		return
	}
}

func (r *logFieldMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFieldMapping(ctx, req, resp)
}

func (r *defaultLogFieldMappingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFieldMapping(ctx, req, resp)
}

func importLogFieldMapping(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
