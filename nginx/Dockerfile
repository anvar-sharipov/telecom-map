# 1️⃣ Берём официальный nginx образ (стабильный, прод)
FROM nginx:1.25-alpine

# 2️⃣ Удаляем дефолтный конфиг nginx
# Он нам не нужен, чтобы не конфликтовал
RUN rm /etc/nginx/conf.d/default.conf

# 3️⃣ Копируем наш nginx.conf
# Он станет основным конфигом
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 4️⃣ Копируем SSL сертификаты внутрь контейнера
COPY ssl /etc/nginx/ssl

# 5️⃣ Открываем порты
EXPOSE 80
EXPOSE 443
