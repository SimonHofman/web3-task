## 1. 路由初始化函数

### [InitRouter](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\router\router.go#L9-L45) 函数
- **功能**：初始化并配置 Gin 路由引擎，定义所有 API 端点
- **返回值**：`*gin.Engine` - 配置好的 Gin 路由引擎实例
- **主要流程**：
    1. 设置 Gin 运行模式为调试模式
    2. 创建新的 Gin 引擎实例
    3. 注册全局中间件
    4. 配置 API v1 版本的路由组及各业务模块路由

## 2. 中间件配置

### 全局中间件
- `gin.Recovery()`：Gin 内置的崩溃恢复中间件
- [middleware.GlobalErrorHandlerMiddleware()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\globalErrorHandlerMiddleware.go#L13-L42)：全局错误处理中间件
- [middleware.LoggerMiddleware()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\loggerMiddleware.go#L11-L28)：请求日志记录中间件

### 路由组中间件
- [middleware.AuthMiddleware()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\authMiddleware.go#L11-L38)：认证鉴权中间件，应用于需要登录验证的路由组

## 3. 路由组结构

### API 版本组
- **路径前缀**：`/api/v1`
- **子组**：
    - `authGroup`：认证相关接口组（无需认证）
    - `userGroup`：用户相关接口组（需认证）
    - `postGroup`：文章相关接口组（需认证）
    - `commentGroup`：评论相关接口组（需认证）

## 4. 接口端点定义

### 认证接口组 (`/api/v1/auth`)
- `POST /register`：用户注册接口，绑定 [handler.Register](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\userHandler.go#L26-L37)
- `POST /login`：用户登录接口，绑定 [handler.Login](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\userHandler.go#L40-L53)

### 用户接口组 (`/api/v1/user`)
- `GET /page`：获取用户分页列表，绑定 [handler.UserPage](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\userHandler.go#L11-L23)

### 文章接口组 (`/api/v1/post`)
- `POST /create`：创建文章，绑定 [handler.CreatePost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L14-L28)
- `GET /page`：获取文章分页列表，绑定 [handler.PostPage](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L31-L44)
- `GET /byId`：根据ID获取文章详情，绑定 [handler.PostById](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L47-L60)
- `POST /edit`：编辑文章，绑定 [handler.EditPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L63-L80)
- `GET /delete`：删除文章，绑定 [handler.DelPost](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\postHandler.go#L83-L99)

### 评论接口组 (`/api/v1/comment`)
- `POST /create`：创建评论，绑定 [handler.CreateComment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L12-L25)
- `GET /byPostId`：根据文章ID获取评论列表，绑定 [handler.CommentByPostId](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L28-L39)

## 5. 依赖组件

### 外部框架
- `github.com/gin-gonic/gin`：Gin Web 框架

### 内部模块
- `mission04/internal/handler`：HTTP 请求处理器
- `mission04/internal/middleware`：HTTP 中间件组件

## 6. 安全设计

### 认证保护
- 认证、用户、文章、评论接口组中，除认证接口外均需通过 [middleware.AuthMiddleware](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\authMiddleware.go#L11-L38) 验证
- 确保只有已登录用户才能访问受保护的资源