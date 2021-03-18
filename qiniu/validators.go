package qiniu

import "fmt"

func validateRegionID(v interface{}, attributeName string) (warns []string, errs []error) {
	regionId := v.(string)
	switch regionId {
	case "z0", "z1", "z2", "na0", "as0":
		return
	default:
		errs = append(errs, fmt.Errorf("%q must be one of 'z0', 'z1', 'z2', 'na0' or 'as0'", attributeName))
		return
	}
}
