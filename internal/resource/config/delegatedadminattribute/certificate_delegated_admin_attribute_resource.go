package delegatedadminattribute

import (
	"context"
	"strings"
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
	_ resource.Resource                = &certificateDelegatedAdminAttributeResource{}
	_ resource.ResourceWithConfigure   = &certificateDelegatedAdminAttributeResource{}
	_ resource.ResourceWithImportState = &certificateDelegatedAdminAttributeResource{}
	_ resource.Resource                = &defaultCertificateDelegatedAdminAttributeResource{}
	_ resource.ResourceWithConfigure   = &defaultCertificateDelegatedAdminAttributeResource{}
	_ resource.ResourceWithImportState = &defaultCertificateDelegatedAdminAttributeResource{}
)

// Create a Certificate Delegated Admin Attribute resource
func NewCertificateDelegatedAdminAttributeResource() resource.Resource {
	return &certificateDelegatedAdminAttributeResource{}
}

func NewDefaultCertificateDelegatedAdminAttributeResource() resource.Resource {
	return &defaultCertificateDelegatedAdminAttributeResource{}
}

// certificateDelegatedAdminAttributeResource is the resource implementation.
type certificateDelegatedAdminAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCertificateDelegatedAdminAttributeResource is the resource implementation.
type defaultCertificateDelegatedAdminAttributeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *certificateDelegatedAdminAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificate_delegated_admin_attribute"
}

func (r *defaultCertificateDelegatedAdminAttributeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_certificate_delegated_admin_attribute"
}

// Configure adds the provider configured client to the resource.
func (r *certificateDelegatedAdminAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultCertificateDelegatedAdminAttributeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type certificateDelegatedAdminAttributeResourceModel struct {
	Id                    types.String `tfsdk:"id"`
	LastUpdated           types.String `tfsdk:"last_updated"`
	Notifications         types.Set    `tfsdk:"notifications"`
	RequiredActions       types.Set    `tfsdk:"required_actions"`
	RestResourceTypeName  types.String `tfsdk:"rest_resource_type_name"`
	AllowedMIMEType       types.Set    `tfsdk:"allowed_mime_type"`
	Description           types.String `tfsdk:"description"`
	AttributeType         types.String `tfsdk:"attribute_type"`
	DisplayName           types.String `tfsdk:"display_name"`
	Mutability            types.String `tfsdk:"mutability"`
	MultiValued           types.Bool   `tfsdk:"multi_valued"`
	AttributeCategory     types.String `tfsdk:"attribute_category"`
	DisplayOrderIndex     types.Int64  `tfsdk:"display_order_index"`
	ReferenceResourceType types.String `tfsdk:"reference_resource_type"`
	AttributePresentation types.String `tfsdk:"attribute_presentation"`
	DateTimeFormat        types.String `tfsdk:"date_time_format"`
}

// GetSchema defines the schema for the resource.
func (r *certificateDelegatedAdminAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificateDelegatedAdminAttributeSchema(ctx, req, resp, false)
}

func (r *defaultCertificateDelegatedAdminAttributeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	certificateDelegatedAdminAttributeSchema(ctx, req, resp, true)
}

func certificateDelegatedAdminAttributeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Certificate Delegated Admin Attribute.",
		Attributes: map[string]schema.Attribute{
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"allowed_mime_type": schema.SetAttribute{
				Description: "The list of file types allowed to be uploaded. If no types are specified, then all types will be allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Attribute",
				Optional:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Attribute.",
				Required:    true,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the attribute can be written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Indicates whether this Delegated Admin Attribute may have multiple values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"attribute_category": schema.StringAttribute{
				Description: "Specifies which attribute category this attribute belongs to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"display_order_index": schema.Int64Attribute{
				Description: "This property determines a display order for attributes within a given attribute category. Attributes are ordered within their category based on this index from least to greatest.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"reference_resource_type": schema.StringAttribute{
				Description: "For LDAP attributes with DN syntax, specifies what kind of resource is referenced.",
				Optional:    true,
			},
			"attribute_presentation": schema.StringAttribute{
				Description: "Indicates how the attribute is presented to the user of the app.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"date_time_format": schema.StringAttribute{
				Description: "Specifies the format string that is used to present a date and/or time value to the user of the app. This property only applies to LDAP attribute types whose LDAP syntax is GeneralizedTime and is ignored if the attribute type has any other syntax.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"attribute_type", "rest_resource_type_name"})
	}
	config.AddCommonSchema(&schema, false)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalCertificateDelegatedAdminAttributeFields(ctx context.Context, addRequest *client.AddCertificateDelegatedAdminAttributeRequest, plan certificateDelegatedAdminAttributeResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedMIMEType) {
		var slice []string
		plan.AllowedMIMEType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumdelegatedAdminAttributeCertificateAllowedMIMETypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumdelegatedAdminAttributeCertificateAllowedMIMETypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedMIMEType = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Mutability) {
		mutability, err := client.NewEnumdelegatedAdminAttributeMutabilityPropFromValue(plan.Mutability.ValueString())
		if err != nil {
			return err
		}
		addRequest.Mutability = mutability
	}
	if internaltypes.IsDefined(plan.MultiValued) {
		boolVal := plan.MultiValued.ValueBool()
		addRequest.MultiValued = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributeCategory) {
		stringVal := plan.AttributeCategory.ValueString()
		addRequest.AttributeCategory = &stringVal
	}
	if internaltypes.IsDefined(plan.DisplayOrderIndex) {
		intVal := int32(plan.DisplayOrderIndex.ValueInt64())
		addRequest.DisplayOrderIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReferenceResourceType) {
		stringVal := plan.ReferenceResourceType.ValueString()
		addRequest.ReferenceResourceType = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AttributePresentation) {
		attributePresentation, err := client.NewEnumdelegatedAdminAttributeAttributePresentationPropFromValue(plan.AttributePresentation.ValueString())
		if err != nil {
			return err
		}
		addRequest.AttributePresentation = attributePresentation
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DateTimeFormat) {
		stringVal := plan.DateTimeFormat.ValueString()
		addRequest.DateTimeFormat = &stringVal
	}
	return nil
}

// Read a CertificateDelegatedAdminAttributeResponse object into the model struct
func readCertificateDelegatedAdminAttributeResponse(ctx context.Context, r *client.CertificateDelegatedAdminAttributeResponse, state *certificateDelegatedAdminAttributeResourceModel, expectedValues *certificateDelegatedAdminAttributeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.RestResourceTypeName = expectedValues.RestResourceTypeName
	state.AllowedMIMEType = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminAttributeCertificateAllowedMIMETypeProp(r.AllowedMIMEType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, internaltypes.IsEmptyString(expectedValues.AttributeCategory))
	state.DisplayOrderIndex = types.Int64Value(int64(r.DisplayOrderIndex))
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, internaltypes.IsEmptyString(expectedValues.ReferenceResourceType))
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), internaltypes.IsEmptyString(expectedValues.AttributePresentation))
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, internaltypes.IsEmptyString(expectedValues.DateTimeFormat))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCertificateDelegatedAdminAttributeOperations(plan certificateDelegatedAdminAttributeResourceModel, state certificateDelegatedAdminAttributeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedMIMEType, state.AllowedMIMEType, "allowed-mime-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.DisplayName, state.DisplayName, "display-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Mutability, state.Mutability, "mutability")
	operations.AddBoolOperationIfNecessary(&ops, plan.MultiValued, state.MultiValued, "multi-valued")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributeCategory, state.AttributeCategory, "attribute-category")
	operations.AddInt64OperationIfNecessary(&ops, plan.DisplayOrderIndex, state.DisplayOrderIndex, "display-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.ReferenceResourceType, state.ReferenceResourceType, "reference-resource-type")
	operations.AddStringOperationIfNecessary(&ops, plan.AttributePresentation, state.AttributePresentation, "attribute-presentation")
	operations.AddStringOperationIfNecessary(&ops, plan.DateTimeFormat, state.DateTimeFormat, "date-time-format")
	return ops
}

