#!/usr/bin/env bash
# ============================================================
# Дамп проекта: структура + содержимое текстовых файлов
# Проверено на macOS (утилита tree опциональна).
# Использование:
#   ./project_dump.sh [выходной_файл]
# ============================================================

set -euo pipefail

# --- Настройки исключений ---
EXCLUDE_DIRS=(".git" "node_modules" "__pycache__" ".venv" "venv" "build" "dist" ".idea" ".vscode" ".DS_Store")
SKIP_EXT=("png" "jpg" "jpeg" "gif" "bmp" "svg" "ico" "woff" "woff2" "ttf" "eot" "mp3" "mp4" "avi" "mov" "pdf" "zip" "gz" "tar" "exe" "dll" "so" "o" "a" "class" "pyc")
SKIP_FILES=("*.min.js" "*.min.css" "package-lock.json" "yarn.lock" "*.map")

OUTPUT_FILE="${1:-project_dump_$(date +%Y%m%d_%H%M%S).txt}"
SCRIPT_NAME="$(basename "$0")"

if [ "$OUTPUT_FILE" = "$SCRIPT_NAME" ]; then
    echo "Ошибка: выходной файл совпадает с именем скрипта. Укажите другое имя." >&2
    exit 1
fi

# --- Дерево каталогов ---
print_tree() {
    if command -v tree &> /dev/null; then
        # tree устанавливается через brew install tree
        local ignore
        ignore=$(printf "%s|" "${EXCLUDE_DIRS[@]}" | sed 's/|$//')
        echo "Project structure (tree):"
        tree -a -I "$ignore" --charset utf-8 --dirsfirst .
    else
        echo "Project structure (find, tree не установлен):"
        find . \( $(printf -- "-name '%s' -prune -o " "${EXCLUDE_DIRS[@]}" | sed 's/ -o $//') \) -print | sed -e 's|[^/]*/| |g'
    fi
}

# --- Проверка: текстовый ли файл ---
is_text_file() {
    local f="$1"
    local ext="${f##*.}"
    for skip in "${SKIP_EXT[@]}"; do
        if [[ "$ext" == "$skip" ]]; then
            return 1
        fi
    done
    # Используем корректную для macOS команду (file -bI)
    local mime
    mime=$(file -bI "$f" 2>/dev/null | cut -d';' -f1 || true)
    [[ "$mime" == text/* ]] && return 0
    return 1
}

# --- Основной блок ---
{
    echo "=============== PROJECT DUMP ==============="
    echo "Generated: $(date)"
    echo "Root: $(pwd)"
    echo "============================================"
    echo ""
    print_tree
    echo ""
    echo "=============== FILE CONTENTS ==============="
    echo ""

    # Формируем аргументы для find – исключаем каталоги через -not -path
    find_args=()
    for dir in "${EXCLUDE_DIRS[@]}"; do
        find_args+=(-not -path "*/${dir}/*")
    done

    # Рекурсивно проходим по всем файлам
    while IFS= read -r -d '' file; do
        # Пропускаем сам выходной файл
        if [[ "$file" == "./$OUTPUT_FILE" ]]; then
            continue
        fi

        # Пропускаем по маскам имён (например, *.min.js)
        base=$(basename "$file")
        local_skip=0
        for pattern in "${SKIP_FILES[@]}"; do
            if [[ "$base" == $pattern ]]; then
                local_skip=1
                break
            fi
        done
        [[ $local_skip -eq 1 ]] && continue

        if is_text_file "$file"; then
            echo "--- FILE: $file ---"
            cat "$file"
            echo ""
        else
            echo "--- SKIPPED BINARY: $file ---"
        fi
    done < <(find . "${find_args[@]}" -type f -print0)

    echo ""
    echo "=============== END OF DUMP ==============="
} > "$OUTPUT_FILE"

echo "Дамп сохранён в: $OUTPUT_FILE"