package provider

import (
	"context"
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
func NewDataProductResource() resource.Resource {
	return &dataProductResource{}
}

// dataProductResource is the resource implementation.
type dataProductResource struct {
	client       *neos.DataProductClient
	schemaClient *neos.DataProductSchemaClient
}

var (
	_ resource.Resource                = &dataProductResource{}
	_ resource.ResourceWithConfigure   = &dataProductResource{}
	_ resource.ResourceWithImportState = &dataProductResource{}
)

// Metadata returns the resource type name.
func (r *dataProductResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_product"
}

// Schema defines the schema for the resource.
func (r *dataProductResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The Unique ID of the data product",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"urn": schema.StringAttribute{
				Computed:    true,
				Required:    false,
				Optional:    false,
				Description: "The URN of the data product which is read only",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Computed:    false,
				Required:    true,
				Optional:    false,
				Description: "Name of the data product",
			},
			"description": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Description of the data product",
			},
			"label": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "Label for the data product",
			},
			"owner": schema.StringAttribute{
				Computed:    false,
				Optional:    true,
				Required:    false,
				Description: "The owner of the data product",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Optional:    false,
				Required:    false,
				Description: "when the data product was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"contact_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "list of contacts Ids",
			},
			"links": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    false,
				Optional:    false,
				Required:    true,
				Description: "list of links",
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},

			"schema": schema.SingleNestedAttribute{
				Computed: false,
				Optional: true,
				Required: false,
				Attributes: map[string]schema.Attribute{
					"product_type": schema.StringAttribute{
						Computed:    false,
						Optional:    true,
						Required:    false,
						Description: "product type 'stored' etc",
					},
					"fields": schema.ListNestedAttribute{
						Computed: false,
						Optional: true,
						Required: false,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									Computed:    false,
									Required:    false,
									Optional:    true,
									Description: "Name of the schema field",
								},
								"description": schema.StringAttribute{
									Computed:    false,
									Optional:    false,
									Required:    true,
									Description: "Description of the schema field",
								},
								"primary": schema.BoolAttribute{
									Computed:    false,
									Optional:    true,
									Required:    false,
									Description: "set the schmea field to be a primary key",
								},
								"optional": schema.BoolAttribute{
									Computed:    false,
									Optional:    true,
									Required:    false,
									Description: "set the schmea field to be a optional",
								},
								"data_type": schema.SingleNestedAttribute{
									Computed: false,
									Optional: true,
									Required: false,
									Attributes: map[string]schema.Attribute{
										"column_type": schema.StringAttribute{
											Computed:    false,
											Optional:    true,
											Required:    false,
											Description: "set the schmea field column type",
										},
										"meta": schema.MapAttribute{
											ElementType: types.StringType,
											Computed:    false,
											Optional:    true,
											Required:    false,
										},
									},
									Description: "set the schmea field data type",
								},
							},
						},
					},
				},
			},
		},
	}
}

// dataProductResourceModel maps the resource schema data.
type dataProductResourceModel struct {
	ID          types.String           `tfsdk:"id"`
	URN         types.String           `tfsdk:"urn"`
	Name        types.String           `tfsdk:"name"`
	Label       types.String           `tfsdk:"label"`
	Description types.String           `tfsdk:"description"`
	Owner       types.String           `tfsdk:"owner"`
	CreatedAt   types.String           `tfsdk:"created_at"`
	Links       types.List             `tfsdk:"links"`
	ContactIds  types.List             `tfsdk:"contact_ids"`
	LastUpdated types.String           `tfsdk:"last_updated"`
	Schema      DataProductSchemaModel `tfsdk:"schema"`
}

type DataProductSchemaModel struct {
	ProductType types.String                    `tfsdk:"product_type"`
	Fields      []DataProductFieldResourceModel `tfsdk:"fields"`
}

type DataProductFieldResourceModel struct {
	Name        types.String                     `tfsdk:"name"`
	Description types.String                     `tfsdk:"description"`
	Primary     types.Bool                       `tfsdk:"primary"`
	Optional    types.Bool                       `tfsdk:"optional"`
	DataType    DataProductDataTypeResourceModel `tfsdk:"data_type"`
}

type DataProductDataTypeResourceModel struct {
	Meta       types.Map    `tfsdk:"meta"`
	ColumnType types.String `tfsdk:"column_type"`
}

