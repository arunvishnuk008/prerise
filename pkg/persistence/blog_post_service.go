package persistence

import (
	"crimson-sunrise.site/pkg/common"
	"crimson-sunrise.site/pkg/model"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"
)

func GetAllBlogPosts(db *sql.DB) ([]model.BlogPost, error) {
	const query = "SELECT bp.id, bp.title, bp.tag_line, bp.tags, bp.content, " +
		"bp.created_at, bp.updated_at from blog_posts bp ORDER BY bp.created_at DESC "
	result, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(resultSet *sql.Rows) {
		err := resultSet.Close()
		if err != nil {
			log.Printf("error closing resultset %s", err.Error())
		}
	}(result)
	var posts []model.BlogPost
	for result.Next() {
		var post model.BlogPost
		err = result.Scan(&post.ID, &post.Title, &post.TagLine, &post.Tags, &post.Content, &post.CreatedAtUnix, &post.UpdatedAtUnix)
		if err != nil {
			return nil, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt, err := common.EpochSecondToTime(strconv.Itoa(post.CreatedAtUnix))
		if err != nil {
			return nil, err
		}
		createdAt = createdAt.In(loc)
		post.CreatedAt = createdAt.Format(time.RFC3339)
		updatedAt, err := common.EpochSecondToTime(strconv.Itoa(post.UpdatedAtUnix))
		if err != nil {
			return nil, err
		}
		updatedAt = updatedAt.In(loc)
		post.UpdatedAt = updatedAt.Format(time.RFC3339)
		posts = append(posts, post)
	}
	return posts, nil
}

func GetBlogPostByID(db *sql.DB, id int64) (*model.BlogPost, error) {
	const query =
		"SELECT bp.id, bp.title, bp.tag_line, bp.tags, bp.content, " +
		"bp.created_at, bp.updated_at from blog_posts bp where bp.id =?"
	result, err := db.Query(query,id)
	if err != nil {
		return nil, err
	}
	defer func(resultSet *sql.Rows) {
		err := resultSet.Close()
		if err != nil {
			log.Printf("error closing resultset %s", err.Error())
		}
	}(result)
	if result.Next() {
		var post model.BlogPost
		err = result.Scan(&post.ID, &post.Title, &post.TagLine, &post.Tags, &post.Content, &post.CreatedAtUnix, &post.UpdatedAtUnix)
		if err != nil {
			return nil, err
		}
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt, err := common.EpochSecondToTime(strconv.Itoa(post.CreatedAtUnix))
		if err != nil {
			return nil, err
		}
		createdAt = createdAt.In(loc)
		post.CreatedAt = createdAt.Format(time.RFC3339)
		updatedAt, err := common.EpochSecondToTime(strconv.Itoa(post.UpdatedAtUnix))
		if err != nil {
			return nil, err
		}
		updatedAt = updatedAt.In(loc)
		post.UpdatedAt = updatedAt.Format(time.RFC3339)
		return &post, nil
	}
	return nil, errors.New("blog post not found")
}

func AddNewPost(db *sql.DB, post model.BlogPost) (*model.BlogPost, error) {
	const query =
		"INSERT INTO blog_posts(title,tag_line,tags,content,created_at,updated_at)" +
			" VALUES(?,?,?,?,?,?)"
	result, err := db.Exec(query, post.Title, post.TagLine, post.Tags, post.Content, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return nil, err
	}
	lastinsertId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	if lastinsertId == 0 {
		return nil, errors.New("faield to insert blog post to db")
	}
	return GetBlogPostByID(db, lastinsertId)
}

func UpdateBlogPost(db *sql.DB, post model.BlogPost) (*model.BlogPost, error) {
	const query =
		"UPDATE blog_posts SET title=?, tag_line=?, tags=?, content=?, created_at=?, updated_at=? where id=?"
	result, err := db.Exec(query, post.Title, post.TagLine, post.Tags, post.Content, post.CreatedAtUnix, time.Now().Unix(), post.ID)
	if err != nil {
		return nil, err
	}
	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsUpdated == 0 {
		return nil, errors.New("failed to update blog post")
	}
	return GetBlogPostByID(db, post.ID)
}

func DeleteBlogPostByID(db *sql.DB, id int64) (bool, error) {
	const query =
		"DELETE FROM blog_posts where id=?"
	result, err := db.Exec(query,id)
	if err != nil {
		return false, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsAffected > 0, nil
}