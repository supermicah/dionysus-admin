package sys

import (
	"github.com/google/wire"

	"github.com/supermicah/dionysus-admin/internal/mods/sys/api"
	"github.com/supermicah/dionysus-admin/internal/mods/sys/biz"
	"github.com/supermicah/dionysus-admin/internal/mods/sys/dal"
)

var Set = wire.NewSet(
	wire.Struct(new(SYS), "*"),
	wire.Struct(new(dal.Logger), "*"),
	wire.Struct(new(biz.Logger), "*"),
	wire.Struct(new(api.Logger), "*"),
)
