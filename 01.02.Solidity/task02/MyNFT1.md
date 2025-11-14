## 作业2：在测试网上发行一个图文并茂的 NFT

以下是完成这个任务的详细步骤：

### 1. 编写 NFT 合约

使用 OpenZeppelin 的 `ERC721` 库创建一个简单的 NFT 合约：

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract MyNFT is ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor() ERC721("MyNFT", "MNFT") {}

    function mintNFT(address recipient, string memory tokenURI)
        public onlyOwner
        returns (uint256)
    {
        _tokenIds.increment();

        uint256 newItemId = _tokenIds.current();
        _mint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);

        return newItemId;
    }
}
```


**合约说明：**
- 继承自 `ERC721URIStorage` 和 `Ownable`
- 使用 `Counters` 来管理 token ID 的递增
- `mintNFT` 函数只能由合约所有者调用

### 2. 准备图文数据

#### 步骤一：准备图片
- 准备一张图片文件（例如：my-nft-image.png）
- 将其上传到 IPFS 平台，如 Pinata
- 获取图片的 IPFS CID（例如：Qm...）

CID：bafybeichqfoiembbt7er2nsmonf33pnkhdu4khk3fi3b7putsoylk6tmf4

#### 步骤二：创建元数据 JSON 文件
创建一个符合 OpenSea 标准的 JSON 文件：

```json
{
  "name": "My Awesome NFT",
  "description": "This is my first NFT created on Ethereum testnet",
  "image": "ipfs://QmYourImageCIDHere"
}
```


#### 步骤三：上传 JSON 到 IPFS
- 将上面的 JSON 文件保存为 metadata.json
- 上传到 IPFS，获得元数据的 CID（例如：QmMetadataCIDHere）

CID：bafkreifryxc2c7ftj6t7yxs4k55jvaidmb62izb7e2vcbvrpzyha35is2a

最终的 tokenURI 格式为：`ipfs://bafkreifryxc2c7ftj6t7yxs4k55jvaidmb62izb7e2vcbvrpzyha35is2a`

### 3. 部署合约到测试网

#### 在 Remix 中部署：
1. 打开 [Remix IDE](https://remix.ethereum.org/)
2. 创建新文件并粘贴上述合约代码
3. 编译合约（确保编译版本与 pragma 指定的版本一致）
4. 切换环境到 "Injected Web3"
5. 确保 MetaMask 已连接到 Goerli 或 Sepolia 测试网且有足够 ETH
6. 部署合约
7. 记录部署后的合约地址

### 4. 铸造 NFT

在 Remix 中调用 `mintNFT` 函数：
- `recipient`: 填入你的钱包地址
- `tokenURI`: 填入之前生成的 IPFS 元数据链接 (`ipfs://QmMetadataCIDHere`)
- 点击 transact 并在 MetaMask 中确认交易

### 5. 查看 NFT

#### 方法一：OpenSea 测试网
1. 访问 [testnets.opensea.io](https://testnets.opensea.io)
2. 连接你的钱包
3. 搜索你的合约地址或者直接访问你的个人资料页面查看 NFT

#### 方法二：Etherscan 测试网
1. 访问对应测试网的 Etherscan 网站（如 sepolia.etherscan.io）
2. 搜索你的合约地址
3. 在 "Contract" 标签下点击 "Read Contract" 可以查看 token 信息

这样你就成功完成了在测试网上发行一个图文并茂的 NFT 的全部流程。