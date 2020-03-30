package storage

import "database/sql"

// StopSearch - disable receive message about new offers.
func (c *Connector) StopSearch(id int64) error {
	resp, err := c.DB.Exec("UPDATE chat SET enable = 0 WHERE id = ?;", id)
	if err != nil {
		return err
	}

	affect, err := resp.RowsAffected()
	if err != nil {
		return err
	}

	if affect == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// StartSearch - register new user or group if not exist or enable receive new
//  offers.
func (c *Connector) StartSearch(id int64, username, title, cType string) error {
	chat, err := c.ReadChat(id)
	if err != nil {
		if err == sql.ErrNoRows {
			_, err := c.DB.Exec(
				`INSERT INTO chat (id, username, title, enable, c_type)
						VALUES (?, ?, ?, ?, ?);`,
				id,
				username,
				title,
				true,
				cType,
			)
			return err
		}
		return err
	}

	_, err = c.DB.Exec("UPDATE chat SET enable = 1 WHERE id = ?;", chat.Id)
	return err
}
