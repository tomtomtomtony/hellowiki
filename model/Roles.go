package model

import "container/list"

// 所有访问网站的，都是用户
// 任何用户都有角色
type registerUser struct {
	id          int8
	name        string
	permissions list.List
}

type visitor struct {
	id          int8
	name        string
	permissions list.List
}

type admin struct {
	id          int8
	name        string
	permissions list.List
}
