# Gin基础
## 1. Hello World
- 使用gin.Default()创建一个默认的gin实例
```go
func main() {
	r := gin.Default()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	r.Run(":8080")
}
```
- 使用gin.New()创建一个空的gin实例
```go
func main() {
	r := gin.New()
	r.GET("/hello", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	r.Run(":8080")
}

```
- r.Run()：启动服务器，默认端口 8080
- ctx.JSON()：返回 JSON 响应
- gin.Default()和gin.New()的区别
    - gin.Default()创建的实例会默认使用Logger和Recovery中间件
    - gin.New()创建的实例不会使用任何中间件，需要手动添加需要的中间件
    - gin.Default()和gin.New()创建的实例都会使用默认的错误处理函数
    - gin.Default()和gin.New()创建的实例都会使用默认的日志记录器
    - gin.Default()和gin.New()创建的实例都会使用默认的模板引擎
    - gin.Default()和gin.New()创建的实例都会使用默认的路由器
    - gin.Default()和gin.New()创建的实例都会使用默认的配置

## 2. 路由

HTTP 方法：`GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`OPTIONS`/`HEAD`/`Any`/`Handle`

```go
// 方式1:
func getReqDemo(ctx *gin.Context){}
// 方式2
func getReqDemo() gin.HandlerFunc {
    return func(ctx *gin.Context) {}
}

// 路由注册
r.GET("/hello", getReqDemo())
r.POST("/hello", func (ctx *gin.Context) {

})
```
### 2.1 路径参数携带
- 单个路径参数: `/user/:name`
```go
r.GET("/user/:id", func(ctx *gin.Context) {
    id := ctx.Param("id")
    ctx.JSON(http.StatusOK, gin.H{"user_id": id})
})
```
- 多个路径参数: `/user/:id/profile/:name`
```go
r.GET("/user/:id/profile/:name", func(ctx *gin.Context) {
    id := ctx.Param("id")
    name := ctx.Param("name")
    ctx.JSON(http.StatusOK, gin.H{"user_id": id, "profile_name": name})
})
```
- 通配符参数

通配符参数使用 /* 符号，可以匹配路径中剩余的所有部分（包括多个斜杠）。

通配符参数与普通路径参数的区别：

- :id 只能匹配单个路径段（如 /users/:id 匹配 /users/123，但不匹配 /users/123/posts）
- /*filepath 可以匹配多个路径段（如 /files/*filepath 可以匹配 /files/docs/readme.md

使用场景：

- 文件路径：/files/*filepath 可以匹配 /files/docs/readme.md、/files/images/photo.jpg 等
- 静态资源：/static/*filepath 可以匹配所有静态资源路径
- 代理转发：将剩余路径转发到其他服务

注意:

- 通配符参数必须放在路由路径的最后
- 获取参数时使用 c.Param("filepath")，参数名不包含 * 号
- 获取到的值会包含前导斜杠（如 "/docs/readme.md"）
```go

```

- 查询参数：c.Query("name") `/path?name=value`
```go
r.GET("/search", func(ctx *gin.Context) {
    keyword := ctx.Query("kw")
    page := ctx.DefaultQuery("page", "1")
    size := ctx.DefaultQuery("size", "10") //参数默认值
    ctx.JSON(http.StatusOK, gin.H{
        "keyword": keyword,
        "page":    page,
        "size":    size,
    })
})
```
### 2.2 表单参数（POST 请求）


- 使用 PostForm 获取表单参数: 表单参数通常用于 HTML 表单提交或 application/x-www-form-urlencoded 格式的 POST 请求。适用于HTML 表单、文件上传
> 使用PostForm 获取表单参数
```go
r.POST("/login", func(ctx *gin.Context) {
    uname := ctx.PostForm("username")
    upwd := ctx.PostForm("userpwd")
    remember := ctx.DefaultPostForm("remember", "false")
    ctx.JSON(http.StatusOK, gin.H{
        "username": uname,
        "userpwd":  upwd,
        "remember": remember,
    })
})
```
> 使用ShouldBind 获取参数
```go
type LoginRequest struct {
    Username string `form:"username" binding:"required"`
    Password string `form:"password" binding:"required"`
    Remember bool   `form:"remember"`
}

