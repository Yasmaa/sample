package templateutils

func Add(in ...interface{}) (int64, error) {
	r := int64(0)

	for _, v := range in {
		vv, err := ToInt64(v)
		if err != nil {
			return 0, err
		}

		r += vv
	}
	return r, nil
}
