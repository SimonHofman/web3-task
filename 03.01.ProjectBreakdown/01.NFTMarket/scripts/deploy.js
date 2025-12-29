const {ethers, upgrades} = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("deployer: ", deployer.address);

    let esVault = await ethers.getContractFactory("EasySwapVault");
    esVault = await upgrades.deployProxy(esVault, {initializer: 'initialize'});
    // 等待部署完成
    if (esVault.deployTransaction) {
        await esVault.deployed();
    } else {
        await esVault.waitForDeployment();
    }
    const esVaultAddress = typeof esVault.getAddress === 'function' ? await esVault.getAddress() : esVault.address;
    console.log("EasySwapVault deployed to: ", esVaultAddress);
    console.log(await upgrades.erc1967.getImplementationAddress(esVaultAddress), " esVault getImplementationAddress");
    console.log(await upgrades.erc1967.getAdminAddress(esVaultAddress), " esVault getAdminAddress")

    newProtocolShare = 200;
    newESVault = esVaultAddress;
    EIP721Name = "EasySwapOrderBook";
    EIP712Version = "1";
    let esDex = await ethers.getContractFactory("EasySwapOrderBook");
    esDex = await upgrades.deployProxy(esDex, [newProtocolShare, newESVault, EIP721Name, EIP712Version], {initializer: 'initialize'});
    if (esDex.deployTransaction) {
        await esDex.deployed();
    } else {
        await esDex.waitForDeployment();
    }
    const esDexAddress = typeof esDex.getAddress === 'function' ? await esDex.getAddress() : esDex.address;
    console.log("EasySwapOrderBook deployed to: ", esDexAddress);
    console.log(await upgrades.erc1967.getImplementationAddress(esDexAddress), " esDex getImplementationAddress");
    console.log(await upgrades.erc1967.getAdminAddress(esDexAddress), " esDex getAdminAddress");

    tx = await esVault.setOrderBook(esDexAddress);
    await  tx.wait();
    console.log("esVault setOrderBook tx: ", tx.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error)
        process.exit(1)
    })