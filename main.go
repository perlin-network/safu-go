package main

import (
	"github.com/perlin-network/safu-go/database"
	"log"
	"os"
)

// test address 0xDdd4E8279F3D5CEF259869F9866fC26817727aEA
// test address 0x14450b13B03D97B686A5eC40671D12c0963fd9bF

func main() {
	store := database.NewTieDotStore("/tmp/safu-level")

	defer os.RemoveAll("/tmp/safu-level")

	err := store.BFSEthAccount("0x4a966d2Ad06F980cD7f8fDc4c4360641aB2C9852", "4EIR7V4K5QBWDUGJKHFK4BGZ6HWD1NIFT1")
	if err != nil {
		log.Panicf("BFS error: %s", err)
	}
}
