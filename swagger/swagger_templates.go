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
      summary: Returns the {{ .Title }} resources with the given IDs, or all of them if no IDs are given
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
    post:
      operationId: Create{{ .Definition.Name }}
      summary: Creates a {{ .Title }}
      parameters:
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/{{ .Definition.Name }}'
      responses:
        200:
          description: Success
          schema:
            properties:
              id:
                type: integer
                description: The ID of the {{ .Title }} that was created
        400:
          description: Bad request
          schema:
            $ref: '#/definitions/Error'
		
        500:
          description: Internal server error
  /{{ .Path }}/{id}:
    get:
      operationId: Get{{ .Definition.Name }}
      summary: Returns the {{ .Title }} with the given ID
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
      summary: Patches the {{ .Title }} with the given ID
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
      operationId: Put{{ .Definition.Name }}
      summary: Replaces the {{ .Title }} with the given ID
      parameters:
        - name: id
          in: path
          type: integer
          required: true
        - name: resource
          in: body
          required: true
          schema:
            $ref: '#/definitions/{{ .Definition.Name }}'
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
      operationId: Delete{{ .Definition.Name }}
      summary: Deletes the {{ .Title }} with the given ID
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
      {{- end }}`
