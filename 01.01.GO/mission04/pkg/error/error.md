## 1. 数据结构

### [AppError](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L8-L12) 结构体
- **功能**：定义应用程序错误的统一格式
- **组成字段**：
    - `Code int`：HTTP状态码，JSON标签为 `-`（不序列化到响应中）
    - `ErrCode string`：错误代码，JSON标签为 `error_code`
    - `Message string`：错误消息，JSON标签为 `message`

### [Error()](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L14-L16) 方法
- **功能**：实现 `error` 接口，返回错误消息
- **返回值**：`string` - 错误描述信息

## 2. 错误处理函数

### [ThrowErr](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L18-L24) 函数
- **功能**：在Gin上下文中抛出应用错误并终止请求处理
- **参数**：
    - `c *gin.Context`：Gin上下文
    - `appErr *AppError`：应用错误对象
    - `message string`：可选的自定义错误消息
- **行为**：
    1. 若提供自定义消息则更新错误对象的 [Message](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L11-L11) 字段
    2. 将错误添加到Gin上下文的错误列表中
    3. 使用 `c.Abort()` 终止请求处理链

## 3. 预定义错误变量

### 系统级错误
- [ErrSystem](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L27-L27)：系统内部错误，HTTP 500状态码
- [ErrUserNotFound](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L28-L28)：用户不存在错误，HTTP 404状态码

### 认证授权错误
- [ErrInvalidCredentials](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L29-L29)：认证失败错误，HTTP 401状态码
- [ErrUnauthorized](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L30-L30)：权限不足错误，HTTP 403状态码

### 参数验证错误
- [ErrInvalidParams](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L31-L31)：请求参数错误，HTTP 400状态码

## 4. 依赖组件

### 外部依赖
- `net/http`：HTTP协议支持和状态码定义
- `github.com/gin-gonic/gin`：Gin Web框架

## 5. 设计特点

### 统一错误格式
- 所有应用错误都遵循 [AppError](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L8-L12) 结构
- 包含错误代码和错误消息，便于前端处理

### 上下文集成
- 通过 [ThrowErr](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L18-L24) 函数与Gin上下文无缝集成
- 支持错误链和请求终止机制

### 预定义错误集合
- 提供常用的预定义错误变量，便于统一管理和使用