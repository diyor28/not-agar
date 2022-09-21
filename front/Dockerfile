FROM node:12.18.2-alpine3.9 as build
WORKDIR /app
COPY package.json ./
COPY yarn.lock ./

RUN yarn install --pure-lockfile --non-interactive

COPY . .

RUN NODE_OPTIONS=--max_old_space_size=4096 yarn build

FROM nginx:1.21.3-alpine as production
COPY nginx/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/build /usr/share/nginx/html
CMD ["nginx", "-g", "daemon off;"]
