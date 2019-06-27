package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/hexya/src/models/types"
	"github.com/hexya-erp/hexya/src/models/types/dates"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

func init() {
	h.StockScrap().DeclareModel()

	h.StockScrap().Methods().GetDefaultScrapLocationId().DeclareMethod(
		`GetDefaultScrapLocationId`,
		func(rs m.StockScrapSet) {
			//        return self.env['stock.location'].search([('scrap_location', '=', True), ('company_id', 'in', [self.env.user.company_id.id, False])], limit=1).id
		})
	h.StockScrap().Methods().GetDefaultLocationId().DeclareMethod(
		`GetDefaultLocationId`,
		func(rs m.StockScrapSet) {
			//        company_user = self.env.user.company_id
			//        warehouse = self.env['stock.warehouse'].search(
			//            [('company_id', '=', company_user.id)], limit=1)
			//        if warehouse:
			//            return warehouse.lot_stock_id.id
			//        return None
		})
	h.StockScrap().AddFields(map[string]models.FieldDefinition{
		"Name": models.CharField{
			String:   "Reference",
			Default:  func(env models.Environment) interface{} { return odoo._() },
			NoCopy:   true,
			ReadOnly: true,
			Required: true,
			//states={'done': [('readonly', True)]}
		},
		"Origin": models.CharField{
			String: "Source Document",
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Required:      true,
			//states={'done': [('readonly', True)]}
		},
		"ProductUomId": models.Many2OneField{
			RelationModel: h.ProductUom(),
			String:        "Unit of Measure",
			Required:      true,
			//states={'done': [('readonly', True)]}
		},
		"Tracking": models.SelectionField{
			Selection: "Product Tracking",
			ReadOnly:  true,
			Related:   `ProductId.Tracking`,
		},
		"LotId": models.Many2OneField{
			RelationModel: h.StockProductionLot(),
			String:        "Lot",
			//states={'done': [('readonly', True)]}
			Filter: q.ProductId().Equals(product_id),
		},
		"PackageId": models.Many2OneField{
			RelationModel: h.StockQuantPackage(),
			String:        "Package",
			//states={'done': [('readonly', True)]}
		},
		"OwnerId": models.Many2OneField{
			RelationModel: h.Partner(),
			String:        "Owner",
			//states={'done': [('readonly', True)]}
		},
		"MoveId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Scrap Move",
			ReadOnly:      true,
		},
		"PickingId": models.Many2OneField{
			RelationModel: h.StockPicking(),
			String:        "Picking",
			//states={'done': [('readonly', True)]}
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Location",
			Filter:        q.Usage().Equals("internal"),
			Required:      true,
			//states={'done': [('readonly', True)]}
			Default: models.DefaultValue(_get_default_location_id),
		},
		"ScrapLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Scrap Location",
			Default:       models.DefaultValue(_get_default_scrap_location_id),
			Filter:        q.ScrapLocation().Equals(True),
			//states={'done': [('readonly', True)]}
		},
		"ScrapQty": models.FloatField{
			String:   "Quantity",
			Default:  models.DefaultValue(1),
			Required: true,
			//states={'done': [('readonly', True)]}
		},
		"State": models.SelectionField{
			Selection: types.Selection{
				"draft": "Draft",
				"done":  "Done",
			},
			String:  "Status",
			Default: models.DefaultValue("draft"),
		},
		"DateExpected": models.DateTimeField{
			String:  "Expected Date",
			Default: func(env models.Environment) interface{} { return dates.Now() },
		},
	})
	h.StockScrap().Methods().OnchangePickingId().DeclareMethod(
		`OnchangePickingId`,
		func(rs m.StockScrapSet) {
			//        if self.picking_id:
			//            self.location_id = (
			//                self.picking_id.state == 'done') and self.picking_id.location_dest_id.id or self.picking_id.location_id.id
		})
	h.StockScrap().Methods().OnchangeProductId().DeclareMethod(
		`OnchangeProductId`,
		func(rs m.StockScrapSet) {
			//        if self.product_id:
			//            self.product_uom_id = self.product_id.uom_id.id
		})
	h.StockScrap().Methods().Create().Extend(
		`Create`,
		func(rs m.StockScrapSet, vals models.RecordData) {
			//        if 'name' not in vals or vals['name'] == _('New'):
			//            vals['name'] = self.env['ir.sequence'].next_by_code(
			//                'stock.scrap') or _('New')
			//        scrap = super(StockScrap, self).create(vals)
			//        scrap.do_scrap()
			//        return scrap
		})
	h.StockScrap().Methods().Unlink().Extend(
		`Unlink`,
		func(rs m.StockScrapSet) {
			//        if 'done' in self.mapped('state'):
			//            raise UserError(_('You cannot delete a scrap which is done.'))
			//        return super(StockScrap, self).unlink()
		})
	h.StockScrap().Methods().GetOriginMoves().DeclareMethod(
		`GetOriginMoves`,
		func(rs m.StockScrapSet) {
			//        return self.picking_id and self.picking_id.move_lines.filtered(lambda x: x.product_id == self.product_id)
		})
	h.StockScrap().Methods().DoScrap().DeclareMethod(
		`DoScrap`,
		func(rs m.StockScrapSet) {
			//        for scrap in self:
			//            moves = scrap._get_origin_moves() or self.env['stock.move']
			//            move = self.env['stock.move'].create(scrap._prepare_move_values())
			//            if move.product_id.type == 'product':
			//                quants = self.env['stock.quant'].quants_get_preferred_domain(
			//                    move.product_qty, move,
			//                    domain=[
			//                        ('qty', '>', 0),
			//                        ('lot_id', '=', self.lot_id.id),
			//                        ('package_id', '=', self.package_id.id)],
			//                    preferred_domain_list=scrap._get_preferred_domain())
			//                if any([not x[0] for x in quants]):
			//                    raise UserError(
			//                        _('You cannot scrap a move without having available stock for %s. You can correct it with an inventory adjustment.') % move.product_id.name)
			//                self.env['stock.quant'].quants_reserve(quants, move)
			//            move.action_done()
			//            scrap.write({'move_id': move.id, 'state': 'done'})
			//            moves.recalculate_move_state()
			//        return True
		})
	h.StockScrap().Methods().PrepareMoveValues().DeclareMethod(
		`PrepareMoveValues`,
		func(rs m.StockScrapSet) {
			//        self.ensure_one()
			//        return {
			//            'name': self.name,
			//            'origin': self.origin or self.picking_id.name,
			//            'product_id': self.product_id.id,
			//            'product_uom': self.product_uom_id.id,
			//            'product_uom_qty': self.scrap_qty,
			//            'location_id': self.location_id.id,
			//            'scrapped': True,
			//            'location_dest_id': self.scrap_location_id.id,
			//            'restrict_lot_id': self.lot_id.id,
			//            'restrict_partner_id': self.owner_id.id,
			//            'picking_id': self.picking_id.id
			//        }
		})
	h.StockScrap().Methods().GetPreferredDomain().DeclareMethod(
		`GetPreferredDomain`,
		func(rs m.StockScrapSet) {
			//        if not self.picking_id:
			//            return []
			//        if self.picking_id.state == 'done':
			//            preferred_domain = [('history_ids', 'in', self.picking_id.move_lines.filtered(
			//                lambda x: x.state == 'done')).ids]
			//            preferred_domain2 = [('history_ids', 'not in', self.picking_id.move_lines.filtered(
			//                lambda x: x.state == 'done')).ids]
			//            return [preferred_domain, preferred_domain2]
			//        else:
			//            preferred_domain = [
			//                ('reservation_id', 'in', self.picking_id.move_lines.ids)]
			//            preferred_domain2 = [('reservation_id', '=', False)]
			//            preferred_domain3 = [
			//                '&', ('reservation_id', 'not in', self.picking_id.move_lines.ids), ('reservation_id', '!=', False)]
			//            return [preferred_domain, preferred_domain2, preferred_domain3]
		})
	h.StockScrap().Methods().ActionGetStockPicking().DeclareMethod(
		`ActionGetStockPicking`,
		func(rs m.StockScrapSet) {
			//        action = self.env.ref('stock.action_picking_tree_all').read([])[0]
			//        action['domain'] = [('id', '=', self.picking_id.id)]
			//        return action
		})
	h.StockScrap().Methods().ActionGetStockMove().DeclareMethod(
		`ActionGetStockMove`,
		func(rs m.StockScrapSet) {
			//        action = self.env.ref('stock.stock_move_action').read([])[0]
			//        action['domain'] = [('id', '=', self.move_id.id)]
			//        return action
		})
	h.StockScrap().Methods().ActionDone().DeclareMethod(
		`ActionDone`,
		func(rs m.StockScrapSet) {
			//        return {'type': 'ir.actions.act_window_close'}
		})
}
