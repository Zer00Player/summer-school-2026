import json
import os

def load_test_cases(file_path):
    """Загружает тест-кейсы из JSON файла"""
    with open(file_path, 'r', encoding='utf-8') as f:
        return json.load(f)

def print_test_case(tc):
    """Красиво выводит один тест-кейс в консоль"""
    print(f"\n{'='*60}")
    print(f"📋 {tc['id']}: {tc['title']}")
    print(f"{'='*60}")
    print(f"Приоритет: {tc['priority']}")
    print(f"Тип: {tc['type']}")
    
    print(f"\n📌 Предусловия:")
    for key, value in tc.get('preconditions', {}).items():
        print(f"  • {key}: {value}")
    
    print(f"\n Шаги:")
    for step in tc.get('steps', []):
        print(f"  {step['step_number']}. {step['action']}")
        print(f"     ✅ Ожидается: {step['expected']}")
    
    print(f"\n🎯 Ожидаемый результат:")
    for key, value in tc.get('expected_result', {}).items():
        print(f"  • {key}: {value}")
    
    print()

def run_all_tests():
    """Запускает все тест-кейсы из папки json-test-cases"""
    
    # 🔧 ИСПРАВЛЕНИЕ ПУТИ:
    # os.path.abspath(__file__) получает полный путь к этому скрипту
    # os.path.dirname(...) берёт папку, где лежит скрипт (03-test/scripts/)
    script_dir = os.path.dirname(os.path.abspath(__file__))
    
    # Поднимаемся на уровень вверх (в 03-test/) и заходим в json-test-cases/
    test_dir = os.path.join(script_dir, '..', 'json-test-cases')
    
    # Нормализуем путь (убирает символы ../)
    test_dir = os.path.normpath(test_dir)
    
    print("🚀 Запуск тест-кейсов...")
    print(f"📁 Ищем файлы в: {test_dir}\n")
    
    # Проверяем, существует ли папка
    if not os.path.exists(test_dir):
        print(f"❌ ОШИБКА: Папка '{test_dir}' не найдена!")
        print("💡 Проверьте, что структура папок соответствует документации.")
        return

    # Проходим по всем JSON файлам в папке
    for filename in sorted(os.listdir(test_dir)):
        if filename.endswith('.json'):
            file_path = os.path.join(test_dir, filename)
            
            print(f"\n{'#'*60}")
            print(f"# ФАЙЛ: {filename}")
            print(f"{'#'*60}")
            
            try:
                data = load_test_cases(file_path)
                print(f"\n📦 Тестовый набор: {data['test_suite']}")
                print(f"📅 Версия: {data['version']}")
                
                for tc in data['test_cases']:
                    print_test_case(tc)
                    
            except json.JSONDecodeError as e:
                print(f"❌ Ошибка формата JSON в файле {filename}: {e}")
            except KeyError as e:
                print(f"❌ Ошибка структуры JSON в файле {filename}: Отсутствует обязательное поле {e}")
            except Exception as e:
                print(f"❌ Неизвестная ошибка при чтении {filename}: {e}")

if __name__ == '__main__':
    run_all_tests()