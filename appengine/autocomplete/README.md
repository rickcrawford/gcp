autocomplete
------------


App engine sample autocomplete engine


Create test products

```bash


curl -H "Content-Type: application/json" -X POST https://typeahead-183622.appspot.com/catalog -d '{"catalogId":"bestbuy.com","name":"Best Buy Test Catalog", "batch":"11111"}'


curl -H "Content-Type: application/json" -X POST http://localhost:8080/catalog -d '{"catalogId":"bestbuy.com","name":"Best Buy Test Catalog", "batch":"1"}'

curl -H "Content-Type: application/json" -X POST http://localhost:8080/catalog/test/product/batch -d '[{"productId":"1","name":"test1"},{"productId":"2","name":"test2"},{"productId":"3","name":"test3"},{"productId":"4","name":"test4"}]'

```