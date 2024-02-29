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

func NewLinksDataSource() datasource.DataSource {
	return &linksDataSourceV2{}
}

var (
	_ datasource.DataSource              = &linksDataSourceV2{}
	_ datasource.DataSourceWithConfigure = &linksDataSourceV2{}
)

type linksDataSourceV2 struct {
	client *neos.LinksClient
}

func (d *linksDataSourceV2) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_links"
}

func (d *linksDataSourceV2) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "Links READ")

	var state LinksDataSourceModelV2

	list, err := d.client.Get()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Data System List",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "Links READ Post converting to model")

	//tflog.Info(ctx, fmt.Sprintf("Abi READ length %d", len(list.Entities)))

	// Map response body to model
	for _, ds := range list.Links {
		tflog.Info(ctx, fmt.Sprintf("Links read ds parent id: %s child id: %s ", ds.Parent.Identifier, ds.Child.Identifier))

		// state := LinksDataSourceModelV2{}

		parent := LinksParentModel{
			Identifier:  types.StringValue(ds.Parent.Identifier),
			Urn:         types.StringValue(ds.Parent.Urn),
			Name:        types.StringValue(ds.Parent.Name),
			IsSystem:    types.BoolValue(ds.Parent.IsSystem),
			Description: types.StringValue(ds.Parent.Description),
			Label:       types.StringValue(ds.Parent.Label),
			CreatedAt:   types.StringValue(ds.Parent.CreatedAt.String()),
			//State       LinksStateModelV2 `tfsdk:"state"`
			Owner:      types.StringValue(ds.Parent.Owner),
			EntityType: types.StringValue(ds.Parent.EntityType),
			OutputType: types.StringValue(ds.Parent.OutputType),
		}
		child := LinksChildModel{
			Identifier:  types.StringValue(ds.Child.Identifier),
			Urn:         types.StringValue(ds.Child.Urn),
			Name:        types.StringValue(ds.Child.Name),
			IsSystem:    types.BoolValue(ds.Child.IsSystem),
			Description: types.StringValue(ds.Child.Description),
			Label:       types.StringValue(ds.Child.Label),
			CreatedAt:   types.StringValue(ds.Child.CreatedAt.String()),
			//State       LinksStateModelV2 `tfsdk:"state"`
			Owner:      types.StringValue(ds.Child.Owner),
			EntityType: types.StringValue(ds.Child.EntityType),
			OutputType: types.StringValue(ds.Child.OutputType),
		}

		linksState := LinksModelV2{
			Tmp:    types.StringValue("abc123"),
			Parent: parent,
			Child:  child,
		}

		state.Links = append(state.Links, linksState)
	}

	tflog.Info(ctx, "NEOS Link Setting state")
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *linksDataSourceV2) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Links Data source v2 configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected linksDataSourceV2 Configure Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))

		return
	}

	d.client = &client.LinksClient
}

func (d *linksDataSourceV2) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"links": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tmp": schema.StringAttribute{
							Computed: true,
						},
						"parent": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"urn": schema.StringAttribute{
									Computed: true,
								},
								"is_system": schema.BoolAttribute{
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
								"entity_type": schema.StringAttribute{
									Computed: true,
								},
								"output_type": schema.StringAttribute{
									Computed: true,
								},
								// "state": schema.SingleNestedAttribute{
								// 	Computed: true,
								// 	Attributes: map[string]schema.Attribute{
								// 		"code": schema.StringAttribute{
								// 			Computed: true,
								// 		},
								// 		"state": schema.StringAttribute{
								// 			Computed: true,
								// 		},
								// 		"healthy": schema.BoolAttribute{
								// 			Computed: true,
								// 		},
								// 	},
								// },
							},
						},
						"child": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"identifier": schema.StringAttribute{
									Computed: true,
								},
								"name": schema.StringAttribute{
									Computed: true,
								},
								"urn": schema.StringAttribute{
									Computed: true,
								},
								"description": schema.StringAttribute{
									Computed: true,
								},
								"is_system": schema.BoolAttribute{
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
								"entity_type": schema.StringAttribute{
									Computed: true,
								},
								"output_type": schema.StringAttribute{
									Computed: true,
								},
								// "state": schema.SingleNestedAttribute{
								// 	Computed: true,
								// 	Attributes: map[string]schema.Attribute{
								// 		"code": schema.StringAttribute{
								// 			Computed: true,
								// 		},
								// 		"state": schema.StringAttribute{
								// 			Computed: true,
								// 		},
								// 		"healthy": schema.BoolAttribute{
								// 			Computed: true,
								// 		},
								// 	},
								// },
							},
						},
					},
				},
			},
		},
	}
}

type LinksDataSourceModelV2 struct {
	Links []LinksModelV2 `tfsdk:"links"`
}

type LinksStateModelV2 struct {
	Code    types.String `tfsdk:"code"`
	Reason  types.String `tfsdk:"reason"`
	Healthy types.Bool   `tfsdk:"healthy"`
}

type LinksModelV2 struct {
	Tmp    types.String     `tfsdk:"tmp"`
	Parent LinksParentModel `tfsdk:"parent"`
	Child  LinksChildModel  `tfsdk:"child"`
}

type LinksChildModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	Urn         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	Description types.String `tfsdk:"description"`
	Label       types.String `tfsdk:"label"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Owner       types.String `tfsdk:"owner"`
	EntityType  types.String `tfsdk:"entity_type"`
	OutputType  types.String `tfsdk:"output_type"`
	//State       LinksStateModelV2 `tfsdk:"state"`
}

type LinksParentModel struct {
	Identifier  types.String `tfsdk:"identifier"`
	Urn         types.String `tfsdk:"urn"`
	Name        types.String `tfsdk:"name"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	Description types.String `tfsdk:"description"`
	Label       types.String `tfsdk:"label"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Owner       types.String `tfsdk:"owner"`
	EntityType  types.String `tfsdk:"entity_type"`
	OutputType  types.String `tfsdk:"output_type"`
	//State       LinksStateModelV2 `tfsdk:"state"`
}
