package zeq

import (
	"github.com/GoelandProver/Goeland/global"
)

func Enable() {
	global.PrintInfo("ZEQ", "ZEQ plugin enabled")
	global.SetZeq(true)
	// Ici PrintInfo printera 9223372036.854776s, juste ici, idk why

}
