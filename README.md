# Описание (в процессе написания)


### Запуск

```make
docker compose -f ./deployment/docker-compose.yml up
```
после этого сервис будер доступен на `localhost:8080`


Сервис также развернут на  
https://cnrprod1725726044-team-79436-32727.avito2024.codenrock.com


### Доступные endpoins:    

1. `GET /api/ping`          `Проверка доступности`  
2. `GET /api/tenders`          `Получить список опубликованныхх тендеров`   
3. `POST /api/tenders`          `Создать новый тендер`  
4. `GET /api/tenders/my`          `Список пользователя в его зоне ответственности`     
5. `GET /api/tenders/{tenderId}/status`          `Получить статус тендера`   
6. `POST /api/tenders/{tenderId}/status`          `Изменить статус тендера`   
7. `PATCH /api/tenders/{tenderId}/edit`        `Редактирование тендера`  
8. `PATCH /api/tenders/{tenderId}/edit`          `Редактирование тендера`   
9. `POST /api/bids/new`          `Создание предложения`   
10. `GET /api/bids/my`          `Получение списка ваших предложений`   
11. `GET /api/bids/{tenderId}/list`          `Получение предложений, связанных с указанным тендером`   
12. `GET /api/bids/{bidId}/status`          `Статус предложения`   
13. `PUT /api/bids/{bidId}/status`          `Изменить статус предложения`   
14. `PATCH /api/bids/{bidId}/edit`          `Редактирование предложения`   
15. `PUT /api/bids/{bidId}/submit_decision`          `Отправка решения по предложению`   



