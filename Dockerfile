FROM python:3.9-alpine

WORKDIR /app

RUN pip install flask requests

COPY . .

ENV PORT=8880

EXPOSE ${PORT}

CMD ["python", "app.py"]