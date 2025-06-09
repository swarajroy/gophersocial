A 3-tier architecture mono repo containing

backend built with
Go
Postgresql
Docker

frontend built with
React
Typescript


autocannon -- to perf test with redis on
npx autocannon http://localhost:8080/v1/users/51 --connections 10 --duration 5 -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJnb3BoZXJzb2NpYWwiLCJleHAiOjE3NDk1NTI1NDIsImlhdCI6MTc0OTQ2NjE0MiwiaXNzIjoiZ29waGVyc29jaWFsIiwibmJmIjoxNzQ5NDY2MTQyLCJzdWIiOjUxfQ.fGoVTiyuCSkO0KNtA8kZmduCPel51FgaMlJ4CLjEQSQ"
