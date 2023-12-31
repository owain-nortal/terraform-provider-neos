package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/owain-nortal/neos-client-go"
)

func NewRegistryCoreDataSource() datasource.DataSource {
	return &registryCoreDataSource{}
}

var (
	_ datasource.DataSource              = &registryCoreDataSource{}
	_ datasource.DataSourceWithConfigure = &registryCoreDataSource{}
)

type registryCoreDataSource struct {
	client *neos.NeosClient
}

func (d *registryCoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_core"
}

func (d *registryCoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state RegistryCoreDataSourceModel

	list, err := d.client.RegistryCoreGet()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read registry core List",
			err.Error(),
		)
		return
	}

	// tflog.Info(ctx, "Abi READ Post error ")

	// tflog.Info(ctx, fmt.Sprintf("Abi READ length %d", len(list.Entities)))

	// Map response body to model
	for _, ds := range list.Cores {
		registryCoreState := RegistryCoreModel{
			Host: types.StringValue(ds.Host),
			Name: types.StringValue(ds.Name),
			Urn:  types.StringValue(ds.Urn),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Urn))
		state.RegistryCores = append(state.RegistryCores, registryCoreState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *registryCoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Registry Core Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *registryCoreDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"registry_cores": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"urn": schema.StringAttribute{
							Computed: true,
						},
						"host": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type RegistryCoreDataSourceModel struct {
	RegistryCores []RegistryCoreModel `tfsdk:"registry_cores"`
}

type RegistryCoreModel struct {
	Host types.String `tfsdk:"host"`
	Urn  types.String `tfsdk:"urn"`
	Name types.String `tfsdk:"name"`
}
