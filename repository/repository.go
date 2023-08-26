package repository

import (
	"github.com/shimodatkh/contestboard-public/log"
	"github.com/shimodatkh/contestboard-public/model"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteRepository struct {
	db *gorm.DB
}

func NewSqliteRepository(useInmemoryDB bool) *SqliteRepository {
	log := log.InitLogger(zapcore.DebugLevel)

	var dbFilename string
	if useInmemoryDB {
		dbFilename = ":memory:"
	} else {
		dbFilename = "gorm.db"
	}
	db, err := gorm.Open(sqlite.Open(dbFilename), &gorm.Config{})
	if err != nil {
		log.Errorw("db connection error", "err", err.Error())
	} else {
		log.Infow("db connected", "dbFilename", dbFilename)
	}
	db.AutoMigrate(&model.Measurement{})
	db.AutoMigrate(&model.TimeResult{})
	return &SqliteRepository{
		db: db,
	}
}

func (s *SqliteRepository) InsertMeasurement(measurement *model.Measurement) {
	log := log.InitLogger(zapcore.DebugLevel)

	result := s.db.Create(measurement)
	if result.Error != nil {
		log.Errorw("InsertMeasurement error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
	}
}

func (s *SqliteRepository) GetNewMeasurementID() int64 {
	log := log.InitLogger(zapcore.DebugLevel)
	var measurement *model.Measurement
	result := s.db.Order("measurementid desc").First(&measurement)
	if result.Error != nil {
		log.Errorw("GetNewMeasurementID error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
		return 0
	}
	return measurement.Measurementid + 1
}

func (s *SqliteRepository) GetAllMeasurements() []*model.Measurement {
	log := log.InitLogger(zapcore.DebugLevel)

	var measurements []*model.Measurement
	result := s.db.Preload("TimeResults").Order("measurementid desc").Find(&measurements)
	if result.Error != nil {
		log.Errorw("GetAllMeasurements error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
		return nil
	}
	return measurements
}

func (s *SqliteRepository) DeleteMeasurement(measurementid int64) {
	log := log.InitLogger(zapcore.DebugLevel)

	result := s.db.Where("measurementid = ?", measurementid).Delete(&model.Measurement{})
	if result.Error != nil {
		log.Errorw("DeleteMeasurement error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
	}
}

func (s *SqliteRepository) UpdateMeasurementByMemo(measurementid int64, memo string) {
	log := log.InitLogger(zapcore.DebugLevel)

	result := s.db.Model(&model.Measurement{}).Where("measurementid = ?", measurementid).Update("memo", memo)
	if result.Error != nil {
		log.Errorw("UpdateMeasurementByMemo error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
	}
}

func (s *SqliteRepository) UpdateMeasurementByParam(measurementid int64, param string) {
	log := log.InitLogger(zapcore.DebugLevel)

	result := s.db.Model(&model.Measurement{}).Where("measurementid = ?", measurementid).Update("param", param)
	if result.Error != nil {
		log.Errorw("UpdateMeasurementByParam error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
	}
}

func (s *SqliteRepository) UpdateMeasurementByScore(measurementid int64, score float64) {
	log := log.InitLogger(zapcore.DebugLevel)

	result := s.db.Model(&model.Measurement{}).Where("measurementid = ?", measurementid).Update("score", score)
	if result.Error != nil {
		log.Errorw("UpdateMeasurementByScore error", "err", result.Error.Error(), "RowsAffected", result.RowsAffected)
	}
}
