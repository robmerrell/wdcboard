package models

// Index adds indexes to the database
func Index(conn *MgoConnection) error {
	prices := mainConnection.DB.C("prices")
	if err := prices.EnsureIndexKey("generatedAt"); err != nil {
		return err
	}

	return nil
}
