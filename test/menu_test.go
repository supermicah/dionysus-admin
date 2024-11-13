package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

func TestMenu(t *testing.T) {
	e := tester(t)

	menuFormItem := schema.MenuForm{
		Code:        "menu",
		Name:        "Menu management",
		Description: "Menu management",
		Sequence:    9,
		Type:        "page",
		Path:        "/system/menu",
		Properties:  `{"icon":"menu"}`,
		Status:      schema.MenuStatusEnabled,
	}

	var menu schema.Menu
	e.POST(baseAPI + "/menus").WithJSON(menuFormItem).
		Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menu})

	as := assert.New(t)
	as.NotEmpty(menu.ID)
	as.Equal(menuFormItem.Code, menu.Code)
	as.Equal(menuFormItem.Name, menu.Name)
	as.Equal(menuFormItem.Description, menu.Description)
	as.Equal(menuFormItem.Sequence, menu.Sequence)
	as.Equal(menuFormItem.Type, menu.Type)
	as.Equal(menuFormItem.Path, menu.Path)
	as.Equal(menuFormItem.Properties, menu.Properties)
	as.Equal(menuFormItem.Status, menu.Status)

	var menus schema.Menus
	e.GET(baseAPI + "/menus").Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &menus})
	as.GreaterOrEqual(len(menus), 1)

	newName := "Menu management 1"
	newStatus := schema.MenuStatusDisabled
	menu.Name = newName
	menu.Status = newStatus
	e.PUT(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).WithJSON(menu).Expect().Status(http.StatusOK)

	var getMenu schema.Menu
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getMenu})
	as.Equal(newName, getMenu.Name)
	as.Equal(newStatus, getMenu.Status)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).Expect().Status(http.StatusNotFound)
}
