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
	MERGE INTO likes l
	USING (
		SELECT CAST($1 AS bigint) AS user_id, CAST($2 AS bigint) AS quote_id, CAST($3 AS smallint) AS val
	) AS n
	ON l.user_id = n.user_id AND l.quote_id = n.quote_id
	WHEN NOT MATCHED THEN
		INSERT VALUES(n.user_id, n.quote_id, n.val)
	WHEN MATCHED AND l.val != n.val THEN
		UPDATE SET val = n.val
	WHEN MATCHED THEN
		DELETE;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{like.UserID, like.QuoteID, like.Val}

	// TODO: upgrade to postgres 17 to support merge returning statements
	// then return the merge_action()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m LikesDatabaseModel) GetLikeDislikeNumForQuote(quoteID int64) (*LikeCount, error) {
	query := `
	SELECT COUNT (CASE WHEN val = 0 THEN 1 ELSE NULL END) AS dislikes,
		   COUNT (CASE WHEN val = 1 THEN 1 ELSE NULL END) AS likes
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
