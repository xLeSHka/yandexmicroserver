package orchestrator

import "time"

type Expression struct {
	Id string `json:"id"` /* Уникальный идентификатор (time * 100 + 0 < N < 100) */

	Expression string `json:"expression"` /* арифметическое выражение */

	Status string `json:"expression_status"` /* статус выражения */

	CreatedTime time.Time `json:"created_at"` /* время создания запроса */

	/*
		Возможные варианты ответа:
		*200. Выражение успешно принято, распаршено и принято к обработке
		*400. Выражение невалидно
		*500. Что-то не так на back-end. В качестве ответа нужно возвращать id принятного к выполнению выражения.
	*/
}

type Operation struct {
	Operation                   string `json:"name"` /* арифметическая операция */
	ExecutionTimeByMilliseconds int    `json:"execution_time_by_milliseconds"`
}
