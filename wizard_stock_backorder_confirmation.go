package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockBackorderConfirmation().DeclareModel()

	h.StockBackorderConfirmation().AddFields(map[string]models.FieldDefinition{
		"PickId": models.Many2OneField{
			RelationModel: h.StockPicking(),
		},
	})
	h.StockBackorderConfirmation().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockBackorderConfirmationSet, fields interface{}) {
			//        res = super(StockBackorderConfirmation, self).default_get(fields)
			//        if 'pick_id' in fields and self._context.get('active_id') and not res.get('pick_id'):
			//            res = {'pick_id': self._context['active_id']}
			//        return res
		})
	h.StockBackorderConfirmation().Methods().Process().DeclareMethod(
		`Process`,
		func(rs m.StockBackorderConfirmationSet, cancel_backorder interface{}) {
			//        operations_to_delete = self.pick_id.pack_operation_ids.filtered(
			//            lambda o: o.qty_done <= 0)
			//        for pack in self.pick_id.pack_operation_ids - operations_to_delete:
			//            pack.product_qty = pack.qty_done
			//        operations_to_delete.unlink()
			//        self.pick_id.do_transfer()
			//        if cancel_backorder:
			//            backorder_pick = self.env['stock.picking'].search(
			//                [('backorder_id', '=', self.pick_id.id)])
			//            backorder_pick.action_cancel()
			//            self.pick_id.message_post(
			//                body=_("Back order <em>%s</em> <b>cancelled</b>.") % (backorder_pick.name))
		})
	h.StockBackorderConfirmation().Methods().Process().DeclareMethod(
		`Process`,
		func(rs m.StockBackorderConfirmationSet) {
			//        self._process()
		})
	h.StockBackorderConfirmation().Methods().ProcessCancelBackorder().DeclareMethod(
		`ProcessCancelBackorder`,
		func(rs m.StockBackorderConfirmationSet) {
			//        self._process(cancel_backorder=True)
		})
}
