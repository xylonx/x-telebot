package session

import (
	"time"

	"github.com/xylonx/x-telebot/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SessionModel struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ID     string `json:"session_id" gorm:"column:session_id;primaryKey"`
	Value  string `json:"session_value" gorm:"column:session_value"`
	Expire int64  `json:"session_expire" gorm:"column:session_expire"`
}

type SqliteSession struct {
	db        *gorm.DB
	SessionID string
	encoder   Encoder
	decoder   Decoder
}

var _ Session = &SqliteSession{}

func (s *SessionModel) TableName() string {
	return "session"
}

func (s *SqliteSession) Set(value interface{}, opt Option) (err error) {
	if s.encoder == nil {
		return ErrorNilEncoder
	}
	if s.SessionID == "" {
		return ErrorNilSessionID
	}

	model := &SessionModel{
		ID:     s.SessionID,
		Expire: int64(opt.Expire),
	}

	model.Value, err = s.encoder.Encode(value)
	if err != nil {
		return err
	}

	return s.db.Create(model).Error
}

func (s *SqliteSession) Get() (interface{}, error) {
	if s.decoder == nil {
		return nil, ErrorNilDecoder
	}
	if s.SessionID == "" {
		return nil, ErrorNilSessionID
	}

	model := &SessionModel{
		ID: s.SessionID,
	}

	tmp := s.db.Where(model).First(model)
	if err := tmp.Error; err != nil {
		return nil, err
	}
	if tmp.RowsAffected == 0 {
		return nil, ErrorNilSession
	}

	return s.decoder.Decode(model.Value)
}

func (s *SqliteSession) Close() error {
	if s.SessionID == "" {
		return ErrorNilSessionID
	}

	return s.db.Delete(&SessionModel{ID: s.SessionID}).Error
}

type SQLiteSessionManager struct {
	db *gorm.DB
}

var _ SessionManager = &SQLiteSessionManager{}

func NewSQLLiteSessionManager(dsn string) (SessionManager, error) {
	// if dsn(*.db) does not exist, create it.
	db, err := database.NewSqliteConn(dsn)
	if err != nil {
		return nil, err
	}
	return &SQLiteSessionManager{
		db: db,
	}, nil
}

func (s *SQLiteSessionManager) New(sessionID string, opts ...OptionFunc) (Session, error) {
	opt := defaultOption()
	for i := range opts {
		opts[i](opt)
	}

	sess := &SqliteSession{
		db: s.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}),
		SessionID: sessionID,
		encoder:   opt.Encoder,
		decoder:   opt.Decoder,
	}

	if err := sess.Set("", *opt); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *SQLiteSessionManager) Load(sessionID string, opts ...OptionFunc) (Session, error) {
	opt := defaultOption()
	for i := range opts {
		opts[i](opt)
	}

	sess := &SqliteSession{
		db:        s.db,
		SessionID: sessionID,
		encoder:   opt.Encoder,
		decoder:   opt.Decoder,
	}

	// judge whether session exists
	if _, err := sess.Get(); err != nil {
		return nil, err
	}

	return sess, nil
}
