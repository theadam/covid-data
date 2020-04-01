FROM node:13.12.0-alpine3.11

WORKDIR /app

RUN npm install --force -g yarn

COPY ./package.json ./yarn.lock /app/
RUN yarn

COPY . /app/

CMD ["npm", "start"]
