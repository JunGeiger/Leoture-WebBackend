# 项目名称：LeotureWeb
*go版本：1.26.2*

---

### 文档说明：本文档是项目的架构契约与开发规范。所有代码组织、包职责划分、依赖方向均以此为准
### 维护原则：文档与代码同步演进。新增业务模块或者调整包结构时，必须同步更新本文档。

### 顶层目录概览与职责划分
    LeotureWeb
        cmd                 - 项目入口，仅包含启动逻辑  
        config              - 项目配置文件目录，包含casbin配置文件
        docs                - 项目说明文档目录，包含生成的Swagger文档文件
            /swagger
        internal            - 私有业务代码
            /auth           - 权限业务
            /cache          - redis初始化，缓存操作
            /config         - 配置初始化，热加载处理
            /database       - 数据库初始化
            /errors         - 全局错误定义
            /handler        - 接收web请求，请求基础校验，不包含业务校验
            /logger         - 日志初始化
            /middleware     - gin中间件定义，权限、跨域、全局错误处理、全局请求返回日志、全局请求和返回包装
            /model          - 数据库映射
            /repository     - 数据库操作
            /router         - 路由注册
            /scheduler      - 定时任务初始化与注册
                /jobs       - 定时任务业务逻辑
            /service        - 业务逻辑，根据领域划分包，包内包含接口、业务校验、数据转换、业务处理、返回结果包装
            /types          - 请求与返回数据定义
                /request
                /response
            /utils          - 工具类，不可包含其他业务包依赖
            /migrations
            /static
            /tests
            /third_party

### 事务处理规范：
- 嵌套事务：Service 层用 Savepoint
- 网络超时：pgxpool / DSN 配置超时参数
- 业务超时：尊重ctx，在service中使用context.WithTimeout

### 代码规范
- BO 永不出 Service，DTO 永不进 Repository，Model 永不进 Handler
- handler、service、repository 统一注册到set，由主程序统一New

### 单行权限Policy
 - 允许：p, role, /api/*, GET,POST,PUT,DELETE,PATCH, allow
 - 禁止：p, role, /api/*, GET,POST,PUT,DELETE,PATCH, deny
 - 全量删插 + Watcher + 缓存清理
 - p = sub, obj, act, eft    sub=角色, obj=路径, act=HTTP方法, eft=allow/deny
 - g = _, _    用户 -> 角色

### swagger 常见参数
```
数据类型
    string（这包括日期和文件）
    number
    integer
    boolean
    array
    object
```
```
参数类型
    query
    path
    header
    body
    formData
```
```

```
### 减少结构体定义，service调用request结构体自身方法转为model类型结构体，传到repo层保存使用，分页参数是单独结构体
### repo层insert、update只允许传入model结构体，查询结构体和保存结构体定义需要分离
### 为了确保结构体统一可用，需要区分查询和保存时使用的结构体字段需要单独传入
### 其他的外部查询条件都是其他通用结构体，不膨胀model，非通用就要重新定义结构体，尽量使用通用结构体
### gin req参数标签 binding：指定校验规则，是最核心的校验标签。
### 常用规则：required、email、min=1、max=100、gte=0、oneof=admin user、omitempty（为空跳过校验）、-（跳过该字段校验）。