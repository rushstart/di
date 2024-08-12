package di

func must1[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

func use(a any) {

}
