GOOS=linux GOARCH=amd64 go build -o bin/main .
serverless deploy --stage local
serverless deploy --verbose --stage local
curl http://localhost:4566/restapis/jns8lfvku1/local/_user_request_/hello
sls deploy —- stage local
sls invoke -f hello --stage local


aws --endpoint-url=http://localhost:4566 dynamodb create-table \
  --attribute-definitions \
    AttributeName=userId,AttributeType=N \
  --key-schema \
    AttributeName=userId,KeyType=HASH \
  --table-name dynamo-stream-to-sns \
  --provisioned-throughput \
    ReadCapacityUnits=1,WriteCapacityUnits=1 \
    --region us-east-2
    LocationConstraint=us-east-2
aws --endpoint-url=http://localhost:4566 dynamodb describe-table --table-name dynamo-stream-to-sns --region us-east-2 


aws --endpoint-url=http://localhost:4566 s3api create-bucket --bucket deployment-bucket --region us-east-2 --create-bucket-configuration LocationConstraint=us-east-2
aws --endpoint-url=http://localhost:4566 s3api get-bucket-location --bucket deployment-bucket

aws --endpoint-url=http://localhost:4566 sns create-topic --name ExampleTopic

aws --endpoint-url=http://localhost:4566 dynamodb put-item \
--table-name dynamo-stream-to-sns  \
--item \
    '{"userId": {"N": "1"}}' \
--region us-east-2 