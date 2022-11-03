package types

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type Uint32List []uint32

func (l Uint32List) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "[]", nil
	}

	strIDs := make([]string, len(l))
	for idx, id := range l {
		strIDs[idx] = strconv.Itoa(int(id))
	}

	return "[" + strings.Join(strIDs, ",") + "]", nil
}
