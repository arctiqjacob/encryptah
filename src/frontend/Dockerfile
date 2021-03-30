# Use the official lightweight Python image.
FROM python:3.9-slim

# Copy local code to the container image.
COPY app /app

# Install production dependencies.
RUN pip install -r /app/requirements.txt

# Listens on port 8080
EXPOSE 8080

# Run the web service on container startup. Here we use the gunicorn
# webserver, with one worker process and 8 threads.
CMD exec gunicorn --bind :8080 --workers 1 --threads 8 --timeout 0 app.app:app