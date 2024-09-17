# Описание (в процессе написания)


### Запуск

```make
docker compose -f ./deployment/docker-compose.yml up
```
после этого сервис будер доступен на `localhost:8080`


Сервис также развернут на  
https://cnrprod1725726044-team-79436-32727.avito2024.codenrock.com


### Доступные endpoins:    

`GET /api/ping`          `Проверка доступности`  
`GET /api/tenders`          `Получить список опубликованныхх тендеров`   
`POST /api/tenders`          `Создать новый тендер` 
`GET /api/tenders/my`          `Список пользователя в его зоне ответственности`     
`GET /api/tenders/{tenderId}/status`          `Получить статус тендера`   
`POST /api/tenders/{tenderId}/status`          `Изменить статус тендера`   
`PATCH /api/tenders/{tenderId}/edit`        `Редактирование тендера`  
`PATCH /api/tenders/{tenderId}/edit`          `Редактирование тендера`   
`POST /api/bids/new`          `Создание предложения`   
`GET /api/bids/my`          `Получение списка ваших предложений`   
`GET /api/bids/{tenderId}/list`          `Получение предложений, связанных с указанным тендером`   
`GET /api/bids/{bidId}/status`          `Статус предложения`   
`PUT /api/bids/{bidId}/status`          `Изменить статус предложения`   
`PATCH /api/bids/{bidId}/edit`          `Редактирование предложения`   
`PUT /api/bids/{bidId}/submit_decision`          `Отправка решения по предложению`   



