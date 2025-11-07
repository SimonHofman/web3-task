## 1. 中间件函数

### [GlobalErrorHandlerMiddleware](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\middleware\globalErrorHandlerMiddleware.go#L13-L42) 函数
- **功能**：全局错误处理中间件，捕获并统一处理请求处理过程中发生的各种错误
- **返回值**：`gin.HandlerFunc` - Gin处理器函数
- **处理机制**：
    1. 先执行后续处理器链（`c.Next()`）
    2. 检查上下文中是否包含错误信息
    3. 根据错误类型进行分类处理并返回相应响应

## 2. 错误处理流程

### 错误检测
- 通过 `len(c.Errors) > 0` 检测是否有错误发生
- 遍历 `c.Errors` 处理每个错误

### 错误分类处理
1. **认证错误**：`error2.ErrInvalidCredentials` - 返回认证失败响应
2. **参数错误**：`error2.ErrInvalidParams` - 返回参数无效响应
3. **权限错误**：`error2.ErrUnauthorized` - 返回未授权响应
4. **数据不存在**：`gorm.ErrRecordNotFound` - 返回数据不存在响应
5. **默认错误**：其他未分类错误 - 记录日志并返回系统错误响应

## 3. 响应处理

### 错误响应
- 使用 [response.FailStop](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L39-L45) 终止请求处理并返回错误响应
- 不同类型错误返回对应的错误码和错误信息
- 系统错误会记录详细日志 `log.Logger.Error`

## 4. 依赖组件

### 内部模块
- `mission04/pkg/error`：自定义错误定义
- `mission04/pkg/log`：日志记录
- `mission04/pkg/response`：HTTP响应处理

### 外部库
- `github.com/gin-gonic/gin`：Gin Web框架
- `gorm.io/gorm`：GORM数据库ORM库
- `errors`：Go标准库错误处理
- `net/http`：HTTP协议支持

## 5. 设计特点

### 统一错误处理
- 集中处理所有类型的错误，保证响应格式一致性
- 通过 `c.Next()` 确保在处理器链执行完毕后处理错误

### 错误日志记录
- 系统错误会自动记录到日志中便于问题排查
- 使用 `log.Logger.Error` 记录错误详细信息

### 请求终止机制
- 使用 [response.FailStop](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\response\response.go#L39-L45) 确保错误发生时立即终止请求处理
- 防止错误继续传播和处理