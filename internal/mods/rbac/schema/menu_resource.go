package schema

import (
	"time"

	"github.com/supermicah/dionysus-admin/internal/config"
	"github.com/supermicah/dionysus-admin/pkg/util"
)

// MenuResource Menu resource management for RBAC
type MenuResource struct {
	ID        int64     `json:"id" gorm:"size:64;primarykey;autoIncrement;"` // Unique ID
	MenuID    int64     `json:"menu_id" gorm:"size:64;index"`                // From Menu.ID
	Method    string    `json:"method" gorm:"size:64;"`                      // HTTP method
	Path      string    `json:"path" gorm:"size:255;"`                       // API request path (e.g. /api/v1/users/:id)
	CreatedAt time.Time `json:"created_at" gorm:"index;"`                    // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`                    // Update time
}

func (a *MenuResource) TableName() string {
	return config.C.FormatTableName("menu_resource")
}

// MenuResourceQueryParam Defining the query parameters for the `MenuResource` struct.
type MenuResourceQueryParam struct {
	util.PaginationParam
	MenuID  int64   `form:"-"` // From Menu.ID
	MenuIDs []int64 `form:"-"` // From Menu.ID
}

// MenuResourceQueryOptions Defining the query options for the `MenuResource` struct.
type MenuResourceQueryOptions struct {
	util.QueryOptions
}

// MenuResourceQueryResult Defining the query result for the `MenuResource` struct.
type MenuResourceQueryResult struct {
	Data       MenuResources
	PageResult *util.PaginationResult
}

// MenuResources Defining the slice of `MenuResource` struct.
type MenuResources []*MenuResource

// MenuResourceForm Defining the data structure for creating a `MenuResource` struct.
type MenuResourceForm struct {
}

// Validate A validation function for the `MenuResourceForm` struct.
func (a *MenuResourceForm) Validate() error {
	return nil
}

func (a *MenuResourceForm) FillTo(menuResource *MenuResource) error {
	return nil
}
