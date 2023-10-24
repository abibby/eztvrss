package models

// func FirstOrCreate[T model.Model](tx *sqlx.Tx, q *builder.Builder[T], defaultValue T) (T, error) {
// 	s, err := q.First(tx)
// 	if err != nil {
// 		var zero T
// 		return zero, err
// 	}
// 	if reflect.ValueOf(s).IsZero() {
// 		return s, nil
// 	}
// 	err = model.Save(tx, defaultValue)
// 	if err != nil {
// 		var zero T
// 		return zero, err
// 	}
// 	return defaultValue, nil
// }
