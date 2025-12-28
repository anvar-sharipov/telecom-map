# Telecom Map

Backend + Frontend приложение для управления пользователями и авторизацией.

---

## 🚀 Стек технологий

### Backend

- Go (net/http)
- PostgreSQL
- JWT авторизация
- bcrypt (хеширование паролей)

### Frontend

- React
- Vite
- TypeScript
- Fetch API

---

## ℹ️ О проекте

Этот проект — мой первый полноценный backend на Go после перехода с Django.
Цель проекта — изучить архитектуру Go-приложений, работу с JWT, PostgreSQL
и построение чистого backend API с нуля.

## 📁 Структура проекта

├── backend
│ ├── cmd
│ │ └── main.go
│ ├── config
│ │ └── config.go
│ ├── internal
│ │ ├── db
│ │ │ └── postgres.go
│ │ ├── domain
│ │ │ └── user.go
│ │ ├── handler
│ │ │ └── auth.go
│ │ ├── middleware
│ │ ├── repository
│ │ │ ├── postgres
│ │ │ │ └── user_repository.go
│ │ │ └── user_repository.go
│ │ ├── service
│ │ └── utils
│ │ └── jwt.go
│ ├── migrations
│ │ └── 001_create_users.sql
│ ├── pkg
│ ├── tmp
│ │ ├── app.exe
│ │ └── build-errors.log
│ ├── .air.toml
│ ├── .dockerignore
│ ├── .env.example
│ ├── .env.local
│ ├── .env.prod
│ ├── Dockerfile.dev
│ ├── go.mod
│ ├── go.sum
│ ├── main
│ └── netstat
├── frontend
│ ├── node_modules
│ ├── public
│ ├── src
│ │ ├── app
│ │ │ ├── App.tsx
│ │ │ ├── hooks.ts
│ │ │ ├── LanguageSwitcher.tsx
│ │ │ ├── Layout.tsx
│ │ │ ├── router.tsx
│ │ │ └── store.ts
│ │ ├── components
│ │ │ ├── Button.tsx
│ │ │ ├── Input.tsx
│ │ │ └── ThemeToggle.tsx
│ │ ├── features
│ │ │ ├── auth
│ │ │ │ ├── authSlice.ts
│ │ │ │ └── types.ts
│ │ │ └── theme
│ │ │ └── themeSlice.ts
│ │ ├── i18n
│ │ │ ├── ru
│ │ │ │ ├── auth.json
│ │ │ │ ├── common.json
│ │ │ │ └── errors.json
│ │ │ ├── tm
│ │ │ │ ├── auth.json
│ │ │ │ ├── common.json
│ │ │ │ └── errors.json
│ │ │ ├── index.ts
│ │ │ ├── ru.json
│ │ │ └── tm.json
│ │ ├── pages
│ │ │ ├── Home.tsx
│ │ │ ├── Login.tsx
│ │ │ └── Register.tsx
│ │ ├── services
│ │ │ └── api.ts
│ │ ├── styles
│ │ │ └── index.css
│ │ ├── utils
│ │ └── main.tsx
│ ├── .env
│ ├── .gitignore
│ ├── eslint.config.js
│ ├── index.html
│ ├── package-lock.json
│ ├── package.json
│ ├── postcss.config.js
│ ├── README.md
│ ├── tailwind.config.js
│ ├── tsconfig.app.json
│ ├── tsconfig.json
│ ├── tsconfig.node.json
│ └── vite.config.ts
├── .env.example
├── .env.local
├── .env.prod
├── .gitignore
├── docker-compose.yml
├── README.md
