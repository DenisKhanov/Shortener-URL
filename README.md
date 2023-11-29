# RU
# Сокращатель Ссылок  
Сокращатель ссылок - это проект, разработанный как часть моей учебной программы Яндекс Практикума. Этот сервис позволяет пользователям сокращать длинные URL-адреса до более удобного формата, облегчая их использование и распространение.

Основные Функции  
Сокращение URL: Пользователи могут преобразовать длинные URL в короткие ссылки, которые легче обменивать и использовать.
Перенаправление по коротким ссылкам: Каждая сокращенная ссылка перенаправляет пользователя на оригинальный URL.
Хранение и Управление Ссылками: Сервис предоставляет интерфейс для управления сокращенными ссылками, включая возможность отметить ссылку как удаленную.
Поддержка Асинхронных Задач: Сервис поддерживает асинхронное удаление ссылок и способен обрабатывать запросы на удаление в фоновом режиме.

Технологии  
Go (Golang): Основной язык программирования проекта.
PostgreSQL: Система управления базами данных для хранения информации о ссылках.
Gin Web Framework: Используется для реализации веб-сервера и маршрутизации.
Docker: Используется для контейнеризации приложения.

Начало Работы  
Чтобы начать использовать сервис, выполните следующие шаги:

Клонирование репозитория: Склонируйте репозиторий проекта на свой локальный компьютер.
Настройка базы данных: Настройте и запустите локальный экземпляр PostgreSQL (Напремер подняв образ в Docker)
Запуск приложения: Запустите сервис сокращения ссылок на вашем компьютере.

Контрибуция  
Этот проект был разработан как часть учебной программы, и мы приветствуем любые предложения и улучшения. Если у вас есть идеи по улучшению проекта, не стесняйтесь отправлять Pull Requests или создавать Issues.


# EN
# URL Shortener  
The URL Shortener is a project developed as part of the Yandex Practicum my training program. This service allows users to shorten long URLs into a more convenient format, making them easier to share and use.

Key Features  
URL Shortening: Users can transform long URLs into short links for easier sharing and usage.
Redirection to Original URLs: Each shortened link redirects the user to the original URL.
Link Management: The service provides an interface for managing shortened links, including the ability to mark a link as deleted.
Asynchronous Task Support: The service supports asynchronous link deletion, capable of handling deletion requests in the background.

Technologies  
Go (Golang): The primary programming language used in the project.
PostgreSQL: Database management system for storing link information.
Gin Web Framework: Used for implementing the web server and routing.
Docker: Used for containerizing the application.

Getting Started  
To start using the service, follow these steps:

Clone the Repository: Clone the project repository to your local machine.
Database Setup: Set up and run a local instance of PostgreSQL.
Launch the Application: Run the URL Shortening service on your machine.
Contribution
This project was developed as part of an educational program, and we welcome any suggestions and improvements. If you have ideas for enhancing the project, feel free to submit Pull Requests or create Issues.
