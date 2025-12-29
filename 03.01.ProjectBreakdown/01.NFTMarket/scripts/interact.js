const {ethers, upgrades} = require("hardhat")
const {Side, SaleKind} = require("../test/common")
const {toBn} = require("evm-bn");

const esDex_name = "EasySwapOrderBook";
let esDex_address

const esVault_name = "EasySwapVault";
let esVault_address

const erc721_name = "TestERC721"
let erc721_address

let esDex, esVault, testERC721
let deployer

async function main() {
    [deployer, trader] = await ethers.getSigners();
    console.log("deployer:", deployer.address);
    console.log("trader:", trader.address);

    esVault = await ethers.getContractFactory(esVault_name);
    esVault = await upgrades.deployProxy(esVault, {initializer: 'initialize'});
    // 等待部署完成
    if (esVault.deployTransaction) {
        await esVault.deployed();
    } else {
        await esVault.waitForDeployment();
    }
    esVault_address = typeof esVault.getAddress === 'function' ? await esVault.getAddress() : esVault.address;
    console.log("EasySwapVault deployed to: ", esVault_address);
    console.log(await upgrades.erc1967.getImplementationAddress(esVault_address), " esVault getImplementationAddress");
    console.log(await upgrades.erc1967.getAdminAddress(esVault_address), " esVault getAdminAddress")

    newProtocolShare = 200;
    newESVault = esVault_address;
    EIP721Name = esDex_name;
    EIP712Version = "1";
    esDex = await ethers.getContractFactory(esDex_name);
    esDex = await upgrades.deployProxy(esDex, [newProtocolShare, newESVault, EIP721Name, EIP712Version], {initializer: 'initialize'});
    if (esDex.deployTransaction) {
        await esDex.deployed();
    } else {
        await esDex.waitForDeployment();
    }
    esDex_address = typeof esDex.getAddress === 'function' ? await esDex.getAddress() : esDex.address;
    console.log("EasySwapOrderBook deployed to: ", esDex_address);
    console.log(await upgrades.erc1967.getImplementationAddress(esDex_address), " esDex getImplementationAddress");
    console.log(await upgrades.erc1967.getAdminAddress(esDex_address), " esDex getAdminAddress");

    const TestERC721 = await ethers.getContractFactory(erc721_name);
    testERC721 = await TestERC721.deploy();
    await testERC721.waitForDeployment();
    const erc721_address = await testERC721.getAddress();
    console.log("ERC721Address::", erc721_address);

    await approvalForVault();

    await testMakeOrder();
}

async function approvalForVault() {
    let isApproved = await testERC721.isApprovedForAll(deployer.address, esVault_address);

    if (!isApproved) {
        console.log("Already approved");
        return;
    }

    let tx = await testERC721.setApprovalForAll(esVault_address, true);
    await tx.wait();
    console.log("Approval tx:", tx.hash);
}

async function testMakeOrder(tokenId = 0) {
    let now = parseInt(new Date() / 1000) + 100000;
    let salt = 1;
    let nftAddress = erc721_address;
    let order = {
        side: Side.List,
        saleKind: SaleKind.FixedPriceForItem,
        maker: deployer.address,
        nft: [tokenId, nftAddress, 1],
        price: toBn("0.002"),
        expiry: now,
        salt: salt,
    }

    tx = await esDex.makeOrders([order]);
    txRec = await tx.wait();
    console.log(tx.hash);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    })