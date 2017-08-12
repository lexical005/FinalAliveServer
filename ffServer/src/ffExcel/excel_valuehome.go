package ffExcel

import "strings"

// 配置形式为server或client或server|client或直接为空
type valueHome struct {
	server, client bool
}

func (vh valueHome) String() string {
	if vh.server && vh.client {
		return "c|s"
	} else if vh.server {
		return "s"
	} else if vh.client {
		return "c"
	}
	return ""
}

func newValueHome(v string) *valueHome {
	v = strings.ToLower(v)
	server := strings.Index(string(v), "server") != -1
	client := strings.Index(string(v), "client") != -1
	if !server && !client && len(v) > 0 {
		return nil
	}
	return &valueHome{
		server: server,
		client: client,
	}
}
