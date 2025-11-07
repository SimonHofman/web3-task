## 1. 数据结构

### [QueryParams](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L5-L8) 结构体
- **功能**：定义分页查询的参数结构
- **组成字段**：
    - `Page int`：当前页码，表单标签为 `page`，最小值验证为1
    - `PageSize int`：每页记录数，表单标签为 `pageSize`，最小值验证为1

### [PagedResult](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L11-L24) 结构体
- **功能**：封装分页查询的结果信息
- **组成字段**：
    - `Total int64`：总记录数
    - `Page int`：当前页码
    - `PageSize int`：每页记录数
    - `TotalPages int`：总页数
    - `HasNextPage bool`：是否有下一页
    - `Data interface{}`：查询结果数据

## 2. 核心函数

### [Paginate](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L27-L53) 函数
- **功能**：通用分页查询函数
- **参数**：
    - `db *gorm.DB`：GORM数据库连接实例
    - `params QueryParams`：分页查询参数
    - `dest interface{}`：查询结果存储目标
- **返回值**：
    - `*PagedResult`：分页结果
    - `error`：可能的错误信息
- **主要流程**：
    1. 查询符合条件的总记录数
    2. 计算总页数和偏移量
    3. 执行分页查询获取当前页数据
    4. 构造并返回分页结果

## 3. 计算逻辑

### 偏移量计算
- 公式：`offset = (params.Page - 1) * params.PageSize`

### 总页数计算
- 公式：`totalPages = (total + int64(params.PageSize) - 1) / int64(params.PageSize)`

### 下一页判断
- 条件：`params.Page < totalPages`

## 4. 依赖组件

### 外部依赖
- `gorm.io/gorm`：GORM数据库ORM库

## 5. 验证规则

### 参数验证
- [Page](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L6-L6) 字段最小值为1
- [PageSize](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L7-L7) 字段最小值为1
- 通过Gin的binding标签实现参数自动验证