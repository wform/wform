## Wform  Documentation
[点我查看中文版](README.cn.md)

### 1, About
 - A Simple orm for go

### 2, Get Start
```
   go get github.com/wform/wform
```

### 3, How to define table model?
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
### 4, How to query?
```go
    // Get first record and order by primary key.
    wform.E(user).One()
    wform.E(user).One(11)
    // SELECT * FROM users ORDER BY id LIMIT 1;
    
    // Get all records from table.
    wform.E(user).All()
    //// SELECT * FROM users;
    
    // Get record by id.
     wform.E(user).All(10,11,12)
    //// SELECT * FROM users WHERE id = 10;
    wform.E(user).GroupsOne("id")
    wform.E(user).GroupsAll("name")
    
    wform.E().Find(&user) 
    wform.Find(&user)
    wform.E().Find(&users)
    wform.Find(&users)
```
### 5, Type Asserts
```go
    // Get first record and order by primary key.
    var user User
    user = wform.E(user).One().(User)
    fmt.Println(user.Id)
    
    // Get all records from table.
    var userList []User{}
    userList = wform.E(user).All().([]User{})
```

### 6, How to query records from a condition.
```go
    // Get first record from query condition.
    wform.E(user).Where("name = ?", "wnote").One()
    //// SELECT * FROM users WHERE name = 'wnote' limit 1;
    
    // Get all records from condition.
     wform.E(user).Where("name = ?", "wnote").All()
    //// SELECT * FROM users WHERE name = 'wnote';
    
     wform.E(user).Where("name <> ?", "wnote").All()
    
    // IN
     wform.E(user).Where("name in (?)", []string{"wnote", "wnote 2"}).All()
    
    // LIKE
     wform.E(user).Where("name LIKE ?", "%wform%").All()
    
    // AND
     wform.E(user).Where("name = ? AND age >= ?", "wnote", "22").All()
    
```
### 7, How to query records from mixed conditions.（Struct & Map）
```go
    // Struct
    wform.E(user).Where(&User{Name: "wnote", Age: 20}).One()
    //// SELECT * FROM users WHERE name = "wnote" AND age = 20 LIMIT 1;
    
    // Map
    wform.E(user).Where(map[string]interface{}{"name": "wnote", "age": 20}).All()
    //// SELECT * FROM users WHERE name = "wnote" AND age = 20;
    
    // slice combined with primary keys 
    wform.E(user).Where([]int64{20, 21, 22}).All()
    //// SELECT * FROM users WHERE id IN (20, 21, 22);
```
### 8, Not Condition A
```go
    wform.E(user).Not("name", "wnote").One()
    //// SELECT * FROM users WHERE name <> "wnote" LIMIT 1;
    
    // Not In
    wform.E(user).Not("name", []string{"wnote", "wnote 2"}).All()
    //// SELECT * FROM users WHERE name NOT IN ("wnote", "wnote 2");
    
    // Not In slice of primary keys
    wform.E(user).Not([]int64{1,2,3}).One()
    //// SELECT * FROM users WHERE id NOT IN (1,2,3);
    
    // Plain SQL
    wform.E(user).Not("name = ?", "wnote").One()
    //// SELECT * FROM users WHERE NOT(name = "wnote");
    
    // Struct
    wform.E(user).Not(User{Name: "wnote"}).One()
    //// SELECT * FROM users WHERE name <> "wnote";
```

### 9, Condition A or Condition B
```go
    wform.E(user).Where("role = ?", "admin").Or("role = ?", "super_admin").All()
    //// SELECT * FROM users WHERE role = 'admin' OR role = 'super_admin';
    
    // Struct
    wform.E(user).Where("name = 'wnote'").Or(User{Name: "wnote 2"}).All()
    //// SELECT * FROM users WHERE name = 'wnote' OR name = 'wnote 2';
    
    // Map
    wform.E(user).Where("name = 'wnote'").Or(map[string]interface{}{"name": "wnote 2"}).All()
```

