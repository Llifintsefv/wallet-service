# wallet-service

Это тестовый проект API, написанный на Go, с использованием Docker и Docker Compose.

## Установка и запуск

1. Склонируйте репозиторий:

   ```sh
   git clone https://github.com/Llifintsefv/wallet-service.git
   ```

2. Файл `config.env` в корневой директории со следующим содержимым:

   ```env
    DB_HOST=db
    DB_PORT=5432
    DB_USER=postgres
    DB_PASSWORD=mysecretpassword
    DB_NAME=postgres-db
    DB_SSL_MODE=disable
    APP_PORT=:8080

   ```

   Вы можете изменить значения по своему усмотрению.

3. Соберите и запустите приложение с помощью Docker Compose:

   ```sh
    docker-compose --env-file config.env up --build 
   ```

4. API будет доступно по адресу: `http://localhost:8080`.

## Дополнительно

- Миграции базы данных автоматически применяются при запуске.
- Добавлен эндпоинт POST api/v1/wallets для создания счета и генерации UUID
