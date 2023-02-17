### 1：什么是orm

|                  orm是一种术语而并非软件，                   |
| :----------------------------------------------------------: |
| orm的全称是（object relational mapping）就是对象映射关系程序 |
| 简单来说类似于Python这种面向对象的程序来说一切皆对象，但是我们使用的数据库却是关系型数据库 |
| 为了保证一致的使用习惯，通过orm将编程语言的对象模型和数据库的关系模型建立映射关系 |
| 这样我们就直接使用编程语言的对象模型进行操作数据库就可以了，而不用去写SQL语言了 |

### 2：什么是gorm

```shell
参考文档：https://gorm.io/zh_CN/docs/index.html
```

|     gorm是一个神奇的，对开发人员友好的Golang ORM库      |
| :-----------------------------------------------------: |
|              1：全特性（几乎包含所有特性）              |
| 2：模型关系（一对一，一对多，多对一，多对多，多态关联） |
| 3：钩子（Before/After，Create/Save/Update/Delete/Find） |
|                        4：预加载                        |
|                         5：事务                         |
|                       6：复合主键                       |
|                      7：SQL构造器                       |
|                       8：自动迁移                       |
|                         9：日志                         |
|             10：基于GORM回调编写可扩展插件              |
|                   11：全特性测试覆盖                    |
|                    12：对开发者友好                     |



### 3：gorm（v2）基本使用

```shell
安装：go get -u gorm.io/gorm

下面我们来看我们的第一个GORM程序，前提是我们得有一个MYSQL的数据库服务，我这里是Docker启动的一个MySQL的服务
账号：root
密码：123456

# 其次我们需要去创建一个提供给gorm链接的数据库

mysql> create database gorm_service charset utf8;

# 然后创建项目初始化项目
go mod init gorm_service

# 创建main.go并开始写我们的第一个gorm程序
```

```go
package main

import (
	"fmt"
	// 我们使用gorm的时候记得银日驱动包
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 连接成功
	fmt.Println("连接成功", db)
}
```

```shell
# 这样我们的第一个gorm程序就写好了，当然这个时候我们拿到的结果是一个内存地址，我们可以执行一下看看结果

PS E:\code\goland\gorm_service> go run .\main.go              
连接成功 &{0xc00011a6c0 <nil> 0 0xc0001c8380 1}

# 这样我们的第一个程序写好之后我们就可以来操作创建表了，前提是我们前面说了，它是Golang对象和数据库对象的映射，所以我们要有一个结构体映射给数据库，所以我们来看下面的代码
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})
}
```

```shell
# 我们尝试执行程序并去查看数据库

PS E:\code\goland\gorm_service> go run .\main.go
# 注：我们没有处理错误，所以没有返回是正常的，直接去查看数据库就好

MySQL root@10.0.0.11:(none)> show databases;
+--------------------+
| Database           |
+--------------------+
| information_schema |
| gorm_service       |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
5 rows in set
Time: 0.008s

MySQL root@10.0.0.11:(none)> use gorm_service;
You are now connected to database "gorm_service" as user "root"
Time: 0.000s
MySQL root@10.0.0.11:gorm_service> show tables;
+------------------------+
| Tables_in_gorm_service |
+------------------------+
| users                  |
+------------------------+
1 row in set
Time: 0.004s

MySQL root@10.0.0.11:gorm_service> select * from users;
+----+------+----------+----------+
| id | name | username | password |
+----+------+----------+----------+
+----+------+----------+----------+
0 rows in set
Time: 0.004s

# 可以看到和我们结构体定义的完全一致

# 当然了，这里我们还会提到的一个就是自定义表
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// TableName 自定义表名
func (*Users) TableName() string {
	return "projects"
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})
}
```

```shell
# 随后我们来查看结果

PS E:\code\goland\gorm_service> go run .\main.go

MySQL root@10.0.0.11:(none)> show tables from gorm_service;
+------------------------+
| Tables_in_gorm_service |
+------------------------+
| projects               |
| users                  |
+------------------------+
2 rows in set
Time: 0.008s
MySQL root@10.0.0.11:(none)>

# 可以看到结果和我们定义的一致
```

### 4：基于gorm的crud

```shell
# 上面我们学习了如何去创建表以及自定义表，然后我们下面就开始针对表进行增删改查的操作了，这里说明一下，增删改查也就是在业内说的crud
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// TableName 自定义表名
func (*Users) TableName() string {
	return "projects"
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	// 新增数据，因为我们讲了，我们用GORM只需要关注结构体的定义，所以我们只需要定义好结构体，然后传入结构体就可以了
	// 不传ID，因为ID是自增的主键
	db.Create(&Users{
		Name:     "张三",
		Username: "admin",
		Password: "123456",
	})
}
```

