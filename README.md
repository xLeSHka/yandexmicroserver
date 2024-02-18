# yandexmicroserver
 Distributed calculator of arithmetic expressions
curl -X POST -H 'Content-Type:application/json' -d "{""name"": ""*"", ""execution_time_by_milliseconds"": 3000}" "http://localhost:8082/setOperation"
curl -X POST -d "2 + 2" http://localhost:8082/add
curl -X GET http://localhost:8082/initialize
curl -X GET http://localhost:8082/