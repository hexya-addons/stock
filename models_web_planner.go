package stock

import (
	"github.com/hexya-erp/pool/h"
)

func init() {
	h.WebPlanner().DeclareModel()

	h.WebPlanner().Methods().GetPlannerApplication().DeclareMethod(
		`GetPlannerApplication`,
		func(rs m.WebPlannerSet) {
			//        planner = super(PlannerInventory, self)._get_planner_application()
			//        planner.append(['planner_inventory', 'Inventory Planner'])
			//        return planner
		})
}
