package node

type node interface {
	NewInstance()
	NewCB()
	CompileHeader()
}
