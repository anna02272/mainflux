// Copyright (c) Mainflux
// SPDX-License-Identifier: Apache-2.0

package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MainfluxLabs/mainflux/pkg/dbutil"
	"github.com/MainfluxLabs/mainflux/pkg/errors"
	"github.com/MainfluxLabs/mainflux/pkg/transformers/senml"
	"github.com/MainfluxLabs/mainflux/readers"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
)

const (
	defTable  = "messages"
	jsonTable = "json"
)

var _ readers.MessageRepository = (*postgresRepository)(nil)

var (
	errInvalidMessage = errors.New("invalid message representation")
	errTransRollback  = errors.New("failed to rollback transaction")
)

type postgresRepository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) readers.MessageRepository {
	return &postgresRepository{
		db: db,
	}
}

func (tr postgresRepository) ListAllMessages(rpm readers.PageMetadata) (readers.MessagesPage, error) {
	return tr.readAll(rpm)
}

func (tr postgresRepository) Backup(rpm readers.PageMetadata) (readers.MessagesPage, error) {
	return tr.readAll(rpm)
}

func (tr postgresRepository) DeleteMessages(ctx context.Context, rpm readers.PageMetadata, table string) error {
	tx, err := tr.db.BeginTxx(ctx, nil)
	if err != nil {
		return errors.Wrap(errors.ErrSaveMessages, err)
	}

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = errors.Wrap(err, errors.Wrap(errTransRollback, txErr))
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = errors.Wrap(errors.ErrDeleteMessages, err)
		}
	}()

	condition := fmtCondition(rpm, table)
	q := fmt.Sprintf("DELETE FROM %s %s", table, condition)

	params := map[string]interface{}{
		"subtopic":     rpm.Subtopic,
		"publisher":    rpm.Publisher,
		"name":         rpm.Name,
		"protocol":     rpm.Protocol,
		"value":        rpm.Value,
		"bool_value":   rpm.BoolValue,
		"string_value": rpm.StringValue,
		"data_value":   rpm.DataValue,
		"from":         rpm.From,
		"to":           rpm.To,
	}

	_, err = tx.NamedExecContext(ctx, q, params)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok {
			switch pgErr.Code {
			case pgerrcode.UndefinedTable:
				return errors.Wrap(errors.ErrDeleteMessages, err)
			case pgerrcode.InvalidTextRepresentation:
				return errors.Wrap(errors.ErrDeleteMessages, errInvalidMessage)
			default:
				return errors.Wrap(errors.ErrDeleteMessages, err)
			}
		}
		return errors.Wrap(errors.ErrDeleteMessages, err)
	}

	return nil
}

func (tr postgresRepository) Restore(ctx context.Context, messages ...senml.Message) error {
	q := `INSERT INTO messages (subtopic, publisher, protocol,
          name, unit, value, string_value, bool_value, data_value, sum,
          time, update_time)
          VALUES (:subtopic, :publisher, :protocol, :name, :unit,
          :value, :string_value, :bool_value, :data_value, :sum,
          :time, :update_time);`

	tx, err := tr.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return errors.Wrap(errors.ErrSaveMessages, err)
	}

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = errors.Wrap(err, errors.Wrap(errTransRollback, txErr))
			}
			return
		}

		if err = tx.Commit(); err != nil {
			err = errors.Wrap(errors.ErrSaveMessages, err)
		}
	}()

	for _, msg := range messages {
		m := senmlMessage{Message: msg}
		if _, err := tx.NamedExec(q, m); err != nil {
			pgErr, ok := err.(*pgconn.PgError)
			if ok && pgErr.Code == pgerrcode.InvalidTextRepresentation {
				return errors.Wrap(errors.ErrSaveMessages, errInvalidMessage)
			}
			return errors.Wrap(errors.ErrSaveMessages, err)
		}
	}

	return err
}

