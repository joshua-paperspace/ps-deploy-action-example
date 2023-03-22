FROM paperspace/fastapi-deployment:latest

WORKDIR /app

COPY main.py preprocess.py resnet.py requirements.txt ./
COPY config ./config
COPY models ./models

RUN pip3 install -U pip && pip3 install -r requirements.txt

ENV MODEL_DIR = models
ENV MODEL_FILE = resnet50.pt
ENV MODEL_NAME = resnet50

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "80"]