package dom

import (
	"errors"
	"strings"
)

// Errors.
var (
	ErrInvalidName = errors.New("invalid name")
)

type (
	// Name is an identifier of some entity in format "{collection}/{id}".
	Name struct {
		collection, id string
	}
	// UserName is a Name within "users" collection.
	UserName struct{ Name }
)

// Empty values.
//nolint:gochecknoglobals // Const.
var (
	NoName Name
	NoUser UserName
)

// NewName creates and returns Name with given collection and id.
func NewName(collection, id string) Name {
	switch {
	case collection == "":
		panic("require collection")
	case collection[len(collection)-1] == '/':
		panic("collection should not ends with /")
	case id == "":
		panic("require id")
	case id[0] == '/':
		panic("id should not begins with /")
	}
	return Name{
		collection: collection,
		id:         id,
	}
}

// ParseName returns ErrInvalidName if name doesn't belong to collection.
func ParseName(collection, name string) (*Name, error) {
	id := strings.TrimPrefix(name, collection+"/")
	if id == "" || id[0] == '/' {
		return nil, ErrInvalidName
	}
	n := NewName(collection, id)
	if n.String() != name {
		return nil, ErrInvalidName
	}
	return &n, nil
}

// String returns "{collection}/{id}" or empty string for NoName.
func (name Name) String() string {
	if name == NoName {
		return ""
	}
	return name.collection + "/" + name.id
}

// ID returns "{id}" part of the name.
func (name Name) ID() string { return name.id }

// NewUserName converts user's ID to UserName.
func NewUserName(userID string) UserName {
	return UserName{Name: NewName("users", userID)}
}

// ParseUserName returns ErrInvalidName if name doesn't belong to "users" collection.
func ParseUserName(name string) (*UserName, error) {
	n, err := ParseName("users", name)
	if err != nil {
		return nil, err
	}
	return &UserName{Name: *n}, nil
}
