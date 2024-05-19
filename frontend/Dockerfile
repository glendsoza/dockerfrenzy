FROM node:16.20.2-buster as build

WORKDIR /app

COPY package*.json ./

RUN npm install --force

COPY . .

RUN npm run build --force

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf 

EXPOSE 80

ENTRYPOINT ["nginx", "-g", "daemon off;"]

