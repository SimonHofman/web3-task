```shell
同步区块数据
curl -X GET 'http://localhost:8888/api/v1/block/syncByNumber?blockNumber=5671744'
curl -X GET 'http://localhost:8888/api/v1/block/syncByNumber?blockNumber=5671744'
curl -X GET 'http://localhost:8888/api/v1/block/syncByNumber?blockNumber=9625660'
curl -X GET 'http://localhost:8888/api/v1/block/syncByNumber?blockNumber=9625763'

同步最新的区块
curl -X GET 'http://localhost:8888/api/v1/block/syncLatest'

批量同步区块数据
curl -X GET 'http://localhost:8888/api/v1/block/syncBatch' \
  -H 'Content-Type: application/json' \
  -d '{
    "blockNumbers": [9695279, 9695305]
  }'

查询address数据
curl -X GET 'http://localhost:8888/api/v1/address/search?address=0x2cda41645f2dbffb852a605e92b185501801fc28'

查询contract数据
curl -X GET 'http://localhost:8888/api/v1/address/search?address=0x8eb196e77ee0edbe3d75c44a0423438f29f52e9b'
```