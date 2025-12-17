const {ethers, upgrades} = require("hardhat");

async function main() {
    const [signer] = await ethers.getSigners();

    const MetaNodeToken = await ethers.getContractFactory("MetaNodeToken");
    const metaNodeToken = await MetaNodeToken.deploy();
    await metaNodeToken.waitForDeployment();
    const metaNodeTokenAddress = await metaNodeToken.getAddress();
    console.log("MetaNodeToken deployed to:", metaNodeTokenAddress);

    // 1. 获取合约工厂
    const MetaNodeStake = await ethers.getContractFactory("MetaNodeStake");

    // 2. 设置初始化参数
    const startBlock = 1;
    const endBlock = 9999999;
    const metaNodePerBlock = ethers.parseUnits("1", 18);

    // 3. 部署可升级代理合约
    const stake = await upgrades.deployProxy(
        MetaNodeStake,
        [metaNodeTokenAddress, startBlock, endBlock, metaNodePerBlock],
        {initializer: "initialize"}
    );
    await stake.waitForDeployment();

    // todo
    const stakeAddress = await stake.getAddress();
    const tokenAmount = await metaNodeToken.balanceOf(signer.address);
    let tx = await metaNodeToken.connect(signer).transfer(stakeAddress, tokenAmount);
    await tx.wait();
    console.log("MetaNodeToken transfer to MetaNodeStake:", stakeAddress);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });