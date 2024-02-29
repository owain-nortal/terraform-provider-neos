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

func NewDataSystemDataSource() datasource.DataSource {
	return &dataSystemDataSourceV2{}
}

var (
	_ datasource.DataSource              = &dataSystemDataSourceV2{}
	_ datasource.DataSourceWithConfigure = &dataSystemDataSourceV2{}
)

type dataSystemDataSourceV2 struct {
	client *neos.DataSystemClient
}

func (d *dataSystemDataSourceV2) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_system"
}

func (d *dataSystemDataSourceV2) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "Abi READ")

	var state DataSystemDataSourceModelV2

	list, err := d.client.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data System List",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Abi READ Post error ")

	tflog.Info(ctx, fmt.Sprintf("Abi READ length %d", len(list.Entities)))

	// Map response body to model
	for _, ds := range list.Entities {
		dataSystemState := DataSystemModelV2{
			Identifier:  types.StringValue(ds.Identifier),
			Name:        types.StringValue(ds.Name),
			Description: types.StringValue(ds.Description),
			Label:       types.StringValue(ds.Label),
			Owner:       types.StringValue(ds.Owner),
			Urn:         types.StringValue(ds.Urn),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Identifier))
		state.DataSystems = append(state.DataSystems, dataSystemState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *dataSystemDataSourceV2) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected dataSystemDataSourceV2 Configure Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))

		return
	}

	d.client = &client.DataSystemClient
}

func (d *dataSystemDataSourceV2) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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

type DataSystemDataSourceModelV2 struct {
	DataSystems []DataSystemModelV2 `tfsdk:"datasystems"`
}

// coffeesModel maps coffees schema data.
type DataSystemModelV2 struct {
	Identifier  types.String           `tfsdk:"id"`
	Urn         types.String           `tfsdk:"urn"`
	Name        types.String           `tfsdk:"name"`
	Description types.String           `tfsdk:"description"`
	Label       types.String           `tfsdk:"label"`
	Owner       types.String           `tfsdk:"owner"`
	CreatedAt   types.String           `tfsdk:"created_at"`
	State       DataSystemStateModelV2 `tfsdk:"state"`
}

type DataSystemStateModelV2 struct {
	State   types.String `tfsdk:"state"`
	Healthy types.Bool   `tfsdk:"healthy"`
}
