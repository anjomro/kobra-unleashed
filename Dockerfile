# Use node.js image as builder to build the frontend
FROM node:21-alpine as builder
WORKDIR /js
COPY kobra-client /js
RUN pwd && npm install && npm run build

FROM python:alpine

# Set work directory
WORKDIR /app
# Set environment variables
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONUNBUFFERED 1

ENV MQTT_HOST=10.0.2.249

EXPOSE 5000

# Copy project
COPY Pipfile Pipfile.lock main.py /app/
COPY --from=builder /js/dist /app/kobra-client/dist
# Copy ./certs to /app/certs
COPY ./certs /app/certs

VOLUME /app/certs

# Install python dependencies
RUN pip install pipenv && \
    pipenv requirements > ./requirements.txt && \
    pip install -r ./requirements.txt



CMD gunicorn --bind :5000 --worker-class eventlet --workers 1 'main:app'