const {config} = require("chai");
const {ethers} = require("hardhat")

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with:", deployer.address);

    const BinarySearch1 = await ethers.getContractFactory("BinarySearch1");
    const contract = await BinarySearch1.deploy();
    await contract.waitForDeployment();
    console.log("Contract deployed to:", await contract.getAddress());

    // // Test the binarySearch function
    // const testArray1 = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];
    // const target1 = 7;

    // try {
    //     const result = await contract.binarySearch(testArray1, target1);
    //     console.log("Result:", result);
    // } catch (error) {
    //     console.error("Error:", error);
    // }

    // const testArray2 = [[1, 2, 3, 4, 5, 6, 7, 8, 9, 10], [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]];
    // const target2 = [7,11];

    // try {
    //     const result = await contract.binarySearch(testArray2, target2);
    //     console.log("Result:", result);
    // } catch (error) {
    //     console.error("Error:", error);
    // }

    const testArray = [[1, 2, 3, 4, 5, 6, 7, 8, 9, 10], [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]];
    const target = [7, 11];

    for (let i = 0; i < testArray.length; i++) {
        try {
            const result = await contract.binarySearch(testArray[i], target[i]);
            console.log("Result:", result);
        } catch (error) {
            console.error("Error:", error);
        }
    }
}

main()