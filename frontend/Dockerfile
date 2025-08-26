FROM node:20
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
# dev server agar hot reload jalan
EXPOSE 5173
CMD ["npm","run","dev","--","--host","0.0.0.0"]
