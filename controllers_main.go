package stock

import (
	"net/http"

	"github.com/hexya-erp/hexya/src/controllers"
)

//import logging
//_logger = logging.getLogger(__name__)
func init() {
	root := controllers.Registry
	var ok bool
	var stock *controllers.Group
	stock, ok = root.GetGroup("/stock")
	if !ok {
		stock = root.AddGroup("/stock")
	}
	var barcode *controllers.Group
	barcode, ok = stock.GetGroup("/barcode")
	if !ok {
		barcode = stock.AddGroup("/barcode")
	}
	if barcode.HasController(http.MethodGet, "/") {
		barcode.ExtendController(http.MethodPost, "/", A)
	} else {
		barcode.AddController(http.MethodPost, "/", A)
	}
}
func A(self interface{}, debug interface{}) {
	//        if not request.session.uid:
	//            return http.local_redirect('/web/login?redirect=/stock/barcode/')
	//        return request.render('stock.barcode_index')
}
