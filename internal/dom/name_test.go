package dom_test

import (
	"testing"

	"github.com/powerman/check"

	"github.com/powerman/go-monolith-example/internal/dom"
)

func TestEmpty(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	t.Zero(dom.NoName)
	t.Zero(dom.NoName.String())
	t.Zero(dom.NoName.ID())

	t.Zero(dom.NoUser)
	t.Zero(dom.NoUser.String())
	t.Zero(dom.NoUser.ID())
}

func TestName(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	tests := []struct {
		collection string
		id         string
		want       string
		wantPanic  string
	}{
		{"", "", "", `require collection`},
		{"/", "", "", `collection should not ends with /`},
		{"a", "", "", `require id`},
		{"a", "/", "", `id should not begins with /`},
		{"a", "b", "a/b", ``},
		{"/a", "b", "/a/b", ``},
		{"a", "b/", "a/b/", ``},
		{"/a", "b/", "/a/b/", ``},
		{"a/b", "c/d", "a/b/c/d", ``},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			if tc.wantPanic != "" {
				t.PanicMatch(func() { dom.NewName(tc.collection, tc.id) }, tc.wantPanic)
			} else {
				res := dom.NewName(tc.collection, tc.id)
				t.Equal(res.String(), tc.want)
				t.Equal(res.ID(), tc.id)
				res2, err := dom.ParseName(tc.collection, res.String())
				t.Nil(err)
				t.DeepEqual(&res, res2)
			}
		})
	}
}

func TestParseName(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	t.PanicMatch(func() { dom.ParseName("", "a/b") }, `require collection`)
	t.PanicMatch(func() { dom.ParseName("a/", "a/b") }, `collection should not ends with /`)

	var (
		name1 = dom.NewName("a", "b")
		name2 = dom.NewName("a", "b/")
		name3 = dom.NewName("a", "b/c")
	)
	tests := []struct {
		given   string
		want    *dom.Name
		wantErr error
	}{
		{"", nil, dom.ErrInvalidName},
		{"x", nil, dom.ErrInvalidName},
		{"x/y", nil, dom.ErrInvalidName},
		{"a", nil, dom.ErrInvalidName},
		{"a/", nil, dom.ErrInvalidName},
		{"/a/b", nil, dom.ErrInvalidName},
		{"a//b", nil, dom.ErrInvalidName},
		{"a/b", &name1, nil},
		{"a/b/", &name2, nil},
		{"a/b/c", &name3, nil},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(tt *testing.T) {
			t := check.T(tt)
			res, err := dom.ParseName("a", tc.given)
			t.Err(err, tc.wantErr)
			t.DeepEqual(res, tc.want)
		})
	}
}

func TestUserName(tt *testing.T) {
	t := check.T(tt)
	t.Parallel()

	t.PanicMatch(func() { dom.NewUserName("") }, `require id`)

	name1 := dom.NewName("users", "1")
	user1 := dom.NewUserName("1")
	t.Equal(user1.String(), name1.String())
	t.Equal(user1.ID(), name1.ID())

	res, err := dom.ParseUserName("")
	t.Err(err, dom.ErrInvalidName)
	t.Nil(res)

	res, err = dom.ParseUserName("users/1")
	t.Nil(err)
	t.DeepEqual(res, &user1)
}
