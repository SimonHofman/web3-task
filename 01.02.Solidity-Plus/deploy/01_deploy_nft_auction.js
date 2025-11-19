const {deployments, upgrades, ethers} = require("hardhat")

const fs = require('fs');
const path = require('path');

module.exports = async ({getNamedAccounts, deployments}) => {
    const {save} = deployments;
    const {deployer} = await getNamedAccounts();
    console.log("部署用户地址:", deployer);

    const NftAuction = await ethers.getContractFactory("NftAuction");
    const NftAuctionProxy = await upgrades.deployProxy(
        NftAuction, [], {initializer: "initialize"}
    );
    await NftAuctionProxy.deployed();

    const proxyAddress = await nftAuctionProxy.getAddress();
    console.log("代理合约地址:", proxyAddress);
    const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log("实现合约地址:", implAddress);

    const storePath = path.resolve(_dirname, "./.cache/proxyNftAuction.json");

    // 将部署信息写入存储文件
    // 包含代理地址、实现地址和ABI接口定义
    fs.writeFileSync(
        storePath,
        JSON.stringify({
            proxyAddress,
            implAddress,
            abi: NftAuction.interface.format("json")
        })
    );

    await save("NftAuctionProxy", {
        abi: NftAuction.interface.format("json"),
        address: proxyAddress,
    })
}

module.exports.tags = ["deployNftAuction"];