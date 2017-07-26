package main

// toml: https://github.com/toml-lang/toml

import (
	"ffCommon/log/log"
	"fmt"
	"time"

	"github.com/lexical005/toml"
)

var tomlString = `
# This is a TOML document. Boom.

title = "TOML Example"

[owner]
name = "Tom Preston-Werner"
organization = "GitHub"
bio = "GitHub Cofounder & CEO\nLikes tater tots and beer."
dob = 1979-05-27T07:32:00Z # First class dates? Why not?

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[ServersA]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [ServersA.1]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [ServersA.2]
  ip = "10.0.0.2"
  dc = "eqdc10"

[ServersB]

  # You can indent as you please. Tabs or spaces. TOML don't care.
  [ServersB.alpha]
  ip = "10.0.0.1"
  dc = "eqdc10"

  [ServersB.beta]
  ip = "10.0.0.2"
  dc = "eqdc10"

[clients]
data = [ ["gamma", "delta"], [1, 2] ] # just an update to make sure parsers support it

# Line breaks are OK when inside arrays
hosts = [
  "alpha",
  "omega"
]
`

type tomlConfig struct {
	Title    string
	Owner    ownerInfo
	DB       database `toml:"database"`
	ServersA map[int]server
	ServersB map[string]server
	Clients  clients
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
	Bio  string
	DOB  time.Time
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `toml:"connection_max"`
	Enabled bool
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Data  [][]interface{}
	Hosts []string
}

func main() {
	config := &tomlConfig{}
	err := toml.Unmarshal([]byte(tomlString), config)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.RunLogger.Println(config)
}
