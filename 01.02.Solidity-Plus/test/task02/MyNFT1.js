const {config} = require("chai");
const hre = require("hardhat");

async function main() {
    const {deployer} = await hre.getNamedAccounts()
    console.log("deployer account:", deployer);

    const MyNFT1 = await hre.ethers.getContractFactory("MyNFT1");
    const contract = await MyNFT1.deploy(deployer);
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress());

    const tx = await contract.mintNFT(deployer, `ipfs://bafkreifryxc2c7ftj6t7yxs4k55jvaidmb62izb7e2vcbvrpzyha35is2a`)
    await tx.wait();
    console.log("Mint transaction hash:", tx.hash)
}

main()