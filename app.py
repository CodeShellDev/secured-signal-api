from flask import Flask, Response, request, jsonify, make_response
import os
import json
import requests
import re
import base64
import logging
from urllib.parse import unquote, urlencode, parse_qs

app = Flask("Secured Signal Api")

app.logger.setLevel(logging.INFO)

DEFAULT_BLOCKED_ENDPOINTS = [
    "/v1/about",
    "/v1/configuration",
    "/v1/devices",
    "/v1/register",
    "/v1/unregister",
    "/v1/qrcodelink",
    "/v1/accounts",
    "/v1/contacts"
]

SENDER = os.getenv("SENDER")
DEFAULT_RECIPIENTS = os.getenv("DEFAULT_RECIPIENTS")
SIGNAL_API_URL = os.getenv("SIGNAL_API_URL")
API_TOKEN = os.getenv("API_TOKEN")

BLOCKED_ENDPOINTS = os.getenv("BLOCKED_ENDPOINTS")
VARIABLES = os.getenv("VARIABLES")

secure = False

def fillInVars(obj):
    if isinstance(obj, dict):
        for key, value in obj.items():
            obj[key] = fillInVars(value)
    elif isinstance(obj, list):
        for i in range(len(obj)):
            obj[i] = fillInVars(obj[i])
    elif isinstance(obj, str):
        matches = re.findall(r"\${(.*?)}", obj)
        for match in matches:
            if match in VARIABLES:
                value = VARIABLES[match]
            
                if isinstance(value, str):
                    newValue = obj.replace(f"${{{match}}}", str(value))
                    return newValue
                else:
                    return value
    return obj

def UnauthorizedResponse(prompt=None):
    headers = {}

    if prompt:
        headers = {
            "WWW-Authenticate": 'Basic realm="Login Required", Bearer realm="Access Token Required"'
        }

    return Response(
        "Unauthorized", 401,
        headers
    )

@app.before_request
def middlewares():
    for blockedPath in BLOCKED_ENDPOINTS:
        if request.path.startswith(blockedPath):
            infoLog(f"Client tried to access Blocked Endpoint [{blockedPath}]")
            return Response("Forbidden", 401)

    query_string = request.query_string.decode()

    if secure:
        auth_header = request.headers.get("Authorization", "")
        
        if auth_header.startswith("Bearer "):
            token = auth_header.split(" ", 1)[1]

            token = unquote(token)

            if token != API_TOKEN:
                infoLog(f"Client failed Bearer Auth [token: {token}]")
                return UnauthorizedResponse()
        elif auth_header.startswith("Basic "):
            try:
                decoded = base64.b64decode(auth_header.split(" ", 1)[1]).decode()
                username, password = decoded.split(":", 1)

                username = unquote(username)
                password = unquote(password)

                if username != "api" or password != API_TOKEN:
                    infoLog(f"Client failed Basic Auth [user: {username}, pw:{password}]")
                    return UnauthorizedResponse()
            except Exception as error:
                errorLog(f"Unexpected Error during Basic Auth: {error}")
                return UnauthorizedResponse()
        elif request.args.get("authorization", None):
            token = request.args.get("authorization", "")

            token = unquote(token)

            if token != API_TOKEN:
                infoLog(f"Client failed Query Auth [query: {token}]")
                return UnauthorizedResponse()

            args = parse_qs(query_string)

            args.pop('authorization', None)
            query_string = urlencode(args, doseq=True)
        else:
            infoLog(f"Client did not provide any Auth Method")
            return UnauthorizedResponse(True)

    g.query_string = query_string  

@app.route('/', defaults={'path': ''}, methods=['GET', 'POST', 'PUT'])
@app.route('/<path:path>', methods=['GET', 'POST', 'PUT'])
def proxy(path):
    method = request.method
    incomingJSON = request.get_json(force=True, silent=True)
    jsonData = incomingJSON
    headers = {k: v for k, v in request.headers if k.lower() != 'host'}

    if incomingJSON:
        jsonData = fillInVars(incomingJSON)

    if "${NUMBER}" in path:
        path = path.replace("${NUMBER}", SENDER)

    query_string = g.query_string

    if query_string:
        query_string = "?" + query_string

    targetURL = f"{SIGNAL_API_URL}/{path}{query_string}"

    resp = requests.request(
        method=method,
        url=targetURL,
        headers=headers,
        json=jsonData
    )

    infoLog(f"Forwarded {resp.text} to {targetURL} [{method}]")

    # return Response(resp.content, status=resp.status_code, headers=dict(resp.headers))

    response = make_response(resp.json())
    response.status_code = resp.status_code

    return response

def infoLog(msg):
    app.logger.info(msg)

def errorLog(msg):
    app.logger.error(msg)

if __name__ == '__main__':
    if SENDER and SIGNAL_API_URL:
        if API_TOKEN == None or API_TOKEN == "":
            infoLog("No API Token set (API_TOKEN), this is not recommended!")
        else:
            secure = True
            infoLog("API Token set, use Bearer or Basic Auth (user: api) for authentication")
        
        if DEFAULT_RECIPIENTS != None and DEFAULT_RECIPIENTS != "":
            DEFAULT_RECIPIENTS = json.loads(DEFAULT_RECIPIENTS)

        if BLOCKED_ENDPOINTS != None and BLOCKED_ENDPOINTS != "":
            BLOCKED_ENDPOINTS = json.loads(BLOCKED_ENDPOINTS)
        else:
            BLOCKED_ENDPOINTS = DEFAULT_BLOCKED_ENDPOINTS

        if VARIABLES != None and VARIABLES != "":
            VARIABLES = json.loads(VARIABLES)
        else:
            VARIABLES = {
                "NUMBER": SENDER,
                "RECIPIENTS": DEFAULT_RECIPIENTS
            }

        app.run(debug=False, port=8880, host='0.0.0.0')