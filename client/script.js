// Находим кнопку и место для результата
const button = document.getElementById('myButton');
const resultDiv = document.getElementById('result');

// Добавляем событие "клик"
button.addEventListener('click', () => {
    // Генерируем случайное число (как в твоих JS задачах)
    const randomNum = Math.floor(Math.random() * 100);
    
    // Выводим результат на экран
    resultDiv.textContent = `Случайное число: ${randomNum}`;
    
    console.log('Кнопка нажата! Число:', randomNum);
});