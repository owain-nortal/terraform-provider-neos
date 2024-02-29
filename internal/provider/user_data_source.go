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

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

type userDataSource struct {
	client *neos.UserClient
}

func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	tflog.Info(ctx, "userDataSource READ")

	var state UserDataSourceModel

	list, err := d.client.List("","","")
	if err != nil {
		resp.Diagnostics.AddError("Unable to Read User List", err.Error())
		return
	}

	// Map response body to model
	for _, ds := range list.Users {
		userState := UserModel{
			ID:        types.StringValue(ds.Identifier),
			FirstName: types.StringValue(ds.FirstName),
			LastName:  types.StringValue(ds.LastName),
			Username:  types.StringValue(ds.Username),
			URN:       types.StringValue(ds.Urn),
			Email:     types.StringValue(ds.Email),
			Enabled:   types.BoolValue(ds.Enabled),
			IsSystem:  types.BoolValue(ds.IsSystem),
			Account:   types.StringValue(ds.Account),
		}
		state.UserModel = append(state.UserModel, userState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	tflog.Info(ctx, "Data source configure")

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)
	if !ok {
		resp.Diagnostics.AddError("Unexpected userDataSource Type", fmt.Sprintf("Expected *neos.NeosClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	d.client = &client.UserClient
}

func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"user": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"first_name": schema.StringAttribute{
							Computed: true,
						},
						"last_name": schema.StringAttribute{
							Computed: true,
						},
						"username": schema.StringAttribute{
							Computed: true,
						},
						"email": schema.StringAttribute{
							Computed: true,
						},
						"urn": schema.StringAttribute{
							Computed: true,
						},
						"id": schema.StringAttribute{
							Computed: true,
						},
						"enabled": schema.BoolAttribute{
							Computed: true,
						},
						"is_system": schema.BoolAttribute{
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

type UserDataSourceModel struct {
	UserModel []UserModel `tfsdk:"user"`
}

type UserModel struct {
	FirstName   types.String `tfsdk:"first_name"`
	LastName    types.String `tfsdk:"last_name"`
	Username    types.String `tfsdk:"username"`
	Email       types.String `tfsdk:"email"`
	URN         types.String `tfsdk:"urn"`
	ID          types.String `tfsdk:"id"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	Account     types.String `tfsdk:"account"`
	//LastUpdated types.String `tfsdk:"last_updated"`	
}
