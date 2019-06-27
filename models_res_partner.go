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

h.Partner().AddFields(map[string]models.FieldDefinition{
"PropertyStockCustomer": models.Many2OneField{
RelationModel: h.StockLocation(),
String: "Customer Location",
//company_dependent=True
Help: "This stock location will be used, instead of the default" + 
"one, as the destination location for goods you send to this partner",
},
"PropertyStockSupplier": models.Many2OneField{
RelationModel: h.StockLocation(),
String: "Vendor Location",
//company_dependent=True
Help: "This stock location will be used, instead of the default" + 
"one, as the source location for goods you receive from" + 
"the current partner",
},
"PickingWarn": models.SelectionField{
Selection: WARNING_MESSAGE,
String: "Stock Picking",
{"_fields":["id","ctx"],"ctx":null,"id":"WARNING_HELP","loc":{"end":{"column":59,"line":18},"start":{"column":47,"line":18}},"type":"Name"}
Default: models.DefaultValue("no-message"),
Required: true,
},
"PickingWarnMsg": models.TextField{
String: "Message for Stock Picking",
},
})
}