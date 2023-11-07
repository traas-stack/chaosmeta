FROM node:16-alpine as builder

WORKDIR /app

COPY . .

RUN yarn install

RUN  yarn build

RUN rm -rf node_modules && \
  NODE_ENV=production yarn install \
  --prefer-offline \
  --pure-lockfile \
  --non-interactive \
  --production=true

FROM nginx:alpine

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder  app/nginx/nginx.conf /etc/nginx/conf.d/nginx.conf

COPY --from=builder  app/dist/ /usr/share/nginx/dist/html/

EXPOSE 8000

CMD ["nginx", "-g", "daemon off;"]