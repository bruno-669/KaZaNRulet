import os
os.environ['PADDLE_PDX_CACHE_HOME'] = '/home/sonya/.paddlex'
os.environ['PADDLE_PDX_DISABLE_MODEL_SOURCE_CHECK'] = 'True'

from paddleocr import PaddleOCR
from pdf2image import convert_from_path

ocr = PaddleOCR(lang='ru')
print("OCR готов, обрабатываю PDF...")

images = convert_from_path("test_document.pdf")
print(f"Страниц: {len(images)}")

raw_text = ""

for i, img in enumerate(images):
    img_path = f"/tmp/page_{i}.png"
    img.save(img_path)
    result = ocr.predict(img_path)
    print(f"\n=== Страница {i+1} ===")
    
    # Новый формат результата в PaddleOCR 3.5
    if isinstance(result, dict) and 'rec_texts' in result:
        for txt in result['rec_texts']:
            if txt:
                print(txt)
                raw_text += txt + "\n"
    elif isinstance(result, list):
        for line in result:
            if isinstance(line, dict) and 'rec_texts' in line:
                for txt in line['rec_texts']:
                    if txt:
                        print(txt)
                        raw_text += txt + "\n"
            elif isinstance(line, dict) and 'rec_text' in line:
                txt = line.get('rec_text', '')
                if txt:
                    print(txt)
                    raw_text += txt + "\n"
            elif isinstance(line, list):
                for item in line:
                    if isinstance(item, tuple) and len(item) >= 2:
                        txt = item[1][0] if isinstance(item[1], list) else item[1]
                        if txt:
                            print(txt)
                            raw_text += txt + "\n"

print("\n=== ИТОГО ТЕКСТ ===")
print(raw_text)