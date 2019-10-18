# Swagen
Swagen reads your MySQL database and generates Swagger/OpenAPI YAML from table meta data. This results in a boilerplate swaggerfile (YML) that should probably be edited to completely suit your needs. Use [go-swagger](https://github.com/go-swagger/go-swagger) or [swagger-codegen](https://github.com/swagger-api/swagger-codegen) to generate the API boilerplate Go code.

## Swagger version
Swagen generates Swagger/OpenAPI 2.0 configuration.

## Installation
Swagen can be acquired and installed by running

```
go get -v -u github.com/minitauros/swagen
```

## How to build
Or, alternatively, the binary can be manually built:

```
go build -o ~/bin/swagen main.go
```

## Usage

### Creating a config file
You will need to create a small config file that Swagen uses to generate YML. For example:

```
db:
  dns: user:password@tcp(host:port)/db_name?parseTime=true

service:
  name: MyService
  host: 1.1.1.1:1234

resources:
  table_name:
    # The written name of the resource, e.g. "product". Will be used in summaries.
    title: written name of resource
    definition:
      # This is how the resource will be named in the definitions list.
      name: ResourceDefinition
  # Continue listing all the tables you want to generate Swagger YAML for.
```

### Running Swagen
Swagger YML can then be generated using the following command:

```
swagen -conf=conf.yml > swagger.yml
```

### Example
Products table:

```sql
CREATE TABLE products (
    id INT NOT NULL AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    price INT NOT NULL,
    PRIMARY KEY (id)
); 
```

Config:

```
db:
  dns: root:@tcp(0.0.0.0:3306)/products?parseTime=true

service:
  name: TestService
  host: 1.1.1.1:1234

resources:
  products:
    title: product
    definition:
      name: Product
```

Result:

```
swagger: '2.0'
info:
  title: TestService
  version: 1.0.0
host: 1.1.1.1:1234
basePath: /v1
schemes:
  - http
consumes:
  - application/json
produces:
  - application/json
paths:
  /products:
    get:
      summary: Returns the product resources with the given IDs
      parameters:
        - in: query
          name: ids
          type: array
          items:
            type: integer
      responses:
        200:
          description: List of product resources
          schema:
            type: array
            items:
              $ref: '#/definitions/Product'
        500:
          description: Internal server error
    post:
      summary: Creates a product
      parameters:
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/Product'
      responses:
        200:
          description: Success
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
        500:
          description: Internal server error
  /products/{id}:
    get:
      summary: Returns the product with the given ID
      parameters:
        - in: path
          name: id
          type: integer
          required: true
      responses:
        200:
          description: Single product
          schema:
            $ref: '#/definitions/Product'
        404:
          description: Not found
        500:
          description: Internal server error
    patch:
      summary: Patches the product with the given ID
      parameters:
        - name: id
          in: path
          type: integer
          required: true
        - name: patch
          in: body
          required: true
          schema:
            $ref: '#/definitions/Patch'
      responses:
        200:
          description: Success
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
        404:
          description: Not found
        500:
          description: Internal server error
    put:
      summary: Replaces the product with the given ID
      parameters:
        - name: id
          in: path
          type: integer
          required: true
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/Product'
      responses:
        200:
          description: Success
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
        404:
          description: Not found
        500:
          description: Internal server error
    delete:
      summary: Deletes the product with the given ID
      parameters:
        - name: id
          in: path
          type: integer
          required: true
      responses:
        200:
          description: Success
        404:
          description: Not found
        500:
          description: Internal server error
definitions:
  Product:
    properties:
      id:
        type: integer
      name:
        type: string
      price:
        type: integer
  Patch:
    type: array
    description: Patch instructions
    items:
      type: object
      required:
        - op
        - path
        - value
      properties:
        op:
          type: string
          description: Operation
        path:
          type: string
          description: Path to field to operate on
        value:
          $ref: '#/definitions/AnyValue'
  AnyValue:
    description: Any type of value
  Error:
    type: object
    properties:
      message:
        type: string

```

## Tests
This package does not contain tests. Since it would be foolish to run this directly in any production environment, I will assume that every time this tool might be used, the person using it will check and modify the output (as the tool attow only generates boilerplate YML). On top of that, if invalid YML is generated, other tools (like swagger-codegen) will break, also spotting any defects.