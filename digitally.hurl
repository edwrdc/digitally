# GET Healthcheck
GET http://localhost:6969/v1/healthz
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.status" == "OK"
jsonpath "$.version" == "0.0.1"
jsonpath "$.env" == "development"

# POST Create Product
POST http://localhost:6969/v1/products
Content-Type: application/json
{
    "name": "Sample Product",
    "price": "99.99",
    "description": "This is a sample product description",
    "categories": ["electronics", "gadgets"]
}

# GET Product by ID
GET http://localhost:6969/v1/products/1
HTTP 200
