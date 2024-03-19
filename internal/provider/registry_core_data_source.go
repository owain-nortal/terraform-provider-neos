package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
)

func NewRegistryCoreDataSource() datasource.DataSource {
	return &registryCoreDataSource{}
}

var (
	_ datasource.DataSource              = &registryCoreDataSource{}
	_ datasource.DataSourceWithConfigure = &registryCoreDataSource{}
)

type registryCoreDataSource struct {
	client *neos.RegistryCoreClient
}

func (d *registryCoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registry_core"
}

func (d *registryCoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	mod := RegistryCoreDataSourceModel{}
	req.Config.Get(ctx, &mod)

	var state RegistryCoreDataSourceModel

	//state.RegistryCores = append(state.RegistryCores, RegistryCoreModel{})
	state.RegistryCores = []RegistryCoreModel{}

	//state.RegistryCores[0].ID = types.StringValue("00000000-0000-0000-0000-000000000000")

	list, err := d.client.Get(mod.Account.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read registry core List",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, ds := range list.Cores {
		id := d.splitUrnGetId(ds.Urn)
		registryCoreState := RegistryCoreModel{
			ID:      types.StringValue(id),
			Host:    types.StringValue(ds.Host),
			Name:    types.StringValue(ds.Name),
			Urn:     types.StringValue(ds.Urn),
			Account: types.StringValue(ds.Account),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Urn))
		if (mod.Name.ValueString() == ds.Name) || mod.Name.IsNull() {
			state.RegistryCores = append(state.RegistryCores, registryCoreState)
		}
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d registryCoreDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data RegistryCoreDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

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

	var client *neos.NeosClient
	var ok bool
	client, ok = req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected registryCoreDataSource Configure Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))

		return
	}

	d.client = &client.RegistryCoreClient
}

func (d *registryCoreDataSource) Schema(_ context.Context, foo datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account": schema.StringAttribute{
				Required: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
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
						"id": schema.StringAttribute{
							Computed: true,
						},
						"account": schema.StringAttribute{
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
	Account       types.String        `tfsdk:"account"`
	Name          types.String        `tfsdk:"name"`
}

type RegistryCoreModel struct {
	ID      types.String `tfsdk:"id"`
	Host    types.String `tfsdk:"host"`
	Urn     types.String `tfsdk:"urn"`
	Name    types.String `tfsdk:"name"`
	Account types.String `tfsdk:"account"`
}

func (d *registryCoreDataSource) splitUrnGetId(urn string) string {
	bits := strings.Split(urn, ":")
	return bits[len(bits)-1]
}
