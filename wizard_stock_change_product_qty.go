package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

func init() {
	h.StockChangeProductQty().DeclareModel()

	h.StockChangeProductQty().AddFields(map[string]models.FieldDefinition{
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
			String:  "Variant Count",
			Related: `ProductTmplId.ProductVariantCount`,
		},
		"NewQuantity": models.FloatField{
			String:  "New Quantity on Hand",
			Default: models.DefaultValue(1),
			//digits=dp.get_precision('Product Unit of Measure')
			Required: true,
			Help: "This quantity is expressed in the Default Unit of Measure" +
				"of the product.",
		},
		"LotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Lot/Serial Number",
			Filter:        q.ProductId().Equals(product_id),
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Required:      true,
			Filter:        q.Usage().Equals("internal"),
		},
	})
	h.StockChangeProductQty().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockChangeProductQtySet, fields interface{}) {
			//        res = super(ProductChangeQuantity, self).default_get(fields)
			//        if not res.get('product_id') and self.env.context.get('active_id') and self.env.context.get('active_model') == 'product.template' and self.env.context.get('active_id'):
			//            res['product_id'] = self.env['product.product'].search(
			//                [('product_tmpl_id', '=', self.env.context['active_id'])], limit=1).id
			//        elif not res.get('product_id') and self.env.context.get('active_id') and self.env.context.get('active_model') == 'product.product' and self.env.context.get('active_id'):
			//            res['product_id'] = self.env['product.product'].browse(
			//                self.env.context['active_id']).id
			//        if 'location_id' in fields and not res.get('location_id'):
			//            res['location_id'] = self.env.ref('stock.stock_location_stock').id
			//        return res
		})
	h.StockChangeProductQty().Methods().OnchangeLocationId().DeclareMethod(
		`OnchangeLocationId`,
		func(rs m.StockChangeProductQtySet) {
			//        if self.location_id and self.product_id:
			//            availability = self.product_id.with_context(
			//                compute_child=False)._product_available()
			//            self.new_quantity = availability[self.product_id.id]['qty_available']
		})
	h.StockChangeProductQty().Methods().OnchangeProductId().DeclareMethod(
		`OnchangeProductId`,
		func(rs m.StockChangeProductQtySet) {
			//        if self.product_id:
			//            self.product_tmpl_id = self.onchange_product_id_dict(self.product_id.id)[
			//                'product_tmpl_id']
		})
	h.StockChangeProductQty().Methods().PrepareInventoryLine().DeclareMethod(
		`PrepareInventoryLine`,
		func(rs m.StockChangeProductQtySet) {
			//        product = self.product_id.with_context(
			//            location=self.location_id.id, lot_id=self.lot_id.id)
			//        th_qty = product.qty_available
			//        res = {
			//            'product_qty': self.new_quantity,
			//            'location_id': self.location_id.id,
			//            'product_id': self.product_id.id,
			//            'product_uom_id': self.product_id.uom_id.id,
			//            'theoretical_qty': th_qty,
			//            'prod_lot_id': self.lot_id.id,
			//        }
			//        return res
		})
	h.StockChangeProductQty().Methods().OnchangeProductIdDict().DeclareMethod(
		`OnchangeProductIdDict`,
		func(rs m.StockChangeProductQtySet, product_id interface{}) {
			//        return {
			//            'product_tmpl_id': self.env['product.product'].browse(product_id).product_tmpl_id.id,
			//        }
		})
	h.StockChangeProductQty().Methods().Create().Extend(
		`Create`,
		func(rs m.StockChangeProductQtySet, values models.RecordData) {
			//        if values.get('product_id'):
			//            values.update(self.onchange_product_id_dict(values['product_id']))
			//        return super(ProductChangeQuantity, self).create(values)
		})
	h.StockChangeProductQty().Methods().CheckNewQuantity().DeclareMethod(
		`CheckNewQuantity`,
		func(rs m.StockChangeProductQtySet) {
			//        if any(wizard.new_quantity < 0 for wizard in self):
			//            raise UserError(_('Quantity cannot be negative.'))
		})
	h.StockChangeProductQty().Methods().ChangeProductQty().DeclareMethod(
		` Changes the Product Quantity by making a Physical Inventory. `,
		func(rs m.StockChangeProductQtySet) {
			//        Inventory = self.env['stock.inventory']
			//        for wizard in self:
			//            product = wizard.product_id.with_context(
			//                location=wizard.location_id.id, lot_id=wizard.lot_id.id)
			//            line_data = wizard._prepare_inventory_line()
			//
			//            if wizard.product_id.id and wizard.lot_id.id:
			//                inventory_filter = 'none'
			//            elif wizard.product_id.id:
			//                inventory_filter = 'product'
			//            else:
			//                inventory_filter = 'none'
			//            inventory = Inventory.create({
			//                'name': _('INV: %s') % tools.ustr(wizard.product_id.name),
			//                'filter': inventory_filter,
			//                'product_id': wizard.product_id.id,
			//                'location_id': wizard.location_id.id,
			//                'lot_id': wizard.lot_id.id,
			//                'line_ids': [(0, 0, line_data)],
			//            })
			//            inventory.action_done()
			//        return {'type': 'ir.actions.act_window_close'}
		})
}
