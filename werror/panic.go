package werror

type Panic struct {
	Msg string
	Err error
}

func WformPanic(str interface{}) {
	switch str.(type) {
	case string:
		panic(Panic{Msg: str.(string)})
	case error:
		panic(Panic{Err: str.(error)})
	}

}
