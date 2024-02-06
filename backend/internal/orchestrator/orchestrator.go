package orchestrator

import "time"

type Expression struct {
	Id int `json:"id"` /* Уникальный идентификатор (time * 100 + 0 < N < 100) */

	Expression string `json:"expression"` /* арифметическое выражение */

	Status string `json:"status"` /* статус выражения */

	CreatedTime time.Time `json:"created_time"` /* время создания запроса */

	CompletedTime time.Time `json:"updated_time"` /* время завершения вычислений */

	/*
		Возможные варианты ответа:
		*200. Выражение успешно принято, распаршено и принято к обработке
		*400. Выражение невалидно
		*500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
	*/
}

type Operation struct {
	Name     string        `json:"name"`     /* имя арифметической операции */
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
	Wait = ""
)