```shell
# 随后我们来查看一下执行并查看一下数据是否写入数据，不过注意，这里我们如果要写到users，需要把自定义表明那段代码注释哦

MySQL root@10.0.0.11:gorm_service> select * from projects;
+----+------+----------+----------+
| id | name | username | password |
+----+------+----------+----------+
| 1  | 张三 | admin    | 123456   |
+----+------+----------+----------+
1 row in set
Time: 0.005s

# 我这里是直接写到了那个表里了

# 这是更正后的结果
MySQL root@10.0.0.11:gorm_service> select * from users;
+----+------+----------+----------+
| id | name | username | password |
+----+------+----------+----------+
| 1  | 张三 | admin    | 123456   |
+----+------+----------+----------+
1 row in set
Time: 0.005s

# 当然了实例化对象的方式有多种，所以就等于我们去往表里写数据的方式也就是多种
```

```shell
# 下面我们来做的就是修改这些数据，那么我们如何对数据进行修改呢？我们来看代码
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	// 新增数据，因为我们讲了，我们用GORM只需要关注结构体的定义，所以我们只需要定义好结构体，然后传入结构体就可以了
	// 不传ID，因为ID是自增的主键
	db.Create(&Users{
		Name:     "wangwu",
		Username: "wangwu",
		Password: "123456",
	})

	// 修改数据(根据条件查询出数据，然后对指定的字段进行修改)
	db.Where("id = ?", 1).Updates(&Users{
		Name: "zhangsan",
	})
}
```

```shell
PS E:\code\goland\gorm_service> go run .\main.go    
# 我们前面的数据就看到我的name其实是中文的，而后我们用更新字段的方式变更成了英文

MySQL root@10.0.0.11:gorm_service> select * from users;
+----+----------+----------+-----------------+
| id | name     | username | password        |
+----+----------+----------+-----------------+
| 1  | zhangsan | admin    | 123456          |
| 2  | lisi     | lisi     | woshilsii123456 |
| 3  | wangwu   | wangwu   | 123456          |
| 4  | wangwu   | wangwu   | 123456          |
+----+----------+----------+-----------------+

4 rows in set
Time: 0.005s

# 当然这里有个问题，就是我们不一定想修改这个表内的数据，怎么办？当然这个gorm早就给我们想到了，下面看代码
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	// 新增数据，因为我们讲了，我们用GORM只需要关注结构体的定义，所以我们只需要定义好结构体，然后传入结构体就可以了
	// 不传ID，因为ID是自增的主键
	// db.Create(&Users{
	//	Name:     "wangwu",
	//	Username: "wangwu",
	//	Password: "123456",
	// })

	// 修改数据(根据条件查询出数据，然后对指定的字段进行修改)，其实好像并没有特别大的变化，但是其实这个已经指定了是Users表了，因为Users结构体对应了Users表，所以我们只需要给它Users结构体就OK了
	db.Model(&Users{}).Where("id = ?", 1).Updates(&Users{
		Name: "zhangsan",
	})
    
    // 当然这样写也行，这个时候的Model就等于是Where去找Id为4的数据然后更新数据了
    db.Model(&Users{Id: 4}).Update("name", "zhaoliu")
    
    // 这样操作也行
    db.Table("users").Where("id = ?", 1).Updates(&Users{
		Name: "zhangsan",
    })
}
```

```shell
# 我们来看看结果
# 我们通常使用第一种和第三种方法去更新数据。
MySQL root@10.0.0.11:gorm_service> select * from users;
+----+----------+----------+-----------------+
| id | name     | username | password        |
+----+----------+----------+-----------------+
| 1  | zhangsan | admin    | 123456          |
| 2  | lisi     | lisi     | woshilsii123456 |
| 3  | wangwu   | wangwu   | 123456          |
| 4  | zhaoliu  | wangwu   | 123456          |
| 5  | wangwu   | wangwu   | 123456          |
| 6  | wangwu   | wangwu   | 123456          |
| 7  | wangwu   | wangwu   | 123456          |
+----+----------+----------+-----------------+

7 rows in set
Time: 0.005s
```

```shell
# 下面我们来看的就是查询数据，我们一般主要是根据Id来查询数据，我们来看看GORM会如何帮我们实现这个操作
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	// 查询数据-1
	info := &Users{Id: 1}
	db.First(info)
	fmt.Println(info)
    
	// 查询数据-2
	data := &Users{}
	db.Where("id= ?", 3).First(data)
	fmt.Println(data)
    
    // 查询所有数据
	var users []Users
	db.Find(&users)
	fmt.Println(users)
}
```

```shell
# 我们查看一下运行结果

PS E:\code\goland\gorm_service> go run .\main.go
&{1 zhangge admin 123456}
&{3 wangwu wangwu 123456}
[{1 zhangge admin 123456} {2 lisi lisi woshilsii123456} {3 wangwu wangwu 123456} {4 zhaoliu wangwu 123456} {5 wangwu wangwu 123456} {6 wangwu wangwu 123456} {7 wangwu wangwu 123456} {8 wangwu wangwu 123456}]


# 可以看到数据全部被列出来了

# 最后我们要看的就是删除数据了，直接看代码
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	// 删除数据-1
	db.Delete(&Users{
		Id: 6,
	})

	// 删除数据-2
	db.Where("id = ?", 5).Delete(&Users{})

	// 删除数据-3
	db.First(&Users{}, 7).Delete(&Users{})

	// 删除数据-4
	db.Find(&Users{}, "id = ?", 8).Delete(&Users{})
}
```

