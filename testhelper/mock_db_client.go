package testhelper

import (
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetNewDbMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	// sqlmock.New()で新しいmockデータベースを作成
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, mock, err
	}

	// SELECT VERSION() に対する期待値を追加
	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version()"}).AddRow("8.0.21")) // 例としてMySQLのバージョンを返す

	// Gormに適切なドライバを指定して、mockデータベースを開く
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db, // sqlmockから取得した*sql.DBを使用
	}), &gorm.Config{})
	if err != nil {
		return nil, mock, err
	}

	return gormDB, mock, nil
}
