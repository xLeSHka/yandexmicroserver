# yandexmicroserver
 Distributed calculator of arithmetic expressions

 ## Starting
  ### Step 1 
   in terminal 
   `go run .\backend\cmd\app\main.go` 
  ### Step 2
  in terminal
   `go run .\backend\cmd\agent\agent.go` 
  ### Step 3
   in cmd
   `curl -X GET http://localhost:8082/initialize`
 ## Post expression
 Добавить выражение, в скобках указывать выражение
 > Обязательно пробелы между числами и операциями  
 ### `curl -X POST -d "2 + 2" http://localhost:8082/add` 
 ## Post operation`s executing time 
 Выставить время выполнения операции
 ### `curl -X POST -H 'Content-Type:application/json' -d "{""name"": ""*"", ""execution_time_by_milliseconds"": 3000}" "http://localhost:8082/setOperation" `
 - `name` - арифметическая операция(`+` `-` `*` `/`)
 - `execution_time_by_milliseconds` - время выполнения операции в милисекундах
 
 ## Get all expressions
 Получить все выражения, в виде json с полями 
 - `id` string - уникальный id выражения
 - `expression`string - выражение
 - `expression_status`string - статус выражения(wait, proccess, completed)
 - `created_at` time.Time - когда было создано выражение
 - `completed_at` time.Time - когда выражение было посчитанно
 ## Get all operations
 Получить все операции, в виде json c полями
 - `name` string - арифметическая операция
 - `execution_time_by_milliseconds` time.Duration - время выполнения операции в милисекундах
 ## Get all agents
 Получить всех агентов, в виде json c полями
 - `id` string - уникальный id агента
 - `address` string - адресс агента
 - `status_code` string - статус агента
 - `last_heartbeat` time.Time - время последнего heartbeat`a агента
curl -X GET http://localhost:8082/