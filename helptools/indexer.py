import requests
import json

url = "http://54.224.191.25:8000/rawrequest"

payload = json.dumps({
  "request": {
    "type": "ListTransaction",
    "count": 1
  }
})
headers = {
  'Content-Type': 'application/json'
}

response = requests.request("POST", url, headers=headers, data=payload)
data = json.loads(response.text)["data"]

print(data)
