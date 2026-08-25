基于 Go 实现的公交充电场站调度平台项目，一款后端服务，完成公交场站的桩位分配、插枪鉴权、功率限制与充电计量调度。

## 运行

```bash
go build -mod=vendor -o buscharge ./cmd/buscharge
./buscharge -addr :8080 -data ./buscharge-data -web ./web
```

控制台页面：`/buses`、`/piles`、`/power`、`/alarms`。
