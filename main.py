import os
import sys
import time
import sqlite3
import multiprocessing

import schedule
import requests
from bs4 import BeautifulSoup
from flask import Flask


def pull():
    try:
        r = requests.get("https://github.com/gin-gonic/gin", timeout=(5, 5))
        r.raise_for_status()
        soup = BeautifulSoup(r.text, "html.parser")
        a = soup.find("a", attrs={"href": "/gin-gonic/gin/stargazers"})
        star = a.text.strip().split()[0]
        create_time = time.strftime("%Y-%m-%d %H:%M:%S", time.localtime())
        return star, create_time
    except Exception as e:
        print(e)


def insert():
    try:
        star, create_time = pull()
        conn = sqlite3.connect("stars.db")
        cursor = conn.cursor()
        cursor.execute(
            "insert into stars (star, created_at) values (?, ?)", (star, create_time)
        )
        conn.commit()
        cursor.close()
        conn.close()
        print(f"star: {star}, create_at: {create_time}, 录入完成...")
    except Exception as e:
        print(e)


# 创建 Flask 应用
app = Flask(__name__)


@app.route("/", methods=["GET"])
def get_today_star_count():
    # 连接到数据库，查询 star 数量
    conn = sqlite3.connect("stars.db")
    cursor = conn.cursor()
    cursor.execute(
        "SELECT star, created_at FROM stars;",
    )
    raw_list = []
    for raw in cursor.fetchall():
        raw_dict = {"star_count": raw[0], "create_at": raw[1]}
        raw_list.append(raw_dict)
    cursor.close()
    conn.close()
    return raw_list


def flask_app():
    # 启动 Flask 服务器
    app.run(host="0.0.0.0", port=8080)


def scheduled_task():
    # 每天 00:00 执行一次 insert() 函数
    schedule.every().day.at("15:43").do(insert)

    # 循环执行定时任务
    while True:
        schedule.run_pending()
        time.sleep(1)


def check_file_exists(filename):
    if not os.path.exists(filename):
        print(f"The file '{filename}' does not exist in the current directory.")
        print("Execute the command:")
        print(
            """
sqlite3 stars.db "CREATE TABLE stars (
    id INTEGER PRIMARY KEY,
    star TEXT NOT NULL,
    created_at DATETIME
); 
INSERT INTO stars (star, created_at) VALUES ('0k', datetime('now'));
select * from stars;" -header -column
"""
        )
        print("Next execute: python3 main.py")
        sys.exit()


if __name__ == "__main__":
    # 检查数据库文件是否在程序执行目录下
    filename = "stars.db"
    check_file_exists(filename)

    # 创建两个进程：一个运行 Flask 应用，另一个运行定时任务
    flask_process = multiprocessing.Process(target=flask_app)
    scheduled_process = multiprocessing.Process(target=scheduled_task)

    # 启动两个进程
    flask_process.start()
    scheduled_process.start()

    # 等待两个进程结束
    flask_process.join()
    scheduled_process.join()
