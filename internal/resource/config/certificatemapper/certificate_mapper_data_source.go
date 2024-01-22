package certificatemapper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificateMapperDataSource{}
	_ datasource.DataSourceWithConfigure = &certificateMapperDataSource{}
)

// Create a Certificate Mapper data source
func NewCertificateMapperDataSource() datasource.DataSource {
	return &certificateMapperDataSource{}
}

// certificateMapperDataSource is the datasource implementation.
type certificateMapperDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *certificateMapperDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_mapper"
}

// Configure adds the provider configured client to the data source.
func (r *certificateMapperDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificateMapperDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	Type                    types.String `tfsdk:"type"`
	ExtensionClass          types.String `tfsdk:"extension_class"`
	ExtensionArgument       types.Set    `tfsdk:"extension_argument"`
	FingerprintAttribute    types.String `tfsdk:"fingerprint_attribute"`
	FingerprintAlgorithm    types.String `tfsdk:"fingerprint_algorithm"`
	SubjectAttributeMapping types.Set    `tfsdk:"subject_attribute_mapping"`
	ScriptClass             types.String `tfsdk:"script_class"`
	ScriptArgument          types.Set    `tfsdk:"script_argument"`
	SubjectAttribute        types.String `tfsdk:"subject_attribute"`
	UserBaseDN              types.Set    `tfsdk:"user_base_dn"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *certificateMapperDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Certificate Mapper.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Certificate Mapper resource. Options are ['subject-equals-dn', 'subject-dn-to-user-attribute', 'groovy-scripted', 'subject-attribute-to-user-attribute', 'fingerprint', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Certificate Mapper.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"fingerprint_attribute": schema.StringAttribute{
				Description: "Specifies the attribute in which to look for the fingerprint.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"fingerprint_algorithm": schema.StringAttribute{
				Description: "Specifies the name of the digest algorithm to compute the fingerprint of client certificates.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_attribute_mapping": schema.SetAttribute{
				Description: "Specifies a mapping between certificate attributes and user attributes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Certificate Mapper.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Certificate Mapper. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"subject_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute whose value should exactly match the certificate subject DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"user_base_dn": schema.SetAttribute{
				Description:         "When the `type` attribute is set to  one of [`subject-dn-to-user-attribute`, `subject-attribute-to-user-attribute`]: Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry. When the `type` attribute is set to `fingerprint`: Specifies the set of base DNs below which to search for users.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`subject-dn-to-user-attribute`, `subject-attribute-to-user-attribute`]: Specifies the base DNs that should be used when performing searches to map the client certificate to a user entry.\n  - `fingerprint`: Specifies the set of base DNs below which to search for users.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Certificate Mapper",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Certificate Mapper is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SubjectEqualsDnCertificateMapperResponse object into the model struct
func readSubjectEqualsDnCertificateMapperResponseDataSource(ctx context.Context, r *client.SubjectEqualsDnCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-equals-dn")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SubjectDnToUserAttributeCertificateMapperResponse object into the model struct
func readSubjectDnToUserAttributeCertificateMapperResponseDataSource(ctx context.Context, r *client.SubjectDnToUserAttributeCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-dn-to-user-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubjectAttribute = types.StringValue(r.SubjectAttribute)
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedCertificateMapperResponse object into the model struct
func readGroovyScriptedCertificateMapperResponseDataSource(ctx context.Context, r *client.GroovyScriptedCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SubjectAttributeToUserAttributeCertificateMapperResponse object into the model struct
func readSubjectAttributeToUserAttributeCertificateMapperResponseDataSource(ctx context.Context, r *client.SubjectAttributeToUserAttributeCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("subject-attribute-to-user-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SubjectAttributeMapping = internaltypes.GetStringSet(r.SubjectAttributeMapping)
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a FingerprintCertificateMapperResponse object into the model struct
func readFingerprintCertificateMapperResponseDataSource(ctx context.Context, r *client.FingerprintCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fingerprint")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FingerprintAttribute = types.StringValue(r.FingerprintAttribute)
	state.FingerprintAlgorithm = types.StringValue(r.FingerprintAlgorithm.String())
	state.UserBaseDN = internaltypes.GetStringSet(r.UserBaseDN)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyCertificateMapperResponse object into the model struct
func readThirdPartyCertificateMapperResponseDataSource(ctx context.Context, r *client.ThirdPartyCertificateMapperResponse, state *certificateMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *certificateMapperDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state certificateMapperDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CertificateMapperAPI.GetCertificateMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SubjectEqualsDnCertificateMapperResponse != nil {
		readSubjectEqualsDnCertificateMapperResponseDataSource(ctx, readResponse.SubjectEqualsDnCertificateMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectDnToUserAttributeCertificateMapperResponse != nil {
		readSubjectDnToUserAttributeCertificateMapperResponseDataSource(ctx, readResponse.SubjectDnToUserAttributeCertificateMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedCertificateMapperResponse != nil {
		readGroovyScriptedCertificateMapperResponseDataSource(ctx, readResponse.GroovyScriptedCertificateMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse != nil {
		readSubjectAttributeToUserAttributeCertificateMapperResponseDataSource(ctx, readResponse.SubjectAttributeToUserAttributeCertificateMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FingerprintCertificateMapperResponse != nil {
		readFingerprintCertificateMapperResponseDataSource(ctx, readResponse.FingerprintCertificateMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyCertificateMapperResponse != nil {
		readThirdPartyCertificateMapperResponseDataSource(ctx, readResponse.ThirdPartyCertificateMapperResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
