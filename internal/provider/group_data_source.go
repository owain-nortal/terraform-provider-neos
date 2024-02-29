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

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

type groupDataSource struct {
	client *neos.GroupClient
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state GroupDataSourceModel

	list, err := d.client.List("")
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read Data System List", err.Error())
		return
	}

	// Map response body to model
	for _, ds := range list.Groups {
		group := GroupModel{
			Identifier:  types.StringValue(ds.Identifier),
			Name:        types.StringValue(ds.Name),
			Description: types.StringValue(ds.Description),
			IsSystem:    types.BoolValue(ds.IsSystem),
		}
		for _, v := range ds.Principals {
			group.Principals = append(group.Principals, types.StringValue(v))
		}
		tflog.Info(ctx, fmt.Sprintf("NEOS - ID: %s ", ds.Identifier))
		state.GroupModel = append(state.GroupModel, group)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected groupDataSource Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	d.client = &client.GroupClient
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"groups": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"is_system": schema.BoolAttribute{
							Computed: true,
						},
						"principals": schema.ListAttribute{
							ElementType: types.StringType,
							Computed:    false,
							Optional:    true,
							Required:    false,
							Description: "list of principals",
						},
					},
				},
			},
		},
	}
}

type GroupDataSourceModel struct {
	GroupModel []GroupModel `tfsdk:"group"`
}

// coffeesModel maps coffees schema data.
type GroupModel struct {
	Identifier  types.String   `tfsdk:"id"`
	Description types.String   `tfsdk:"description"`
	Name        types.String   `tfsdk:"name"`
	IsSystem    types.Bool     `tfsdk:"is_system"`
	Principals  []types.String `tfsdk:"printcipals"`
}
