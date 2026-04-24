import httpx
import json
import re
import time

GIGACHAT_AUTH_URL = "https://ngw.devices.sberbank.ru:9443/api/v2/oauth"
GIGACHAT_API_URL = "https://gigachat.devices.sberbank.ru/api/v1/chat/completions"

# Это Authorization Key (из вашего скриншота)
AUTHORIZATION_KEY = "MDE5YWU4ZWYtMDFhYi03YjcyLTgwZDItMDZiYjkyMDFjZGE0OmQ5OTkzNzljLWY0ZGItNDliOC04MDU5LWY4ZmZmNjA5NzY0Yw=="

class GigaChatClient:
    def __init__(self):
        self.access_token = None
        self.token_expires_at = 0
    
    async def get_access_token(self) -> str:
        """Получает или обновляет Access Token"""
        # Если токен ещё действителен (осталось больше 5 минут), используем его
        if self.access_token and time.time() * 1000 < self.token_expires_at - 300000:
            return self.access_token
        
        # Запрос на получение токена
        headers = {
            "Authorization": f"Basic {AUTHORIZATION_KEY}",
            "Content-Type": "application/x-www-form-urlencoded",
            "RqUID": str(__import__('uuid').uuid4())
        }
        
        data = {
            "scope": "GIGACHAT_API_PERS"  # Или GIGACHAT_API_B2B / GIGACHAT_API_CORP
        }
        
        async with httpx.AsyncClient() as client:
            response = await client.post(
                GIGACHAT_AUTH_URL,
                headers=headers,
                data=data,
                timeout=30.0
            )
            response.raise_for_status()
            
            token_data = response.json()
            self.access_token = token_data["access_token"]
            self.token_expires_at = token_data["expires_at"]
            
            return self.access_token


async def extract_fields_from_text(raw_text: str) -> dict:
    client = GigaChatClient()
    
    # Получаем Access Token
    token = await client.get_access_token()
    
    prompt = f"""
Ты бухгалтерский ассистент. Из текста накладной извлеки поля строго в JSON.

Верни только JSON, без комментариев. Ключи на русском.

Извлеки:
- номер_документа
- дата (ДД.ММ.ГГГГ)
- поставщик_название
- поставщик_инн
- покупатель_название
- покупатель_инн
- товары (массив объектов: название, количество, цена)
- сумма_итого (число)

Если поле не найдено — ставь null.
Товары — извлеки ВСЕ строки таблицы.

Текст документа:
{raw_text}
"""

    payload = {
        "model": "GigaChat",
        "messages": [{"role": "user", "content": prompt}],
        "temperature": 0.1
    }

    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"  # Используем Access Token
    }

    async with httpx.AsyncClient() as http_client:
        for attempt in range(2):
            try:
                response = await http_client.post(
                    GIGACHAT_API_URL, json=payload, headers=headers, timeout=60.0
                )
                response.raise_for_status()
                data = response.json()
                content = data["choices"][0]["message"]["content"]
                
                content = re.sub(r'```(?:json)?', '', content).strip()
                
                return json.loads(content)
            except (json.JSONDecodeError, KeyError) as e:
                if attempt == 1:
                    return {}
                continue
            except Exception as e:
                return {}
    
    return {}