```shell
# 查看运行结果

PS E:\code\goland\gorm_service> go run .\main.go

# 可以看到5-8都没了

MySQL root@10.0.0.11:gorm_service> select * from users;
+----+---------+----------+-----------------+
| id | name    | username | password        |
+----+---------+----------+-----------------+
| 1  | zhangge | admin    | 123456          |
| 2  | lisi    | lisi     | woshilsii123456 |
| 3  | wangwu  | wangwu   | 123456          |
| 4  | zhaoliu | wangwu   | 123456          |
+----+---------+----------+-----------------+

4 rows in set
Time: 0.006s

# 这样增删改查我们就写完了，其次我们需要了解下它的错误处理，前面我们是没有用到错误处理的。
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type Users struct {
	Id       int64  `gorm:"primary_key" json:"id"`
	Name     string `gorm:"type:varchar(20);not null" json:"name"`
	Username string `gorm:"type:varchar(15);not null" json:"username"`
	Password string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&Users{})

	data := db.Where("id = ?", 1).First(&Users{})
	// 记住一个问题，我们查询的数据结果如果为空也是error，所以我们需要用如下方法判断
	if data.Error != nil && !errors.Is(data.Error, gorm.ErrRecordNotFound) {
		fmt.Println(data.Error)
	} else {
		fmt.Println("查询成功")
	}
}
```

### 5：gorm模型定义

```shell
参考连接：https://gorm.io/zh_CN/docs/models.html

模型一般都是普通的Golang结构体，Go的基本数据类型，或者指针，我们下面来简单的看下例子
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// User Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type User struct {
	Id           int64      `gorm:"primary_key" json:"id"`                  // 设置主键
	CreatedAt    *time.Time `gorm:"column:created_at" json:"createdAt"`     // 创建时间
	Email        string     `gorm:"type:varchar(30);not null;unique_index"` // 唯一索引
	Role         string     `gorm:"size:255"`                               //设置字段大小为255
	MemberNumber string     `gorm:"unique:not null"`                        // 设置 member_number 字段唯一且不为空
	Num          int        `gorm:"AUTO_INCREMENT"`                         // 设置 num 为自增类型
	Address      string     `gorm:"index:addr"`                             // 给address字段创建名为addr的索引
	IgnoreMe     int        `gorm:"-"`                                      // 忽略这个字段
	Name         string     `gorm:"type:varchar(20);not null" json:"name"`
	Username     string     `gorm:"type:varchar(15);not null" json:"username"`
	Password     string     `gorm:"type:varchar(30);not null" json:"password"`
}

func main() {

	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})
}
```

```shell
# 我们执行完来看看表的信息是不是和我们定义的一模一样

MySQL root@10.0.0.11:gorm_service> desc users;
+---------------+--------------+------+-----+---------+----------------+
| Field         | Type         | Null | Key | Default | Extra          |
+---------------+--------------+------+-----+---------+----------------+
| id            | bigint(20)   | NO   | PRI | <null>  | auto_increment |
| created_at    | datetime(3)  | YES  |     | <null>  |                |
| email         | varchar(30)  | NO   |     | <null>  |                |
| role          | varchar(255) | YES  |     | <null>  |                |
| member_number | varchar(191) | YES  | UNI | <null>  |                |
| num           | bigint(20)   | YES  |     | <null>  |                |
| address       | varchar(191) | YES  | MUL | <null>  |                |
| name          | varchar(20)  | NO   |     | <null>  |                |
| username      | varchar(15)  | NO   |     | <null>  |                |
| password      | varchar(30)  | NO   |     | <null>  |                |
+---------------+--------------+------+-----+---------+----------------+

10 rows in set
Time: 0.006s

# 当然了 gorm帮我们提供了一个普遍都会用的字段，我们来看看
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User Users 定义与数据库的映射关系，这里定义的标签就相当于数据库中的字段
type User struct {
	gorm.Model          // gorm.Model 包含了ID、CreatedAt、UpdatedAt、DeletedAt四个字段
	Email        string `gorm:"type:varchar(30);not null;unique_index"` // 唯一索引
	Role         string `gorm:"size:255"`                               //设置字段大小为255
	MemberNumber string `gorm:"unique:not null"`                        // 设置 member_number 字段唯一且不为空
	Num          int    `gorm:"AUTO_INCREMENT"`                         // 设置 num 为自增类型
	Address      string `gorm:"index:addr"`                             // 给address字段创建名为addr的索引
	IgnoreMe     int    `gorm:"-"`                                      // 忽略这个字段
	Name         string `gorm:"type:varchar(20);not null" json:"name"`
	Username     string `gorm:"type:varchar(15);not null" json:"username"`
	Password     string `gorm:"type:varchar(30);not null" json:"password"`
}

// 初始化数据库
func main() {
	// 定义连接地址
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"
	// 连接数据库
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// 使用GORM的自动创建表功能,但是它只是创建表，并没有建立映射关系
	db.AutoMigrate(&User{})
}
```

### 6：gorm一对多

