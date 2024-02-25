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
`curl -X POST -d "2 + 2" http://localhost:8082/add`
 ## Post operation`s executing time
 Выставить время выполнения операции, где name - арифметическая операция, execution_time_by_milliseconds - время выполнения операции в милисекундах

 `curl -X POST -H 'Content-Type:application/json' -d "{""name"": ""*"", ""execution_time_by_milliseconds"": 3000}" "http://localhost:8082/setOperation" ` `#0969DA`


curl -X POST -d "2 + 2" http://localhost:8082/add
curl -X GET http://localhost:8082/initialize
curl -X GET http://localhost:8082/