package data

import (
	"context"
	"database/sql"
	"strconv"
	"time"
)

type LikeType int8

type Like struct {
	UserID  int64    `json:"user_id"`
	QuoteID int64    `json:"quote_id"`
	Val     LikeType `json:"like_type"`
}

type LikeCount struct {
	LikeNum    int64 `json:"likes"`
	DislikeNum int64 `json:"dislikes"`
}

const (
	DislikeValue = iota
	LikeValue
)

func (l LikeType) MarshalJSON() ([]byte, error) {
	switch l {
	case LikeValue:
		return []byte(strconv.Quote("like")), nil
	case DislikeValue:
		return []byte(strconv.Quote("dislike")), nil
	default:
		panic("invalid like type during JSON marshal")
	}
}

type LikesDatabaseModel struct {
	DB *sql.DB
}

func (m LikesDatabaseModel) LikeOrDislikeQuote(like Like) error {
	query := `
	INSERT INTO likes (user_id, quote_id, val)
	VALUES ($1, $2, $3)
	RETURNING val`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{like.UserID, like.QuoteID, like.Val}

	return m.DB.QueryRowContext(ctx, query, args).Scan(&like.Val)
}

func (m LikesDatabaseModel) GetLikeDislikeNumForQuote(quoteID int64) (*LikeCount, error) {
	query := `
	SELECT SUM(val = 0) AS dislike_num,
		   SUM(val = 1) AS like_num
	FROM likes`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var likeCount LikeCount
	err := m.DB.QueryRowContext(ctx, query).Scan(&likeCount.DislikeNum, &likeCount.LikeNum)
	if err != nil {
		return nil, err
	}
	return &likeCount, nil
}
