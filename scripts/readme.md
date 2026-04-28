# 构建docker镜像

## 主平台

```bash
# 进入项目根目录执行以下命令
./scripts/build.sh -n microapp-store -v 1.0.0 --push \
-r registry.cn-hangzhou.aliyuncs.com/hslr \
-p linux/amd64,linux/arm64
```
