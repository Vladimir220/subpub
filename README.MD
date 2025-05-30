# Библиотека SubPub + Сервис подписок #
---
### Содержание ###
- [Как работает библиотека SubPub](#как-работает-библиотека-subpub)
- [Запуск сервера подписок](#запуск-сервера-подписок)
- [Особенности](#особенности)
- [TODO](#todo)
- [Пример клиентского приложения](#пример-клиентского-приложения)
---

### Как работает библиотека SubPub ###

В основе моей самодельной шины:
- Слайсы сообщений для каждого канала -- map[string][]any
- Горутины активных подписок, которые:
    - Ожидают Broadcast выбранного канала при поступлении новых сообщений по подписке -- map[string]*sync.Cond
    - Хранят индекс последнего считанного сообщения для гарантии порядка чтения
    - Завершаются при отмене подписки (генерируются сигнал отмены подписки и служебный Broadcast выбранного канала)

Библиотеку можете подключить к своему проекту, выполнив команду: 
```
go get github.com/Vladimir220/subpub/subpub_lib 
```
Для запуска тестов выполните слшедующую команду из папки ``` subpub_lib ```:
```
go test .
```

### Запуск сервера подписок ###
- Перейдите в папку subpub_service: ``` cd subpub_service ```
- Задайте порт, который будет прослушивать сервер из контейнера:
    - В ```docker-compose.yml``` измените следующую секцию:
        ```
        ports:
        - "1234:1234"
        ```
    - В ``` ./server/.env ``` измените переменную ``` SERVER_URL=":1234" ```
- Выполните команду ``` docker-compose up -d ``` чтобы фоново запустить сервер

### Особенности ###
- Логи будут сохраняться у вас на хосте в папке ``` ./server/Log/ ```
- При изменении кода на хосте нет необходимости пересобирать контейнер
- ``` graceful shutdown ``` обеспечивается использованием ``` defer ``` в сочетании с методами ``` Cancel ``` различных объектов

### TODO ###
- Авторизация/аутентификация (возможно, для внутреннего использования сервиса это и не нужно)
- Агрегация логов  
- Внешняя конфигурация
- По достижению устойчивой версии заменить маппинг томов на полное копирование файлов в контейнер

### Пример клиентского приложения ###
Простой пример клиентского приложения можете найти в ``` ./client/ ```

Перед запуском в ``` ./client/.env ``` измените переменную ``` SERVER_URL=":1234" ```, чтобы адрес соответствовал адресу запущенного сервера

Для запуска выполните следующие команды на своём хосте (из папки ``` subpub_service ```):
```
go mod tidy
go run ./client/.
```


