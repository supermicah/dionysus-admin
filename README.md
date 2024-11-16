# dionysus-admin

> A test API service based on golang.

## Project Create
```shell
go-framework-cli new -d ./Golang/ --name dionysus-admin --app-name dionysus-admin --desc 'A test API service based on golang.' --pkg 'github.com/supermicah/dionysus-admin' --fe-dir ./React/ --fe-name dionysus-web


#gen // 主命令
#    --dir(-d) // 文件地址
#    --module(-m) // 模块名称，需要大写
#    --module-path // 默认在internal/mods下
#    --wire-path // 默认在internal/wirex下
#    --internal/swagger //默认在internal/swagger下
#    --config(-c) // 配置文件地址
#    --structs(-s) // 结构体名称
#    --structs-comment // 结构体的备注
#    --structs-router-prefix // 结构体的路由前缀 默认是/api/v1
#    --structs-output // 默认列表是：schema,dal,biz,api
#    --tpl-path // 模版地址，默认：tpls
#    --tpl-type // 模版类型（default：go、react：fe、ant-design-pro-v5：fe）
#    --fe-dir // 前端项目地址
```

## Quick Start

```bash
make start
```

## Build

```bash
make build
```

## Generate wire inject files

```bash
make wire
```

## Generate swagger documents

```bash
make swagger
```

