## 1. 数据结构

### [Response](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L9-L13) 结构体
- **功能**：定义统一的HTTP响应格式
- **组成字段**：
    - `Code int`：响应状态码，JSON标签为 `code`
    - `Message string`：响应消息，JSON标签为 `message`
    - `Data interface{}`：响应数据，JSON标签为 `data`

## 2. 响应处理函数

### [Success](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L15-L21) 函数
- **功能**：返回成功响应
- **参数**：
    - `c *gin.Context`：Gin上下文
    - `data interface{}`：响应数据
    - `msg string`：成功消息
- **行为**：返回HTTP 200状态码，[Code](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L9-L9)字段设置为200

### [Fail](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L23-L29) 函数
- **功能**：返回失败响应
- **参数**：
    - `c *gin.Context`：Gin上下文
    - `code int`：错误状态码
    - `msg string`：错误消息
- **行为**：返回HTTP 200状态码，[Data](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\db\pagination.go#L23-L23)字段设置为nil

### [Error](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L31-L37) 函数
- **功能**：返回应用错误响应
- **参数**：
    - `c *gin.Context`：Gin上下文
    - `appErr *error2.AppError`：应用错误对象
- **行为**：使用 [AppError](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L8-L12) 中的 [Code](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L9-L9) 和 [Message](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\error\error.go#L11-L11) 字段构建响应

### [FailStop](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L39-L45) 函数
- **功能**：返回失败响应并终止请求处理
- **参数**：
    - `c *gin.Context`：Gin上下文
    - `code int`：错误状态码
    - `msg string`：错误消息
- **行为**：使用 `AbortWithStatusJSON` 终止请求链并返回响应

## 3. 依赖组件

### 外部依赖
- `github.com/gin-gonic/gin`：Gin Web框架
- `net/http`：HTTP协议支持

### 内部依赖
- `mission04/pkg/error`：应用错误处理模块

## 4. 设计特点

### 统一响应格式
- 所有HTTP响应都遵循相同的 [Response](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L9-L13) 结构
- 包含状态码、消息和数据三个核心字段

### 请求处理控制
- [Fail](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L23-L29) 与 [FailStop](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L39-L45) 的区别在于是否终止请求处理链
- [FailStop](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L39-L45) 使用 `AbortWithStatusJSON` 确保请求不会继续处理