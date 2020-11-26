# build
GOOS=linux GOARCH=amd64 go build -o bin/main .
# verbose deploy
serverless deploy --verbose --stage local
# use generated endpoint, it may vary
curl http://localhost:4566/restapis/jns8lfvku1/local/_user_request_/hello
# lambda invocation
sls invoke -f hello --stage local --path resources/api_gateway_request.json

