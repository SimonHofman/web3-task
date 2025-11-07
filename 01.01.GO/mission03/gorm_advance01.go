package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//进阶gorm
//题目1：模型定义
//假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//要求 ：
//使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//编写Go代码，使用Gorm创建这些模型对应的数据库表。
//题目2：关联查询
//基于上述博客系统的模型定义。
//要求 ：
//编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//编写Go代码，使用Gorm查询评论数量最多的文章信息。
//题目3：钩子函数
//继续使用博客系统的模型。
//要求 ：
//为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

type User struct {
	ID        int `gorm:"primaryKey"`
	Name      string
	Posts     []Post `gorm:"foreignKey:UserID""`
	PostCount int
}

type Post struct {
	Id        int `gorm:"primaryKey"`
	Title     string
	Content   string
	UserID    int
	Comments  []Comment `gorm:"foreignKey:PostID"`
	WordCount int
	Status    string
}

func (post *Post) BeforeCreate(db *gorm.DB) (err error) {
	post.WordCount = len([]rune(post.Content))
	post.Status = "正常"
	return
}

type Comment struct {
	Id      int `gorm:"primaryKey"`
	Content string
	PostID  int
}

func (comment *Comment) AfterDelete(db *gorm.DB) (err error) {
	var count int64
	err = db.Model(&Comment{}).Where("post_id = ?", comment.PostID).Count(&count).Error
	if err != nil {
		return
	}
	if count == 0 {
		db.Model(&Post{}).Where("id = ?", comment.PostID).Update("status", "无评论")
	}
	return
}

func initDB(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Post{}, &Comment{})
}

func insertTestData(db *gorm.DB) {
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return
	}

	users := []User{
		{
			Name: "张三",
			Posts: []Post{
				{
					Title:   "文章1",
					Content: "这是文章1的内容",
					Comments: []Comment{
						{Content: "文章1评论1"},
						{Content: "文章1评论2"},
					},
				},
				{
					Title:   "文章2",
					Content: "这是文章2的内容",
					Comments: []Comment{
						{Content: "文章2评论1"},
						{Content: "文章2评论2"},
						{Content: "文章2评论3"},
					},
				},
			},
			PostCount: 2,
		},
		{
			Name: "李四",
			Posts: []Post{
				{
					Title:   "文章3",
					Content: "这是文章3的内容",
					Comments: []Comment{
						{Content: "文章3评论1"},
						{Content: "文章3评论2"},
					},
				},
			},
			PostCount: 1,
		},
		{
			Name: "王五",
			Posts: []Post{
				{
					Title:   "文章4",
					Content: "这是文章4的内容",
				},
			},
			PostCount: 1,
		},
	}

	for _, user := range users {
		db.Create(&user)
	}
}

func QueryUserPostsAndComments(db *gorm.DB, name string) {
	var user User
	err := db.Preload("Posts").Preload("Posts.Comments").Where("name = ?", name).Find(&user).Error
	if err != nil {
		panic("查询失败：" + err.Error())
	}
	//fmt.Printf("用户信息：%v\n", user)
	fmt.Printf("用户：%s\n", user.Name)
	for _, post := range user.Posts {
		fmt.Printf("    标题：%s\n    内容：%s\n", post.Title, post.Content)
		fmt.Printf("    评论：\n")
		for _, comment := range post.Comments {
			fmt.Printf("         %v\n", comment.Content)
		}
	}
}

// 查询评论数量最多的文章信息
func QueryMostCommentsPost(db *gorm.DB) {
	var post Post
	err := db.Model(&Post{}).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON posts.id = comments.post_id").
		Group("posts.id").
		Order("comment_count DESC").
		Limit(1).
		First(&post).Error
	if err != nil {
		panic("查询失败：" + err.Error())
	}
	fmt.Printf("评论数量最多的文章：%s\n", post.Title)
}

func DeleteCommentByID(db *gorm.DB, commentID int) error {
	var comment Comment
	if err := db.Where("id = ?", commentID).First(&comment).Error; err != nil {
		return err
	}
	if err := db.Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

func main() {
	db, err := gorm.Open(sqlite.Open("blog1.db"), &gorm.Config{})
	if err != nil {
		panic("连接数据库失败：" + err.Error())
	}

	// 初始化数据库
	initDB(db)

	// 插入测试数据
	insertTestData(db)

	//// 测试查询
	QueryUserPostsAndComments(db, "张三")
	QueryMostCommentsPost(db)
	DeleteCommentByID(db, 6)
	DeleteCommentByID(db, 7)
}
