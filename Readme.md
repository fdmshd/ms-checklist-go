# Пример микросервиса "чеклист"

Остальные сервисы:
* [Авторизация](https://github.com/fdmshd/ms-user-go)
* [Gateway](https://github.com/fdmshd/ms-gateway)

---
Для запуска использовать:
```
    make up
```

Для остановки контейнеров использовать:
```
    make down
```

Для запуска миграций:
```
    make migrateup
```

Для отката миграций:
```
    make migratedown
```
Для создания миграции:
```
    make migration name=your_migration_name
```
---