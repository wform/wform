## Worm 中文文档

### 概述
```go
package worm

```
### 模型
```go
    type Product struct {
        // define primary key
    	Id        int    `db:"primary_key"`
        // mapping customized table column_name
    	FieldName string `db:"column:column_name_test"` 
    }

    type User struct {
        // default column is name
    	Name string
        // default column is age
        Age  int
    }
    func (u *User)BeforeSave(mod interface{})  {
        
    }
    func (u *User)AfterSave(mod interface{})  {
            
    }
    func (u *User)BeforeCreate(mod interface{})  {
        
    }  
    func (u *User)AfterCreate(mod interface{})  {
    	
    }
```
### CURD
#### 查询
```go
    // 获取第一条记录，按主键排序
    worm.E(user).One()
    worm.E(user).One(11)
    // SELECT * FROM users ORDER BY id LIMIT 1;
    
    // 获取所有记录
    worm.E(user).All()
    //// SELECT * FROM users;
    
    // 使用主键获取记录
     worm.E(user).All(10,11,12)
    //// SELECT * FROM users WHERE id = 10;
    worm.E(user).GroupsOne("id")
    worm.E(user).GroupsAll("name")
    
    worm.E().Find(&user) 
    worm.Find(&user)
    worm.E().Find(&users)
    worm.Find(&users)
```
#### 类型断言
```go
    // 获取第一条记录，按主键排序
    var user User
    user = worm.E(user).One().(User)
    fmt.Println(user.Id)
    
    // 获取所有记录
    var userList []User{}
    userList = worm.E(user).All().([]User{})
```

#### Where查询条件(简单条件)
```go
    // 获取第一个匹配记录
    worm.E(user).Where("name = ?", "wnote").One()
    //// SELECT * FROM users WHERE name = 'wnote' limit 1;
    
    // 获取所有匹配记录
     worm.E(user).Where("name = ?", "wnote").All()
    //// SELECT * FROM users WHERE name = 'wnote';
    
     worm.E(user).Where("name <> ?", "wnote").All()
    
    // IN
     worm.E(user).Where("name in (?)", []string{"wnote", "wnote 2"}).All()
    
    // LIKE
     worm.E(user).Where("name LIKE ?", "%worm%").All()
    
    // AND
     worm.E(user).Where("name = ? AND age >= ?", "wnote", "22").All()
    
```
#### Where复杂条件（Struct & Map）
```go
    // Struct
    worm.E(user).Where(&User{Name: "wnote", Age: 20}).One()
    //// SELECT * FROM users WHERE name = "wnote" AND age = 20 LIMIT 1;
    
    // Map
    worm.E(user).Where(map[string]interface{}{"name": "wnote", "age": 20}).All()
    //// SELECT * FROM users WHERE name = "wnote" AND age = 20;
    
    // 主键的Slice
    worm.E(user).Where([]int64{20, 21, 22}).All()
    //// SELECT * FROM users WHERE id IN (20, 21, 22);
```
#### Where Not查询
```go
    worm.E(user).Not("name", "wnote").One()
    //// SELECT * FROM users WHERE name <> "wnote" LIMIT 1;
    
    // Not In
    worm.E(user).Not("name", []string{"wnote", "wnote 2"}).All()
    //// SELECT * FROM users WHERE name NOT IN ("wnote", "wnote 2");
    
    // Not In slice of primary keys
    worm.E(user).Not([]int64{1,2,3}).One()
    //// SELECT * FROM users WHERE id NOT IN (1,2,3);
    
    // Plain SQL
    worm.E(user).Not("name = ?", "wnote").One()
    //// SELECT * FROM users WHERE NOT(name = "wnote");
    
    // Struct
    worm.E(user).Not(User{Name: "wnote"}).One()
    //// SELECT * FROM users WHERE name <> "wnote";
```

