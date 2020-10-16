## sources
```
- https://dgraph.io/docs/badger/get-started/
- https://github.com/dgraph-io/badger
```


## get badger kv value
```
//kv store - in memory mode
curl --location --request GET 'localhost:8080/item/volatile'
//kv store - persistence mode
curl --location --request GET 'localhost:8080/item/persistence'
```

## add new item to kv store
```
curl --location --request POST 'localhost:8080/item' \
--header 'Content-Type: application/json' \
--data-raw '{
    "file_name":"sample-file.csv",
	"user_id":123456,
	"known_category":false,
	"status":"pending",
	"file_path":"report/22334455.csv"
}'
```

## update item in kv store
```
curl --location --request POST 'localhost:8080/item/update' \
--header 'Content-Type: application/json' \
--data-raw '{
    "file_name":"inactive-gunpla-kit6.csv",
	"user_id":123456,
	"known_category":false,
	"status":"in-progress",
	"file_path":"report/22334455.csv",
    "index": "m_16_10_2020_10_59_01"
}'
```

## set new item with ttl 15seconds
```
curl --location --request POST 'localhost:8080/item/ttl' \
--header 'Content-Type: application/json' \
--data-raw '{
    "file_name":"repot-ttl.csv",
	"user_id":123456,
	"known_category":false,
	"status":"pending",
	"file_path":"report/22334455.csv"
}'
```