package web

import (
	"fmt"
	"testing"
)

func TestNil(t *testing.T) {
	testTypeAssert(nil)
}

func testTypeAssert(c any) {
	claims := c.(*UserClaims)
	fmt.Println(claims.Uid)
}
