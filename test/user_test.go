package test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/crypto/hash"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

func TestUser(t *testing.T) {
	e := tester(t)

	menuFormItem := schema.MenuForm{
		Code:        "user",
		Name:        "User management",
		Description: "User management",
		Sequence:    7,
		Type:        "page",
		Path:        "/system/user",
		Properties:  `{"icon":"user"}`,
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
		Code: "user",
		Name: "Normal",
		Menus: schema.RoleMenus{
			{MenuID: menu.ID},
		},
		Description: "Normal",
		Sequence:    8,
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

	userFormItem := schema.UserForm{
		Username: "test",
		Name:     "Test",
		Password: hash.MD5String("test"),
		Phone:    "0720",
		Email:    "test@gmail.com",
		Remark:   "test user",
		Status:   schema.UserStatusActivated,
		Roles:    schema.UserRoles{{RoleID: role.ID}},
	}

	var user schema.User
	e.POST(baseAPI + "/users").WithJSON(userFormItem).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &user})
	as.NotEmpty(user.ID)
	as.Equal(userFormItem.Username, user.Username)
	as.Equal(userFormItem.Name, user.Name)
	as.Equal(userFormItem.Phone, user.Phone)
	as.Equal(userFormItem.Email, user.Email)
	as.Equal(userFormItem.Remark, user.Remark)
	as.Equal(userFormItem.Status, user.Status)
	as.Equal(len(userFormItem.Roles), len(user.Roles))

	var users schema.Users
	e.GET(baseAPI+"/users").WithQuery("username", userFormItem.Username).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &users})
	as.GreaterOrEqual(len(users), 1)

	newName := "Test 1"
	newStatus := schema.UserStatusFreezed
	user.Name = newName
	user.Status = newStatus
	e.PUT(fmt.Sprintf("%s%s%d", baseAPI, "/users/", user.ID)).WithJSON(user).Expect().Status(http.StatusOK)

	var getUser schema.User
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/users/", user.ID)).Expect().Status(http.StatusOK).JSON().Decode(&util.ResponseResult{Data: &getUser})
	as.Equal(newName, getUser.Name)
	as.Equal(newStatus, getUser.Status)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/users/", user.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/users/", user.ID)).Expect().Status(http.StatusNotFound)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusNotFound)

	e.DELETE(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusOK)
	e.GET(fmt.Sprintf("%s%s%d", baseAPI, "/roles/", role.ID)).Expect().Status(http.StatusNotFound)
}