#### Where Or条件查询
```go
    worm.E(user).Where("role = ?", "admin").Or("role = ?", "super_admin").All()
    //// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
    
    // Struct
    worm.E(user).Where("name = 'wnote'").Or(User{Name: "wnote 2"}).All()
    //// SELECT * FROM users WHERE name = 'wnote' OR name = 'wnote 2';
    
    // Map
    worm.E(user).Where("name = 'wnote'").Or(map[string]interface{}{"name": "wnote 2"}).All()
```

#### 扩展查询选项
````go
    // 为Select语句添加扩展SQL选项
    worm.E(user).Option("FOR UPDATE").Where(10).Find(&user)
    //// SELECT * FROM users WHERE id = 10 FOR UPDATE;
    
    worm.E().Option("ON CONFLICT").Create(&product)
    // INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
````

### Select 选项
```go
    worm.E(user).Select("name, age").All()
    // SELECT name, age FROM users;
    
    worm.E(user).Select([]string{"name", "age"}).All()
    // SELECT name, age FROM users;
    
    worm.E(user).Select("COALESCE(age,?)", 42).All()
    // SELECT COALESCE(age,'42') FROM users;
```

### Order 选项
```go
    worm.E(user).Order("age desc, name").All()
    // SELECT * FROM users ORDER BY age desc, name;
    
    // Multiple orders
    worm.E(user).Order("age desc").Order("name").All()
    // SELECT * FROM users ORDER BY age desc, name;
    
    // ReOrder
    worm.E(user).Order("age desc").Order("age", true).All()
    // SELECT * FROM users ORDER BY age desc; (users1)
    // SELECT * FROM users ORDER BY age; (users2)
```
#### Limit
```go
    worm.E(user).Limit(3).All()
    // SELECT * FROM users LIMIT 3;
    
    // Cancel limit condition with -1
    worm.E(user).Limit(10).Find(&users1).Limit(-1).All()
    // SELECT * FROM users LIMIT 10; (users1)
    // SELECT * FROM users; (users2)
```

#### Offset 
```go
    worm.E(user).Offset(3).All()
    // SELECT * FROM users OFFSET 3;
    
    // Cancel offset condition with -1
    worm.E(user).Offset(10).Find(&users1).Offset(-1).All()
    // SELECT * FROM users OFFSET 10; (users1)
    // SELECT * FROM users; (users2)
```

#### Count
```go
    worm.E().Where("name = ?", "wnote").Or("name = ?", "wnote 2").Find(&users).Count(&count)
    // SELECT * from USERS WHERE name = 'wnote' OR name = 'wnote 2'; (users)
    // SELECT count(*) FROM users WHERE name = 'wnote' OR name = 'wnote 2'; (count)
    
    worm.E().Where("name = ?", "wnote").Count(&count)
    // SELECT count(*) FROM users WHERE name = 'wnote'; (count)
    
    worm.E().Table("deleted_users").Count(&count)
    // SELECT count(*) FROM deleted_users;
```

#### Group & Having
```go
    type Result struct {
        Date  time.Time
        Total int64
    }
    results := []Result{}
    worm.E(user).Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

#### Join
```go
    rows, err := worm.E(user).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
    for rows.Next() {
    	rows.Scan(&user_name, &email)
        ...
    }
    
    worm.E(user).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
    
    // 多个连接与参数
    worm.E(user).Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "wnote@wnote.net").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

#### Pluck
```go
    var ages []int64
    worm.E(user).Pluck("age", &ages)
    
    var names []string
    worm.E(user).Pluck("name", &names)
    
    worm.E(user).Pluck("name", &names)
    
    // 要返回多个列，做这样：
    worm.E(user).Select("name, age").Find(&users)
```

#### Find
```go
    type Result struct {
        Name string
        Age  int
    }
    
    var result Result
    
    // Raw SQL
    worm.E(user).Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
```

#### Rows 
```go
    worm.E(order).Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
    for rows.Next() {
        ...
    }
```


