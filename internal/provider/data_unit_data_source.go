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

func NewDataUnitDataSource() datasource.DataSource {
	return &dataUnitDataSourceV2{}
}

var (
	_ datasource.DataSource              = &dataUnitDataSourceV2{}
	_ datasource.DataSourceWithConfigure = &dataUnitDataSourceV2{}
)

type dataUnitDataSourceV2 struct {
	client *neos.NeosClient
}

func (d *dataUnitDataSourceV2) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_unit"
}

func (d *dataUnitDataSourceV2) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "Abi READ")

	var state DataUnitDataSourceModelV2

	list, err := d.client.DataUnitGet()
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

		configJson, err := d.client.DataUnitConfigGetBase(ctx, ds.Identifier)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Data System List",
				err.Error(),
			)
			return
		}
		dataUnitState := DataUnitModelV2{
			Identifier:  types.StringValue(ds.Identifier),
			Name:        types.StringValue(ds.Name),
			Description: types.StringValue(ds.Description),
			Label:       types.StringValue(ds.Label),
			Owner:       types.StringValue(ds.Owner),
			Urn:         types.StringValue(ds.Urn),
			ConfigJson:  types.StringValue(string(configJson)),
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Identifier))
		state.Entities = append(state.Entities, dataUnitState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *dataUnitDataSourceV2) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "dataUnitDataSourceV2  Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *neos.DataUnitClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *dataUnitDataSourceV2) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"entities": schema.ListNestedAttribute{
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
						"config_json": schema.StringAttribute{
							Computed: true,
							//Optional:    true,
							//Required:    false,
							//Description: "json that describes the configuration of the data unit",
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

type DataUnitDataSourceModelV2 struct {
	Entities []DataUnitModelV2 `tfsdk:"entities"`
}

// coffeesModel maps coffees schema data.
type DataUnitModelV2 struct {
	Identifier  types.String         `tfsdk:"id"`
	Urn         types.String         `tfsdk:"urn"`
	Name        types.String         `tfsdk:"name"`
	Description types.String         `tfsdk:"description"`
	Label       types.String         `tfsdk:"label"`
	Owner       types.String         `tfsdk:"owner"`
	CreatedAt   types.String         `tfsdk:"created_at"`
	State       DataUnitStateModelV2 `tfsdk:"state"`
	ConfigJson  types.String         `tfsdk:"config_json"`
}

type DataUnitStateModelV2 struct {
	State   types.String `tfsdk:"state"`
	Healthy types.Bool   `tfsdk:"healthy"`
}
