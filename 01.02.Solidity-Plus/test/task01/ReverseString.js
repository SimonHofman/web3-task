const {config} = require("chai");
const {ethers} = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with:", deployer.address);

    const ReverseString = await  ethers.getContractFactory("ReverseString");
    const contract = await ReverseString.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress());

    const str = "asdf";
    try {
        const result  = await contract.reverseString(str);
        console.log("reverse string is : ",  result);
    } catch (error) {
        console.log("Error:", error);
    }
}

main()