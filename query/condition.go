package query

const (
	LogicAnd    = uint8(1)
	LogicOr     = uint8(2)
	LogicAndNot = uint8(3)
)

type conditionExp struct {
	query interface{}
	args  []interface{}
	logic uint8
}

type Condition struct {
	exp []conditionExp
}

func (cond *Condition) Where(query interface{}, values ...interface{}) {
	cond.addCondition(query, LogicAnd, values...)
}

func (cond *Condition) Or(query interface{}, values ...interface{}) {
	cond.addCondition(query, LogicOr, values...)
}

func (cond *Condition) Not(query interface{}, values ...interface{}) {
	cond.addCondition(query, LogicAndNot, values...)
}

func (cond *Condition) addCondition(query interface{}, logic uint8, values ...interface{}) {
	cond.exp = append(cond.exp, conditionExp{
		query: query,
		args:  values,
		logic: logic,
	})
}

func (cond *Condition) Empty() {
	cond.exp = []conditionExp{}
}

func (cond *Condition) GetExp() []conditionExp {
	return cond.exp
}
func (cond *Condition) SetExp(exp []conditionExp) {
	cond.exp = exp
}

func (cond *Condition) hasCondition() bool {
	return len(cond.exp) > 0
}
