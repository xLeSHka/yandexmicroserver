package orchestrator

import "time"

type Expression struct {
	Id int `json:"id"` /* Уникальный идентификатор (time * 100 + 0 < N < 100) */

	Expression string `json:"expression"` /* арифметическое выражение */

	Status string `json:"status"` /* статус выражения */

	CreatedTime time.Time `json:"created_time"` /* время создания запроса */

	CompletedTime time.Time `json:"completed_time"` /* время завершения вычислений */

	ExecutionTimeByMilliseconds int `json:"execution_time"`

	/*
		Возможные варианты ответа:
		*200. Выражение успешно принято, распаршено и принято к обработке
		*400. Выражение невалидно
		*500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
	*/
}

type Operation struct {
	Name     rune          `json:"name"`     /* арифметическая операция */
	Duration time.Duration `json:"duration"` /* Время выполнения операции*/
}
type Todo interface {
	addExpression()
	getExpressions()
	getExpressionByID()
	getOperations()
	getTask()
	getResult()
}

const (
	Wait  = "the exprassion will be calculated soon"
	Ready = "the expression calculated"
	Error = "exprassion parsing error"
)
