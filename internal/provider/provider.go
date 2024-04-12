package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &neosProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &neosProvider{
			version: version,
		}
	}
}

type neosProviderModel struct {
	HubHost   types.String `tfsdk:"hub_host"`
	CoreHost  types.String `tfsdk:"core_host"`
	Account   types.String `tfsdk:"account"`
	Partition types.String `tfsdk:"partition"`

	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// neosProvider is the provider implementation.
type neosProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *neosProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "neos"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *neosProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hub_host": schema.StringAttribute{
				Optional: true,
			},
			"core_host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"account": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"partition": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a NEOS API client for data sources and resources.
func (p *neosProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	tflog.Info(ctx, "Configuring NEOS client")

	// Retrieve provider data from configuration
	var config neosProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.HubHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hub_host"),
			"Unknown NEOS HUB Host",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS HUB host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.Account.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("account"),
			"Unknown NEOS Account",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS Account "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.Partition.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("partition"),
			"Unknown NEOS partition",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS partition "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.CoreHost.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("core_host"),
			"Unknown NEOS Core Host",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS Core host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown NEOS API Username",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown NEOS API Password",
			"The provider cannot create the NEOS API client as there is an unknown configuration value for the NEOS API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	hubhost := os.Getenv("NEOS_HUB_HOST")
	corehost := os.Getenv("NEOS_CORE_HOST")
	username := os.Getenv("NEOS_USERNAME")
	password := os.Getenv("NEOS_PASSWORD")
	account := os.Getenv("NEOS_ACCOUNT")
	partition := os.Getenv("NEOS_PARTITION")

	if !config.HubHost.IsNull() {
		hubhost = config.HubHost.ValueString()
	}

	if !config.CoreHost.IsNull() {
		corehost = config.CoreHost.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Account.IsNull() {
		account = config.Account.ValueString()
	}

	if !config.Partition.IsNull() {
		partition = config.Partition.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if hubhost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("hub_host"),
			"Missing NEOS API hub host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS hub host. "+
				"Set the host value in the configuration or use the NEOS_HUB_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if corehost == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("core_host"),
			"Missing NEOS API core host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS core host. "+
				"Set the host value in the configuration or use the NEOS_CORE_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if account == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("account_host"),
			"Missing NEOS account host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS account. "+
				"Set the host value in the configuration or use the NEOS_ACCOUNT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if partition == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("partition_host"),
			"Missing NEOS partition host",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS partition. "+
				"Set the host value in the configuration or use the NEOS_PARTITION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing NEOS API Username",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS username. "+
				"Set the username value in the configuration or use the NEOS_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing NEOS API Password",
			"The provider cannot create the NEOS API client as there is a missing or empty value for the NEOS password. "+
				"Set the password value in the configuration or use the NEOS_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "neos_hub", hubhost)
	ctx = tflog.SetField(ctx, "neos_account", account)
	ctx = tflog.SetField(ctx, "neos_partition", partition)
	ctx = tflog.SetField(ctx, "neos_username", username)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "neos_password")
	tflog.Info(ctx, "Creating NEOS IAM client")

	iamUrl := fmt.Sprintf("%s/api/iam", hubhost)
	iamClient := neos.NewIAMClient(iamUrl, username, password)
	loginResponse, err := iamClient.Login()

	tflog.Info(ctx, loginResponse.AccessToken)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create NEOS API Client",
			"An unexpected error occurred when creating the NEOS API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"NEOS Client Error: "+err.Error(),
		)
		return
	}

	neos.AccessToken = loginResponse.AccessToken

	client, err := neos.NewNeosClient(hubhost, corehost, "https", account, partition)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create NewNeosClient",
			"An unexpected error occurred when creating the NewNeosClient. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"NEOS Client Error: "+err.Error(),
		)
		return
	}

	// Make the NEOS client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = &client
	resp.ResourceData = &client
}

// DataSources defines the data sources implemented in the provider.
func (p *neosProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewDataSystemDataSource,
		NewDataProductDataSource,
		NewDataSourceDataSource,
		NewDataUnitDataSource,
		NewGroupDataSource,
		NewLinksDataSource,
		NewRegistryCoreDataSource,
		NewUserDataSource,
		NewUserPolicyDataSource,
	}
}

func (p *neosProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAccountResource,
		NewDataProductResource,
		NewDataProductBuilderResource,
		NewDataSourceResource,
		NewDataSystemResource,
		NewDataUnitResource,
		NewGroupResource,
		NewLinkDataSourceDataUnitResource,
		NewLinkDataSystemDataSourceResource,
		NewLinkDataUnitDataProductResource,
		NewLinkDataProductOutputResource,
		NewLinkDataProductDataProductResource,
		NewOutputResource,
		NewRegistryCoreResource,
		NewSecretResource,
		NewUserResource,
		NewUserPolicyResource,
	}
}
