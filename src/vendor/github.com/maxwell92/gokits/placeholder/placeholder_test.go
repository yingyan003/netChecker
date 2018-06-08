package placeholder

import (
	"fmt"
	"testing"
)

func Test_NewPlaceHolder(*testing.T) {
	replacer := "rbd create <image> -s <size> -p <pool>"

	ph := NewPlaceHolder(replacer)

	fmt.Println(ph.Replace("<image>", "rbd", "<size>", "1024", "<pool>", "rbd"))

}

func Test_Replace_RBDCREATE(*testing.T) {
	fs := "mkfs.<fs> <path>"

	ph := NewPlaceHolder(fs)

	fmt.Println(ph.Replace("<fs>", "ext4", "<path>", "/dev/rbd0"))
}