```shell
参考文档：https://gorm.io/zh_CN/docs/has_many.html

# 首先我们来学习的就是一对多入门，我们一般称之为（has many），那么我们来看看什么是一对多。

1：has many 关联就是创建和另一个模型的一对多关系
2：例如，一个人有多个信用卡，这就是生活中常见的一对多的关系
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	CreditCards []CreditCard
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint // 默认会在CreditCard表中添加user_id字段，与User表中的ID字段关联外键
}

func main() {
    
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

}
```

```shell
# 这里面我们提到了外键，那么什么是外键呢
1：为了定义一对多关系，外键是必须存在的，默认外键的名字是所有者名字加上它的主键（UserID）
2：可以看到我门上面的例子，定义了一个User的模型那么CrditCard里面去定义外键就应该是UserID。
3：使用其他字段名作为外键，你可以通过标签使用foreignkey来定制它，例如如下代码
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name       string
	CreditCards []CreditCard ``
}

type CreditCard struct {
	gorm.Model
	Number   string
	UserRefer uint 
}

func main() {
    
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

}
```

```shell
# 这个时候其实UserRefer拿到的值还是与UserID关联的，我们还需要配置一个值，然后我们来看看外键关联
1：GORM通常使用所有者的主键作为外键的值，在上面的例子中默认用的就是UserID。
2：当你分配一张Card给用户时，GORM将用户ID存到信用卡的UserID字段中。
3：你可以通过association_foreignkey来改变它

# 这也就是说我们不一定要用UserID来作为外键，可以自定义外键，下面我们来看代码
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	CreditCard   []CreditCard `gorm:"foreignKey:UserMemberNumber;references:MemberNumber"`
	MemberNumber string
}

type CreditCard struct {
	gorm.Model
	Member           string
	UserMemberNumber string
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})
}
```

```shell
# 执行得到如下结果

MySQL root@10.0.0.11:gorm_service> show tables;
+------------------------+
| Tables_in_gorm_service |
+------------------------+
| credit_cards           |
| users                  |
+------------------------+
2 rows in set
Time: 0.005s

MySQL root@10.0.0.11:gorm_service> select * from credit_cards;
+----+------------+------------+------------+--------+--------------------+
| id | created_at | updated_at | deleted_at | member | user_member_number |
+----+------------+------------+------------+--------+--------------------+
+----+------------+------------+------------+--------+--------------------+

0 rows in set
Time: 0.005s
MySQL root@10.0.0.11:gorm_service> select * from users;
+----+------------+------------+------------+---------------+
| id | created_at | updated_at | deleted_at | member_number |
+----+------------+------------+------------+---------------+
+----+------------+------------+------------+---------------+

0 rows in set
Time: 0.006s


# 那么我们知道了这些之后就开始搞看看如何创建数据并且如何关联查询数据了
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username"`
	CreditCard []CreditCard `gorm:"constraint:OpUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

	// 创建数据
	db.Create(&User{
		Username: "zhangsan",
		CreditCard: []CreditCard{
			{Number: "000001"},
		},
	})
}
```

```shell
# 执行后查看数据

MySQL root@10.0.0.11:gorm_service> select * from users;
+----+----------------------------+----------------------------+------------+----------+
| id | created_at                 | updated_at                 | deleted_at | username |
+----+----------------------------+----------------------------+------------+----------+
| 1  | 2023-02-16 01:38:41.390000 | 2023-02-16 01:38:41.390000 | <null>     | zhangsan |
+----+----------------------------+----------------------------+------------+----------+

1 row in set
Time: 0.005s
MySQL root@10.0.0.11:gorm_service> select * from credit_cards;
+----+----------------------------+----------------------------+------------+--------+---------+
| id | created_at                 | updated_at                 | deleted_at | number | user_id |
+----+----------------------------+----------------------------+------------+--------+---------+
| 1  | 2023-02-16 01:38:41.411000 | 2023-02-16 01:38:41.411000 | <null>     | 000001 | 1       |
+----+----------------------------+----------------------------+------------+--------+---------+

# 可以看到这里的关联就生效了，看User表的ID，然后再看CreditCard的表的user_id，发现是一致的，这就说明这个一对多的关系是成功了，随后我们去查询数据
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username"`
	CreditCard []CreditCard `gorm:"constraint:OpUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

	// 创建数据
	db.Create(&User{
		Username: "zhangsan",
		CreditCard: []CreditCard{
			{Number: "000001"},
		},
	})

	// 查询数据
	u := &User{
		Username: "zhangsan",
	}
	db.First(&u)
	fmt.Println(u)
}
```

```shell
# 查看执行结果

PS E:\code\goland\gorm_service> go run .\main.go
&{{1 2023-02-16 01:49:13.887 +0800 CST 2023-02-16 01:49:13.887 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} zhangsan []}

# 这个时候其实还是有缺陷的，那么既然是关联表，有没有一个操作是专门针对关联表增删改查的呢？这个当然有，其实也就是Association。
1：使用Association方法，需要把User查询好，然后根据User定义中指定AssociationForeignKey去查找CreditCard
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username"`
	CreditCard []CreditCard `gorm:"constraint:OpUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

	// 创建数据
	//db.Create(&User{
	//	Username: "zhangsan",
	//	CreditCard: []CreditCard{
	//		{Number: "000001"},
	//	},
	//})

	// 查询数据
	u := &User{
		Username: "zhangsan",
	}
	db.First(&u)
	fmt.Println(u)

	// 关联查询
	db.Model(&u).Association("CreditCard").Find(&u.CreditCard)
	fmt.Println(u)
}
```

