package main

import "fmt"

func onClientRequestDefault(reqClient string, ReqContent string, dictData map[string]string) error {
	if pveLevel, ok := dictData["pve_level"]; !ok {
		return fmt.Errorf("onSetupIAPvivo not contain pve_level")
	} else if pveExp, ok := dictData["pve_exp"]; !ok {
		return fmt.Errorf("onSetupIAPvivo not contain pve_exp")
	} else if ReqID, ok := dictData["ReqID"]; !ok {
		return fmt.Errorf("onSetupIAPvivo not contain ReqID")
	} else if ReqType, ok := dictData["ReqType"]; !ok {
		return fmt.Errorf("onSetupIAPvivo not contain ReqType")
	} else if ReqTime, ok := dictData["ReqTime"]; !ok {
		return fmt.Errorf("onSetupIAPvivo not contain ReqTime")
	} else if true {
		mysql.query(0, 200, nil, reqClient, pveLevel, pveExp, ReqID, ReqType, ReqTime, ReqContent)
		return nil
	}
	return nil
}
