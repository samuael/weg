package permission

import (
	// "net/http"
	// "strings"

	// "github.com/Projects/ScientificNRS/internal/pkg/entity"
)

type permission struct {
	roles   []string
	methods []string
}
type authority map[string]permission

var authorities = authority{
	"/admin/controll/": permission{
		roles:   []string{"SUPERADMIN", "SECRETARY"},
		methods: []string{"GET", "POST"},
	},
	"/admin/registration/": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"GET"},
	},

	"/admin/teacher/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/admin/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/admin/fieldman/new": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
	"/api/vehicles/": permission{
		roles:   []string{"SUPERADMIN", "SECRETARY"},
		methods: []string{"POST", "GET"},
	},
	"/admin/trainer/new/": permission{
		roles:   []string{"SUPERADMIN"},
		methods: []string{"POST"},
	},
}

// // HasPermission checks if a given role has permission to access a given route for a given method
// func HasPermission(path string, role string, method string) bool {
// 	perm := authorities[path]
// 	checkedRole := checkRole(role, perm.roles)
// 	checkedMethod := checkMethod(method, perm.methods)
// 	if !checkedRole || !checkedMethod {
// 		return false
// 	}
// 	return true
// }

// func checkRole(role string, roles []string) bool {
// 	for _, r := range roles {
// 		if strings.ToUpper(r) == strings.ToUpper(role) {
// 			return true
// 		}
// 	}
// 	return false
// }
// func checkMethod(method string, methods []string) bool {
// 	for _, m := range methods {
// 		if strings.ToUpper(m) == strings.ToUpper(method) {
// 			return true
// 		}
// 	}
// 	return false
// }
