
# Quick Start
#### 安装依赖
```
pip3 install -r requirements.txt
```


#### 创建库和表
```
sqlite3 stars.db "CREATE TABLE stars (
    id INTEGER PRIMARY KEY,
    star TEXT NOT NULL,
    created_at DATETIME
); 
INSERT INTO stars (star, created_at) VALUES ('0k', datetime('now'));
select * from stars;" -header -column
```


#### 运行程序
```
python3 main.py &
```


#### 测试
```
curl http://127.0.0.1:8080
```