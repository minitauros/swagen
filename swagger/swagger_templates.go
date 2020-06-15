package swagger

var fullTemplate = `swagger: '2.0'
info:
  title: {{ .ServiceInfo.Name }}
  version: 1.0.0
host: {{ .ServiceInfo.Host }}
basePath: /v1
schemes:
  - http
consumes:
  - application/json
produces:
  - application/json
paths:
  /health-check:
    get:
      operationId: HealthCheck
      description: Returns 200 if the service is healthy.
      responses:
        200:
          description: Healthy
        500:
          description: Not healthy
  {{- .Resources }}
definitions:
  {{- .Definitions }}
  Patch:
    type: array
    description: Patch instructions
    items:
      type: object
      required:
        - op
        - path
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
      code:
        type: integer
        format: int64
        x-nullable: true
      message:
        type: string
  Principal:
    type: object
    description: Security principal for validating that a user is authorized to execute certain actions
    properties:
      userId:
        type: string
      permissions:
        type: array
        items:
          type: string
`

var resourceTemplate = `
  /{{ .Path }}:
    get:
      operationId: Get{{ .Definition.Name }}s
      summary: Get {{ .Title }}s
      description: Returns all {{ .Title }} resources.
      parameters:
        - in: query
          name: limit
          type: integer
        - in: query
          name: offset
          type: integer
      responses:
        200:
          description: List of {{ .Title }} resources.
          schema:
            type: array
            items:
              $ref: '#/definitions/{{ .Definition.Name }}'
        500:
          description: Internal server error
    post:
      operationId: Create{{ .Definition.Name }}
      summary: Create {{ .Title }}
      description: Creates a {{ .Title }}.
      parameters:
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/{{ .Definition.Name }}Input'
      responses:
        201:
          description: Created
          schema:
            properties:
              id:
                type: integer
                description: The ID of the {{ .Title }} that was created.
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
        422:
          description: Unprocessable entity
          schema:
            $ref: '#/definitions/Error'
        500:
          description: Internal server error
  /{{ .Path }}/batch:
    get:
      operationId: Get{{ .Definition.Name }}sByID
      summary: Get {{ .Title }}s by ID
      description: Returns the {{ .Title }} resources with the given IDs.
      parameters:
        - in: query
          name: ids
          type: array
          items:
            type: integer
      responses:
        200:
          description: List of {{ .Title }} resources
          schema:
            type: array
            items:
              $ref: '#/definitions/{{ .Definition.Name }}'
        500:
          description: Internal server error
  /{{ .Path }}/{id}:
    get:
      operationId: Get{{ .Definition.Name }}
      summary: Get {{ .Title }} by ID
      description: Returns the {{ .Title }} with the given ID.
      parameters:
        - in: path
          name: id
          type: integer
          required: true
      responses:
        200:
          description: Single {{ .Title }}
          schema:
            $ref: '#/definitions/{{ .Definition.Name }}'
        404:
          description: Not found
        500:
          description: Internal server error
    patch:
      operationId: Patch{{ .Definition.Name }}
      summary: Patch {{ .Title }}
      description: Patches the {{ .Title }} with the given ID.
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
        422:
          description: Unprocessable entity
          schema:
            $ref: '#/definitions/Error'
        500:
          description: Internal server error
    put:
      operationId: Put{{ .Definition.Name }}
      summary: Put {{ .Title }}
      description: Replaces the {{ .Title }} with the given ID.
      parameters:
        - name: id
          in: path
          type: integer
          required: true
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/{{ .Definition.Name }}Input'
      responses:
        200:
          description: Success
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
        404:
          description: Not found
        422:
          description: Unprocessable entity
          schema:
            $ref: '#/definitions/Error'
        500:
          description: Internal server error
    delete:
      operationId: Delete{{ .Definition.Name }}
      summary: Delete {{ .Title }}
      description: Deletes the {{ .Title }} with the given ID.
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
          description: Internal server error`

var definitionTemplate = `
  {{ .Name }}:
    properties:
      {{- range $index, $field := .Fields }}
      {{ $field.Name }}:
        type: {{ $field.Type.Name }}
        {{- if $field.Type.ExtraProperties }}
          {{- range $key, $value := $field.Type.ExtraProperties }}
        {{ $key }}: {{ $value }}
          {{- end }}
        {{- end }}
		{{- if $field.IsNullable }}
        x-nullable: true
		{{- end }}
      {{- end }}
  {{ .Name }}Input:
    required:
    {{- range $index, $field := .Fields }}
      {{- if and (ne $field.Name "id") (not $field.IsNullable) }}
      - {{ $field.Name }}
      {{- end }}
    {{- end }}
    properties:
      {{- range $index, $field := .Fields }}
        {{- if  ne $field.Name "id" }}
      {{ $field.Name }}:
        type: {{ $field.Type.Name }}
        {{- if $field.Type.ExtraProperties }}
          {{- range $key, $value := $field.Type.ExtraProperties }}
        {{ $key }}: {{ $value }}
          {{- end }}
        {{- end }}
		{{- if $field.IsNullable }}
        x-nullable: true
		{{- end }}
        {{- end }}
      {{- end }}`