func (tr postgresRepository) readAll(rpm readers.PageMetadata) (readers.MessagesPage, error) {
	order := "time"
	format := defTable
	olq := dbutil.GetOffsetLimitQuery(rpm.Limit)
	interval := rpm.Interval

	if rpm.Format == jsonTable {
		order = "created"
		format = rpm.Format
	}

	var q string
	condition := fmtCondition(rpm, format)

	if interval != "" {
		switch format {
		case defTable:
			q = fmt.Sprintf(`
				SELECT * FROM (
					SELECT DISTINCT ON (date_trunc('%[1]s', to_timestamp(%[2]s))) *
					FROM %[3]s %[4]s
					ORDER BY date_trunc('%[1]s', to_timestamp(%[2]s)), %[2]s DESC
				) sub
				%[5]s;`, interval, order, format, condition, olq)

		case jsonTable:
			q = fmt.Sprintf(`
				SELECT * FROM (
					SELECT DISTINCT ON (date_trunc('%[1]s', to_timestamp(created / 1000000000))) *
					FROM %[2]s %[3]s
					ORDER BY date_trunc('%[1]s', to_timestamp(created / 1000000000)), created DESC
				) sub
				%[4]s;`, interval, format, condition, olq)
		}
	} else {
		q = fmt.Sprintf(`SELECT * FROM %s %s ORDER BY %s DESC %s;`, format, condition, order, olq)
	}

	params := map[string]interface{}{
		"limit":        rpm.Limit,
		"offset":       rpm.Offset,
		"subtopic":     rpm.Subtopic,
		"publisher":    rpm.Publisher,
		"name":         rpm.Name,
		"protocol":     rpm.Protocol,
		"value":        rpm.Value,
		"bool_value":   rpm.BoolValue,
		"string_value": rpm.StringValue,
		"data_value":   rpm.DataValue,
		"from":         rpm.From,
		"to":           rpm.To,
	}

	rows, err := tr.db.NamedQuery(q, params)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == pgerrcode.UndefinedTable {
				return readers.MessagesPage{}, nil
			}
		}
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	page := readers.MessagesPage{
		PageMetadata: rpm,
		Messages:     []readers.Message{},
	}

	switch format {
	case defTable:
		for rows.Next() {
			msg := senmlMessage{Message: senml.Message{}}
			if err := rows.StructScan(&msg); err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}
			page.Messages = append(page.Messages, msg.Message)
		}
	default:
		for rows.Next() {
			msg := jsonMessage{}
			if err := rows.StructScan(&msg); err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}
			m, err := msg.toMap()
			if err != nil {
				return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
			}
			page.Messages = append(page.Messages, m)
		}
	}

	// Count total for pagination
	if interval != "" {
		switch format {
		case defTable:
			q = fmt.Sprintf(`
				SELECT COUNT(*) FROM (
					SELECT DISTINCT ON (date_trunc('%[1]s', to_timestamp(%[2]s))) *
					FROM %[3]s %[4]s
					ORDER BY date_trunc('%[1]s', to_timestamp(%[2]s)), %[2]s DESC
				) sub;`, interval, order, format, condition)
		case jsonTable:
			q = fmt.Sprintf(`
				SELECT COUNT(*) FROM (
					SELECT DISTINCT ON (date_trunc('%[1]s', to_timestamp(created / 1000000000))) *
					FROM %[2]s %[3]s
					ORDER BY date_trunc('%[1]s', to_timestamp(created / 1000000000)), created DESC
				) sub;`, interval, format, condition)
		}
	} else {
		q = fmt.Sprintf(`SELECT COUNT(*) FROM %s %s;`, format, condition)
	}

	rows, err = tr.db.NamedQuery(q, params)
	if err != nil {
		return readers.MessagesPage{}, errors.Wrap(readers.ErrReadMessages, err)
	}
	defer rows.Close()

	total := uint64(0)
	if rows.Next() {
		if err := rows.Scan(&total); err != nil {
			return page, err
		}
	}
	page.Total = total

	return page, nil
}

func fmtCondition(rpm readers.PageMetadata, table string) string {
	var query map[string]interface{}
	meta, err := json.Marshal(rpm)
	if err != nil {
		return ""
	}
	json.Unmarshal(meta, &query)

	condition := ""
	op := "WHERE"
	timeColumn := "time"

	if table != "" && table == jsonTable {
		timeColumn = "created"
	}

	for name := range query {
		switch name {
		case "subtopic", "publisher", "name", "protocol":
			condition = fmt.Sprintf(`%s %s %s = :%s`, condition, op, name, name)
			op = "AND"
		case "v":
			comparator := readers.ParseValueComparator(query)
			condition = fmt.Sprintf(`%s %s value %s :value`, condition, op, comparator)
			op = "AND"
		case "vb":
			condition = fmt.Sprintf(`%s %s bool_value = :bool_value`, condition, op)
			op = "AND"
		case "vs":
			condition = fmt.Sprintf(`%s %s string_value = :string_value`, condition, op)
			op = "AND"
		case "vd":
			condition = fmt.Sprintf(`%s %s data_value = :data_value`, condition, op)
			op = "AND"
		case "from":
			condition = fmt.Sprintf(`%s %s %s >= :from`, condition, op, timeColumn)
			op = "AND"
		case "to":
			condition = fmt.Sprintf(`%s %s %s <= :to`, condition, op, timeColumn)
			op = "AND"
		}
	}
	return condition
}

type senmlMessage struct {
	senml.Message
}

type jsonMessage struct {
	Created   int64  `db:"created"`
	Subtopic  string `db:"subtopic"`
	Publisher string `db:"publisher"`
	Protocol  string `db:"protocol"`
	Payload   []byte `db:"payload"`
}

func (msg jsonMessage) toMap() (map[string]interface{}, error) {
	ret := map[string]interface{}{
		"created":   msg.Created,
		"subtopic":  msg.Subtopic,
		"publisher": msg.Publisher,
		"protocol":  msg.Protocol,
		"payload":   map[string]interface{}{},
	}
	pld := make(map[string]interface{})
	if err := json.Unmarshal(msg.Payload, &pld); err != nil {
		return nil, err
	}
	ret["payload"] = pld
	return ret, nil
}