### 10, Extend query options.
````go
    // Add options for query.
    wform.E(user).Option("FOR UPDATE").Where(10).Find(&user)
    //// SELECT * FROM users WHERE id = 10 FOR UPDATE;
    
    wform.E().Option("ON CONFLICT").Create(&product)
    // INSERT INTO products (name, code) VALUES ("name", "code") ON CONFLICT;
````

### 11, How to query with distinct fields. 
```go
    wform.E(user).Select("name, age").All()
    // SELECT name, age FROM users;
    
    wform.E(user).Select([]string{"name", "age"}).All()
    // SELECT name, age FROM users;
    
    wform.E(user).Select("COALESCE(age,?)", 42).All()
    // SELECT COALESCE(age,'42') FROM users;
```

### 12, How to query with order by. 
```go
    wform.E(user).Order("age desc, name").All()
    // SELECT * FROM users ORDER BY age desc, name;
    
    // Multiple orders
    wform.E(user).Order("age desc").Order("name").All()
    // SELECT * FROM users ORDER BY age desc, name;
    
    // ReOrder
    wform.E(user).Order("age desc").Order("age", true).All()
    // SELECT * FROM users ORDER BY age desc; (users1)
    // SELECT * FROM users ORDER BY age; (users2)
```
### 13, Limit
```go
    wform.E(user).Limit(3).All()
    // SELECT * FROM users LIMIT 3;
    
    // Cancel limit condition with -1
    wform.E(user).Limit(10).Find(&users1).Limit(-1).All()
    // SELECT * FROM users LIMIT 10; (users1)
    // SELECT * FROM users; (users2)
```

### 14, Offset 
```go
    wform.E(user).Offset(3).All()
    // SELECT * FROM users OFFSET 3;
    
    // Cancel offset condition with -1
    wform.E(user).Offset(10).Find(&users1).Offset(-1).All()
    // SELECT * FROM users OFFSET 10; (users1)
    // SELECT * FROM users; (users2)
```

### 15, Count
```go
    wform.E().Where("name = ?", "wnote").Or("name = ?", "wnote 2").Find(&users).Count(&count)
    // SELECT * from USERS WHERE name = 'wnote' OR name = 'wnote 2'; (users)
    // SELECT count(*) FROM users WHERE name = 'wnote' OR name = 'wnote 2'; (count)
    
    wform.E().Where("name = ?", "wnote").Count(&count)
    // SELECT count(*) FROM users WHERE name = 'wnote'; (count)
    
    wform.E().Table("deleted_users").Count(&count)
    // SELECT count(*) FROM deleted_users;
```

### 16, Group & Having
```go
    type Result struct {
        Date  time.Time
        Total int64
    }
    results := []Result{}
    wform.E(user).Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Scan(&results)
```

### 17, Join
```go
    rows, err := wform.E(user).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Rows()
    for rows.Next() {
    	rows.Scan(&user_name, &email)
        ...
    }
    
    wform.E(user).Select("users.name, emails.email").Joins("left join emails on emails.user_id = users.id").Scan(&results)
    
    // Multiple join queries.
    wform.E(user).Joins("JOIN emails ON emails.user_id = users.id AND emails.email = ?", "wnote@wnote.net").Joins("JOIN credit_cards ON credit_cards.user_id = users.id").Where("credit_cards.number = ?", "411111111111").Find(&user)
```

### 18, Pluck
```go
    var ages []int64
    wform.E(user).Pluck("age", &ages)
    
    var names []string
    wform.E(user).Pluck("name", &names)
    
    wform.E(user).Pluck("name", &names)
    
    // How to select multiple users.
    wform.E(user).Select("name, age").Find(&users)
```

### 19, Find
```go
    type Result struct {
        Name string
        Age  int
    }
    
    var result Result
    
    // Raw SQL
    wform.E(user).Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)
```

### 20, Rows 
```go
    wform.E(order).Select("date(created_at) as date, sum(amount) as total").Group("date(created_at)").Having("sum(amount) > ?", 100).Rows()
    for rows.Next() {
        ...
    }
```


