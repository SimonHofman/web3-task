const {config} = require("chai");
const {ethers} = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    console.log("Deploying with:", deployer.address);

    const RomanToInt1 = await ethers.getContractFactory("RomanToInt1");
    const contract = await RomanToInt1.deploy();
    await contract.waitForDeployment();
    console.log("Contract is deployed to:", await contract.getAddress())

    const strArray =["MMMDCCXLIX", "LVIII", "MCMXCIV", "III", "IV", "IX"];
    const intArray = [3749n, 58n, 1994n, 3n, 4n, 9n];
    for (let i = 0; i < strArray.length; i++) {
        try {
            const result = await contract.romanToInt(strArray[i]);
            console.log("Result:", result);
            console.log("Expected:", intArray[i]);
            console.log("Is correct:", result === intArray[i]);
        } catch (error) {
            console.error("Error:", error);
        }
    }
}

main()