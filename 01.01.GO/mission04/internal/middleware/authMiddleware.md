## 1. 中间件函数

### [AuthMiddleware](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\authMiddleware.go#L11-L38) 函数
- **功能**：认证中间件，验证请求中的JWT token
- **返回值**：`gin.HandlerFunc` - Gin处理器函数
- **主要流程**：
    1. 从请求头获取 `Authorization` 字段
    2. 验证token格式是否为Bearer类型
    3. 调用 [auth.ParseToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L69-L83) 验证token有效性
    4. 将用户信息存入上下文供后续处理使用

### [RoleMiddleware](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\authMiddleware.go#L41-L52) 函数
- **功能**：角色权限中间件，验证用户是否具有指定角色
- **参数**：`requiredRole string` - 所需的角色权限
- **返回值**：`gin.HandlerFunc` - Gin处理器函数
- **主要流程**：
    1. 从上下文获取用户角色信息
    2. 验证用户角色是否匹配所需角色
    3. 不匹配时中断请求处理并返回未授权错误

## 2. 认证流程

### Token验证流程
1. **获取token**：从 `Authorization` 请求头获取token
2. **格式检查**：验证是否为 `Bearer <token>` 格式
3. **解析验证**：调用 [auth.ParseToken](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L69-L83) 解析并验证token
4. **信息存储**：将解析出的 [UserID](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L11-L11) 和 [Role](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L12-L12) 存储到上下文中

### 错误处理
- 缺少token：返回 `error2.ErrInvalidCredentials` 错误
- 格式错误：返回 `error2.ErrInvalidCredentials` 错误
- token无效：返回 `error2.ErrInvalidCredentials` 错误
- 权限不足：返回 `error2.ErrUnauthorized` 错误

## 3. 上下文操作

### 用户信息存储
- 使用 `c.Set("userID", claims.UserID)` 存储用户ID
- 使用 `c.Set("userRole", claims.Role)` 存储用户角色

### 用户信息获取
- 使用 `c.Get("userRole")` 获取用户角色进行权限验证

## 4. 依赖组件

### 内部模块
- `mission04/pkg/auth`：JWT token解析和验证
- `mission04/pkg/error`：错误处理和响应

### 外部库
- `github.com/gin-gonic/gin`：Gin Web框架
- `strings`：字符串处理
- `fmt`：格式化处理

## 5. 安全机制

### Bearer Token验证
- 严格按照HTTP Bearer Token规范解析token
- 验证token格式的完整性

### 用户信息保护
- 仅将必要用户信息（ID和角色）存入上下文
- 避免敏感信息泄露

### 请求中断机制
- 验证失败时使用 `c.Abort()` 中断请求处理链
- 防止未认证请求继续处理