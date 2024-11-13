package biz

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/supermicah/dionysus-admin/internal/config"
	"github.com/supermicah/dionysus-admin/internal/mods/rbac/dal"
	"github.com/supermicah/dionysus-admin/internal/mods/rbac/schema"
	"github.com/supermicah/dionysus-admin/pkg/cachex"
	"github.com/supermicah/dionysus-admin/pkg/encoding/json"
	"github.com/supermicah/dionysus-admin/pkg/encoding/yaml"
	"github.com/supermicah/dionysus-admin/pkg/errors"
	"github.com/supermicah/dionysus-admin/pkg/logging"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

// Menu management for RBAC
type Menu struct {
	Cache           cachex.Cacher
	Trans           *util.Trans
	MenuDAL         *dal.Menu
	MenuResourceDAL *dal.MenuResource
	RoleMenuDAL     *dal.RoleMenu
}

func (a *Menu) InitFromFile(ctx context.Context, menuFile string) error {
	f, err := os.ReadFile(menuFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			logging.Context(ctx).Warn("Menu data file not found, skip init menu data from file", zap.String("file", menuFile))
			return nil
		}
		return err
	}

	var menus schema.Menus
	if ext := filepath.Ext(menuFile); ext == ".json" {
		if err := json.Unmarshal(f, &menus); err != nil {
			return errors.Wrapf(err, "Unmarshal JSON file '%s' failed", menuFile)
		}
	} else if ext == ".yaml" || ext == ".yml" {
		if err := yaml.Unmarshal(f, &menus); err != nil {
			return errors.Wrapf(err, "Unmarshal YAML file '%s' failed", menuFile)
		}
	} else {
		return errors.Errorf("Unsupported file type '%s'", ext)
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		return a.createInBatchByParent(ctx, menus, nil)
	})
}

func (a *Menu) createInBatchByParent(ctx context.Context, items schema.Menus, parent *schema.Menu) error {
	total := len(items)
	for i, item := range items {
		var parentID int64
		if parent != nil {
			parentID = parent.ID
		}

		exist := false

		if item.ID > 0 {
			exists, err := a.MenuDAL.Exists(ctx, item.ID)
			if err != nil {
				return err
			} else if exists {
				exist = true
			}
		} else if item.Code != "" {
			exists, err := a.MenuDAL.ExistsCodeByParentID(ctx, item.Code, parentID)
			if err != nil {
				return err
			} else if exists {
				exist = true
				existItem, err := a.MenuDAL.GetByCodeAndParentID(ctx, item.Code, parentID)
				if err != nil {
					return err
				}
				if existItem != nil {
					item.ID = existItem.ID
				}
			}
		} else if item.Name != "" {
			exists, err := a.MenuDAL.ExistsNameByParentID(ctx, item.Name, parentID)
			if err != nil {
				return err
			} else if exists {
				exist = true
				existItem, err := a.MenuDAL.GetByNameAndParentID(ctx, item.Name, parentID)
				if err != nil {
					return err
				}
				if existItem != nil {
					item.ID = existItem.ID
				}
			}
		}

		if !exist {
			if item.Status == "" {
				item.Status = schema.MenuStatusEnabled
			}
			if item.Sequence == 0 {
				item.Sequence = total - i
			}

			item.ParentID = parentID
			if parent != nil {
				item.ParentPath = fmt.Sprintf("%s%d%s", parent.ParentPath, parentID, util.TreePathDelimiter)
			}
			item.CreatedAt = time.Now()

			if err := a.MenuDAL.Create(ctx, item); err != nil {
				return err
			}
		}

		for _, res := range item.Resources {
			if res.ID > 0 {
				exists, err := a.MenuResourceDAL.Exists(ctx, res.ID)
				if err != nil {
					return err
				} else if exists {
					continue
				}
			}

			if res.Path != "" {
				exists, err := a.MenuResourceDAL.ExistsMethodPathByMenuID(ctx, res.Method, res.Path, item.ID)
				if err != nil {
					return err
				} else if exists {
					continue
				}
			}

			res.MenuID = item.ID
			res.CreatedAt = time.Now()
			if err := a.MenuResourceDAL.Create(ctx, res); err != nil {
				return err
			}
		}

		if item.Children != nil {
			if err := a.createInBatchByParent(ctx, *item.Children, item); err != nil {
				return err
			}
		}
	}
	return nil
}

