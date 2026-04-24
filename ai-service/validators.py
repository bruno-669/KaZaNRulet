import re
from datetime import datetime


def validate_inn(inn: str) -> str | None:
    """Возвращает сообщение об ошибке или None, если всё ок."""
    if not inn:
        return "ИНН не указан"
    if not inn.isdigit():
        return "ИНН должен содержать только цифры"
    if len(inn) not in (10, 12):
        return f"ИНН должен быть 10 (юрлицо) или 12 (ИП) цифр, сейчас {len(inn)}"
    return None


def validate_date(date_str: str) -> str | None:
    if not date_str:
        return "Дата не указана"
    if not re.match(r'^\d{2}\.\d{2}\.\d{4}$', date_str):
        return "Дата должна быть в формате ДД.ММ.ГГГГ"
    try:
        datetime.strptime(date_str, '%d.%m.%Y')
    except ValueError:
        return "Некорректная дата"
    return None


def validate_positive_number(value) -> str | None:
    if value is None:
        return None
    try:
        num = float(value)
        if num <= 0:
            return "Значение должно быть положительным"
    except (ValueError, TypeError):
        return "Должно быть числом"
    return None


def validate_fields(data: dict) -> dict:
    errors = {}

    err = validate_inn(data.get("поставщик_инн"))
    if err:
        errors["поставщик_инн"] = err

    err = validate_inn(data.get("покупатель_инн"))
    if err:
        errors["покупатель_инн"] = err

    err = validate_date(data.get("дата"))
    if err:
        errors["дата"] = err

    err = validate_positive_number(data.get("сумма_итого"))
    if err:
        errors["сумма_итого"] = err

    items = data.get("товары", [])
    if items:
        for i, item in enumerate(items):
            err = validate_positive_number(item.get("количество"))
            if err:
                errors[f"товары[{i}].количество"] = err
            err = validate_positive_number(item.get("цена"))
            if err:
                errors[f"товары[{i}].цена"] = err

    return errors