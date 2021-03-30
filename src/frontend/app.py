from flask import Flask, render_template, request, flash
import json
import os
import requests
import logging

app = Flask(__name__)
app.secret_key = b'_5#y2L"F4Q8z\n\xec]/'

gunicorn_logger = logging.getLogger('gunicorn.error')
app.logger.handlers = gunicorn_logger.handlers
app.logger.setLevel(gunicorn_logger.level)

backendHostname = "encryptah-be" if (os.environ.get("BACKEND_HOSTNAME") is None) else os.environ.get("BACKEND_HOSTNAME")
backendPort = "5678" if (os.environ.get("BACKEND_PORT_NUMBER") is None) else os.environ.get("BACKEND_PORT_NUMBER")

backend = {
  "name": "http://{0}:{1}".format(backendHostname, backendPort),
  "endpoint": "encryptah-be",
}

@app.route('/', methods=['GET', 'POST'])
def index():
  ciphertext = ''
  plaintext = ''

  # Ensure connectivity to backend service
  try:
    app.logger.info('trying backend at {}'.format(backend['name']))
    requests.get(backend['name'] + '/health')
    app.logger.info('successfully connected to backend at {}'.format(backend['name']))
  except:
    app.logger.error('failed connecting to backend at {}'.format(backend['name']))
    flash ('Cannot connect to backend service.')


  if request.method == 'POST':
    if request.form['submit'] == 'encrypt':
      # send a POST request to the backend service with the plaintext data to get encrypted
      payload = {'plaintext': request.form['message']}
      response = requests.post(backend['name'] + '/encrypt', data = json.dumps(payload))

      # retrieve the ciphertext key in JSON response
      ciphertext = response.json()['ciphertext']
    elif request.form['submit'] == 'decrypt':
      # send a POST request to the backend service with the ciphertext data to get decrypted
      payload = {'ciphertext': request.form['message']}
      response = requests.post(backend['name'] + '/decrypt', data = json.dumps(payload))

      # retrieve the plaintext key in JSON response
      plaintext = response.json()['plaintext']
    else:
      # need to add failure handling here
      pass
  elif request.method == 'GET':
    # temporary pass
    pass

  return render_template('index.html', ciphertext=ciphertext, plaintext=plaintext)

@app.route('/health')
def health():
  return json.loads('{ "status": "OK" }')

if __name__ == '__main__':
  app.run(debug = True, host = '0.0.0.0', port = '8080')
