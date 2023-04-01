package repository

import (
	"github.com/kuzkuss/VK_DB_Project/app/models"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryI interface {
	CreateThread(thread *models.Thread) (error)
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error)
	SelectThreadById(id uint64) (*models.Thread, error)
	UpdateThread(thread *models.Thread) (error)
	CreateVote(vote *models.Vote) (error)
}

type dataBase struct {
	db *gorm.DB
}

func New(db *gorm.DB) RepositoryI {
	return &dataBase{
		db: db,
	}
}

func (dbThread *dataBase) CreateThread(thread *models.Thread) (error) {
	if thread.Slug == "" {
		tx := dbThread.db.Omit("votes", "slug").Create(thread)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table threads)")
		}
	} else {
		tx := dbThread.db.Omit("votes").Create(thread)
		if tx.Error != nil {
			return errors.Wrap(tx.Error, "database error (table threads)")
		}
	}

	return nil
}

func (dbThread *dataBase) UpdateThread(thread *models.Thread) (error) {
	tx := dbThread.db.Model(thread).Clauses(clause.Returning{}).Updates(models.Thread{Message:thread.Message, Title: thread.Title})
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table threads)")
	}

	return nil
}

func (dbThread *dataBase) SelectThreadBySlug(slug string) (*models.Thread, error) {
	thread := models.Thread{}

	tx := dbThread.db.Where("slug = ?", slug).Take(&thread)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table threads)")
	}

	return &thread, nil
}

func (dbThread *dataBase) SelectThreadById(id uint64) (*models.Thread, error) {
	thread := models.Thread{}

	tx := dbThread.db.Where("id = ?", id).Take(&thread)
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return nil, models.ErrNotFound
	} else if tx.Error != nil {
		return nil, errors.Wrap(tx.Error, "database error (table threads)")
	}

	return &thread, nil
}

func (dbThread *dataBase) SelectForumThreads(slug string, limit int, since string, desc bool) ([]*models.Thread, error) {
	threads := make([]*models.Thread, 0, 10)

	if desc {
		if since != "" {
			tx := dbThread.db.Limit(limit).Where("forum = ? AND created <= ?", slug, since).
			Order("created desc").Find(&threads)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table threads)")
			}
		} else {
			tx := dbThread.db.Limit(limit).Where("forum = ?", slug).Order("created desc").
			Find(&threads)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table threads)")
			}
		}
	} else {
		if since != "" {
			tx := dbThread.db.Limit(limit).Where("forum = ? AND created >= ?", slug, since).
			Order("created").Find(&threads)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table threads)")
			}
		} else {
			tx := dbThread.db.Limit(limit).Where("forum = ?", slug).Order("created").
			Find(&threads)
			if tx.Error != nil {
				return nil, errors.Wrap(tx.Error, "database error (table threads)")
			}
		}
	}

	return threads, nil
}

func (dbThread *dataBase) CreateVote(vote *models.Vote) (error) {
	tx := dbThread.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "thread_id"}, {Name: "nickname"}},
		DoUpdates: clause.AssignmentColumns([]string{"voice"}),
	  }).Create(vote)
	if tx.Error != nil {
		return errors.Wrap(tx.Error, "database error (table votes)")
	}

	return nil
}

