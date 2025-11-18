const {config} = require("chai");
const {ethers} = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with:", deployer.address);

    const MergeSortedArray1 = await ethers.getContractFactory("MergeSortedArray1");
    const contract = await MergeSortedArray1.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress());

    const intArray1 = [1, 2, 3]
    const intArray2 = [2, 5, 6]
    try {
        const result = await contract.mergeSortedArray(intArray1, intArray2);
        console.log("Result:", result);
    } catch (error) {
        console.error("Error:", error);
    }
}

main()