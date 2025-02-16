# 🚀 Event Registration API

Этот API позволяет управлять пользователями и событиями: регистрировать пользователей, создавать события, регистрировать участников, а также управлять аккаунтами.

## 📌 1️⃣ Регистрация пользователя
**URL:** `/api/register`  
**Method:** `POST`  
**Body:**
```json
{
    "email": "testuser@example.com",
    "password": "testpassword"
}
```

---

## 📌 2️⃣ Авторизация (получение `accessToken`)
**URL:** `/api/login`  
**Method:** `POST`  
**Body:**
```json
{
    "email": "testuser@example.com",
    "password": "testpassword"
}
```

---

## 📌 3️⃣ Создание события
**URL:** `/events/create`  
**Method:** `POST`  
**Body:**
```json
{
    "user_id": "1",
    "title": "Tech Meetup",
    "date": "2025-05-10",
    "time": "18:00:00",
    "venue": "New York Conference Center",
    "description": "Annual meetup for developers and entrepreneurs.",
    "note": "Bring business cards.",
    "price": "Free",
    "image_url": "https://example.com/event.jpg",
    "attendees": [],
    "is_active": true
}
```

---

## 📌 4️⃣ Получение списка всех событий
**URL:** `/events/list`  
**Method:** `GET`  
**Body:** _нет_

---

## 📌 5️⃣ Получение события по ID
**URL:** `/events/get?id=1`  
**Method:** `GET`  
**Body:** _нет_

---

## 📌 6️⃣ Регистрация пользователя на событие
**URL:** `/events/register`  
**Method:** `POST`  
**Body:**
```json
{
    "event_id": "1",
    "email": "testuser@example.com"
}
```

---

## 📌 7️⃣ Обновление события
**URL:** `/events/update?id=1`  
**Method:** `PUT`  
**Body:**
```json
{
    "title": "Tech Conference 2025",
    "date": "2025-06-20",
    "time": "10:00:00",
    "venue": "San Francisco Tech Center",
    "description": "An updated event with more tech talks.",
    "note": "Early registration recommended.",
    "price": "$50",
    "image_url": "https://example.com/conference.jpg",
    "attendees": ["testuser@example.com"],
    "is_active": true
}
```

---

## 📌 8️⃣ Удаление события
**URL:** `/events/delete?id=1`  
**Method:** `DELETE`  
**Body:** _нет_

---

## 📌 9️⃣ Обновление пользователя
**URL:** `/admin/update?id=1`  
**Method:** `PUT`  
**Body:**
```json
{
    "new_email": "updateduser@example.com",
    "new_password": "newsecurepassword"
}
```

---

## 📌 🔟 Удаление пользователя
**URL:** `/admin/delete?id=1`  
**Method:** `DELETE`  
**Body:** _нет_
