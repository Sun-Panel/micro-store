# build frontend
FROM node:18.20.5 AS web_image

# 华为源
# RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

# 使用阿里云源
RUN npm config set registry https://registry.npmmirror.com

RUN npm install pnpm -g

WORKDIR /build

COPY . /build

RUN pnpm install 

RUN pnpm run build

# build backend
# 最新alpine3.19导致sqlite3编译失败(https://github.com/mattn/go-sqlite3/issues/1164，
# 临时解决方案:https://github.com/mattn/go-sqlite3/pull/1177)
# sun-panel-auth-server暂时解决方案使用golang:1.21-alpine3.18（因旧版本使用没问题，短期内较稳定） 
FROM golang:1.21-alpine3.18 AS server_image

WORKDIR /build

COPY ./service .

# 中国国内源
RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
    && go env -w GOPROXY=https://goproxy.cn,direct

RUN apk add --no-cache bash curl gcc git musl-dev

RUN go env -w GO111MODULE=on \
    && export PATH=$PATH:/go/bin \
    && go install -a -v github.com/go-bindata/go-bindata/...@latest \
    && go install -a -v github.com/elazarl/go-bindata-assetfs/...@latest \
    && go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... \
    && go build -o sun-panel-auth-server --ldflags="-X sun-panel/global.RUNCODE=release -X sun-panel/global.ISDOCKER=true" main.go
    # && go build -o sun-panel-auth-server --ldflags="-X sun-panel/global.RUNCODE=debug -X sun-panel/global.ISDOCKER=true" main.go



# run_image
FROM alpine

WORKDIR /app

COPY --from=web_image /build/dist /app/web

COPY --from=server_image /build/sun-panel-auth-server /app/sun-panel-auth-server

# 中国国内源
RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories

EXPOSE 3002

RUN apk add --no-cache bash ca-certificates su-exec tzdata \
    && chmod +x ./sun-panel-auth-server

CMD ./sun-panel-auth-server
