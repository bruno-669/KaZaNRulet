import os
os.environ['PADDLE_PDX_CACHE_HOME'] = '/home/sonya/.paddlex'
os.environ['PADDLE_PDX_DISABLE_MODEL_SOURCE_CHECK'] = 'True'

from paddleocr import PaddleOCR

ocr = PaddleOCR(lang='ru')
print("PaddleOCR готов!")