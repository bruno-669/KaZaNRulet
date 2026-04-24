import os
import tempfile
from contextlib import asynccontextmanager

from fastapi import FastAPI, UploadFile, File
from paddleocr import PaddleOCR
from pdf2image import convert_from_path

from gigachat import extract_fields_from_text
from validators import validate_fields

os.environ['PADDLE_PDX_CACHE_HOME'] = '/home/sonya/.paddlex'
os.environ['PADDLE_PDX_DISABLE_MODEL_SOURCE_CHECK'] = 'True'

ocr = None


@asynccontextmanager
async def lifespan(app: FastAPI):
    global ocr
    ocr = PaddleOCR(lang='ru')
    print("OCR готов")
    yield
    ocr = None


app = FastAPI(lifespan=lifespan)


def extract_text(file_path: str) -> str:
    raw_text = ""

    if file_path.lower().endswith('.pdf'):
        images = convert_from_path(file_path, dpi=300)
        for i, img in enumerate(images):
            img_path = f"/tmp/page_{i}.png"
            img.save(img_path)
            result = ocr.predict(img_path)

            if isinstance(result, list):
                for page in result:
                    if isinstance(page, dict) and 'rec_texts' in page:
                        for t in page['rec_texts']:
                            if t:
                                raw_text += t + "\n"
            elif isinstance(result, dict) and 'rec_texts' in result:
                for t in result['rec_texts']:
                    if t:
                        raw_text += t + "\n"
    else:
        result = ocr.predict(file_path)
        if isinstance(result, list):
            for page in result:
                if isinstance(page, dict) and 'rec_texts' in page:
                    for t in page['rec_texts']:
                        if t:
                            raw_text += t + "\n"
        elif isinstance(result, dict) and 'rec_texts' in result:
            for t in result['rec_texts']:
                if t:
                    raw_text += t + "\n"

    return raw_text.strip()


@app.post("/extract")
async def extract(file: UploadFile = File(...)):
    suffix = ".pdf" if file.filename.lower().endswith('.pdf') else ".png"
    with tempfile.NamedTemporaryFile(delete=False, suffix=suffix) as tmp:
        content = await file.read()
        tmp.write(content)
        tmp_path = tmp.name

    try:
        raw_text = extract_text(tmp_path)

        if not raw_text:
            return {
                "status": "error",
                "error": "Не удалось распознать текст",
                "fields": {},
                "errors": {}
            }

        fields = await extract_fields_from_text(raw_text)
        errors = validate_fields(fields)

        return {
            "status": "ok",
            "fields": fields,
            "errors": errors
        }

    except Exception as e:
        return {
            "status": "error",
            "error": str(e),
            "fields": {},
            "errors": {}
        }
    finally:
        if os.path.exists(tmp_path):
            os.unlink(tmp_path)