r.POST("/login", func(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBind(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, gin.H{"message": "Login success"})
})
```
- 使用Json 获取参数: JSON 参数用于接收 application/json 格式的请求体, 适用于RESTful API、前后端分离
```go
// JSON表单参数
/*
    测试：curl -X POST http://localhost:8080/api/users -H "Content-Type: application/json"  -d '{"name":"Alice","email":"alice@example.com","age":25}'
*/
r.POST("/api/users", func(ctx *gin.Context) {
    type CreateUserReq struct {
        Name  string `json:"name" binding:"required"`
        Email string `json:"email" binding:"required"`
        Age   int    `json:"age" binding:"gte=0,lte=130"`
    }
    var req CreateUserReq
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "message": "创建用户成功",
        "id":      1,
        "name":    req.Name,
        "email":   req.Email,
        "age":     req.Age,
    })
})
```
- `ShouldBind`与`ShouldBindJSON`验证标签
    - required：必填
    - email：邮箱格式
    - url：URL 格式
    - min=10：最小值
    - max=100：最大值
    - gte=0：大于等于
    - lte=120：小于等于
    - len=10：长度等于
    oneof=red green blue：枚举值
- `ShouldBind`与`ShouldBindJSON`的区别：
1. `ShouldBind` 会根据请求的 `Content-Type` 自动选择合适的绑定方式，支持表单`（form）`、`JSON`、`XML`、`Query` 等多种格式。例如：`Content-Type` 为 `application/x-www-form-urlencoded` 时绑定表单，为 `application/json` 时绑定 JSON。
2. `ShouldBindJSON` 只会解析并绑定 JSON 格式的请求体（`Content-Type` 必须为 `application/json`），不会处理表单或其他格式的数据。

### 2.3 路由分组

- 基础分组

```go
func main() {
    ...
    v1 := r.Group("/api/v1")
    {
        v1.GET("/users", getUsers)
        v1.GET("/user/:id", getUserByID)
    }
    ...
}
func getUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"users": []gin.H{
			{
				"id":   1,
				"name": "张三",
				"age":  18,
			},
			{
				"id":   2,
				"name": "李四",
				"age":  19,
			},
		},
	})
}
func getUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.JSON(http.StatusOK, gin.H{
		"message": "获取用户详情成功",
		"user": gin.H{
			"id":   id,
			"name": "张三",
			"age":  18,
		},
	})
}
```
- 带中间件的分组
```go
v1 := r.Group("/api/v1")
v1.Use(authMiddleware())
{
    v1.GET("/users", getUsers)
    v1.GET("/user/:id", getUserByID)
}
```

## 3. 中间件
### 3.1 基础中间件
- 创建中间件：
```go
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		fmt.Printf("[%s] %s %d %v\n", c.Request.Method, path, status, latency)
	}
}
```
- 注册中间件
```go
func middlewareRegister(r *gin.Engine) {
	r.Use(loggerMiddleware())
}
```
- 多个中间件的使用
```go
r := gin.Default()
r.Use(middleware1())
r.Use(middleware2())
r.Use(middleware3())

r.GET("/test", handler)
```
执行顺序是：middleware1 → middleware2 → middleware3 → handler → middleware3 → middleware2 → middleware1

- 路由分组中使用中间件
```go
v1 := r.Group("/api/v1")
v1.Use(authMiddleware())
{
    v1.GET("/users", getUsers)
    v1.GET("/user/:id", getUserByID)
}
```
### 3.2 常用中间件
- Logger：日志中间件
```go
func loggerMiddleware() gin.HandlerFunc {
    return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
        return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
            param.ClientIP,
            param.TimeStamp.Format(time.RFC1123),
            param.Method,
            param.Path,
            param.Request.Proto,
            param.StatusCode,
            param.Latency,
            param.Request.UserAgent(),
            param.ErrorMessage,
        )
    })
}
```
- 恢复中间件
恢复中间件（Recovery Middleware）是 Web 应用中必不可少的中间件，用于捕获和处理 panic，防止程序崩溃。

作用：
    
1. 防止程序崩溃：当处理函数（Handler）或中间件发生 panic 时，恢复中间件会捕获它，避免整个服务崩溃
2. 优雅的错误处理：捕获 panic 后，返回统一的错误响应（如 500 错误），而不是直接崩溃或返回空响应
3. 提升服务稳定性：即使某个请求处理出错，也不会影响其他请求，服务可以继续正常运行
4. 生产环境必需：在生产环境中，恢复中间件是必需的，可以避免单个请求错误导致整个服务不可用
```go
func recoveryMiddleware() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        // recovered 是 panic 的值，可以用于日志记录
        // 例如：log.Printf("Panic recovered: %v", recovered)
        
        c.JSON(500, gin.H{
            "error": "Internal Server Error",
        })
        c.Abort()
    })
}
```
Gin 的默认行为：
1. gin.Default()：自动包含日志和恢复中间件
2. gin.New()：不包含任何中间件，需要手动添加恢复中间件

- CORS 中间件
跨域资源共享（CORS）是一种机制，用于允许 Web 应用程序从非同源的域请求资源。
```go
func corsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        if origin != "" {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
            c.Header("Access-Control-Allow-Credentials", "true")
        }
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```
或使用官方 CORS 中间件：
```shell
go get github.com/gin-contrib/cors
```
```go
import "github.com/gin-contrib/cors"

