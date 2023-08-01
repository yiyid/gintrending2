
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
python3 main.py > output.log 2>&1 &
jobs  # 列出后台任务
fg <n>   # 将任务编号为n的任务切换到前台，可以通过 Ctrl+C 终止脚本运行。
```


#### 测试
```
curl http://127.0.0.1:8080
```