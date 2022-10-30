package model

// 模拟枚举行为，设计读写等权限
type read struct {
	id   int8
	name string
}

func (it read) NewRead() read {
	it.id = 1
	it.name = "read"
	return it
}

func (it read) GetName() string {
	return it.name
}

type write struct {
	id   int8
	name string
}

func (it write) NewWrite() write {
	it.id = 2
	it.name = "write"
	return it
}

func (it write) GetName() string {
	return it.name
}
