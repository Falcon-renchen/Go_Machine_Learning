package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)
//连接数据库用于存储，查询，提取数据。
func main() {
	pgURL := os.Getenv("PGURL")
	if pgURL == "" {
		log.Fatal("PGURL empty")
	}
	db, err := sql.Open("postgres",pgURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/*if err := db.Ping(); err !=nil {
		log.Fatal(err)
	}*/

	// Query the database.
	rows, err := db.Query(`
    SELECT
       sepal_length as sLength,
       sepal_width as sWidth,
       petal_length as pLength,
       petal_width as pWidth
    FROM iris
    WHERE species = $1`, "Iris-setosa")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()


	//遍历行，将结果发送到标准输出。
	for rows.Next() {
		var (
			sLength float64
			sWidth float64
			pLength float64
			pWidth float64
		)

		if err := rows.Scan(&sLength, &sWidth, &pLength, &pWidth); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%.2f, %.2f, %.2f, %.2f\n", sLength, sWidth, pLength, pWidth)
	}

	//在完成对行的迭代之后检查错误。
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	//更新一些值
	res, err := db.Exec("UPDATE iris SET species = 'setosa' WHERE species = 'Iris-setosa'")
	if err != nil {
		log.Fatal(err)
	}

	//查看更新了多少行。
	rowCount, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	//将行数输出到标准输出。
	log.Printf("affected = %d\n", rowCount)
}