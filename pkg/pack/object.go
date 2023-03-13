package pack

//go:generate stringer -type=ObjectType -trimprefix=Obj

type ObjectType int8

const (
	ObjNone ObjectType = iota
	ObjCommit
	ObjTree
	ObjBlob
	ObjTag
	ObjFake
	ObjOfsDelta
	ObjRefDelta
	ObjAny
	ObjMax
)

type Object struct {
	size   uint64
	offset uint64
	_type  ObjectType
}
