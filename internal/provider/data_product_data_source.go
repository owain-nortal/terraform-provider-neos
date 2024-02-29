package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
)

func NewDataProductDataSource() datasource.DataSource {
	return &dataProductDataSource{}
}

var (
	_ datasource.DataSource              = &dataProductDataSource{}
	_ datasource.DataSourceWithConfigure = &dataProductDataSource{}
)

type dataProductDataSource struct {
	client *neos.DataProductClient
}

func (d *dataProductDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_product"
	//resp.TypeName = "neos_data_product"
}

func (d *dataProductDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "Abi READ")

	var state DataProductDataSourceModelV2

	list, err := d.client.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data Product List",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Abi READ Post error ")

	tflog.Info(ctx, fmt.Sprintf("Abi READ length %d", len(list.Entities)))

	// Map response body to model
	for _, ds := range list.Entities {
		dataProductStateV2 := DataProductModelV2{
			Identifier:  types.StringValue(ds.Identifier),
			Name:        types.StringValue(ds.Name),
			Description: types.StringValue(ds.Description),
			Label:       types.StringValue(ds.Label),
			Owner:       types.StringValue(ds.Owner),
			Urn:         types.StringValue(ds.Urn),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Identifier))
		state.DataProducts = append(state.DataProducts, dataProductStateV2)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *dataProductDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected dataProductDataSource Configure Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))

		return
	}

	d.client = &client.DataProductClient
}

func (d *dataProductDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

type DataProductDataSourceModelV2 struct {
	DataProducts []DataProductModelV2 `tfsdk:"datasystems"`
}

// coffeesModel maps coffees schema data.
type DataProductModelV2 struct {
	Identifier  types.String            `tfsdk:"id"`
	Urn         types.String            `tfsdk:"urn"`
	Name        types.String            `tfsdk:"name"`
	Description types.String            `tfsdk:"description"`
	Label       types.String            `tfsdk:"label"`
	Owner       types.String            `tfsdk:"owner"`
	CreatedAt   types.String            `tfsdk:"created_at"`
	State       DataProductStateModelV2 `tfsdk:"state"`
}

type DataProductStateModelV2 struct {
	State   types.String `tfsdk:"state"`
	Healthy types.Bool   `tfsdk:"healthy"`
}
