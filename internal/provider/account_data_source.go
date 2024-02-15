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

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

var (
	_ datasource.DataSource              = &accountDataSource{}
	_ datasource.DataSourceWithConfigure = &accountDataSource{}
)

type accountDataSource struct {
	client *neos.DataSourceClient
}

func (d *accountDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_source"
}

func (d *accountDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "accountDataSource READ")

	var state AccountDataSourceModel

	list, err := d.client.Get()
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Data System List", err.Error())
		return
	}

	 // Map response body to model
	for _, ds := range list.Entities {
		dataSourceState := DataSourceModelV2{
			Identifier:  types.StringValue(ds.Identifier),
			Name:        types.StringValue(ds.Name),
			Description: types.StringValue(ds.Description),
			Label:       types.StringValue(ds.Label),
			Owner:       types.StringValue(ds.Owner),
			Urn:         types.StringValue(ds.Urn),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Identifier))
		state.DataSources = append(state.DataSources, dataSourceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *accountDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.DataSourceClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *neos.DataSourceClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"datasystems": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"urn": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"label": schema.StringAttribute{
							Computed: true,
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"owner": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"state": schema.StringAttribute{
									Computed: true,
								},
								"healthy": schema.BoolAttribute{
									Computed: true,
								},
							},
						},
					},
				},
			},
		},
	}
}

type AccountDataSourceModel struct {
	DataSources []DataSourceModelV2 `tfsdk:"datasource"`
}

// coffeesModel maps coffees schema data.
type AccountModel struct {
	Identifier  types.String           `tfsdk:"id"`
	Urn         types.String           `tfsdk:"urn"`
	Name        types.String           `tfsdk:"name"`
	Description types.String           `tfsdk:"description"`
	Label       types.String           `tfsdk:"label"`
	Owner       types.String           `tfsdk:"owner"`
	CreatedAt   types.String           `tfsdk:"created_at"`
	State       DataSourceStateModelV2 `tfsdk:"state"`
}

type AccountStateModel struct {
	State   types.String `tfsdk:"state"`
	Healthy types.Bool   `tfsdk:"healthy"`
}
