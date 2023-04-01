package repository

import (
	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreatePosts(posts []*models.Post) (error)
	UpdatePost(post *models.Post) (error)
	SelectPostById(id uint64) (*models.Post, error)
	SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbPost *dataBase) CreatePosts(posts []*models.Post) (error) {
	tx := dbPost.db.Create(&posts)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (dbPost *dataBase) UpdatePost(post *models.Post) (error) {
	tx := dbPost.db.Model(post).Clauses(clause.Returning{}).Updates(models.Post{Message:post.Message, IsEdited: true})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table posts)")
	}

	return nil
}

func (dbPost *dataBase) SelectPostById(id uint64) (*models.Post, error) {
	post := models.Post{}

	tx := dbPost.db.Where("id = ?", id).Take(&post)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table posts)")
	}

	return &post, nil
}

func (dbPost *dataBase) SelectThreadPosts(id uint64, limit int, since int, desc bool, sort string) ([]*models.Post, error) {
	posts := make([]*models.Post, 0, 10)

	if sort == "flat" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND id < ?", id, since).
				Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("id desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND id > ?", id, since).
				Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("id").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "tree" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND post_tree < (?)", id,
				dbPost.db.Table("posts").Select("post_tree").Where("id = ?", since)).
				Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("post_tree desc").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Limit(limit).Where("thread = ? AND post_tree > (?)", id,
				dbPost.db.Table("posts").Select("post_tree").Where("id = ?", since)).
				Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Limit(limit).Where("thread = ?", id).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	} else if sort == "parent_tree" {
		if desc {
			if since != 0 {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id < (?)", id, 
				dbPost.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
				Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
				Order("id desc").Select("id")).Order("post_tree[1] desc, post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		} else {
			if since != 0 {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ? AND id > (?)", id,
				dbPost.db.Table("posts").Select("post_tree[1]").Where("id = ?", since)).
				Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			} else {
				tx := dbPost.db.Where("post_tree[1] IN (?)", dbPost.db.
				Table("posts").Limit(limit).Where("parent = 0 AND thread = ?", id).
				Order("id").Select("id")).Order("post_tree").Find(&posts)
				if tx.Error != nil {
					return nil, errors.Wrap(tx.Error, "database error (table posts)")
				}
			}
		}
	}

	return posts, nil
}

