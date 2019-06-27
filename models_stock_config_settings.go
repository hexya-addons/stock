package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockConfigSettings().DeclareModel()

	h.StockConfigSettings().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockConfigSettingsSet, fields interface{}) {
			//        res = super(StockSettings, self).default_get(fields)
			//        if 'warehouse_and_location_usage_level' in fields or not fields:
			//            res['warehouse_and_location_usage_level'] = int(res.get(
			//                'group_stock_multi_locations', False)) + int(res.get('group_stock_multi_warehouses', False))
			//        return res
		})
	h.StockConfigSettings().AddFields(map[string]models.FieldDefinition{
		"GroupProductVariant": models.SelectionField{
			Selection: types.Selection{
				"": "No variants on products",
				"": "Products can have several attributes, defining variants (Example: size, color,...)",
			},
			String: "Product Variants",
			//implied_group='product.group_product_variant'
			Help: "Work with product variant allows you to define some variant" +
				"of the same products, an ease the product management in" +
				"the ecommerce for example",
		},
		"CompanyId": models.Many2OneField{
			RelationModel: h.Company(),
			String:        "Company",
			Default:       func(env models.Environment) interface{} { return env.Uid().company_id },
			Required:      true,
		},
		"ModuleProcurementJit": models.SelectionField{
			Selection: types.Selection{
				"": "Reserve products immediately after the sale order confirmation",
				"": "Reserve products manually or based on automatic scheduler",
			},
			String: "Procurements",
			Help: "Allows you to automatically reserve the available" +
				"        products when confirming a sale order." +
				"            This installs the module procurement_jit.",
		},
		"ModuleProductExpiry": models.SelectionField{
			Selection: types.Selection{
				"": "Do not use Expiration Date on serial numbers",
				"": "Define Expiration Date on serial numbers",
			},
			String: "Expiration Dates",
			Help: "Track different dates on products and serial numbers." +
				"                The following dates can be tracked:" +
				"                - end of life" +
				"                - best before date" +
				"                - removal date" +
				"                - alert date." +
				"                This installs the module product_expiry.",
		},
		"GroupUom": models.SelectionField{
			Selection: types.Selection{
				"": "Products have only one unit of measure (easier)",
				"": "Some products may be sold/purchased in different units of measure (advanced)",
			},
			String: "Units of Measure",
			//implied_group='product.group_uom'
			Help: "Allows you to select and maintain different units of measure" +
				"for products.",
		},
		"GroupStockPackaging": models.SelectionField{
			Selection: types.Selection{
				"": "Do not manage packaging",
				"": "Manage available packaging options per products",
			},
			String: "Packaging Methods",
			//implied_group='product.group_stock_packaging'
			Help: "Allows you to create and manage your packaging dimensions" +
				"and types you want to be maintained in your system.",
		},
		"GroupStockProductionLot": models.SelectionField{
			Selection: types.Selection{
				"": "Do not track individual product items",
				"": "Track lots or serial numbers",
			},
			String: "Lots and Serial Numbers",
			//implied_group='stock.group_production_lot'
			Help: "This allows you to assign a lot (or serial number) to the" +
				"pickings and moves.  This can make it possible to know" +
				"which production lot was sent to a certain client, ...",
		},
		"GroupStockTrackingLot": models.SelectionField{
			Selection: types.Selection{
				"": "Do not manage packaging",
				"": "Record packages used on packing: pallets, boxes, ...",
			},
			String: "Packages",
			//implied_group='stock.group_tracking_lot'
			Help: "This allows to manipulate packages.  You can put something" +
				"in, take something from a package, but also move entire" +
				"packages and put them even in another package.  ",
		},
		"GroupStockTrackingOwner": models.SelectionField{
			Selection: types.Selection{
				"": "All products in your warehouse belong to your company",
				"": "Manage consignee stocks (advanced)",
			},
			String: "Product Owners",
			//implied_group='stock.group_tracking_owner'
			Help: "This way you can receive products attributed to a certain owner. ",
		},
		"GroupStockAdvLocation": models.SelectionField{
			Selection: types.Selection{
				"": "No automatic routing of products",
				"": "Advanced routing of products using rules",
			},
			String: "Routes",
			//implied_group='stock.group_adv_location'
			Help: "This option supplements the warehouse application by effectively" +
				"implementing Push and Pull inventory flows through Routes.",
		},
		"GroupWarningStock": models.SelectionField{
			Selection: types.Selection{
				"": "All the partners can be used in pickings",
				"": "An informative or blocking warning can be set on a partner",
			},
			String: "Warning",
			//implied_group='stock.group_warning_stock'
		},
		"DecimalPrecision": models.IntegerField{
			String: "Decimal precision on weight",
			Help: "As an example, a decimal precision of 2 will allow weights" +
				"like: 9.99 kg, whereas a decimal precision of 4 will allow" +
				"weights like:  0.0231 kg.",
		},
		"PropagationMinimumDelta": models.IntegerField{
			String:  "Minimum days to trigger a propagation of date change in pushed/pull flows.",
			Related: `CompanyId.PropagationMinimumDelta`,
		},
		"ModuleStockDropshipping": models.SelectionField{
			Selection: types.Selection{
				"": "Suppliers always deliver to your warehouse(s)",
				"": "Allow suppliers to deliver directly to your customers",
			},
			String: "Dropshipping",
			Help: "" +
				"Creates the dropship route and add more complex tests" +
				"-This installs the module stock_dropshipping.",
		},
		"ModuleStockPickingWave": models.SelectionField{
			Selection: types.Selection{
				"": "Manage pickings one at a time",
				"": "Manage picking in batch per worker",
			},
			String: "Picking Waves",
			Help: "Install the picking wave module which will help you grouping" +
				"your pickings and processing them in batch",
		},
		"ModuleStockCalendar": models.SelectionField{
			Selection: types.Selection{
				"": "Set lead times in calendar days (easy)",
				"": "Adapt lead times using the suppliers' open days calendars (advanced)",
			},
			String: "Minimum Stock Rules",
			Help: "This allows you to handle minimum stock rules differently" +
				"by the possibility to take into account the purchase and" +
				"delivery calendars " +
				"-This installs the module stock_calendar.",
		},
		"ModuleStockBarcode": models.BooleanField{
			String: "Barcode scanner support",
		},
		"ModuleDeliveryDhl": models.BooleanField{
			String: "DHL integration",
		},
		"ModuleDeliveryFedex": models.BooleanField{
			String: "Fedex integration",
		},
		"ModuleDeliveryTemando": models.BooleanField{
			String: "Temando integration",
		},
		"ModuleDeliveryUps": models.BooleanField{
			String: "UPS integration",
		},
		"ModuleDeliveryUsps": models.BooleanField{
			String: "USPS integration",
		},
		"WarehouseAndLocationUsageLevel": models.SelectionField{
			Selection: types.Selection{
				"": "Manage only 1 Warehouse with only 1 stock location",
				"": "Manage only 1 Warehouse, composed by several stock locations",
				"": "Manage several Warehouses, each one composed by several stock locations",
			},
			String: "Warehouses and Locations usage level",
		},
		"GroupStockMultiLocations": models.BooleanField{
			String: "Manage several stock locations",
			//implied_group='stock.group_stock_multi_locations'
		},
		"GroupStockMultiWarehouses": models.BooleanField{
			String: "Manage several warehouses",
			//implied_group='stock.group_stock_multi_warehouses'
		},
		"ModuleQuality": models.BooleanField{
			String: "Quality",
			Help:   "This module allows you to generate quality alerts and quality check",
		},
	})
	h.StockConfigSettings().Methods().OnchangeWarehouseAndLocationUsageLevel().DeclareMethod(
		`OnchangeWarehouseAndLocationUsageLevel`,
		func(rs m.StockConfigSettingsSet) {
			//        self.group_stock_multi_locations = self.warehouse_and_location_usage_level > 0
			//        self.group_stock_multi_warehouses = self.warehouse_and_location_usage_level > 1
		})
	h.StockConfigSettings().Methods().OnchangeAdvLocation().DeclareMethod(
		`OnchangeAdvLocation`,
		func(rs m.StockConfigSettingsSet) {
			//        if self.group_stock_adv_location and self.warehouse_and_location_usage_level == 0:
			//            self.warehouse_and_location_usage_level = 1
		})
	h.StockConfigSettings().Methods().SetGroupStockMultiLocations().DeclareMethod(
		` If we are not in multiple locations, we can deactivate the internal
        picking types of the warehouses, so they won't
appear in the dashboard.
        Otherwise, activate them.
        `,
		func(rs m.StockConfigSettingsSet) {
			//        for config in self:
			//            if config.group_stock_multi_locations:
			//                active = True
			//                domain = []
			//            else:
			//                active = False
			//                domain = [('reception_steps', '=', 'one_step'),
			//                          ('delivery_steps', '=', 'ship_only')]
			//
			//            warehouses = self.env['stock.warehouse'].search(domain)
			//            warehouses.mapped('int_type_id').write({'active': active})
			//        return True
		})
	h.StockConfigSettings().Methods().GetDefaultDecimalPrecision().DeclareMethod(
		`GetDefaultDecimalPrecision`,
		func(rs m.StockConfigSettingsSet, fields interface{}) {
			//        digits = self.env.ref('product.decimal_stock_weight').digits
			//        return {'decimal_precision': digits}
		})
	h.StockConfigSettings().Methods().SetDecimalPrecision().DeclareMethod(
		`SetDecimalPrecision`,
		func(rs m.StockConfigSettingsSet) {
			//        for record in self:
			//            self.env.ref('product.decimal_stock_weight').write(
			//                {'digits': record.decimal_precision})
		})
}
