package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
	"github.com/hexya-erp/pool/q"
)

func init() {
	h.StockReturnPickingLine().DeclareModel()

	h.StockReturnPickingLine().AddFields(map[string]models.FieldDefinition{
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			Required:      true,
		},
		"Quantity": models.FloatField{
			String: "Quantity",
			//digits=dp.get_precision('Product Unit of Measure')
			Required: true,
		},
		"WizardId": models.Many2OneField{
			RelationModel: h.StockReturnPicking(),
			String:        "Wizard",
		},
		"MoveId": models.Many2OneField{
			RelationModel: h.StockMove(),
			String:        "Move",
		},
	})
	h.StockReturnPicking().DeclareModel()

	h.StockReturnPicking().AddFields(map[string]models.FieldDefinition{
		"ProductReturnMoves": models.One2ManyField{
			RelationModel: h.StockReturnPickingLine(),
			ReverseFK:     "",
			String:        "Moves",
		},
		"MoveDestExists": models.BooleanField{
			String:   "Chained Move Exists",
			ReadOnly: true,
		},
		"OriginalLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
		},
		"ParentLocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
		},
		"LocationId": models.Many2OneField{
			RelationModel: h.StockLocation(),
			String:        "Return Location",
			Filter:        q.Id().Equals(original_location_id).Or().ReturnLocation().Equals(True).And().Id().ChildOf(parent_location_id),
		},
	})
	h.StockReturnPicking().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockReturnPickingSet, fields interface{}) {
			//        if len(self.env.context.get('active_ids', list())) > 1:
			//            raise UserError("You may only return one picking at a time!")
			//        res = super(ReturnPicking, self).default_get(fields)
			//        Quant = self.env['stock.quant']
			//        move_dest_exists = False
			//        product_return_moves = []
			//        picking = self.env['stock.picking'].browse(
			//            self.env.context.get('active_id'))
			//        if picking:
			//            if picking.state != 'done':
			//                raise UserError(_("You may only return Done pickings"))
			//            for move in picking.move_lines:
			//                if move.scrapped:
			//                    continue
			//                if move.move_dest_id:
			//                    move_dest_exists = True
			//                # Sum the quants in that location that can be returned (they should have been moved by the moves that were included in the returned picking)
			//                quantity = sum(quant.qty for quant in Quant.search([
			//                    ('history_ids', 'in', move.id),
			//                    ('qty', '>', 0.0), ('location_id',
			//                                        'child_of', move.location_dest_id.id)
			//                ]).filtered(
			//                    lambda quant: not quant.reservation_id or quant.reservation_id.origin_returned_move_id != move)
			//                )
			//                quantity = move.product_id.uom_id._compute_quantity(
			//                    quantity, move.product_uom)
			//                product_return_moves.append(
			//                    (0, 0, {'product_id': move.product_id.id, 'quantity': quantity, 'move_id': move.id}))
			//
			//            if not product_return_moves:
			//                raise UserError(
			//                    _("No products to return (only lines in Done state and not fully returned yet can be returned)!"))
			//            if 'product_return_moves' in fields:
			//                res.update({'product_return_moves': product_return_moves})
			//            if 'move_dest_exists' in fields:
			//                res.update({'move_dest_exists': move_dest_exists})
			//            if 'parent_location_id' in fields and picking.location_id.usage == 'internal':
			//                res.update(
			//                    {'parent_location_id': picking.picking_type_id.warehouse_id and picking.picking_type_id.warehouse_id.view_location_id.id or picking.location_id.location_id.id})
			//            if 'original_location_id' in fields:
			//                res.update({'original_location_id': picking.location_id.id})
			//            if 'location_id' in fields:
			//                location_id = picking.location_id.id
			//                if picking.picking_type_id.return_picking_type_id.default_location_dest_id.return_location:
			//                    location_id = picking.picking_type_id.return_picking_type_id.default_location_dest_id.id
			//                res['location_id'] = location_id
			//        return res
		})
	h.StockReturnPicking().Methods().CreateReturns().DeclareMethod(
		`CreateReturns`,
		func(rs m.StockReturnPickingSet) {
			//        picking = self.env['stock.picking'].browse(
			//            self.env.context['active_id'])
			//        return_moves = self.product_return_moves.mapped('move_id')
			//        unreserve_moves = self.env['stock.move']
			//        for move in return_moves:
			//            to_check_moves = self.env['stock.move'] | move.move_dest_id
			//            while to_check_moves:
			//                current_move = to_check_moves[-1]
			//                to_check_moves = to_check_moves[:-1]
			//                if current_move.state not in ('done', 'cancel') and current_move.reserved_quant_ids:
			//                    unreserve_moves |= current_move
			//                split_move_ids = self.env['stock.move'].search(
			//                    [('split_from', '=', current_move.id)])
			//                to_check_moves |= split_move_ids
			//        if unreserve_moves:
			//            unreserve_moves.do_unreserve()
			//            # break the link between moves in order to be able to fix them later if needed
			//            unreserve_moves.write({'move_orig_ids': False})
			//        picking_type_id = picking.picking_type_id.return_picking_type_id.id or picking.picking_type_id.id
			//        new_picking = picking.copy({
			//            'move_lines': [],
			//            'picking_type_id': picking_type_id,
			//            'state': 'draft',
			//            'origin': picking.name,
			//            'location_id': picking.location_dest_id.id,
			//            'location_dest_id': self.location_id.id})
			//        new_picking.message_post_with_view('mail.message_origin_link',
			//                                           values={'self': new_picking,
			//                                                   'origin': picking},
			//                                           subtype_id=self.env.ref('mail.mt_note').id)
			//        returned_lines = 0
			//        for return_line in self.product_return_moves:
			//            if not return_line.move_id:
			//                raise UserError(
			//                    _("You have manually created product lines, please delete them to proceed"))
			//            new_qty = return_line.quantity
			//            if new_qty:
			//                # The return of a return should be linked with the original's destination move if it was not cancelled
			//                if return_line.move_id.origin_returned_move_id.move_dest_id.id and return_line.move_id.origin_returned_move_id.move_dest_id.state != 'cancel':
			//                    move_dest_id = return_line.move_id.origin_returned_move_id.move_dest_id.id
			//                else:
			//                    move_dest_id = False
			//
			//                returned_lines += 1
			//                return_line.move_id.copy({
			//                    'product_id': return_line.product_id.id,
			//                    'product_uom_qty': new_qty,
			//                    'picking_id': new_picking.id,
			//                    'state': 'draft',
			//                    'location_id': return_line.move_id.location_dest_id.id,
			//                    'location_dest_id': self.location_id.id or return_line.move_id.location_id.id,
			//                    'picking_type_id': picking_type_id,
			//                    'warehouse_id': picking.picking_type_id.warehouse_id.id,
			//                    'origin_returned_move_id': return_line.move_id.id,
			//                    'procure_method': 'make_to_stock',
			//                    'move_dest_id': move_dest_id,
			//                })
			//        if not returned_lines:
			//            raise UserError(
			//                _("Please specify at least one non-zero quantity."))
			//        new_picking.action_confirm()
			//        new_picking.action_assign()
			//        return new_picking.id, picking_type_id
		})
	h.StockReturnPicking().Methods().CreateReturns().DeclareMethod(
		`CreateReturns`,
		func(rs m.StockReturnPickingSet) {
			//        for wizard in self:
			//            new_picking_id, pick_type_id = wizard._create_returns()
			//        ctx = dict(self.env.context)
			//        ctx.update({
			//            'search_default_picking_type_id': pick_type_id,
			//            'search_default_draft': False,
			//            'search_default_assigned': False,
			//            'search_default_confirmed': False,
			//            'search_default_ready': False,
			//            'search_default_late': False,
			//            'search_default_available': False,
			//        })
			//        return {
			//            'name': _('Returned Picking'),
			//            'view_type': 'form',
			//            'view_mode': 'form,tree,calendar',
			//            'res_model': 'stock.picking',
			//            'res_id': new_picking_id,
			//            'type': 'ir.actions.act_window',
			//            'context': ctx,
			//        }
		})
}
