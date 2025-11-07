## 1. 数据结构

### [postLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\postLogic.go#L8-L8) 结构体
- **功能**：实现文章相关业务逻辑的结构体
- **特点**：空结构体，仅用于方法绑定

## 2. 全局变量

### [PostLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\postLogic.go#L10-L10) 变量
- **类型**：`*postLogic`
- **功能**：文章业务逻辑的单例实例
- **初始化**：使用 `new(postLogic)` 创建

## 3. 业务逻辑方法

### [CreatePost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L14-L28) 方法
- **功能**：创建新文章的业务逻辑
- **参数**：`post *model.Post` - 文章信息
- **返回值**：`error` - 数据库操作错误
- **流程**：直接将文章信息保存到数据库

### [PostPage](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L31-L44) 方法
- **功能**：分页查询用户的文章列表
- **参数**：
    - `c *db.QueryParams` - 分页查询参数
    - `userId uint` - 用户ID
- **返回值**：
    - `*db.PagedResult` - 分页结果
    - `error` - 查询错误
- **流程**：根据用户ID分页查询文章

### [PostById](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L47-L60) 方法
- **功能**：根据ID查询文章详情
- **参数**：
    - `postId string` - 文章ID
    - `userId uint` - 用户ID
- **返回值**：
    - `*model.Post` - 文章详情
    - `error` - 查询错误
- **流程**：查询指定用户的文章详情

### [EditPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L63-L80) 方法
- **功能**：修改文章的业务逻辑
- **参数**：
    - `post *model.Post` - 更新的文章信息
    - `userId *uint` - 当前用户ID
- **返回值**：`error` - 操作错误或权限错误
- **流程**：
    1. 验证文章是否属于当前用户
    2. 若无权限返回 `error2.ErrUnauthorized`
    3. 执行文章更新操作

### [DelPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L83-L99) 方法
- **功能**：删除文章的业务逻辑
- **参数**：
    - `postId *string` - 要删除的文章ID
    - `userId *uint` - 当前用户ID
- **返回值**：`error` - 操作错误或权限错误
- **流程**：
    1. 验证文章是否属于当前用户
    2. 若无权限返回 `error2.ErrUnauthorized`
    3. 执行文章删除操作

## 4. 依赖组件

### 内部模块
- `mission04/internal/model`：数据模型定义
- `mission04/pkg/db`：数据库访问层和分页工具
- `mission04/pkg/error`：错误定义模块

## 5. 权限控制

### 用户文章权限验证
- 在 [EditPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L63-L80) 和 [DelPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L83-L99) 方法中验证操作权限
- 查询条件：`WHERE id = ? and user_id = ?`
- 无权限时返回 `error2.ErrUnauthorized` 错误