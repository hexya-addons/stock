package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

//import datetime
func init() {
	h.MakeProcurement().DeclareModel()

	h.MakeProcurement().AddFields(map[string]models.FieldDefinition{
		"Qty": models.FloatField{
			String:  "Quantity",
			Default: models.DefaultValue(1),
			//digits=(16, 2)
			Required: true,
		},
		"ResModel": models.CharField{
			String: "Res Model",
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Required:      true,
		},
		"ProductTmplId": models.Many2OneField{
			RelationModel: h.ProductTemplate(),
			String:        "Template",
			Required:      true,
		},
		"ProductVariantCount": models.IntegerField{
			String:  "Variant Number",
			Related: `ProductTmplId.ProductVariantCount`,
		},
		"UomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Unit of Measure",
			Required:      true,
		},
		"WarehouseId": models.Many2OneField{
			RelationModel: h.StockWarehouse(),
			String:        "Warehouse",
			Required:      true,
		},
		"DatePlanned": models.DateField{
			String:   "Planned Date",
			Default:  func(env models.Environment) interface{} { return odoo.fields.Date.context_today },
			Required: true,
		},
		"RouteIds": models.Many2ManyField{
			RelationModel: h.StockLocationRoute(),
			String:        "Preferred Routes",
		},
	})
	h.MakeProcurement().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.MakeProcurementSet, fields interface{}) {
			//        res = super(MakeProcurement, self).default_get(fields)
			//        if self.env.context.get('active_id') and self.env.context.get('active_model') == 'product.template':
			//            product = self.env['product.product'].search(
			//                [('product_tmpl_id', '=', self.env.context['active_id'])], limit=1)
			//        elif self.env.context.get('active_id') and self.env.context.get('active_model') == 'product.product':
			//            product = self.env['product.product'].browse(
			//                self.env.context['active_id'])
			//        else:
			//            product = self.env['product.product']
			//        if 'product_id' in fields and not res.get('product_id') and product:
			//            res['product_id'] = product.id
			//        if 'product_tmpl_id' in fields and not res.get('product_tmpl_id') and product:
			//            res['product_tmpl_id'] = product.product_tmpl_id.id
			//        if 'uom_id' in fields and not res.get('uom_id') and product:
			//            res['uom_id'] = product.uom_id.id
			//        if 'warehouse_id' in fields and not res.get('warehouse_id'):
			//            res['warehouse_id'] = self.env['stock.warehouse'].search(
			//                [], limit=1).id
			//        return res
		})
	h.MakeProcurement().Methods().OnchangeProductIdDict().DeclareMethod(
		`OnchangeProductIdDict`,
		func(rs m.MakeProcurementSet, product_id interface{}) {
			//        product = self.env['product.product'].browse(product_id)
			//        return {
			//            'uom_id': product.uom_id.id,
			//            'product_tmpl_id': product.product_tmpl_id.id,
			//            'product_variant_count': product.product_tmpl_id.product_variant_count
			//        }
		})
	h.MakeProcurement().Methods().OnchangeProductId().DeclareMethod(
		`OnchangeProductId`,
		func(rs m.MakeProcurementSet) {
			//        if self.product_id:
			//            for key, value in self.onchange_product_id_dict(self.product_id.id).iteritems():
			//                setattr(self, key, value)
		})
	h.MakeProcurement().Methods().Create().Extend(
		`Create`,
		func(rs m.MakeProcurementSet, values models.RecordData) {
			//        if values.get('product_id'):
			//            values.update(self.onchange_product_id_dict(values['product_id']))
			//        return super(MakeProcurement, self).create(values)
		})
	h.MakeProcurement().Methods().MakeProcurement().DeclareMethod(
		` Creates procurement order for selected product. `,
		func(rs m.MakeProcurementSet) {
			//        ProcurementOrder = self.env['procurement.order']
			//        for wizard in self:
			//            # we set the time to noon to avoid the date to be changed because of timezone issues
			//            date = fields.Datetime.from_string(wizard.date_planned)
			//            date = date + datetime.timedelta(hours=12)
			//            date = fields.Datetime.to_string(date)
			//
			//            procurement = ProcurementOrder.create({
			//                'name': 'INT: %s' % (self.env.user.login),
			//                'date_planned': date,
			//                'product_id': wizard.product_id.id,
			//                'product_qty': wizard.qty,
			//                'product_uom': wizard.uom_id.id,
			//                'warehouse_id': wizard.warehouse_id.id,
			//                'location_id': wizard.warehouse_id.lot_stock_id.id,
			//                'company_id': wizard.warehouse_id.company_id.id,
			//                'route_ids': [(6, 0, wizard.route_ids.ids)]})
			//        return {
			//            'view_type': 'form',
			//            'view_mode': 'tree,form',
			//            'res_model': 'procurement.order',
			//            'res_id': procurement.id,
			//            'views': [(False, 'form'), (False, 'tree')],
			//            'type': 'ir.actions.act_window',
			//        }
		})
}
