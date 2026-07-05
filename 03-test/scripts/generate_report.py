import json
import os
from datetime import datetime

def generate_html_report():
    """Генерирует HTML отчёт из JSON тест-кейсов"""
    test_dir = '03-test/json-test-cases'
    
    html = f"""
    <!DOCTYPE html>
    <html lang="ru">
    <head>
        <meta charset="UTF-8">
        <title>Отчёт по тест-кейсам</title>
        <style>
            body {{ font-family: Arial, sans-serif; max-width: 1200px; margin: 0 auto; padding: 20px; }}
            h1 {{ color: #2c3e50; }}
            .test-suite {{ background: #f8f9fa; padding: 20px; border-radius: 8px; margin-bottom: 20px; }}
            .test-case {{ background: white; border: 1px solid #ddd; border-radius: 8px; padding: 15px; margin: 10px 0; }}
            .test-case h3 {{ color: #4ca1af; margin-top: 0; }}
            .priority {{ display: inline-block; padding: 4px 8px; border-radius: 4px; font-size: 12px; font-weight: bold; }}
            .P0 {{ background: #ffebee; color: #c62828; }}
            .P1 {{ background: #fff3e0; color: #e65100; }}
            .P2 {{ background: #e8f5e9; color: #2e7d32; }}
            .step {{ margin: 10px 0; padding-left: 20px; }}
            .expected {{ color: #2e7d32; font-style: italic; }}
            .preconditions {{ background: #f0f0f0; padding: 10px; border-radius: 4px; }}
        </style>
    </head>
    <body>
        <h1>📋 Отчёт по тест-кейсам</h1>
        <p>Сгенерировано: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
    """
    
    for filename in sorted(os.listdir(test_dir)):
        if filename.endswith('.json'):
            with open(os.path.join(test_dir, filename), 'r', encoding='utf-8') as f:
                data = json.load(f)
            
            html += f"""
            <div class="test-suite">
                <h2>📦 {data['test_suite']}</h2>
                <p>Версия: {data['version']}</p>
            """
            
            for tc in data['test_cases']:
                html += f"""
                <div class="test-case">
                    <h3>{tc['id']}: {tc['title']}</h3>
                    <span class="priority {tc['priority']}">{tc['priority']}</span>
                    <span>Тип: {tc['type']}</span>
                    
                    <div class="preconditions">
                        <strong>📌 Предусловия:</strong>
                        <ul>
                """
                
                for key, value in tc.get('preconditions', {}).items():
                    html += f"<li>{key}: {value}</li>"
                
                html += """
                        </ul>
                    </div>
                    
                    <strong>👣 Шаги:</strong>
                    <ol>
                """
                
                for step in tc.get('steps', []):
                    html += f"""
                        <li class="step">
                            {step['action']}
                            <div class="expected">✅ {step['expected']}</div>
                        </li>
                    """
                
                html += """
                    </ol>
                    
                    <strong>🎯 Ожидаемый результат:</strong>
                    <ul>
                """
                
                for key, value in tc.get('expected_result', {}).items():
                    html += f"<li>{key}: {value}</li>"
                
                html += """
                    </ul>
                </div>
                """
            
            html += "</div>"
    
    html += """
    </body>
    </html>
    """
    
    # Сохраняем HTML файл
    output_file = '03-test/test-report.html'
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(html)
    
    print(f"✅ Отчёт создан: {output_file}")
    print(f"🌐 Открой файл в браузере для просмотра")

if __name__ == '__main__':
    generate_html_report()