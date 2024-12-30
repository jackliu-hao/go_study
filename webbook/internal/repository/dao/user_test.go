package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		//mock
		mock func(t *testing.T) *sql.DB

		// except
		ctx      context.Context
		wantUser User
		wantErr  error
	}{
		// 插入成功
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				// mockDb
				mockDb, mock, err := sqlmock.New()
				// 第一个参数是主键， 第二个参数是影响行数
				res := sqlmock.NewResult(3, 1)

				mock.ExpectExec("INSERT INTO `users` . *").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDb
			},
			wantUser: User{},
		},
		// 用户名冲突
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				// mockDb
				mockDb, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` . *").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				require.NoError(t, err)
				return mockDb
			},
			wantUser: User{},
			wantErr:  ErrUserDuplicate,
		},
		// 数据库错误
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				// mockDb
				mockDb, mock, err := sqlmock.New()
				mock.ExpectExec("INSERT INTO `users` . *").
					WillReturnError(errors.New("数据库错误"))
				require.NoError(t, err)
				return mockDb
			},
			wantUser: User{},
			wantErr:  errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			// mockDb
			mockDb := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn: mockDb,
				// 跳过select version的调用
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				// 不需要ping
				DisableAutomaticPing: true,
				// gorm会默认开启事务，这里关闭
				SkipDefaultTransaction: true,
			})

			d := NewGormUserDAO(db)
			err = d.Insert(tc.ctx, tc.wantUser)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