// Query menus from the data access object based on the provided parameters and options.
func (a *Menu) Query(ctx context.Context, params schema.MenuQueryParam) (*schema.MenuQueryResult, error) {
	params.Pagination = false

	if err := a.fillQueryParam(ctx, &params); err != nil {
		return nil, err
	}

	result, err := a.MenuDAL.Query(ctx, params, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			OrderFields: schema.MenusOrderParams,
		},
	})
	if err != nil {
		return nil, err
	}

	if params.LikeName != "" || params.CodePath != "" {
		result.Data, err = a.appendChildren(ctx, result.Data)
		if err != nil {
			return nil, err
		}
	}

	if params.IncludeResources {
		for i, item := range result.Data {
			resResult, err := a.MenuResourceDAL.Query(ctx, schema.MenuResourceQueryParam{
				MenuID: item.ID,
			})
			if err != nil {
				return nil, err
			}
			result.Data[i].Resources = resResult.Data
		}
	}

	result.Data = result.Data.ToTree()
	return result, nil
}

func (a *Menu) fillQueryParam(ctx context.Context, params *schema.MenuQueryParam) error {
	if params.CodePath != "" {
		var (
			lastMenu schema.Menu
		)
		codes := strings.Split(params.CodePath, util.TreePathDelimiter)
		length := len(codes)
		for i := 0; i < length; i++ {
			code := codes[i]
			if code == "" {
				continue
			}
			if i == length-1 {
				params.Code = code
				continue
			}
			menu, err := a.MenuDAL.GetByCodeAndParentID(ctx, code, lastMenu.ID, schema.MenuQueryOptions{
				QueryOptions: util.QueryOptions{
					SelectFields: []string{"id", "parent_id", "parent_path"},
				},
			})
			if err != nil {
				return err
			} else if menu == nil {
				return errors.NotFound("", "Menu not found by code '%s'", strings.Join(codes, util.TreePathDelimiter))
			}
			lastMenu = *menu
		}
		params.ParentPathPrefix = fmt.Sprintf("%s%d%s", lastMenu.ParentPath, lastMenu.ID, util.TreePathDelimiter)
	}
	return nil
}

func (a *Menu) appendChildren(ctx context.Context, data schema.Menus) (schema.Menus, error) {
	if len(data) == 0 {
		return data, nil
	}

	existsInData := func(id int64) bool {
		for _, item := range data {
			if item.ID == id {
				return true
			}
		}
		return false
	}

	for _, item := range data {
		childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			ParentPathPrefix: fmt.Sprintf("%s%d%s", item.ParentPath, item.ID, util.TreePathDelimiter),
		})
		if err != nil {
			return nil, err
		}
		for _, child := range childResult.Data {
			if existsInData(child.ID) {
				continue
			}
			data = append(data, child)
		}
	}

	if parentIDs := data.SplitParentIDs(); len(parentIDs) > 0 {
		parentResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			InIDs: parentIDs,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range parentResult.Data {
			if existsInData(p.ID) {
				continue
			}
			data = append(data, p)
		}
	}
	sort.Sort(data)

	return data, nil
}

// Get the specified menu from the data access object.
func (a *Menu) Get(ctx context.Context, id int64) (*schema.Menu, error) {
	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return nil, err
	} else if menu == nil {
		return nil, errors.NotFound("", "Menu not found")
	}

	menuResResult, err := a.MenuResourceDAL.Query(ctx, schema.MenuResourceQueryParam{
		MenuID: menu.ID,
	})
	if err != nil {
		return nil, err
	}
	menu.Resources = menuResResult.Data

	return menu, nil
}