config := cors.DefaultConfig()
config.AllowOrigins = []string{"http://localhost:3000"}
config.AllowCredentials = true
r.Use(cors.New(config))
```
- 限流中间件
```go
import "golang.org/x/time/rate"

func rateLimitMiddleware() gin.HandlerFunc {
    limiter := rate.NewLimiter(10, 20) // 每秒10个请求，突发20个
    
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Too many requests"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```
### 3.3. JWT认证中间件
- JWT Token生成
```go
var jwtSecret = []byte("your-secret-key")

type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
func generateJWT(userId int, userName string) (string, error) {
	claims := JWTClaims{
		UserID:   uint(userId),
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(token)
}
```
- JWT Token验证
```go
func parseJWT(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Invalid token")
}
```
- 认证中间件
```go
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{"error": "未授权"})
			c.Abort()
			return
		}
		if auth[:7] != "Bearer " {
			c.JSON(401, gin.H{"error": "无效的授权头"})
			c.Abort()
			return
		}
		tokenStr := auth[7:]
		claims, err := parseJWT(tokenStr)
		if err != nil {
			c.JSON(401, gin.H{"error": "无效的token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
```
## 4. 文件上传与下载
- 单文件上传
```go
func uploadSingleFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	filename := file.Filename
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "文件上传成功", "filename": filename})
}
```
- 多文件上传
```go
func uploadMutipleFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	files := form.File["files"]
	var filenames []string
	for _, file := range files {
		filename := file.Filename
		dst := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		filenames = append(filenames, filename)
	}
	c.JSON(200, gin.H{"message": "文件上传成功", "filenames": filenames})
}
```
- 文件下载
```go
func downloadFile(c *gin.Context) {
    filename := c.Param("filename")
    filepath := filepath.Join("uploads", filename)
    
    // 检查文件是否存在
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"error": "File not found"})
        return
    }
    
    // 设置响应头
    c.Header("Content-Description", "File Transfer")
    c.Header("Content-Transfer-Encoding", "binary")
    c.Header("Content-Disposition", "attachment; filename="+filename)
    c.Header("Content-Type", "application/octet-stream")
    
    // 返回文件
    c.File(filepath)
}
```
- 静态文件服务
```go
r.Static("/uploads", "./uploads") // 提供静态文件服务
r.StaticFS("/files", http.Dir("./uploads")) // 提供静态文件系统 curl http://localhost:8080/files/test.txt

```

1. 文件大小限制：Gin 默认最大请求体大小为 32MB，可通过 MaxMultipartMemory 调整
2. 文件类型验证：生产环境应添加文件类型和大小验证
3. 文件名安全：使用 filepath.Base() 防止目录遍历攻击
4. 存储路径：生产环境建议使用对象存储服务（如 S3、OSS）
5. 权限控制：添加认证中间件保护上传和下载接口
## 5. Viper配置管理
Viper 是 Go 语言的配置管理库，支持：
- JSON、TOML、YAML、HCL、envfile、Java properties
- 环境变量
- 命令行参数
- 远程配置系统（etcd、Consul）
- 配置热加载

### 5.1 基础使用
```go
func init() {
	viper.SetConfigName("config") // 配置文件名（不带扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(".")      // 配置文件路径
	viper.AddConfigPath("$HOME/.app")
	if err := viper.ReadInConfig(); err != nil {
		panic("读取配置文件失败: " + err.Error())
	}
}
func main() {
	port := viper.GetString("server.port")
	dbhost := viper.GetString("database.host")
	fmt.Printf("Server will run on port: %s\n", port)
	fmt.Printf("Database host: %s\n", dbhost)
}
```
### 5.2 配置文件加载到环境变量中
- 使用 SetEnvPrefix（自动映射）
```go
func loadConfigToEnv() {
	viper.AutomaticEnv()      // 从环境变量加载配置
	viper.SetEnvPrefix("APP") // 环境变量前缀，例如 APP_SERVER_PORT,设置前缀后，Viper 会自动将配置键转换为环境变量名，例如 server.port 会转换为 APP_SERVER_PORT
}
```
- 使用 BindEnv
```go
viper.BindEnv("server.port", "PORT")
```
如：
```shell
export PORT=9000
```

```go
port := viper.GetString("server.port")  // 从 PORT 环境变量读取，值为 "9000"
```
### 5.3 使用配置结构体
- 如：有以下yaml文件
```yaml
server:
  port: 8080
  host: "0.0.0.0"
  mode: "debug"

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"
  dbname: "mydb"

jwt:
  secret: "your-secret-key"
  expire: 24h
```
- 配置结构体
```go
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}
type ServerConfig struct {
	Port string `mapstructure:"port"`
	Host string `mapstructure:"host"`
	Mode string `mapstructure:"mode"`
}
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
}

func NewConfig() *Config {
	return &Config{}
}
func (c *Config) LoadConfig(path string) error {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}
```
## 6. Protocol Buffers(Protobuf)
Protobuf 是 Google 的一种语言无关、平台无关、可扩展的序列化结构格式，用于在网络传输中传输结构化数据。它具有以下特点：
- 高效：比 JSON/XML 更小、更快
- 跨语言：支持多种编程语言
- 类型安全：强类型定义，编译时检查
- 向后兼容：支持字段版本演进
- 广泛应用：gRPC、微服务通信、数据存储
### 6.1 安装 protoc 编译器
```shell
# macOS
brew install protobuf

# Linux
apt-get install protobuf-compiler

// go插件安装
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
// 依赖安装
go get google.golang.org/protobuf/proto

```

- 创建 proto 文件(如user.proto)
```proto
syntax="proto3";
package user;
option go_package="./pb";
// 用户自定义消息
message User{
    int64 id=1;
    string username=2;
    string email=3;
    int32 age=4;
    bool active=5;
    repeated string tags=6;
    map<string,string> metadata=7;
}

message UserList{
    repeated User users=1;
    int32 total=2;
}

// 创建用户请求
message CreateUserRequest{
    string username=1;
    string email=2;
    int32 age=3;
}

message CreateUserResponse{
    User user=1;
    bool success=2;
    string message=3;
}
```
- 生成 Go 代码
```shell
protoc --go_out=. --go_opt=paths=source_relative user.proto
```
- 使用grpc生成
```shell
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto
```
将会生成user.pb.go

## 7. GRPC 服务
### 7.1 GRPC简介
gRPC（gRPC Remote Procedure Calls）是 Google 开发的高性能、开源的 RPC 框架，具有以下特点：
- 高性能：基于 HTTP/2，支持多路复用、流式传输
- 跨语言：支持多种编程语言（Go、Java、Python、C++ 等）
- 类型安全：使用 Protobuf 定义接口，编译时检查
- 流式传输：支持服务端流、客户端流、双向流
- 自动代码生成：从 .proto 文件自动生成客户端和服务端代码

应用场景：
- 微服务通信
- 高性能 API
- 实时数据推送
- 跨语言服务调用

安装grpc
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```
### 7.2 创建一个微服务服务端
- user.proto
```go
syntax = "proto3";
package user;
option go_package="./pb";
service UserService{
    rpc GetUser(GetUserRequest) returns (User);
    rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
    rpc ListUser(ListUserRequest) returns (ListUserResponse);
    rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
    rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
    // 流式获取用户（服务端流）
    rpc StreamUsers(StreamUsersRequest) returns (stream User);
    
    // 批量创建用户（客户端流）
    rpc BatchCreateUsers(stream CreateUserRequest) returns (BatchCreateUsersResponse);
    
    // 双向流（聊天式交互）
    rpc ChatUsers(stream ChatMessage) returns (stream ChatMessage);
}

message User{ 
    int64 id=1;
    string username=2;
    string email=3;
    int32 age=4;
    bool active=5;
}
message GetUserRequest{
    int64 id=1;
}
message CreateUserRequest{
    string username=1;
    string email=2;
    int32 age=3;
    repeated string tags=4;
    map<string,string> metadata=5;
}
message CreateUserResponse{
    User user=1;
    bool success=2;
    string message=3;
}
message ListUserRequest{
    int32 page=1;
    int32 page_size=10;
    string filter=3;
}
message ListUserResponse{
    repeated User users=1;
    int32 total=2;
    int32 page=3;
    int32 page_size=4;
}
message UpdateUserRequest{
    int32 id=1;
    string username=2;
    string email=3;
    int32 age=4;
    bool active=5;
    repeated string tags=6;
    map<string,string> metadata=7;
}
message UpdateUserResponse{
    User user=1;
    bool success=2;
    string message=3;
}
message DeleteUserRequest{
    int32 id=1;
}
message DeleteUserResponse{
    bool success=1;
    string message=2;
}
// 流式获取用户请求
message StreamUsersRequest {
  int32 limit = 1;
  int32 interval_ms = 2;  // 发送间隔（毫秒）
}

// 批量创建用户响应
message BatchCreateUsersResponse {
  repeated User users = 1;
  int32 success_count = 2;
  int32 fail_count = 3;
  string message = 4;
}

// 聊天消息（用于双向流示例）
message ChatMessage {
  string user_id = 1;
  string message = 2;
  int64 timestamp = 3;
}

```
- 生成代码
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       user.proto
```
- 实现服务Service
```go
type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	user   map[int64]*pb.User
	mu     sync.RWMutex
	nextID int64
}

func NewUserServiceServer() *UserServiceServer {
	return &UserServiceServer{
		user: make(map[int64]*pb.User),
	}
}
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.user[req.Id]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "user %d not found", req.Id)
	}
	return user, nil
}
....
```
在`pb.UnimplementedUserServiceServer`结构体中包含了刚才`proto`定义的所有的`service`，我们仅仅需要实现想要的方法即可

`NewUserServiceServer()`方法适用于创建一个UserServiceServer实例

`func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error)`是对结构体中继承的方法的实现

- 注册服务
```go
func startServer(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	grpcServer := grpc.NewServer()

	// 注册服务
	userService := NewUserServiceServer()
	pb.RegisterUserServiceServer(grpcServer, userService)
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}
	return nil
}
```
- 启动服务
```go
func main() {
	log.Println("Starting gRPC server...")

	if err := startServer(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
```
### 7.3 服务端常见的四种模式
- 一元 RPC（Unary RPC）：客户端发送一个请求，服务端返回一个响应：
```go
// 服务端
func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // 处理请求并返回响应
    return user, nil
}

// 客户端
user, err := client.GetUser(ctx, &pb.GetUserRequest{Id: 1})
```
- 服务端流（Server Streaming）: 客户端发送一个请求，服务端返回一个流：
```go
// 服务端
func (s *UserServiceServer) StreamUsers(req *pb.StreamUsersRequest, stream pb.UserService_StreamUsersServer) error {
    for i := 0; i < int(req.Limit); i++ {
        user := &pb.User{Id: int64(i + 1), Username: fmt.Sprintf("user%d", i+1)}
        if err := stream.Send(user); err != nil {
            return err
        }
        time.Sleep(time.Duration(req.IntervalMs) * time.Millisecond)
    }
    return nil
}

// 客户端
stream, err := client.StreamUsers(ctx, &pb.StreamUsersRequest{Limit: 10})
for {
    user, err := stream.Recv()
    if err == io.EOF {
        break
    }
    if err != nil {
        log.Fatalf("StreamUsers failed: %v", err)
    }
    log.Printf("Received user: %+v", user)
}
```
- 客户端流（Client Streaming）: 客户端发送一个流，服务端返回一个响应
```go
// 服务端
func (s *UserServiceServer) BatchCreateUsers(stream pb.UserService_BatchCreateUsersServer) error {
    var users []*pb.User
    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        // 处理请求
        user := &pb.User{Username: req.Username, Email: req.Email}
        users = append(users, user)
    }
    return stream.SendAndClose(&pb.BatchCreateUsersResponse{
        Users: users,
        SuccessCount: int32(len(users)),
    })
}

// 客户端
stream, err := client.BatchCreateUsers(ctx)
for _, userReq := range userRequests {
    if err := stream.Send(userReq); err != nil {
        log.Fatalf("Send failed: %v", err)
    }
}
resp, err := stream.CloseAndRecv()

```
- 双向流（Bidirectional Streaming）: 客户端和服务端同时发送一个流
```go
// 服务端
func (s *UserServiceServer) ChatUsers(stream pb.UserService_ChatUsersServer) error {
    for {
        msg, err := stream.Recv()
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        // 处理消息并回复
        reply := &pb.ChatMessage{
            UserId:    "server",
            Message:   "Echo: " + msg.Message,
            Timestamp: time.Now().Unix(),
        }
        if err := stream.Send(reply); err != nil {
            return err
        }
    }
}

// 客户端
stream, err := client.ChatUsers(ctx)
// 发送消息
stream.Send(&pb.ChatMessage{Message: "Hello"})
// 接收消息
msg, err := stream.Recv()
```
## 8. 最佳实践
- 在生产环境中，必须使用恢复中间件
- 可以在恢复中间件中记录 panic 信息，便于排查问题
- 返回统一的错误格式，不要暴露内部错误详情给客户端
