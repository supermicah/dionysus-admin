package dal

import (
	"context"

	"gorm.io/gorm"

	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/errors"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

// GetRoleMenuDB Get role menu storage instance
func GetRoleMenuDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDB(ctx, defDB).Model(new(schema.RoleMenu))
}

// RoleMenu permissions for RBAC
type RoleMenu struct {
	DB *gorm.DB
}

// Query role menus from the database based on the provided parameters and options.
func (a *RoleMenu) Query(ctx context.Context, params schema.RoleMenuQueryParam, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenuQueryResult, error) {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	db := GetRoleMenuDB(ctx, a.DB)
	if v := params.RoleID; v > 0 {
		db = db.Where("role_id = ?", v)
	}

	var list schema.RoleMenus
	pageResult, err := util.WrapPageQuery(ctx, db, params.PaginationParam, opt.QueryOptions, &list)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	queryResult := &schema.RoleMenuQueryResult{
		PageResult: pageResult,
		Data:       list,
	}
	return queryResult, nil
}

// Get the specified role menu from the database.
func (a *RoleMenu) Get(ctx context.Context, id int64, opts ...schema.RoleMenuQueryOptions) (*schema.RoleMenu, error) {
	var opt schema.RoleMenuQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	item := new(schema.RoleMenu)
	ok, err := util.FindOne(ctx, GetRoleMenuDB(ctx, a.DB).Where("id=?", id), opt.QueryOptions, item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}
	return item, nil
}

// Exists checks if the specified role menu exists in the database.
func (a *RoleMenu) Exists(ctx context.Context, id int64) (bool, error) {
	ok, err := util.Exists(ctx, GetRoleMenuDB(ctx, a.DB).Where("id=?", id))
	return ok, errors.WithStack(err)
}

// Create a new role menu.
func (a *RoleMenu) Create(ctx context.Context, item *schema.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Create(item)
	return errors.WithStack(result.Error)
}

// Update the specified role menu in the database.
func (a *RoleMenu) Update(ctx context.Context, item *schema.RoleMenu) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", item.ID).Select("*").Omit("created_at").Updates(item)
	return errors.WithStack(result.Error)
}

// Delete the specified role menu from the database.
func (a *RoleMenu) Delete(ctx context.Context, id int64) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("id=?", id).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}

// DeleteByRoleID Deletes role menus by role id.
func (a *RoleMenu) DeleteByRoleID(ctx context.Context, roleID int64) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("role_id=?", roleID).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}

// DeleteByMenuID Deletes role menus by menu id.
func (a *RoleMenu) DeleteByMenuID(ctx context.Context, menuID int64) error {
	result := GetRoleMenuDB(ctx, a.DB).Where("menu_id=?", menuID).Delete(new(schema.RoleMenu))
	return errors.WithStack(result.Error)
}