```shell
# 查看执行结果

PS E:\code\goland\gorm_service> go run .\main.go
&{{1 2023-02-16 01:49:13.887 +0800 CST 2023-02-16 01:49:13.887 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} zhangsan []}
&{{1 2023-02-16 01:49:13.887 +0800 CST 2023-02-16 01:49:13.887 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} zhangsan [{{1 2023-02-16 01:49:13.909 +0800 CST 2023-02-16 01:49:13.909 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} 000001 1}]}


# 那么既然都关联查询了，是不是也可以关联追加数据呢？这个其实也是可以的，直接看代码
```

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username"`
	CreditCard []CreditCard `gorm:"constraint:OpUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

	// 创建数据
	//db.Create(&User{
	//	Username: "zhangsan",
	//	CreditCard: []CreditCard{
	//		{Number: "000001"},
	//	},
	//})

	// 查询数据
	u := &User{
		Username: "zhangsan",
	}
	db.First(&u)
	fmt.Println(u)

	// 关联追加数据
	db.Model(&u).Association("CreditCard").Append([]CreditCard{
		{Number: "000002"},
	})
	fmt.Println(u)
}
```

```shell
# 当然还有什么关联的删除，清空等操作，这个在GORM的官网都可以找到，我就不多敲了。

# 最后我们了解一下Preload（预加载），它的主要作用也是关联查询，但是它与Association的区别
1：Preload只有查询功能，而Association支持CRUD的功能
2：Preload查询时能将User和CreditCard都查询出来，而Association只能查CreditCard数据
```

```go
package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string       `json:"username" gorm:"column:username"`
	CreditCard []CreditCard `gorm:"constraint:OpUpdate:CASCADE,OnDelete:CASCADE;"`
}

type CreditCard struct {
	gorm.Model
	Number string
	UserID uint
}

func main() {
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, CreditCard{})

	// 创建数据
	//db.Create(&User{
	//	Username: "zhangsan",
	//	CreditCard: []CreditCard{
	//		{Number: "000001"},
	//	},
	//})

	// 查询数据
	u := &User{
		Username: "zhangsan",
	}

	// 关联查询
	db.Model(&u).Preload("CreditCard").First(&u)
	data, _ := json.Marshal(u)
	fmt.Println(string(data))
}
```

```shell
# 这里我们用json格式化了一下数据，然后我们来看看数据是什么样的

PS E:\code\goland\gorm_service> go run .\main.go
{"ID":1,"CreatedAt":"2023-02-16T01:49:13.887+08:00","UpdatedAt":"2023-02-16T02:01:38.199+08:00","DeletedAt":null,"username":"zhangsan","CreditCard":[{"ID":1,"CreatedAt":"2023-02-16T01:49:13.909+08:00","UpdatedAt":"2023-02-16T01:49:13.909+08:00","DeletedAt":null,"Number":"000001","UserID":1},{"ID":2,"CreatedAt":"2023-02-16T02:01:38.2+08:00","UpdatedAt":"2023-02-16T02:01:38.2+08:00","DeletedAt":null,"Number":"000002","UserID":1}]}
```

### 7：gorm多对多

```shell
参考文档：https://gorm.io/zh_CN/docs/many_to_many.html

1：many to many 会在两个Model中添加一张链表，也就是会生成第三张表来维护多对多关系。
2：例如，你的应用包含了user和language，且一个user可以对应多个language，多个user也可以对应一个language。
3：当使用GORM的AotuMigrate为User创建表时，GORM会自动创建链表
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Languages []Language `gorm:"many2many:user_languages"`
}

type Language struct {
	gorm.Model
	Name  string
	Users []User `gorm:"many2many:user_languages"`
}

func main() {
	
	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Language{})
}
```

