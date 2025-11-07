## 1. 数据结构

### [commentLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\commentLogic.go#L8-L8) 结构体
- **功能**：实现评论相关业务逻辑的结构体
- **特点**：空结构体，仅用于方法绑定

## 2. 全局变量

### [CommentLogic](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\logic\commentLogic.go#L10-L10) 变量
- **类型**：`*commentLogic`
- **功能**：评论业务逻辑的单例实例
- **初始化**：使用 `new(commentLogic)` 创建

## 3. 业务逻辑方法

### [CreateComment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L12-L25) 方法
- **功能**：创建新评论的业务逻辑
- **参数**：
    - `comment *model.Comment`：评论信息
    - `userId *uint`：用户ID
- **主要流程**：
    1. 验证关联的文章是否存在
    2. 设置评论的用户ID
    3. 将评论保存到数据库
- **返回值**：`error` - 可能的错误信息
- **错误处理**：
    - 文章不存在时返回自定义错误
    - 数据库操作失败时返回相应错误

### [CommentByPostId](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\handler\commentHandler.go#L28-L39) 方法
- **功能**：根据文章ID查询所有评论
- **参数**：`postId *string` - 文章ID
- **主要流程**：
    1. 查询数据库中指定文章的所有评论
    2. 返回评论列表
- **返回值**：
    - `*[]model.Comment`：评论列表指针
    - `error`：可能的错误信息

## 4. 依赖组件

### 内部模块
- `mission04/internal/model`：数据模型定义
- `mission04/pkg/db`：数据库访问层

### 标准库
- `errors`：错误处理

## 5. 数据库操作

### 文章存在性验证
- **查询对象**：`model.Post{}`
- **查询条件**：`WHERE id = ?`
- **目的**：确保评论关联的文章存在

### 评论创建
- **操作类型**：`CREATE`
- **操作对象**：[model.Comment](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\internal\model\comment.go#L4-L11)
- **字段设置**：自动设置 [UserID](file://D:\sourcecode\go-projects\src\050.actual_combat\mission04\pkg\auth\jwt.go#L11-L11) 字段

### 评论查询
- **操作类型**：`FIND`
- **查询对象**：`model.Comment{}`
- **查询条件**：`WHERE post_id = ?`
- **结果**：返回匹配条件的评论列表