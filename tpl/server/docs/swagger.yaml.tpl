info:
  contact: {}
paths:
  /health:
    get:
      consumes:
      - application/json
      description: /health to make sure the server is health
      operationId: health
      produces:
      - application/json
      responses:
        "200":
          description: OK
          headers:
            X-Request-Id:
              description: the request id
              type: string
          schema:
            type: string
        "500":
          description: Internal Server Error
          schema:
            type: string
      summary: health for server health check
      tags:
      - basic
  /ping:
    get:
      consumes:
      - application/json
      description: /ping to get response from <xpfo{ .ProjectName }xpfo>, make sure the server is
        alive
      operationId: ping
      produces:
      - application/json
      responses:
        "200":
          description: OK
          headers:
            X-Request-Id:
              description: the request id
              type: string
      summary: ping-pong for alive test
      tags:
      - basic
  /version:
    get:
      consumes:
      - application/json
      description: /version to get the version of iam
      operationId: version
      produces:
      - application/json
      responses:
        "200":
          description: OK
          headers:
            X-Request-Id:
              description: the request id
              type: string
      summary: version for identify
      tags:
      - basic
swagger: "2.0"
