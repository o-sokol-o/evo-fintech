## Тестове завдання

#### Потрібно реалізувати REST API з двома ендпоінтами:
1. Завантаження example.csv файлу, парсинг його і збереження результатів парсингу в базу даних. 
2. Фільтрація і вивантаження попередньо збережених даних в JSON форматі в респонсі.
+ Вимоги до фільтрів:
    + пошук по transaction_id
    + пошук по terminal_id (можливість вказати декілька одночасно id)
    + пошук по status (accepted/declined)
    + пошук по payment_type (cash/card)
    + пошук по date_post по періодам (from/to), наприклад: from 2022-08-12, to 2022-09-01 повинен повернути всі транзакції за вказанийперіод
    + пошук по частково вказаному payment_narrative

#### Технічні вимоги:
- База даних повинна бути реляційна: MySQL, PostgrSQL, тощо
- Документація API повинна бути присутня (Swagger, OpenAPI чи просто в README.md)
- Документація до запуску і використання проекту (в README.md)
- Використання docker та docker-compose
- Можна використовувати будь-які бібліотеки чи фреймворки доступні в опенсорсі.
- Юніт та/або інтеграційні тести
- Зробити третій ендпоінт, котрий буде вивантажувати дані не в JSON, а у вигляді CSV файлу.
- Передбачити можливість завантажувати файл великого розміру (від 1гб) при умові, що ресурс виданий сервісу буде обмежений в 100мб ОЗУ


### Download app
```shell
git clone https://github.com/o-sokol-o/evo-fintech
```

### Build and run app
```shell
docker-compose up -d --build
```

Browse to http://localhost:8080/swagger/index.html. You will see Swagger (only on the local network) Api documents as shown below:

![swagger-image](../main/assets/swagger-image.png)

