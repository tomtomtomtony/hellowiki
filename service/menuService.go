package service

import (
	"hellowiki/api/v1/menu/vo"
	"hellowiki/model"
)

type MenuNode struct {
	MenuID   uint       `json:"menuID"`
	NodeType string     `json:"nodeType"`
	Children []MenuNode `json:"children"`
	Name     string     `json:"name"`
	ParentId uint       `json:"parentId"`
}

func GetDirectChildren(categoryId uint) []vo.MenuVO {
	resRaw := model.GetAllDirectMenuChildren(categoryId)
	var res = make([]vo.MenuVO, 0, len(resRaw))
	for _, item := range resRaw {
		res = append(res, data2Vo(item))
	}
	return res
}

func GetAllTopCategory() []vo.TopMenu {
	resRaw := model.GetTopLevelMenu()
	var res = make([]vo.TopMenu, 0, len(resRaw))
	for _, item := range resRaw {
		res = append(res, data2TopVo(item))
	}
	return res
}

func data2Vo(menu model.Menu) vo.MenuVO {
	var res vo.MenuVO
	res.ParentMenuId = menu.ParentId
	res.Id = menu.ID
	res.ParentName = menu.ParentName
	res.Name = menu.Name
	if menu.Type == 1 {
		res.Type = "category"
	} else {
		res.Type = "article"
	}
	return res
}

func data2TopVo(menu model.Menu) vo.TopMenu {
	var res vo.TopMenu
	res.Id = menu.ID
	res.Name = menu.Name
	res.ParentMenuId = menu.ParentId
	res.ParentName = menu.ParentName
	if menu.Type == 1 {
		res.Type = "category"
	} else {
		res.Type = "article"
	}
	return res
}