### 21, How to update all the none empty fields.
 ```go
    user.Name = "wnote 2"
    user.Age = 100
    wform.E(user).Save()
    // UPDATE users SET name='wnote 2', age=100, birthday='2016-01-01', updated_at = '2013-11-17 21:34:10' WHERE id=111;
```

### 22, How to update limit fields.
```go
    // Update only one filed.
    wform.E(user).Update("name", "hello")
    // UPDATE users SET name='hello' WHERE id=111;
    
    // Update only one filed with mixed query conditions.
    wform.E(user).Where("active = ?", true).Update("name", "hello")
    // UPDATE users SET name='hello' WHERE id=111 AND active=true;
    
    // How to update more than one fields with map.
    wform.E(user).Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET name='hello', age=18, actived=false WHERE id=111;
    
    // How to update more than one fields with struct, ant it will update only the fields that is not empty.
    wform.E(user).Update(User{Name: "hello", Age: 18})
    //// UPDATE users SET name='hello', age=18 WHERE id = 111;
    
    // Warn:It will update only the fields that is not empty.
    // And next code will not execute the update query because all the fields in the struct 'User' is empty.
    wform.E(user).Update(User{Name: "", Age: 0, Actived: false})
```

### 23, How to update only the selected fields.
***If you want to update or ignore some fields from the update struct & map,you can use Select or use Omit***
```go
    wform.E(user).Field("name").Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET name='hello' WHERE id=99;
    
    wform.E(user).Omit("name").Update(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
    // UPDATE users SET age=18, actived=false WHERE id=88;
```

### 24,  How to update multiple records.
```go
    wform.E(user).Where("id IN (?)", []int{10, 11}).Update(map[string]interface{}{"name": "hello", "age": 18})
    // UPDATE users SET name='hello', age=18 WHERE id IN (10, 11);
    
    // Use `RowsAffected` get the affected rows from the last update.
    wform.E(user).Update(User{Name: "hello", Age: 18}).RowsAffected
```

### 25, How to update by sql expression.
```go
    wform.E(product).Update("price", wform.Expr("price * ? + ?", 2, 100))
    // UPDATE products SET price = price * 2 + 100 WHERE id = '2';
    
    wform.E(product).Update(map[string]interface{}{"price": wform.Expr("price * ? + ?", 2, 100)})
    // UPDATE products SET price= price * 2 + 100 WHERE id = '2';
    
    wform.E(product).Update("quantity", wform.Expr("quantity - ?", 1), "price": 1)
    // UPDATE products SET quantity = quantity - 1, price=1 WHERE id = '2';
    
    wform.E(product).Where("quantity > 1").Update("quantity", wform.Expr("quantity - ?", 1))
    // UPDATE products SET quantity = quantity - 1 WHERE id = '2' AND quantity > 1;
```

### 26, Soft delete.
```go
    wform.E(user).Delete()
    //UPDATE users SET deleted=1528686522776 WHERE id = 111;

    // Batch delete
    wform.E(user).Where("age = ?", 20).Delete()
    // UPDATE users SET deleted=1528686522776 WHERE age = 20;
    
    unscoped := true
    wform.E(user).Delete(unscoped)
    // DELETE FROM orders WHERE id=10;
```


### 27, How to create a record for a table.
+ Create
```go
    wform.Create(&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776);
    // update users set age=0,deleted=0 where id=1;
    wform.Insert(&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776);
    wform.InsertMany([]&User{})
    // insert into users(age,deleted,created)values(0,0,1528686522776),(0,0,1528686522776);
```


### 28, Clear conditions from last query.


### 29, Db transaction.

```go
// Begin a transaction.
ts := wform.E()
ts.Commit()

// Begin a transaction.
trans := wform.Begin()
// If you want to execute a query from a transaction, please use trans and not the code wform.E()
trans.Create(...)
// ...
// If error occurred, use Rollback to rollback the transaction.
trans.Rollback()
// How to commit a transaction.
trans.Commit()
```
