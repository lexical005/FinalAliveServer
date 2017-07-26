package main

import (
	"ffCommon/util"

	"encoding/json"
	"errors"
	"strconv"
)

func loadServerConfig(cfgFilePath string) (*tServerConfig, error) {

	buff, err := util.ReadFile(cfgFilePath)
	if err != nil {
		return nil, err
	}

	var dataOri map[string]interface{}
	if err := json.Unmarshal(buff, &dataOri); err != nil {
		return nil, errors.New("invalid cfg_file_path:" + cfgFilePath + "\n" + err.Error())
	}
	result := &tServerConfig{}

	result.devChannel, _ = dataOri["DEV_CHANNEL"].(string)

	result.listenIPPort, _ = dataOri["SERVE_IP_PORT"].(string)
	result.outerIPPort, _ = dataOri["NTF_CLIENT_IP_PORT"].(string)

	maxConnectionCount, _ := dataOri["MAX_CONNECTION_COUNT"].(string)
	result.connectionLimit, _ = strconv.Atoi(maxConnectionCount)

	mapChannelConfig, _ := dataOri["CHANNEL"].(map[string]interface{})
	result.channelConfig = make(map[string]map[string]string, len(mapChannelConfig))

	for channelName, channelConfig := range mapChannelConfig {
		mapChannelConfigDetail, _ := channelConfig.(map[string]interface{})

		fullDownURL, _ := mapChannelConfigDetail["FULL_DOWN_URL"].(string)
		selectServerIP, _ := mapChannelConfigDetail["SELECT_SERVER_IP"].(string)

		result.channelConfig[channelName] = make(map[string]string, len(mapChannelConfigDetail))
		result.channelConfig[channelName]["FULL_DOWN_URL"] = fullDownURL
		result.channelConfig[channelName]["SELECT_SERVER_IP"] = selectServerIP
	}

	return result, nil
}