// Create a new resource.
func (r *dataProductResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "dataProductResource Create Get plan")
	var plan dataProductResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "dataProductResource Create building data product")

	linkList, diag := plan.Links.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var links = make([]string, 0)
	diag = linkList.ElementsAs(ctx, &links, false)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactIDs, diag := plan.ContactIds.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var contacts = make([]string, 0)
	diag = contactIDs.ElementsAs(ctx, &contacts, false)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() == "retentiondataproduct" {
		fmt.Println("retentiondataproduct")
	}

	item := neos.DataProductPostRequest{
		Entity: neos.DataProductPostRequestEntity{
			Name:        plan.Name.ValueString(),
			Label:       plan.Label.ValueString(),
			Description: plan.Description.ValueString(),
		},
		EntityInfo: neos.DataProductPostRequestEntityInfo{
			Owner:      plan.Owner.ValueString(),
			ContactIds: contacts,
			Links:      links,
		},
	}

	tflog.Info(ctx, "dataProductResource Create date product post")
	result, err := r.client.Post(ctx, item)
	if err != nil {
		resp.Diagnostics.AddError("Error creating data product", "Could not create data product, unexpected error: "+err.Error())
		return
	}

	id := result.Identifier
	tflog.Info(ctx, fmt.Sprintf("dataProductResource Create date product post id %s", id))

	fields := []neos.DataProductSchemaFieldPutRequest{}

	for _, v := range plan.Schema.Fields {

		meta := make(map[string]string)
		diag := v.DataType.Meta.ElementsAs(ctx, &meta, true)
		if diag.HasError() {
			tflog.Info(ctx, fmt.Sprintf("%v", diag.Errors()))
			resp.Diagnostics.AddError("Error creating v.DataType.Meta.ToMapValue", "Error")
			return
		}

		dataType := neos.DataProductSchemaDataTypePutRequest{
			Meta:       meta,
			ColumnType: v.DataType.ColumnType.ValueString(),
		}

		f := neos.DataProductSchemaFieldPutRequest{
			Description: v.Description.ValueString(),
			Name:        v.Name.ValueString(),
			Primary:     v.Primary.ValueBool(),
			Optional:    v.Optional.ValueBool(),
			DataType:    dataType,
		}
		fields = append(fields, f)
	}

	if plan.Schema.ProductType.ValueString() != "" && len(fields) != 0 {
		schemaPutRequest := neos.DataProductSchemaPutRequest{
			Details: neos.DataProductSchemaDetailsPutRequest{
				ProductType: plan.Schema.ProductType.ValueString(),
				Fields:      fields,
			},
		}

		tflog.Info(ctx, fmt.Sprintf("dataProductResource Create schema put %s", id))
		schemaResult, err := r.schemaClient.Put(ctx, id, schemaPutRequest)
		if err != nil {
			resp.Diagnostics.AddError("Error putting data product schema ", "Could not create data product schema, unexpected error: "+err.Error())
			return
		}

		pfields := []DataProductFieldResourceModel{}

		for _, v := range schemaResult.Fields {

			tflog.Info(ctx, "dataProductResource into schemaResult.Fields")
			tflog.Info(ctx, fmt.Sprintf("dataProductResource ColumnType values [%s] ", types.StringValue(v.DataType.ColumnType)))

			meta, diag := types.MapValueFrom(ctx, types.StringType, v.DataType.Meta)
			resp.Diagnostics.Append(diag...)
			if resp.Diagnostics.HasError() {
				resp.Diagnostics.AddError("ErrorMapping values for datatype meta ", "Could not create data product schema, unexpected error: "+err.Error())
				return
			}

			dt := DataProductDataTypeResourceModel{
				Meta:       meta,
				ColumnType: types.StringValue(v.DataType.ColumnType),
			}

			tflog.Info(ctx, fmt.Sprintf("dataProductResource field values [%s] [%v] [%v] [%s] ", types.StringValue(v.Name), types.BoolValue(v.Primary), types.BoolValue(v.Optional), types.StringValue(v.Description)))

			i := DataProductFieldResourceModel{
				Name:        types.StringValue(v.Name),
				Primary:     types.BoolValue(v.Primary),
				Optional:    types.BoolValue(v.Optional),
				Description: types.StringValue(v.Description),
				DataType:    dt,
			}
			pfields = append(pfields, i)
		}
		plan.Schema = DataProductSchemaModel{
			ProductType: types.StringValue(schemaPutRequest.Details.ProductType),
			Fields:      pfields,
		}
	}

	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.Label = types.StringValue(result.Label)
	plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
