package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.StockImmediateTransfer().DeclareModel()

	h.StockImmediateTransfer().AddFields(map[string]models.FieldDefinition{
		"PickId": models.Many2OneField{
			RelationModel: h.StockPicking(),
		},
	})
	h.StockImmediateTransfer().Methods().DefaultGet().Extend(
		`DefaultGet`,
		func(rs m.StockImmediateTransferSet, fields interface{}) {
			//        res = super(StockImmediateTransfer, self).default_get(fields)
			//        if not res.get('pick_id') and self._context.get('active_id'):
			//            res['pick_id'] = self._context['active_id']
			//        return res
		})
	h.StockImmediateTransfer().Methods().Process().DeclareMethod(
		`Process`,
		func(rs m.StockImmediateTransferSet) {
			//        self.ensure_one()
			//        if self.pick_id.state == 'draft':
			//            self.pick_id.action_confirm()
			//            if self.pick_id.state != 'assigned':
			//                self.pick_id.action_assign()
			//                if self.pick_id.state != 'assigned':
			//                    raise UserError(
			//                        _("Could not reserve all requested products. Please use the \'Mark as Todo\' button to handle the reservation manually."))
			//        for pack in self.pick_id.pack_operation_ids:
			//            if pack.product_qty > 0:
			//                pack.write({'qty_done': pack.product_qty})
			//            else:
			//                pack.unlink()
			//        return self.pick_id.do_transfer()
		})
}
