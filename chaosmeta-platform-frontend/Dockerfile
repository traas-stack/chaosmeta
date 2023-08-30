# 设置 Node.js 版本
FROM node:16-alpine

# 设置工作目录
WORKDIR /app

# 将 package.json 和 package-lock.json 复制到工作目录中
COPY . .

RUN yarn install

# 设置容器启动时的命令
CMD ["yarn", "start"]