# Gorm基础

---

## 目录

- [1. 安装与使用](#1-安装与使用)
- [2. 模型定义与Migrate](#2-模型定义与migrate)
- [3. CRUD操作](#3-crud操作)
  - [3.1 增加数据](#31-增加数据)
  - [3.2 查询数据](#32-查询数据)
  - [3.3 修改数据](#33-修改数据)
  - [3.4 删除数据](#34-删除数据)
- [4. GORM 的一些最佳实践](#4-gorm-的一些最佳实践)

---

## 1. 安装与使用

Gorm常见驱动
- SQLite: gorm.io/driver/sqlite
- MySQL: gorm.io/driver/mysql
- PostgreSQL: gorm.io/driver/postgres

Logger：控制 SQL 日志输出

- logger.Silent：无日志
- logger.Error：仅错误
- logger.Warn：错误和警告
- logger.Info：所有 SQL 查询（开发推荐）

NamingStrategy：自定义命名策略
- TablePrefix：表名前缀
- SingularTable：使用单数表名
- NoLowerCase：禁用自动小写

连接池配置（在 *sql.DB 上配置）：

- SetMaxIdleConns：最大空闲连接数
- SetMaxOpenConns：最大打开连接数
- SetConnMaxLifetime：连接最大生存时间

环境变量配置：通过 .env 文件配置数据库类型和连接信息：
```go
# 数据库类型：sqlite, mysql, postgres
TEST_DB_TYPE=sqlite

# MySQL 连接字符串
TEST_MYSQL_DSN=root:password@tcp(localhost:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local

# PostgreSQL 连接字符串
TEST_POSTGRES_DSN=host=localhost user=postgres password=password dbname=testdb port=5432 sslmode=disable
```
## 2. 模型定义与Migrate

模型定义：

```go
type User struct {
    ID        uint      `gorm:"column:id;type:int;primaryKey;autoIncrement"`
    Name      string    `gorm:"size:64;not null"`
    Email     string    `gorm:"size:128;unique;not null"`
    Age       int       `gorm:"type:int;not null"`
    Status    string    `gorm:"size:16;default:active;index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```
gorm模型定义时，需要使用gorm标签来指定字段的属性，如字段长度、是否为空、默认值、索引等。

常见的gorm标签有：
- primaryKey：主键
- size:64：字段大小
- not null：非空约束
- uniqueIndex：唯一索引
- index：普通索引
- default:value：默认值
- autoCreateTime：自动创建时间
- autoUpdateTime：自动更新时间

Migrate：自动迁移数据库结构，自动迁移数据库结构是指自动在数据库中中创建、修改、删除表结构。

> 注意：不会删除已存在的列，不会修改现有数据，不会删除索引
```go
if err := db.AutoMigrate(&User{}); err != nil {
    fmt.Println(err)
}
```
## 3. CRUD操作

以User模型为例，以增删改查为例
```go
type User struct {
    ID        uint      `gorm:"column:id;type:int;primaryKey;autoIncrement"`
    Name      string    `gorm:"size:64;not null"`
    Email     string    `gorm:"size:128;unique;not null"`
    Age       int       `gorm:"type:int;not null"`
    Status    string    `gorm:"size:16;default:active;index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```
### 3.1 增加数据
- 使用Create方法，批量增加数据
```go
users := []User{
    {Name: "Alice", Email: "alice@example.com", Age: 18, Status: "active"},
    {Name: "Bob", Email: "bob@example.com", Age: 19, Status: "active"},
    {Name: "Charlie", Email: "charlie@example.com", Age: 20, Status: "inactive"},
    {Name: "David", Email: "david@example.com", Age: 21, Status: "active"},
    {Name: "Edward", Email: "edward@example.com", Age: 22, Status: "active"},
    {Name: "Frank", Email: "frank@example.com", Age: 23, Status: "inactive"},
    {Name: "George", Email: "george@example.com", Age: 24, Status: "active"},
    {Name: "Helen", Email: "helen@example.com", Age: 25, Status: "active"},
    {Name: "Irene", Email: "irene@example.com", Age: 26, Status: "active"},
    {Name: "James", Email: "james@example.com", Age: 27, Status: "active"},
    {Name: "Kevin", Email: "kevin@example.com", Age: 28, Status: "active"},
    {Name: "Lucas", Email: "lucas@example.com", Age: 29, Status: "active"},
    }
if err := db.Create(&users).Error; err != nil {
    fmt.Println(err)
}
```
- 使用Create方法，单条增加数据
```go
u := User{Name: "Mike", Email: "mike@example.com", Age: 30, Status: "active"}
if err := db.Create(&u).Error; err != nil {
    fmt.Println(err)
}
```
### 3.2 查询数据
- 条件查询`where("字段=?","值")`： 查询某个字段为某个值的数据
- 查询第一条数据 `First()`

查询第一条status为active的用户
```go 
var user User
if err := db.Where("status = ?", "active").First(&user).Error; err != nil {
    fmt.Println(err)
}
```
- 获取一条记录(不要求条件，找不到不会报错) `Take()`
```go
var user User
if err:=db.Take(&user).Error; err != nil {
    fmt.Println(err)
}
```
- 查询所有匹配条件的记录`Find()`

查询所有status为active的用户
```go
var user []User
if err := db.Where("status = ?", "active").Find(&user).Error; err != nil {
    fmt.Println(err)
}
```

- 扫描数据到自定义结构体中
```go
type UserInfo struct {
    Name string
    Age int
    Status string
}
if err:=db.Model(&User{}).Select("name", "age", "status").Where("status = ?", "active").Scan(&userInfo).Error; err != nil {
    fmt.Println(err)
}
```
### 3.3 修改数据
- 更新所有数据的字段`Save()`
```go
var user User
user.Name = "Mike"
db.Save(&user)
```
- 更新指定的字段`Updates()`(可以使用`Select()`/`map`优化性能)
```go
var user User
// 查询满足条件的数据存放在user中
db.Model(&user).Updates(map[string]any{"name": "Mike","status": "vip"})
```
- 批量更新
```go
db.Model(&User{}).Where("status = ?", "active").Updates(map[string]any{"status": "inactive"})
```
### 3.4 删除数据
- 删除一条数据`Delete()`
```go
var user User
// 查询出user
db.Delete(&user) // 使用instance的方式删除单条数据
db.Delete(&User{}, user.ID) // 使用ID的方式删除单条数据
```
- 批量删除`Delete()`
```go
db.Where("status = ?", "inactive").Delete(&User{})

```
- 验证删除
```go
err := db.First(&User{}, user.ID).Error
if !errors.Is(err, gorm.ErrRecordNotFound) {
    // 记录仍然存在
}
```
## 4. 高级查询运用
### 4.1 Where条件基础查询
- 查询所有数据`db.Where("status = ?", "active").Find(&users)`
```go
var users []User
if err := db.Where("status = ?", "active").Find(&users).Error; err != nil {
    fmt.Println(err)
}
fmt.Printf("basic query: %v\n", users)
```
- 多条件查询`db.Where("status = ? AND age > ?", "active", 18).Find(&users)`
```go
var users []User
if err := db.Where("status = ? and age > ?", "active", 25).Find(&users).Error; err != nil {
    fmt.Printf("query user err: %v", err)
}
fmt.Printf("basic query: %v\n", users)
```
- 模糊查询`db.Where("name LIKE ?", "%jon%").Find(&users)`
```go
db.Where("email LIKE ?", "a%").Find(&users)  // 以 'a' 开头
```
- IN 查询
```go
db.Where("status IN ?", []string{"active", "pending"}).Find(&users)

```
- Between查询
```go
db.Where("age BETWEEN ? AND ?", 20, 30).Find(&users)

```
- Order 排序：`desc`——降序   `asc`——升序
```go
db.Order("created_at desc").Find(&users)  // 降序
db.Order("age asc").Find(&users)          // 升序
```
- Limit 限制查询数量与Offset 偏移量控制分页
```go
pageSize := 5
page := 2
offset := (page - 1) * pageSize
var users []User
if err := db.Order("name desc").Limit(pageSize).Offset(offset).Find(&users).Error; err != nil {
    fmt.Printf("query user err: %v", err)
}
```
### 4.2 查询Scope的应用
Scopes 允许将通用查询条件提取为可复用函数。

Scopes函数基本格式
```go
func funcname(param1, param2, ...) func(db *gorm.DB) *gorm.DB{
    return func(db *gorm.DB) *gorm.DB{
        return db.Where("name = ?", "jinzhu").Find(&User{})
    }
}
```

- 查询active user
```go
func activeUser() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active")
	}
}
// 调用
var users []User
if err := db.Scopes(activeUser()).Find(&users).Error; err != nil {
    fmt.Printf("query user err: %v", err)
}
```
- 分页
```go
func paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(pageSize).Offset(pageSize * (page - 1))
	}
}
// 调用
var users []User
if err := db.Scopes(paginate(1,10)).Find(&users).Error; err != nil {
    fmt.Printf("query user err: %v", err)
}
```
- 多个Scopes组合
```go
if err := db.Scopes(activeUser(),paginate(1,10)).Find(&users).Error; err != nil {
    fmt.Printf("query user err: %v", err)
}
```
### 4.3 聚合查询
- 统计数量 `Count`
```go
var count int64
db.Model(&User{}).Where("status = ?", "active").Count(&count)
```
- 分组统计`Group`
```go
type StatusCount struct {
    Status string
    Total  int64
}
var counts []StatusCount
db.Model(&User{}).Select("status, COUNT(*) as total").Group("status").Scan(&counts)
```
### 4.4 原生sql运用
当链式查询无法满足复杂需求时，可以直接执行原生 SQL。常用方式如下：

- db.Raw(sql, args...).Scan(dest)：查询语句，将结果映射到结构体或切片
- db.Exec(sql, args...)：执行更新、删除、插入等非查询语句
- db.Raw(sql).Rows() / db.Raw(sql).Row()：需要手动遍历 *sql.Rows 或获取单行记录时使用

- 查询示例
```go
type StatusSummary struct {
    Status string
    Total  int64
    AvgAge float64
}

var stats []StatusSummary
start := time.Now().AddDate(0, -1, 0)
end := time.Now()

err := db.Raw(`
    SELECT status, COUNT(*) AS total, AVG(age) AS avg_age
    FROM users
    WHERE created_at BETWEEN ? AND ?
    GROUP BY status
`, start, end).Scan(&stats).Error

if err != nil {
    log.Fatalf("query failed: %v", err)
}
```
- 执行示例
```go
threshold := time.Now().AddDate(0, 0, -30)
result := db.Exec(
    "UPDATE users SET status = ? WHERE last_login_at < ?",
    "inactive",
    threshold,
)

if result.Error != nil {
    log.Fatalf("exec failed: %v", result.Error)
}

fmt.Printf("affected rows: %d\n", result.RowsAffected)
```
## 5. GORM 的一些最佳实践

- ORM 三大优势
    - 生产效率：减少重复代码，提高开发速度
    - 模型同步：代码模型与数据库结构保持一致
    - 组合能力：通过链式调用灵活组合查询
- Create / Save / Updates 的区别：

    - Create：插入新记录
    - Save：保存记录（插入或更新所有字段）
    - Updates：更新指定字段（忽略零值）
- Find vs Scan 的区别：

    - Find：查询到相同模型的结构体
    - Scan：查询到自定义结构体、map 或原始值（用于聚合查询）
- 最佳实践
    - 错误处理：始终检查错误，特别是使用 First 时要检查 gorm.ErrRecordNotFound
    - 连接池配置：根据应用负载合理配置连接池参数
    - 使用 Select：只查询需要的字段，提高性能
    - 使用 Scopes：提取通用查询逻辑，提高代码复用性
    - 环境变量：使用 .env 文件管理配置，便于不同环境切换