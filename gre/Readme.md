* go 使用的是多模块工作区，可以让开发者更容易同时处理多个模块的开发。在 Go 1.17 之前，只能使用 go.mod replace 指令来实现，如果你正巧是同时进行多个模块的开发，使用它可能是很痛苦的。每次当你想要提交代码的时候，都不得不删除掉 go.mod 中的 replace 才能使模块稳定的发布版本。•在使用 go 1.18 多模块工作区功能的时候，就使用这项工作变得简单容易处理<br>
```
go work init 初始化工作区文件，用于生成 go.work 工作区文件

初始化并写入一个新的 go.work 到当前路径下，可以指定需要添加的代码模块
示例: go work init ./hello 将本地仓库 hello 添加到工作区
hello 仓库必须是 go mod 依赖管理的仓库(./hello/go.mod 文件必须存在)


go work use 添加新的模块到工作区

命令示例:
go work use ./example 添加一个模块到工作区
go work use ./example ./example1 添加多个模块到工作区
go work use -r ./example 递归 ./example 目录到当前工作区
删除命令使用 go work edit -dropuse=./example 功能


go work edit 用于编辑 go.work 文件

可以使用 edit 命令编辑和手动编辑 go.work 文件效果是相同的 示例:
go work edit -fmt go.work 重新格式化 go.work 文件
go work edit -replace=github.com/link1st/example=./example go.work 替换代码模块
go work edit -dropreplace=github.com/link1st/example 删除替换代码模块
go work edit -use=./example go.work 添加新的模块到工作区
go work edit -dropuse=./example go.work 从工作区中删除模块

go work sync 将工作区的构建列表同步到工作区的模块

go env GOWORK

查看环境变量，查看当前工作区文件路径 可以排查工作区文件是否设置正确，go.work 路径找不到可以使用 GOWORK 指定

go env GOWORK
$GOPATH/src/link1st/link1st/workspaces/go.work

```