```shell
# 这里是我们常用的一种方法，让两张表都与对方建立多对多关系就OK了，其次就是指定tag来声明多对多，然后指定一下维护这两张表关系的第三张表的名称，然后我们运行看结果

MySQL root@10.0.0.11:gorm_service> show tables;
+------------------------+
| Tables_in_gorm_service |
+------------------------+
| languages              |
| user_languages         |
| users                  |
+------------------------+
3 rows in set
Time: 0.004s

# 可以看到按照我们的预期生成了第三张表了

MySQL root@10.0.0.11:gorm_service> desc users;
+------------+---------------------+------+-----+---------+----------------+
| Field      | Type                | Null | Key | Default | Extra          |
+------------+---------------------+------+-----+---------+----------------+
| id         | bigint(20) unsigned | NO   | PRI | <null>  | auto_increment |
| created_at | datetime(3)         | YES  |     | <null>  |                |
| updated_at | datetime(3)         | YES  |     | <null>  |                |
| deleted_at | datetime(3)         | YES  | MUL | <null>  |                |
+------------+---------------------+------+-----+---------+----------------+

4 rows in set
Time: 0.007s
MySQL root@10.0.0.11:gorm_service> desc languages;
+------------+---------------------+------+-----+---------+----------------+
| Field      | Type                | Null | Key | Default | Extra          |
+------------+---------------------+------+-----+---------+----------------+
| id         | bigint(20) unsigned | NO   | PRI | <null>  | auto_increment |
| created_at | datetime(3)         | YES  |     | <null>  |                |
| updated_at | datetime(3)         | YES  |     | <null>  |                |
| deleted_at | datetime(3)         | YES  | MUL | <null>  |                |
| name       | longtext            | YES  |     | <null>  |                |
+------------+---------------------+------+-----+---------+----------------+

5 rows in set
Time: 0.006s
MySQL root@10.0.0.11:gorm_service> desc user_languages;
+-------------+---------------------+------+-----+---------+-------+
| Field       | Type                | Null | Key | Default | Extra |
+-------------+---------------------+------+-----+---------+-------+
| language_id | bigint(20) unsigned | NO   | PRI | <null>  |       |
| user_id     | bigint(20) unsigned | NO   | PRI | <null>  |       |
+-------------+---------------------+------+-----+---------+-------+

2 rows in set
Time: 0.005s

# 可以看到第三张表内维护了前两张表的主键，当然还有重写外键的方法，不过这个我就不做了，有需要的可以去gorm的官网去看

# 当然这里我们还是要讲一下的，我们是可以自定义第三张表的，也就是说我们可以把任意一张表作为维护关系的表，下面我们来看看
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID   uint
	Name string
}

type UserLanguage struct {
	UserID     int `gorm:"primaryKey"`
	LanguageID int `gorm:"primaryKey"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func main() {

	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Language{}, &UserLanguage{})
}
```

```shell
# 查看运行结果与数据库

Time: 0.002s
MySQL root@10.0.0.11:gorm_service> show tables;
+------------------------+
| Tables_in_gorm_service |
+------------------------+
| languages              |
| user_languages         |
| users                  |
+------------------------+
3 rows in set
Time: 0.005s

MySQL root@10.0.0.11:gorm_service> desc users;
+-------+------------+------+-----+---------+----------------+
| Field | Type       | Null | Key | Default | Extra          |
+-------+------------+------+-----+---------+----------------+
| id    | bigint(20) | NO   | PRI | <null>  | auto_increment |
| name  | longtext   | YES  |     | <null>  |                |
+-------+------------+------+-----+---------+----------------+

2 rows in set
Time: 0.005s
MySQL root@10.0.0.11:gorm_service> desc languages;
+-------+---------------------+------+-----+---------+----------------+
| Field | Type                | Null | Key | Default | Extra          |
+-------+---------------------+------+-----+---------+----------------+
| id    | bigint(20) unsigned | NO   | PRI | <null>  | auto_increment |
| name  | longtext            | YES  |     | <null>  |                |
+-------+---------------------+------+-----+---------+----------------+

2 rows in set
Time: 0.006s
MySQL root@10.0.0.11:gorm_service> desc user_languages;
+-------------+-------------+------+-----+---------+-------+
| Field       | Type        | Null | Key | Default | Extra |
+-------------+-------------+------+-----+---------+-------+
| user_id     | bigint(20)  | NO   | PRI | <null>  |       |
| language_id | bigint(20)  | NO   | PRI | <null>  |       |
| created_at  | datetime(3) | YES  |     | <null>  |       |
| deleted_at  | datetime(3) | YES  |     | <null>  |       |
+-------------+-------------+------+-----+---------+-------+

4 rows in set
Time: 0.007s

# 可以看到和我们的预期是一模一样的，那么后面就是基于多对多的CRUD了
```

```go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID   uint
	Name string
}