// Create a new resource
func (r *certificateDelegatedAdminAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificateDelegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddCertificateDelegatedAdminAttributeRequest(plan.AttributeType.ValueString(),
		[]client.EnumcertificateDelegatedAdminAttributeSchemaUrn{client.ENUMCERTIFICATEDELEGATEDADMINATTRIBUTESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0DELEGATED_ADMIN_ATTRIBUTECERTIFICATE},
		plan.DisplayName.ValueString())
	err := addOptionalCertificateDelegatedAdminAttributeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Certificate Delegated Admin Attribute", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.RestResourceTypeName.ValueString())
	apiAddRequest = apiAddRequest.AddDelegatedAdminAttributeRequest(
		client.AddCertificateDelegatedAdminAttributeRequestAsAddDelegatedAdminAttributeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.AddDelegatedAdminAttributeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Certificate Delegated Admin Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state certificateDelegatedAdminAttributeResourceModel
	readCertificateDelegatedAdminAttributeResponse(ctx, addResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultCertificateDelegatedAdminAttributeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan certificateDelegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Delegated Admin Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state certificateDelegatedAdminAttributeResourceModel
	readCertificateDelegatedAdminAttributeResponse(ctx, readResponse.CertificateDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttribute(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString())
	ops := createCertificateDelegatedAdminAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Certificate Delegated Admin Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCertificateDelegatedAdminAttributeResponse(ctx, updateResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *certificateDelegatedAdminAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCertificateDelegatedAdminAttributeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCertificateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCertificateDelegatedAdminAttribute(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state certificateDelegatedAdminAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.DelegatedAdminAttributeApi.GetDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.AttributeType.ValueString(), state.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Certificate Delegated Admin Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCertificateDelegatedAdminAttributeResponse(ctx, readResponse.CertificateDelegatedAdminAttributeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *certificateDelegatedAdminAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCertificateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCertificateDelegatedAdminAttributeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCertificateDelegatedAdminAttribute(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCertificateDelegatedAdminAttribute(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan certificateDelegatedAdminAttributeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state certificateDelegatedAdminAttributeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.AttributeType.ValueString(), plan.RestResourceTypeName.ValueString())

	// Determine what update operations are necessary
	ops := createCertificateDelegatedAdminAttributeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.DelegatedAdminAttributeApi.UpdateDelegatedAdminAttributeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Certificate Delegated Admin Attribute", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCertificateDelegatedAdminAttributeResponse(ctx, updateResponse.CertificateDelegatedAdminAttributeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCertificateDelegatedAdminAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *certificateDelegatedAdminAttributeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state certificateDelegatedAdminAttributeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.DelegatedAdminAttributeApi.DeleteDelegatedAdminAttributeExecute(r.apiClient.DelegatedAdminAttributeApi.DeleteDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AttributeType.ValueString(), state.RestResourceTypeName.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Certificate Delegated Admin Attribute", err, httpResp)
		return
	}
}

func (r *certificateDelegatedAdminAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCertificateDelegatedAdminAttribute(ctx, req, resp)
}

func (r *defaultCertificateDelegatedAdminAttributeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCertificateDelegatedAdminAttribute(ctx, req, resp)
}

func importCertificateDelegatedAdminAttribute(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [rest-resource-type-name]/[delegated-admin-attribute-attribute-type]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("rest_resource_type_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("attribute_type"), split[1])...)
}