// Read resource information.
func (r *dataProductResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Info(ctx, "DP READ Get current state")

	var state dataProductResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	foo := fmt.Sprintf("DP READ state id: [%s]  Desc [%s]", state.ID.ValueString(), state.Description.ValueString())
	tflog.Info(ctx, foo)

	dataProductList, err := r.client.Get()
	if err != nil {
		resp.Diagnostics.AddError("Error Reading NEOS data product", "Could not read NEOS  data product ID "+state.ID.ValueString()+": "+err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("DP READ iterate over list looking for: %s", state.ID.ValueString()))
	for _, ds := range dataProductList.Entities {
		//		tflog.Info(ctx, fmt.Sprintf("££ READ ITEM: [%s] [%s] %v", ds.Identifier, state.ID.ValueString(), (ds.Identifier == state.ID.ValueString())))
		if ds.Identifier == state.ID.ValueString() {
			state.ID = types.StringValue(ds.Identifier)
			state.Name = types.StringValue(ds.Name)
			state.Label = types.StringValue(ds.Label)
			state.URN = types.StringValue(ds.Urn)
			state.Description = types.StringValue(ds.Description)
			state.Owner = types.StringValue(ds.Owner)
			state.CreatedAt = types.StringValue(ds.CreatedAt.String())

			dataProductSchema, err := r.schemaClient.Get(ds.Identifier)
			if err != nil {
				// no schema so assume its not be created rather than an error
				dataProductSchema = neos.DataProductSchema{}
				//resp.Diagnostics.AddError("Error Reading NEOS data product", "Could not read NEOS schema data product ID "+state.ID.ValueString()+": "+err.Error())
				//return
			}
			dpsm, shouldReturn := convertSchemaToModel(ctx, dataProductSchema, resp)
			if shouldReturn {
				return
			}

			state.Schema = dpsm
			//state.Schema.ProductType = types.StringValue("stored")

			break
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Info(ctx, "data product Read Has error")
		return
	}

}

func convertSchemaToModel(ctx context.Context, dataProductSchema neos.DataProductSchema, resp *resource.ReadResponse) (DataProductSchemaModel, bool) {
	fields := []DataProductFieldResourceModel{}
	for _, v := range dataProductSchema.Fields {
		meta, diag := types.MapValueFrom(ctx, types.StringType, v.DataType.Meta)
		resp.Diagnostics.Append(diag...)
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("Error Mapping values for datatype meta ", "unexpected error")
			return DataProductSchemaModel{}, true
		}

		field := DataProductFieldResourceModel{
			Name:        types.StringValue(v.Name),
			Description: types.StringValue(v.Description),
			Primary:     types.BoolValue(v.Primary),
			Optional:    types.BoolValue(v.Optional),
			DataType: DataProductDataTypeResourceModel{
				ColumnType: types.StringValue(v.DataType.ColumnType),
				Meta:       meta,
			},
		}
		fields = append(fields, field)
	}

	dpsm := DataProductSchemaModel{
		Fields: fields,
	}
	return dpsm, false
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *dataProductResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	tflog.Info(ctx, "dataProductResource Update called ")

	var plan dataProductResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	linkList, diag := plan.Links.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var links = make([]string, 0)
	diag = linkList.ElementsAs(ctx, &links, false)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	contactIDs, diag := plan.ContactIds.ToListValue(ctx)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}
	var contacts = make([]string, 0)
	diag = contactIDs.ElementsAs(ctx, &contacts, false)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	item := neos.DataProductPutRequest{
		Entity: neos.DataProductPutRequestEntity{
			Name:        plan.Name.ValueString(),
			Label:       plan.Label.ValueString(),
			Description: plan.Description.ValueString(),
		},
	}

	eItem := neos.DataProductPutRequestEntityInfo{
		Owner:      plan.Owner.ValueString(),
		ContactIds: contacts,
		Links:      links,
	}

	result, err := r.client.Put(ctx, plan.ID.ValueString(), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data product",
			"Could not put data product, unexpected error: "+err.Error(),
		)
		return
	}
	infoResult, err := r.client.DataProductPutInfo(ctx, plan.ID.ValueString(), eItem)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating data product",
			"Could not put data product, unexpected error: "+err.Error(),
		)
		return
	}

	contactsList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.ContactIds)
	linksList, _ := types.ListValueFrom(ctx, types.StringType, infoResult.Links)

	plan.ID = types.StringValue(result.Identifier)
	plan.Name = types.StringValue(result.Name)
	plan.URN = types.StringValue(result.Urn)
	plan.Description = types.StringValue(result.Description)
	plan.Label = types.StringValue(result.Label)
	plan.CreatedAt = types.StringValue(result.CreatedAt.String())
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ContactIds = contactsList
	plan.Links = linksList
	plan.Owner = types.StringValue(infoResult.Owner)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	id := result.Identifier

	tflog.Info(ctx, fmt.Sprintf("dataProductResource Update date product post id %s", id))
	fields := []neos.DataProductSchemaFieldPutRequest{}

	for _, v := range plan.Schema.Fields {

		meta := make(map[string]string)

		diag := v.DataType.Meta.ElementsAs(ctx, &meta, true)
		if diag.HasError() {
			tflog.Info(ctx, fmt.Sprintf("%v", diag.Errors()))
			resp.Diagnostics.AddError("Error creating v.DataType.Meta.ToMapValue", "Error")
			return
		}

		dataType := neos.DataProductSchemaDataTypePutRequest{
			Meta:       meta,
			ColumnType: v.DataType.ColumnType.ValueString(),
		}

		f := neos.DataProductSchemaFieldPutRequest{
			Description: v.Description.ValueString(),
			Name:        v.Name.ValueString(),
			Primary:     v.Primary.ValueBool(),
			Optional:    v.Optional.ValueBool(),
			DataType:    dataType,
		}
		fields = append(fields, f)
	}

	schemaPutRequest := neos.DataProductSchemaPutRequest{
		Details: neos.DataProductSchemaDetailsPutRequest{
			ProductType: plan.Schema.ProductType.ValueString(),
			Fields:      fields,
		},
	}

	if plan.Schema.ProductType.ValueString() != "" && len(fields) != 0 {

		tflog.Info(ctx, fmt.Sprintf("dataProductResource update schema put %s", id))
		schemaResult, err := r.schemaClient.Put(ctx, id, schemaPutRequest)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error putting data product schema ",
				"Could not update data product schema, unexpected error: "+err.Error(),
			)
			return
		}

		pfields := []DataProductFieldResourceModel{}

		for _, v := range schemaResult.Fields {
			tflog.Info(ctx, fmt.Sprintf("dataProductResource ColumnType values [%s] ", types.StringValue(v.DataType.ColumnType)))
			meta, diag := types.MapValueFrom(ctx, types.StringType, v.DataType.Meta)
			resp.Diagnostics.Append(diag...)
			if resp.Diagnostics.HasError() {
				resp.Diagnostics.AddError(
					"ErrorMapping values for datatype meta ",
					"Could not update data product schema, unexpected error: "+err.Error(),
				)
				return
			}

			dt := DataProductDataTypeResourceModel{
				Meta:       meta,
				ColumnType: types.StringValue(v.DataType.ColumnType),
			}

			tflog.Info(ctx, fmt.Sprintf("dataProductResource field values [%s] [%v] [%v] [%s]", types.StringValue(v.Name), types.BoolValue(v.Primary), types.BoolValue(v.Optional), types.StringValue(v.Description)))

			i := DataProductFieldResourceModel{
				Name:        types.StringValue(v.Name),
				Primary:     types.BoolValue(v.Primary),
				Optional:    types.BoolValue(v.Optional),
				Description: types.StringValue(v.Description),
				DataType:    dt,
			}
			pfields = append(pfields, i)
		}

		plan.Schema = DataProductSchemaModel{
			ProductType: types.StringValue(schemaPutRequest.Details.ProductType),
			Fields:      pfields,
		}
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *dataProductResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan dataProductResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := plan.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("DP Delete iterate plan ID: %s", id))

	err := r.client.Delete(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting data product", "Could not delete data product, unexpected error: "+err.Error())
		return
	}

}

func (r *dataProductResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*neos.NeosClient)

	if !ok {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.DataProductClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	r.client = &client.DataProductClient

	// schemaClient, ok := req.ProviderData.(*neos.DataProductSchemaClient)

	// if !ok {
	// 	resp.Diagnostics.AddError("Unexpected Data Source Configure Type", fmt.Sprintf("Expected *neos.DataProductClient, got: %T. Please report this issue to the provider developers.", req.ProviderData))
	// 	return
	// }

	r.schemaClient = &client.DataProductSchemaClient
}

func (r *dataProductResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
