package either

import (
	"encoding/json"
	"fmt"
)

func (e *Either[L, R]) UnmarshalJSON(bytes []byte) error {
	l := new(L)
	lErr := json.Unmarshal(bytes, &l)
	if lErr == nil {
		*e = NewLeft[L, R](*l)
		return nil
	}

	r := new(R)
	rErr := json.Unmarshal(bytes, &r)
	if rErr == nil {
		*e = NewRight[L, R](*r)
		return nil
	}

	return fmt.Errorf("failed to unmarshal either to %T or %T", l, r)
}

func (e Either[L, R]) MarshalJSON() (res []byte, err error) {
	e.
		MapLeft(func(l L) {
			res, err = json.Marshal(l)
		}).
		MapRight(func(r R) {
			res, err = json.Marshal(r)
		})

	return res, err
}
