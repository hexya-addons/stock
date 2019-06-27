package stock

	import (
		"net/http"

		"github.com/hexya-erp/hexya/src/controllers"
		"github.com/hexya-erp/hexya/src/models"
		"github.com/hexya-erp/hexya/src/models/types"
		"github.com/hexya-erp/hexya/src/models/types/dates"
		"github.com/hexya-erp/pool/h"
		"github.com/hexya-erp/pool/q"
	)
	
func init() {
h.StockProductionLot().DeclareModel()
h.StockProductionLot().AddSQLConstraint("name_ref_uniq", "unique (name, product_id)", "The combination of serial number and product must be unique !")



h.StockProductionLot().AddFields(map[string]models.FieldDefinition{
"Name": models.CharField{
String: "Lot/Serial Number",
Default: func (env models.Environment) interface{} { return env["ir.sequence"].next_by_code() },
Required: true,
Help: "Unique Lot/Serial Number",
},
"Ref": models.CharField{
String: "Internal Reference",
Help: "Internal reference number in case it differs from the manufacturer's" + 
"lot/serial number",
},
"ProductId": models.Many2OneField{
RelationModel: h.ProductProduct(),
String: "Product",
Filter: q.Type().In(%!s(<nil>)),
Required: true,
},
"ProductUomId": models.Many2OneField{
RelationModel: h.ProductUom(),
String: "Unit of Measure",
Related: `ProductId.UomId`,
Stored: true,
},
"QuantIds": models.One2ManyField{
RelationModel: h.StockQuant(),
ReverseFK: "",
String: "Quants",
ReadOnly: true,
},
"ProductQty": models.FloatField{
String: "Quantity",
Compute: h.StockProductionLot().Methods().ProductQty(),
},

})
h.StockProductionLot().Fields().CreateDate().setString( "Creation Date",)
h.StockProductionLot().Methods().Create().Extend(
`Create`,
func(rs m.StockProductionLotSet, vals models.RecordData)  {
//        pack_id = self.env.context.get('active_pack_operation', False)
//        if pack_id:
//            pack = self.env['stock.pack.operation'].browse(pack_id)
//            if pack.picking_id and not pack.picking_id.picking_type_id.use_create_lots:
//                raise UserError(
//                    _("You are not allowed to create a lot for this picking type"))
//        return super(ProductionLot, self).create(vals)
})
h.StockProductionLot().Methods().ProductQty().DeclareMethod(
`ProductQty`,
func(rs h.StockProductionLotSet) h.StockProductionLotData {
//        self.product_qty = sum(self.quant_ids.mapped('qty'))
})
h.StockProductionLot().Methods().ActionTraceability().DeclareMethod(
`ActionTraceability`,
func(rs m.StockProductionLotSet)  {
//        move_ids = self.mapped('quant_ids').mapped('history_ids').ids
//        if not move_ids:
//            return False
//        return {
//            'domain': [('id', 'in', move_ids)],
//            'name': _('Traceability'),
//            'view_mode': 'tree,form',
//            'view_type': 'form',
//            'context': {'tree_view_ref': 'stock.view_move_tree'},
//            'res_model': 'stock.move',
//            'type': 'ir.actions.act_window'}
})
}