### 更新
#### 更新全部非空字段
 ```go
    user.Name = "wnote 2"
    user.Age = 100
    worm.E(user).Save()
    // UPDATE users SET name='wnote 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
```

#### 更新更改字段
```go
    // 更新单个属性（如果更改）
    worm.E(user).Update("name", "hello")
    // UPDATE users SET name='hello' WHERE id=111;
    
    // 使用组合条件更新单个属性
    worm.E(user).Where("active = ?", true).Update("name", "hello")
    // UPDATE users SET name='hello' WHERE id=111 AND active=true;
    
    // 使用`map`更新多个属性，只会更新这些更改的字段
    worm.E(user).Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET name='hello', age=18, actived=false WHERE id=111;
    
    // 使用`struct`更新多个属性，只会更新这些更改的和非空白字段
    worm.E(user).Update(User{Name: "hello", Age: 18})
    //// UPDATE users SET name='hello', age=18 WHERE id = 111;
    
    // 警告:当使用struct更新时，FORM将仅更新具有非空值的字段
    // 对于下面的更新，什么都不会更新为""，0，false是其类型的空白值
    worm.E(user).Update(User{Name: "", Age: 0, Actived: false})
```

#### 更新选择的字段
***如果您只想在更新时更新或忽略某些字段，可以使用Select, Omit***
```go
    worm.E(user).Field("name").Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET name='hello' WHERE id=99;
    
    worm.E(user).Omit("name").Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET age=18, actived=false WHERE id=88;
```

####  Batch Updates 批量更新
```go
    worm.E(user).Where("id IN (?)", []int{10, 11}).Update(map[string]interface{}{"name": "hello", "age": 18})
    // UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
    
    // 使用struct更新仅适用于非零值，或使用map[string]interface{}
    worm.E(user).Update(User{Name: "hello", Age: 18})
    // UPDATE users SET name='hello', age=18;
    
    // 使用`RowsAffected`获取更新记录计数
    worm.E(user).Update(User{Name: "hello", Age: 18}).RowsAffected
```

#### 使用SQL表达式更新
```go
    worm.E(product).Update("price", worm.Expr("price * ? + ?", 2, 100))
    // UPDATE products SET price = price * 2 + 100 WHERE id = '2';
    
    worm.E(product).Update(map[string]interface{}{"price": worm.Expr("price * ? + ?", 2, 100)})
    // UPDATE products SET price= price * 2 + 100 WHERE id = '2';
    
    worm.E(product).Update("quantity", worm.Expr("quantity - ?", 1), "price": 1)
    // UPDATE products SET quantity = quantity - 1, price=1 WHERE id = '2';
    
    worm.E(product).Where("quantity > 1").Update("quantity", worm.Expr("quantity - ?", 1))
    // UPDATE products SET quantity = quantity - 1 WHERE id = '2' AND quantity > 1;
```

#### 软删除
```go
    worm.E(user).Delete()
    //UPDATE users SET deleted=1528686522776 WHERE id = 111;

    // 批量删除
    worm.E(user).Where("age = ?", 20).Delete()
    // UPDATE users SET deleted=1528686522776 WHERE age = 20;
    
    unscoped := true
    worm.E(user).Delete(unscoped)
    // DELETE FROM orders WHERE id=10;
```


#### 创建
+ 创建
```go
    worm.Create(&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776);
    // update users set age=0,deleted=0 where id=1;
    worm.Insert(&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776);
    worm.InsertMany([]&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776),(0,0,1528686522776);
```


#### 清空模型的查询条件等信息


### 数据库事务

```go
// 开始事务
ts := worm.E()
ts.Commit()

// 开始事务
trans := worm.Begin()
// 在事务中做一些数据库操作（从这一点使用'tx'，而不是'db'）
trans.Create(...)
// ...
// 发生错误时回滚事务
trans.Rollback()
// 或提交事务
trans.Commit()
```
