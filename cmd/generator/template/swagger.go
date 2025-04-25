package template

const Swagger = `openapi: 3.0.0
info:
  title: goserve
  description: A generated base goserver project
  version: 1.0.0
servers:
  - url: http://localhost:8080/
    description: Local host app endpoint

paths:
  /hello:
    get:
      tags:
        - Hello
      summary: Hello resource
      operationId: hello
      parameters:
        - $ref: "#/components/parameters/username"
      responses:
        '200':
          description: Operation executed successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/BaseResponse'
        '400':
          description: Invalid input
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '401':
          description: Unauthorized
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '500':
          description: Unexpected error
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  parameters:
    username:
      name: username
      in: query
      required: true
      schema:
        type: string
        example: goserve

  schemas:
    BaseResponse:
      type: object
      properties:
        message:
          type: string
          example: "Operation executed successful"
        timestamp:
          type: integer
          format: int64
    ErrorResponse:
      type: object
      allOf:
        - $ref: "#/components/schemas/BaseResponse"
        - type: object
          properties:
            errors:
              type: array
              items:
                type: string
                example: "Offensive language"

  securitySchemes:
    auth:
      type: oauth2
      flows:
        implicit:
          authorizationUrl: http://go-timeline/api/v1/authorization
          scopes:
            admin:resource:usage: Role allowed for admin resource access
`
