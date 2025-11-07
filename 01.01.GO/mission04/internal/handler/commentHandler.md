## 1. 处理函数

### [CreateComment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L12-L25) 函数
- **功能**：处理创建评论的HTTP请求
- **参数**：`c *gin.Context` - Gin上下文
- **主要流程**：
    1. 解析请求体中的JSON数据到 [model.Comment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\model\comment.go#L4-L11) 对象
    2. 从上下文中获取用户ID（通过认证中间件设置）
    3. 调用 `logic.CommentLogic.CreateComment` 执行业务逻辑
    4. 根据结果返回相应响应

### [CommentByPostId](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L28-L41) 函数
- **功能**：处理查询某篇文章所有评论的HTTP请求
- **参数**：`c *gin.Context` - Gin上下文
- **主要流程**：
    1. 从URL查询参数中获取 `postId`
    2. 调用 `logic.CommentLogic.CommentByPostId` 查询评论列表
    3. 根据结果返回相应响应

## 2. 依赖组件

### 内部模块
- `mission04/internal/logic`：业务逻辑层
- `mission04/internal/model`：数据模型定义
- `mission04/pkg/error`：错误处理模块
- `mission04/pkg/response`：HTTP响应处理模块

### 外部框架
- `github.com/gin-gonic/gin`：Gin Web框架

## 3. 数据处理流程

### 评论创建流程
1. **数据绑定**：使用 `c.ShouldBindJSON` 将请求体绑定到 [model.Comment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\model\comment.go#L4-L11)
2. **参数验证**：若绑定失败则返回 `error2.ErrInvalidParams` 错误
3. **身份获取**：从上下文获取用户ID `c.MustGet("userID")`
4. **业务调用**：调用 `logic.CommentLogic.CreateComment` 执行创建操作
5. **结果响应**：根据执行结果返回成功或失败响应

### 评论查询流程
1. **参数获取**：使用 `c.GetQuery` 获取URL查询参数 `postId`
2. **参数验证**：若参数不存在则返回 `error2.ErrInvalidParams` 错误
3. **业务调用**：调用 `logic.CommentLogic.CommentByPostId` 查询评论列表
4. **结果响应**：将查询结果包装在响应中返回

## 4. 错误处理机制

### 参数错误
- 使用 [response.Error(c, error2.ErrInvalidParams)](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L31-L37) 处理请求参数绑定或缺失错误

### 系统错误
- 使用 `response.Fail(c, error2.ErrSystem.Code, "错误信息")` 处理业务逻辑执行失败

### 成功响应
- 创建成功：`response.Success(c, nil, "创建评论成功")`
- 查询成功：`response.Success(c, list, "查询文章评论成功")`