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
- [4. 高级查询运用](#4-高级查询运用)
  - [4.1 Where条件基础查询](#41-where条件基础查询)
  - [4.2 查询Scope的应用](#42-查询scope的应用)
  - [4.3 聚合查询](#43-聚合查询)
  - [4.4 原生sql运用](#44-原生sql运用)
- [5. 模型关系](#5-模型关系)
  - [5.1 一对多关系](#51-一对多关系)
  - [5.2 多对一关系](#52-多对一关系)
  - [5.3 多对多关系](#53-多对多关系)
  - [5.4 预加载（Preload）](#54-预加载preload)
  - [5.5 创建带关联的记录](#55-创建带关联的记录)
- [6. 事务](#6-事务)
- [7. 钩子（Hooks）](#7-钩子hooks)
- [8. 软删除](#8-软删除)
- [9. 乐观锁](#9-乐观锁)
- [10. 审计字段](#10-审计字段)
- [11. GORM 的一些最佳实践](#11-gorm-的一些最佳实践)

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
## 5. 模型关系

### 5.1 一对多关系

例如：有一张用户表(user)，对应着有一个用户的详细资料表(profile),用户表和用户详情表是一对一的关系,那么gorm模型可以定义为如下:
```go
type User struct {
    ID        uint   `gorm:"primary_key"`
    Name      string `gorm:"not null"`
    Profile   Profile
    ProfileID uint
    CreatedAt time.Time
    UpdatedAt time.Time
}
typer Profile struct {
    ID        uint   `gorm:"primary_key"`
    Age        int
    Gender    string
    User    User
    UserID     uint
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
通过如上的定义
- GORM 会为 User 和 Profile 两张表自动建表。
- User 表有 profile_id 字段，Profile 表有 user_id 字段。
- User 结构体中的 Profile 字段表示一对一关系（has one），GORM 会自动建立外键关联。
- 你可以通过 db.Preload("Profile").Find(&users) 预加载用户的 Profile 信息。
- 可以通过user.Profile 获取用户的 Profile 信息。
- 可以通过profile.User 获取 Profile 对应的用户信息。
```go
// 查询用户并预加载 Profile
var users []User
db.Preload("Profile").Find(&users)
var profile Profile

var profiles []Profile
db.Preload("User").Find(&profiles)
```
### 5.2 多对一关系
例如：有一张用户表(user)，对应着有一个用户的部门表(dept),部门表和用户详情表是一对一的关系,那么gorm模型可以定义为如下:

```go
type Dept struct {
    ID        uint   `gorm:"primary_key"`
    Name      string `gorm:"not null"`
    Users     []User
    CreatedAt time.Time
    UpdatedAt time.Time
}
type User struct { 
    ID        uint   `gorm:"primary_key"`
    Name      string `gorm:"not null"`
    Dept      Dept
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
通过如上的定义,实现了多对一关系
- GORM 会为 Dept 和 User 两张表自动建表。
- Dept 表有 users 字段，User 表有 dept_id 字段。
- User 模型中的 Dept 字段表示多对一关系（belongs to），GORM 会自动建立外键关联。
- 你可以通过 db.Preload("Dept").Find(&users) 预加载用户的部门信息。
- 获取用户对应的部门信息：user.Dept
- 获取部门对应的用户信息：dept.Users
```go
var users []User
db.Preload("Dept").Find(&users)
var dept Dept
db.Preload("Users").Find(&dept)
```

### 5.3 多对多关系
例如：有一张用户表(user)，对应着有一个用户的角色表(role),角色表和用户详情表是多对多的关系,那么gorm模型可以定义为如下:
```go
type Role struct {
    ID        uint   `gorm:"primary_key"`
    Name      string `gorm:"not null"`
    Users     []User `gorm:"many2many:user_roles;"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
type User struct { 
    ID        uint   `gorm:"primary_key"`
    Name      string `gorm:"not null"`
    Roles     []Role `gorm:"many2many:user_roles;"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
通过如上的定义,实现了多对多关系
- GORM 会为 Role 和 User 两张表自动建表。
- Role 表有 users 字段，User 表有 roles 字段。
- User 模型中的 Roles 字段表示多对多关系（belongs to many），GORM 会自动建立中间表（user_roles）关联。
- 你可以通过 db.Preload("Roles").Find(&users) 预加载用户的角色信息。
- 获取用户对应的角色信息：user.Roles
- 获取角色对应的用户信息：role.Users
```go
var users []User
db.Preload("Roles").Find(&users)
var role Role
db.Preload("Users").Find(&role)
```
### 5.4 预加载（Preload）
预加载（Eager Loading）是在查询主记录时，同时加载关联记录的技术。相比懒加载（Lazy Loading），预加载可以避免 N+1 查询问题。

1. N+1 查询问题：
N+1 查询问题是 ORM 中常见的性能问题，指的是：
- 1 次查询：获取主记录列表（如查询所有用户）
- N 次查询：在循环中为每条主记录查询关联数据（如为每个用户查询订单）
- 总计：1 + N 次数据库查询

当主记录数量（N）很大时，会产生大量数据库查询，严重影响性能。

- 懒加载示例
```go
var users []User
db.Find(&users)
for _, user := range users {
    // 为每个用户单独查询订单（100 次查询）
    db.Model(&user).Association("Orders").Find(&user.Orders)
    // 第 2 次查询：SELECT * FROM orders WHERE user_id = 1
    // 第 3 次查询：SELECT * FROM orders WHERE user_id = 2
    // ...
    // 第 101 次查询：SELECT * FROM orders WHERE user_id = 100
}
// 总共：1 + 100 = 101 次查询！

```
- 性能分析
假设每次数据库查询耗时 10ms：

    - 懒加载：101 次查询 × 10ms = 1010ms（超过 1 秒）

    - 预加载：2 次查询 × 10ms = 20ms（仅 20 毫秒）

    性能提升：50 倍！

- 为什么会出现 N+1 问题？

- 懒加载机制：ORM 默认使用懒加载，只有在访问关联字段时才执行查询
- 循环访问：在循环中访问关联字段，每次循环都会触发一次查询
- 缺乏批量优化：ORM 无法自动识别批量查询场景，无法合并查询

- 预加载如何解决 N+1 问题？

预加载（Eager Loading）的核心思想是：在查询主记录时，一次性批量加载所有关联数据。

    - 预加载解决方案
    ```go
    var users []User
    db.Preload("Orders").Find(&users)
    for _, user := range users {
        fmt.Println(user.Orders)
    }
    ```
    - 预加载工作原理
    1. 查询主记录
    ```SQL
    SELECT * FROM users;
    ```
    2. 批量加载关联数据
    ```SQL
    SELECT * FROM orders WHERE user_id IN (1, 2, 3, ..., 100);
    ```
    3. 在内存中将关联数据匹配到对应的主记录
- 最佳实践
    - ✅ 总是使用 Preload：当需要访问关联数据时，优先使用 Preload
    - ✅ 条件预加载：使用条件预加载减少数据量（如只加载已支付的订单）
    - ❌ 避免循环查询：不要在循环中使用 Association().Find() 查询关联数据
- 预加载单个关联
```go
var user User
db.Preload("Profile").First(&user, 1)
```
- 预加载多个关联
```go
db.Preload("Profile").Preload("Orders").Find(&user, 1)
```
- 条件预加载
```go
db.Preload("Profile", "status = ?", "active").Preload("Orders", "status = ?", "paid").Find(&user, 1)
```
- 预加载关联的嵌套结构
```go
db.Preload("Orders.Items.Product").First(&user, 1)
```
- 预加载所有关联
```go
db.Preload(clause.Associations).First(&user, 1)
```
### 5.5 创建带关联的记录
- 使用Create
```go
user:=User{
    Name:"Alice",
    Profile:Profile{
        Age:30,
        Sex:"F",
        Address:"New York",
        Email:"<EMAIL>",
    },
    Orders:[]Order{
        {Product:"Book", Price:10},
        {Product:"Pen", Price:5},
    },
    Roles:[]Role{
        {Name:"Admin"},
        {Name:"User"},
    },
}
db.Session(&gorm.Session{FullSaveAssociations: true}).Create(&user)
```
- `FullSaveAssociations: true` 的作用：
    - 保存所有关联记录，即使它们是零值
    - 确保嵌套的关联结构被正确创建
    - 自动设置外键关系

- 使用Select/Omit 控制关联的保存
```go
// 只保存指定的关联
db.Select("Profile", "Orders").Create(&user)
// 排除指定的关联
db.Omit("Orders").Create(&user)
```
## 6. 事务

事务用于确保一组数据库操作要么全部成功，要么全部失败，保证数据一致性。
- 典型场景：
    - 转账操作：扣款和加款必须同时成功或失败
    - 订单创建：创建订单和扣减库存必须原子性
    - 批量操作：多个相关操作必须作为一个整体
- 自动事务（推荐）——`Transaction`
```go
err := db.Transaction(func(tx *gorm.DB) error {
    // 所有操作都在事务中
    if err := tx.Create(&order).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", quantisty)).Error; err != nil {
        return err  // 返回错误会自动回滚
    }
    return nil  // 返回 nil 会自动提交
})
```
- 手动事务——`Begin`
```go
tx := db.Begin()
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

if err := tx.Create(&order).Error; err != nil {
    tx.Rollback()
    return err
}

if err := tx.Commit().Error; err != nil {
    tx.Rollback()
    return err
}
```
- SavePoint（保存点）
SavePoint 允许在事务中创建检查点，可以回滚到特定点而不回滚整个事务：
```go
db.Transaction(func(tx *gorm.DB) error {
    // 操作1
    tx.Create(&order)
    
    // 创建保存点
    tx.SavePoint("sp1")
    
    // 操作2
    if err := tx.Create(&item).Error; err != nil {
        // 回滚到保存点（不影响操作1）
        tx.RollbackTo("sp1")
        return err
    }
    
    return nil
})
```
- 嵌套事务
GORM 支持嵌套事务，内层事务失败不会影响外层事务：

```go
db.Transaction(func(tx1 *gorm.DB) error {
    // 外层事务
    tx1.Create(&order)
    
    return tx1.Transaction(func(tx2 *gorm.DB) error {
        // 内层事务
        tx2.Create(&item)
        return nil
    })
})
```
- 幂等性设计
在事务中实现幂等性，防止重复操作：
```go
func transfer(db *gorm.DB, input transferInput) error {
    return db.Transaction(func(tx *gorm.DB) error {
        // 检查是否已存在相同的转账记录（幂等性）
        var existing transferRecord
        if err := tx.Where("reference = ?", input.Reference).First(&existing).Error; err == nil {
            return errDuplicateTransfer  // 已存在，返回错误
        }
        
        // 执行转账逻辑
        // ...
        return nil
    })
}
```
## 7. 钩子（Hooks）
GORM 提供了丰富的钩子（Hooks），可以在执行 CRUD 操作时触发自定义逻辑。
- 钩子类型：
1. 创建钩子：在创建记录时触发
2. 更新钩子：在更新记录时触发
3. 删除钩子：在删除记录时触发
- 钩子函数
```go
func beforeCreate(u *User) (err error) {
    return nil
}
func afterCreate(u *User) (err error) {
    return nil
}
func beforeUpdate(u *User) (err error) {
    return nil
}
func afterUpdate(u *User) (err error) {
    return nil
}
func beforeDelete(u *User) (err error) {
    return nil
}
func afterDelete(u *User) (err error) {
    return nil
}
```

注意：对于软删除，BeforeDelete 钩子会被触发，但 SetColumn 在软删除的 UPDATE 语句中可能不起作用。实际使用时，建议在删除前先更新 deleted_by 字段，然后再执行删除操作。

- 钩子执行顺序
```go
BeforeCreate → SQL INSERT → AfterCreate
BeforeUpdate → SQL UPDATE → AfterUpdate
BeforeDelete → SQL DELETE → AfterDelete
```
## 8. 软删除
GORM 默认支持软删除，即将记录的 deleted_at 字段设置为当前时间，而不是从数据库中删除记录。

- 启用软删除
使用 gorm.DeletedAt 类型启用软删除：
```go
type User struct {
    ID uint
    Name string
    Age int
    Address string
    Email string
    Profile Profile
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt
}
```
- 使用软删除
```go
db.Delete(&user)  // 不会真正删除，而是设置 deleted_at
// SQL: UPDATE user SET deleted_at = NOW() WHERE id = 1
db.First(&user, 1)  // 自动过滤已删除的记录
// SQL: SELECT * FROM articles WHERE id = 1 AND deleted_at IS NULL
// 如果记录已被软删除，会返回 ErrRecordNotFound
db.Unscoped().First(&article, 1)  // 包含已删除的记录
// SQL: SELECT * FROM articles WHERE id = 1
db.Unscoped().Delete(&article)  // 真正删除记录
// SQL: DELETE FROM articles WHERE id = 1
```
- 软删除与 BeforeDelete 钩子
软删除时会触发 BeforeDelete 钩子，但需要注意：
    - SetColumn 的限制：在软删除的 UPDATE 语句中，SetColumn 可能不起作用
    - 推荐做法：在删除前先更新 deleted_by 字段，然后再执行删除操作
```go
// 推荐做法：先更新 deleted_by，再删除
ctx := withOperator("charlie")
db.WithContext(ctx).
    Model(&article{}).
    Where("id = ?", articleID).
    Update("deleted_by", "charlie")
db.WithContext(ctx).Delete(&article{}, articleID)
```
- 软删除执行流程
```go
db.Delete(&article, 1)
  ↓
检查模型是否有 gorm.DeletedAt 字段 → 有
  ↓
调用 SoftDeleteDeleteClause.ModifyStatement
  ↓
检查 Unscoped 标志 → false
  ↓
构建 UPDATE 语句：UPDATE articles SET deleted_at = NOW() WHERE id = 1
```
- 硬删除执行流程（使用 Unscoped()）
```go
db.Unscoped().Delete(&article, 1)
  ↓
设置 Statement.Unscoped = true
  ↓
检查模型是否有 gorm.DeletedAt 字段 → 有
  ↓
调用 SoftDeleteDeleteClause.ModifyStatement
  ↓
检查 Unscoped 标志 → true（跳过软删除逻辑）
  ↓
构建 DELETE 语句：DELETE FROM articles WHERE id = 1
```
## 9. 乐观锁
乐观锁是一种并发控制机制，通过版本号来检测数据是否被其他事务修改。
- 实现乐观锁
使用  `gorm:"version"` 标签启用乐观锁：
```go
type User struct {
    ID        uint      `gorm:"primary_key"`
    Name      string    `gorm:"not null"`
    Age       int       `gorm:"not null"`
    Version   int       `gorm:"version"`  // 乐观锁版本字段
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
- 使用乐观锁
使用 Updates 方法时，需要手动递增版本号，并通过检查 RowsAffected 来判断是否更新成功：
```go
// 第一次更新，版本号从 0 变为 1
if err := db.First(&user, userId).Error; err != nil {
    return err
}

result := db.Model(&user).
    Select("content", "updated_by", "version").
    Updates(map[string]any{
        "content": "新内容",
        "version": gorm.Expr("version + 1"), // 手动递增版本号
    })

if result.Error != nil {
    return result.Error
}
// 版本号已更新为 1

// 使用旧版本号尝试更新（会失败）
stale := User{
    ID: userId,
    Version: 0,
    Name: "旧内容",
}
result = db.Model(&stale).
    Where("version = ?", 0). // 使用旧版本号
    Updates(map[string]any{"content": "尝试使用旧版本更新"})

if result.Error != nil {
    return result.Error
}
// 检查更新的行数，如果版本号不匹配，应该更新 0 行
if result.RowsAffected == 0 {
    return errors.New("optimistic lock conflict: version mismatch")
}
```
- 注意点
    - Updates 方法不会自动递增版本号，需要手动使用 gorm.Expr("version + 1")
    - 版本号不匹配时，更新会影响 0 行，但不会返回错误
    - 需要通过检查 RowsAffected 来判断是否更新成功

- 乐观锁与悲观锁
    - 乐观锁
    1. 假设冲突很少发生
    2. 通过版本号检测冲突
    3. 性能较好，适合读多写少场景
    4. 推荐使用
    - 悲观锁
    1. 假设冲突经常发生
    2. 通过数据库锁机制（SELECT FOR UPDATE）
    3. 性能较低，适合读少写多场景
    4. 适合高并发场景
    5. ⚠️ 不推荐使用，原因如下：
        - 死锁风险：多个事务以不同顺序锁定资源时容易产生死锁
        - 性能问题：
            - 阻塞其他事务：被锁定的记录会阻塞其他需要修改该记录的事务
            - 锁持有时间长：锁会持续到事务结束，如果事务中有慢操作，锁持有时间会更长
            - 并发性能差：高并发场景下，大量事务会排队等待锁释放
            - 资源浪费：即使事务最终可能失败，锁也会一直持有直到事务结束
            - 扩展性差：随着并发量增加，性能会急剧下降
## 10. 审计字段
审计字段用于记录数据的创建人、更新人、删除人等信息，便于追踪数据变更历史。
```go
type AuditFields struct {
    CreatedBy string
    UpdatedBy string
    DeletedBy string
}

type Article struct {
    ID        uint
    Title     string
    Audit     AuditFields `gorm:"embedded"`  // 嵌入审计字段
    CreatedAt time.Time
    UpdatedAt time.Time
}
```
在钩子中设置审计字段
```go
func (a *Article) BeforeCreate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.CreatedBy = user
    a.Audit.UpdatedBy = user
    // 对于 embedded 字段，使用扁平化的字段名（snake_case）
    tx.Statement.SetColumn("created_by", user)
    tx.Statement.SetColumn("updated_by", user)
    return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.UpdatedBy = user
    tx.Statement.SetColumn("updated_by", user)
    return nil
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
    user := getCurrentUser(tx)
    a.Audit.DeletedBy = user
    tx.Statement.SetColumn("deleted_by", user)
    return nil
}
```

## 11. GORM 的一些最佳实践

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