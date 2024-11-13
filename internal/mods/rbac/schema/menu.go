package schema

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/supermicah/dionysus-admin/internal/config"
	"github.com/supermicah/dionysus-admin/pkg/errors"
	"github.com/supermicah/dionysus-admin/pkg/logging"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

const (
	MenuStatusDisabled = "disabled"
	MenuStatusEnabled  = "enabled"
)

var (
	MenusOrderParams = []util.OrderByParam{
		{Field: "sequence", Direction: util.DESC},
		{Field: "created_at", Direction: util.DESC},
	}
)

// Menu management for RBAC
type Menu struct {
	ID          int64         `json:"id" gorm:"size:64;primarykey;autoIncrement;"` // Unique ID
	Code        string        `json:"code" gorm:"size:32;index;"`                  // Code of menu (unique for each level)
	Name        string        `json:"name" gorm:"size:128;index"`                  // Display name of menu
	Description string        `json:"description" gorm:"size:1024"`                // Details about menu
	Sequence    int           `json:"sequence" gorm:"index;"`                      // Sequence for sorting (Order by desc)
	Type        string        `json:"type" gorm:"size:20;index"`                   // Type of menu (page, button)
	Path        string        `json:"path" gorm:"size:255;"`                       // Access path of menu
	Properties  string        `json:"properties" gorm:"type:text;"`                // Properties of menu (JSON)
	Status      string        `json:"status" gorm:"size:20;index"`                 // Status of menu (enabled, disabled)
	ParentID    int64         `json:"parent_id" gorm:"size:64;index;"`             // Parent ID (From Menu.ID)
	ParentPath  string        `json:"parent_path" gorm:"size:255;index;"`          // Parent path (split by .)
	Children    *Menus        `json:"children" gorm:"-"`                           // Child menus
	CreatedAt   time.Time     `json:"created_at" gorm:"index;"`                    // Create time
	UpdatedAt   time.Time     `json:"updated_at" gorm:"index;"`                    // Update time
	Resources   MenuResources `json:"resources" gorm:"-"`                          // Resources of menu
}

func (a *Menu) TableName() string {
	return config.C.FormatTableName("menu")
}

// MenuQueryParam Defining the query parameters for the `Menu` struct.
type MenuQueryParam struct {
	util.PaginationParam
	CodePath         string  `form:"code"`             // Code path (like xxx.xxx.xxx)
	LikeName         string  `form:"name"`             // Display name of menu
	IncludeResources bool    `form:"includeResources"` // Include resources
	InIDs            []int64 `form:"-"`                // Include menu IDs
	Status           string  `form:"-"`                // Status of menu (disabled, enabled)
	ParentID         int64   `form:"-"`                // Parent ID (From Menu.ID)
	ParentPathPrefix string  `form:"-"`                // Parent path (split by .)
	Code             string  `form:"-"`                // Code (like xxx)
	UserID           int64   `form:"-"`                // User ID
	RoleID           int64   `form:"-"`                // Role ID
}

// MenuQueryOptions Defining the query options for the `Menu` struct.
type MenuQueryOptions struct {
	util.QueryOptions
}

// MenuQueryResult Defining the query result for the `Menu` struct.
type MenuQueryResult struct {
	Data       Menus
	PageResult *util.PaginationResult
}

// Menus Defining the slice of `Menu` struct.
type Menus []*Menu

func (a Menus) Len() int {
	return len(a)
}

func (a Menus) Less(i, j int) bool {
	if a[i].Sequence == a[j].Sequence {
		return a[i].CreatedAt.Unix() > a[j].CreatedAt.Unix()
	}
	return a[i].Sequence > a[j].Sequence
}

func (a Menus) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Menus) ToMap() map[int64]*Menu {
	m := make(map[int64]*Menu)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

func (a Menus) SplitParentIDs() []int64 {
	parentIDs := make([]int64, 0, len(a))
	idMapper := make(map[int64]struct{})
	for _, item := range a {
		if _, ok := idMapper[item.ID]; ok {
			continue
		}
		idMapper[item.ID] = struct{}{}
		if pp := item.ParentPath; pp != "" {
			for _, pid := range strings.Split(pp, util.TreePathDelimiter) {
				if pid == "" {
					continue
				}
				parentID, err := strconv.ParseInt(pid, 10, 64)
				if err != nil {
					logging.Context(context.Background()).Error("Failed to parse pid value", zap.Error(err), zap.String("pid", pid))
					continue
				}
				if _, ok := idMapper[parentID]; ok {
					continue
				}
				parentIDs = append(parentIDs, parentID)
				idMapper[parentID] = struct{}{}
			}
		}
	}
	return parentIDs
}

func (a Menus) ToTree() Menus {
	var list Menus
	m := a.ToMap()
	for _, item := range a {
		if item.ParentID <= 0 {
			list = append(list, item)
			continue
		}
		if parent, ok := m[item.ParentID]; ok {
			if parent.Children == nil {
				children := Menus{item}
				parent.Children = &children
				continue
			}
			*parent.Children = append(*parent.Children, item)
		}
	}
	return list
}

// MenuForm Defining the data structure for creating a `Menu` struct.
type MenuForm struct {
	Code        string        `json:"code" binding:"required,max=32"`                   // Code of menu (unique for each level)
	Name        string        `json:"name" binding:"required,max=128"`                  // Display name of menu
	Description string        `json:"description"`                                      // Details about menu
	Sequence    int           `json:"sequence"`                                         // Sequence for sorting (Order by desc)
	Type        string        `json:"type" binding:"required,oneof=page button"`        // Type of menu (page, button)
	Path        string        `json:"path"`                                             // Access path of menu
	Properties  string        `json:"properties"`                                       // Properties of menu (JSON)
	Status      string        `json:"status" binding:"required,oneof=disabled enabled"` // Status of menu (enabled, disabled)
	ParentID    int64         `json:"parent_id"`                                        // Parent ID (From Menu.ID)
	Resources   MenuResources `json:"resources"`                                        // Resources of menu
}

// Validate A validation function for the `MenuForm` struct.
func (a *MenuForm) Validate() error {
	if v := a.Properties; v != "" {
		if !json.Valid([]byte(v)) {
			return errors.BadRequest("", "invalid properties")
		}
	}
	return nil
}

func (a *MenuForm) FillTo(menu *Menu) error {
	menu.Code = a.Code
	menu.Name = a.Name
	menu.Description = a.Description
	menu.Sequence = a.Sequence
	menu.Type = a.Type
	menu.Path = a.Path
	menu.Properties = a.Properties
	menu.Status = a.Status
	menu.ParentID = a.ParentID
	return nil
}