type UserLanguage struct {
	UserID     int `gorm:"primaryKey"`
	LanguageID int `gorm:"primaryKey"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func main() {

	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Language{}, &UserLanguage{})

	user := &User{
		Name: "zhangsan",
		Languages: []Language{
			{1, "golang"},
			{2, "python"},
			{3, "java"},
		},
	}
	db.Create(user)
}
```

```shell
# 我们运行完看看结果

MySQL root@10.0.0.11:gorm_service> select * from users;
+----+----------+
| id | name     |
+----+----------+
| 1  | zhangsan |
+----+----------+
1 row in set
Time: 0.005s

MySQL root@10.0.0.11:gorm_service> select * from languages;
+----+--------+
| id | name   |
+----+--------+
| 1  | golang |
| 2  | python |
| 3  | java   |
+----+--------+
3 rows in set
Time: 0.005s

MySQL root@10.0.0.11:gorm_service> select * from user_languages;
+---------+-------------+------------+------------+
| user_id | language_id | created_at | deleted_at |
+---------+-------------+------------+------------+
| 1       | 1           | <null>     | <null>     |
| 1       | 2           | <null>     | <null>     |
| 1       | 3           | <null>     | <null>     |
+---------+-------------+------------+------------+

3 rows in set
Time: 0.008s


# 那么数据是有了，但是我们如何取数据呢，其实就是多对多的查询
```

```go
package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        int
	Name      string
	Languages []Language `gorm:"many2many:user_languages;"`
}

type Language struct {
	ID   uint
	Name string
}

type UserLanguage struct {
	UserID     int `gorm:"primaryKey"`
	LanguageID int `gorm:"primaryKey"`
	CreatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

func main() {

	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_service?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{}, &Language{}, &UserLanguage{})

	// preload
	var users []*User
	db.Preload("Languages").Find(&users)
	data, _ := json.Marshal(users)
	fmt.Println(string(data))
}
```

```shell
# 查看运行结果
PS E:\code\goland\gorm_service> go run .\main.go

# 这里我们格式化一下

[
    {
        "ID": 1,
        "Name": "zhangsan",
        "Languages": [
            {
                "ID": 1,
                "Name": "golang"
            },
            {
                "ID": 2,
                "Name": "python"
            },
            {
                "ID": 3,
                "Name": "java"
            }
        ]
    }
]

# 可以看到擦汗寻出来的数据就是我们创建出来的结果
```

### 8：gorm + gin

```shell
# 前面我们学了Gin框架，上面又学了gorm，那么我们把这两个结合一下形成一个后端API服务来看看如何实现把

# 新建项目创建main.go文件，其次创建分层目录分别是dao，config，controller，service，model，logging，middleware，pkg，router

# 创建数据库，create database gorm_and_gin charset urf8;

# 划分项目目录
1：controller：服务入口，负责处理路由，参数校验，请求转发
2：dao：负责数据与存储相关功能（MySQL，Redis，ES等）
3：middleware：中间件（JWT认证）
4：model：模型定义
5：router：路由
```

#### 8.1：dao/db_init.go (old)

```go
package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
)

func InitMySQL() {

	conn := "root:123456@tcp(10.0.0.11:3306)/gorm_and_gin?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}

	Database = db

	fmt.Println("数据库连接成功")
}
```

#### 8.2：model/user.go

```go
package model

type User struct {
	ID       int        `gorm:"primary_key" json:"id"`
	Username string     `gorm:"not null" json:"username"`
	Password string     `gorm:"not null" json:"password"`
	Token    string     `gorm:"not null" json:"token"`
	Projects []*Project `gorm:"many2many:project_users"`
}

func (*User) TableName() string {
	return "user"
}
```

#### 8.3：model/project

```go
package model

type Project struct {
	ID    int     `gorm:"primary_key" json:"id"`
	Name  string  `gorm:"not null" json:"name"`
	Desc  string  `json:"desc"`
	Users []*User `gorm:"many2many:project_users"`
}

func (*Project) TableName() string {
	return "project"
}
```

#### 8.4：config/mysql.go

```go
package config

var (
	DbUser = "root"
	DbPass = "123456"
	DbHost = "10.0.0.11"
	DbPort = "3306"
	DbName = "gorm_and_gin"
	DbLang = "utf8"
	DbTime = "True"
	DbLoc  = "Local"
)
```

#### 8.5：dao/db_init.go (new)

```go
package dao

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm_gin_service/config"
	"gorm_gin_service/model"
)

func InitMySQL() *gorm.DB {

	// 引用变量并实例化conn
	var conn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName,
		config.DbLang,
		config.DbTime,
		config.DbLoc,
	)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败", err)
	}

	fmt.Println("数据库连接成功")

	if err := db.AutoMigrate(&model.Project{}, &model.User{}); err != nil {
		fmt.Println("表创建失败", err)
	}

	return db
}
```

#### 8.6：controller/user.go

```go
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
)

// RegisterHandler 注册
func RegisterHandler(c *gin.Context) {
	p := new(model.User)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 使用全局的db对象
	if data := dao.InitMySQL().Create(p); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "注册成功",
		"data": p.Username,
	})
}

// LoginHandler 登录
func LoginHandler(c *gin.Context) {
	// 获取参数并与数据库中的数据进行比对绑定
	p := new(model.User)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 定义一个用户对象
	info := &model.User{
		Username: p.Username,
		Password: p.Password,
	}
	// 查询用户是否存在
	if rows := dao.InitMySQL().Where(&info).First(&info); rows == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "用户不存在",
		})
		return
	}

	// 生成token
	token := uuid.New().String()
	// 存入数据库
	if data := dao.InitMySQL().Model(&info).Update("token", token); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	// 登录成功并返回token
	c.JSON(http.StatusOK, gin.H{
		"msg":  "登录成功",
		"data": token,
	})
}
```

#### 8.7：controller/project.go

```go
package controller

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
	"strconv"
)

// CreateProjectHandler 新增项目
func CreateProjectHandler(c *gin.Context) {
	p := new(model.Project)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 数据落库
	if data := dao.InitMySQL().Create(&p); data.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": data.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目创建成功",
		"data": p,
	})
}

// GetProjectHandler 查询项目
func GetProjectHandler(c *gin.Context) {
	projects := make([]model.Project, 0)
	// 列出所有项目
	if err := dao.InitMySQL().Find(&projects); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目列表查询成功",
		"data": projects,
	})
}

// GetProjectDetailHandler 查看项目详情
func GetProjectDetailHandler(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	p := new(model.Project)
	// 查看项目详情
	if err := dao.InitMySQL().Where("id = ?", projectId).First(&p); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg":  "项目信息查询成功",
		"data": p,
	})
}

// UpdateProjectHandler 更新项目
func UpdateProjectHandler(c *gin.Context) {
	p := new(model.Project)
	if err := c.ShouldBindJSON(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 数据落库
	if err := dao.InitMySQL().Save(&p); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg": "项目更新成功",
	})
}

// DeleteProjectHandler 删除项目
func DeleteProjectHandler(c *gin.Context) {
	projectId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 删除关联表的数据
	if err := dao.InitMySQL().Select("Users").Delete(&model.Project{ID: projectId}); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 数据库删除
	if err := dao.InitMySQL().Delete(&model.Project{}, projectId); err.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error.Error(),
		})
		return
	}

	// 返回数据
	c.JSON(http.StatusOK, gin.H{
		"msg": "项目删除成功",
	})
}
```

#### 8.8：middleware/auth.go

```go
package middleware

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/dao"
	"gorm_gin_service/model"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		selectToken := &model.User{}
		row := dao.InitMySQL().Where("token = ?", token).First(selectToken).RowsAffected
		if row == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"msg": "Token is invalid",
			})
			c.Abort()
			return
		} else {
			c.Next()
		}
	}
}
```

#### 8.9：router/test_router.go

```go
package router

import "github.com/gin-gonic/gin"

func TestRouter(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	v1.GET("/ping", TestHandler)
}

func TestHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
```

#### 8.10：router/setup_router.go

```go
package router

import (
	"github.com/gin-gonic/gin"
	"gorm_gin_service/controller"
	"gorm_gin_service/middleware"
)

// SetupApiRouters 用于注册登录的路由
func SetupApiRouters(r *gin.Engine) {
	r.POST("/register", controller.RegisterHandler)
	r.POST("/login", controller.LoginHandler)
	v1 := r.Group("/api/v1")
	r.Use(middleware.AuthMiddleware())
	v1.POST("/project", controller.CreateProjectHandler)
	v1.GET("/project", controller.GetProjectHandler)
	v1.GET("/project/:id", controller.GetProjectDetailHandler)
	v1.PUT("/project", controller.UpdateProjectHandler)
	v1.DELETE("/project/:id", controller.DeleteProjectHandler)
}
```

#### 8.11：router/init_router.go

```go
package router

import "github.com/gin-gonic/gin"

// InitRouter 集合所有的路由
func InitRouter() *gin.Engine {

	r := gin.Default()
	
	TestRouter(r)

	SetupApiRouters(r)

	return r
}
```

#### 8.13：main.go

```go
package main

import (
	"gorm_gin_service/dao"
	"gorm_gin_service/router"
)

func main() {
	dao.InitMySQL()
	r := router.InitRouter()
	r.Run(":80")
}
```

```shell
# 这样我们的代码就完成了，然后我们一步步测试API是否可用
```

### 9：测试API

```apl
GET /ping
```

![image](https://user-images.githubusercontent.com/77761224/219657938-c75a1fc8-0bc0-4bfa-86d5-151b2957e40b.png)

```apl
POST /register

{
	"username": "zhangsan",
	"password": "123456"
}
```

![image](https://user-images.githubusercontent.com/77761224/219659840-bd783411-6882-4183-a2d3-1a7638d64e7c.png)

```apl
POST /login

{
	"username": "zhangsan",
	"password": "123456"
}
```

![image](https://user-images.githubusercontent.com/77761224/219696944-d158ae95-c58c-4ef9-9923-2ea6e7d480eb.png)

```apl
POST /api/v1/project

# 因为涉及到项目的增删改查，所以这个时候中间件就起作用了，所以要在请求的Header内带Token

{
    "name": "微商城",
    "desc": "这是关于电商的一个项目",
    "users": [
        {
            "id": 1
        }
    ]
}
```

![image](https://user-images.githubusercontent.com/77761224/219713770-e06e63b0-a4a9-491f-8818-bc93a31ca4c9.png)

```apl
GET /api/v1/project

# 需要添加Token
```

![image](https://user-images.githubusercontent.com/77761224/219716482-77fa3a7f-5194-4153-89a6-60b345e96a8d.png)

```apl
GET /api/v1/project/<id>

# 需要添加Token
```

![image](https://user-images.githubusercontent.com/77761224/219716656-4c597e74-5cb4-4a0a-bfe5-e876e429f332.png)

```apl
PUT /api/v1/project

# 需要携带Token

{
    "id": 1,
    "name": "微商城",
    "desc": "这是关于电商的一个项目，新!",
    "users": [
        {
            "id": 1
        }
    ]
}
```

![image](https://user-images.githubusercontent.com/77761224/219718556-67fe23be-e61c-4eef-a6eb-ec07841710cc.png)

```apl
DELETE /api/v1/project/<id>

# 需要传递Token
```

![image](https://user-images.githubusercontent.com/77761224/219724534-f41b64e8-52a7-4d10-b495-d44c75e2787a.png)