// Create a new menu in the data access object.
func (a *Menu) Create(ctx context.Context, formItem *schema.MenuForm) (*schema.Menu, error) {
	menu := &schema.Menu{
		CreatedAt: time.Now(),
	}

	if parentID := formItem.ParentID; parentID > 0 {
		parent, err := a.MenuDAL.Get(ctx, parentID)
		if err != nil {
			return nil, err
		} else if parent == nil {
			return nil, errors.NotFound("", "Parent not found")
		}
		menu.ParentPath = fmt.Sprintf("%s%d%s", parent.ParentPath, parent.ID, util.TreePathDelimiter)
	}

	if exists, err := a.MenuDAL.ExistsCodeByParentID(ctx, formItem.Code, formItem.ParentID); err != nil {
		return nil, err
	} else if exists {
		return nil, errors.BadRequest("", "Menu code already exists at the same level")
	}

	if err := formItem.FillTo(menu); err != nil {
		return nil, err
	}

	err := a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.MenuDAL.Create(ctx, menu); err != nil {
			return err
		}

		for _, res := range formItem.Resources {
			res.MenuID = menu.ID
			res.CreatedAt = time.Now()
			if err := a.MenuResourceDAL.Create(ctx, res); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return menu, nil
}

// Update the specified menu in the data access object.
func (a *Menu) Update(ctx context.Context, id int64, formItem *schema.MenuForm) error {
	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	oldParentPath := menu.ParentPath
	oldStatus := menu.Status
	var childData schema.Menus
	if menu.ParentID != formItem.ParentID {
		if parentID := formItem.ParentID; parentID > 0 {
			parent, err := a.MenuDAL.Get(ctx, parentID)
			if err != nil {
				return err
			} else if parent == nil {
				return errors.NotFound("", "Parent not found")
			}
			menu.ParentPath = fmt.Sprintf("%s%d%s", parent.ParentPath, parent.ID, util.TreePathDelimiter)
		} else {
			menu.ParentPath = ""
		}

		childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
			ParentPathPrefix: fmt.Sprintf("%s%d%s", oldParentPath, menu.ID, util.TreePathDelimiter),
		}, schema.MenuQueryOptions{
			QueryOptions: util.QueryOptions{
				SelectFields: []string{"id", "parent_path"},
			},
		})
		if err != nil {
			return err
		}
		childData = childResult.Data
	}

	if menu.Code != formItem.Code {
		if exists, err := a.MenuDAL.ExistsCodeByParentID(ctx, formItem.Code, formItem.ParentID); err != nil {
			return err
		} else if exists {
			return errors.BadRequest("", "Menu code already exists at the same level")
		}
	}

	if err := formItem.FillTo(menu); err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if oldStatus != formItem.Status {
			oldPath := fmt.Sprintf("%s%d%s", oldParentPath, menu.ID, util.TreePathDelimiter)
			if err := a.MenuDAL.UpdateStatusByParentPath(ctx, oldPath, formItem.Status); err != nil {
				return err
			}
		}

		for _, child := range childData {
			oldPath := fmt.Sprintf("%s%d%s", oldParentPath, menu.ID, util.TreePathDelimiter)
			newPath := fmt.Sprintf("%s%d%s", menu.ParentPath, menu.ID, util.TreePathDelimiter)
			err := a.MenuDAL.UpdateParentPath(ctx, child.ID, strings.Replace(child.ParentPath, oldPath, newPath, 1))
			if err != nil {
				return err
			}
		}

		if err := a.MenuDAL.Update(ctx, menu); err != nil {
			return err
		}

		if err := a.MenuResourceDAL.DeleteByMenuID(ctx, id); err != nil {
			return err
		}
		for _, res := range formItem.Resources {
			res.MenuID = id
			if res.CreatedAt.IsZero() {
				res.CreatedAt = time.Now()
			}
			res.UpdatedAt = time.Now()
			if err := a.MenuResourceDAL.Create(ctx, res); err != nil {
				return err
			}
		}

		return a.syncToCasbin(ctx)
	})
}

// Delete the specified menu from the data access object.
func (a *Menu) Delete(ctx context.Context, id int64) error {
	if config.C.General.DenyDeleteMenu {
		return errors.BadRequest("", "Menu deletion is not allowed")
	}

	menu, err := a.MenuDAL.Get(ctx, id)
	if err != nil {
		return err
	} else if menu == nil {
		return errors.NotFound("", "Menu not found")
	}

	childResult, err := a.MenuDAL.Query(ctx, schema.MenuQueryParam{
		ParentPathPrefix: fmt.Sprintf("%s%d%s", menu.ParentPath, menu.ID, util.TreePathDelimiter),
	}, schema.MenuQueryOptions{
		QueryOptions: util.QueryOptions{
			SelectFields: []string{"id"},
		},
	})
	if err != nil {
		return err
	}

	return a.Trans.Exec(ctx, func(ctx context.Context) error {
		if err := a.delete(ctx, id); err != nil {
			return err
		}

		for _, child := range childResult.Data {
			if err := a.delete(ctx, child.ID); err != nil {
				return err
			}
		}

		return a.syncToCasbin(ctx)
	})
}

func (a *Menu) delete(ctx context.Context, id int64) error {
	if err := a.MenuDAL.Delete(ctx, id); err != nil {
		return err
	}
	if err := a.MenuResourceDAL.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	if err := a.RoleMenuDAL.DeleteByMenuID(ctx, id); err != nil {
		return err
	}
	return nil
}

func (a *Menu) syncToCasbin(ctx context.Context) error {
	return a.Cache.Set(ctx, config.CacheNSForRole, config.CacheKeyForSyncToCasbin, fmt.Sprintf("%d", time.Now().Unix()))
}
