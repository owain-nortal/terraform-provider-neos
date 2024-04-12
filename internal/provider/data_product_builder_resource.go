package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	neos "github.com/owain-nortal/neos-client-go"
	"time"
)

// New data productResource is a helper function to simplify the provider implementation.
func NewDataProductBuilderResource() resource.Resource {
	return &DataProductBuilderResource{}
}

// DataProductBuilderResource is the resource implementation.
type DataProductBuilderResource struct {
	client *neos.DataProductClient
}

var (
	_ resource.Resource                = &DataProductBuilderResource{}
	_ resource.ResourceWithConfigure   = &DataProductBuilderResource{}
	_ resource.ResourceWithImportState = &DataProductBuilderResource{}
)

// Metadata returns the resource type name.
func (r *DataProductBuilderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_product_builder"
}

// Schema defines the schema for the resource.
func (r *DataProductBuilderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The Unique ID of the data product",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"dataunit_datasource_linkids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "The link ids of the data unit data source, to ensure the correct dependency graph is created",
			},

			"builder_json": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "builder json",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// DataProductBuilderResourceModel maps the resource schema data.
type DataProductBuilderResourceModel struct {
	ID                        types.String `tfsdk:"id"`
	DataUnitDataSourceLinkIds types.List   `tfsdk:"dataunit_datasource_linkids"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	BuilderJson               types.String `tfsdk:"builder_json"`
}

// type DataProductSchemaModel struct {
// 	ProductType types.String                    `tfsdk:"product_type"`
// 	Fields      []DataProductFieldResourceModel `tfsdk:"fields"`
// }

// type DataProductFieldResourceModel struct {
// 	Name        types.String                     `tfsdk:"name"`
// 	Description types.String                     `tfsdk:"description"`
// 	Primary     types.Bool                       `tfsdk:"primary"`
// 	Optional    types.Bool                       `tfsdk:"optional"`
// 	DataType    DataProductDataTypeResourceModel `tfsdk:"data_type"`
// }

// Create a new resource.
func (r *DataProductBuilderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "DataProductBuilderResource Create Get plan")
	var plan DataProductBuilderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "DataProductBuilderResource Create building data product")

	builderJson := plan.BuilderJson.ValueString()

	// d1 := []byte(builderJson)
	// os.WriteFile(fmt.Sprintf("/tmp/%s.json", plan.Name.ValueString()), d1, 0644)

	tflog.Debug(ctx, builderJson)

	if builderJson != "" {
		if json.Valid([]byte(builderJson)) {
			tflog.Info(ctx, fmt.Sprintf("DataProductBuilderResource Create builder put %s", plan.ID.ValueString()))
			_, err := r.client.DataProductBuilderPut(ctx, plan.ID.ValueString(), builderJson)
			if err != nil {
				resp.Diagnostics.AddError("Error putting data product builder ", "Could not create data product builder, unexpected error: "+err.Error())
				return
			}
			plan.BuilderJson = types.StringValue(builderJson)
		} else {
			resp.Diagnostics.AddError("Error invalid json data product builder ", "the builder json is invalid")
			return
		}
	} else {
		tflog.Info(ctx, "DataProductBuilderResource no builder json found ")
	}

	//plan.ID = types.StringValue(plan.ID.ValueString())

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	foo := plan.DataUnitDataSourceLinkIds
	plan.DataUnitDataSourceLinkIds = foo
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *DataProductBuilderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Info(ctx, "DP READ Get current state")

	var state DataProductBuilderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataProductbuilderJson, err := r.client.DataProductBuilderGet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS data product builder ", "Could not read NEOS data product builder ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	state.BuilderJson = types.StringValue(dataProductbuilderJson)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "data product Read Has error")
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *DataProductBuilderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Info(ctx, "DataProductBuilderResource Update called ")

	var plan DataProductBuilderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DataProductBuilderPut(ctx, plan.ID.ValueString(), plan.BuilderJson.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating data product builder ", "Could not put data product builder, unexpected error: "+err.Error())
		return
	}

	dpbj, err := r.client.DataProductBuilderGet(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS data product builder after update", "Could not read NEOS data product builder ID "+plan.ID.ValueString()+": "+err.Error())
		return
	}
	foo := plan.DataUnitDataSourceLinkIds
	plan.DataUnitDataSourceLinkIds = foo
	plan.BuilderJson = types.StringValue(dpbj)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *DataProductBuilderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan DataProductBuilderResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("DP Builder Delete ID: %s", plan.ID.ValueString()))

	// Delete the data product builder json not currently supported

	// err := r.client.DataProductBuilderDelete(plan.ID.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError("Error deleting data product", "Could not delete data product, unexpected error: "+err.Error())
	// 	return
	// }

}

func (r *DataProductBuilderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.DataProductClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.DataProductClient
}

func (r *DataProductBuilderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
