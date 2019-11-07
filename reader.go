package graphdb

func (c *Client) ReadSingleRow(cypher string, params map[string]interface{}, blankTypes ...interface{}) ([]interface{}, error) {
	row, err := c.QuerySingle(cypher, params, convertRecordToTypesFunc(blankTypes))
	if err != nil {
		return nil, err
	}
	if conversionErr, ok := row.(error); ok {
		return nil, conversionErr
	}
	if row == nil {
		return nil, nil
	}
	return row.([]interface{}), nil
}

func (c *Client) ReadRows(cypher string, params map[string]interface{}, blankTypes ...interface{}) ([][]interface{}, error) {
	var validatedRows [][]interface{}
	rows, err := c.Query(cypher, params, convertRecordToTypesFunc(blankTypes))
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if conversionErr, ok := row.(error); ok {
			return nil, conversionErr
		}
		if row == nil {
			return nil, nil
		}
		validatedRows = append(validatedRows, row.([]interface{}))
	}
	return validatedRows, nil
}
