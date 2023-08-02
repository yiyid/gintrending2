package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robfig/cron"
)

func pull() (string, string, error) {
	// 爬取网页数据
	resp, err := http.Get("https://github.com/gin-gonic/gin")
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", err
	}

	starText := doc.Find("a[href='/gin-gonic/gin/stargazers']").First().Text()
	re := regexp.MustCompile(`\s+`)
	starRexp := re.ReplaceAllString(starText, " ")
	starSlice := strings.Fields(starRexp) // 使用Fields函数直接切割字符串并去除多余空白字符
	star := starSlice[0]
	createTime := time.Now().Format("2006-01-02 15:04:05")
	if len(starSlice) > 0 {
		return star, createTime, nil
	}
	return "", "", err
}

func insert(db string) {
	fmt.Println("开始录入...")
	star, createTime, err := pull()
	if err != nil {
		fmt.Println("Error fetching star and create time:", err)
		return
	}

	dbClient, err := sql.Open("sqlite3", db)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer dbClient.Close()

	_, err = dbClient.Exec("INSERT INTO stars (star, created_at) VALUES (?, ?)", star, createTime)
	if err != nil {
		fmt.Println("Error executing INSERT query:", err)
		return
	}
	fmt.Printf("Star count: %s, Create time: %s\n", star, createTime)
	fmt.Println("录入完成...")
	fmt.Println("定时任务执行完毕...")
}

func setupGinServer() {
	r := gin.Default()

	// 定义 / 接口
	r.GET("/", func(c *gin.Context) {
		db, err := sql.Open("sqlite3", "stars.db")
		if err != nil {
			fmt.Println("Error opening database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer db.Close()

		rows, err := db.Query("select star, created_at from stars;")
		if err != nil {
			fmt.Println("Error querying database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		var star string
		var createTime string
		var starSlice []map[string]string
		for rows.Next() {
			// 在每次循环中创建新的 starMap
			starMap := make(map[string]string)
			rows.Scan(&star, &createTime)
			starMap["star"] = star
			starMap["createTime"] = createTime
			starSlice = append(starSlice, starMap)
		}

		// 响应JSON数据给浏览器
		c.JSON(http.StatusOK, starSlice)
	})

	// 启动 Gin 服务器
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("Error starting Gin server:", err)
	}
}

func scheduledTask() {
	// 创建一个定时任务调度器
	c := cron.New()

	// 添加定时任务，每天23:59执行一次 insert() 函数
	c.AddFunc("0 18 16 * * *", func() {
		fmt.Println("定时任务开始执行...")
		insert("stars.db")
	})

	// 启动定时任务调度器
	c.Start()

	// 阻塞主线程，保持程序运行
	select {}
}

func main() {
	// 检查数据库文件是否在程序执行目录下
	dbFilename := "stars.db"
	if _, err := os.Stat("stars.db"); os.IsNotExist(err) {
		fmt.Printf("The file %s does not exist in the current directory.\n", dbFilename)
		fmt.Printf("Execute the command:\n")
		fmt.Printf(`
sqlite3 stars.db "CREATE TABLE stars ( 
	id INTEGER PRIMARY KEY,
	star TEXT NOT NULL,
	created_at DATETIME
); 
INSERT INTO stars (star, created_at) VALUES ('0k', datetime('now'));
select * from stars;" -header -column
`)
		fmt.Printf("\nNext execute: ./main\n")
		os.Exit(0)
	}

	// 创建两个进程：一个运行 Gin 服务器，另一个运行定时任务
	go setupGinServer()
	go scheduledTask()

	// 阻塞主线程，保持程序运行
	select {}
}
