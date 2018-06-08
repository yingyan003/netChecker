package uuid

import (
	"fmt"
	"github.com/pborman/uuid"
	"testing"
)

func Test_UUID(*testing.T) {
	uuid := uuid.New()

	fmt.Println(uuid)
}
