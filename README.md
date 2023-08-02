
# Quick Start

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


#### 运行 Python 程序
```
cd py-gintrending && pip3 install -r requirements.txt
python3 main.py > output.log 2>&1 &
jobs     # 列出后台任务
fg <n>   # 将任务编号为n的任务切换到前台，可以通过 Ctrl+C 终止脚本运行。
```


#### 运行 Go 程序
```
cd go-gintrending && go mod tidy && go build -o main main.go
./main > output.log 2>&1 &
jobs     # 列出后台任务
fg <n>   # 将任务编号为n的任务切换到前台，可以通过 Ctrl+C 终止脚本运行。
```



#### 验证
```
curl http://127.0.0.1:8080
```