const { ethers, upgrades } = require("hardhat");
const fs = require("fs");
const path = require("path");

module.exports = async function ({ getNamedAccounts, deployments }) {
    const { save } = deployments;
    const { deployer } = await getNamedAccounts();
    console.log("部署用户地址:", deployer);

    // 读取 .cache/proxyNftAuction.json 文件
    const storePath = path.resolve(__dirname, "./.cache/proxyMyNftAuction.json");
    const storeData = fs.readFileSync(storePath, "utf-8");
    const { proxyAddress, implAddress, abi } = JSON.parse(storeData);

    // 升级版的业务合约
    const MyNftAuctionV2 = await ethers.getContractFactory("MyNftAuctionV2");

    // 升级代理合约
    const myNftAuctionProxyV2 = await upgrades.upgradeProxy(proxyAddress, MyNftAuctionV2, {call: "admin"});
    await myNftAuctionProxyV2.waitForDeployment();
    const proxyAddressV2 = await myNftAuctionProxyV2.getAddress();

    await  save("myNftAuctionProxyV2", {
        abi,
        address: proxyAddressV2,
    });
}

module.exports.tags = ["upgradeMyNftAuction"];