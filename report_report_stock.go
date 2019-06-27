package stock

import (
	"github.com/hexya-erp/hexya/src/models"
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.ReportStockLinesDate().DeclareModel()

	h.ReportStockLinesDate().AddFields(map[string]models.FieldDefinition{
		"Id": models.IntegerField{
			String:   "Product Id",
			ReadOnly: true,
		},
		"ProductId": models.Many2OneField{
			RelationModel: h.ProductProduct(),
			String:        "Product",
			ReadOnly:      true,
			Index:         true,
		},
		"Date": models.DateTimeField{
			String:   "Date of latest Inventory",
			ReadOnly: true,
		},
		"MoveDate": models.DateTimeField{
			String:   "Date of latest Stock Move",
			ReadOnly: true,
		},
		"Active": models.BooleanField{
			String:   "Active",
			ReadOnly: true,
		},
	})
	h.ReportStockLinesDate().Methods().Init().DeclareMethod(
		`Init`,
		func(rs m.ReportStockLinesDateSet) {
			//        drop_view_if_exists(self._cr, 'report_stock_lines_date')
			//        self._cr.execute("""
			//            create or replace view report_stock_lines_date as (
			//                select
			//                p.id as id,
			//                p.id as product_id,
			//                max(s.date) as date,
			//                max(m.date) as move_date,
			//                p.active as active
			//            from
			//                product_product p
			//                    left join (
			//                        stock_inventory_line l
			//                        inner join stock_inventory s on (l.inventory_id=s.id and s.state = 'done')
			//                    ) on (p.id=l.product_id)
			//                    left join stock_move m on (m.product_id=p.id and m.state = 'done')
			//                group by p.id
			//            )""")
		})
}
