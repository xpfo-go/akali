{
    "swagger": "2.0",
    "info": {
        "contact": {}
    },
    "paths": {
        "/health": {
            "get": {
                "description": "/health to make sure the server is health",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "basic"
                ],
                "summary": "health for server health check",
                "operationId": "health",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        },
                        "headers": {
                            "X-Request-Id": {
                                "type": "string",
                                "description": "the request id"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "description": "/ping to get response from <xpfo{ .ProjectName }xpfo>, make sure the server is alive",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "basic"
                ],
                "summary": "ping-pong for alive test",
                "operationId": "ping",
                "responses": {
                    "200": {
                        "description": "OK",
                        "headers": {
                            "X-Request-Id": {
                                "type": "string",
                                "description": "the request id"
                            }
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "/version to get the version of iam",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "basic"
                ],
                "summary": "version for identify",
                "operationId": "version",
                "responses": {
                    "200": {
                        "description": "OK",
                        "headers": {
                            "X-Request-Id": {
                                "type": "string",
                                "description": "the request id"
                            }
                        }
                    }
                }
            }
        }
    }
}