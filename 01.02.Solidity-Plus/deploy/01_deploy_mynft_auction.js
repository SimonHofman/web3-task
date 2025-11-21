const {deployments, upgrades, ethers} = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async ({getNamedAccounts, deployments}) => {
    const {save} = deployments;
    const {deployer} = await getNamedAccounts();
    console.log("部署用户地址:", deployer);

    const MyNftAuction = await ethers.getContractFactory("MyNftAuction");
    const myNftAuctionProxy = await upgrades.deployProxy(
        MyNftAuction,
        [],
        {initializer: "initialize"}
    );
    await myNftAuctionProxy.waitForDeployment();

    const proxyAddress = await myNftAuctionProxy.getAddress();
    console.log("代理合约地址:", proxyAddress);
    const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);
    console.log("实现合约地址:", implAddress);

    const storePath = path.resolve(__dirname, "./.cache/proxyMyNftAuction.json");

    // 将部署信息写入存储文件
    // 包含代理地址、实现地址和ABI接口定义
    fs.writeFileSync(
        storePath,
        JSON.stringify({
            proxyAddress,
            implAddress,
            abi: MyNftAuction.interface.format("json"),
        })
    );

    await save("MyNftAuctionProxy", {
        abi: MyNftAuction.interface.format("json"),
        address: proxyAddress,
    })
};

module.exports.tags = ["deployMyNftAuction"];