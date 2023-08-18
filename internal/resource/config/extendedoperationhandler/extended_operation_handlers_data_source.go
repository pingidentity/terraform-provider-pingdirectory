package extendedoperationhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &extendedOperationHandlersDataSource{}
	_ datasource.DataSourceWithConfigure = &extendedOperationHandlersDataSource{}
)

// Create a Extended Operation Handlers data source
func NewExtendedOperationHandlersDataSource() datasource.DataSource {
	return &extendedOperationHandlersDataSource{}
}

// extendedOperationHandlersDataSource is the datasource implementation.
type extendedOperationHandlersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *extendedOperationHandlersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extended_operation_handlers"
}

// Configure adds the provider configured client to the data source.
func (r *extendedOperationHandlersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type extendedOperationHandlersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *extendedOperationHandlersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Extended Operation Handler objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Extended Operation Handler objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *extendedOperationHandlersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state extendedOperationHandlersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.ExtendedOperationHandlerApi.ListExtendedOperationHandlers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.ListExtendedOperationHandlersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Extended Operation Handler objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.CancelExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.CancelExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("cancel")
		}
		if response.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ValidateTotpPasswordExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("validate-totp-password")
		}
		if response.InteractiveTransactionsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.InteractiveTransactionsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("interactive-transactions")
		}
		if response.ReplaceCertificateExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ReplaceCertificateExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("replace-certificate")
		}
		if response.BackupCompatibilityExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.BackupCompatibilityExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("backup-compatibility")
		}
		if response.GetConnectionIdExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetConnectionIdExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-connection-id")
		}
		if response.BatchedTransactionsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.BatchedTransactionsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("batched-transactions")
		}
		if response.AdministrativeSessionExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.AdministrativeSessionExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("administrative-session")
		}
		if response.SingleUseTokensExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SingleUseTokensExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("single-use-tokens")
		}
		if response.StreamDirectoryValuesExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.StreamDirectoryValuesExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("stream-directory-values")
		}
		if response.StartTlsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.StartTlsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("start-tls")
		}
		if response.SubtreeAccessibilityExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SubtreeAccessibilityExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("subtree-accessibility")
		}
		if response.GetSymmetricKeyExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetSymmetricKeyExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-symmetric-key")
		}
		if response.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetPasswordQualityRequirementsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-password-quality-requirements")
		}
		if response.DeliverOtpExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.DeliverOtpExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("deliver-otp")
		}
		if response.ThirdPartyExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		if response.ThirdPartyProxiedExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyProxiedExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("third-party-proxied")
		}
		if response.SynchronizeEncryptionSettingsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.SynchronizeEncryptionSettingsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("synchronize-encryption-settings")
		}
		if response.MultiUpdateExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.MultiUpdateExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("multi-update")
		}
		if response.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.NotificationSubscriptionExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("notification-subscription")
		}
		if response.PasswordModifyExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.PasswordModifyExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("password-modify")
		}
		if response.CustomExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.CustomExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.CollectSupportDataExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.CollectSupportDataExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("collect-support-data")
		}
		if response.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.ExportReversiblePasswordsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("export-reversible-passwords")
		}
		if response.GetChangelogBatchExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetChangelogBatchExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-changelog-batch")
		}
		if response.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-supported-otp-delivery-mechanisms")
		}
		if response.GeneratePasswordExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GeneratePasswordExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("generate-password")
		}
		if response.WhoAmIExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.WhoAmIExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("who-am-i")
		}
		if response.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.DeliverPasswordResetTokenExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("deliver-password-reset-token")
		}
		if response.StreamProxyValuesExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.StreamProxyValuesExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("stream-proxy-values")
		}
		if response.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.PasswordPolicyStateExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("password-policy-state")
		}
		if response.GetConfigurationExtendedOperationHandlerResponse != nil {
			attributes["id"] = types.StringValue(response.GetConfigurationExtendedOperationHandlerResponse.Id)
			attributes["type"] = types.StringValue("get-configuration")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
