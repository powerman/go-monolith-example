package app

import (
	"fmt"
	"reflect"

	"github.com/powerman/go-monolith-example/internal/dom"
)

// MatchUser provides gomock.Matcher for matching User and *User:
//   - empty Name match any Name
//   - PassHash is ignored
//   - CreateTime is ignored
// On successful match Name/PassHash/CreateTime of matching User will be
// set to values from matched User (which means empty Name will match any
// Name only once, but next time it'll match only previous Name).
type MatchUser struct{ *User }

// String implements gomock.Matcher.
func (m MatchUser) String() string { return fmt.Sprint(m.User) }

// Matches implements gomock.Matcher.
func (m MatchUser) Matches(x interface{}) bool {
	u, ok := x.(User)
	if !ok {
		uPtr, ok := x.(*User)
		if !ok {
			return false
		}
		u = *uPtr
	}
	v := *m.User
	if v.Name.ID() == "" && u.Name.ID() != "" {
		v.Name = dom.NewUserName(u.Name.ID())
	}
	v.PassHash = u.PassHash
	v.CreateTime = u.CreateTime
	if !reflect.DeepEqual(v, u) {
		return false
	}
	*m.User = v
	return true
}
