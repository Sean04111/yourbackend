// Code generated by goctl. DO NOT EDIT.

package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	userFieldNames          = builder.RawFieldNames(&User{})
	userRows                = strings.Join(userFieldNames, ",")
	userRowsExpectAutoSet   = strings.Join(stringx.Remove(userFieldNames, "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), ",")
	userRowsWithPlaceHolder = strings.Join(stringx.Remove(userFieldNames, "`email`", "`create_at`", "`create_time`", "`created_at`", "`update_at`", "`update_time`", "`updated_at`"), "=?,") + "=?"
)

type (
	userModel interface {
		Insert(ctx context.Context, data *User) (sql.Result, error)
		FindOne(ctx context.Context, email string) (*User, error)
		Update(ctx context.Context, data *User) error
		Delete(ctx context.Context, email string) error
	}

	defaultUserModel struct {
		conn  sqlx.SqlConn
		table string
	}

	User struct {
		Uid        string  `db:"uid"`
		Email      string `db:"email"`
		Password   string `db:"password"`
		Name       string `db:"name"`
		AvatarLink string `db:"avatarLink"`
		Profession string `db:"profession"`
		Type       string `db:"type"`
	}
)

func newUserModel(conn sqlx.SqlConn) *defaultUserModel {
	return &defaultUserModel{
		conn:  conn,
		table: "`user`",
	}
}

func (m *defaultUserModel) Delete(ctx context.Context, email string) error {
	query := fmt.Sprintf("delete from %s where `email` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, email)
	return err
}

func (m *defaultUserModel) FindOne(ctx context.Context, email string) (*User, error) {
	query := fmt.Sprintf("select %s from %s where `email` = ? limit 1", userRows, m.table)
	var resp User
	err := m.conn.QueryRowCtx(ctx, &resp, query, email)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultUserModel) Insert(ctx context.Context, data *User) (sql.Result, error) {
	query := fmt.Sprintf("insert into %s (%s) values (?, ?, ?, ?, ?, ?, ?)", m.table, userRowsExpectAutoSet)
	ret, err := m.conn.ExecCtx(ctx, query, data.Uid, data.Email, data.Password, data.Name, data.AvatarLink, data.Profession, data.Type)
	return ret, err
}

func (m *defaultUserModel) Update(ctx context.Context, data *User) error {
	query := fmt.Sprintf("update %s set %s where `email` = ?", m.table, userRowsWithPlaceHolder)
	_, err := m.conn.ExecCtx(ctx, query, data.Uid, data.Password, data.Name, data.AvatarLink, data.Profession, data.Type, data.Email)
	return err
}

func (m *defaultUserModel) tableName() string {
	return m.table
}
