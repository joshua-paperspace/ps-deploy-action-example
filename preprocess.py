import torch
import torchvision.transforms as transforms
import cv2
import numpy as np
from PIL import Image
from io import BytesIO

def imgToTensor(image):

    pil_image = Image.open(BytesIO(image))
    img = np.array(pil_image)
    img_stretch = cv2.resize(img, (224, 224))
    transform = transforms.ToTensor()
    tensor = transform(img_stretch)
    tensor_shaped = torch.reshape(tensor, [1,3,224,224])

    return tensor_shaped