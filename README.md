# Описание (в процессе написания)


### Запуск

```make
docker compose -f ./deployment/docker-compose.yml up
```
после этого сервис будер доступен на `localhost:8080`


Сервис также развернут на  
https://cnrprod1725726044-team-79436-32727.avito2024.codenrock.com


### Доступные endpoints (из задания):    

1. `GET /api/ping`          `Проверка доступности`  
2. `GET /api/tenders`          `Получить список опубликованныхх тендеров`   
3. `POST /api/tenders`          `Создать новый тендер`  
4. `GET /api/tenders/my`          `Список пользователя в его зоне ответственности`     
5. `GET /api/tenders/{tenderId}/status`          `Получить статус тендера`   
6. `PUT /api/tenders/{tenderId}/status`          `Изменить статус тендера`   
7. `PATCH /api/tenders/{tenderId}/edit`        `Редактирование тендера`  
8. `PATCH /api/tenders/{tenderId}/edit`          `Редактирование тендера`   
9. `POST /api/bids/new`          `Создание предложения`   
10. `GET /api/bids/my`          `Получение списка ваших предложений`   
11. `GET /api/bids/{tenderId}/list`          `Получение предложений, связанных с указанным тендером`   
12. `GET /api/bids/{bidId}/status`          `Статус предложения`   
13. `PUT /api/bids/{bidId}/status`          `Изменить статус предложения`   
14. `PATCH /api/bids/{bidId}/edit`          `Редактирование предложения`   
15. `PUT /api/bids/{bidId}/submit_decision`          `Отправка решения по предложению`   


### Дополнительные endpoints (для удобства):    
1. `GET /api/users`          `Получить список пользователей`  
2. `GET /api/organizations/all`          `Получить список организаций`  
2. `GET /api/organizations/responsible`          `Получить список организаций и их ответсвенных пользователей `  

### Описание таблиц базы данных:

```sql

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE bids (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    description text,
    tender_id uuid NOT NULL,
    author_type text NOT NULL,
    author_id uuid NOT NULL,
    status text NOT NULL,
    decision text DEFAULT 'Pending'::text NOT NULL,
    version integer DEFAULT 1,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT bids_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES tender(id)
);

CREATE TABLE employee (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    username character varying(50) NOT NULL UNIQUE,
    first_name character varying(50),
    last_name character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name character varying(100) NOT NULL,
    description text,
    type organization_type,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    organization_id uuid,
    user_id uuid,
    CONSTRAINT organization_responsible_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT organization_responsible_user_id_fkey FOREIGN KEY (user_id) REFERENCES employee(id) ON DELETE CASCADE
);

CREATE TABLE tender (
    id uuid DEFAULT uuid_generate_v4() NOT NULL PRIMARY KEY,
    name character varying(100) NOT NULL,
    description text,
    servicetype text NOT NULL,
    organization_id uuid,
    creator_id uuid NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    current_version integer DEFAULT 1,
    status text DEFAULT 'Created'::text NOT NULL,
    CONSTRAINT status_check CHECK ((status = ANY (ARRAY['Created'::text, 'Published'::text, 'Closed'::text]))),
    CONSTRAINT tender_servicetype_check CHECK ((servicetype = ANY (ARRAY['Construction'::text, 'Delivery'::text, 'Manufacture'::text]))),
    CONSTRAINT tender_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES organization(id) ON DELETE CASCADE,
    CONSTRAINT tender_creator_id_fkey FOREIGN KEY (creator_id) REFERENCES employee(id) ON DELETE SET NULL
);

```


