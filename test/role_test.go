package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

func TestRole(t *testing.T) {
	e := tester(t)

	menuFormItem := schema.MenuForm{
		Code:        "role",
		Name:        "Role management",
		Description: "Role management",
		Sequence:    8,
		Type:        "page",
		Path:        "/system/role",
		Properties:  `{"icon":"role"}`,
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

	roleFormItem := schema.RoleForm{
		Code: "admin",
		Name: "Administrator",
		Menus: schema.RoleMenus{
			{MenuID: menu.ID},
		},
		Description: "Administrator",
		Sequence:    9,
		Status:      schema.RoleStatusEnabled,
	}

	var role schema.Role
	e.POST(baseAPI + "/roles").WithJSON(roleFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &role})
	as.NotEmpty(role.ID)
	as.Equal(roleFormItem.Code, role.Code)
	as.Equal(roleFormItem.Name, role.Name)
	as.Equal(roleFormItem.Description, role.Description)
	as.Equal(roleFormItem.Sequence, role.Sequence)
	as.Equal(roleFormItem.Status, role.Status)
	as.Equal(len(roleFormItem.Menus), len(role.Menus))

	var roles schema.Roles
	e.GET(baseAPI + "/roles").Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &roles})
	as.GreaterOrEqual(len(roles), 1)

	newName := "Administrator 1"
	newStatus := schema.RoleStatusDisabled
	role.Name = newName
	role.Status = newStatus
	e.PUT(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).WithJSON(role).Expect().Status(http.StatusOK)

	var getRole schema.Role
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getRole})
	as.Equal(newName, getRole.Name)
	as.Equal(newStatus, getRole.Status)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusNotFound)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/menus/", menu.ID)).Expect().Status(http.StatusNotFound)
}
